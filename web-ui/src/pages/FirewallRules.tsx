import React from 'react';

const FirewallRules: React.FC = () => {
    return (
        <div>
            <div className="flex justify-between items-center mb-6">
                <h2 className="text-3xl font-semibold">Firewall Rules</h2>
                <button className="bg-cyan-600 hover:bg-cyan-500 px-4 py-2 rounded">Add Rule</button>
            </div>

            <table className="w-full text-left bg-slate-800 rounded-lg overflow-hidden">
                <thead>
                    <tr className="bg-slate-700">
                        <th className="p-4">Action</th>
                        <th className="p-4">Source</th>
                        <th className="p-4">Destination</th>
                        <th className="p-4">Port</th>
                        <th className="p-4">Description</th>
                    </tr>
                </thead>
                <tbody>
                    <tr className="border-t border-slate-700">
                        <td className="p-4">PASS</td>
                        <td className="p-4">LAN Net</td>
                        <td className="p-4">*</td>
                        <td className="p-4">80, 443</td>
                        <td className="p-4">Allow Internet Access</td>
                    </tr>
                </tbody>
            </table>
        </div>
    );
};

export default FirewallRules;
