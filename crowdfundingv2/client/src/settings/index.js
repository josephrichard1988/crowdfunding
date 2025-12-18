import dotenv from 'dotenv';
import { fileURLToPath } from 'url';
import { dirname, resolve } from 'path';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

dotenv.config({ path: resolve(__dirname, '../../.env.example') });

export default {
    port: process.env.PORT || 3001,
    nodeEnv: process.env.NODE_ENV || 'development',

    fabric: {
        channelName: process.env.CHANNEL_NAME || 'crowdfunding-channel',
        chaincodeName: process.env.CHAINCODE_NAME || 'crowdfunding',
        gatewayPath: resolve(__dirname, '../../../_gateways'),
        walletPath: resolve(__dirname, '../../../_wallets'),
        mspPath: resolve(__dirname, '../../../_msp'),
    },

    peers: {
        orderer: process.env.ORDERER_ENDPOINT || 'orderer-api.127-0-0-1.nip.io:9090',
        startup: process.env.STARTUP_PEER || 'startuporgpeer-api.127-0-0-1.nip.io:9090',
        investor: process.env.INVESTOR_PEER || 'investororgpeer-api.127-0-0-1.nip.io:9090',
        platform: process.env.PLATFORM_PEER || 'platformorgpeer-api.127-0-0-1.nip.io:9090',
        validator: process.env.VALIDATOR_PEER || 'validatororgpeer-api.127-0-0-1.nip.io:9090',
    },

    orgs: {
        startup: {
            mspId: 'StartupOrgMSP',
            peerEndpoint: 'startuporgpeer-api.127-0-0-1.nip.io:9090',
            gatewayFile: 'startuporggateway.json',
            adminUser: 'startuporgadmin',
        },
        investor: {
            mspId: 'InvestorOrgMSP',
            peerEndpoint: 'investororgpeer-api.127-0-0-1.nip.io:9090',
            gatewayFile: 'investororggateway.json',
            adminUser: 'investororgadmin',
        },
        platform: {
            mspId: 'PlatformOrgMSP',
            peerEndpoint: 'platformorgpeer-api.127-0-0-1.nip.io:9090',
            gatewayFile: 'platformorggateway.json',
            adminUser: 'platformorgadmin',
        },
        validator: {
            mspId: 'ValidatorOrgMSP',
            peerEndpoint: 'validatororgpeer-api.127-0-0-1.nip.io:9090',
            gatewayFile: 'validatororggateway.json',
            adminUser: 'validatororgadmin',
        },
    },
};
