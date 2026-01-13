import express from 'express';
import * as validatorController from '../controllers/validator.controller.js';

const router = express.Router();

// Campaign Validation
router.get('/campaigns/:campaignId', validatorController.getCampaign);
router.get('/pending-validations', validatorController.getPendingValidations);
router.post('/validate/:campaignId', validatorController.validateCampaign);
router.post('/approve/:campaignId', validatorController.approveCampaign);

// Investor Validation Requests
router.get('/validation-requests', validatorController.getValidationRequests);
router.post('/validation-requests/:requestId/respond', validatorController.provideValidationDetails);

// Milestone Verification
router.post('/milestones/:milestoneId/verify', validatorController.verifyMilestone);

export default router;
