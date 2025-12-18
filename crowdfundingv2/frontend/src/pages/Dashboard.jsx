import { Link } from 'react-router-dom';
import { Rocket, Shield, LayoutDashboard, TrendingUp, ArrowRight, CheckCircle } from 'lucide-react';

const roles = [
    {
        path: '/startup',
        icon: Rocket,
        title: 'Startup',
        description: 'Create and manage crowdfunding campaigns',
        color: 'from-blue-500 to-blue-600',
        features: ['Create Campaigns', 'Track Investments', 'Manage Milestones'],
    },
    {
        path: '/validator',
        icon: Shield,
        title: 'Validator',
        description: 'Verify and validate campaign legitimacy',
        color: 'from-purple-500 to-purple-600',
        features: ['Review Campaigns', 'Risk Assessment', 'Approve/Reject'],
    },
    {
        path: '/platform',
        icon: LayoutDashboard,
        title: 'Platform',
        description: 'Publish campaigns and manage platform',
        color: 'from-primary-500 to-primary-600',
        features: ['Publish Campaigns', 'Manage Fees', 'Handle Disputes'],
    },
    {
        path: '/investor',
        icon: TrendingUp,
        title: 'Investor',
        description: 'Browse and invest in campaigns',
        color: 'from-accent-500 to-accent-600',
        features: ['Browse Campaigns', 'Make Investments', 'Track Portfolio'],
    },
];

export default function Dashboard() {
    return (
        <div className="space-y-8">
            {/* Hero Section */}
            <div className="text-center py-12 bg-gradient-to-br from-primary-500 via-primary-600 to-accent-600 rounded-2xl shadow-xl">
                <h1 className="text-4xl md:text-5xl font-bold text-white mb-4">
                    Welcome to CrowdFundChain
                </h1>
                <p className="text-xl text-primary-100 max-w-2xl mx-auto mb-8">
                    Blockchain-powered crowdfunding platform built on Hyperledger Fabric.
                    Secure, transparent, and trustworthy.
                </p>
                <div className="flex justify-center gap-4">
                    <Link to="/investor" className="btn bg-white text-primary-700 hover:bg-gray-100">
                        Start Investing
                    </Link>
                    <Link to="/startup" className="btn bg-primary-800 text-white hover:bg-primary-900">
                        Launch Campaign
                    </Link>
                </div>
            </div>

            {/* Role Selection */}
            <div>
                <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-6 text-center">
                    Select Your Role
                </h2>
                <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-6">
                    {roles.map(({ path, icon: Icon, title, description, color, features }) => (
                        <Link
                            key={path}
                            to={path}
                            className="card hover:shadow-xl transition-all duration-300 hover:-translate-y-1 group"
                        >
                            <div className={`w-14 h-14 rounded-xl bg-gradient-to-br ${color} flex items-center justify-center mb-4`}>
                                <Icon className="text-white" size={28} />
                            </div>
                            <h3 className="text-xl font-bold text-gray-900 dark:text-white mb-2">
                                {title}
                            </h3>
                            <p className="text-gray-600 dark:text-gray-400 mb-4 text-sm">
                                {description}
                            </p>
                            <ul className="space-y-2 mb-4">
                                {features.map((feature) => (
                                    <li key={feature} className="flex items-center text-sm text-gray-500 dark:text-gray-400">
                                        <CheckCircle size={14} className="text-accent-500 mr-2" />
                                        {feature}
                                    </li>
                                ))}
                            </ul>
                            <div className="flex items-center text-primary-600 dark:text-primary-400 font-medium group-hover:translate-x-1 transition-transform">
                                Enter <ArrowRight size={16} className="ml-1" />
                            </div>
                        </Link>
                    ))}
                </div>
            </div>

            {/* Stats Section */}
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                {[
                    { label: 'Total Campaigns', value: '150+' },
                    { label: 'Funds Raised', value: '$2.5M' },
                    { label: 'Investors', value: '1,200+' },
                    { label: 'Success Rate', value: '94%' },
                ].map(({ label, value }) => (
                    <div key={label} className="card text-center">
                        <div className="text-3xl font-bold text-primary-600 dark:text-primary-400">{value}</div>
                        <div className="text-sm text-gray-500 dark:text-gray-400">{label}</div>
                    </div>
                ))}
            </div>
        </div>
    );
}
