import { useState } from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { UserPlus, Mail, Lock, User, ArrowLeft, Rocket, Shield, LayoutDashboard, TrendingUp, Check } from 'lucide-react';

const roles = [
    { id: 'STARTUP', icon: Rocket, label: 'Startup', description: 'Create campaigns', color: 'from-blue-500 to-blue-600' },
    { id: 'INVESTOR', icon: TrendingUp, label: 'Investor', description: 'Invest in campaigns', color: 'from-accent-500 to-accent-600' },
    { id: 'VALIDATOR', icon: Shield, label: 'Validator', description: 'Validate campaigns', color: 'from-purple-500 to-purple-600' },
    { id: 'PLATFORM', icon: LayoutDashboard, label: 'Platform', description: 'Manage platform', color: 'from-primary-500 to-primary-600' }
];

export default function Signup() {
    const { signup } = useAuth();
    const navigate = useNavigate();
    const location = useLocation();

    const [formData, setFormData] = useState({
        name: '',
        email: '',
        password: '',
        confirmPassword: ''
    });
    const [selectedRole, setSelectedRole] = useState(location.state?.role || '');
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);

    const handleChange = (e) => {
        setFormData(prev => ({
            ...prev,
            [e.target.name]: e.target.value
        }));
        setError('');
    };

    const handleSubmit = async (e) => {
        e.preventDefault();

        if (!selectedRole) {
            setError('Please select a role');
            return;
        }

        if (formData.password !== formData.confirmPassword) {
            setError('Passwords do not match');
            return;
        }

        if (formData.password.length < 6) {
            setError('Password must be at least 6 characters');
            return;
        }

        setLoading(true);
        setError('');

        try {
            const user = await signup(formData.name, formData.email, formData.password, selectedRole);
            navigate(`/${user.role.toLowerCase()}`);
        } catch (err) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="min-h-[80vh] flex items-center justify-center py-8">
            <div className="w-full max-w-lg">
                <div className="card p-8">
                    {/* Header */}
                    <div className="text-center mb-6">
                        <div className="w-16 h-16 mx-auto rounded-2xl bg-gradient-to-br from-primary-500 to-accent-500 flex items-center justify-center mb-4">
                            <UserPlus className="text-white" size={32} />
                        </div>
                        <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
                            Create Account
                        </h1>
                        <p className="text-gray-600 dark:text-gray-400 mt-1">
                            Join CrowdFundChain today
                        </p>
                    </div>

                    {/* Role Selection */}
                    <div className="mb-6">
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                            Select Your Role
                        </label>
                        <div className="grid grid-cols-2 gap-3">
                            {roles.map(({ id, icon: Icon, label, description, color }) => (
                                <button
                                    key={id}
                                    type="button"
                                    onClick={() => setSelectedRole(id)}
                                    className={`relative p-4 rounded-xl border-2 transition-all text-left ${selectedRole === id
                                            ? 'border-primary-500 bg-primary-50 dark:bg-primary-900/20'
                                            : 'border-gray-200 dark:border-gray-700 hover:border-gray-300 dark:hover:border-gray-600'
                                        }`}
                                >
                                    {selectedRole === id && (
                                        <div className="absolute top-2 right-2">
                                            <Check className="text-primary-500" size={16} />
                                        </div>
                                    )}
                                    <div className={`w-10 h-10 rounded-lg bg-gradient-to-br ${color} flex items-center justify-center mb-2`}>
                                        <Icon className="text-white" size={20} />
                                    </div>
                                    <div className="font-medium text-gray-900 dark:text-white">{label}</div>
                                    <div className="text-xs text-gray-500 dark:text-gray-400">{description}</div>
                                </button>
                            ))}
                        </div>
                    </div>

                    {/* Error Message */}
                    {error && (
                        <div className="mb-4 p-3 bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-300 rounded-lg text-sm">
                            {error}
                        </div>
                    )}

                    {/* Signup Form */}
                    <form onSubmit={handleSubmit} className="space-y-4">
                        <div>
                            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                                Full Name
                            </label>
                            <div className="relative">
                                <User className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" size={18} />
                                <input
                                    type="text"
                                    name="name"
                                    value={formData.name}
                                    onChange={handleChange}
                                    required
                                    className="w-full pl-10 pr-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500"
                                    placeholder="John Doe"
                                />
                            </div>
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                                Email
                            </label>
                            <div className="relative">
                                <Mail className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" size={18} />
                                <input
                                    type="email"
                                    name="email"
                                    value={formData.email}
                                    onChange={handleChange}
                                    required
                                    className="w-full pl-10 pr-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500"
                                    placeholder="you@example.com"
                                />
                            </div>
                        </div>

                        <div className="grid grid-cols-2 gap-4">
                            <div>
                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                                    Password
                                </label>
                                <div className="relative">
                                    <Lock className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" size={18} />
                                    <input
                                        type="password"
                                        name="password"
                                        value={formData.password}
                                        onChange={handleChange}
                                        required
                                        className="w-full pl-10 pr-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500"
                                        placeholder="••••••"
                                    />
                                </div>
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                                    Confirm
                                </label>
                                <div className="relative">
                                    <Lock className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" size={18} />
                                    <input
                                        type="password"
                                        name="confirmPassword"
                                        value={formData.confirmPassword}
                                        onChange={handleChange}
                                        required
                                        className="w-full pl-10 pr-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500"
                                        placeholder="••••••"
                                    />
                                </div>
                            </div>
                        </div>

                        <button
                            type="submit"
                            disabled={loading || !selectedRole}
                            className="w-full py-3 rounded-lg font-medium text-white bg-gradient-to-br from-primary-500 to-accent-500 hover:opacity-90 transition-opacity disabled:opacity-50"
                        >
                            {loading ? 'Creating account...' : 'Create Account'}
                        </button>
                    </form>

                    {/* Footer */}
                    <div className="mt-6 text-center">
                        <p className="text-gray-600 dark:text-gray-400">
                            Already have an account?{' '}
                            <Link to="/login" className="text-primary-600 dark:text-primary-400 font-medium hover:underline">
                                Sign in
                            </Link>
                        </p>
                    </div>

                    <div className="mt-4">
                        <Link to="/" className="flex items-center justify-center text-gray-500 dark:text-gray-400 hover:text-gray-700">
                            <ArrowLeft size={16} className="mr-1" />
                            Back to Home
                        </Link>
                    </div>
                </div>
            </div>
        </div>
    );
}
