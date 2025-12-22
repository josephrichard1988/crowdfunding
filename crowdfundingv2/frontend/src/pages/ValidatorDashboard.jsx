import { useState, useEffect } from 'react';
import { validatorApi } from '../services/api';
import { useAuth } from '../context/AuthContext';
import { Link } from 'react-router-dom';
import { Shield, CheckCircle, XCircle, Clock, Loader2, RefreshCw, FileText, AlertTriangle, LogIn, Coins } from 'lucide-react';

// Token constants
const FEES = {
    validationFee: 500,
    disputeFee: 750
};

export default function ValidatorDashboard() {
    const { user, isAuthenticated } = useAuth();
    const [pendingValidations, setPendingValidations] = useState([]);
    const [loading, setLoading] = useState(true);
    const [selectedCampaign, setSelectedCampaign] = useState(null);
    const [approving, setApproving] = useState(false);
    const [formData, setFormData] = useState({
        dueDiligenceScore: 8.5,
        riskScore: 3.0,
        riskLevel: 'LOW',
        comments: '',
    });

    const cftBalance = user?.wallet?.cftBalance || 0;
    const isValidatorUser = isAuthenticated && user?.role === 'VALIDATOR';
    const isPreviewMode = !isAuthenticated || user?.role !== 'VALIDATOR';

    const fetchPendingValidations = async () => {
        setLoading(true);
        try {
            const res = await validatorApi.getPendingValidations();
            setPendingValidations(res.data?.data || []);
        } catch (err) {
            setPendingValidations([]);
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchPendingValidations();
    }, []);

    const handleApprove = async (campaignId) => {
        if (isPreviewMode) {
            alert('Please login as a validator to approve campaigns');
            return;
        }
        setApproving(true);
        try {
            await validatorApi.approveCampaign(campaignId, {
                validationId: `VAL_${Date.now()}`,
                status: 'APPROVED',
                dueDiligenceScore: formData.dueDiligenceScore,
                riskScore: formData.riskScore,
                riskLevel: formData.riskLevel,
                comments: formData.comments ? [formData.comments] : ['Approved'],
                issues: [],
                requiredDocuments: '',
            });
            setSelectedCampaign(null);
            fetchPendingValidations();
            alert('Campaign approved successfully!');
        } catch (err) {
            console.error('Failed to approve:', err);
            alert('Failed to approve: ' + err.message);
        } finally {
            setApproving(false);
        }
    };

    const handleReject = async (campaignId) => {
        if (isPreviewMode) {
            alert('Please login as a validator to reject campaigns');
            return;
        }
        setApproving(true);
        try {
            await validatorApi.approveCampaign(campaignId, {
                validationId: `VAL_${Date.now()}`,
                status: 'REJECTED',
                dueDiligenceScore: formData.dueDiligenceScore,
                riskScore: formData.riskScore,
                riskLevel: formData.riskLevel,
                comments: formData.comments ? [formData.comments] : ['Rejected'],
                issues: ['Did not meet requirements'],
                requiredDocuments: '',
            });
            setSelectedCampaign(null);
            fetchPendingValidations();
            alert('Campaign rejected.');
        } catch (err) {
            console.error('Failed to reject:', err);
            alert('Failed to reject: ' + err.message);
        } finally {
            setApproving(false);
        }
    };

    return (
        <div className="space-y-6">
            {/* Preview Mode Banner */}
            {isPreviewMode && (
                <div className="bg-gradient-to-r from-purple-500 to-purple-600 text-white p-4 rounded-xl flex flex-col sm:flex-row items-center justify-between gap-4">
                    <div className="flex items-center gap-3">
                        <Shield size={24} />
                        <div>
                            <h3 className="font-bold">Validator Dashboard Preview</h3>
                            <p className="text-sm opacity-90">Login as a validator to review and validate campaigns</p>
                        </div>
                    </div>
                    <Link to="/login" state={{ role: 'VALIDATOR' }} className="btn bg-white text-purple-700 hover:bg-gray-100 flex items-center gap-2">
                        <LogIn size={18} />
                        Login as Validator
                    </Link>
                </div>
            )}

            {/* Header */}
            <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
                <div>
                    <h1 className="text-3xl font-bold text-gray-900 dark:text-white">
                        Validator Dashboard
                    </h1>
                    <p className="text-gray-600 dark:text-gray-400 mt-1">
                        Review and validate campaign submissions
                    </p>
                </div>
                <div className="flex gap-3 items-center">
                    {isValidatorUser && (
                        <Link to="/wallet" className="flex items-center gap-2 px-4 py-2 bg-purple-100 dark:bg-purple-900/30 text-purple-700 dark:text-purple-300 rounded-lg hover:bg-purple-200">
                            <Coins size={18} />
                            <span className="font-medium">{cftBalance.toLocaleString()} CFT</span>
                        </Link>
                    )}
                    <button onClick={fetchPendingValidations} className="btn btn-secondary flex items-center gap-2">
                        <RefreshCw size={18} />
                        Refresh
                    </button>
                </div>
            </div>

            {/* Stats */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div className="card flex items-center gap-4">
                    <div className="p-3 bg-yellow-100 dark:bg-yellow-900 rounded-lg">
                        <Clock className="text-yellow-600 dark:text-yellow-400" size={24} />
                    </div>
                    <div>
                        <h3 className="font-semibold text-gray-900 dark:text-white">Pending</h3>
                        <p className="text-2xl font-bold text-yellow-600">{pendingValidations.length}</p>
                    </div>
                </div>
                <div className="card flex items-center gap-4">
                    <div className="p-3 bg-green-100 dark:bg-green-900 rounded-lg">
                        <CheckCircle className="text-green-600 dark:text-green-400" size={24} />
                    </div>
                    <div>
                        <h3 className="font-semibold text-gray-900 dark:text-white">Approved Today</h3>
                        <p className="text-2xl font-bold text-green-600">--</p>
                    </div>
                </div>
                <div className="card flex items-center gap-4">
                    <div className="p-3 bg-red-100 dark:bg-red-900 rounded-lg">
                        <XCircle className="text-red-600 dark:text-red-400" size={24} />
                    </div>
                    <div>
                        <h3 className="font-semibold text-gray-900 dark:text-white">Rejected Today</h3>
                        <p className="text-2xl font-bold text-red-600">--</p>
                    </div>
                </div>
            </div>

            {/* Validation Modal */}
            {selectedCampaign && (
                <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
                    <div className="bg-white dark:bg-gray-800 rounded-xl shadow-2xl max-w-lg w-full">
                        <div className="p-6 border-b border-gray-200 dark:border-gray-700">
                            <h2 className="text-xl font-bold text-gray-900 dark:text-white">
                                Validate Campaign
                            </h2>
                            <p className="text-gray-500">{selectedCampaign.campaignId}</p>
                        </div>
                        <div className="p-6 space-y-4">
                            <div className="grid grid-cols-2 gap-4">
                                <div>
                                    <label className="label">Due Diligence Score</label>
                                    <input
                                        type="number"
                                        step="0.1"
                                        min="0"
                                        max="10"
                                        value={formData.dueDiligenceScore}
                                        onChange={(e) => setFormData({ ...formData, dueDiligenceScore: parseFloat(e.target.value) })}
                                        className="input"
                                    />
                                </div>
                                <div>
                                    <label className="label">Risk Score</label>
                                    <input
                                        type="number"
                                        step="0.1"
                                        min="0"
                                        max="10"
                                        value={formData.riskScore}
                                        onChange={(e) => setFormData({ ...formData, riskScore: parseFloat(e.target.value) })}
                                        className="input"
                                    />
                                </div>
                            </div>
                            <div>
                                <label className="label">Risk Level</label>
                                <select
                                    value={formData.riskLevel}
                                    onChange={(e) => setFormData({ ...formData, riskLevel: e.target.value })}
                                    className="input"
                                >
                                    <option value="LOW">LOW</option>
                                    <option value="MEDIUM">MEDIUM</option>
                                    <option value="HIGH">HIGH</option>
                                </select>
                            </div>
                            <div>
                                <label className="label">Comments</label>
                                <textarea
                                    value={formData.comments}
                                    onChange={(e) => setFormData({ ...formData, comments: e.target.value })}
                                    className="input"
                                    rows={3}
                                    placeholder="Add validation comments..."
                                />
                            </div>
                        </div>
                        <div className="p-6 border-t border-gray-200 dark:border-gray-700 flex justify-end gap-3">
                            <button onClick={() => setSelectedCampaign(null)} className="btn btn-secondary">
                                Cancel
                            </button>
                            <button
                                onClick={() => handleReject(selectedCampaign.campaignId)}
                                disabled={approving}
                                className="btn bg-red-600 text-white hover:bg-red-700"
                            >
                                {approving ? <Loader2 className="animate-spin" size={18} /> : <XCircle size={18} />}
                                <span className="ml-2">Reject</span>
                            </button>
                            <button
                                onClick={() => handleApprove(selectedCampaign.campaignId)}
                                disabled={approving}
                                className="btn btn-primary"
                            >
                                {approving ? <Loader2 className="animate-spin" size={18} /> : <CheckCircle size={18} />}
                                <span className="ml-2">Approve</span>
                            </button>
                        </div>
                    </div>
                </div>
            )}

            {/* Pending Validations */}
            <div className="card">
                <h2 className="text-xl font-bold text-gray-900 dark:text-white mb-4">
                    Pending Validations
                </h2>

                {loading ? (
                    <div className="flex justify-center py-12">
                        <Loader2 className="animate-spin text-primary-600" size={40} />
                    </div>
                ) : pendingValidations.length === 0 ? (
                    <div className="text-center py-12">
                        <Shield size={48} className="mx-auto mb-4 text-gray-300" />
                        <p className="text-gray-500">No pending validations</p>
                        <p className="text-sm text-gray-400 mt-2">Campaigns submitted for validation will appear here</p>
                    </div>
                ) : (
                    <div className="space-y-4">
                        {pendingValidations.map((campaign) => (
                            <div key={campaign.campaignId} className="border border-gray-200 dark:border-gray-700 rounded-lg p-4">
                                <div className="flex justify-between items-start">
                                    <div className="flex-1">
                                        <h3 className="font-semibold text-gray-900 dark:text-white">
                                            {campaign.projectName || campaign.project_name || 'Untitled'}
                                        </h3>
                                        <p className="text-sm text-gray-500">{campaign.campaignId}</p>
                                        <p className="text-sm text-gray-600 dark:text-gray-400 mt-2">
                                            {campaign.description?.substring(0, 150)}...
                                        </p>
                                        <div className="flex gap-2 mt-3">
                                            <span className="badge badge-info">{campaign.category}</span>
                                            <span className="badge badge-info">${(campaign.goalAmount || campaign.goal_amount || 0).toLocaleString()}</span>
                                        </div>
                                    </div>
                                    <button
                                        onClick={() => setSelectedCampaign(campaign)}
                                        className="btn btn-primary text-sm"
                                    >
                                        <FileText size={16} className="mr-1" />
                                        Review
                                    </button>
                                </div>
                            </div>
                        ))}
                    </div>
                )}
            </div>
        </div>
    );
}
