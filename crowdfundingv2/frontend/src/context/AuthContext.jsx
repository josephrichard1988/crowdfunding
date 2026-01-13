import { createContext, useContext, useState, useEffect } from 'react';

const AuthContext = createContext(null);

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:3001/api';

// ============================================================================
// STORAGE TOGGLE - Set to true for per-tab sessions (multi-user testing)
// Set to false for shared sessions across tabs (production behavior)
// TODO: Remove this toggle when no longer needed for testing
// ============================================================================
const USE_SESSION_STORAGE = true;  // Change to false for production

// Helper functions for storage (easy to remove later - just replace with localStorage)
const storage = {
    getItem: (key) => USE_SESSION_STORAGE ? sessionStorage.getItem(key) : localStorage.getItem(key),
    setItem: (key, value) => USE_SESSION_STORAGE ? sessionStorage.setItem(key, value) : localStorage.setItem(key, value),
    removeItem: (key) => USE_SESSION_STORAGE ? sessionStorage.removeItem(key) : localStorage.removeItem(key),
};
// ============================================================================

export function AuthProvider({ children }) {
    const [user, setUser] = useState(null);
    const [token, setToken] = useState(storage.getItem('token'));
    const [loading, setLoading] = useState(true);

    // Check auth status on mount
    useEffect(() => {
        if (token) {
            fetchCurrentUser();
        } else {
            setLoading(false);
        }
    }, []);

    const fetchCurrentUser = async () => {
        try {
            const response = await fetch(`${API_URL}/auth/me`, {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });

            if (response.ok) {
                const data = await response.json();
                setUser(data.user);
            } else {
                // Token invalid - clear it
                logout();
            }
        } catch (error) {
            console.error('Auth check failed:', error);
        } finally {
            setLoading(false);
        }
    };

    const login = async (email, password) => {
        const response = await fetch(`${API_URL}/auth/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password })
        });

        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error || 'Login failed');
        }

        storage.setItem('token', data.token);
        setToken(data.token);
        setUser(data.user);

        return data.user;
    };

    const signup = async (name, email, password, role) => {
        const response = await fetch(`${API_URL}/auth/signup`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name, email, password, role })
        });

        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error || 'Signup failed');
        }

        storage.setItem('token', data.token);
        setToken(data.token);
        setUser(data.user);

        return data.user;
    };

    const logout = () => {
        storage.removeItem('token');
        setToken(null);
        setUser(null);
    };

    const updateWallet = async (walletData) => {
        if (!user) return;

        try {
            const response = await fetch(`${API_URL}/auth/wallet`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify(walletData)
            });

            if (response.ok) {
                const data = await response.json();
                setUser(prev => ({ ...prev, wallet: data.wallet }));
            }
        } catch (error) {
            console.error('Wallet update failed:', error);
        }
    };

    const value = {
        user,
        token,
        loading,
        isAuthenticated: !!user,
        role: user?.role || null,
        login,
        signup,
        logout,
        updateWallet,
        refreshUser: fetchCurrentUser
    };

    return (
        <AuthContext.Provider value={value}>
            {children}
        </AuthContext.Provider>
    );
}

export function useAuth() {
    const context = useContext(AuthContext);
    if (!context) {
        throw new Error('useAuth must be used within an AuthProvider');
    }
    return context;
}

export default AuthContext;
