import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Layout from './components/Layout';
import Dashboard from './pages/Dashboard';
import StartupDashboard from './pages/StartupDashboard';
import StartupCampaignDetails from './pages/StartupCampaignDetails';
import ValidatorDashboard from './pages/ValidatorDashboard';
import PlatformDashboard from './pages/PlatformDashboard';
import InvestorDashboard from './pages/InvestorDashboard';
import CampaignDetails from './pages/CampaignDetails';

function App() {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<Layout />}>
                    <Route index element={<Dashboard />} />
                    <Route path="startup" element={<StartupDashboard />} />
                    <Route path="startup/campaign/:campaignId" element={<StartupCampaignDetails />} />
                    <Route path="validator" element={<ValidatorDashboard />} />
                    <Route path="platform" element={<PlatformDashboard />} />
                    <Route path="investor" element={<InvestorDashboard />} />
                    <Route path="campaign/:campaignId" element={<CampaignDetails />} />
                </Route>
            </Routes>
        </BrowserRouter>
    );
}

export default App;
