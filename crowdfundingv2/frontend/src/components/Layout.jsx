import { Outlet, Link, useLocation, useNavigate } from 'react-router-dom';
import { useTheme } from '../context/ThemeContext';
import { useAuth } from '../context/AuthContext';
import {
    Sun, Moon, Home, Rocket, Shield, LayoutDashboard,
    TrendingUp, Menu, X, LogIn, LogOut, User, Wallet
} from 'lucide-react';
import { useState } from 'react';

// Navigation items by role
const roleNavItems = {
    STARTUP: [
        { path: '/startup', icon: Rocket, label: 'Dashboard' },
    ],
    INVESTOR: [
        { path: '/investor', icon: TrendingUp, label: 'Dashboard' },
    ],
    VALIDATOR: [
        { path: '/validator', icon: Shield, label: 'Dashboard' },
    ],
    PLATFORM: [
        { path: '/platform', icon: LayoutDashboard, label: 'Dashboard' },
    ]
};

// Guest navigation (can preview all dashboards)
const guestNavItems = [
    { path: '/', icon: Home, label: 'Home' },
    { path: '/startup', icon: Rocket, label: 'Startup' },
    { path: '/validator', icon: Shield, label: 'Validator' },
    { path: '/platform', icon: LayoutDashboard, label: 'Platform' },
    { path: '/investor', icon: TrendingUp, label: 'Investor' },
];

export default function Layout() {
    const { darkMode, toggleDarkMode } = useTheme();
    const { user, isAuthenticated, logout, role } = useAuth();
    const location = useLocation();
    const navigate = useNavigate();
    const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

    // Select nav items based on auth status
    const navItems = isAuthenticated
        ? [{ path: '/', icon: Home, label: 'Home' }, ...roleNavItems[role]]
        : guestNavItems;

    const handleLogout = () => {
        logout();
        navigate('/');
    };

    return (
        <div className="min-h-screen bg-gray-50 dark:bg-gray-900 transition-colors duration-200">
            {/* Header */}
            <header className="bg-white dark:bg-gray-800 shadow-sm border-b border-gray-200 dark:border-gray-700 sticky top-0 z-50">
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                    <div className="flex justify-between items-center h-16">
                        {/* Logo */}
                        <Link to="/" className="flex items-center space-x-2">
                            <div className="w-10 h-10 bg-gradient-to-br from-primary-500 to-accent-500 rounded-lg flex items-center justify-center">
                                <span className="text-white font-bold text-xl">C</span>
                            </div>
                            <span className="text-xl font-bold text-gray-900 dark:text-white hidden sm:block">
                                CrowdFund<span className="text-primary-600">Chain</span>
                            </span>
                        </Link>

                        {/* Desktop Navigation */}
                        <nav className="hidden md:flex items-center space-x-1">
                            {navItems.map(({ path, icon: Icon, label }) => {
                                const isActive = location.pathname === path ||
                                    (path !== '/' && location.pathname.startsWith(path));
                                return (
                                    <Link
                                        key={path}
                                        to={path}
                                        className={`flex items-center space-x-2 px-4 py-2 rounded-lg transition-all duration-200 ${isActive
                                            ? 'bg-primary-100 dark:bg-primary-900/50 text-primary-700 dark:text-primary-300'
                                            : 'text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700'
                                            }`}
                                    >
                                        <Icon size={18} />
                                        <span className="font-medium">{label}</span>
                                    </Link>
                                );
                            })}
                        </nav>

                        {/* Right side */}
                        <div className="flex items-center space-x-3">
                            {/* Wallet Balance (authenticated only) */}
                            {isAuthenticated && user?.wallet && (
                                <div className="hidden sm:flex items-center space-x-2 px-3 py-1.5 bg-accent-100 dark:bg-accent-900/30 rounded-lg">
                                    <Wallet size={16} className="text-accent-600 dark:text-accent-400" />
                                    <span className="text-sm font-medium text-accent-700 dark:text-accent-300">
                                        {user.wallet.cftBalance?.toLocaleString() || 0} CFT
                                    </span>
                                </div>
                            )}

                            {/* Auth Buttons */}
                            {isAuthenticated ? (
                                <div className="flex items-center space-x-2">
                                    <div className="hidden sm:flex items-center space-x-2 px-3 py-1.5 bg-gray-100 dark:bg-gray-700 rounded-lg">
                                        <User size={16} className="text-gray-500" />
                                        <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
                                            {user?.name?.split(' ')[0]}
                                        </span>
                                        <span className="text-xs px-1.5 py-0.5 bg-primary-100 dark:bg-primary-900 text-primary-700 dark:text-primary-300 rounded">
                                            {role}
                                        </span>
                                    </div>
                                    <button
                                        onClick={handleLogout}
                                        className="p-2 rounded-lg text-gray-500 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
                                        title="Logout"
                                    >
                                        <LogOut size={20} />
                                    </button>
                                </div>
                            ) : (
                                <Link
                                    to="/login"
                                    className="flex items-center space-x-2 px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors"
                                >
                                    <LogIn size={18} />
                                    <span className="font-medium">Login</span>
                                </Link>
                            )}

                            <button
                                onClick={toggleDarkMode}
                                className="p-2 rounded-lg bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors"
                                aria-label="Toggle dark mode"
                            >
                                {darkMode ? <Sun size={20} /> : <Moon size={20} />}
                            </button>

                            {/* Mobile menu button */}
                            <button
                                onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
                                className="md:hidden p-2 rounded-lg bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300"
                            >
                                {mobileMenuOpen ? <X size={20} /> : <Menu size={20} />}
                            </button>
                        </div>
                    </div>
                </div>

                {/* Mobile Navigation */}
                {mobileMenuOpen && (
                    <nav className="md:hidden border-t border-gray-200 dark:border-gray-700 py-2 px-4">
                        {navItems.map(({ path, icon: Icon, label }) => {
                            const isActive = location.pathname === path;
                            return (
                                <Link
                                    key={path}
                                    to={path}
                                    onClick={() => setMobileMenuOpen(false)}
                                    className={`flex items-center space-x-3 px-4 py-3 rounded-lg ${isActive
                                        ? 'bg-primary-100 dark:bg-primary-900/50 text-primary-700 dark:text-primary-300'
                                        : 'text-gray-600 dark:text-gray-300'
                                        }`}
                                >
                                    <Icon size={20} />
                                    <span className="font-medium">{label}</span>
                                </Link>
                            );
                        })}
                        {!isAuthenticated && (
                            <Link
                                to="/login"
                                onClick={() => setMobileMenuOpen(false)}
                                className="flex items-center space-x-3 px-4 py-3 text-primary-600 dark:text-primary-400"
                            >
                                <LogIn size={20} />
                                <span className="font-medium">Login</span>
                            </Link>
                        )}
                    </nav>
                )}
            </header>

            {/* Main Content */}
            <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
                <Outlet />
            </main>

            {/* Footer */}
            <footer className="bg-white dark:bg-gray-800 border-t border-gray-200 dark:border-gray-700 mt-auto">
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
                    <div className="flex flex-col md:flex-row justify-between items-center">
                        <p className="text-gray-500 dark:text-gray-400 text-sm">
                            © 2025 CrowdFundChain. Powered by Hyperledger Fabric.
                        </p>
                        <div className="flex items-center space-x-4 mt-4 md:mt-0">
                            <span className="text-xs text-gray-400 dark:text-gray-500">
                                Blockchain-secured crowdfunding
                            </span>
                        </div>
                    </div>
                </div>
            </footer>
        </div>
    );
}
