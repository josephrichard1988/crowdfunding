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

            // CRITICAL: Use NETWORK_SCOPE_ANYFORTX to wait for commit confirmation
            // This waits for ANY peer in the network to confirm commit, which works
            // reliably in Microfab where discovery doesn't properly expose MSP-specific peers
            const eventStrategy = eventStrategies.NETWORK_SCOPE_ANYFORTX;

            logger.info(`✅ Using NETWORK_SCOPE_ANYFORTX event strategy - waiting for commit confirmation`);

            // Create and connect gateway
            const gateway = new Gateway();

            await gateway.connect(connectionProfile, {
                wallet,
                identity: org.adminUser,
                // Enable discovery with asLocalhost for .nip.io URL translation
                // The warning about 0 endorsers is normal in Microfab - fallback works fine
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
                    // Get network instance
                    const network = await this.gateways[orgKey].getNetwork(config.channelName);
                    const channel = network.getChannel();
                    
                    // Try discovery first
                    let orgEndorsers = channel.getEndorsers().filter(peer => peer.mspid === mspId);

                    // If discovery fails, build peers manually from connection profile
                    if (orgEndorsers.length === 0) {
                        logger.warn(`⚠️ Discovery returned 0 endorsers for ${mspId}. Building from connection profile...`);
                        
                        const orgConfig = config.orgs[orgKey];
                        const profilePath = require('path').resolve(config.gatewaysDir, orgConfig.gatewayFile);
                        const profile = require(profilePath);

                        // Get peer URLs from profile
                        const peerUrls = profile.organizations?.[orgConfig.name]?.peers || [];
                        
                        for (const peerName of peerUrls) {
                            const peerConfig = profile.peers?.[peerName];
                            if (peerConfig && peerConfig.url) {
                                // Build endorser from connection profile data
                                const endorser = channel.client.newEndorser(peerName);
                                endorser.endpoint = channel.client.newEndpoint(peerConfig);
                                endorser.mspid = mspId;
                                
                                await endorser.connect();
                                orgEndorsers.push(endorser);
                                logger.info(`✅ Connected to peer: ${peerName} (${peerConfig.url})`);
                            }
                        }
                    }

                    if (orgEndorsers.length > 0) {
                        logger.info(`🎯 Using ${orgEndorsers.length} peer(s) for endorsement and events`);
                        const transaction = contract.createTransaction(fcn);
                        transaction.setEndorsingPeers(orgEndorsers);

                        const result = await transaction.submit(...args);
                        const resultStr = result.toString();
                        logger.info(`📥 Result: ${resultStr.substring(0, 100)}...`);

                        if (!resultStr || resultStr === '') return { success: true };
                        try { return JSON.parse(resultStr); } catch { return { success: true, message: resultStr }; }
                    }
                    
                    logger.warn(`⚠️ Could not find any peers, using default transaction`);
                } catch (pe) {
                    logger.warn(`⚠️ Peer setup failed: ${pe.message}, using default transaction`);
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
     * Attempts peer targeting for PDC queries, always falls back to default on any error
     */
    async evaluateTransaction(orgKey, contractName, functionName, ...args) {
        const contract = await this.getContract(orgKey);
        const fcn = `${contractName}:${functionName}`;

        logger.info(`🔍 Query: ${fcn} | Org: ${orgKey} | Args: ${JSON.stringify(args).substring(0, 100)}...`);

        // Try peer targeting for PDC access, but ALWAYS fall back if anything fails
        let peerTargetingWorked = false;
        try {
            const orgMspMap = {
                'startup': 'StartupOrgMSP',
                'investor': 'InvestorOrgMSP',
                'validator': 'ValidatorOrgMSP',
                'platform': 'PlatformOrgMSP'
            };
            const mspId = orgMspMap[orgKey];

            const network = await this.gateways[orgKey].getNetwork(config.channelName);
            const channel = network.getChannel();
            
            let orgEndorsers = channel.getEndorsers().filter(peer => peer.mspid === mspId);

            if (orgEndorsers.length === 0) {
                const orgConfig = config.orgs[orgKey];
                const profilePath = require('path').resolve(config.gatewaysDir, orgConfig.gatewayFile);
                const profile = require(profilePath);
                const peerUrls = profile.organizations?.[orgConfig.name]?.peers || [];
                
                for (const peerName of peerUrls) {
                    const peerConfig = profile.peers?.[peerName];
                    if (peerConfig && peerConfig.url) {
                        const endorser = channel.client.newEndorser(peerName);
                        endorser.endpoint = channel.client.newEndpoint(peerConfig);
                        endorser.mspid = mspId;
                        await endorser.connect();
                        orgEndorsers.push(endorser);
                    }
                }
            }

            if (orgEndorsers.length > 0) {
                const transaction = contract.createTransaction(fcn);
                transaction.setEndorsingPeers(orgEndorsers);
                const result = await transaction.evaluate(...args);
                const resultStr = result.toString();
                
                peerTargetingWorked = true;
                logger.info(`📥 Result (peer-targeted): ${resultStr.substring(0, 100)}...`);

                if (!resultStr || resultStr === '') return [];
                try { return JSON.parse(resultStr); } catch { return resultStr; }
            }
        } catch (err) {
            logger.warn(`⚠️ Peer targeting failed: ${err.message}, falling back to default`);
        }

        // Fallback to default query (always executes if peer targeting didn't work)
        if (!peerTargetingWorked) {
            try {
                const result = await contract.evaluateTransaction(fcn, ...args);
                const resultStr = result.toString();
                logger.info(`📥 Result (default): ${resultStr.substring(0, 100)}...`);

                if (!resultStr || resultStr === '') return [];
                try { return JSON.parse(resultStr); } catch { return resultStr; }
            } catch (error) {
                logger.error(`❌ Query failed: ${error.message}`);
                throw error;
            }
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
