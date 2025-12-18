import axios from 'axios';

const API_BASE = '/api';

const api = axios.create({
    baseURL: API_BASE,
    headers: {
        'Content-Type': 'application/json',
    },
});

// Startup API
export const startupApi = {
    createCampaign: (data) => api.post('/startup/campaigns', data),
    getCampaign: (campaignId) => api.get(`/startup/campaigns/${campaignId}`),
    getAllCampaigns: () => api.get('/startup/campaigns'),
    submitForValidation: (campaignId, data) => api.post(`/startup/campaigns/${campaignId}/submit-validation`, data),
    shareToPlatform: (campaignId, data) => api.post(`/startup/campaigns/${campaignId}/share-to-platform`, data),
    checkPublishNotification: (campaignId) => api.get(`/startup/campaigns/${campaignId}/publish-notification`),
    getProposals: () => api.get('/startup/proposals'),
    respondToProposal: (proposalId, data) => api.post(`/startup/proposals/${proposalId}/respond`, data),
};

// Validator API
export const validatorApi = {
    getCampaign: (campaignId) => api.get(`/validator/campaigns/${campaignId}`),
    getPendingValidations: () => api.get('/validator/pending-validations'),
    validateCampaign: (campaignId, data) => api.post(`/validator/validate/${campaignId}`, data),
    approveCampaign: (campaignId, data) => api.post(`/validator/approve/${campaignId}`, data),
    getValidationRequests: () => api.get('/validator/validation-requests'),
    provideValidationDetails: (requestId, data) => api.post(`/validator/validation-requests/${requestId}/respond`, data),
};

// Platform API
export const platformApi = {
    getSharedCampaign: (campaignId) => api.get(`/platform/shared-campaigns/${campaignId}`),
    getAllSharedCampaigns: () => api.get('/platform/shared-campaigns'),
    publishCampaign: (campaignId, data) => api.post(`/platform/publish/${campaignId}`, data),
    getPublishedCampaign: (campaignId) => api.get(`/platform/published-campaigns/${campaignId}`),
    createWallet: (data) => api.post('/platform/wallets', data),
    triggerFundRelease: (data) => api.post('/platform/release-funds', data),
    createDispute: (data) => api.post('/platform/disputes', data),
};

// Investor API
export const investorApi = {
    getAvailableCampaigns: () => api.get('/investor/campaigns'),
    viewCampaignDetails: (campaignId) => api.get(`/investor/campaigns/${campaignId}`),
    viewCampaign: (campaignId, data) => api.post(`/investor/view/${campaignId}`, data),
    requestValidationDetails: (campaignId, data) => api.post(`/investor/request-validation/${campaignId}`, data),
    getViewedCampaigns: () => api.get('/investor/viewed-campaigns'),
    makeInvestment: (data) => api.post('/investor/investments', data),
    getMyInvestments: () => api.get('/investor/my-investments'),
    createProposal: (data) => api.post('/investor/proposals', data),
};

export default api;
