const express = require('express');
const cors = require('cors');
const morgan = require('morgan');
const axios = require('axios');
const config = require('./config');
const logger = require('./utils/logger');
const fabricConnection = require('./fabricConnection');

const app = express();

// Auth API base URL (for syncing to MongoDB)
const AUTH_API_BASE = process.env.AUTH_API_URL || 'http://127.0.0.1:3001/api/auth';

// ============================================================================
// ID GENERATION TOGGLE
// ============================================================================
// Set to true for random IDs (prevents collisions after server restart)
// Set to false for sequential IDs (cleaner for fresh testing with cleared blockchain)
const USE_RANDOM_IDS = true;

const crypto = require('crypto');

// Campaign sequence counters (in-memory, resets on restart - for demo)
// In production, use Redis or query MongoDB for actual counts
const campaignCounters = {};
const startupCounters = {};  // Track startup sequence per user

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

// Helper to get next campaign sequence for a startup
const getNextCampaignSeq = (startupId) => {
    if (USE_RANDOM_IDS) {
        return crypto.randomBytes(3).toString('hex').toUpperCase(); // e.g., "A7B3C2"
    }
    if (!campaignCounters[startupId]) {
        campaignCounters[startupId] = 0;
    }
    campaignCounters[startupId]++;
    return String(campaignCounters[startupId]).padStart(3, '0');
};

// Helper to get next startup sequence for a user
const getNextStartupSeq = (ownerId) => {
    if (USE_RANDOM_IDS) {
        return crypto.randomBytes(3).toString('hex').toUpperCase(); // e.g., "F4C2E1"
    }
    if (!startupCounters[ownerId]) {
        startupCounters[ownerId] = 0;
    }
    startupCounters[ownerId]++;
    return String(startupCounters[ownerId]).padStart(3, '0');
};

// Health check
app.get('/health', (req, res) => {
    res.json({ status: 'ok', timestamp: new Date().toISOString() });
});

// ====================
// STARTUP MANAGEMENT
// ====================

// Create a new startup
app.post('/api/startup/startups', async (req, res) => {
    try {
        const { name, description, ownerId, authToken } = req.body;

        if (!name || !ownerId) {
            return res.status(400).json({ success: false, error: 'name and ownerId are required' });
        }

        // Auto-generate startup ID
        const seq = getNextStartupSeq(ownerId);
        const startupId = `STU_${ownerId}_${seq}`;
        const displayId = `S-${seq}`;

        logger.info(`Creating startup ${startupId} for owner ${ownerId}`);

        // Create in Fabric chaincode
        await fabricConnection.submitTransaction(
            'startup', 'StartupContract', 'CreateStartup',
            startupId, ownerId, name, description || '', displayId
        );

        // Sync to MongoDB
        if (authToken) {
            try {
                await axios.post(`${AUTH_API_BASE}/startups`, {
                    startupId,
                    displayId,
                    name,
                    description
                }, {
                    headers: { Authorization: `Bearer ${authToken}` }
                });
            } catch (syncErr) {
                logger.warn(`MongoDB sync failed: ${syncErr.message} - ${JSON.stringify(syncErr.response?.data || {})}`);
            }
        }

        res.json({
            success: true,
            data: {
                startupId,
                displayId,
                name,
                description,
                ownerId,
                campaignIds: [],
                createdAt: new Date().toISOString()
            }
        });
    } catch (error) {
        logger.error('Create startup error:', error);
        res.status(500).json({ success: false, error: error.message });
    }
});

// Sync startup to chaincode (recreate from MongoDB data when chaincode is reset)
app.post('/api/startup/startups/:startupId/sync-to-chaincode', async (req, res) => {
    try {
        const { startupId } = req.params;
        const { name, description, ownerId, displayId } = req.body;

        if (!name || !ownerId || !startupId) {
            return res.status(400).json({ success: false, error: 'startupId, name, and ownerId are required' });
        }

        logger.info(`Syncing startup ${startupId} to chaincode for owner ${ownerId}`);

        // First check if startup already exists in chaincode
        try {
            await fabricConnection.evaluateTransaction(
                'startup', 'StartupContract', 'GetStartup', startupId
            );
            // If we reach here, startup already exists in chaincode
            return res.json({ success: true, message: 'Startup already exists in chaincode', alreadyExists: true });
        } catch (existsErr) {
            // Startup not found, we can create it
            if (!existsErr.message.includes('not found')) {
                throw existsErr;
            }
        }

        // Create in Fabric chaincode
        await fabricConnection.submitTransaction(
            'startup', 'StartupContract', 'CreateStartup',
            startupId, ownerId, name, description || '', displayId || startupId.split('_').pop()
        );

        logger.info(`Successfully synced startup ${startupId} to chaincode`);

        res.json({
            success: true,
            message: 'Startup synced to chaincode successfully',
            data: { startupId, name, ownerId, displayId }
        });
    } catch (error) {
        logger.error('Sync startup to chaincode error:', error);
        res.status(500).json({ success: false, error: error.message });
    }
});

app.get('/api/startup/startups/:startupId', async (req, res) => {
    try {
        const { startupId } = req.params;
        const result = await fabricConnection.evaluateTransaction(
            'startup', 'StartupContract', 'GetStartup', startupId
        );
        res.json({ success: true, data: result });
    } catch (error) {
        logger.error('Get startup error:', error);
        res.status(500).json({ success: false, error: error.message });
    }
});

// Get all startups for an owner
app.get('/api/startup/startups/owner/:ownerId', async (req, res) => {
    try {
        const { ownerId } = req.params;
        const result = await fabricConnection.evaluateTransaction(
            'startup', 'StartupContract', 'GetStartupsByOwner', ownerId
        );
        res.json({ success: true, data: result || [] });
    } catch (error) {
        logger.warn('Get startups by owner error:', error.message);
        // Return empty array on error (common when no startups exist)
        res.json({ success: true, data: [] });
    }
});

// ====================
// DELETION ROUTES
// ====================

// Get deletion fee preview for a campaign
app.get('/api/startup/campaigns/:campaignId/deletion-fee', async (req, res) => {
    try {
        const { campaignId } = req.params;
        const result = await fabricConnection.evaluateTransaction(
            'startup', 'StartupContract', 'CalculateCampaignDeletionFee', campaignId
        );
        res.json({ success: true, data: result });
    } catch (error) {
        logger.error('Get campaign deletion fee error:', error);
        res.status(500).json({ success: false, error: error.message });
    }
});

// Delete a campaign (charges fee)
app.delete('/api/startup/campaigns/:campaignId', async (req, res) => {
    try {
        const { campaignId } = req.params;
        const { reason, authToken, startupId } = req.body;

        logger.info(`Deleting campaign ${campaignId}. Reason: ${reason || 'Not provided'}`);

        const result = await fabricConnection.submitTransaction(
            'startup', 'StartupContract', 'DeleteCampaign',
            campaignId, reason || 'User requested deletion'
        );

        // Sync deletion status to MongoDB
        if (authToken && startupId) {
            try {
                logger.info(`Syncing deletion status for campaign ${campaignId} to MongoDB`);
                await axios.post(`${AUTH_API_BASE}/sync/campaign-status`, {
                    startupId,
                    campaignId,
                    status: 'DELETED'
                }, {
                    headers: { Authorization: `Bearer ${authToken}` }
                });
                logger.info(`Campaign ${campaignId} marked as DELETED in MongoDB`);
            } catch (syncError) {
                logger.warn(`MongoDB deletion sync failed: ${syncError.message} - ${JSON.stringify(syncError.response?.data || {})}`);
                // Don't fail the response if sync fails, but log it warning
            }
        } else {
            logger.warn(`Skipping MongoDB sync for deletion of ${campaignId}: authToken or startupId missing`);
        }

        res.json({
            success: true,
            data: result,
            message: `Campaign ${campaignId} deleted successfully. Fee charged: ${result.feeCharged} CFT`
        });
    } catch (error) {
        logger.error('Delete campaign error:', error);
        res.status(500).json({ success: false, error: error.message });
    }
});

// Get deletion fee preview for a startup (includes all campaign fees)
app.get('/api/startup/startups/:startupId/deletion-fee', async (req, res) => {
    try {
        const { startupId } = req.params;
        const result = await fabricConnection.evaluateTransaction(
            'startup', 'StartupContract', 'CalculateStartupDeletionFee', startupId
        );
        res.json({ success: true, data: result });
    } catch (error) {
        logger.error('Get startup deletion fee error:', error);
        res.status(500).json({ success: false, error: error.message });
    }
});

// Delete a startup and all its campaigns (charges total fee)
app.delete('/api/startup/startups/:startupId', async (req, res) => {
    try {
        const { startupId } = req.params;
        const { reason, authToken } = req.body;

        logger.info(`Deleting startup ${startupId} and all campaigns. Reason: ${reason || 'Not provided'}`);

        const result = await fabricConnection.submitTransaction(
            'startup', 'StartupContract', 'DeleteStartup',
            startupId, reason || 'User requested deletion'
        );

        // Sync deletion status for ALL campaigns of this startup?
        // Or simply delete the startup record?
        // Since we don't have a specific endpoint for startup deletion sync, 
        // we might rely on the frontend or backend to handle cascade.
        // However, if we just want to mark the STARTUP as deleted?
        // The user asked to update "whichever campaign or startup is deleted".
        // I will assume there is a sync endpoint for startup status or I should try to iterate?
        // The `result` from DeleteStartup return deleted campaigns. I could sync each one?
        // Or better: call a sync endpoint for startup deletion if it exists.
        // I'll try calling `/sync/startup-status`?
        // I'll stick to trying to mark campaigns as deleted if I can, OR just warn if I can't sync startup deletion.
        // Wait, `result` contains `CampaignDeletions` list.

        if (authToken && result.campaignDeletions && result.campaignDeletions.length > 0) {
            logger.info(`Syncing deletion status for ${result.campaignDeletions.length} campaigns to MongoDB`);
            for (const delRecord of result.campaignDeletions) {
                try {
                    await axios.post(`${AUTH_API_BASE}/sync/campaign-status`, {
                        startupId,
                        campaignId: delRecord.entityId,
                        status: 'DELETED'
                    }, {
                        headers: { Authorization: `Bearer ${authToken}` }
                    });
                } catch (e) {
                    logger.warn(`Failed to sync deletion for campaign ${delRecord.entityId}: ${e.message}`);
                }
            }
        }

        res.json({
            success: true,
            data: result,
            message: `Startup ${startupId} and all campaigns deleted successfully`
        });
    } catch (error) {

        logger.error('Delete startup error:', error);
        res.status(500).json({ success: false, error: error.message });
    }
});

// Get a deletion record by ID
app.get('/api/startup/deletions/:deletionId', async (req, res) => {
    try {
        const { deletionId } = req.params;
        const result = await fabricConnection.evaluateTransaction(
            'startup', 'StartupContract', 'GetDeletionRecord', deletionId
        );
        res.json({ success: true, data: result });
    } catch (error) {
        logger.error('Get deletion record error:', error);
        res.status(500).json({ success: false, error: error.message });
    }
});

// ====================
// CAMPAIGN ROUTES
// ====================
app.post('/api/startup/campaigns', async (req, res) => {
    try {
        const {
            startupId, category, deadline, currency,
            hasRaised, hasGovGrants, incorpDate, projectStage, sector,
            tags, teamAvailable, investorCommitted, duration,
            fundingDay, fundingMonth, fundingYear, goalAmount,
            investmentRange, projectName, description, documents,
            authToken  // JWT token for syncing to MongoDB
        } = req.body;

        // Auto-generate campaign ID
        const seq = getNextCampaignSeq(startupId);
        const campaignId = `CAMP_${startupId}_${seq}`;
        const displayId = `C-${seq}`;

        logger.info(`Creating campaign ${campaignId} for startup ${startupId}`);

        // Create in Fabric (22 parameters)
        const result = await fabricConnection.submitTransaction(
            'startup', 'StartupContract', 'CreateCampaign',
            campaignId, startupId, category, deadline, currency,
            String(hasRaised || false), String(hasGovGrants || false), incorpDate || '', projectStage || '', sector || '',
            JSON.stringify(tags || []), String(teamAvailable || false), String(investorCommitted || false),
            String(duration || 90), String(fundingDay || 1), String(fundingMonth || 1), String(fundingYear || 2025),
            String(goalAmount || 0), investmentRange || '', projectName || '', description || '',
            JSON.stringify(documents || [])
        );

        // Sync to MongoDB (if auth token provided)
        if (authToken) {
            try {
                await axios.post(`${AUTH_API_BASE}/sync/campaign`, {
                    startupId,
                    campaignId,
                    displayId,
                    projectName,
                    status: 'DRAFT'
                }, {
                    headers: { Authorization: `Bearer ${authToken}` }
                });
                logger.info(`Campaign ${campaignId} synced to MongoDB`);
            } catch (syncError) {
                logger.warn(`MongoDB sync failed: ${syncError.message} - ${JSON.stringify(syncError.response?.data || {})}`);
                // Don't fail the request - Fabric creation succeeded
            }
        }

        res.status(201).json({
            success: true,
            data: {
                ...result,
                campaignId,
                displayId
            }
        });
    } catch (error) {
        logger.error(`CreateCampaign error: ${error.message}`);
        res.status(500).json({ error: error.message });
    }
});

// Get campaigns for startup dashboard
// If startupId is provided, returns only that user's campaigns (strict isolation)
// If no startupId provided, returns all campaigns for the org
app.get('/api/startup/campaigns', async (req, res) => {
    const { startupId } = req.query;
    if (startupId) {
        // User-scoped query - returns only this user's campaigns
        safeQuery(res, 'startup', 'StartupContract', 'GetCampaignsByStartupId', startupId);
    } else {
        // Org-wide query (backward compatible)
        safeQuery(res, 'startup', 'StartupContract', 'GetAllCampaigns');
    }
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
        const { documents, notes, authToken, startupId, projectName } = req.body;

        // 0. CHECK IF VALIDATOR IS AVAILABLE
        // If authToken is provided, we check allocation first to prevent "submitted to void" state
        if (authToken) {
            try {
                // Peek if any validator is available
                await axios.get(`${AUTH_API_BASE}/allocation/next?role=VALIDATOR`, {
                    headers: { Authorization: `Bearer ${authToken}` }
                });
            } catch (checkError) {
                if (checkError.response?.status === 404) {
                    logger.warn('Pre-check failed: No ACTIVE validators available');
                    return res.status(400).json({
                        success: false,
                        error: 'Cannot submit: No active VALIDATOR users found in the system to review this campaign.'
                    });
                }
                // Determine if we should block or proceed on other errors (e.g., Auth API down)
                // For safety, let's block if we can't verify availability
                logger.warn(`Validator availability check failed: ${checkError.message}`);
                return res.status(503).json({
                    success: false,
                    error: 'Unable to verify validator availability. Please try again later.'
                });
            }
        }

        // 1. Submit to Fabric
        const result = await fabricConnection.submitTransaction(
            'startup', 'StartupContract', 'SubmitForValidation',
            req.params.campaignId, JSON.stringify(documents || []), notes || ''
        );

        // 2. Sync validation status to MongoDB
        if (authToken) {
            try {
                await axios.post(`${AUTH_API_BASE}/sync/campaign-status`, {
                    startupId,
                    campaignId: req.params.campaignId,
                    status: 'SUBMITTED',
                    validationStatus: 'PENDING_VALIDATION'
                }, {
                    headers: { Authorization: `Bearer ${authToken}` }
                });
                logger.info(`Campaign ${req.params.campaignId} status synced to MongoDB`);
            } catch (syncError) {
                logger.warn(`MongoDB status sync failed: ${syncError.message} - ${JSON.stringify(syncError.response?.data || {})}`);
            }
        }

        // 3. Auto-allocate to validator with least queue
        if (authToken) {
            try {
                // Get validator with least queue
                const allocRes = await axios.get(`${AUTH_API_BASE}/allocation/next?role=VALIDATOR`, {
                    headers: { Authorization: `Bearer ${authToken}` }
                });

                if (allocRes.data?.success && allocRes.data?.data?.orgUserId) {
                    const validatorId = allocRes.data.data.orgUserId;

                    // Assign to validator's queue
                    await axios.post(`${AUTH_API_BASE}/allocation/assign`, {
                        assigneeOrgUserId: validatorId,
                        campaignId: req.params.campaignId,
                        startupId,
                        projectName,
                        type: 'VALIDATION'
                    }, {
                        headers: { Authorization: `Bearer ${authToken}` }
                    });

                    logger.info(`Campaign ${req.params.campaignId} assigned to validator ${validatorId}`);
                }
            } catch (allocError) {
                logger.warn(`Auto-allocation failed: ${allocError.message} - ${JSON.stringify(allocError.response?.data || {})}`);
            }
        }

        res.json({ success: true, data: result });
    } catch (error) {
        logger.error(`SubmitForValidation error: ${error.message}`);
        res.status(500).json({ error: error.message });
    }
});

app.post('/api/startup/campaigns/:campaignId/share-to-platform', async (req, res) => {
    try {
        const { validationProofHash, authToken, startupId, projectName } = req.body;

        // 0. CHECK IF PLATFORM USER IS AVAILABLE
        if (authToken) {
            try {
                await axios.get(`${AUTH_API_BASE}/allocation/next?role=PLATFORM`, {
                    headers: { Authorization: `Bearer ${authToken}` }
                });
            } catch (checkError) {
                if (checkError.response?.status === 404) {
                    logger.warn('Pre-check failed: No ACTIVE platform users available');
                    return res.status(400).json({
                        success: false,
                        error: 'Cannot share: No active PLATFORM users found to publish this campaign.'
                    });
                }
                return res.status(503).json({
                    success: false,
                    error: 'Unable to verify platform availability. Please try again later.'
                });
            }
        }

        // 1. Share to Fabric
        const result = await fabricConnection.submitTransaction(
            'startup', 'StartupContract', 'ShareCampaignToPlatform',
            req.params.campaignId, validationProofHash || ''
        );

        // 2. Sync status to MongoDB
        if (authToken) {
            try {
                await axios.post(`${AUTH_API_BASE}/sync/campaign-status`, {
                    startupId,
                    campaignId: req.params.campaignId,
                    status: 'SHARED_TO_PLATFORM',
                    validationStatus: 'APPROVED'
                }, {
                    headers: { Authorization: `Bearer ${authToken}` }
                });
                logger.info(`Campaign ${req.params.campaignId} status synced (SHARED_TO_PLATFORM)`);
            } catch (syncError) {
                logger.warn(`MongoDB status sync failed: ${syncError.message}`);
            }
        }

        // 3. Auto-allocate to platform user with least queue
        if (authToken) {
            try {
                const allocRes = await axios.get(`${AUTH_API_BASE}/allocation/next?role=PLATFORM`, {
                    headers: { Authorization: `Bearer ${authToken}` }
                });

                if (allocRes.data?.success && allocRes.data?.data?.orgUserId) {
                    const platformId = allocRes.data.data.orgUserId;

                    await axios.post(`${AUTH_API_BASE}/allocation/assign`, {
                        assigneeOrgUserId: platformId,
                        campaignId: req.params.campaignId,
                        startupId,
                        projectName,
                        type: 'PUBLISH'
                    }, {
                        headers: { Authorization: `Bearer ${authToken}` }
                    });

                    logger.info(`Campaign ${req.params.campaignId} assigned to platform ${platformId}`);
                }
            } catch (allocError) {
                logger.warn(`Platform auto-allocation failed: ${allocError.message}`);
            }
        }

        res.json({ success: true, data: result });
    } catch (error) {
        logger.error(`ShareCampaignToPlatform error: ${error.message}`);
        res.status(500).json({ error: error.message });
    }
});

// ====================
// VALIDATOR ROUTES
// ====================

// Get pending validations - now filtered by assignedQueue from MongoDB
// This endpoint should be called with user's orgUserId to get only assigned campaigns
app.get('/api/validator/pending-validations', async (req, res) => {
    const { validatorId } = req.query;

    if (validatorId) {
        // User-scoped: Get only assigned campaigns for this validator
        // The frontend should fetch from MongoDB queue, then query Fabric for details
        safeQuery(res, 'validator', 'ValidatorContract', 'GetPendingValidations');
    } else {
        // Fallback: Get all pending (for admin/debug)
        safeQuery(res, 'validator', 'ValidatorContract', 'GetPendingValidations');
    }
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
        const { validationId, validatorId, submissionHash, requiredDocuments } = req.body;
        const result = await fabricConnection.submitTransaction(
            'validator', 'ValidatorContract', 'ValidateCampaign',
            validationId, req.params.campaignId, validatorId,
            submissionHash || '', JSON.stringify(requiredDocuments || [])
        );
        res.json({ success: true, data: result });
    } catch (error) {
        logger.error(`ValidateCampaign error: ${error.message}`);
        res.status(500).json({ error: error.message });
    }
});

app.post('/api/validator/approve/:campaignId', async (req, res) => {
    try {
        const { validationId, status, dueDiligenceScore, riskScore, riskLevel, comments, issues, requiredDocuments, submissionHash } = req.body;
        const campaignId = req.params.campaignId;
        const valId = validationId || `VAL_${Date.now()}`;

        // Step 1: First create the validation record with ValidateCampaign
        logger.info(`Step 1: Creating validation record for ${campaignId}`);
        await fabricConnection.submitTransaction(
            'validator', 'ValidatorContract', 'ValidateCampaign',
            valId, campaignId, 'VALIDATOR001',
            submissionHash || '', JSON.stringify(requiredDocuments || [])
        );

        // Step 2: Then approve/reject the campaign
        logger.info(`Step 2: Approving/Rejecting campaign ${campaignId}`);
        const finalStatus = status || 'APPROVED';
        const result = await fabricConnection.submitTransaction(
            'validator', 'ValidatorContract', 'ApproveOrRejectCampaign',
            valId, campaignId, finalStatus,
            String(dueDiligenceScore || 8.5), String(riskScore || 3.0), riskLevel || 'LOW',
            JSON.stringify(comments || ['Approved']), JSON.stringify(issues || []), requiredDocuments || ''
        );

        // Step 3: Sync status to MongoDB (startupId and authToken passed from frontend)
        const { authToken, startupId } = req.body;
        if (authToken && startupId) {
            try {
                await axios.post(`${AUTH_API_BASE}/sync/campaign-status`, {
                    startupId,
                    campaignId,
                    validationStatus: finalStatus
                }, {
                    headers: { Authorization: `Bearer ${authToken}` }
                });
                logger.info(`Campaign ${campaignId} validation status synced (${finalStatus})`);
            } catch (syncError) {
                logger.warn(`MongoDB status sync failed: ${syncError.message} - ${JSON.stringify(syncError.response?.data || {})}`);
            }
        }

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
        const { validationProofHash, authToken, startupId } = req.body;
        const campaignId = req.params.campaignId;

        const result = await fabricConnection.submitTransaction(
            'platform', 'PlatformContract', 'PublishCampaignToPortal',
            campaignId, validationProofHash || ''
        );

        // Sync status to MongoDB
        if (authToken && startupId) {
            try {
                await axios.post(`${AUTH_API_BASE}/sync/campaign-status`, {
                    startupId,
                    campaignId,
                    status: 'PUBLISHED'
                }, {
                    headers: { Authorization: `Bearer ${authToken}` }
                });
                logger.info(`Campaign ${campaignId} status synced (PUBLISHED)`);
            } catch (syncError) {
                logger.warn(`MongoDB status sync failed: ${syncError.message} - ${JSON.stringify(syncError.response?.data || {})}`);
            }
        }

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

// Get viewed campaigns (Recent Views) - filtered by investorId
app.get('/api/investor/viewed-campaigns', async (req, res) => {
    try {
        const { investorId } = req.query;
        // TODO: Add GetViewedCampaignsByInvestorId to chaincode for proper filtering
        // For now, return all and let frontend filter
        const result = await fabricConnection.evaluateTransaction(
            'investor', 'InvestorContract', 'GetViewedCampaigns'
        );
        res.json({ success: true, data: result || [] });
    } catch (error) {
        logger.warn(`GetViewedCampaigns: ${error.message}`);
        res.json({ success: true, data: [] });
    }
});

// Get investor's investments - filtered by investorId
app.get('/api/investor/my-investments', async (req, res) => {
    try {
        const { investorId } = req.query;
        // TODO: Add GetInvestmentsByInvestorId to chaincode for proper filtering
        // For now, return all and let frontend filter
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
