const { Gateway, Wallets } = require('fabric-network');
const fs = require('fs');
const path = require('path');
const config = require('./config');
const logger = require('./utils/logger');

/**
 * FabricConnection - Manages connections to Fabric network
 * To run in production mode: NODE_ENV=production npm run dev
 * To run in development mode: npm run dev (uses NONE strategy by default)
 */
class FabricConnection {
    constructor() {
        this.gateways = {};
        this.contracts = {};
    }

    /**
     * Connect to an organization's network
     */
    async connect(orgKey) {
        if (this.gateways[orgKey]) {
            return this.contracts[orgKey];
        }

        const org = config.orgs[orgKey];
        if (!org) {
            throw new Error(`Unknown organization: ${orgKey}`);
        }

        try {
            // Load wallet
            const walletPath = path.join(config.walletsDir, org.name);
            logger.info(`📂 Loading wallet from: ${walletPath}`);
            const wallet = await Wallets.newFileSystemWallet(walletPath);

            // Check identity exists
            const identity = await wallet.get(org.adminUser);
            if (!identity) {
                throw new Error(`Identity ${org.adminUser} not found. Run 'npm run enroll' first.`);
            }
            logger.info(`✅ Found identity: ${org.adminUser}`);

            // Load connection profile
            const gatewayPath = path.join(config.gatewaysDir, org.gatewayFile);
            logger.info(`📂 Loading gateway from: ${gatewayPath}`);
            const connectionProfile = JSON.parse(fs.readFileSync(gatewayPath, 'utf8'));

            // Environment-aware configuration
            const isDevelopment = process.env.NODE_ENV !== 'production';
            const eventStrategies = require('fabric-network').DefaultEventHandlerStrategies;

            // For production: use MSPID_SCOPE_ANYFORTX (wait for commit confirmation)
            // For development (Microfab): use NONE (Microfab doesn't expose event listeners properly)
            const eventStrategy = isDevelopment
                ? eventStrategies.NONE
                : eventStrategies.MSPID_SCOPE_ANYFORTX;

            if (isDevelopment) {
                logger.info(`⚠️  Using NONE event strategy (dev mode) - commits still happen`);
            }

            // Create and connect gateway
            const gateway = new Gateway();

            await gateway.connect(connectionProfile, {
                wallet,
                identity: org.adminUser,
                // Discovery enabled - required for Microfab URL translation (asLocalhost)
                // In production, asLocalhost should be false
                discovery: { enabled: true, asLocalhost: isDevelopment },
                eventHandlerOptions: {
                    commitTimeout: isDevelopment ? 30 : 300,
                    endorseTimeout: isDevelopment ? 30 : 300,
                    strategy: eventStrategy
                }
            });

            // Get network and contract
            const network = await gateway.getNetwork(config.channelName);
            const contract = network.getContract(config.chaincodeName);

            this.gateways[orgKey] = gateway;
            this.contracts[orgKey] = contract;

            logger.info(`✅ Connected to ${org.name}`);
            return contract;
        } catch (error) {
            logger.error(`❌ Failed to connect to ${org.name}: ${error.message}`);
            throw error;
        }
    }

    /**
     * Get contract for organization
     */
    async getContract(orgKey) {
        return this.connect(orgKey);
    }

    /**
     * Disconnect and clear cached connections for an org (useful when connections become stale)
     */
    async disconnect(orgKey) {
        if (this.gateways[orgKey]) {
            try {
                this.gateways[orgKey].disconnect();
            } catch (e) {
                // Ignore disconnect errors
            }
            delete this.gateways[orgKey];
            delete this.contracts[orgKey];
            logger.info(`🔌 Disconnected and cleared cache for ${orgKey}`);
        }
    }

    /**
     * Submit transaction with PDC-aware endorsement
     * For PDC transactions, we need to specify which orgs should endorse
     * to prevent non-member orgs from being selected by discovery
     * Includes auto-retry with reconnection for stale connection errors
     */
    async submitTransaction(orgKey, contractName, functionName, ...args) {
        const fcn = `${contractName}:${functionName}`;
        const maxRetries = 2;

        for (let attempt = 1; attempt <= maxRetries; attempt++) {
            try {
                const contract = await this.getContract(orgKey);
                logger.info(`📤 Submit: ${fcn} | Args: ${JSON.stringify(args).substring(0, 100)}... (attempt ${attempt})`);

                // Get the org's MSP ID for endorsing
                // For PDC transactions, we ONLY use the calling org's peer
                // Cross-org endorsement causes mismatch because PDC data isn't synced immediately
                const org = config.orgs[orgKey];
                const endorsingOrgs = [org.mspId];

                // Create transaction with explicit endorsing orgs for PDC support
                const transaction = contract.createTransaction(fcn);
                transaction.setEndorsingOrganizations(...endorsingOrgs);

                const result = await transaction.submit(...args);
                const resultStr = result.toString();

                logger.info(`📥 Result: ${resultStr.substring(0, 100)}...`);

                if (!resultStr || resultStr === '') {
                    return { success: true };
                }

                try {
                    return JSON.parse(resultStr);
                } catch {
                    return { success: true, message: resultStr };
                }
            } catch (error) {
                logger.error(`❌ Submit failed (attempt ${attempt}): ${error.message}`);

                // If "No peers for strategy" error, disconnect and retry
                if (error.message.includes('No peers for strategy') && attempt < maxRetries) {
                    logger.info(`🔄 Reconnecting ${orgKey} due to stale connection...`);
                    await this.disconnect(orgKey);
                    continue;
                }

                throw error;
            }
        }
    }

    /**
     * Evaluate transaction (query)
     */
    async evaluateTransaction(orgKey, contractName, functionName, ...args) {
        const contract = await this.getContract(orgKey);
        const fcn = `${contractName}:${functionName}`;

        logger.info(`🔍 Query: ${fcn} | Args: ${JSON.stringify(args).substring(0, 100)}...`);

        try {
            const result = await contract.evaluateTransaction(fcn, ...args);
            const resultStr = result.toString();

            logger.info(`📥 Result: ${resultStr.substring(0, 100)}...`);

            if (!resultStr || resultStr === '') {
                return [];
            }

            try {
                return JSON.parse(resultStr);
            } catch {
                return resultStr;
            }
        } catch (error) {
            logger.error(`❌ Query failed: ${error.message}`);
            throw error;
        }
    }

    /**
     * Disconnect all gateways
     */
    async disconnect() {
        for (const [orgKey, gateway] of Object.entries(this.gateways)) {
            gateway.disconnect();
            logger.info(`🔌 Disconnected from ${orgKey}`);
        }
        this.gateways = {};
        this.contracts = {};
    }
}

// Export singleton
module.exports = new FabricConnection();
