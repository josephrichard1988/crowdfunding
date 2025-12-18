import fabricGateway from '../services/fabricGateway.js';

const ORG = 'startup';
const CONTRACT = 'StartupContract';

/**
 * Create a new campaign with 22 parameters
 */
export async function createCampaign(req, res, next) {
    try {
        const {
            campaignId, startupId, category, deadline, currency,
            hasRaised, hasGovGrants, incorpDate, projectStage, sector,
            tags, teamAvailable, investorCommitted, duration,
            fundingDay, fundingMonth, fundingYear, goalAmount,
            investmentRange, projectName, description, documents
        } = req.body;

        const result = await fabricGateway.submitTransaction(
            ORG, CONTRACT, 'CreateCampaign',
            campaignId, startupId, category, deadline, currency,
            String(hasRaised), String(hasGovGrants), incorpDate, projectStage, sector,
            JSON.stringify(tags), String(teamAvailable), String(investorCommitted),
            String(duration), String(fundingDay), String(fundingMonth), String(fundingYear),
            String(goalAmount), investmentRange, projectName, description,
            JSON.stringify(documents)
        );

        res.status(201).json({ success: true, data: result });
    } catch (error) {
        next(error);
    }
}

/**
 * Get campaign by ID
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
 * Get all campaigns for the startup
 */
export async function getAllCampaigns(req, res, next) {
    try {
        const result = await fabricGateway.evaluateTransaction(
            ORG, CONTRACT, 'GetAllCampaigns'
        );
        res.json({ success: true, data: result });
    } catch (error) {
        next(error);
    }
}

/**
 * Submit campaign for validation
 */
export async function submitForValidation(req, res, next) {
    try {
        const { campaignId } = req.params;
        const { documents, notes } = req.body;

        const result = await fabricGateway.submitTransaction(
            ORG, CONTRACT, 'SubmitForValidation',
            campaignId, JSON.stringify(documents || []), notes || ''
        );

        res.json({ success: true, data: result });
    } catch (error) {
        next(error);
    }
}

/**
 * Get validation status
 */
export async function getValidationStatus(req, res, next) {
    try {
        const { campaignId } = req.params;
        const result = await fabricGateway.evaluateTransaction(
            ORG, CONTRACT, 'GetValidationStatus', campaignId
        );
        res.json({ success: true, data: result });
    } catch (error) {
        next(error);
    }
}

/**
 * Share campaign to platform
 */
export async function shareToPlatform(req, res, next) {
    try {
        const { campaignId } = req.params;
        const { validationHash } = req.body;

        const result = await fabricGateway.submitTransaction(
            ORG, CONTRACT, 'ShareCampaignToPlatform',
            campaignId, validationHash
        );

        res.json({ success: true, data: result });
    } catch (error) {
        next(error);
    }
}

/**
 * Check publish notification from platform
 */
export async function checkPublishNotification(req, res, next) {
    try {
        const { campaignId } = req.params;
        const result = await fabricGateway.evaluateTransaction(
            ORG, CONTRACT, 'CheckPublishNotification', campaignId
        );
        res.json({ success: true, data: result });
    } catch (error) {
        next(error);
    }
}

/**
 * Get investments for a campaign
 */
export async function getCampaignInvestments(req, res, next) {
    try {
        const { campaignId } = req.params;
        const result = await fabricGateway.evaluateTransaction(
            ORG, CONTRACT, 'GetCampaignInvestments', campaignId
        );
        res.json({ success: true, data: result });
    } catch (error) {
        next(error);
    }
}

/**
 * Get all proposals for startup
 */
export async function getProposals(req, res, next) {
    try {
        const result = await fabricGateway.evaluateTransaction(
            ORG, CONTRACT, 'GetProposals'
        );
        res.json({ success: true, data: result });
    } catch (error) {
        next(error);
    }
}

/**
 * Respond to investment proposal
 */
export async function respondToProposal(req, res, next) {
    try {
        const { proposalId } = req.params;
        const { action, counterAmount, counterTerms } = req.body;

        const result = await fabricGateway.submitTransaction(
            ORG, CONTRACT, 'RespondToProposal',
            proposalId, action, String(counterAmount || 0), counterTerms || ''
        );

        res.json({ success: true, data: result });
    } catch (error) {
        next(error);
    }
}
