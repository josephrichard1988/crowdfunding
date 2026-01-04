import { useState, useEffect } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { startupApi, startupMgmtApi } from '../services/api';
import { Plus, ArrowLeft, Rocket, Loader2, X, RefreshCw, Coins, CreditCard, CheckCircle, FileText, Trash2, AlertTriangle } from 'lucide-react';

// Token fee constants
const FEES = {
    campaignCreation: 10,
    validationSubmission: 50,
    platformPublishing: 50,
};

export default function StartupDetail() {
    const { startupId } = useParams();
    const navigate = useNavigate();
    const { user, isAuthenticated, updateWallet, token } = useAuth();

    const [startup, setStartup] = useState(null);
    const [campaigns, setCampaigns] = useState([]);
    const [loading, setLoading] = useState(true);
    const [showCreateForm, setShowCreateForm] = useState(false);
    const [creating, setCreating] = useState(false);

    // Payment modal state
    const [showPaymentModal, setShowPaymentModal] = useState(false);
    const [pendingAction, setPendingAction] = useState(null);
    const [processing, setProcessing] = useState(false);

    // Delete modal state
    const [showDeleteModal, setShowDeleteModal] = useState(false);
    const [deleteTarget, setDeleteTarget] = useState(null); // { type: 'campaign' | 'startup', id, name }
    const [deletionFee, setDeletionFee] = useState(null);
    const [deleteReason, setDeleteReason] = useState('');
    const [deleting, setDeleting] = useState(false);

    const [formData, setFormData] = useState({
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

    const cftBalance = user?.wallet?.cftBalance || 0;
    const canCreateCampaign = cftBalance >= FEES.campaignCreation;
    const isStartupUser = isAuthenticated && user?.role === 'STARTUP';

    const fetchData = async () => {
        setLoading(true);
        try {
            // Fetch startup details
            try {
                const startupRes = await startupMgmtApi.getStartup(startupId);
                setStartup(startupRes.data?.data || null);
            } catch (e) {
                console.warn('Failed to fetch startup:', e.message);
            }

            // Fetch campaigns for this startup
            try {
                const campaignsRes = await startupApi.getAllCampaigns(startupId);
                setCampaigns(campaignsRes.data?.data || []);
            } catch (e) {
                console.warn('Failed to fetch campaigns:', e.message);
                setCampaigns([]);
            }
        } catch (err) {
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        if (startupId) {
            fetchData();
        }
    }, [startupId]);

    const handleInputChange = (e) => {
        const { name, value, type, checked } = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: type === 'checkbox' ? checked : value
        }));
    };

    const handleCreateCampaign = async (e) => {
        e.preventDefault();
        if (!canCreateCampaign) {
            alert(`Insufficient balance. You need ${FEES.campaignCreation} CFT.`);
            return;
        }

        setCreating(true);
        try {
            const authToken = sessionStorage.getItem('token') || localStorage.getItem('token');
            await startupApi.createCampaign({
                startupId: startupId,
                projectName: formData.projectName,
                description: formData.description,
                category: formData.category,
                goalAmount: parseFloat(formData.goalAmount) || 0,
                currency: formData.currency,
                deadline: formData.deadline,
                hasRaised: formData.hasRaised,
                hasGovGrants: formData.hasGovGrants,
                incorpDate: formData.incorpDate,
                projectStage: formData.projectStage,
                sector: formData.sector,
                tags: formData.tags.split(',').map(t => t.trim()).filter(Boolean),
                teamAvailable: formData.teamAvailable,
                investorCommitted: formData.investorCommitted,
                duration: parseInt(formData.duration) || 90,
                fundingDay: parseInt(formData.fundingDay) || 1,
                fundingMonth: parseInt(formData.fundingMonth) || 1,
                fundingYear: parseInt(formData.fundingYear) || new Date().getFullYear(),
                investmentRange: formData.investmentRange,
                documents: formData.documents.split(',').map(d => d.trim()).filter(Boolean),
                authToken
            });

            // Deduct fee
            if (updateWallet) {
                updateWallet({ cftBalance: cftBalance - FEES.campaignCreation });
            }

            alert('Campaign created successfully!');
            setShowCreateForm(false);
            setFormData({
                projectName: '', description: '', category: 'Technology', goalAmount: '',
                currency: 'USD', deadline: '', hasRaised: false, hasGovGrants: false,
                incorpDate: new Date().toISOString().split('T')[0], projectStage: 'Idea',
                sector: 'Technology', tags: '', teamAvailable: false, investorCommitted: false,
                duration: 90, fundingDay: 1, fundingMonth: new Date().getMonth() + 1,
                fundingYear: new Date().getFullYear(), investmentRange: '10K-100K', documents: '',
            });
            fetchData();
        } catch (err) {
            console.error(err);
            alert('Failed to create campaign: ' + err.message);
        } finally {
            setCreating(false);
        }
    };

    const initiateValidationSubmit = (campaignId) => {
        if (cftBalance < FEES.validationSubmission) {
            alert(`Insufficient balance. You need ${FEES.validationSubmission} CFT.`);
            return;
        }
        setPendingAction({ type: 'validation', campaignId, fee: FEES.validationSubmission });
        setShowPaymentModal(true);
    };

    const initiateShareToPlatform = (campaignId, validationProofHash) => {
        if (cftBalance < FEES.platformPublishing) {
            alert(`Insufficient balance. You need ${FEES.platformPublishing} CFT.`);
            return;
        }
        setPendingAction({ type: 'publishing', campaignId, validationProofHash, fee: FEES.platformPublishing });
        setShowPaymentModal(true);
    };

    const confirmPayment = async () => {
        if (!pendingAction) return;
        setProcessing(true);
        try {
            const authToken = sessionStorage.getItem('token') || localStorage.getItem('token');
            if (pendingAction.type === 'validation') {
                await startupApi.submitForValidation(pendingAction.campaignId, {
                    documents: ['validation_request.pdf'],
                    notes: 'Submitted for validation',
                    authToken,
                    startupId,
                    projectName: campaigns.find(c => c.campaignId === pendingAction.campaignId)?.projectName
                });
            } else {
                await startupApi.shareToPlatform(pendingAction.campaignId, {
                    validationProofHash: pendingAction.validationProofHash,
                    authToken,
                    startupId,
                    projectName: campaigns.find(c => c.campaignId === pendingAction.campaignId)?.projectName
                });
            }
            if (updateWallet) {
                updateWallet({ cftBalance: cftBalance - pendingAction.fee });
            }
            alert(`${pendingAction.type === 'validation' ? 'Submitted for validation' : 'Shared to platform'} successfully!`);
            fetchData();
        } catch (err) {
            alert('Action failed: ' + err.message);
        } finally {
            setProcessing(false);
            setShowPaymentModal(false);
            setPendingAction(null);
        }
    };

    // Delete handlers
    const initiateDeleteCampaign = async (campaign) => {
        try {
            setDeleteTarget({
                type: 'campaign',
                id: campaign.campaignId,
                name: campaign.projectName || campaign.project_name,
                fundsRaised: campaign.fundsRaisedAmount || campaign.funds_raised_amount || 0
            });
            // Fetch deletion fee
            const res = await startupMgmtApi.getCampaignDeletionFee(campaign.campaignId);
            setDeletionFee(res.data?.data || { feeAmount: 100, isFixedFee: true });
            setShowDeleteModal(true);
        } catch (err) {
            // If can't get fee, show fixed 100 CFT
            setDeletionFee({ feeAmount: 100, isFixedFee: true, fundsRaised: 0 });
            setShowDeleteModal(true);
        }
    };



    const initiateDeleteStartup = async () => {
        try {
            setDeleteTarget({
                type: 'startup',
                id: startupId,
                name: startup.name,
                fundsRaised: 0 // Will clearly see aggregated fee in modal
            });
            const res = await startupMgmtApi.getStartupDeletionFee(startupId);
            setDeletionFee(res.data?.data || { feeAmount: 100 });
            setShowDeleteModal(true);
        } catch (err) {
            setDeletionFee({ feeAmount: 100, isFixedFee: false });
            setShowDeleteModal(true);
        }
    };

    const confirmDelete = async () => {
        if (!deleteTarget) return;

        // Check balance
        const requiredFee = deletionFee?.feeAmount || 100;
        if (cftBalance < requiredFee) {
            alert(`⚠️ Insufficient Balance!\n\nYou have: ${cftBalance} CFT\nRequired: ${requiredFee} CFT\n\nPlease top up your wallet to perform this deletion.`);
            return;
        }

        setDeleting(true);
        try {
            if (deleteTarget.type === 'campaign') {
                await startupMgmtApi.deleteCampaign(deleteTarget.id, deleteReason || 'User requested deletion');
                alert(`Campaign "${deleteTarget.name}" deleted. Fee charged: ${deletionFee?.feeAmount || 100} CFT`);
            }
            if (deleteTarget.type === 'startup') {
                await startupMgmtApi.deleteStartup(deleteTarget.id, deleteReason || 'User requested deletion');
                alert(`Startup "${deleteTarget.name}" and all campaigns deleted. Fee charged: ${deletionFee?.feeAmount || 100} CFT`);
                navigate('/startup'); // Redirect to dashboard
                return;
            }
            setShowDeleteModal(false);
            setDeleteTarget(null);
            setDeletionFee(null);
            setDeleteReason('');
            fetchData(); // Refresh list
        } catch (err) {
            alert('Delete failed: ' + err.message);
        } finally {
            setDeleting(false);
        }
    };

    const cancelDelete = () => {
        setShowDeleteModal(false);
        setDeleteTarget(null);
        setDeletionFee(null);
        setDeleteReason('');
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

    if (loading) {
        return (
            <div className="flex justify-center items-center py-20">
                <Loader2 className="animate-spin text-primary-600" size={48} />
            </div>
        );
    }

    return (
        <div className="space-y-6">
            {/* Header with back button */}
            <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
                <div className="flex items-center gap-4">
                    <button onClick={() => navigate('/startup')} className="p-2 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-lg">
                        <ArrowLeft size={24} />
                    </button>
                    <div>
                        <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
                            {startup?.name || 'Startup'}
                        </h1>
                        <p className="text-sm text-gray-500">ID: {startup?.displayId || startupId}</p>
                    </div>
                </div>
                <div className="flex gap-3 items-center">
                    {isStartupUser && (
                        <Link to="/wallet" className="flex items-center gap-2 px-4 py-2 bg-primary-100 dark:bg-primary-900/30 text-primary-700 dark:text-primary-300 rounded-lg">
                            <Coins size={18} />
                            <span className="font-medium">{cftBalance.toLocaleString()} CFT</span>
                        </Link>
                    )}
                    <button onClick={fetchData} className="btn btn-secondary flex items-center gap-2">
                        <RefreshCw size={18} />
                        Refresh
                    </button>
                    {isStartupUser && (
                        <>
                            <button
                                onClick={() => {
                                    if (!canCreateCampaign) {
                                        alert(`Insufficient balance. Need ${FEES.campaignCreation} CFT.`);
                                        return;
                                    }
                                    setShowCreateForm(true);
                                }}
                                className={`btn flex items-center gap-2 ${canCreateCampaign ? 'btn-primary' : 'bg-gray-400 cursor-not-allowed'}`}
                            >
                                <Plus size={18} />
                                New Campaign
                            </button>
                            <button
                                onClick={initiateDeleteStartup}
                                className="btn bg-red-100 text-red-600 hover:bg-red-200 dark:bg-red-900/30 dark:hover:bg-red-900/50 flex items-center gap-2"
                            >
                                <Trash2 size={18} />
                                Delete Startup
                            </button>
                        </>
                    )}
                </div>
            </div>

            {/* Startup Description */}
            {startup?.description && (
                <div className="card">
                    <p className="text-gray-600 dark:text-gray-400">{startup.description}</p>
                </div>
            )}

            {/* Campaigns List */}
            <div className="card">
                <h2 className="text-xl font-bold text-gray-900 dark:text-white mb-4">
                    Campaigns ({campaigns.length})
                </h2>

                {campaigns.length === 0 ? (
                    <div className="text-center py-12">
                        <Rocket size={48} className="mx-auto mb-4 text-gray-300" />
                        <p className="text-gray-500 mb-4">No campaigns yet. Create your first campaign!</p>
                        {isStartupUser && (
                            <button
                                onClick={() => setShowCreateForm(true)}
                                className="btn btn-primary"
                                disabled={!canCreateCampaign}
                            >
                                <Plus size={18} className="mr-2" />
                                Create Campaign ({FEES.campaignCreation} CFT)
                            </button>
                        )}
                    </div>
                ) : (
                    <div className="overflow-x-auto">
                        <table className="w-full">
                            <thead>
                                <tr className="border-b border-gray-200 dark:border-gray-700">
                                    <th className="text-left py-3 px-4 text-gray-600 dark:text-gray-400 font-medium">Campaign</th>
                                    <th className="text-left py-3 px-4 text-gray-600 dark:text-gray-400 font-medium">Goal</th>
                                    <th className="text-left py-3 px-4 text-gray-600 dark:text-gray-400 font-medium">Raised</th>
                                    <th className="text-left py-3 px-4 text-gray-600 dark:text-gray-400 font-medium">Status</th>
                                    <th className="text-left py-3 px-4 text-gray-600 dark:text-gray-400 font-medium">Actions</th>
                                </tr>
                            </thead>
                            <tbody>
                                {campaigns.map((campaign) => (
                                    <tr key={campaign.campaignId} className="border-b border-gray-100 dark:border-gray-800 hover:bg-gray-50 dark:hover:bg-gray-800/50">
                                        <td className="py-4 px-4">
                                            <div>
                                                <p className="font-medium text-gray-900 dark:text-white">
                                                    {campaign.project_name || campaign.projectName}
                                                </p>
                                                <p className="text-xs text-gray-500">{campaign.campaignId}</p>
                                            </div>
                                        </td>
                                        <td className="py-4 px-4 text-gray-900 dark:text-white">
                                            ${(campaign.goal_amount || campaign.goalAmount || 0).toLocaleString()}
                                        </td>
                                        <td className="py-4 px-4">
                                            <div className="text-gray-900 dark:text-white">
                                                ${(campaign.funds_raised_amount || campaign.fundsRaisedAmount || 0).toLocaleString()}
                                                <span className="text-xs text-gray-500 ml-1">
                                                    ({parseFloat(campaign.funds_raised_percent || campaign.fundsRaisedPercent || 0).toFixed(2)}%)
                                                </span>
                                            </div>
                                        </td>
                                        <td className="py-4 px-4">
                                            <div className="flex flex-col gap-1">
                                                <span className={`badge ${getStatusBadge(campaign.status)}`}>
                                                    {campaign.status}
                                                </span>
                                                {campaign.validationStatus && campaign.validationStatus !== 'NOT_SUBMITTED' && (
                                                    <span className={`badge ${getStatusBadge(campaign.validationStatus)} text-xs`}>
                                                        {campaign.validationStatus.replace(/_/g, ' ')}
                                                    </span>
                                                )}
                                            </div>
                                        </td>
                                        <td className="py-4 px-4">
                                            <div className="flex gap-2 flex-wrap">
                                                <Link to={`/startup/campaign/${campaign.campaignId}`} className="text-primary-600 hover:text-primary-800 text-sm font-medium">
                                                    View
                                                </Link>
                                                {campaign.status === 'DRAFT' && !campaign.validationStatus && (
                                                    <button onClick={() => initiateValidationSubmit(campaign.campaignId)} className="text-yellow-600 hover:text-yellow-800 text-sm font-medium flex items-center gap-1">
                                                        <Coins size={14} /> Submit ({FEES.validationSubmission} CFT)
                                                    </button>
                                                )}


                                                {/* Delete button - always visible */}
                                                <button
                                                    onClick={() => initiateDeleteCampaign(campaign)}
                                                    className="text-red-500 hover:text-red-700 text-sm font-medium flex items-center gap-1"
                                                    title="Delete Campaign"
                                                >
                                                    <Trash2 size={14} /> Delete
                                                </button>
                                            </div>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                )}
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
                            {/* Startup Info */}
                            <div className="p-3 bg-primary-50 dark:bg-primary-900/30 rounded-lg">
                                <p className="text-sm text-primary-700 dark:text-primary-300">
                                    <strong>Creating campaign for:</strong> {startup?.name}
                                </p>
                                <p className="text-xs text-primary-600 dark:text-primary-400 mt-1">
                                    Startup ID: {startup?.displayId || startupId} | Campaign ID will be auto-generated
                                </p>
                            </div>

                            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                                <div className="md:col-span-2">
                                    <label className="label">Project Name *</label>
                                    <input type="text" name="projectName" value={formData.projectName} onChange={handleInputChange} className="input w-full" required />
                                </div>
                                <div className="md:col-span-2">
                                    <label className="label">Description *</label>
                                    <textarea name="description" value={formData.description} onChange={handleInputChange} rows={3} className="input w-full" required />
                                </div>
                                <div>
                                    <label className="label">Goal Amount *</label>
                                    <input type="number" name="goalAmount" value={formData.goalAmount} onChange={handleInputChange} min="0" className="input w-full" required />
                                </div>
                                <div>
                                    <label className="label">Currency</label>
                                    <select name="currency" value={formData.currency} onChange={handleInputChange} className="input w-full">
                                        <option value="USD">USD</option>
                                        <option value="EUR">EUR</option>
                                        <option value="INR">INR</option>
                                    </select>
                                </div>
                                <div>
                                    <label className="label">Category</label>
                                    <select name="category" value={formData.category} onChange={handleInputChange} className="input w-full">
                                        <option value="Technology">Technology</option>
                                        <option value="Healthcare">Healthcare</option>
                                        <option value="Finance">Finance</option>
                                        <option value="Education">Education</option>
                                        <option value="Other">Other</option>
                                    </select>
                                </div>
                                <div>
                                    <label className="label">Deadline</label>
                                    <input type="date" name="deadline" value={formData.deadline} onChange={handleInputChange} className="input w-full" />
                                </div>
                            </div>

                            <div className="flex items-center justify-between pt-4 border-t">
                                <div className="text-sm text-gray-500">
                                    Fee: <span className="font-bold text-primary-600">{FEES.campaignCreation} CFT</span>
                                </div>
                                <div className="flex gap-3">
                                    <button type="button" onClick={() => setShowCreateForm(false)} className="btn btn-secondary">Cancel</button>
                                    <button type="submit" disabled={creating || !canCreateCampaign} className="btn btn-primary flex items-center gap-2">
                                        {creating ? <Loader2 className="animate-spin" size={18} /> : <Plus size={18} />}
                                        Create Campaign
                                    </button>
                                </div>
                            </div>
                        </form>
                    </div>
                </div>
            )}

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
                                    <h2 className="text-xl font-bold text-gray-900 dark:text-white">Confirm Payment</h2>
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
                            </div>
                        </div>
                        <div className="p-6 border-t border-gray-200 dark:border-gray-700 flex justify-end gap-3">
                            <button onClick={() => { setShowPaymentModal(false); setPendingAction(null); }} className="btn btn-secondary">Cancel</button>
                            <button onClick={confirmPayment} disabled={processing} className="btn btn-primary flex items-center gap-2">
                                {processing ? <Loader2 className="animate-spin" size={18} /> : <CheckCircle size={18} />}
                                Confirm
                            </button>
                        </div>
                    </div>
                </div>
            )}

            {/* Delete Confirmation Modal */}
            {showDeleteModal && deleteTarget && (
                <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
                    <div className="bg-white dark:bg-gray-800 rounded-xl shadow-2xl max-w-md w-full">
                        <div className="p-6 border-b border-gray-200 dark:border-gray-700">
                            <div className="flex items-center gap-3">
                                <div className="p-3 bg-red-100 dark:bg-red-900 rounded-full">
                                    <AlertTriangle className="text-red-600" size={24} />
                                </div>
                                <div>
                                    <h2 className="text-xl font-bold text-gray-900 dark:text-white">Delete {deleteTarget.type === 'campaign' ? 'Campaign' : 'Startup'}</h2>
                                    <p className="text-sm text-gray-500">This action cannot be undone</p>
                                </div>
                            </div>
                        </div>
                        <div className="p-6 space-y-4">
                            <div className="bg-gray-50 dark:bg-gray-700 rounded-lg p-4">
                                <p className="font-medium text-gray-900 dark:text-white mb-2">
                                    "{deleteTarget.name}"
                                </p>
                                <div className="space-y-2 text-sm">
                                    <div className="flex justify-between">
                                        <span className="text-gray-500">Funds Raised</span>
                                        <span className="text-gray-700 dark:text-gray-300">
                                            ${(deletionFee?.fundsRaised || deleteTarget.fundsRaised || 0).toLocaleString()}
                                        </span>
                                    </div>
                                    <div className="flex justify-between">
                                        <span className="text-gray-500">Deletion Fee</span>
                                        <span className="text-2xl font-bold text-red-600">
                                            {deletionFee?.feeAmount || 100} CFT
                                        </span>
                                    </div>
                                    <div className="flex justify-between items-center text-sm pt-2 border-t border-gray-200 dark:border-gray-600">
                                        <span className="text-gray-500">Your Balance</span>
                                        <span className={`font-bold ${cftBalance < (deletionFee?.feeAmount || 100) ? 'text-red-600' : 'text-gray-700 dark:text-gray-300'}`}>
                                            {cftBalance} CFT
                                        </span>
                                    </div>
                                    <div className="text-xs text-gray-400 mt-1">
                                        {deletionFee?.isFixedFee ? 'Fixed fee (no funds raised)' : '60% of funds raised'}
                                    </div>
                                    {cftBalance < (deletionFee?.feeAmount || 100) && (
                                        <div className="mt-2 p-2 bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-300 text-xs rounded">
                                            Insufficient balance to proceed.
                                        </div>
                                    )}
                                </div>
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                                    Reason for deletion (optional)
                                </label>
                                <textarea
                                    value={deleteReason}
                                    onChange={(e) => setDeleteReason(e.target.value)}
                                    placeholder="Why are you deleting this?"
                                    rows={2}
                                    className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white"
                                />
                            </div>
                        </div>
                        <div className="p-6 border-t border-gray-200 dark:border-gray-700 flex justify-end gap-3">
                            <button onClick={cancelDelete} className="btn btn-secondary">Cancel</button>
                            <button
                                onClick={confirmDelete}
                                disabled={deleting}
                                className="btn bg-red-600 text-white hover:bg-red-700 flex items-center gap-2"
                            >
                                {deleting ? <Loader2 className="animate-spin" size={18} /> : <Trash2 size={18} />}
                                Confirm Delete
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}
