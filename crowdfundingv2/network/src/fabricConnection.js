const { Gateway, Wallets } = require('fabric-network');
const fs = require('fs');
const path = require('path');
const config = require('./config');
const logger = require('./utils/logger');

/**
 * FabricConnection - Manages connections to Fabric network
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

            // Create and connect gateway
            const gateway = new Gateway();
            await gateway.connect(connectionProfile, {
                wallet,
                identity: org.adminUser,
                // Enable discovery - required for finding peers
                // asLocalhost: false is critical for microfab's external URLs
                discovery: { enabled: true, asLocalhost: false },
                eventHandlerOptions: {
                    commitTimeout: 300,
                    endorseTimeout: 300
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
     * Submit transaction with PDC-aware endorsement
     * For PDC transactions, we need to specify which orgs should endorse
     * to prevent non-member orgs from being selected by discovery
     */
    async submitTransaction(orgKey, contractName, functionName, ...args) {
        const contract = await this.getContract(orgKey);
        const fcn = `${contractName}:${functionName}`;

        logger.info(`📤 Submit: ${fcn} | Args: ${JSON.stringify(args).substring(0, 100)}...`);

        try {
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
            logger.error(`❌ Submit failed: ${error.message}`);
            throw error;
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
