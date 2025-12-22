import { useState, useEffect } from 'react';
import { useAuth } from '../context/AuthContext';
import { Link, Navigate } from 'react-router-dom';
import {
    Wallet, Coins, ArrowUpRight, ArrowDownRight, RefreshCw,
    CreditCard, DollarSign, Gift, History, AlertCircle, Loader2,
    TrendingUp, Shield
} from 'lucide-react';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:3001/api';

// Fee schedule - All campaign fees paid by Startup
const FEES = {
    campaignCreation: 10,         // Startup pays to create campaign
    validationSubmission: 50,     // Startup pays to submit to validator
    platformPublishing: 50,       // Startup pays to publish on platform
    // Total campaign journey: 110 CFT
    disputeFee: 750,
    investmentFeePercent: 5,
    withdrawalFeePercent: 1
};

export default function WalletDashboard() {
    const { user, isAuthenticated, token, updateWallet } = useAuth();
    const [loading, setLoading] = useState(false);
    const [purchaseAmount, setPurchaseAmount] = useState('');
    const [purchaseCurrency, setPurchaseCurrency] = useState('INR');
    const [withdrawAmount, setWithdrawAmount] = useState('');
    const [redeemAmount, setRedeemAmount] = useState('');
    const [transactions, setTransactions] = useState([]);
    const [activeTab, setActiveTab] = useState('overview');

    // Exchange rates (from API or constants)
    const exchangeRates = {
        INR: 2.5,
        USD: 83.0
    };

    // Redirect if not authenticated
    if (!isAuthenticated) {
        return <Navigate to="/login" state={{ from: '/wallet' }} replace />;
    }

    const wallet = user?.wallet || { cftBalance: 0, cfrtBalance: 0, frozenCft: 0 };

    const handlePurchase = async () => {
        if (!purchaseAmount || parseFloat(purchaseAmount) <= 0) return;

        setLoading(true);
        try {
            const cftAmount = parseFloat(purchaseAmount) * exchangeRates[purchaseCurrency];

            // In real app: call chaincode TokenContract:PurchaseTokens
            // For now, simulate by updating local wallet
            await updateWallet({
                cftBalance: wallet.cftBalance + cftAmount
            });

            // Add to transaction history
            setTransactions(prev => [{
                id: Date.now(),
                type: 'PURCHASE',
                amount: cftAmount,
                currency: 'CFT',
                fiatAmount: parseFloat(purchaseAmount),
                fiatCurrency: purchaseCurrency,
                timestamp: new Date().toISOString(),
                status: 'COMPLETED'
            }, ...prev]);

            setPurchaseAmount('');
            alert(`Successfully purchased ${cftAmount.toFixed(2)} CFT!`);
        } catch (error) {
            console.error('Purchase failed:', error);
            alert('Purchase failed: ' + error.message);
        } finally {
            setLoading(false);
        }
    };

    const handleWithdraw = async () => {
        if (!withdrawAmount || parseFloat(withdrawAmount) <= 0) return;
        if (parseFloat(withdrawAmount) > wallet.cftBalance) {
            alert('Insufficient CFT balance');
            return;
        }

        setLoading(true);
        try {
            const cftAmount = parseFloat(withdrawAmount);
            const fiatAmount = cftAmount / exchangeRates[purchaseCurrency];
            const fee = fiatAmount * (FEES.withdrawalFeePercent / 100);
            const netAmount = fiatAmount - fee;

            // In real app: call chaincode TokenContract:WithdrawToFiat
            await updateWallet({
                cftBalance: wallet.cftBalance - cftAmount
            });

            setTransactions(prev => [{
                id: Date.now(),
                type: 'WITHDRAWAL',
                amount: -cftAmount,
                currency: 'CFT',
                fiatAmount: netAmount,
                fiatCurrency: purchaseCurrency,
                fee: fee,
                timestamp: new Date().toISOString(),
                status: 'PENDING'
            }, ...prev]);

            setWithdrawAmount('');
            alert(`Withdrawal initiated! Net amount: ${netAmount.toFixed(2)} ${purchaseCurrency}`);
        } catch (error) {
            console.error('Withdrawal failed:', error);
            alert('Withdrawal failed: ' + error.message);
        } finally {
            setLoading(false);
        }
    };

    const handleRedeemRewards = async () => {
        if (!redeemAmount || parseFloat(redeemAmount) <= 0) return;
        if (parseFloat(redeemAmount) > wallet.cfrtBalance) {
            alert('Insufficient CFRT balance');
            return;
        }

        setLoading(true);
        try {
            const cfrtAmount = parseFloat(redeemAmount);
            const cftAmount = cfrtAmount * 10; // 1 CFRT = 10 CFT

            // In real app: call chaincode TokenContract:RedeemRewardTokens
            await updateWallet({
                cfrtBalance: wallet.cfrtBalance - cfrtAmount,
                cftBalance: wallet.cftBalance + cftAmount
            });

            setTransactions(prev => [{
                id: Date.now(),
                type: 'REDEMPTION',
                amount: cftAmount,
                currency: 'CFT',
                cfrtAmount: -cfrtAmount,
                timestamp: new Date().toISOString(),
                status: 'COMPLETED'
            }, ...prev]);

            setRedeemAmount('');
            alert(`Redeemed ${cfrtAmount} CFRT for ${cftAmount} CFT!`);
        } catch (error) {
            console.error('Redemption failed:', error);
            alert('Redemption failed: ' + error.message);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="space-y-6">
            {/* Header */}
            <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
                <div>
                    <h1 className="text-3xl font-bold text-gray-900 dark:text-white flex items-center gap-2">
                        <Wallet className="text-primary-600" />
                        My Wallet
                    </h1>
                    <p className="text-gray-600 dark:text-gray-400 mt-1">
                        Manage your CFT and CFRT tokens
                    </p>
                </div>
                <button
                    onClick={() => updateWallet(wallet)}
                    className="btn btn-secondary flex items-center gap-2"
                >
                    <RefreshCw size={18} />
                    Sync Balance
                </button>
            </div>

            {/* Balance Cards */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                {/* CFT Balance */}
                <div className="card bg-gradient-to-br from-primary-500 to-primary-600 text-white">
                    <div className="flex items-center justify-between mb-4">
                        <Coins size={32} />
                        <span className="text-sm opacity-80">CrowdToken</span>
                    </div>
                    <div className="text-4xl font-bold mb-1">
                        {wallet.cftBalance?.toLocaleString() || 0}
                    </div>
                    <div className="text-primary-100">CFT Balance</div>
                    {wallet.frozenCft > 0 && (
                        <div className="mt-2 text-sm text-primary-200">
                            🔒 {wallet.frozenCft.toLocaleString()} CFT frozen
                        </div>
                    )}
                </div>

                {/* CFRT Balance */}
                <div className="card bg-gradient-to-br from-accent-500 to-accent-600 text-white">
                    <div className="flex items-center justify-between mb-4">
                        <Gift size={32} />
                        <span className="text-sm opacity-80">Reward Token</span>
                    </div>
                    <div className="text-4xl font-bold mb-1">
                        {wallet.cfrtBalance?.toLocaleString() || 0}
                    </div>
                    <div className="text-accent-100">CFRT Balance</div>
                    <div className="mt-2 text-sm text-accent-200">
                        = {((wallet.cfrtBalance || 0) * 10).toLocaleString()} CFT value
                    </div>
                </div>

                {/* ML Rating */}
                <div className="card">
                    <div className="flex items-center justify-between mb-4">
                        <Shield size={32} className="text-purple-500" />
                        <span className={`px-2 py-1 rounded text-xs font-medium ${user?.mlRating?.feeTier === 'TRUSTED' ? 'bg-green-100 text-green-700' :
                            user?.mlRating?.feeTier === 'STANDARD' ? 'bg-blue-100 text-blue-700' :
                                'bg-yellow-100 text-yellow-700'
                            }`}>
                            {user?.mlRating?.feeTier || 'STANDARD'}
                        </span>
                    </div>
                    <div className="text-4xl font-bold text-gray-900 dark:text-white mb-1">
                        {user?.mlRating?.overallScore || 70}
                    </div>
                    <div className="text-gray-500">Trust Score</div>
                </div>
            </div>

            {/* Tabs */}
            <div className="border-b border-gray-200 dark:border-gray-700">
                <nav className="flex gap-4">
                    {['overview', 'purchase', 'withdraw', 'redeem', 'history'].map(tab => (
                        <button
                            key={tab}
                            onClick={() => setActiveTab(tab)}
                            className={`py-3 px-1 border-b-2 font-medium text-sm capitalize ${activeTab === tab
                                ? 'border-primary-500 text-primary-600 dark:text-primary-400'
                                : 'border-transparent text-gray-500 hover:text-gray-700'
                                }`}
                        >
                            {tab}
                        </button>
                    ))}
                </nav>
            </div>

            {/* Tab Content */}
            <div className="card">
                {activeTab === 'overview' && (
                    <div className="space-y-6">
                        <h3 className="text-lg font-bold text-gray-900 dark:text-white">Fee Schedule</h3>
                        <div className="grid md:grid-cols-2 gap-4">
                            {Object.entries(FEES).map(([key, value]) => (
                                <div key={key} className="flex justify-between items-center p-3 bg-gray-50 dark:bg-gray-800 rounded-lg">
                                    <span className="text-gray-600 dark:text-gray-400 capitalize">
                                        {key.replace(/([A-Z])/g, ' $1').replace('Fee', '')}
                                    </span>
                                    <span className="font-medium text-gray-900 dark:text-white">
                                        {typeof value === 'number' && value < 100 ? `${value}%` : `${value} CFT`}
                                    </span>
                                </div>
                            ))}
                        </div>
                        <div className="p-4 bg-blue-50 dark:bg-blue-900/30 rounded-lg">
                            <p className="text-sm text-blue-700 dark:text-blue-300">
                                <strong>Exchange Rate:</strong> 1 INR = {exchangeRates.INR} CFT | 1 USD = {exchangeRates.USD} CFT
                            </p>
                        </div>
                    </div>
                )}

                {activeTab === 'purchase' && (
                    <div className="space-y-6 max-w-md">
                        <h3 className="text-lg font-bold text-gray-900 dark:text-white flex items-center gap-2">
                            <CreditCard className="text-primary-500" />
                            Purchase CFT Tokens
                        </h3>
                        <div className="space-y-4">
                            <div>
                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                                    Amount
                                </label>
                                <div className="flex gap-2">
                                    <input
                                        type="number"
                                        value={purchaseAmount}
                                        onChange={(e) => setPurchaseAmount(e.target.value)}
                                        placeholder="Enter amount"
                                        className="flex-1 px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800"
                                    />
                                    <select
                                        value={purchaseCurrency}
                                        onChange={(e) => setPurchaseCurrency(e.target.value)}
                                        className="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800"
                                    >
                                        <option value="INR">INR</option>
                                        <option value="USD">USD</option>
                                    </select>
                                </div>
                            </div>
                            {purchaseAmount && (
                                <div className="p-4 bg-gray-50 dark:bg-gray-800 rounded-lg">
                                    <p className="text-sm text-gray-600 dark:text-gray-400">You will receive:</p>
                                    <p className="text-2xl font-bold text-primary-600">
                                        {(parseFloat(purchaseAmount) * exchangeRates[purchaseCurrency]).toFixed(2)} CFT
                                    </p>
                                </div>
                            )}
                            <button
                                onClick={handlePurchase}
                                disabled={loading || !purchaseAmount}
                                className="w-full btn btn-primary flex items-center justify-center gap-2"
                            >
                                {loading ? <Loader2 className="animate-spin" size={18} /> : <ArrowUpRight size={18} />}
                                Purchase CFT
                            </button>
                        </div>
                    </div>
                )}

                {activeTab === 'withdraw' && (
                    <div className="space-y-6 max-w-md">
                        <h3 className="text-lg font-bold text-gray-900 dark:text-white flex items-center gap-2">
                            <ArrowDownRight className="text-accent-500" />
                            Withdraw to Fiat
                        </h3>
                        <div className="p-3 bg-yellow-50 dark:bg-yellow-900/30 rounded-lg flex items-start gap-2">
                            <AlertCircle size={18} className="text-yellow-600 mt-0.5" />
                            <p className="text-sm text-yellow-700 dark:text-yellow-300">
                                {FEES.withdrawalFeePercent}% withdrawal fee applies
                            </p>
                        </div>
                        <div className="space-y-4">
                            <div>
                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                                    CFT Amount
                                </label>
                                <input
                                    type="number"
                                    value={withdrawAmount}
                                    onChange={(e) => setWithdrawAmount(e.target.value)}
                                    placeholder="Enter CFT amount"
                                    max={wallet.cftBalance}
                                    className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800"
                                />
                            </div>
                            {withdrawAmount && (
                                <div className="p-4 bg-gray-50 dark:bg-gray-800 rounded-lg space-y-1">
                                    <p className="text-sm text-gray-600 dark:text-gray-400">
                                        Gross: {(parseFloat(withdrawAmount) / exchangeRates[purchaseCurrency]).toFixed(2)} {purchaseCurrency}
                                    </p>
                                    <p className="text-sm text-gray-600 dark:text-gray-400">
                                        Fee ({FEES.withdrawalFeePercent}%): {((parseFloat(withdrawAmount) / exchangeRates[purchaseCurrency]) * 0.01).toFixed(2)} {purchaseCurrency}
                                    </p>
                                    <p className="text-lg font-bold text-accent-600">
                                        Net: {((parseFloat(withdrawAmount) / exchangeRates[purchaseCurrency]) * 0.99).toFixed(2)} {purchaseCurrency}
                                    </p>
                                </div>
                            )}
                            <button
                                onClick={handleWithdraw}
                                disabled={loading || !withdrawAmount || parseFloat(withdrawAmount) > wallet.cftBalance}
                                className="w-full btn btn-primary flex items-center justify-center gap-2"
                            >
                                {loading ? <Loader2 className="animate-spin" size={18} /> : <ArrowDownRight size={18} />}
                                Withdraw
                            </button>
                        </div>
                    </div>
                )}

                {activeTab === 'redeem' && (
                    <div className="space-y-6 max-w-md">
                        <h3 className="text-lg font-bold text-gray-900 dark:text-white flex items-center gap-2">
                            <Gift className="text-purple-500" />
                            Redeem CFRT Rewards
                        </h3>
                        <div className="p-3 bg-purple-50 dark:bg-purple-900/30 rounded-lg">
                            <p className="text-sm text-purple-700 dark:text-purple-300">
                                <strong>Rate:</strong> 1 CFRT = 10 CFT
                            </p>
                        </div>
                        <div className="space-y-4">
                            <div>
                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                                    CFRT Amount (Available: {wallet.cfrtBalance || 0})
                                </label>
                                <input
                                    type="number"
                                    value={redeemAmount}
                                    onChange={(e) => setRedeemAmount(e.target.value)}
                                    placeholder="Enter CFRT amount"
                                    max={wallet.cfrtBalance}
                                    className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800"
                                />
                            </div>
                            {redeemAmount && (
                                <div className="p-4 bg-gray-50 dark:bg-gray-800 rounded-lg">
                                    <p className="text-sm text-gray-600 dark:text-gray-400">You will receive:</p>
                                    <p className="text-2xl font-bold text-purple-600">
                                        {(parseFloat(redeemAmount) * 10).toFixed(0)} CFT
                                    </p>
                                </div>
                            )}
                            <button
                                onClick={handleRedeemRewards}
                                disabled={loading || !redeemAmount || parseFloat(redeemAmount) > wallet.cfrtBalance}
                                className="w-full btn btn-primary flex items-center justify-center gap-2"
                            >
                                {loading ? <Loader2 className="animate-spin" size={18} /> : <Gift size={18} />}
                                Redeem CFRT
                            </button>
                        </div>
                    </div>
                )}

                {activeTab === 'history' && (
                    <div className="space-y-4">
                        <h3 className="text-lg font-bold text-gray-900 dark:text-white flex items-center gap-2">
                            <History className="text-gray-500" />
                            Transaction History
                        </h3>
                        {transactions.length === 0 ? (
                            <div className="text-center py-12 text-gray-500">
                                <History size={48} className="mx-auto mb-4 opacity-30" />
                                <p>No transactions yet</p>
                            </div>
                        ) : (
                            <div className="space-y-3">
                                {transactions.map(tx => (
                                    <div key={tx.id} className="flex items-center justify-between p-4 bg-gray-50 dark:bg-gray-800 rounded-lg">
                                        <div className="flex items-center gap-3">
                                            <div className={`p-2 rounded-lg ${tx.type === 'PURCHASE' ? 'bg-green-100 text-green-600' :
                                                tx.type === 'WITHDRAWAL' ? 'bg-red-100 text-red-600' :
                                                    'bg-purple-100 text-purple-600'
                                                }`}>
                                                {tx.type === 'PURCHASE' ? <ArrowUpRight size={18} /> :
                                                    tx.type === 'WITHDRAWAL' ? <ArrowDownRight size={18} /> :
                                                        <Gift size={18} />}
                                            </div>
                                            <div>
                                                <p className="font-medium text-gray-900 dark:text-white capitalize">
                                                    {tx.type.toLowerCase()}
                                                </p>
                                                <p className="text-xs text-gray-500">
                                                    {new Date(tx.timestamp).toLocaleString()}
                                                </p>
                                            </div>
                                        </div>
                                        <div className="text-right">
                                            <p className={`font-bold ${tx.amount > 0 ? 'text-green-600' : 'text-red-600'}`}>
                                                {tx.amount > 0 ? '+' : ''}{tx.amount.toFixed(2)} CFT
                                            </p>
                                            <span className={`text-xs px-2 py-0.5 rounded ${tx.status === 'COMPLETED' ? 'bg-green-100 text-green-700' : 'bg-yellow-100 text-yellow-700'
                                                }`}>
                                                {tx.status}
                                            </span>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>
                )}
            </div>
        </div>
    );
}
