const express = require('express');
const cors = require('cors');
const morgan = require('morgan');
const config = require('./config');
const logger = require('./utils/logger');
const fabricConnection = require('./fabricConnection');

const app = express();

// Middleware
app.use(cors());
app.use(express.json());
app.use(morgan('dev'));

// Helper to handle empty data gracefully
const safeQuery = async (res, orgKey, contractName, functionName, ...args) => {
    try {
        const result = await fabricConnection.evaluateTransaction(orgKey, contractName, functionName, ...args);
        res.json({ success: true, data: result || [] });
    } catch (error) {
        logger.warn(`${functionName}: ${error.message}`);
        res.json({ success: true, data: [] });
    }
};

// Health check
app.get('/health', (req, res) => {
    res.json({ status: 'ok', timestamp: new Date().toISOString() });
});

// ====================
// STARTUP ROUTES
// ====================
app.post('/api/startup/campaigns', async (req, res) => {
    try {
        const {
            campaignId, startupId, category, deadline, currency,
            hasRaised, hasGovGrants, incorpDate, projectStage, sector,
            tags, teamAvailable, investorCommitted, duration,
            fundingDay, fundingMonth, fundingYear, goalAmount,
            investmentRange, projectName, description, documents
        } = req.body;

        const result = await fabricConnection.submitTransaction(
            'startup', 'StartupContract', 'CreateCampaign',
            campaignId, startupId, category, deadline, currency,
            String(hasRaised), String(hasGovGrants), incorpDate, projectStage, sector,
            JSON.stringify(tags || []), String(teamAvailable), String(investorCommitted),
            String(duration), String(fundingDay), String(fundingMonth), String(fundingYear),
            String(goalAmount), investmentRange, projectName, description,
            JSON.stringify(documents || [])
        );
        res.status(201).json({ success: true, data: result });
    } catch (error) {
        logger.error(`CreateCampaign error: ${error.message}`);
        res.status(500).json({ error: error.message });
    }
});

// Get all campaigns for startup dashboard
app.get('/api/startup/campaigns', (req, res) => {
    safeQuery(res, 'startup', 'StartupContract', 'GetAllCampaigns');
});

app.get('/api/startup/campaigns/:campaignId', async (req, res) => {
    try {
        const result = await fabricConnection.evaluateTransaction(
            'startup', 'StartupContract', 'GetCampaign', req.params.campaignId
        );
        res.json({ success: true, data: result });
    } catch (error) {
        logger.warn(`GetCampaign: ${error.message}`);
        res.status(404).json({ error: 'Campaign not found' });
    }
});

app.post('/api/startup/campaigns/:campaignId/submit-validation', async (req, res) => {
    try {
        const { documents, notes } = req.body;
        const result = await fabricConnection.submitTransaction(
            'startup', 'StartupContract', 'SubmitForValidation',
            req.params.campaignId, JSON.stringify(documents || []), notes || ''
        );
        res.json({ success: true, data: result });
    } catch (error) {
        logger.error(`SubmitForValidation error: ${error.message}`);
        res.status(500).json({ error: error.message });
    }
});

app.post('/api/startup/campaigns/:campaignId/share-to-platform', async (req, res) => {
    try {
        const { validationHash } = req.body;
        const result = await fabricConnection.submitTransaction(
            'startup', 'StartupContract', 'ShareCampaignToPlatform',
            req.params.campaignId, validationHash || ''
        );
        res.json({ success: true, data: result });
    } catch (error) {
        logger.error(`ShareCampaignToPlatform error: ${error.message}`);
        res.status(500).json({ error: error.message });
    }
});

// ====================
// VALIDATOR ROUTES
// ====================
// Get pending validations for validator dashboard
app.get('/api/validator/pending-validations', (req, res) => {
    safeQuery(res, 'validator', 'ValidatorContract', 'GetPendingValidations');
});

app.get('/api/validator/campaigns/:campaignId', async (req, res) => {
    try {
        const result = await fabricConnection.evaluateTransaction(
            'validator', 'ValidatorContract', 'GetCampaign', req.params.campaignId
        );
        res.json({ success: true, data: result });
    } catch (error) {
        logger.warn(`ValidatorContract:GetCampaign: ${error.message}`);
        res.status(404).json({ error: 'Campaign not found' });
    }
});

app.post('/api/validator/validate/:campaignId', async (req, res) => {
    try {
        const { validationId, validatorId, dueDiligenceScore, riskScore, riskLevel, comments, issues, requiredDocuments } = req.body;
        const result = await fabricConnection.submitTransaction(
            'validator', 'ValidatorContract', 'ValidateCampaign',
            validationId, req.params.campaignId, validatorId,
            String(dueDiligenceScore), String(riskScore), riskLevel,
            JSON.stringify(comments || []), JSON.stringify(issues || []), requiredDocuments || ''
        );
        res.json({ success: true, data: result });
    } catch (error) {
        logger.error(`ValidateCampaign error: ${error.message}`);
        res.status(500).json({ error: error.message });
    }
});

app.post('/api/validator/approve/:campaignId', async (req, res) => {
    try {
        const { validationId, status, dueDiligenceScore, riskScore, riskLevel, comments, issues, requiredDocuments } = req.body;
        const campaignId = req.params.campaignId;
        const valId = validationId || `VAL_${Date.now()}`;

        // Step 1: First create the validation record with ValidateCampaign
        logger.info(`Step 1: Creating validation record for ${campaignId}`);
        await fabricConnection.submitTransaction(
            'validator', 'ValidatorContract', 'ValidateCampaign',
            valId, campaignId, 'VALIDATOR001',
            String(dueDiligenceScore || 8.5), String(riskScore || 3.0), riskLevel || 'LOW',
            JSON.stringify(comments || ['Validated']), JSON.stringify(issues || []), requiredDocuments || ''
        );

        // Step 2: Then approve/reject the campaign
        logger.info(`Step 2: Approving/Rejecting campaign ${campaignId}`);
        const result = await fabricConnection.submitTransaction(
            'validator', 'ValidatorContract', 'ApproveOrRejectCampaign',
            valId, campaignId, status || 'APPROVED',
            String(dueDiligenceScore || 8.5), String(riskScore || 3.0), riskLevel || 'LOW',
            JSON.stringify(comments || ['Approved']), JSON.stringify(issues || []), requiredDocuments || ''
        );

        res.json({ success: true, data: result, validationId: valId });
    } catch (error) {
        logger.error(`ApproveOrRejectCampaign error: ${error.message}`);
        res.status(500).json({ error: error.message });
    }
});

// ====================
// PLATFORM ROUTES
// ====================
// Get all shared campaigns for platform dashboard
app.get('/api/platform/shared-campaigns', (req, res) => {
    safeQuery(res, 'platform', 'PlatformContract', 'GetAllSharedCampaigns');
});

app.get('/api/platform/shared-campaigns/:campaignId', async (req, res) => {
    try {
        const result = await fabricConnection.evaluateTransaction(
            'platform', 'PlatformContract', 'GetSharedCampaign', req.params.campaignId
        );
        res.json({ success: true, data: result });
    } catch (error) {
        logger.warn(`GetSharedCampaign: ${error.message}`);
        res.status(404).json({ error: 'Campaign not found' });
    }
});

app.post('/api/platform/publish/:campaignId', async (req, res) => {
    try {
        const { validationHash } = req.body;
        const result = await fabricConnection.submitTransaction(
            'platform', 'PlatformContract', 'PublishCampaignToPortal',
            req.params.campaignId, validationHash || ''
        );
        res.json({ success: true, data: result });
    } catch (error) {
        logger.error(`PublishCampaignToPortal error: ${error.message}`);
        res.status(500).json({ error: error.message });
    }
});

// ====================
// INVESTOR ROUTES
// ====================
app.get('/api/investor/campaigns', (req, res) => {
    safeQuery(res, 'investor', 'InvestorContract', 'GetAvailableCampaigns');
});

app.get('/api/investor/campaigns/:campaignId', async (req, res) => {
    try {
        const result = await fabricConnection.evaluateTransaction(
            'investor', 'InvestorContract', 'ViewCampaignDetails', req.params.campaignId
        );
        res.json({ success: true, data: result });
    } catch (error) {
        logger.warn(`ViewCampaignDetails: ${error.message}`);
        res.status(404).json({ error: 'Campaign not found' });
    }
});

app.post('/api/investor/view/:campaignId', async (req, res) => {
    try {
        const { viewRecordId, investorId } = req.body;
        const result = await fabricConnection.submitTransaction(
            'investor', 'InvestorContract', 'ViewCampaign',
            viewRecordId || `VIEW_${Date.now()}`, req.params.campaignId, investorId
        );
        res.json({ success: true, data: result });
    } catch (error) {
        logger.error(`ViewCampaign error: ${error.message}`);
        res.status(500).json({ error: error.message });
    }
});

app.post('/api/investor/investments', async (req, res) => {
    try {
        const { investmentId, campaignId, investorId, amount, currency } = req.body;
        const result = await fabricConnection.submitTransaction(
            'investor', 'InvestorContract', 'MakeInvestment',
            investmentId, campaignId, investorId, String(amount), currency
        );
        res.status(201).json({ success: true, data: result });
    } catch (error) {
        logger.error(`MakeInvestment error: ${error.message}`);
        res.status(500).json({ error: error.message });
    }
});

app.post('/api/investor/proposals', async (req, res) => {
    try {
        const { proposalId, campaignId, investorId, startupId, investmentAmount, currency, equity, duration, milestones, proposedTerms } = req.body;
        const result = await fabricConnection.submitTransaction(
            'investor', 'InvestorContract', 'CreateInvestmentProposal',
            proposalId, campaignId, investorId, startupId,
            String(investmentAmount), currency, equity || '10', duration || '3 years',
            JSON.stringify(milestones || []), proposedTerms || ''
        );
        res.status(201).json({ success: true, data: result });
    } catch (error) {
        logger.error(`CreateInvestmentProposal error: ${error.message}`);
        res.status(500).json({ error: error.message });
    }
});

// Request validation details from validator
app.post('/api/investor/request-validation/:campaignId', async (req, res) => {
    try {
        const { requestId, investorId } = req.body;
        const result = await fabricConnection.submitTransaction(
            'investor', 'InvestorContract', 'RequestValidationDetails',
            requestId || `REQ_${Date.now()}`, req.params.campaignId, investorId || 'INVESTOR001'
        );
        res.json({ success: true, data: result });
    } catch (error) {
        logger.error(`RequestValidationDetails error: ${error.message}`);
        res.status(500).json({ error: error.message });
    }
});

// Get viewed campaigns (Recent Views)
app.get('/api/investor/viewed-campaigns', async (req, res) => {
    try {
        const result = await fabricConnection.evaluateTransaction(
            'investor', 'InvestorContract', 'GetViewedCampaigns'
        );
        res.json({ success: true, data: result || [] });
    } catch (error) {
        logger.warn(`GetViewedCampaigns: ${error.message}`);
        res.json({ success: true, data: [] });
    }
});

// Get investor's investments
app.get('/api/investor/my-investments', async (req, res) => {
    try {
        const result = await fabricConnection.evaluateTransaction(
            'investor', 'InvestorContract', 'GetMyInvestments'
        );
        res.json({ success: true, data: result || [] });
    } catch (error) {
        logger.warn(`GetMyInvestments: ${error.message}`);
        res.json({ success: true, data: [] });
    }
});

// Start server
app.listen(config.port, () => {
    logger.info(`🚀 Fabric Network API running on port ${config.port}`);
    logger.info(`📡 Channel: ${config.channelName}`);
    logger.info(`📦 Chaincode: ${config.chaincodeName}`);
});

// Graceful shutdown
process.on('SIGINT', async () => {
    logger.info('Shutting down...');
    await fabricConnection.disconnect();
    process.exit(0);
});
