import { Navigate, useLocation } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

// Protects routes that require authentication
export function ProtectedRoute({ children }) {
    const { isAuthenticated, loading } = useAuth();
    const location = useLocation();

    if (loading) {
        return (
            <div className="min-h-[60vh] flex items-center justify-center">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-500"></div>
            </div>
        );
    }

    if (!isAuthenticated) {
        // Redirect to login, preserving the intended destination
        return <Navigate to="/login" state={{ from: location.pathname }} replace />;
    }

    return children;
}

// Guards routes based on user role
export function RoleGuard({ allowedRoles, children }) {
    const { isAuthenticated, role, loading } = useAuth();
    const location = useLocation();

    if (loading) {
        return (
            <div className="min-h-[60vh] flex items-center justify-center">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-500"></div>
            </div>
        );
    }

    if (!isAuthenticated) {
        return <Navigate to="/login" state={{ from: location.pathname }} replace />;
    }

    if (!allowedRoles.includes(role)) {
        // Redirect to user's own dashboard if trying to access wrong role
        return <Navigate to={`/${role.toLowerCase()}`} replace />;
    }

    return children;
}

// Shows dashboard preview for guests, full dashboard for authenticated users
export function DashboardWrapper({ role, PreviewComponent, FullComponent }) {
    const { isAuthenticated, role: userRole, loading } = useAuth();

    if (loading) {
        return (
            <div className="min-h-[60vh] flex items-center justify-center">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-500"></div>
            </div>
        );
    }

    // If not authenticated, show preview
    if (!isAuthenticated) {
        return <PreviewComponent role={role} />;
    }

    // If authenticated but wrong role, redirect to own dashboard
    if (userRole !== role) {
        return <Navigate to={`/${userRole.toLowerCase()}`} replace />;
    }

    // Authenticated and correct role - show full dashboard
    return <FullComponent />;
}

export default { ProtectedRoute, RoleGuard, DashboardWrapper };
