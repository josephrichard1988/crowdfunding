import { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { startupApi } from '../services/api';
import { ArrowLeft, Calendar, DollarSign, Target, FileText, CheckCircle, AlertCircle, Send, Share2, Loader2, X } from 'lucide-react';

export default function StartupCampaignDetails() {
    const { campaignId } = useParams();
    const [campaign, setCampaign] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [submitting, setSubmitting] = useState(false);
    const [showSubmitModal, setShowSubmitModal] = useState(false);
    const [submitNotes, setSubmitNotes] = useState('');

    useEffect(() => {
        fetchCampaign();
    }, [campaignId]);

    const fetchCampaign = async () => {
        try {
            setLoading(true);
            const response = await startupApi.getCampaign(campaignId);
            setCampaign(response.data.data);
            setError(null);
        } catch (err) {
            console.error('Failed to fetch campaign:', err);
            setError('Campaign not found');
        } finally {
            setLoading(false);
        }
    };

    const handleSubmitForValidation = async () => {
        try {
            setSubmitting(true);
            const authToken = sessionStorage.getItem('token') || localStorage.getItem('token');
            await startupApi.submitForValidation(campaignId, {
                documents: campaign.documents || [],
                notes: submitNotes || 'Please validate this campaign',
                authToken,
                startupId: campaign.startupId,
                projectName: campaign.projectName || campaign.project_name
            });
            alert('Campaign submitted for validation!');
            setShowSubmitModal(false);
            setSubmitNotes('');
            fetchCampaign(); // Refresh data
        } catch (err) {
            console.error('Failed to submit:', err);
            alert('Failed to submit: ' + err.message);
        } finally {
            setSubmitting(false);
        }
    };

    const handleShareToPlatform = async () => {
        try {
            setSubmitting(true);
            const authToken = sessionStorage.getItem('token') || localStorage.getItem('token');
            await startupApi.shareToPlatform(campaignId, {
                validationProofHash: campaign.validationProofHash || '',
                authToken,
                startupId: campaign.startupId,
                projectName: campaign.projectName || campaign.project_name
            });
            alert('Campaign shared to platform for publishing!');
            fetchCampaign(); // Refresh data
        } catch (err) {
            console.error('Failed to share:', err);
            alert('Failed to share: ' + err.message);
        } finally {
            setSubmitting(false);
        }
    };

    const getStatusBadge = (status) => {
        const badges = {
            DRAFT: 'badge-info',
            SUBMITTED: 'badge-warning',
            APPROVED: 'badge-success',
            PUBLISHED: 'badge-success',
            REJECTED: 'badge-error',
            PENDING_VALIDATION: 'badge-warning',
            PENDING_PLATFORM_APPROVAL: 'badge-warning',
        };
        return badges[status] || 'badge-info';
    };

    const canSubmitForValidation = campaign?.status === 'DRAFT' &&
        (!campaign?.validationStatus || campaign?.validationStatus === 'NOT_SUBMITTED');

    const canShareToPlatform = campaign?.validationStatus === 'APPROVED' &&
        campaign?.status !== 'PUBLISHED' &&
        campaign?.status !== 'PENDING_PLATFORM_APPROVAL';

    if (loading) {
        return (
            <div className="flex items-center justify-center min-h-96">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
            </div>
        );
    }

    if (error || !campaign) {
        return (
            <div className="card text-center py-12">
                <AlertCircle className="mx-auto text-red-500 mb-4" size={48} />
                <h2 className="text-xl font-bold text-gray-900 dark:text-white mb-2">Campaign Not Found</h2>
                <p className="text-gray-500">{error}</p>
                <Link to="/startup" className="btn btn-primary mt-4">
                    Back to Dashboard
                </Link>
            </div>
        );
    }

    return (
        <div className="space-y-6">
            {/* Submit Validation Modal */}
            {showSubmitModal && (
                <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
                    <div className="bg-white dark:bg-gray-800 rounded-xl shadow-2xl max-w-md w-full">
                        <div className="flex justify-between items-center p-6 border-b border-gray-200 dark:border-gray-700">
                            <h2 className="text-xl font-bold text-gray-900 dark:text-white flex items-center gap-2">
                                <Send className="text-primary-600" size={24} />
                                Submit for Validation
                            </h2>
                            <button onClick={() => setShowSubmitModal(false)} className="text-gray-500 hover:text-gray-700">
                                <X size={24} />
                            </button>
                        </div>
                        <div className="p-6 space-y-4">
                            <p className="text-gray-600 dark:text-gray-400">
                                Submit <strong>{campaign.projectName || campaign.project_name}</strong> to ValidatorOrg for verification.
                            </p>
                            <div>
                                <label className="label">Notes for Validator</label>
                                <textarea
                                    value={submitNotes}
                                    onChange={(e) => setSubmitNotes(e.target.value)}
                                    placeholder="Please validate our campaign for crowdfunding..."
                                    className="input"
                                    rows={3}
                                />
                            </div>
                            <div className="flex gap-3 pt-4">
                                <button
                                    onClick={() => setShowSubmitModal(false)}
                                    className="btn btn-secondary flex-1"
                                    disabled={submitting}
                                >
                                    Cancel
                                </button>
                                <button
                                    onClick={handleSubmitForValidation}
                                    className="btn btn-primary flex-1 flex items-center justify-center gap-2"
                                    disabled={submitting}
                                >
                                    {submitting ? <Loader2 className="animate-spin" size={18} /> : <Send size={18} />}
                                    {submitting ? 'Submitting...' : 'Submit'}
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            )}

            {/* Header with Actions */}
            <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
                <div className="flex items-center gap-4">
                    <Link to="/startup" className="text-gray-500 hover:text-gray-700 dark:hover:text-gray-300">
                        <ArrowLeft size={24} />
                    </Link>
                    <div>
                        <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
                            {campaign.projectName || campaign.project_name}
                        </h1>
                        <p className="text-gray-500">Campaign ID: {campaign.campaignId}</p>
                    </div>
                </div>
                <div className="flex items-center gap-3">
                    <span className={`badge ${getStatusBadge(campaign.status)}`}>
                        {campaign.status}
                    </span>
                    {campaign.validationStatus && campaign.validationStatus !== 'NOT_SUBMITTED' && (
                        <span className={`badge ${getStatusBadge(campaign.validationStatus)}`}>
                            {campaign.validationStatus}
                        </span>
                    )}
                </div>
            </div>

            {/* Action Buttons */}
            <div className="card bg-gradient-to-r from-primary-50 to-accent-50 dark:from-primary-900/20 dark:to-accent-900/20 border-2 border-primary-200 dark:border-primary-800">
                <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Campaign Actions</h3>
                <div className="flex flex-wrap gap-3">
                    {canSubmitForValidation && (
                        <button
                            onClick={() => setShowSubmitModal(true)}
                            className="btn btn-primary flex items-center gap-2"
                            disabled={submitting}
                        >
                            <Send size={18} />
                            Submit for Validation
                        </button>
                    )}

                    {canShareToPlatform && (
                        <button
                            onClick={handleShareToPlatform}
                            className="btn bg-green-600 text-white hover:bg-green-700 focus:ring-green-500 flex items-center gap-2 shadow-md hover:shadow-lg"
                            disabled={submitting}
                        >
                            <Share2 size={18} />
                            Send to Platform for Publishing
                        </button>
                    )}

                    {campaign.validationStatus === 'PENDING_VALIDATION' && (
                        <span className="px-4 py-2 bg-yellow-100 dark:bg-yellow-900/30 text-yellow-700 dark:text-yellow-300 rounded-lg text-sm flex items-center gap-2">
                            <Loader2 className="animate-spin" size={16} />
                            Awaiting Validator Review...
                        </span>
                    )}
                    {campaign.status === 'PENDING_PLATFORM_APPROVAL' && (
                        <span className="px-4 py-2 bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300 rounded-lg text-sm flex items-center gap-2">
                            <Loader2 className="animate-spin" size={16} />
                            Awaiting Platform Publishing...
                        </span>
                    )}
                    {campaign.status === 'PUBLISHED' && (
                        <span className="px-4 py-2 bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-300 rounded-lg text-sm flex items-center gap-2">
                            <CheckCircle size={16} />
                            Published & Live!
                        </span>
                    )}
                    {!canSubmitForValidation && !canShareToPlatform && campaign.status === 'DRAFT' && (
                        <span className="text-gray-500 text-sm">No actions available at this stage.</span>
                    )}
                </div>
            </div>

            {/* Campaign Details Grid */}
            <div className="grid md:grid-cols-2 gap-6">
                {/* Overview Card */}
                <div className="card">
                    <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4 flex items-center gap-2">
                        <Target className="text-primary-600" size={20} />
                        Campaign Overview
                    </h3>
                    <div className="space-y-3">
                        <div className="flex justify-between">
                            <span className="text-gray-500">Category</span>
                            <span className="font-medium text-gray-900 dark:text-white">{campaign.category}</span>
                        </div>
                        <div className="flex justify-between">
                            <span className="text-gray-500">Sector</span>
                            <span className="font-medium text-gray-900 dark:text-white">{campaign.sector}</span>
                        </div>
                        <div className="flex justify-between">
                            <span className="text-gray-500">Project Stage</span>
                            <span className="font-medium text-gray-900 dark:text-white">{campaign.projectStage || campaign.project_stage}</span>
                        </div>
                        <div className="flex justify-between">
                            <span className="text-gray-500">Investment Range</span>
                            <span className="font-medium text-gray-900 dark:text-white">{campaign.investmentRange || campaign.investment_range}</span>
                        </div>
                    </div>
                </div>

                {/* Funding Card */}
                <div className="card">
                    <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4 flex items-center gap-2">
                        <DollarSign className="text-green-600" size={20} />
                        Funding Details
                    </h3>
                    <div className="space-y-3">
                        <div className="flex justify-between">
                            <span className="text-gray-500">Goal Amount</span>
                            <span className="font-medium text-gray-900 dark:text-white">
                                {campaign.currency} {(campaign.goalAmount || campaign.goal_amount || 0).toLocaleString()}
                            </span>
                        </div>
                        <div className="flex justify-between">
                            <span className="text-gray-500">Funds Raised</span>
                            <span className="font-medium text-green-600">
                                {campaign.currency} {(campaign.fundsRaisedAmount || campaign.funds_raised_amount || 0).toLocaleString()}
                            </span>
                        </div>
                        <div className="flex justify-between">
                            <span className="text-gray-500">Duration</span>
                            <span className="font-medium text-gray-900 dark:text-white">{campaign.duration} days</span>
                        </div>
                        <div className="flex justify-between">
                            <span className="text-gray-500">Deadline</span>
                            <span className="font-medium text-gray-900 dark:text-white">{campaign.deadline}</span>
                        </div>
                    </div>
                    {/* Progress Bar */}
                    <div className="mt-4">
                        <div className="flex justify-between text-sm mb-1">
                            <span className="text-gray-500">Progress</span>
                            <span className="font-medium text-gray-900 dark:text-white">
                                {parseFloat(campaign.fundsRaisedPercent || campaign.funds_raised_percent || 0).toFixed(2)}%
                            </span>
                        </div>
                        <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                            <div
                                className="bg-gradient-to-r from-primary-500 to-accent-500 h-2 rounded-full"
                                style={{ width: `${Math.min(campaign.fundsRaisedPercent || 0, 100)}%` }}
                            ></div>
                        </div>
                    </div>
                </div>

                {/* Validation Card */}
                <div className="card">
                    <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4 flex items-center gap-2">
                        <CheckCircle className="text-blue-600" size={20} />
                        Validation Status
                    </h3>
                    <div className="space-y-3">
                        <div className="flex justify-between">
                            <span className="text-gray-500">Validation Status</span>
                            <span className={`badge ${getStatusBadge(campaign.validationStatus)}`}>
                                {campaign.validationStatus || 'NOT_SUBMITTED'}
                            </span>
                        </div>
                        {campaign.validationScore > 0 && (
                            <div className="flex justify-between">
                                <span className="text-gray-500">Validation Score</span>
                                <span className="font-medium text-gray-900 dark:text-white">{campaign.validationScore}/10</span>
                            </div>
                        )}
                        {campaign.riskLevel && (
                            <div className="flex justify-between">
                                <span className="text-gray-500">Risk Level</span>
                                <span className={`font-medium ${campaign.riskLevel === 'LOW' ? 'text-green-600' :
                                    campaign.riskLevel === 'MEDIUM' ? 'text-yellow-600' : 'text-red-600'
                                    }`}>{campaign.riskLevel}</span>
                            </div>
                        )}
                        {campaign.validationProofHash && (
                            <div className="flex justify-between">
                                <span className="text-gray-500">Validation Proof Hash</span>
                                <span className="font-mono text-xs text-gray-600 dark:text-gray-400 truncate max-w-[150px]" title={campaign.validationProofHash}>
                                    {campaign.validationProofHash.substring(0, 20)}...
                                </span>
                            </div>
                        )}
                    </div>
                </div>

                {/* Dates Card */}
                <div className="card">
                    <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4 flex items-center gap-2">
                        <Calendar className="text-purple-600" size={20} />
                        Timeline
                    </h3>
                    <div className="space-y-3">
                        <div className="flex justify-between">
                            <span className="text-gray-500">Incorporation Date</span>
                            <span className="font-medium text-gray-900 dark:text-white">{campaign.incorpDate || campaign.incorp_date}</span>
                        </div>
                        <div className="flex justify-between">
                            <span className="text-gray-500">Campaign Open</span>
                            <span className="font-medium text-gray-900 dark:text-white">{campaign.openDate || campaign.open_date || 'Not set'}</span>
                        </div>
                        <div className="flex justify-between">
                            <span className="text-gray-500">Campaign Close</span>
                            <span className="font-medium text-gray-900 dark:text-white">{campaign.closeDate || campaign.close_date || campaign.deadline}</span>
                        </div>
                    </div>
                </div>
            </div>

            {/* Description */}
            <div className="card">
                <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4 flex items-center gap-2">
                    <FileText className="text-indigo-600" size={20} />
                    Description
                </h3>
                <p className="text-gray-600 dark:text-gray-400">
                    {campaign.description || 'No description provided.'}
                </p>
            </div>

            {/* Tags */}
            {campaign.tags && campaign.tags.length > 0 && (
                <div className="card">
                    <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Tags</h3>
                    <div className="flex flex-wrap gap-2">
                        {campaign.tags.map((tag, index) => (
                            <span key={index} className="px-3 py-1 bg-primary-100 dark:bg-primary-900/30 text-primary-700 dark:text-primary-300 rounded-full text-sm">
                                {tag}
                            </span>
                        ))}
                    </div>
                </div>
            )}

            {/* Documents */}
            {campaign.documents && campaign.documents.length > 0 && (
                <div className="card">
                    <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Documents</h3>
                    <ul className="space-y-2">
                        {campaign.documents.map((doc, index) => (
                            <li key={index} className="flex items-center gap-2 text-gray-600 dark:text-gray-400">
                                <FileText size={16} className="text-gray-400" />
                                {doc}
                            </li>
                        ))}
                    </ul>
                </div>
            )}
        </div>
    );
}

