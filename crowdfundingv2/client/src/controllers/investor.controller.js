import fabricGateway from '../services/fabricGateway.js';

const ORG = 'investor';
const CONTRACT = 'InvestorContract';

/**
 * Get all available published campaigns (list view)
 */
export async function getAvailableCampaigns(req, res, next) {
    try {
        const result = await fabricGateway.evaluateTransaction(
            ORG, CONTRACT, 'GetAvailableCampaigns'
        );
        res.json({ success: true, data: result });
    } catch (error) {
        next(error);
    }
}

/**
 * View campaign details (detail view)
 */
export async function viewCampaignDetails(req, res, next) {
    try {
        const { campaignId } = req.params;
        const result = await fabricGateway.evaluateTransaction(
            ORG, CONTRACT, 'ViewCampaignDetails', campaignId
        );
        res.json({ success: true, data: result });
    } catch (error) {
        next(error);
    }
}

/**
 * View and log campaign (stores in private collection)
 */
export async function viewCampaign(req, res, next) {
    try {
        const { campaignId } = req.params;
        const { viewRecordId, investorId } = req.body;

        const result = await fabricGateway.submitTransaction(
            ORG, CONTRACT, 'ViewCampaign',
            viewRecordId, campaignId, investorId
        );

        res.json({ success: true, data: result });
    } catch (error) {
        next(error);
    }
}

/**
 * Request validation details from validator
 */
export async function requestValidationDetails(req, res, next) {
    try {
        const { requestId, campaignId, investorId } = req.body;

        const result = await fabricGateway.submitTransaction(
            ORG, CONTRACT, 'RequestValidationDetails',
            requestId, campaignId, investorId
        );

        res.json({ success: true, data: result });
    } catch (error) {
        next(error);
    }
}

/**
 * Get validation response
 */
export async function getValidationResponse(req, res, next) {
    try {
        const { requestId } = req.params;
        const result = await fabricGateway.evaluateTransaction(
            ORG, CONTRACT, 'GetValidationResponse', requestId
        );
        res.json({ success: true, data: result });
    } catch (error) {
        next(error);
    }
}

/**
 * Make investment
 */
export async function makeInvestment(req, res, next) {
    try {
        const { investmentId, campaignId, investorId, amount, currency } = req.body;

        const result = await fabricGateway.submitTransaction(
            ORG, CONTRACT, 'MakeInvestment',
            investmentId, campaignId, investorId, String(amount), currency
        );

        res.status(201).json({ success: true, data: result });
    } catch (error) {
        next(error);
    }
}

/**
 * Get investment by ID
 */
export async function getInvestment(req, res, next) {
    try {
        const { investmentId } = req.params;
        const result = await fabricGateway.evaluateTransaction(
            ORG, CONTRACT, 'GetInvestment', investmentId
        );
        res.json({ success: true, data: result });
    } catch (error) {
        next(error);
    }
}

/**
 * Get all my investments
 */
export async function getMyInvestments(req, res, next) {
    try {
        const result = await fabricGateway.evaluateTransaction(
            ORG, CONTRACT, 'GetMyInvestments'
        );
        res.json({ success: true, data: result });
    } catch (error) {
        next(error);
    }
}

/**
 * Create investment proposal
 */
export async function createProposal(req, res, next) {
    try {
        const {
            proposalId, campaignId, investorId, startupId,
            investmentAmount, currency, equity, duration,
            milestones, proposedTerms
        } = req.body;

        const result = await fabricGateway.submitTransaction(
            ORG, CONTRACT, 'CreateInvestmentProposal',
            proposalId, campaignId, investorId, startupId,
            String(investmentAmount), currency, equity, duration,
            JSON.stringify(milestones || []), proposedTerms
        );

        res.status(201).json({ success: true, data: result });
    } catch (error) {
        next(error);
    }
}

/**
 * Respond to counter offer from startup
 */
export async function respondToCounterOffer(req, res, next) {
    try {
        const { proposalId } = req.params;
        const { action, counterAmount, counterTerms } = req.body;

        const result = await fabricGateway.submitTransaction(
            ORG, CONTRACT, 'RespondToCounterOffer',
            proposalId, action, String(counterAmount || 0), counterTerms || ''
        );

        res.json({ success: true, data: result });
    } catch (error) {
        next(error);
    }
}

/**
 * Get proposal by ID
 */
export async function getProposal(req, res, next) {
    try {
        const { proposalId } = req.params;
        const result = await fabricGateway.evaluateTransaction(
            ORG, CONTRACT, 'GetProposal', proposalId
        );
        res.json({ success: true, data: result });
    } catch (error) {
        next(error);
    }
}

/**
 * Request refund
 */
export async function requestRefund(req, res, next) {
    try {
        const { refundId, investmentId, investorId, reason } = req.body;

        const result = await fabricGateway.submitTransaction(
            ORG, CONTRACT, 'RequestRefund',
            refundId, investmentId, investorId, reason
        );

        res.json({ success: true, data: result });
    } catch (error) {
        next(error);
    }
}
