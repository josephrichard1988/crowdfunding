import express from 'express';
import * as investorController from '../controllers/investor.controller.js';

const router = express.Router();

// Campaign Discovery
router.get('/campaigns', investorController.getAvailableCampaigns);
router.get('/campaigns/:campaignId', investorController.viewCampaignDetails);
router.post('/campaigns/:campaignId/view', investorController.viewCampaign);

// Validation Insights
router.post('/validation-requests', investorController.requestValidationDetails);
router.get('/validation-requests/:requestId', investorController.getValidationResponse);

// Investment
router.post('/investments', investorController.makeInvestment);
router.get('/investments/:investmentId', investorController.getInvestment);
router.get('/investments', investorController.getMyInvestments);

// Proposals
router.post('/proposals', investorController.createProposal);
router.post('/proposals/:proposalId/respond', investorController.respondToCounterOffer);
router.get('/proposals/:proposalId', investorController.getProposal);

// Refunds
router.post('/refund-requests', investorController.requestRefund);

export default router;
