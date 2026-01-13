import { Wallets, Gateway } from 'fabric-network';
import { promises as fs } from 'fs';
import { resolve } from 'path';
import config from '../settings/index.js';

/**
 * FabricGateway - Service to connect to Hyperledger Fabric network
 * Uses fabric-network SDK with Wallet pattern for microfab
 */
class FabricGateway {
    constructor() {
        this.connections = {};
    }

    /**
     * Get or create a connection for a specific organization
     */
    async getConnection(orgName) {
        if (this.connections[orgName]) {
            return this.connections[orgName];
        }

        const orgConfig = config.orgs[orgName];
        if (!orgConfig) {
            throw new Error(`Unknown organization: ${orgName}`);
        }

        try {
            // Load wallet from _wallets directory
            const walletPath = resolve(config.fabric.walletPath, orgConfig.mspId.replace('MSP', ''));
            console.log(`📂 Loading wallet from: ${walletPath}`);

            const wallet = await Wallets.newFileSystemWallet(walletPath);

            // Check if identity exists
            const identity = await wallet.get(orgConfig.adminUser);
            if (!identity) {
                throw new Error(`Identity ${orgConfig.adminUser} not found in wallet`);
            }
            console.log(`✅ Found identity: ${orgConfig.adminUser}`);

            // Load connection profile
            const gatewayPath = resolve(config.fabric.gatewayPath, orgConfig.gatewayFile);
            console.log(`📂 Loading gateway from: ${gatewayPath}`);

            const connectionProfile = JSON.parse(await fs.readFile(gatewayPath, 'utf8'));

            // Create gateway
            const gateway = new Gateway();
            await gateway.connect(connectionProfile, {
                wallet,
                identity: orgConfig.adminUser,
                discovery: { enabled: false, asLocalhost: false },
            });

            // Get network and contract
            const network = await gateway.getNetwork(config.fabric.channelName);
            const contract = network.getContract(config.fabric.chaincodeName);

            this.connections[orgName] = {
                gateway,
                network,
                contract,
            };

            console.log(`✅ Connected to Fabric as ${orgName}`);
            return this.connections[orgName];
        } catch (error) {
            console.error(`❌ Failed to connect as ${orgName}:`, error.message);
            throw error;
        }
    }

    /**
     * Submit a transaction (invoke)
     */
    async submitTransaction(orgName, contractName, functionName, ...args) {
        try {
            const { contract } = await this.getConnection(orgName);
            const fullFunctionName = `${contractName}:${functionName}`;

            console.log(`📤 Submitting: ${fullFunctionName}`, args);

            const result = await contract.submitTransaction(fullFunctionName, ...args);
            const resultStr = result.toString();

            console.log(`📥 Result: ${resultStr.substring(0, 200)}...`);

            if (!resultStr || resultStr === '') {
                return { success: true };
            }

            try {
                return JSON.parse(resultStr);
            } catch {
                return { success: true, message: resultStr };
            }
        } catch (error) {
            console.error(`❌ Transaction failed:`, error.message);
            throw error;
        }
    }

    /**
     * Evaluate a transaction (query)
     */
    async evaluateTransaction(orgName, contractName, functionName, ...args) {
        try {
            const { contract } = await this.getConnection(orgName);
            const fullFunctionName = `${contractName}:${functionName}`;

            console.log(`🔍 Querying: ${fullFunctionName}`, args);

            const result = await contract.evaluateTransaction(fullFunctionName, ...args);
            const resultStr = result.toString();

            console.log(`📥 Result: ${resultStr.substring(0, 200)}...`);

            if (!resultStr || resultStr === '') {
                return [];
            }

            try {
                return JSON.parse(resultStr);
            } catch {
                return resultStr;
            }
        } catch (error) {
            console.error(`❌ Query failed:`, error.message);
            throw error;
        }
    }

    /**
     * Close all connections
     */
    async close() {
        for (const [orgName, conn] of Object.entries(this.connections)) {
            conn.gateway.disconnect();
            console.log(`🔌 Disconnected from ${orgName}`);
        }
        this.connections = {};
    }
}

// Export singleton instance
export default new FabricGateway();
