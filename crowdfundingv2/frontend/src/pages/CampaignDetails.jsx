import { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { investorApi } from '../services/api';
import { ArrowLeft, DollarSign, Calendar, Tag, FileText, Loader2, Shield, Check } from 'lucide-react';

export default function CampaignDetails() {
    const { campaignId } = useParams();
    const [campaign, setCampaign] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [investAmount, setInvestAmount] = useState('');
    const [investing, setInvesting] = useState(false);
    const [requesting, setRequesting] = useState(false);
    const [requestSent, setRequestSent] = useState(false);

    useEffect(() => {
        const fetchCampaign = async () => {
            setLoading(true);
            try {
                const res = await investorApi.viewCampaignDetails(campaignId);
                setCampaign(res.data?.data);
                setError(null);

                // Auto-record this view in investor's private data
                investorApi.viewCampaign(campaignId, { investorId: 'INVESTOR001' }).catch(() => { });
            } catch (err) {
                setError('Failed to fetch campaign details');
                console.error(err);
            } finally {
                setLoading(false);
            }
        };

        fetchCampaign();
    }, [campaignId]);

    const handleInvest = async () => {
        if (!investAmount || parseFloat(investAmount) <= 0) return;

        setInvesting(true);
        try {
            await investorApi.makeInvestment({
                investmentId: `INV_${Date.now()}`,
                campaignId,
                investorId: 'INVESTOR001',
                amount: parseFloat(investAmount),
                currency: 'USD',
            });
            alert('Investment successful! Your investment has been recorded.');
            setInvestAmount('');
        } catch (err) {
            console.error('Investment failed:', err);
            alert('Investment failed: ' + (err.response?.data?.error || err.message));
        } finally {
            setInvesting(false);
        }
    };

    const handleRequestValidation = async () => {
        setRequesting(true);
        try {
            await investorApi.requestValidationDetails(campaignId, { investorId: 'INVESTOR001' });
            setRequestSent(true);
            alert('Validation request sent! The validator will respond with detailed risk analysis.');
        } catch (err) {
            console.error('Request failed:', err);
            alert('Request failed: ' + (err.response?.data?.error || err.message));
        } finally {
            setRequesting(false);
        }
    };

    if (loading) {
        return (
            <div className="flex justify-center py-24">
                <Loader2 className="animate-spin text-primary-600" size={48} />
            </div>
        );
    }

    if (error || !campaign) {
        return (
            <div className="text-center py-24">
                <p className="text-red-500">{error || 'Campaign not found'}</p>
                <Link to="/investor" className="btn btn-secondary mt-4">
                    Back to Dashboard
                </Link>
            </div>
        );
    }

    return (
        <div className="space-y-6">
            {/* Back Link */}
            <Link to="/investor" className="flex items-center text-gray-600 dark:text-gray-400 hover:text-primary-600">
                <ArrowLeft size={20} className="mr-2" />
                Back to Campaigns
            </Link>

            {/* Header */}
            <div className="card">
                <div className="flex flex-col md:flex-row justify-between gap-6">
                    <div className="flex-1">
                        <div className="flex items-center gap-3 mb-2">
                            <span className="badge badge-info">{campaign.category}</span>
                            <span className={`badge ${campaign.status === 'PUBLISHED' ? 'badge-success' : 'badge-warning'}`}>
                                {campaign.status}
                            </span>
                            {campaign.riskLevel && (
                                <span className={`badge ${campaign.riskLevel === 'LOW' ? 'badge-success' : 'badge-warning'}`}>
                                    Risk: {campaign.riskLevel}
                                </span>
                            )}
                        </div>
                        <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">
                            {campaign.projectName || campaign.project_name || 'Untitled Campaign'}
                        </h1>
                        <p className="text-gray-600 dark:text-gray-400 mb-4">
                            by {campaign.startupId}
                        </p>
                        <p className="text-gray-700 dark:text-gray-300">
                            {campaign.description}
                        </p>
                    </div>

                    {/* Investment Card */}
                    <div className="w-full md:w-80 bg-gray-50 dark:bg-gray-700/50 rounded-xl p-6">
                        <div className="text-center mb-6">
                            <p className="text-sm text-gray-500 dark:text-gray-400">Goal Amount</p>
                            <p className="text-3xl font-bold text-primary-600">
                                ${(campaign.goalAmount || campaign.goal_amount || 0).toLocaleString()}
                            </p>
                        </div>

                        <div className="space-y-4">
                            <div>
                                <label className="label">Investment Amount (USD)</label>
                                <input
                                    type="number"
                                    value={investAmount}
                                    onChange={(e) => setInvestAmount(e.target.value)}
                                    placeholder="Enter amount"
                                    className="input"
                                />
                            </div>
                            <button
                                onClick={handleInvest}
                                disabled={investing || !investAmount}
                                className="btn btn-primary w-full flex items-center justify-center gap-2"
                            >
                                {investing ? <Loader2 className="animate-spin" size={18} /> : <DollarSign size={18} />}
                                {investing ? 'Processing...' : 'Invest Now'}
                            </button>
                        </div>
                    </div>
                </div>
            </div>

            {/* Details Grid */}
            <div className="grid md:grid-cols-2 gap-6">
                <div className="card">
                    <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4 flex items-center gap-2">
                        <Calendar size={20} className="text-primary-600" />
                        Campaign Details
                    </h3>
                    <dl className="space-y-3">
                        <div className="flex justify-between">
                            <dt className="text-gray-500">Deadline</dt>
                            <dd className="text-gray-900 dark:text-white font-medium">{campaign.deadline || 'N/A'}</dd>
                        </div>
                        <div className="flex justify-between">
                            <dt className="text-gray-500">Currency</dt>
                            <dd className="text-gray-900 dark:text-white font-medium">{campaign.currency}</dd>
                        </div>
                        <div className="flex justify-between">
                            <dt className="text-gray-500">Validation Score</dt>
                            <dd className="text-gray-900 dark:text-white font-medium">{campaign.validationScore || 'N/A'}</dd>
                        </div>
                    </dl>
                </div>

                <div className="card">
                    <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4 flex items-center gap-2">
                        <Tag size={20} className="text-primary-600" />
                        Tags
                    </h3>
                    <div className="flex flex-wrap gap-2">
                        {(campaign.tags || []).map((tag) => (
                            <span key={tag} className="badge bg-primary-100 text-primary-700 dark:bg-primary-900 dark:text-primary-300">
                                {tag}
                            </span>
                        ))}
                    </div>
                </div>

                <div className="card">
                    <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4 flex items-center gap-2">
                        <FileText size={20} className="text-primary-600" />
                        Documents
                    </h3>
                    <ul className="space-y-2">
                        {(campaign.documents || []).map((doc) => (
                            <li key={doc} className="flex items-center text-gray-600 dark:text-gray-400">
                                <FileText size={16} className="mr-2" />
                                {doc}
                            </li>
                        ))}
                    </ul>
                </div>

                <div className="card">
                    <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4 flex items-center gap-2">
                        <Shield size={20} className="text-primary-600" />
                        Request Validation Info
                    </h3>
                    <p className="text-gray-600 dark:text-gray-400 mb-4">
                        Request detailed validation and risk analysis from the validator.
                    </p>
                    {requestSent ? (
                        <div className="flex items-center text-green-600">
                            <Check size={20} className="mr-2" />
                            Request Sent - Awaiting Validator Response
                        </div>
                    ) : (
                        <button
                            onClick={handleRequestValidation}
                            disabled={requesting}
                            className="btn btn-secondary flex items-center gap-2"
                        >
                            {requesting ? <Loader2 className="animate-spin" size={18} /> : <Shield size={18} />}
                            {requesting ? 'Sending...' : 'Request Validation Details'}
                        </button>
                    )}
                </div>
            </div>
        </div>
    );
}
