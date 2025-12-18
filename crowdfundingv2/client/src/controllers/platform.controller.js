import fabricGateway from '../services/fabricGateway.js';

const ORG = 'platform';
const CONTRACT = 'PlatformContract';

/**
 * Get shared campaign from startup
 */
export async function getSharedCampaign(req, res, next) {
    try {
        const { campaignId } = req.params;
        const result = await fabricGateway.evaluateTransaction(
            ORG, CONTRACT, 'GetSharedCampaign', campaignId
        );
        res.json({ success: true, data: result });
    } catch (error) {
        next(error);
    }
}

/**
 * Get all shared campaigns
 */
export async function getAllSharedCampaigns(req, res, next) {
    try {
        const result = await fabricGateway.evaluateTransaction(
            ORG, CONTRACT, 'GetAllSharedCampaigns'
        );
        res.json({ success: true, data: result });
    } catch (error) {
        next(error);
    }
}

/**
 * Publish campaign to portal
 */
export async function publishCampaign(req, res, next) {
    try {
        const { campaignId } = req.params;
        const { validationHash } = req.body;

        const result = await fabricGateway.submitTransaction(
            ORG, CONTRACT, 'PublishCampaignToPortal',
            campaignId, validationHash
        );

        res.json({ success: true, data: result });
    } catch (error) {
        next(error);
    }
}

/**
 * Get published campaign
 */
export async function getPublishedCampaign(req, res, next) {
    try {
        const { campaignId } = req.params;
        const result = await fabricGateway.evaluateTransaction(
            ORG, CONTRACT, 'GetPublishedCampaign', campaignId
        );
        res.json({ success: true, data: result });
    } catch (error) {
        next(error);
    }
}

/**
 * Create wallet
 */
export async function createWallet(req, res, next) {
    try {
        const { walletId, ownerId, ownerType, initialBalance } = req.body;

        const result = await fabricGateway.submitTransaction(
            ORG, CONTRACT, 'CreateWallet',
            walletId, ownerId, ownerType, String(initialBalance)
        );

        res.status(201).json({ success: true, data: result });
    } catch (error) {
        next(error);
    }
}

/**
 * Get wallet
 */
export async function getWallet(req, res, next) {
    try {
        const { walletId } = req.params;
        const result = await fabricGateway.evaluateTransaction(
            ORG, CONTRACT, 'GetWallet', walletId
        );
        res.json({ success: true, data: result });
    } catch (error) {
        next(error);
    }
}

/**
 * Set campaign fee tier
 */
export async function setCampaignFeeTier(req, res, next) {
    try {
        const { tierId, minAmount, maxAmount, feePercentage, description } = req.body;

        const result = await fabricGateway.submitTransaction(
            ORG, CONTRACT, 'SetCampaignFeeTier',
            tierId, String(minAmount), String(maxAmount),
            String(feePercentage), description
        );

        res.json({ success: true, data: result });
    } catch (error) {
        next(error);
    }
}

/**
 * Set dispute fee tier
 */
export async function setDisputeFeeTier(req, res, next) {
    try {
        const { tierId, minAmount, maxAmount, feeAmount, description } = req.body;

        const result = await fabricGateway.submitTransaction(
            ORG, CONTRACT, 'SetDisputeFeeTier',
            tierId, String(minAmount), String(maxAmount),
            String(feeAmount), description
        );

        res.json({ success: true, data: result });
    } catch (error) {
        next(error);
    }
}

/**
 * Trigger fund release
 */
export async function triggerFundRelease(req, res, next) {
    try {
        const {
            releaseId, escrowAgreementId, agreementId, campaignId,
            milestoneId, startupId, amount, reason
        } = req.body;

        const result = await fabricGateway.submitTransaction(
            ORG, CONTRACT, 'TriggerFundRelease',
            releaseId, escrowAgreementId, agreementId, campaignId,
            milestoneId, startupId, String(amount), reason
        );

        res.json({ success: true, data: result });
    } catch (error) {
        next(error);
    }
}

/**
 * Create dispute
 */
export async function createDispute(req, res, next) {
    try {
        const {
            disputeId, complainantType, complainantId, respondentType,
            respondentId, disputeType, campaignId, agreementId,
            subject, description, amount
        } = req.body;

        const result = await fabricGateway.submitTransaction(
            ORG, CONTRACT, 'CreateDispute',
            disputeId, complainantType, complainantId, respondentType,
            respondentId, disputeType, campaignId, agreementId,
            subject, description, String(amount)
        );

        res.status(201).json({ success: true, data: result });
    } catch (error) {
        next(error);
    }
}

/**
 * Assign investigator to dispute
 */
export async function assignInvestigator(req, res, next) {
    try {
        const { disputeId } = req.params;
        const { investigatorId } = req.body;

        const result = await fabricGateway.submitTransaction(
            ORG, CONTRACT, 'AssignInvestigator',
            disputeId, investigatorId
        );

        res.json({ success: true, data: result });
    } catch (error) {
        next(error);
    }
}

/**
 * Get dispute
 */
export async function getDispute(req, res, next) {
    try {
        const { disputeId } = req.params;
        const result = await fabricGateway.evaluateTransaction(
            ORG, CONTRACT, 'GetDispute', disputeId
        );
        res.json({ success: true, data: result });
    } catch (error) {
        next(error);
    }
}
