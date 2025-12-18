import { useState, useEffect } from 'react';
import { platformApi } from '../services/api';
import { LayoutDashboard, Upload, Wallet, AlertTriangle, Loader2, RefreshCw, Globe, CheckCircle } from 'lucide-react';

export default function PlatformDashboard() {
    const [sharedCampaigns, setSharedCampaigns] = useState([]);
    const [loading, setLoading] = useState(true);
    const [publishing, setPublishing] = useState(null);

    const fetchSharedCampaigns = async () => {
        setLoading(true);
        try {
            const res = await platformApi.getAllSharedCampaigns();
            setSharedCampaigns(res.data?.data || []);
        } catch (err) {
            setSharedCampaigns([]);
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchSharedCampaigns();
    }, []);

    const handlePublish = async (campaignId, validationHash) => {
        setPublishing(campaignId);
        try {
            await platformApi.publishCampaign(campaignId, { validationHash: validationHash || '' });
            alert('Campaign published successfully!');
            fetchSharedCampaigns();
        } catch (err) {
            console.error('Failed to publish:', err);
            alert('Failed to publish: ' + err.message);
        } finally {
            setPublishing(null);
        }
    };

    return (
        <div className="space-y-6">
            {/* Header */}
            <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
                <div>
                    <h1 className="text-3xl font-bold text-gray-900 dark:text-white">
                        Platform Dashboard
                    </h1>
                    <p className="text-gray-600 dark:text-gray-400 mt-1">
                        Manage campaign publishing and platform operations
                    </p>
                </div>
                <button onClick={fetchSharedCampaigns} className="btn btn-secondary flex items-center gap-2">
                    <RefreshCw size={18} />
                    Refresh
                </button>
            </div>

            {/* Stats */}
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                <div className="card flex items-center gap-4">
                    <div className="p-3 bg-primary-100 dark:bg-primary-900 rounded-lg">
                        <LayoutDashboard className="text-primary-600 dark:text-primary-400" size={24} />
                    </div>
                    <div>
                        <h3 className="font-semibold text-gray-900 dark:text-white">Pending</h3>
                        <p className="text-2xl font-bold text-primary-600">
                            {sharedCampaigns.filter(c => c.status !== 'PUBLISHED').length}
                        </p>
                    </div>
                </div>
                <div className="card flex items-center gap-4">
                    <div className="p-3 bg-green-100 dark:bg-green-900 rounded-lg">
                        <Globe className="text-green-600 dark:text-green-400" size={24} />
                    </div>
                    <div>
                        <h3 className="font-semibold text-gray-900 dark:text-white">Published</h3>
                        <p className="text-2xl font-bold text-green-600">
                            {sharedCampaigns.filter(c => c.status === 'PUBLISHED').length}
                        </p>
                    </div>
                </div>
                <div className="card flex items-center gap-4">
                    <div className="p-3 bg-blue-100 dark:bg-blue-900 rounded-lg">
                        <Wallet className="text-blue-600 dark:text-blue-400" size={24} />
                    </div>
                    <div>
                        <h3 className="font-semibold text-gray-900 dark:text-white">Total Fees</h3>
                        <p className="text-2xl font-bold text-blue-600">$0</p>
                    </div>
                </div>
                <div className="card flex items-center gap-4">
                    <div className="p-3 bg-red-100 dark:bg-red-900 rounded-lg">
                        <AlertTriangle className="text-red-600 dark:text-red-400" size={24} />
                    </div>
                    <div>
                        <h3 className="font-semibold text-gray-900 dark:text-white">Disputes</h3>
                        <p className="text-2xl font-bold text-red-600">0</p>
                    </div>
                </div>
            </div>

            {/* Shared Campaigns */}
            <div className="card">
                <h2 className="text-xl font-bold text-gray-900 dark:text-white mb-4">
                    Shared Campaigns
                </h2>

                {loading ? (
                    <div className="flex justify-center py-12">
                        <Loader2 className="animate-spin text-primary-600" size={40} />
                    </div>
                ) : sharedCampaigns.length === 0 ? (
                    <div className="text-center py-12">
                        <LayoutDashboard size={48} className="mx-auto mb-4 text-gray-300" />
                        <p className="text-gray-500">No campaigns shared yet</p>
                        <p className="text-sm text-gray-400 mt-2">Validated campaigns shared by startups will appear here</p>
                    </div>
                ) : (
                    <div className="overflow-x-auto">
                        <table className="w-full">
                            <thead>
                                <tr className="border-b border-gray-200 dark:border-gray-700">
                                    <th className="text-left py-3 px-4 text-gray-600 dark:text-gray-400 font-medium">Campaign</th>
                                    <th className="text-left py-3 px-4 text-gray-600 dark:text-gray-400 font-medium">Startup</th>
                                    <th className="text-left py-3 px-4 text-gray-600 dark:text-gray-400 font-medium">Goal</th>
                                    <th className="text-left py-3 px-4 text-gray-600 dark:text-gray-400 font-medium">Risk Level</th>
                                    <th className="text-left py-3 px-4 text-gray-600 dark:text-gray-400 font-medium">Status</th>
                                    <th className="text-left py-3 px-4 text-gray-600 dark:text-gray-400 font-medium">Actions</th>
                                </tr>
                            </thead>
                            <tbody>
                                {sharedCampaigns.map((campaign) => {
                                    const isPublished = campaign.status === 'PUBLISHED';
                                    return (
                                        <tr key={campaign.campaignId} className="border-b border-gray-100 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700/50">
                                            <td className="py-4 px-4">
                                                <div>
                                                    <p className="font-medium text-gray-900 dark:text-white">
                                                        {campaign.projectName || campaign.project_name || 'Untitled'}
                                                    </p>
                                                    <p className="text-sm text-gray-500">{campaign.campaignId}</p>
                                                </div>
                                            </td>
                                            <td className="py-4 px-4 text-gray-900 dark:text-white">
                                                {campaign.startupId}
                                            </td>
                                            <td className="py-4 px-4 text-gray-900 dark:text-white">
                                                ${(campaign.goalAmount || campaign.goal_amount || 0).toLocaleString()}
                                            </td>
                                            <td className="py-4 px-4">
                                                <span className={`badge ${campaign.riskLevel === 'LOW' ? 'badge-success' :
                                                    campaign.riskLevel === 'MEDIUM' ? 'badge-warning' : 'badge-danger'
                                                    }`}>
                                                    {campaign.riskLevel || 'N/A'}
                                                </span>
                                            </td>
                                            <td className="py-4 px-4">
                                                <span className={`badge ${isPublished ? 'badge-success' : 'badge-warning'}`}>
                                                    {isPublished ? 'PUBLISHED TO PORTAL' : 'PENDING REVIEW'}
                                                </span>
                                            </td>
                                            <td className="py-4 px-4">
                                                {isPublished ? (
                                                    <span className="flex items-center text-green-600 text-sm font-medium">
                                                        <CheckCircle size={16} className="mr-1" />
                                                        Already Published
                                                    </span>
                                                ) : (
                                                    <button
                                                        onClick={() => handlePublish(campaign.campaignId, campaign.validationHash)}
                                                        disabled={publishing === campaign.campaignId}
                                                        className="btn btn-primary text-sm flex items-center gap-1"
                                                    >
                                                        {publishing === campaign.campaignId ? (
                                                            <Loader2 className="animate-spin" size={16} />
                                                        ) : (
                                                            <Upload size={16} />
                                                        )}
                                                        Publish to Portal
                                                    </button>
                                                )}
                                            </td>
                                        </tr>
                                    );
                                })}
                            </tbody>
                        </table>
                    </div>
                )}
            </div>
        </div>
    );
}
