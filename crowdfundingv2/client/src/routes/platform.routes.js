import express from 'express';
import * as platformController from '../controllers/platform.controller.js';

const router = express.Router();

// Campaign Management
router.get('/shared-campaigns/:campaignId', platformController.getSharedCampaign);
router.get('/shared-campaigns', platformController.getAllSharedCampaigns);
router.post('/publish/:campaignId', platformController.publishCampaign);
router.get('/published-campaigns/:campaignId', platformController.getPublishedCampaign);

// Wallet & Fee Management
router.post('/wallets', platformController.createWallet);
router.get('/wallets/:walletId', platformController.getWallet);
router.post('/fee-tiers/campaign', platformController.setCampaignFeeTier);
router.post('/fee-tiers/dispute', platformController.setDisputeFeeTier);

// Fund Release
router.post('/release-funds', platformController.triggerFundRelease);

// Disputes
router.post('/disputes', platformController.createDispute);
router.post('/disputes/:disputeId/assign', platformController.assignInvestigator);
router.get('/disputes/:disputeId', platformController.getDispute);

export default router;
