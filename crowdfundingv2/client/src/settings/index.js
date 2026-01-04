import dotenv from 'dotenv';
import { fileURLToPath } from 'url';
import { dirname, resolve } from 'path';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

dotenv.config({ path: resolve(__dirname, '../../.env') });

export default {
    port: process.env.PORT || 3001,
    nodeEnv: process.env.NODE_ENV || 'development',

    // MongoDB configuration
    mongodb: {
        uri: process.env.MONGODB_URI || 'mongodb://localhost:27017/crowdfunding'
    },

    // JWT configuration
    jwt: {
        secret: process.env.JWT_SECRET || 'your-secret-key-change-in-production',
        expiresIn: process.env.JWT_EXPIRES_IN || '7d'
    },

    // Fee configuration (synced from chaincode)
    fees: {
        inrToCftRate: parseFloat(process.env.INR_TO_CFT_RATE) || 2.5,
        usdToCftRate: parseFloat(process.env.USD_TO_CFT_RATE) || 83.0,
        registrationFeeCft: parseInt(process.env.REGISTRATION_FEE_CFT) || 250,
        campaignCreationFeeCft: parseInt(process.env.CAMPAIGN_CREATION_FEE_CFT) || 1250,
        campaignPublishingFeeCft: parseInt(process.env.CAMPAIGN_PUBLISHING_FEE_CFT) || 2500,
        validationFeeCft: parseInt(process.env.VALIDATION_FEE_CFT) || 500,
        disputeFeeCft: parseInt(process.env.DISPUTE_FEE_CFT) || 750,
        investmentFeePercent: parseFloat(process.env.INVESTMENT_FEE_PERCENT) || 5,
        withdrawalFeePercent: parseFloat(process.env.WITHDRAWAL_FEE_PERCENT) || 1
    }
};
