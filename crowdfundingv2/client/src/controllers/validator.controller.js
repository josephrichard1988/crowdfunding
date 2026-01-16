import fabricGateway from '../services/fabricGateway.js';

const ORG = 'validator';
const CONTRACT = 'ValidatorContract';

/**
 * Get campaign details for validation
 */
export async function getCampaign(req, res, next) {
    try {
        const { campaignId } = req.params;
        const result = await fabricGateway.evaluateTransaction(
            ORG, CONTRACT, 'GetCampaign', campaignId
        );
        res.json({ success: true, data: result });
    } catch (error) {
        next(error);
    }
}

/**
 * Get all pending validations
 */
export async function getPendingValidations(req, res, next) {
    try {
        const result = await fabricGateway.evaluateTransaction(
            ORG, CONTRACT, 'GetPendingValidations'
        );
        res.json({ success: true, data: result });
    } catch (error) {
        next(error);
    }
}

/**
 * Validate a campaign
 */
export async function validateCampaign(req, res, next) {
    try {
        const { campaignId } = req.params;
        const {
            validationId, validatorId, dueDiligenceScore, riskScore,
            riskLevel, comments, issues, requiredDocuments
        } = req.body;

        const result = await fabricGateway.submitTransaction(
            ORG, CONTRACT, 'ValidateCampaign',
            validationId, campaignId, validatorId,
            String(dueDiligenceScore), String(riskScore), riskLevel,
            JSON.stringify(comments || []), JSON.stringify(issues || []),
            requiredDocuments || ''
        );

        res.json({ success: true, data: result });
    } catch (error) {
        next(error);
    }
}

/**
 * Approve or reject a campaign
 */
export async function approveCampaign(req, res, next) {
    try {
        const { campaignId } = req.params;
        const {
            validationId, status, dueDiligenceScore, riskScore,
            riskLevel, comments, issues, requiredDocuments
        } = req.body;

        const result = await fabricGateway.submitTransaction(
            ORG, CONTRACT, 'ApproveOrRejectCampaign',
            validationId, campaignId, status,
            String(dueDiligenceScore), String(riskScore), riskLevel,
            JSON.stringify(comments || []), JSON.stringify(issues || []),
            requiredDocuments || ''
        );

        res.json({ success: true, data: result });
    } catch (error) {
        next(error);
    }
}

/**
 * Get validation requests from investors
 */
export async function getValidationRequests(req, res, next) {
    try {
        const result = await fabricGateway.evaluateTransaction(
            ORG, CONTRACT, 'GetValidationRequests'
        );
        res.json({ success: true, data: result });
    } catch (error) {
        next(error);
    }
}

/**
 * Provide validation details to investor
 */
export async function provideValidationDetails(req, res, next) {
    try {
        const { requestId } = req.params;
        const { campaignId } = req.body;

        const result = await fabricGateway.submitTransaction(
            ORG, CONTRACT, 'ProvideValidationDetailsToInvestor',
            requestId, campaignId
        );

        res.json({ success: true, data: result });
    } catch (error) {
        next(error);
    }
}

/**
 * Get all validation history for the organization
 * Shows all validations performed by any validator in the org
 */
export async function getAllValidationHistory(req, res, next) {
    try {
        const result = await fabricGateway.evaluateTransaction(
            ORG, CONTRACT, 'GetAllValidations'
        );
        res.json({ success: true, data: result });
    } catch (error) {
        console.error('Failed to get validation history:', error.message);
        // Return empty array if function doesn't exist yet
        res.json({ success: true, data: [] });
    }
}

/**
 * Verify milestone completion
 */
export async function verifyMilestone(req, res, next) {
    try {
        const { milestoneId } = req.params;
        const {
            verificationId, campaignId, startupId, reportHash,
            verified, score, comments, recommendRelease
        } = req.body;

        const result = await fabricGateway.submitTransaction(
            ORG, CONTRACT, 'VerifyMilestoneCompletion',
            verificationId, milestoneId, campaignId, startupId,
            reportHash, String(verified), String(score),
            comments, String(recommendRelease)
        );

        res.json({ success: true, data: result });
    } catch (error) {
        next(error);
    }
}
