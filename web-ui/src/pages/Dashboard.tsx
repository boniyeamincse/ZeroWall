import React, { useState, useEffect } from 'react';
import { telemetry } from '../services/WebSocketService';

const StatCard = ({ title, value, unit, color }: { title: string, value: string | number, unit: string, color: string }) => (
    <div className="bg-slate-900/50 backdrop-blur-md border border-slate-800 p-6 rounded-2xl hover:border-slate-700 transition-all group">
        <h3 className="text-slate-400 text-sm font-semibold uppercase tracking-wider mb-2">{title}</h3>
        <div className="flex items-baseline gap-2">
            <span className={`text-4xl font-bold bg-gradient-to-r ${color} bg-clip-text text-transparent`}>{value}</span>
            <span className="text-slate-500 font-medium">{unit}</span>
        </div>
        <div className="mt-4 h-1 w-full bg-slate-800 rounded-full overflow-hidden">
            <div className={`h-full bg-gradient-to-r ${color} transition-all duration-500`} style={{ width: `${Math.min(Number(value), 100)}%` }}></div>
        </div>
    </div>
);

const Dashboard: React.FC = () => {
    const [stats, setStats] = useState({
        cpu: 12,
        ram: 45,
        throughputIn: 1.2,
        throughputOut: 0.8
    });

    useEffect(() => {
        telemetry.addListener((data) => {
            if (data.type === 'sys_stats') {
                setStats(prev => ({ ...prev, ...data.payload }));
            }
        });
    }, []);

    return (
        <div className="space-y-10 animate-in fade-in duration-700">
            <header className="flex justify-between items-end">
                <div>
                    <h2 className="text-4xl font-black text-white mb-2">System Dashboard</h2>
                    <p className="text-slate-400">Real-time telemetry and network health overview.</p>
                </div>
                <div className="bg-emerald-500/10 text-emerald-400 px-4 py-2 rounded-full text-sm font-bold border border-emerald-500/20 flex items-center gap-2">
                    <span className="w-2 h-2 bg-emerald-400 rounded-full animate-pulse"></span>
                    SYSTEM ONLINE
                </div>
            </header>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                <StatCard title="CPU Usage" value={stats.cpu} unit="%" color="from-cyan-400 to-blue-500" />
                <StatCard title="RAM Usage" value={stats.ram} unit="%" color="from-purple-400 to-pink-500" />
                <StatCard title="WAN Inbound" value={stats.throughputIn} unit="Mbps" color="from-emerald-400 to-teal-500" />
                <StatCard title="WAN Outbound" value={stats.throughputOut} unit="Mbps" color="from-orange-400 to-red-500" />
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                <div className="lg:col-span-2 bg-slate-900/50 backdrop-blur-md border border-slate-800 rounded-2xl p-8">
                    <h3 className="text-xl font-bold mb-6 flex items-center gap-2 text-white">
                        <span className="w-5 h-5 bg-cyan-500/20 rounded flex items-center justify-center text-cyan-400 text-xs">📈</span>
                        Live Traffic History
                    </h3>
                    <div className="h-64 flex items-end gap-1 px-4">
                        {[40, 60, 45, 80, 70, 90, 65, 55, 75, 85, 95, 60, 40, 50, 65].map((h, i) => (
                            <div key={i} className="flex-1 bg-gradient-to-t from-cyan-600/50 to-cyan-400 rounded-t-sm transition-all hover:scale-x-125" style={{ height: `${h}%` }}></div>
                        ))}
                    </div>
                </div>

                <div className="bg-slate-900/50 backdrop-blur-md border border-slate-800 rounded-2xl p-8">
                    <h3 className="text-xl font-bold mb-6 text-white text-white">Recent Security Alerts</h3>
                    <div className="space-y-4">
                        {[
                            { time: '12:45:01', msg: 'SQLi Attempt Blocked', src: '192.168.1.50' },
                            { time: '12:43:55', msg: 'GeoIP Block: Russia', src: '45.16.2.1' },
                            { time: '12:42:02', msg: 'Port Scan Detected', src: '10.0.0.12' }
                        ].map((alert, i) => (
                            <div key={i} className="text-xs border-l-2 border-red-500/50 pl-4 py-2 group cursor-pointer hover:bg-red-500/5 rounded-r transition-all">
                                <span className="text-slate-500 font-mono block">{alert.time}</span>
                                <span className="text-slate-200 font-bold group-hover:text-red-400 transition-colors uppercase">{alert.msg}</span>
                                <span className="text-slate-500 block">{alert.src}</span>
                            </div>
                        ))}
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Dashboard;
