import axios from 'axios';

const API_BASE = '/api';
const AUTH_API_BASE = 'http://localhost:3001/api/auth';

const api = axios.create({
    baseURL: API_BASE,
    headers: {
        'Content-Type': 'application/json',
    },
});

// Auth API (for startup management, sync, and queue)
const authApi = axios.create({
    baseURL: AUTH_API_BASE,
    headers: {
        'Content-Type': 'application/json',
    },
});

// Add auth token interceptor
authApi.interceptors.request.use((config) => {
    const token = sessionStorage.getItem('token') || localStorage.getItem('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

// Startup Management API (for STARTUP role)
// Creates startups on chaincode via network API and syncs to MongoDB
export const startupMgmtApi = {
    // Create startup on chaincode (via network API)
    createStartup: (data) => api.post('/startup/startups', data),
    // Get startup by ID from chaincode
    getStartup: (startupId) => api.get(`/startup/startups/${startupId}`),
    // Get all startups for an owner from chaincode
    getStartupsByOwner: (ownerId) => api.get(`/startup/startups/owner/${ownerId}`),
    // Sync campaign to MongoDB
    syncCampaign: (data) => authApi.post('/sync/campaign', data),
    // Sync startup to chaincode (recreate in chaincode from MongoDB data)
    syncToChaincode: (startupId, data) => api.post(`/startup/startups/${startupId}/sync-to-chaincode`, data),

    // Deletion APIs
    getCampaignDeletionFee: (campaignId) => api.get(`/startup/campaigns/${campaignId}/deletion-fee`),
    deleteCampaign: (campaignId, reason, startupId) => {
        const token = sessionStorage.getItem('token') || localStorage.getItem('token');
        return api.delete(`/startup/campaigns/${campaignId}`, {
            data: { reason, authToken: token, startupId }
        });
    },
    getStartupDeletionFee: (startupId) => api.get(`/startup/startups/${startupId}/deletion-fee`),
    deleteStartup: (startupId, reason) => {
        const token = sessionStorage.getItem('token') || localStorage.getItem('token');
        return api.delete(`/startup/startups/${startupId}`, {
            data: { reason, authToken: token, startupId }
        });
    },
    getDeletionRecord: (deletionId) => api.get(`/startup/deletions/${deletionId}`),
};

// Queue Management API (for VALIDATOR/PLATFORM roles)
export const queueApi = {
    getQueue: () => authApi.get('/queue'),
    getNextAllocation: (role) => authApi.get(`/allocation/next?role=${role}`),
    assign: (data) => authApi.post('/allocation/assign', data),
    complete: (data) => authApi.post('/allocation/complete', data),
};

// Startup API
export const startupApi = {
    createCampaign: (data) => api.post('/startup/campaigns', data),
    getCampaign: (campaignId) => api.get(`/startup/campaigns/${campaignId}`),
    getAllCampaigns: (startupId) => api.get(`/startup/campaigns${startupId ? `?startupId=${startupId}` : ''}`),
    submitForValidation: (campaignId, data) => api.post(`/startup/campaigns/${campaignId}/submit-validation`, data),
    shareToPlatform: (campaignId, data) => api.post(`/startup/campaigns/${campaignId}/share-to-platform`, data),
    checkPublishNotification: (campaignId) => api.get(`/startup/campaigns/${campaignId}/publish-notification`),
    getProposals: () => api.get('/startup/proposals'),
    respondToProposal: (proposalId, data) => api.post(`/startup/proposals/${proposalId}/respond`, data),
};

// Validator API
export const validatorApi = {
    getCampaign: (campaignId) => api.get(`/validator/campaigns/${campaignId}`),
    getPendingValidations: (validatorId) => api.get(`/validator/pending-validations${validatorId ? `?validatorId=${validatorId}` : ''}`),
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
    getViewedCampaigns: (investorId) => api.get(`/investor/viewed-campaigns${investorId ? `?investorId=${investorId}` : ''}`),
    makeInvestment: (data) => api.post('/investor/investments', data),
    getMyInvestments: (investorId) => api.get(`/investor/my-investments${investorId ? `?investorId=${investorId}` : ''}`),
    createProposal: (data) => api.post('/investor/proposals', data),
};

export { authApi };
export default api;
