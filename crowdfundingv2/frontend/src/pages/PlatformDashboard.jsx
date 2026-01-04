import { useState, useEffect } from 'react';
import { platformApi, queueApi } from '../services/api';
import { useAuth } from '../context/AuthContext';
import { Link } from 'react-router-dom';
import { LayoutDashboard, Upload, Wallet, AlertTriangle, Loader2, RefreshCw, Globe, CheckCircle, LogIn, Coins, ListChecks, History, Settings } from 'lucide-react';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:3001/api';

// Token constants
const FEES = {
    publishingFee: 2500,
    disputeFee: 750
};

export default function PlatformDashboard() {
    const { user, isAuthenticated, token } = useAuth();
    const [sharedCampaigns, setSharedCampaigns] = useState([]);
    const [assignedQueue, setAssignedQueue] = useState([]);
    const [completedTasks, setCompletedTasks] = useState([]);
    const [loading, setLoading] = useState(true);
    const [publishing, setPublishing] = useState(null);

    // CFT Supply Management
    const [cftSupply, setCftSupply] = useState(0);
    const [newSupplyAmount, setNewSupplyAmount] = useState('');
    const [settingSupply, setSettingSupply] = useState(false);

    const cftBalance = user?.wallet?.cftBalance || 0;
    const isPlatformUser = isAuthenticated && user?.role === 'PLATFORM';
    const isPreviewMode = !isAuthenticated || user?.role !== 'PLATFORM';

    const fetchData = async () => {
        if (!isPlatformUser) return;

        setLoading(true);
        try {
            // Fetch user's assigned queue from MongoDB
            let queue = [];
            let completed = [];
            try {
                const queueRes = await queueApi.getQueue();
                queue = queueRes.data?.data?.assignedQueue || [];
                completed = queueRes.data?.data?.completedTasks || [];
            } catch (e) {
                console.warn('Failed to fetch queue:', e.message);
            }

            setAssignedQueue(queue.filter(t => t.type === 'PUBLISH'));
            setCompletedTasks(completed.filter(t => t.type === 'PUBLISH'));

            // Also fetch all shared campaigns for reference
            let allCampaigns = [];
            try {
                const res = await platformApi.getAllSharedCampaigns();
                allCampaigns = res.data?.data || [];
            } catch (e) {
                console.warn('Failed to fetch shared campaigns:', e.message);
            }

            // Filter to only show campaigns assigned to this user
            const assignedCampaignIds = queue.filter(t => t.type === 'PUBLISH').map(t => t.campaignId);
            const myCampaigns = allCampaigns.filter(c => assignedCampaignIds.includes(c.campaignId));

            setSharedCampaigns(myCampaigns);

            // Fetch CFT supply
            try {
                const supplyRes = await fetch(`${API_URL}/auth/cft-supply`, {
                    headers: { 'Authorization': `Bearer ${token}` }
                });
                if (supplyRes.ok) {
                    const supplyData = await supplyRes.json();
                    setCftSupply(supplyData.availableCft || 0);
                }
            } catch (e) {
                console.warn('Failed to fetch CFT supply');
            }
        } catch (err) {
            setSharedCampaigns([]);
            setAssignedQueue([]);
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    // Set CFT Supply (Platform Admin only)
    const handleSetSupply = async () => {
        const amount = parseFloat(newSupplyAmount);
        if (!amount || amount <= 0) {
            alert('Please enter a valid positive amount');
            return;
        }

        setSettingSupply(true);
        try {
            const res = await fetch(`${API_URL}/auth/cft-supply/set`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({ amount })
            });

            const data = await res.json();
            if (res.ok) {
                setCftSupply(data.availableCft);
                setNewSupplyAmount('');
                alert(`CFT supply set to ${amount.toLocaleString()} CFT`);
            } else {
                alert(data.error || 'Failed to set supply');
            }
        } catch (e) {
            alert('Failed to set CFT supply: ' + e.message);
        } finally {
            setSettingSupply(false);
        }
    };

    useEffect(() => {
        if (isPlatformUser) {
            fetchData();
        } else {
            setLoading(false);
        }
    }, [isPlatformUser, user?.orgUserId]);

    const handlePublish = async (campaignId, validationProofHash, startupId) => {
        if (isPreviewMode) {
            alert('Please login as platform admin to publish campaigns');
            return;
        }
        setPublishing(campaignId);
        try {
            const authToken = sessionStorage.getItem('token') || localStorage.getItem('token');
            await platformApi.publishCampaign(campaignId, {
                validationProofHash: validationProofHash || '',
                authToken,
                startupId
            });

            // Mark task as complete in MongoDB queue
            await queueApi.complete({
                campaignId,
                type: 'PUBLISH',
                result: 'PUBLISHED'
            });

            alert('Campaign published successfully!');
            fetchData();
        } catch (err) {
            console.error('Failed to publish:', err);
            alert('Failed to publish: ' + err.message);
        } finally {
            setPublishing(null);
        }
    };

    return (
        <div className="space-y-6">
            {/* Preview Mode Banner */}
            {isPreviewMode && (
                <div className="bg-gradient-to-r from-primary-500 to-primary-600 text-white p-4 rounded-xl flex flex-col sm:flex-row items-center justify-between gap-4">
                    <div className="flex items-center gap-3">
                        <LayoutDashboard size={24} />
                        <div>
                            <h3 className="font-bold">Platform Dashboard Preview</h3>
                            <p className="text-sm opacity-90">Login as platform admin to manage publishing</p>
                        </div>
                    </div>
                    <Link to="/login" state={{ role: 'PLATFORM' }} className="btn bg-white text-primary-700 hover:bg-gray-100 flex items-center gap-2">
                        <LogIn size={18} />
                        Login as Platform
                    </Link>
                </div>
            )}

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
                <div className="flex gap-3 items-center">
                    {isPlatformUser && (
                        <Link to="/wallet" className="flex items-center gap-2 px-4 py-2 bg-primary-100 dark:bg-primary-900/30 text-primary-700 dark:text-primary-300 rounded-lg hover:bg-primary-200">
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

            {/* CFT Supply Management - Platform Admin Only */}
            {isPlatformUser && (
                <div className="card">
                    <div className="flex items-center gap-3 mb-4">
                        <div className="p-2 bg-orange-100 dark:bg-orange-900 rounded-lg">
                            <Settings className="text-orange-600 dark:text-orange-400" size={20} />
                        </div>
                        <h2 className="text-xl font-bold text-gray-900 dark:text-white">
                            CFT Supply Management
                        </h2>
                    </div>

                    <div className="grid md:grid-cols-2 gap-6">
                        {/* Current Supply */}
                        <div className="p-4 bg-gradient-to-br from-green-50 to-green-100 dark:from-green-900/30 dark:to-green-800/30 rounded-xl">
                            <p className="text-sm text-green-700 dark:text-green-300 mb-1">Current Available Supply</p>
                            <p className="text-3xl font-bold text-green-700 dark:text-green-300">
                                {cftSupply.toLocaleString()} CFT
                            </p>
                            <p className="text-xs text-green-600 dark:text-green-400 mt-2">
                                This is the total CFT available for users to purchase
                            </p>
                        </div>

                        {/* Set New Supply */}
                        <div className="space-y-3">
                            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                                Set New CFT Supply
                            </label>
                            <div className="flex gap-2">
                                <input
                                    type="number"
                                    value={newSupplyAmount}
                                    onChange={(e) => setNewSupplyAmount(e.target.value)}
                                    placeholder="e.g., 10000000"
                                    min="0"
                                    className="flex-1 px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800"
                                />
                                <button
                                    onClick={handleSetSupply}
                                    disabled={settingSupply || !newSupplyAmount}
                                    className="btn btn-primary px-6 disabled:opacity-50"
                                >
                                    {settingSupply ? <Loader2 className="animate-spin" size={18} /> : 'Set Supply'}
                                </button>
                            </div>
                            <p className="text-xs text-gray-500">
                                Users can only purchase CFT when supply is available. Enter a positive amount in CFT.
                            </p>
                        </div>
                    </div>
                </div>
            )}

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
                                                        onClick={() => handlePublish(campaign.campaignId, campaign.validationProofHash, campaign.startupId)}
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
