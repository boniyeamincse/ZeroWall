import React from 'react';

const VPNStatus: React.FC = () => {
    return (
        <div>
            <h2 className="text-3xl font-semibold mb-6">VPN Status</h2>
            <div className="space-y-6">
                <div className="bg-slate-800 p-6 rounded-lg border border-slate-700">
                    <h3 className="text-cyan-400 font-bold mb-4">WireGuard (wg0)</h3>
                    <div className="flex justify-between">
                        <span>Peers Connected: 5</span>
                        <span className="text-green-400">STATUS: ACTIVE</span>
                    </div>
                </div>
                <div className="bg-slate-800 p-6 rounded-lg border border-slate-700">
                    <h3 className="text-cyan-400 font-bold mb-4">OpenVPN (ovpns1)</h3>
                    <div className="flex justify-between">
                        <span>Clients Connected: 12</span>
                        <span className="text-green-400">STATUS: ACTIVE</span>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default VPNStatus;
