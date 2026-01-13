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
            logger.info(`🔗 Using cached connection for org: ${orgKey}`);
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
                // Enable discovery with asLocalhost for .nip.io URL translation
                // setEndorsingOrganizations in submitTransaction limits which orgs endorse
                discovery: { enabled: true, asLocalhost: true },
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

                // Explicit MSP ID mapping for each org to ensure correct endorsement
                const orgMspMap = {
                    'startup': 'StartupOrgMSP',
                    'investor': 'InvestorOrgMSP',
                    'validator': 'ValidatorOrgMSP',
                    'platform': 'PlatformOrgMSP'
                };
                const mspId = orgMspMap[orgKey] || config.orgs[orgKey]?.mspId;

                logger.info(`📤 Submit: ${fcn} | Org: ${orgKey} | MSP: ${mspId} | Args: ${JSON.stringify(args).substring(0, 100)}... (attempt ${attempt})`);

                // Use 
                //  to explicitly target the correct peer
                // This is necessary because Microfab discovery ignores setEndorsingOrganizations
                try {
                    // Get network instance (must await if getting fresh)
                    const gateway = this.gateways[orgKey];
                    const network = await gateway.getNetwork(config.channelName);

                    // Get channel and endorsers
                    const channel = network.getChannel();
                    const allEndorsers = channel.getEndorsers();

                    // Filter for endorsers belonging to the target MSP
                    let orgEndorsers = allEndorsers.filter(peer => peer.mspid === mspId);

                    // If no endorsers found via discovery, try to get from connection profile explicit peers
                    if (orgEndorsers.length === 0) {
                        logger.warn(`⚠️ Discovery returned 0 endorsers for ${mspId}. Probing connection profile...`);
                        try {
                            const gateway = this.gateways[orgKey];
                            // Access the internal connection options to get the profile
                            // Note: This relies on internal SDK structure or we need to reload the profile
                            // Safer way: Construct the expected Microfab peer name since we know the pattern
                            // or read from the loaded profile if we saved it. 

                            // Microfab pattern is usually: {orgNameLower}peer-api.127-0-0-1.nip.io:9090
                            // But better to check the profile we loaded
                            const orgConfig = config.orgs[orgKey];
                            const profilePath = require('path').resolve(config.gatewaysDir, orgConfig.gatewayFile);
                            const profile = require(profilePath);

                            if (profile.organizations && profile.organizations[orgConfig.name] && profile.organizations[orgConfig.name].peers) {
                                const peerNames = profile.organizations[orgConfig.name].peers;
                                logger.info(`🔍 Found explicit peers in profile: ${peerNames.join(', ')}`);

                                for (const peerName of peerNames) {
                                    try {
                                        const peer = channel.getEndorser(peerName);
                                        // let peer = channel.getEndorser(peerName);
                                        // // If channel doesn't know it (discovery failed), try getting it from the client directly
                                        // if (!peer && channel.client) {
                                        //     peer = channel.client.getEndorser(peerName, mspId);
                                        // }
                                        if (peer) {
                                            orgEndorsers.push(peer);
                                            logger.info(`✅ Successfully added peer by name: ${peerName}`);
                                        }
                                    } catch (e) {
                                        logger.warn(`Failed to get endorser ${peerName}: ${e.message}`);
                                    }
                                }
                            }
                        } catch (err) {
                            logger.warn(`❌ Failed to parse connection profile for backup peers: ${err.message}`);
                        }
                    }

                    if (orgEndorsers.length > 0) {
                        logger.info(`🎯 Found ${orgEndorsers.length} peer(s) for ${mspId}: ${orgEndorsers.map(p => p.name).join(', ')}`);
                        // Set explicit endorsing peers
                        const transaction = contract.createTransaction(fcn);
                        transaction.setEndorsingPeers(orgEndorsers);

                        const result = await transaction.submit(...args);
                        const resultStr = result.toString();
                        logger.info(`📥 Result: ${resultStr.substring(0, 100)}...`);

                        if (!resultStr || resultStr === '') return { success: true };
                        try { return JSON.parse(resultStr); } catch { return { success: true, message: resultStr }; }
                    } else {
                        logger.warn(`⚠️ No peers found for ${mspId}, falling back to setEndorsingOrganizations`);
                    }
                } catch (pe) {
                    logger.warn(`⚠️ Peer lookup failed: ${pe.message}, falling back`);
                }

                // Fallback (or if peer lookup failed)
                const transaction = contract.createTransaction(fcn);
                transaction.setEndorsingOrganizations(mspId);
                logger.info(`🎯 Endorsing organizations set to: ${mspId}`);

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
     * Parse transaction result
     */
    _parseResult(result) {
        const resultStr = result.toString();
        if (!resultStr || resultStr === '') {
            return { success: true };
        }
        try {
            return JSON.parse(resultStr);
        } catch {
            return { success: true, message: resultStr };
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
