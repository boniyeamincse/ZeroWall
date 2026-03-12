import React, { useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, Link, useLocation, Navigate } from 'react-router-dom';
import Dashboard from './pages/Dashboard';
import FirewallRules from './pages/FirewallRules';
import VPNStatus from './pages/VPNStatus';
import Login from './pages/Login';
import { AuthProvider, useAuth, ProtectedRoute } from './context/AuthContext';
import { telemetry } from './services/WebSocketService';
import { LogOut, Shield, Settings, Wifi, Lock } from 'lucide-react';

const SidebarItem = ({ to, label, icon: Icon }: { to: string, label: string, icon: React.ElementType }) => {
    const location = useLocation();
    const isActive = location.pathname === to;
    return (
        <li>
            <Link to={to} className={`flex items-center gap-3 p-3 rounded-xl transition-all duration-200 ${isActive ? 'bg-cyan-600 text-white shadow-lg shadow-cyan-900/50' : 'text-slate-400 hover:text-cyan-300 hover:bg-slate-800'}`}>
                <Icon className="w-5 h-5" />
                <span className="font-medium">{label}</span>
            </Link>
        </li>
    );
};

const Sidebar = () => {
    const { user, logout } = useAuth();

    return (
        <nav className="w-72 p-6 bg-slate-900/50 backdrop-blur-xl border-r border-slate-800 sticky top-0 h-screen flex flex-col">
            <div className="flex items-center gap-3 mb-10">
                <div className="w-10 h-10 bg-gradient-to-br from-cyan-500 to-blue-600 rounded-xl flex items-center justify-center shadow-lg shadow-cyan-500/20">
                    <Shield className="w-5 h-5 text-white" />
                </div>
                <div>
                    <h1 className="text-xl font-black tracking-tight text-white">ZEROWALL</h1>
                    <p className="text-xs text-slate-500">Firewall Gateway</p>
                </div>
            </div>

            <ul className="space-y-2 flex-1">
                <SidebarItem to="/" label="Dashboard" icon={Settings} />
                <SidebarItem to="/firewall" label="Firewall Rules" icon={Lock} />
                <SidebarItem to="/vpn" label="VPN Status" icon={Wifi} />
                <SidebarItem to="/certs" label="Certificates" icon={Shield} />
            </ul>

            <div className="border-t border-slate-800 pt-4 mt-4">
                <div className="flex items-center gap-3 p-3 mb-3">
                    <div className="w-8 h-8 bg-slate-700 rounded-full flex items-center justify-center text-sm font-bold text-slate-300">
                        {user?.username?.charAt(0).toUpperCase() || 'A'}
                    </div>
                    <div className="flex-1 min-w-0">
                        <p className="text-sm font-medium text-white truncate">{user?.username || 'Admin'}</p>
                        <p className="text-xs text-slate-500 capitalize">{user?.role || 'Administrator'}</p>
                    </div>
                </div>
                <button
                    onClick={logout}
                    className="w-full flex items-center justify-center gap-2 p-3 rounded-xl bg-red-500/10 text-red-400 hover:bg-red-500/20 transition-colors text-sm font-medium"
                >
                    <LogOut className="w-4 h-4" />
                    Sign Out
                </button>
            </div>
        </nav>
    );
};

const AppLayout = ({ children }: { children: React.ReactNode }) => {
    return (
        <div className="flex bg-slate-950 text-slate-100 min-h-screen font-sans selection:bg-cyan-500/30">
            <Sidebar />
            <main className="flex-1 overflow-y-auto">
                <div className="max-w-7xl mx-auto p-10">
                    {children}
                </div>
            </main>
        </div>
    );
};

const AppContent: React.FC = () => {
    const { isAuthenticated, isLoading } = useAuth();

    useEffect(() => {
        if (isAuthenticated) {
            telemetry.connect('ws://' + window.location.host + '/ws/stats');
        }
    }, [isAuthenticated]);

    if (isLoading) {
        return (
            <div className="min-h-screen bg-slate-950 flex items-center justify-center">
                <div className="text-center">
                    <div className="w-12 h-12 border-4 border-cyan-500 border-t-transparent rounded-full animate-spin mx-auto mb-4"></div>
                    <p className="text-slate-400">Loading...</p>
                </div>
            </div>
        );
    }

    return (
        <Routes>
            <Route path="/login" element={
                isAuthenticated ? <Navigate to="/" replace /> : <Login onLogin={() => {}} />
            } />
            <Route path="/" element={
                <ProtectedRoute>
                    <AppLayout>
                        <Dashboard />
                    </AppLayout>
                </ProtectedRoute>
            } />
            <Route path="/firewall" element={
                <ProtectedRoute>
                    <AppLayout>
                        <FirewallRules />
                    </AppLayout>
                </ProtectedRoute>
            } />
            <Route path="/vpn" element={
                <ProtectedRoute>
                    <AppLayout>
                        <VPNStatus />
                    </AppLayout>
                </ProtectedRoute>
            } />
            <Route path="/certs" element={
                <ProtectedRoute>
                    <AppLayout>
                        <div className="text-white">Certificates Page</div>
                    </AppLayout>
                </ProtectedRoute>
            } />
            <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
    );
};

const App: React.FC = () => {
    return (
        <Router>
            <AuthProvider>
                <AppContent />
            </AuthProvider>
        </Router>
    );
};

export default App;
