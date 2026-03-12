import React, { useState } from 'react';

interface Rule {
    id: string;
    action: 'pass' | 'block' | 'reject';
    interface: string;
    protocol: string;
    source: string;
    destination: string;
    description: string;
}

const RulesManager: React.FC = () => {
    const [rules] = useState<Rule[]>([
        { id: '1', action: 'pass', interface: 'WAN', protocol: 'TCP/UDP', source: 'Any', destination: '8.8.8.8', description: 'Allow DNS' },
        { id: '2', action: 'block', interface: 'LAN', protocol: 'Any', source: '192.168.2.0/24', destination: 'Any', description: 'Block Guest Network' }
    ]);

    return (
        <div className="space-y-8 animate-in slide-in-from-bottom-4 duration-500">
            <header className="flex justify-between items-center">
                <div>
                    <h2 className="text-3xl font-black text-white">Firewall Rules</h2>
                    <p className="text-slate-400">Manage incoming and outgoing traffic filters.</p>
                </div>
                <button className="bg-cyan-500 hover:bg-cyan-400 text-slate-950 font-bold px-6 py-2 rounded-xl transition-all shadow-lg shadow-cyan-500/20 active:scale-95">
                    + Create New Rule
                </button>
            </header>

            <div className="bg-slate-900/50 backdrop-blur-md border border-slate-800 rounded-2xl overflow-hidden">
                <table className="w-full text-left">
                    <thead className="bg-slate-800/50 text-slate-400 text-xs font-bold uppercase tracking-widest border-b border-slate-800">
                        <tr>
                            <th className="px-6 py-4">Action</th>
                            <th className="px-6 py-4">Interface</th>
                            <th className="px-6 py-4">Protocol</th>
                            <th className="px-6 py-4">Source</th>
                            <th className="px-6 py-4">Destination</th>
                            <th className="px-6 py-4 text-right">Settings</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-slate-800">
                        {rules.map(rule => (
                            <tr key={rule.id} className="hover:bg-cyan-500/5 transition-colors group cursor-pointer">
                                <td className="px-6 py-4">
                                    <span className={`text-[10px] font-black uppercase px-2 py-1 rounded inline-block ${rule.action === 'pass' ? 'bg-emerald-500/20 text-emerald-400' : 'bg-red-500/20 text-red-400'
                                        }`}>
                                        {rule.action}
                                    </span>
                                </td>
                                <td className="px-6 py-4 text-sm text-slate-300 font-medium">{rule.interface}</td>
                                <td className="px-6 py-4 text-sm text-slate-400 font-mono">{rule.protocol}</td>
                                <td className="px-6 py-4 text-sm text-slate-300 font-mono">{rule.source}</td>
                                <td className="px-6 py-4 text-sm text-slate-300 font-mono">{rule.destination}</td>
                                <td className="px-6 py-4 text-right">
                                    <button className="text-slate-600 hover:text-white transition-colors p-2 rounded-lg hover:bg-slate-800">⚙️</button>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>

            <div className="bg-amber-500/10 border border-amber-500/20 p-4 rounded-xl flex items-center gap-4 text-amber-200 text-sm">
                <span className="text-2xl">⚠️</span>
                <p>Rules are processed top-to-bottom. First matching rule wins. Drag and drop rows to change priority.</p>
            </div>
        </div>
    );
};

export default RulesManager;
