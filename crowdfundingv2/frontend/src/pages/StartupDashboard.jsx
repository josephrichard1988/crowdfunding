import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { startupApi } from '../services/api';
import { Plus, FileText, Share2, Bell, Loader2, RefreshCw, X, Rocket } from 'lucide-react';

export default function StartupDashboard() {
    const [campaigns, setCampaigns] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [showCreateForm, setShowCreateForm] = useState(false);
    const [creating, setCreating] = useState(false);
    const [formData, setFormData] = useState({
        campaignId: '',
        startupId: '',
        projectName: '',
        description: '',
        category: 'Technology',
        goalAmount: '',
        currency: 'USD',
        deadline: '',
        tags: '',
        documents: '',
    });

    const fetchCampaigns = async () => {
        setLoading(true);
        try {
            const res = await startupApi.getAllCampaigns();
            setCampaigns(res.data?.data || []);
            setError(null);
        } catch (err) {
            setCampaigns([]);
            setError(null); // Don't show error, just empty state
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchCampaigns();
    }, []);

    const handleInputChange = (e) => {
        const { name, value } = e.target;
        setFormData(prev => ({ ...prev, [name]: value }));
    };

    const handleCreateCampaign = async (e) => {
        e.preventDefault();
        setCreating(true);
        try {
            const tags = formData.tags.split(',').map(t => t.trim()).filter(Boolean);
            const documents = formData.documents.split(',').map(d => d.trim()).filter(Boolean);

            await startupApi.createCampaign({
                campaignId: formData.campaignId || `CAMP_${Date.now()}`,
                startupId: formData.startupId || 'STARTUP001',
                category: formData.category,
                deadline: formData.deadline,
                currency: formData.currency,
                hasRaised: false,
                hasGovGrants: false,
                incorpDate: new Date().toISOString().split('T')[0],
                projectStage: 'Idea',
                sector: formData.category,
                tags,
                teamAvailable: true,
                investorCommitted: false,
                duration: 90,
                fundingDay: 1,
                fundingMonth: 1,
                fundingYear: new Date().getFullYear(),
                goalAmount: parseFloat(formData.goalAmount) || 50000,
                investmentRange: '10K-100K',
                projectName: formData.projectName,
                description: formData.description,
                documents,
            });

            setShowCreateForm(false);
            setFormData({
                campaignId: '',
                startupId: '',
                projectName: '',
                description: '',
                category: 'Technology',
                goalAmount: '',
                currency: 'USD',
                deadline: '',
                tags: '',
                documents: '',
            });
            fetchCampaigns();
        } catch (err) {
            console.error('Failed to create campaign:', err);
            alert('Failed to create campaign: ' + err.message);
        } finally {
            setCreating(false);
        }
    };

    const handleSubmitValidation = async (campaignId) => {
        try {
            await startupApi.submitForValidation(campaignId, {
                documents: [],
                notes: 'Submitting for validation'
            });
            alert('Campaign submitted for validation!');
            fetchCampaigns();
        } catch (err) {
            console.error('Failed to submit for validation:', err);
            alert('Failed to submit: ' + err.message);
        }
    };

    const handleShareToPlatform = async (campaignId, validationHash) => {
        try {
            await startupApi.shareToPlatform(campaignId, {
                validationHash: validationHash || ''
            });
            alert('Campaign shared to platform!');
            fetchCampaigns();
        } catch (err) {
            console.error('Failed to share to platform:', err);
            alert('Failed to share: ' + err.message);
        }
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
                <div className="flex gap-3">
                    <button onClick={fetchCampaigns} className="btn btn-secondary flex items-center gap-2">
                        <RefreshCw size={18} />
                        Refresh
                    </button>
                    <button onClick={() => setShowCreateForm(true)} className="btn btn-primary flex items-center gap-2">
                        <Plus size={18} />
                        New Campaign
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
                            <div className="grid md:grid-cols-2 gap-4">
                                <div>
                                    <label className="label">Campaign ID</label>
                                    <input
                                        type="text"
                                        name="campaignId"
                                        value={formData.campaignId}
                                        onChange={handleInputChange}
                                        placeholder="CAMP001"
                                        className="input"
                                    />
                                </div>
                                <div>
                                    <label className="label">Startup ID</label>
                                    <input
                                        type="text"
                                        name="startupId"
                                        value={formData.startupId}
                                        onChange={handleInputChange}
                                        placeholder="STARTUP001"
                                        className="input"
                                    />
                                </div>
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

            {/* Campaigns List */}
            <div className="card">
                <h2 className="text-xl font-bold text-gray-900 dark:text-white mb-4">Your Campaigns</h2>

                {loading ? (
                    <div className="flex justify-center py-12">
                        <Loader2 className="animate-spin text-primary-600" size={40} />
                    </div>
                ) : campaigns.length === 0 ? (
                    <div className="text-center py-12">
                        <Rocket size={48} className="mx-auto mb-4 text-gray-300" />
                        <p className="text-gray-500">No campaigns yet. Create your first campaign!</p>
                        <button
                            onClick={() => setShowCreateForm(true)}
                            className="btn btn-primary mt-4"
                        >
                            <Plus size={18} className="mr-2" />
                            Create Campaign
                        </button>
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
                                    <tr key={campaign.campaignId} className="border-b border-gray-100 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700/50">
                                        <td className="py-4 px-4">
                                            <div>
                                                <p className="font-medium text-gray-900 dark:text-white">{campaign.project_name || campaign.projectName}</p>
                                                <p className="text-sm text-gray-500">{campaign.campaignId}</p>
                                            </div>
                                        </td>
                                        <td className="py-4 px-4 text-gray-900 dark:text-white">
                                            ${(campaign.goal_amount || campaign.goalAmount || 0).toLocaleString()}
                                        </td>
                                        <td className="py-4 px-4 text-gray-900 dark:text-white">
                                            ${(campaign.funds_raised_amount || campaign.fundsRaisedAmount || 0).toLocaleString()}
                                        </td>
                                        <td className="py-4 px-4">
                                            <span className={`badge ${getStatusBadge(campaign.status)}`}>
                                                {campaign.status}
                                            </span>
                                        </td>
                                        <td className="py-4 px-4">
                                            <div className="flex gap-2 flex-wrap">
                                                {/* View campaign details */}
                                                <Link
                                                    to={`/startup/campaign/${campaign.campaignId}`}
                                                    className="text-primary-600 hover:text-primary-800 text-sm font-medium"
                                                >
                                                    View
                                                </Link>
                                                {/* Submit for Validation - only for DRAFT campaigns */}
                                                {campaign.status === 'DRAFT' && !campaign.validationStatus && (
                                                    <button
                                                        onClick={() => handleSubmitValidation(campaign.campaignId)}
                                                        className="text-yellow-600 hover:text-yellow-800 text-sm font-medium"
                                                    >
                                                        Submit Validation
                                                    </button>
                                                )}
                                                {/* Share to Platform - only for APPROVED campaigns */}
                                                {campaign.validationStatus === 'APPROVED' && campaign.status !== 'PUBLISHED' && (
                                                    <button
                                                        onClick={() => handleShareToPlatform(campaign.campaignId, campaign.validationHash)}
                                                        className="text-accent-600 hover:text-accent-800 text-sm font-medium"
                                                    >
                                                        Share to Platform
                                                    </button>
                                                )}
                                                {/* Status indicator */}
                                                {campaign.validationStatus === 'PENDING_VALIDATION' && (
                                                    <span className="text-xs text-yellow-600">Pending Validation</span>
                                                )}
                                                {campaign.status === 'PUBLISHED' && (
                                                    <span className="text-xs text-green-600">Published ✓</span>
                                                )}
                                            </div>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                )}
            </div>
        </div>
    );
}
