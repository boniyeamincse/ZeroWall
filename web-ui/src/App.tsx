import React from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import Dashboard from './pages/Dashboard';
import FirewallRules from './pages/FirewallRules';
import VPNStatus from './pages/VPNStatus';

const App: React.FC = () => {
    return (
        <Router>
            <div className="flex bg-slate-900 text-white min-h-screen">
                {/* Sidebar Nav */}
                <nav className="w-64 p-6 border-r border-slate-700">
                    <h1 className="text-2xl font-bold text-cyan-400 mb-8">ZeroWall</h1>
                    <ul className="space-y-4">
                        <li><Link to="/" className="hover:text-cyan-300">Dashboard</Link></li>
                        <li><Link to="/firewall" className="hover:text-cyan-300">Firewall Rules</Link></li>
                        <li><Link to="/vpn" className="hover:text-cyan-300">VPN Status</Link></li>
                        <li><Link to="/certs" className="hover:text-cyan-300">Certificates</Link></li>
                    </ul>
                </nav>

                {/* Main Content */}
                <main className="flex-1 p-8">
                    <Routes>
                        <Route path="/" element={<Dashboard />} />
                        <Route path="/firewall" element={<FirewallRules />} />
                        <Route path="/vpn" element={<VPNStatus />} />
                    </Routes>
                </main>
            </div>
        </Router>
    );
};

export default App;
