import React, { useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, Link, useLocation } from 'react-router-dom';
import Dashboard from './pages/Dashboard';
import FirewallRules from './pages/FirewallRules';
import VPNStatus from './pages/VPNStatus';
import { telemetry } from './services/WebSocketService';

const SidebarItem = ({ to, label }: { to: string, label: string }) => {
    const location = useLocation();
    const isActive = location.pathname === to;
    return (
        <li>
            <Link to={to} className={`block p-2 rounded transition-all duration-200 ${isActive ? 'bg-cyan-600 text-white shadow-lg shadow-cyan-900/50' : 'text-slate-400 hover:text-cyan-300 hover:bg-slate-800'}`}>
                {label}
            </Link>
        </li>
    );
};

const App: React.FC = () => {
    useEffect(() => {
        telemetry.connect('ws://' + window.location.host + '/ws/stats');
    }, []);

    return (
        <Router>
            <div className="flex bg-slate-950 text-slate-100 min-h-screen font-sans selection:bg-cyan-500/30">
                {/* Sidebar Nav */}
                <nav className="w-64 p-6 bg-slate-900/50 backdrop-blur-xl border-r border-slate-800 sticky top-0 h-screen">
                    <div className="flex items-center gap-3 mb-10">
                        <div className="w-8 h-8 bg-cyan-500 rounded-lg shadow-lg shadow-cyan-500/20 animate-pulse"></div>
                        <h1 className="text-2xl font-black tracking-tighter text-white">ZEROWALL</h1>
                    </div>
                    <ul className="space-y-2 text-sm font-medium">
                        <SidebarItem to="/" label="Dashboard" />
                        <SidebarItem to="/firewall" label="Firewall Rules" />
                        <SidebarItem to="/vpn" label="VPN Status" />
                        <SidebarItem to="/certs" label="Certificates" />
                        <SidebarItem to="/settings" label="System Settings" />
                    </ul>
                </nav>

                {/* Main Content */}
                <main className="flex-1 overflow-y-auto">
                    <div className="max-w-7xl mx-auto p-10">
                        <Routes>
                            <Route path="/" element={<Dashboard />} />
                            <Route path="/firewall" element={<FirewallRules />} />
                            <Route path="/vpn" element={<VPNStatus />} />
                            <Route path="/certs" element={<div>Certificates Page (Coming Soon)</div>} />
                        </Routes>
                    </div>
                </main>
            </div>
        </Router>
    );
};

export default App;
