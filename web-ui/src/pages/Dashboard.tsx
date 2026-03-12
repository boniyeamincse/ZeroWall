import React from 'react';

const Dashboard: React.FC = () => {
    return (
        <div>
            <h2 className="text-3xl font-semibold mb-6">System Dashboard</h2>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                <div className="bg-slate-800 p-6 rounded-lg border border-slate-700">
                    <h3 className="text-gray-400 mb-2">CPU Usage</h3>
                    <p className="text-4xl font-bold">12%</p>
                </div>
                <div className="bg-slate-800 p-6 rounded-lg border border-slate-700">
                    <h3 className="text-gray-400 mb-2">Throughput</h3>
                    <p className="text-4xl font-bold">1.2 Gbps</p>
                </div>
                <div className="bg-slate-800 p-6 rounded-lg border border-slate-700">
                    <h3 className="text-gray-400 mb-2">Active VPNs</h3>
                    <p className="text-4xl font-bold">8</p>
                </div>
            </div>
        </div>
    );
};

export default Dashboard;
