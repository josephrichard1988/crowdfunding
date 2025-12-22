import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import Layout from './components/Layout';
import { RoleGuard, ProtectedRoute } from './components/RouteGuards';

// Pages
import Dashboard from './pages/Dashboard';
import Login from './pages/Login';
import Signup from './pages/Signup';
import Wallet from './pages/Wallet';
import StartupDashboard from './pages/StartupDashboard';
import StartupCampaignDetails from './pages/StartupCampaignDetails';
import ValidatorDashboard from './pages/ValidatorDashboard';
import PlatformDashboard from './pages/PlatformDashboard';
import InvestorDashboard from './pages/InvestorDashboard';
import CampaignDetails from './pages/CampaignDetails';

function App() {
    return (
        <AuthProvider>
            <BrowserRouter>
                <Routes>
                    <Route path="/" element={<Layout />}>
                        {/* Public routes */}
                        <Route index element={<Dashboard />} />
                        <Route path="login" element={<Login />} />
                        <Route path="signup" element={<Signup />} />

                        {/* Protected wallet route */}
                        <Route path="wallet" element={
                            <ProtectedRoute>
                                <Wallet />
                            </ProtectedRoute>
                        } />

                        {/* Dashboard routes - accessible as preview for guests, full for auth users */}
                        <Route path="startup" element={<StartupDashboard />} />
                        <Route path="startup/campaign/:campaignId" element={
                            <RoleGuard allowedRoles={['STARTUP']}>
                                <StartupCampaignDetails />
                            </RoleGuard>
                        } />

                        <Route path="validator" element={<ValidatorDashboard />} />

                        <Route path="platform" element={<PlatformDashboard />} />

                        <Route path="investor" element={<InvestorDashboard />} />

                        <Route path="campaign/:campaignId" element={<CampaignDetails />} />
                    </Route>
                </Routes>
            </BrowserRouter>
        </AuthProvider>
    );
}

export default App;
