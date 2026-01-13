import { useState, useEffect } from 'react';
import { investorApi } from '../services/api';
import { useAuth } from '../context/AuthContext';
import { TrendingUp, DollarSign, PieChart, Target, Loader2, RefreshCw, Eye, Clock, Briefcase, LogIn, Coins, Wallet } from 'lucide-react';
import { Link } from 'react-router-dom';

// Token constants
const FEES = {
    investmentFeePercent: 5,
    disputeFee: 750
};

export default function InvestorDashboard() {
    const { user, isAuthenticated } = useAuth();
    const [campaigns, setCampaigns] = useState([]);
    const [investments, setInvestments] = useState([]);
    const [viewedCampaigns, setViewedCampaigns] = useState([]);
    const [loading, setLoading] = useState(true);

    const cftBalance = user?.wallet?.cftBalance || 0;
    const isInvestorUser = isAuthenticated && user?.role === 'INVESTOR';
    const isPreviewMode = !isAuthenticated || user?.role !== 'INVESTOR';

    const fetchData = async () => {
        if (!isInvestorUser) return;

        setLoading(true);
        try {
            const investorId = user?.orgUserId;
            const [campaignsRes, investmentsRes, viewedRes] = await Promise.all([
                investorApi.getAvailableCampaigns().catch(() => ({ data: { data: [] } })),
                investorApi.getMyInvestments(investorId).catch(() => ({ data: { data: [] } })),
                investorApi.getViewedCampaigns(investorId).catch(() => ({ data: { data: [] } }))
            ]);
            setCampaigns(campaignsRes.data?.data || []);
            setInvestments(investmentsRes.data?.data || []);
            setViewedCampaigns(viewedRes.data?.data || []);
        } catch (err) {
            setCampaigns([]);
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        if (isInvestorUser) {
            fetchData();
        } else {
            setLoading(false);
        }
    }, [isInvestorUser, user?.orgUserId]);  // Re-fetch when user changes

    const totalInvested = investments.reduce((sum, inv) => sum + (inv.amount || 0), 0);

    return (
        <div className="space-y-6">
            {/* Preview Mode Banner */}
            {isPreviewMode && (
                <div className="bg-gradient-to-r from-accent-500 to-green-600 text-white p-4 rounded-xl flex flex-col sm:flex-row items-center justify-between gap-4">
                    <div className="flex items-center gap-3">
                        <TrendingUp size={24} />
                        <div>
                            <h3 className="font-bold">Investor Dashboard Preview</h3>
                            <p className="text-sm opacity-90">Login as an investor to make investments</p>
                        </div>
                    </div>
                    <Link to="/login" state={{ role: 'INVESTOR' }} className="btn bg-white text-accent-700 hover:bg-gray-100 flex items-center gap-2">
                        <LogIn size={18} />
                        Login as Investor
                    </Link>
                </div>
            )}

            {/* Header */}
            <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
                <div>
                    <h1 className="text-3xl font-bold text-gray-900 dark:text-white">
                        Investor Dashboard
                    </h1>
                    <p className="text-gray-600 dark:text-gray-400 mt-1">
                        Discover and invest in promising startups
                    </p>
                </div>
                <div className="flex gap-3 items-center">
                    {isInvestorUser && (
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

            {/* Stats */}
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                <div className="card flex items-center gap-4">
                    <div className="p-3 bg-primary-100 dark:bg-primary-900 rounded-lg">
                        <Target className="text-primary-600 dark:text-primary-400" size={24} />
                    </div>
                    <div>
                        <h3 className="font-semibold text-gray-900 dark:text-white">Available</h3>
                        <p className="text-2xl font-bold text-primary-600">{campaigns.length}</p>
                    </div>
                </div>
                <div className="card flex items-center gap-4">
                    <div className="p-3 bg-green-100 dark:bg-green-900 rounded-lg">
                        <DollarSign className="text-green-600 dark:text-green-400" size={24} />
                    </div>
                    <div>
                        <h3 className="font-semibold text-gray-900 dark:text-white">Invested</h3>
                        <p className="text-2xl font-bold text-green-600">${totalInvested.toLocaleString()}</p>
                    </div>
                </div>
                <div className="card flex items-center gap-4">
                    <div className="p-3 bg-blue-100 dark:bg-blue-900 rounded-lg">
                        <PieChart className="text-blue-600 dark:text-blue-400" size={24} />
                    </div>
                    <div>
                        <h3 className="font-semibold text-gray-900 dark:text-white">Portfolio</h3>
                        <p className="text-2xl font-bold text-blue-600">{investments.length}</p>
                    </div>
                </div>
                <div className="card flex items-center gap-4">
                    <div className="p-3 bg-accent-100 dark:bg-accent-900 rounded-lg">
                        <Clock className="text-accent-600 dark:text-accent-400" size={24} />
                    </div>
                    <div>
                        <h3 className="font-semibold text-gray-900 dark:text-white">Recent Views</h3>
                        <p className="text-2xl font-bold text-accent-600">{viewedCampaigns.length}</p>
                    </div>
                </div>
            </div>

            {/* My Investments */}
            {investments.length > 0 && (
                <div className="card">
                    <div className="flex items-center gap-2 mb-4">
                        <Briefcase className="text-green-600" size={20} />
                        <h2 className="text-xl font-bold text-gray-900 dark:text-white">My Investments</h2>
                    </div>
                    <div className="overflow-x-auto">
                        <table className="w-full">
                            <thead>
                                <tr className="border-b border-gray-200 dark:border-gray-700">
                                    <th className="text-left py-2 px-4 text-gray-600 dark:text-gray-400">Campaign</th>
                                    <th className="text-left py-2 px-4 text-gray-600 dark:text-gray-400">Amount</th>
                                    <th className="text-left py-2 px-4 text-gray-600 dark:text-gray-400">Date</th>
                                    <th className="text-left py-2 px-4 text-gray-600 dark:text-gray-400">Status</th>
                                </tr>
                            </thead>
                            <tbody>
                                {investments.map((inv, idx) => (
                                    <tr key={idx} className="border-b border-gray-100 dark:border-gray-700">
                                        <td className="py-3 px-4">{inv.campaignId}</td>
                                        <td className="py-3 px-4 font-semibold text-green-600">
                                            {inv.currency || 'USD'} {(inv.amount || 0).toLocaleString()}
                                        </td>
                                        <td className="py-3 px-4 text-gray-500">{inv.investedAt || 'N/A'}</td>
                                        <td className="py-3 px-4">
                                            <span className="badge badge-success">{inv.status || 'ACTIVE'}</span>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                </div>
            )}

            {/* Recent Views */}
            {viewedCampaigns.length > 0 && (
                <div className="card">
                    <div className="flex items-center gap-2 mb-4">
                        <Clock className="text-blue-600" size={20} />
                        <h2 className="text-xl font-bold text-gray-900 dark:text-white">Recently Viewed</h2>
                    </div>
                    <div className="flex flex-wrap gap-2">
                        {viewedCampaigns.slice(0, 5).map((view, idx) => (
                            <Link
                                key={idx}
                                to={`/campaign/${view.campaignId}`}
                                className="px-3 py-2 bg-blue-50 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300 rounded-lg text-sm hover:bg-blue-100 transition"
                            >
                                {view.projectName || view.campaignId}
                            </Link>
                        ))}
                    </div>
                </div>
            )}

            {/* Available Campaigns */}
            <div className="card">
                <div className="flex justify-between items-center mb-4">
                    <h2 className="text-xl font-bold text-gray-900 dark:text-white">Published Campaigns</h2>
                </div>

                {loading ? (
                    <div className="flex justify-center py-12">
                        <Loader2 className="animate-spin text-primary-600" size={40} />
                    </div>
                ) : campaigns.length === 0 ? (
                    <div className="text-center py-12">
                        <TrendingUp size={48} className="mx-auto mb-4 text-gray-300" />
                        <p className="text-gray-500">No campaigns available yet</p>
                        <p className="text-sm text-gray-400 mt-2">Published campaigns will appear here</p>
                    </div>
                ) : (
                    <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
                        {campaigns.map((campaign) => {
                            const name = campaign.projectName || campaign.ProjectName || campaign.project_name || 'Campaign';
                            const id = campaign.campaignId || campaign.CampaignID;
                            const startup = campaign.startupId || campaign.StartupID || 'Startup';
                            const cat = campaign.category || campaign.Category || 'General';
                            const sector = campaign.sector || campaign.Sector;
                            const desc = campaign.description || campaign.Description;
                            const curr = campaign.currency || campaign.Currency || 'USD';
                            const goal = campaign.goalAmount || campaign.GoalAmount || campaign.goal_amount || 0;
                            const risk = campaign.riskLevel || campaign.RiskLevel;
                            const score = campaign.validationScore || campaign.ValidationScore || 0;

                            return (
                                <div key={id} className="border border-gray-200 dark:border-gray-700 rounded-xl p-5 hover:shadow-xl transition-all duration-300 bg-gradient-to-br from-white to-gray-50 dark:from-gray-800 dark:to-gray-900">
                                    <div className="flex justify-between items-start mb-3">
                                        <div className="flex flex-wrap gap-1">
                                            <span className="badge badge-info">{cat}</span>
                                            {sector && <span className="badge badge-secondary">{sector}</span>}
                                        </div>
                                        <span className="badge badge-success">PUBLISHED</span>
                                    </div>

                                    <h3 className="font-bold text-lg text-gray-900 dark:text-white mb-1">{name}</h3>
                                    <p className="text-sm text-gray-500 mb-3">by {startup}</p>

                                    {desc && (
                                        <p className="text-sm text-gray-600 dark:text-gray-400 mb-3 line-clamp-2">{desc}</p>
                                    )}

                                    <div className="bg-gray-100 dark:bg-gray-700/50 rounded-lg p-3 mb-3">
                                        <div className="flex justify-between items-center">
                                            <span className="text-gray-500 text-sm">Goal</span>
                                            <span className="font-bold text-primary-600 dark:text-primary-400">
                                                {curr} {goal.toLocaleString()}
                                            </span>
                                        </div>
                                    </div>

                                    <div className="flex items-center gap-2 mb-3 flex-wrap">
                                        {risk && (
                                            <span className={`text-xs px-2 py-1 rounded-full font-medium ${risk === 'LOW' ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400' :
                                                risk === 'MEDIUM' ? 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400' :
                                                    'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400'
                                                }`}>
                                                {risk} Risk
                                            </span>
                                        )}
                                        {score > 0 && (
                                            <span className="text-xs px-2 py-1 rounded-full bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400 font-medium">
                                                Score: {score}/10
                                            </span>
                                        )}
                                    </div>

                                    <Link
                                        to={`/campaign/${id}`}
                                        className="flex items-center justify-center w-full btn btn-primary text-sm mt-2"
                                    >
                                        <Eye size={16} className="mr-2" />
                                        View Full Details
                                    </Link>
                                </div>
                            );
                        })}
                    </div>
                )}
            </div>
        </div>
    );
}
