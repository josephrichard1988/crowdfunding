import { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { startupApi, startupMgmtApi } from '../services/api';
import { Plus, FileText, Share2, Bell, Loader2, RefreshCw, X, Rocket, Wallet, AlertCircle, LogIn, Coins, CreditCard, CheckCircle, Building2 } from 'lucide-react';

// Token fee constants (in CFT) - All paid by Startup
const FEES = {
    campaignCreation: 10,      // Create campaign
    validationSubmission: 50,  // Submit to validator
    platformPublishing: 50,    // Share to platform for publishing
    // Total journey: 110 CFT
};

export default function StartupDashboard() {
    const { user, isAuthenticated, updateWallet } = useAuth();
    const navigate = useNavigate();
    const [campaigns, setCampaigns] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [showCreateForm, setShowCreateForm] = useState(false);
    const [creating, setCreating] = useState(false);

    // Payment confirmation modal state
    const [showPaymentModal, setShowPaymentModal] = useState(false);
    const [pendingAction, setPendingAction] = useState(null); // { type: 'validation' | 'publishing', campaignId, fee }
    const [processing, setProcessing] = useState(false);

    // Startup management state
    const [startups, setStartups] = useState([]);
    const [selectedStartup, setSelectedStartup] = useState(null);
    const [showCreateStartupModal, setShowCreateStartupModal] = useState(false);
    const [newStartupName, setNewStartupName] = useState('');
    const [newStartupDesc, setNewStartupDesc] = useState('');

    const [formData, setFormData] = useState({
        // 22 Parameters for campaign creation
        projectName: '',
        description: '',
        category: 'Technology',
        goalAmount: '',
        currency: 'USD',
        deadline: '',
        hasRaised: false,
        hasGovGrants: false,
        incorpDate: new Date().toISOString().split('T')[0],
        projectStage: 'Idea',
        sector: 'Technology',
        tags: '',
        teamAvailable: false,
        investorCommitted: false,
        duration: 90,
        fundingDay: 1,
        fundingMonth: new Date().getMonth() + 1,
        fundingYear: new Date().getFullYear(),
        investmentRange: '10K-100K',
        documents: '',
    });

    // Get user wallet balance
    const cftBalance = user?.wallet?.cftBalance || 0;
    const canCreateCampaign = cftBalance >= FEES.campaignCreation;
    const canSubmitValidation = cftBalance >= FEES.validationSubmission;
    const canPublish = cftBalance >= FEES.platformPublishing;

    // Check if user has correct role
    const isStartupUser = isAuthenticated && user?.role === 'STARTUP';
    const isPreviewMode = !isAuthenticated || user?.role !== 'STARTUP';

    // Fetch startups and campaigns
    const fetchData = async () => {
        if (!isStartupUser) return;

        setLoading(true);
        try {
            // Fetch user's startups from chaincode
            let userStartups = [];
            try {
                const startupsRes = await startupMgmtApi.getStartupsByOwner(user.orgUserId);
                userStartups = startupsRes.data?.data || [];
            } catch (e) {
                console.warn('Failed to fetch startups:', e.message);
            }
            setStartups(userStartups);

            // Auto-select first startup if none selected
            if (userStartups.length > 0 && !selectedStartup) {
                setSelectedStartup(userStartups[0]);
            }

            // Fetch campaigns for selected startup (or all user campaigns)
            try {
                const startupIdToQuery = selectedStartup?.startupId || user?.orgUserId;
                const res = await startupApi.getAllCampaigns(startupIdToQuery);
                setCampaigns(res.data?.data || []);
            } catch (e) {
                console.warn('Failed to fetch campaigns:', e.message);
                setCampaigns([]);
            }
            setError(null);
        } catch (err) {
            setCampaigns([]);
            setStartups([]);
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        if (isStartupUser) {
            fetchData();
        } else {
            // In preview mode, don't show loading spinner
            setLoading(false);
        }
    }, [isStartupUser, user?.orgUserId, selectedStartup?.startupId]);

    // Create a new startup
    const handleCreateStartup = async () => {
        if (!newStartupName.trim()) {
            alert('Startup name is required');
            return;
        }
        try {
            const token = sessionStorage.getItem('token') || localStorage.getItem('token');
            const res = await startupMgmtApi.createStartup({
                name: newStartupName,
                description: newStartupDesc,
                ownerId: user.orgUserId,
                authToken: token  // For MongoDB sync
            });
            const createdStartup = res.data.data;
            setStartups(prev => [...prev, createdStartup]);
            setSelectedStartup(createdStartup);
            setShowCreateStartupModal(false);
            setNewStartupName('');
            setNewStartupDesc('');
            // Show success alert with startup details
            alert(`✅ Startup Created Successfully!\n\nName: ${createdStartup.name}\nStartup ID: ${createdStartup.startupId}\nDisplay ID: ${createdStartup.displayId}\n\nYou can now create campaigns under this startup.`);
        } catch (err) {
            console.error('Failed to create startup:', err);
            alert('Failed to create startup: ' + err.message);
        }
    };

    const handleInputChange = (e) => {
        const { name, value, type, checked } = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: type === 'checkbox' ? checked : value
        }));
    };

    const handleCreateCampaign = async (e) => {
        e.preventDefault();

        if (!selectedStartup) {
            alert('Please create or select a startup first');
            return;
        }

        setCreating(true);
        try {
            const tags = formData.tags.split(',').map(t => t.trim()).filter(Boolean);
            const documents = formData.documents.split(',').map(d => d.trim()).filter(Boolean);
            const token = sessionStorage.getItem('token') || localStorage.getItem('token');

            // Create campaign with all 22 parameters + auth token for MongoDB sync
            const result = await startupApi.createCampaign({
                startupId: selectedStartup.startupId,  // System-managed startup ID
                category: formData.category,
                deadline: formData.deadline,
                currency: formData.currency,
                hasRaised: formData.hasRaised,
                hasGovGrants: formData.hasGovGrants,
                incorpDate: formData.incorpDate,
                projectStage: formData.projectStage,
                sector: formData.sector,
                tags,
                teamAvailable: formData.teamAvailable,
                investorCommitted: formData.investorCommitted,
                duration: parseInt(formData.duration) || 90,
                fundingDay: parseInt(formData.fundingDay) || 1,
                fundingMonth: parseInt(formData.fundingMonth) || 1,
                fundingYear: parseInt(formData.fundingYear) || new Date().getFullYear(),
                goalAmount: parseFloat(formData.goalAmount) || 50000,
                investmentRange: formData.investmentRange,
                projectName: formData.projectName,
                description: formData.description,
                documents,
                authToken: token,  // For MongoDB sync
            });

            setShowCreateForm(false);
            setFormData({
                projectName: '',
                description: '',
                category: 'Technology',
                goalAmount: '',
                currency: 'USD',
                deadline: '',
                hasRaised: false,
                hasGovGrants: false,
                incorpDate: new Date().toISOString().split('T')[0],
                projectStage: 'Idea',
                sector: 'Technology',
                tags: '',
                teamAvailable: false,
                investorCommitted: false,
                duration: 90,
                fundingDay: 1,
                fundingMonth: new Date().getMonth() + 1,
                fundingYear: new Date().getFullYear(),
                investmentRange: '10K-100K',
                documents: '',
            });
            fetchData();

            // Show success with auto-generated ID
            if (result.data?.data?.displayId) {
                alert(`Campaign ${result.data.data.displayId} created successfully!`);
            }
        } catch (err) {
            console.error('Failed to create campaign:', err);
            alert('Failed to create campaign: ' + err.message);
        } finally {
            setCreating(false);
        }
    };

    // Initiate validation with payment confirmation
    const initiateValidationSubmit = (campaignId) => {
        if (!canSubmitValidation) {
            alert(`Insufficient balance. You need ${FEES.validationSubmission} CFT. Current: ${cftBalance} CFT`);
            return;
        }
        setPendingAction({ type: 'validation', campaignId, fee: FEES.validationSubmission });
        setShowPaymentModal(true);
    };

    // Initiate publishing with payment confirmation
    const initiateShareToPlatform = (campaignId, validationProofHash) => {
        if (!canPublish) {
            alert(`Insufficient balance. You need ${FEES.platformPublishing} CFT. Current: ${cftBalance} CFT`);
            return;
        }
        setPendingAction({ type: 'publishing', campaignId, validationProofHash, fee: FEES.platformPublishing });
        setShowPaymentModal(true);
    };

    // Process payment and execute action
    const confirmPayment = async () => {
        if (!pendingAction) return;
        setProcessing(true);

        try {
            // Deduct fee from wallet
            const newBalance = cftBalance - pendingAction.fee;
            await updateWallet({ cftBalance: newBalance });

            const token = sessionStorage.getItem('token') || localStorage.getItem('token');
            const campaign = campaigns.find(c => c.campaignId === pendingAction.campaignId);
            const projectName = campaign?.projectName || campaign?.project_name || '';

            // Execute the actual action with auth token for auto-allocation
            if (pendingAction.type === 'validation') {
                await startupApi.submitForValidation(pendingAction.campaignId, {
                    documents: [],
                    notes: 'Submitting for validation',
                    authToken: token,  // For auto-allocation
                    startupId: selectedStartup?.startupId,
                    projectName
                });
                alert(`Payment of ${pendingAction.fee} CFT confirmed! Campaign submitted for validation (auto-assigned to validator).`);
            } else if (pendingAction.type === 'publishing') {
                await startupApi.shareToPlatform(pendingAction.campaignId, {
                    validationProofHash: pendingAction.validationProofHash || '',
                    authToken: token,  // For auto-allocation
                    startupId: selectedStartup?.startupId,
                    projectName
                });
                alert(`Payment of ${pendingAction.fee} CFT confirmed! Campaign shared to platform (auto-assigned to publisher).`);
            }

            fetchData();
        } catch (err) {
            console.error('Failed:', err);
            alert('Failed: ' + err.message);
        } finally {
            setProcessing(false);
            setShowPaymentModal(false);
            setPendingAction(null);
        }
    };

    const cancelPayment = () => {
        setShowPaymentModal(false);
        setPendingAction(null);
    };

    const getStatusBadge = (status) => {
        const badges = {
            DRAFT: 'badge-info',
            SUBMITTED: 'badge-warning',
            APPROVED: 'badge-success',
            PUBLISHED: 'badge-success',
            REJECTED: 'badge-danger',
        };
        return badges[status] || 'badge-info';
    };

    return (
        <div className="space-y-6">
            {/* Payment Confirmation Modal */}
            {showPaymentModal && pendingAction && (
                <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
                    <div className="bg-white dark:bg-gray-800 rounded-xl shadow-2xl max-w-md w-full">
                        <div className="p-6 border-b border-gray-200 dark:border-gray-700">
                            <div className="flex items-center gap-3">
                                <div className="p-3 bg-accent-100 dark:bg-accent-900 rounded-full">
                                    <CreditCard className="text-accent-600" size={24} />
                                </div>
                                <div>
                                    <h2 className="text-xl font-bold text-gray-900 dark:text-white">
                                        Confirm Payment
                                    </h2>
                                    <p className="text-sm text-gray-500">
                                        {pendingAction.type === 'validation' ? 'Submit to Validator' : 'Share to Platform'}
                                    </p>
                                </div>
                            </div>
                        </div>
                        <div className="p-6 space-y-4">
                            <div className="bg-gray-50 dark:bg-gray-700 rounded-lg p-4">
                                <div className="flex justify-between items-center mb-2">
                                    <span className="text-gray-600 dark:text-gray-400">Fee</span>
                                    <span className="text-2xl font-bold text-accent-600">{pendingAction.fee} CFT</span>
                                </div>
                                <div className="flex justify-between items-center text-sm">
                                    <span className="text-gray-500">Your Balance</span>
                                    <span className="text-gray-700 dark:text-gray-300">{cftBalance} CFT</span>
                                </div>
                                <div className="flex justify-between items-center text-sm mt-1">
                                    <span className="text-gray-500">After Payment</span>
                                    <span className="font-medium text-green-600">{cftBalance - pendingAction.fee} CFT</span>
                                </div>
                            </div>
                            <p className="text-sm text-gray-500 text-center">
                                {pendingAction.type === 'validation'
                                    ? 'This fee goes to the validator for reviewing your campaign.'
                                    : 'This fee goes to the platform for publishing your campaign.'}
                            </p>
                        </div>
                        <div className="p-6 border-t border-gray-200 dark:border-gray-700 flex justify-end gap-3">
                            <button onClick={cancelPayment} className="btn btn-secondary">
                                Cancel
                            </button>
                            <button
                                onClick={confirmPayment}
                                disabled={processing}
                                className="btn btn-primary flex items-center gap-2"
                            >
                                {processing ? <Loader2 className="animate-spin" size={18} /> : <CheckCircle size={18} />}
                                {processing ? 'Processing...' : 'Confirm & Pay'}
                            </button>
                        </div>
                    </div>
                </div>
            )}

            {/* Preview Mode Banner */}
            {isPreviewMode && (
                <div className="bg-gradient-to-r from-blue-500 to-primary-600 text-white p-4 rounded-xl flex flex-col sm:flex-row items-center justify-between gap-4">
                    <div className="flex items-center gap-3">
                        <Rocket size={24} />
                        <div>
                            <h3 className="font-bold">Startup Dashboard Preview</h3>
                            <p className="text-sm opacity-90">Login as a startup to create and manage campaigns</p>
                        </div>
                    </div>
                    <Link to="/login" state={{ role: 'STARTUP' }} className="btn bg-white text-primary-700 hover:bg-gray-100 flex items-center gap-2">
                        <LogIn size={18} />
                        Login as Startup
                    </Link>
                </div>
            )}

            {/* Header */}
            <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
                <div>
                    <h1 className="text-3xl font-bold text-gray-900 dark:text-white">
                        Startup Dashboard
                    </h1>
                    <p className="text-gray-600 dark:text-gray-400 mt-1">
                        Manage your crowdfunding campaigns
                    </p>
                </div>
                <div className="flex gap-3 items-center">
                    {/* Wallet Balance (authenticated only) */}
                    {isStartupUser && (
                        <Link to="/wallet" className="flex items-center gap-2 px-4 py-2 bg-accent-100 dark:bg-accent-900/30 text-accent-700 dark:text-accent-300 rounded-lg hover:bg-accent-200">
                            <Coins size={18} />
                            <span className="font-medium">{cftBalance.toLocaleString()} CFT</span>
                        </Link>
                    )}
                    <button onClick={fetchData} className="btn btn-secondary flex items-center gap-2">
                        <RefreshCw size={18} />
                        Refresh
                    </button>
                </div>
            </div>

            {/* Create Campaign Modal */}
            {showCreateForm && (
                <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
                    <div className="bg-white dark:bg-gray-800 rounded-xl shadow-2xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
                        <div className="flex justify-between items-center p-6 border-b border-gray-200 dark:border-gray-700">
                            <h2 className="text-xl font-bold text-gray-900 dark:text-white flex items-center gap-2">
                                <Rocket className="text-primary-600" size={24} />
                                Create New Campaign
                            </h2>
                            <button onClick={() => setShowCreateForm(false)} className="text-gray-500 hover:text-gray-700">
                                <X size={24} />
                            </button>
                        </div>
                        <form onSubmit={handleCreateCampaign} className="p-6 space-y-4">
                            {/* Selected Startup Info - IDs are system-generated */}
                            <div className="p-3 bg-primary-50 dark:bg-primary-900/30 rounded-lg">
                                <p className="text-sm text-primary-700 dark:text-primary-300">
                                    <strong>Creating campaign for:</strong> {selectedStartup?.name || 'No startup selected'}
                                </p>
                                <p className="text-xs text-primary-600 dark:text-primary-400 mt-1">
                                    Startup ID: {selectedStartup?.startupId || 'N/A'}
                                    <span className="ml-2 text-gray-500">(Campaign ID will be auto-generated)</span>
                                </p>
                            </div>
                            <div>
                                <label className="label">Project Name *</label>
                                <input
                                    type="text"
                                    name="projectName"
                                    value={formData.projectName}
                                    onChange={handleInputChange}
                                    required
                                    placeholder="My Awesome Project"
                                    className="input"
                                />
                            </div>
                            <div>
                                <label className="label">Description *</label>
                                <textarea
                                    name="description"
                                    value={formData.description}
                                    onChange={handleInputChange}
                                    required
                                    rows={3}
                                    placeholder="Describe your project..."
                                    className="input"
                                />
                            </div>
                            <div className="grid md:grid-cols-3 gap-4">
                                <div>
                                    <label className="label">Category</label>
                                    <select
                                        name="category"
                                        value={formData.category}
                                        onChange={handleInputChange}
                                        className="input"
                                    >
                                        <option value="Technology">Technology</option>
                                        <option value="Healthcare">Healthcare</option>
                                        <option value="Finance">Finance</option>
                                        <option value="Education">Education</option>
                                        <option value="E-commerce">E-commerce</option>
                                        <option value="SaaS">SaaS</option>
                                    </select>
                                </div>
                                <div>
                                    <label className="label">Goal Amount *</label>
                                    <input
                                        type="number"
                                        name="goalAmount"
                                        value={formData.goalAmount}
                                        onChange={handleInputChange}
                                        required
                                        placeholder="50000"
                                        className="input"
                                    />
                                </div>
                                <div>
                                    <label className="label">Currency</label>
                                    <select
                                        name="currency"
                                        value={formData.currency}
                                        onChange={handleInputChange}
                                        className="input"
                                    >
                                        <option value="USD">USD</option>
                                        <option value="EUR">EUR</option>
                                        <option value="INR">INR</option>
                                    </select>
                                </div>
                            </div>
                            <div>
                                <label className="label">Deadline</label>
                                <input
                                    type="date"
                                    name="deadline"
                                    value={formData.deadline}
                                    onChange={handleInputChange}
                                    className="input"
                                />
                            </div>
                            <div>
                                <label className="label">Tags (comma-separated)</label>
                                <input
                                    type="text"
                                    name="tags"
                                    value={formData.tags}
                                    onChange={handleInputChange}
                                    placeholder="IoT, AI, Mobile"
                                    className="input"
                                />
                            </div>
                            <div>
                                <label className="label">Documents (comma-separated)</label>
                                <input
                                    type="text"
                                    name="documents"
                                    value={formData.documents}
                                    onChange={handleInputChange}
                                    placeholder="pitch_deck.pdf, financials.xlsx"
                                    className="input"
                                />
                            </div>
                            <div className="flex justify-end gap-3 pt-4">
                                <button type="button" onClick={() => setShowCreateForm(false)} className="btn btn-secondary">
                                    Cancel
                                </button>
                                <button type="submit" disabled={creating} className="btn btn-primary flex items-center gap-2">
                                    {creating ? <Loader2 className="animate-spin" size={18} /> : <Plus size={18} />}
                                    {creating ? 'Creating...' : 'Create Campaign'}
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}

            {/* Quick Actions */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div className="card flex items-center gap-4">
                    <div className="p-3 bg-blue-100 dark:bg-blue-900 rounded-lg">
                        <FileText className="text-blue-600 dark:text-blue-400" size={24} />
                    </div>
                    <div>
                        <h3 className="font-semibold text-gray-900 dark:text-white">Draft Campaigns</h3>
                        <p className="text-2xl font-bold text-blue-600">{campaigns.filter(c => c.status === 'DRAFT').length}</p>
                    </div>
                </div>
                <div className="card flex items-center gap-4">
                    <div className="p-3 bg-yellow-100 dark:bg-yellow-900 rounded-lg">
                        <Share2 className="text-yellow-600 dark:text-yellow-400" size={24} />
                    </div>
                    <div>
                        <h3 className="font-semibold text-gray-900 dark:text-white">Pending Validation</h3>
                        <p className="text-2xl font-bold text-yellow-600">{campaigns.filter(c => c.validationStatus === 'PENDING_VALIDATION').length}</p>
                    </div>
                </div>
                <div className="card flex items-center gap-4">
                    <div className="p-3 bg-green-100 dark:bg-green-900 rounded-lg">
                        <Bell className="text-green-600 dark:text-green-400" size={24} />
                    </div>
                    <div>
                        <h3 className="font-semibold text-gray-900 dark:text-white">Published</h3>
                        <p className="text-2xl font-bold text-green-600">{campaigns.filter(c => c.status === 'PUBLISHED').length}</p>
                    </div>
                </div>
            </div>

            {/* My Startups Section */}
            {isStartupUser && (
                <div className="card">
                    <div className="flex justify-between items-center mb-4">
                        <h2 className="text-xl font-bold text-gray-900 dark:text-white flex items-center gap-2">
                            <Building2 className="text-primary-600" />
                            My Startups
                        </h2>
                        <button
                            onClick={() => setShowCreateStartupModal(true)}
                            className="btn btn-primary flex items-center gap-2"
                        >
                            <Plus size={18} />
                            New Startup
                        </button>
                    </div>

                    {startups.length === 0 ? (
                        <div className="text-center py-12 bg-gray-50 dark:bg-gray-800 rounded-lg">
                            <Building2 size={48} className="mx-auto mb-4 text-gray-300" />
                            <h3 className="text-lg font-medium text-gray-700 dark:text-gray-300 mb-2">No startups yet</h3>
                            <p className="text-gray-500 mb-4">Create your first startup to start launching campaigns</p>
                            <button
                                onClick={() => setShowCreateStartupModal(true)}
                                className="btn btn-primary"
                            >
                                <Plus size={18} className="mr-2" />
                                Create First Startup
                            </button>
                        </div>
                    ) : (
                        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                            {startups.map((startup) => (
                                <div
                                    key={startup.startupId}
                                    onClick={() => navigate(`/startup/${startup.startupId}`)}
                                    className="p-4 rounded-lg border-2 cursor-pointer transition-all border-gray-200 dark:border-gray-700 hover:border-primary-300 hover:shadow-lg bg-white dark:bg-gray-800"
                                >
                                    <div className="flex items-start justify-between">
                                        <div>
                                            <h3 className="font-bold text-gray-900 dark:text-white">{startup.name}</h3>
                                            <p className="text-xs text-gray-500 mt-1">ID: {startup.displayId || startup.startupId}</p>
                                        </div>
                                    </div>
                                    {startup.description && (
                                        <p className="text-sm text-gray-600 dark:text-gray-400 mt-2 line-clamp-2">
                                            {startup.description}
                                        </p>
                                    )}
                                    <div className="flex items-center justify-between mt-3 pt-3 border-t border-gray-200 dark:border-gray-700">
                                        <span className="text-xs text-gray-500">
                                            {(startup.campaignIds || startup.CampaignIDs || startup.campaign_ids || []).length} campaigns
                                        </span>
                                        <span className="text-xs text-primary-600 font-medium">View Campaigns →</span>
                                    </div>
                                </div>
                            ))}
                        </div>
                    )}
                </div>
            )}

            {/* Create Startup Modal */}
            {showCreateStartupModal && (
                <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
                    <div className="bg-white dark:bg-gray-800 rounded-xl shadow-2xl max-w-md w-full">
                        <div className="flex justify-between items-center p-6 border-b border-gray-200 dark:border-gray-700">
                            <h2 className="text-xl font-bold text-gray-900 dark:text-white flex items-center gap-2">
                                <Building2 className="text-primary-600" size={24} />
                                Create New Startup
                            </h2>
                            <button onClick={() => setShowCreateStartupModal(false)} className="text-gray-500 hover:text-gray-700">
                                <X size={24} />
                            </button>
                        </div>
                        <div className="p-6 space-y-4">
                            <div>
                                <label className="label">Startup Name *</label>
                                <input
                                    type="text"
                                    value={newStartupName}
                                    onChange={(e) => setNewStartupName(e.target.value)}
                                    placeholder="Enter startup name"
                                    className="input w-full"
                                    required
                                />
                            </div>
                            <div>
                                <label className="label">Description</label>
                                <textarea
                                    value={newStartupDesc}
                                    onChange={(e) => setNewStartupDesc(e.target.value)}
                                    placeholder="Describe your startup..."
                                    rows={3}
                                    className="input w-full"
                                />
                            </div>
                            <p className="text-xs text-gray-500">
                                A unique Startup ID will be auto-generated for your startup.
                            </p>
                        </div>
                        <div className="p-6 border-t border-gray-200 dark:border-gray-700 flex justify-end gap-3">
                            <button onClick={() => setShowCreateStartupModal(false)} className="btn btn-secondary">
                                Cancel
                            </button>
                            <button
                                onClick={handleCreateStartup}
                                className="btn btn-primary flex items-center gap-2"
                            >
                                <Plus size={18} />
                                Create Startup
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}
