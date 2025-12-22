import { useState } from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { LogIn, Mail, Lock, ArrowLeft, Rocket, Shield, LayoutDashboard, TrendingUp } from 'lucide-react';

const roleIcons = {
    STARTUP: Rocket,
    VALIDATOR: Shield,
    PLATFORM: LayoutDashboard,
    INVESTOR: TrendingUp
};

const roleColors = {
    STARTUP: 'from-blue-500 to-blue-600',
    VALIDATOR: 'from-purple-500 to-purple-600',
    PLATFORM: 'from-primary-500 to-primary-600',
    INVESTOR: 'from-accent-500 to-accent-600'
};

export default function Login() {
    const { login } = useAuth();
    const navigate = useNavigate();
    const location = useLocation();

    const [formData, setFormData] = useState({
        email: '',
        password: ''
    });
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);

    // Get intended role from URL state
    const intendedRole = location.state?.role;
    const from = location.state?.from || '/';

    const handleChange = (e) => {
        setFormData(prev => ({
            ...prev,
            [e.target.name]: e.target.value
        }));
        setError('');
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setLoading(true);
        setError('');

        try {
            const user = await login(formData.email, formData.password);

            // Redirect to user's role dashboard
            const dashboardPath = `/${user.role.toLowerCase()}`;
            navigate(dashboardPath);

        } catch (err) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    };

    const RoleIcon = intendedRole ? roleIcons[intendedRole] : LogIn;
    const roleColor = intendedRole ? roleColors[intendedRole] : 'from-primary-500 to-primary-600';

    return (
        <div className="min-h-[80vh] flex items-center justify-center">
            <div className="w-full max-w-md">
                <div className="card p-8">
                    {/* Header */}
                    <div className="text-center mb-8">
                        <div className={`w-16 h-16 mx-auto rounded-2xl bg-gradient-to-br ${roleColor} flex items-center justify-center mb-4`}>
                            <RoleIcon className="text-white" size={32} />
                        </div>
                        <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
                            Welcome Back
                        </h1>
                        {intendedRole && (
                            <p className="text-gray-600 dark:text-gray-400 mt-1">
                                Login to access {intendedRole.toLowerCase()} dashboard
                            </p>
                        )}
                    </div>

                    {/* Error Message */}
                    {error && (
                        <div className="mb-4 p-3 bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-300 rounded-lg text-sm">
                            {error}
                        </div>
                    )}

                    {/* Login Form */}
                    <form onSubmit={handleSubmit} className="space-y-4">
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
                                    className="w-full pl-10 pr-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                                    placeholder="you@example.com"
                                />
                            </div>
                        </div>

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
                                    className="w-full pl-10 pr-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                                    placeholder="••••••••"
                                />
                            </div>
                        </div>

                        <button
                            type="submit"
                            disabled={loading}
                            className={`w-full py-3 rounded-lg font-medium text-white bg-gradient-to-br ${roleColor} hover:opacity-90 transition-opacity disabled:opacity-50`}
                        >
                            {loading ? 'Logging in...' : 'Login'}
                        </button>
                    </form>

                    {/* Footer */}
                    <div className="mt-6 text-center">
                        <p className="text-gray-600 dark:text-gray-400">
                            Don't have an account?{' '}
                            <Link
                                to="/signup"
                                state={{ role: intendedRole }}
                                className="text-primary-600 dark:text-primary-400 font-medium hover:underline"
                            >
                                Sign up
                            </Link>
                        </p>
                    </div>

                    <div className="mt-4">
                        <Link
                            to="/"
                            className="flex items-center justify-center text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200"
                        >
                            <ArrowLeft size={16} className="mr-1" />
                            Back to Home
                        </Link>
                    </div>
                </div>
            </div>
        </div>
    );
}
