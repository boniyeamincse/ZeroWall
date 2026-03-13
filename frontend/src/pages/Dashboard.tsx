import { useEffect, useState } from 'react';
import {
  Cpu,
  HardDrive,
  MemoryStick,
  Network,
  Shield,
  Activity,
  Lock,
  Clock,
  TrendingUp,
  TrendingDown,
  AlertTriangle,
} from 'lucide-react';
import { LineChart, Line, XAxis, YAxis, Tooltip, ResponsiveContainer } from 'recharts';
import { firewallService } from '../services/api';

const mockTrafficData = [
  { time: '00:00', in: 45, out: 12 },
  { time: '04:00', in: 28, out: 8 },
  { time: '08:00', in: 120, out: 45 },
  { time: '12:00', in: 180, out: 65 },
  { time: '16:00', in: 150, out: 55 },
  { time: '20:00', in: 95, out: 30 },
  { time: '24:00', in: 60, out: 20 },
];

export function Dashboard() {
  const [stats, setStats] = useState({ states: 0, rules: 0, blocked_hour: 0 });
  const [systemInfo, setSystemInfo] = useState({
    cpu: 12,
    memory: 38,
    disk: 45,
    uptime: '14 days',
  });

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const data = await firewallService.getStats();
        setStats(data);
      } catch (error) {
        console.error('Failed to fetch stats:', error);
      }
    };
    fetchStats();
    const interval = setInterval(fetchStats, 30000);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-white">Dashboard</h1>
        <div className="flex items-center gap-2 text-sm text-gray-400">
          <Clock className="w-4 h-4" />
          <span>Last updated: {new Date().toLocaleTimeString()}</span>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <div className="bg-zw-dark border border-zw-light rounded-lg p-6">
          <div className="flex items-center justify-between mb-4">
            <div className="w-12 h-12 bg-zw-accent/10 rounded-lg flex items-center justify-center">
              <Cpu className="w-6 h-6 text-zw-accent" />
            </div>
            <span className="text-xs text-gray-400">Live</span>
          </div>
          <p className="text-2xl font-bold text-white">{systemInfo.cpu}%</p>
          <p className="text-sm text-gray-400">CPU Usage</p>
        </div>

        <div className="bg-zw-dark border border-zw-light rounded-lg p-6">
          <div className="flex items-center justify-between mb-4">
            <div className="w-12 h-12 bg-purple-500/10 rounded-lg flex items-center justify-center">
              <MemoryStick className="w-6 h-6 text-purple-500" />
            </div>
          </div>
          <p className="text-2xl font-bold text-white">{systemInfo.memory}%</p>
          <p className="text-sm text-gray-400">Memory Usage</p>
        </div>

        <div className="bg-zw-dark border border-zw-light rounded-lg p-6">
          <div className="flex items-center justify-between mb-4">
            <div className="w-12 h-12 bg-green-500/10 rounded-lg flex items-center justify-center">
              <HardDrive className="w-6 h-6 text-green-500" />
            </div>
          </div>
          <p className="text-2xl font-bold text-white">{systemInfo.disk}%</p>
          <p className="text-sm text-gray-400">Disk Usage</p>
        </div>

        <div className="bg-zw-dark border border-zw-light rounded-lg p-6">
          <div className="flex items-center justify-between mb-4">
            <div className="w-12 h-12 bg-yellow-500/10 rounded-lg flex items-center justify-center">
              <Clock className="w-6 h-6 text-yellow-500" />
            </div>
          </div>
          <p className="text-2xl font-bold text-white">{systemInfo.uptime}</p>
          <p className="text-sm text-gray-400">Uptime</p>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-zw-dark border border-zw-light rounded-lg p-6">
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-lg font-semibold text-white">Traffic Overview</h2>
            <div className="flex items-center gap-4 text-sm">
              <span className="flex items-center gap-2">
                <TrendingUp className="w-4 h-4 text-zw-accent" />
                <span className="text-gray-400">Inbound</span>
              </span>
              <span className="flex items-center gap-2">
                <TrendingDown className="w-4 h-4 text-green-500" />
                <span className="text-gray-400">Outbound</span>
              </span>
            </div>
          </div>
          <div className="h-64">
            <ResponsiveContainer width="100%" height="100%">
              <LineChart data={mockTrafficData}>
                <XAxis dataKey="time" stroke="#64748b" fontSize={12} />
                <YAxis stroke="#64748b" fontSize={12} />
                <Tooltip
                  contentStyle={{
                    backgroundColor: '#0f2942',
                    border: '1px solid #1e3a5f',
                    borderRadius: '8px',
                  }}
                  labelStyle={{ color: '#fff' }}
                />
                <Line
                  type="monotone"
                  dataKey="in"
                  stroke="#00d4ff"
                  strokeWidth={2}
                  dot={false}
                />
                <Line
                  type="monotone"
                  dataKey="out"
                  stroke="#00c853"
                  strokeWidth={2}
                  dot={false}
                />
              </LineChart>
            </ResponsiveContainer>
          </div>
        </div>

        <div className="bg-zw-dark border border-zw-light rounded-lg p-6">
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-lg font-semibold text-white">Firewall Status</h2>
            <Shield className="w-5 h-5 text-zw-success" />
          </div>
          <div className="space-y-4">
            <div className="flex items-center justify-between p-4 bg-zw-darker rounded-lg">
              <div className="flex items-center gap-3">
                <Activity className="w-5 h-5 text-zw-accent" />
                <span className="text-gray-300">Active States</span>
              </div>
              <span className="text-xl font-bold text-white">{stats.states.toLocaleString()}</span>
            </div>
            <div className="flex items-center justify-between p-4 bg-zw-darker rounded-lg">
              <div className="flex items-center gap-3">
                <Shield className="w-5 h-5 text-purple-500" />
                <span className="text-gray-300">Firewall Rules</span>
              </div>
              <span className="text-xl font-bold text-white">{stats.rules}</span>
            </div>
            <div className="flex items-center justify-between p-4 bg-zw-darker rounded-lg">
              <div className="flex items-center gap-3">
                <AlertTriangle className="w-5 h-5 text-zw-warning" />
                <span className="text-gray-300">Blocked (last hour)</span>
              </div>
              <span className="text-xl font-bold text-white">{stats.blocked_hour}</span>
            </div>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="bg-zw-dark border border-zw-light rounded-lg p-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-semibold text-white">VPN Connections</h2>
            <Lock className="w-5 h-5 text-zw-accent" />
          </div>
          <div className="text-center py-8">
            <p className="text-3xl font-bold text-white">0</p>
            <p className="text-sm text-gray-400">Active connections</p>
          </div>
          <div className="mt-4 text-sm text-gray-400">
            <p>WireGuard: Inactive</p>
            <p>OpenVPN: Inactive</p>
          </div>
        </div>

        <div className="bg-zw-dark border border-zw-light rounded-lg p-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-semibold text-white">Network Interfaces</h2>
            <Network className="w-5 h-5 text-zw-accent" />
          </div>
          <div className="space-y-3">
            <div className="flex items-center justify-between p-3 bg-zw-darker rounded-lg">
              <span className="text-gray-300">WAN (em0)</span>
              <span className="text-green-500 text-sm">Online</span>
            </div>
            <div className="flex items-center justify-between p-3 bg-zw-darker rounded-lg">
              <span className="text-gray-300">LAN (em1)</span>
              <span className="text-green-500 text-sm">Online</span>
            </div>
          </div>
        </div>

        <div className="bg-zw-dark border border-zw-light rounded-lg p-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-semibold text-white">IDS/IPS</h2>
            <Activity className="w-5 h-5 text-zw-accent" />
          </div>
          <div className="text-center py-8">
            <p className="text-3xl font-bold text-white">0</p>
            <p className="text-sm text-gray-400">Alerts (24h)</p>
          </div>
          <div className="mt-4">
            <div className="w-full bg-zw-darker rounded-full h-2">
              <div className="bg-zw-success h-2 rounded-full" style={{ width: '100%' }}></div>
            </div>
            <p className="text-sm text-gray-400 mt-2">IDS Engine: Enabled</p>
          </div>
        </div>
      </div>
    </div>
  );
}
