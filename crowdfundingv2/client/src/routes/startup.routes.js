import express from 'express';
import * as startupController from '../controllers/startup.controller.js';

const router = express.Router();

// Campaign Management
router.post('/campaigns', startupController.createCampaign);
router.get('/campaigns/:campaignId', startupController.getCampaign);
router.get('/campaigns', startupController.getAllCampaigns);

// Validation Flow
router.post('/campaigns/:campaignId/submit-validation', startupController.submitForValidation);
router.get('/campaigns/:campaignId/validation-status', startupController.getValidationStatus);

// Platform Sharing
router.post('/campaigns/:campaignId/share-to-platform', startupController.shareToPlatform);
router.get('/campaigns/:campaignId/publish-notification', startupController.checkPublishNotification);

// Investment Management
router.get('/campaigns/:campaignId/investments', startupController.getCampaignInvestments);
router.get('/proposals', startupController.getProposals);
router.post('/proposals/:proposalId/respond', startupController.respondToProposal);

export default router;
