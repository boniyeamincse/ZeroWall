import React, { useState, useEffect } from 'react';
import { Plus, Trash2, Edit, Power, PowerOff, RefreshCw, Search, Filter } from 'lucide-react';

interface FirewallRule {
  id: string;
  uuid: string;
  sequence: number;
  enabled: boolean;
  action: string;
  interface: string;
  direction: string;
  protocol: string;
  source?: {
    network: string;
    port?: string;
  };
  destination?: {
    network: string;
    port?: string;
  };
  log: boolean;
  description: string;
}

interface APIResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
}

const FirewallRules: React.FC = () => {
  const [rules, setRules] = useState<FirewallRule[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [filterAction, setFilterAction] = useState('all');
  const [showModal, setShowModal] = useState(false);
  const [editingRule, setEditingRule] = useState<FirewallRule | null>(null);

  const fetchRules = async () => {
    setLoading(true);
    setError(null);
    try {
      const token = localStorage.getItem('auth_token');
      const response = await fetch('/api/v1/firewall/rules', {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
      });
      const data: APIResponse<FirewallRule[]> = await response.json();
      if (data.success && data.data) {
        setRules(data.data);
      } else {
        setError(data.error || 'Failed to fetch rules');
      }
    } catch {
      setError('Failed to connect to API');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchRules();
  }, []);

  const handleToggleRule = async (id: string) => {
    try {
      const token = localStorage.getItem('auth_token');
      await fetch(`/api/v1/firewall/rule/toggle/${id}`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
      });
      fetchRules();
    } catch {
      console.error('Failed to toggle rule');
    }
  };

  const handleDeleteRule = async (id: string) => {
    if (!confirm('Are you sure you want to delete this rule?')) return;
    try {
      const token = localStorage.getItem('auth_token');
      await fetch(`/api/v1/firewall/rule/${id}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });
      fetchRules();
    } catch {
      console.error('Failed to delete rule');
    }
  };

  const filteredRules = rules.filter(rule => {
    const matchesSearch = rule.description.toLowerCase().includes(searchTerm.toLowerCase()) ||
      rule.source?.network?.includes(searchTerm) ||
      rule.destination?.network?.includes(searchTerm);
    const matchesFilter = filterAction === 'all' || rule.action === filterAction;
    return matchesSearch && matchesFilter;
  });

  const getActionColor = (action: string) => {
    switch (action) {
      case 'pass': return 'bg-emerald-500/20 text-emerald-400';
      case 'block': return 'bg-red-500/20 text-red-400';
      case 'reject': return 'bg-orange-500/20 text-orange-400';
      default: return 'bg-slate-500/20 text-slate-400';
    }
  };

  return (
    <div className="space-y-6 animate-in fade-in duration-500">
      <header className="flex justify-between items-center">
        <div>
          <h2 className="text-3xl font-black text-white">Firewall Rules</h2>
          <p className="text-slate-400">Manage incoming and outgoing traffic filters.</p>
        </div>
        <div className="flex gap-3">
          <button
            onClick={fetchRules}
            className="flex items-center gap-2 px-4 py-2 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded-xl transition-all"
          >
            <RefreshCw className={`w-4 h-4 ${loading ? 'animate-spin' : ''}`} />
            Refresh
          </button>
          <button
            onClick={() => { setEditingRule(null); setShowModal(true); }}
            className="flex items-center gap-2 px-4 py-2 bg-cyan-600 hover:bg-cyan-500 text-white font-bold rounded-xl transition-all shadow-lg shadow-cyan-500/20"
          >
            <Plus className="w-4 h-4" />
            Add Rule
          </button>
        </div>
      </header>

      <div className="flex gap-4 items-center">
        <div className="relative flex-1 max-w-md">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-500" />
          <input
            type="text"
            placeholder="Search rules..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="w-full bg-slate-900 border border-slate-700 rounded-xl py-2.5 pl-10 pr-4 text-white placeholder-slate-500 focus:outline-none focus:border-cyan-500 transition-all"
          />
        </div>
        <div className="flex items-center gap-2">
          <Filter className="w-5 h-5 text-slate-500" />
          <select
            value={filterAction}
            onChange={(e) => setFilterAction(e.target.value)}
            className="bg-slate-900 border border-slate-700 rounded-xl py-2.5 px-4 text-white focus:outline-none focus:border-cyan-500 transition-all"
          >
            <option value="all">All Actions</option>
            <option value="pass">Pass</option>
            <option value="block">Block</option>
            <option value="reject">Reject</option>
          </select>
        </div>
      </div>

      {error && (
        <div className="bg-red-500/10 border border-red-500/20 p-4 rounded-xl text-red-400 flex items-center gap-3">
          <span>{error}</span>
          <button onClick={fetchRules} className="text-red-300 underline">Retry</button>
        </div>
      )}

      <div className="bg-slate-900/50 backdrop-blur-md border border-slate-800 rounded-2xl overflow-hidden">
        <table className="w-full text-left">
          <thead className="bg-slate-800/50 text-slate-400 text-xs font-bold uppercase tracking-widest border-b border-slate-800">
            <tr>
              <th className="px-6 py-4 w-16">#</th>
              <th className="px-6 py-4">Action</th>
              <th className="px-6 py-4">Interface</th>
              <th className="px-6 py-4">Protocol</th>
              <th className="px-6 py-4">Source</th>
              <th className="px-6 py-4">Destination</th>
              <th className="px-6 py-4">Description</th>
              <th className="px-6 py-4 text-right">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-800">
            {loading ? (
              <tr>
                <td colSpan={8} className="px-6 py-12 text-center">
                  <div className="flex justify-center">
                    <RefreshCw className="w-6 h-6 text-cyan-500 animate-spin" />
                  </div>
                </td>
              </tr>
            ) : filteredRules.length === 0 ? (
              <tr>
                <td colSpan={8} className="px-6 py-12 text-center text-slate-500">
                  No firewall rules found. Click "Add Rule" to create one.
                </td>
              </tr>
            ) : (
              filteredRules.map((rule) => (
                <tr key={rule.id || rule.uuid} className="hover:bg-cyan-500/5 transition-colors">
                  <td className="px-6 py-4 text-slate-500 font-mono text-sm">{rule.sequence}</td>
                  <td className="px-6 py-4">
                    <span className={`text-xs font-black uppercase px-2 py-1 rounded ${getActionColor(rule.action)}`}>
                      {rule.action}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-slate-300 font-medium">{rule.interface}</td>
                  <td className="px-6 py-4 text-slate-400 font-mono uppercase text-sm">{rule.protocol}</td>
                  <td className="px-6 py-4 text-slate-300 font-mono text-sm">
                    {rule.source?.network || 'any'}
                    {rule.source?.port && `:${rule.source.port}`}
                  </td>
                  <td className="px-6 py-4 text-slate-300 font-mono text-sm">
                    {rule.destination?.network || 'any'}
                    {rule.destination?.port && `:${rule.destination.port}`}
                  </td>
                  <td className="px-6 py-4 text-slate-400">{rule.description}</td>
                  <td className="px-6 py-4 text-right">
                    <div className="flex items-center justify-end gap-2">
                      <button
                        onClick={() => handleToggleRule(rule.id || rule.uuid)}
                        className={`p-2 rounded-lg transition-all ${rule.enabled ? 'text-emerald-400 hover:bg-emerald-500/10' : 'text-slate-500 hover:bg-slate-800'}`}
                        title={rule.enabled ? 'Disable' : 'Enable'}
                      >
                        {rule.enabled ? <Power className="w-4 h-4" /> : <PowerOff className="w-4 h-4" />}
                      </button>
                      <button
                        onClick={() => { setEditingRule(rule); setShowModal(true); }}
                        className="p-2 text-slate-500 hover:text-cyan-400 hover:bg-cyan-500/10 rounded-lg transition-all"
                      >
                        <Edit className="w-4 h-4" />
                      </button>
                      <button
                        onClick={() => handleDeleteRule(rule.id || rule.uuid)}
                        className="p-2 text-slate-500 hover:text-red-400 hover:bg-red-500/10 rounded-lg transition-all"
                      >
                        <Trash2 className="w-4 h-4" />
                      </button>
                    </div>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      <div className="bg-cyan-500/10 border border-cyan-500/20 p-4 rounded-xl flex items-center gap-4 text-cyan-200 text-sm">
        <span className="text-xl">💡</span>
        <p>Rules are processed top-to-bottom. First matching rule wins. Drag and drop rows to change priority.</p>
      </div>

      {showModal && (
        <RuleModal
          rule={editingRule}
          onClose={() => setShowModal(false)}
          onSave={() => { setShowModal(false); fetchRules(); }}
        />
      )}
    </div>
  );
};

interface RuleModalProps {
  rule: FirewallRule | null;
  onClose: () => void;
  onSave: () => void;
}

const RuleModal: React.FC<RuleModalProps> = ({ rule, onClose, onSave }) => {
  const [formData, setFormData] = useState({
    action: rule?.action || 'pass',
    interface: rule?.interface || 'lan',
    direction: rule?.direction || 'in',
    protocol: rule?.protocol || 'any',
    sourceNetwork: rule?.source?.network || 'any',
    sourcePort: rule?.source?.port || '',
    destNetwork: rule?.destination?.network || 'any',
    destPort: rule?.destination?.port || '',
    description: rule?.description || '',
    log: rule?.log || false,
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const token = localStorage.getItem('auth_token');
    const ruleData = {
      action: formData.action,
      interface: formData.interface,
      direction: formData.direction,
      protocol: formData.protocol,
      source: { network: formData.sourceNetwork, port: formData.sourcePort || undefined },
      destination: { network: formData.destNetwork, port: formData.destPort || undefined },
      description: formData.description,
      log: formData.log,
      state: 'keep',
      enabled: true,
    };

    const url = rule ? `/api/v1/firewall/rule/${rule.id}` : '/api/v1/firewall/rule';
    const method = rule ? 'PUT' : 'POST';

    try {
      await fetch(url, {
        method,
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(ruleData),
      });
      onSave();
    } catch {
      console.error('Failed to save rule');
    }
  };

  return (
    <div className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50">
      <div className="bg-slate-900 border border-slate-800 rounded-2xl p-6 w-full max-w-2xl max-h-[90vh] overflow-y-auto">
        <h3 className="text-xl font-bold text-white mb-6">{rule ? 'Edit Rule' : 'Create New Rule'}</h3>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="grid grid-cols-3 gap-4">
            <div>
              <label className="block text-sm font-medium text-slate-400 mb-2">Action</label>
              <select
                value={formData.action}
                onChange={(e) => setFormData({ ...formData, action: e.target.value })}
                className="w-full bg-slate-800 border border-slate-700 rounded-xl py-2.5 px-3 text-white"
              >
                <option value="pass">Pass</option>
                <option value="block">Block</option>
                <option value="reject">Reject</option>
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-400 mb-2">Interface</label>
              <select
                value={formData.interface}
                onChange={(e) => setFormData({ ...formData, interface: e.target.value })}
                className="w-full bg-slate-800 border border-slate-700 rounded-xl py-2.5 px-3 text-white"
              >
                <option value="lan">LAN</option>
                <option value="wan">WAN</option>
                <option value="opt1">OPT1</option>
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-400 mb-2">Direction</label>
              <select
                value={formData.direction}
                onChange={(e) => setFormData({ ...formData, direction: e.target.value })}
                className="w-full bg-slate-800 border border-slate-700 rounded-xl py-2.5 px-3 text-white"
              >
                <option value="in">In</option>
                <option value="out">Out</option>
              </select>
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-400 mb-2">Protocol</label>
            <select
              value={formData.protocol}
              onChange={(e) => setFormData({ ...formData, protocol: e.target.value })}
              className="w-full bg-slate-800 border border-slate-700 rounded-xl py-2.5 px-3 text-white"
            >
              <option value="any">Any</option>
              <option value="tcp">TCP</option>
              <option value="udp">UDP</option>
              <option value="icmp">ICMP</option>
            </select>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-slate-400 mb-2">Source</label>
              <input
                type="text"
                value={formData.sourceNetwork}
                onChange={(e) => setFormData({ ...formData, sourceNetwork: e.target.value })}
                placeholder="any or IP/CIDR"
                className="w-full bg-slate-800 border border-slate-700 rounded-xl py-2.5 px-3 text-white"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-400 mb-2">Source Port</label>
              <input
                type="text"
                value={formData.sourcePort}
                onChange={(e) => setFormData({ ...formData, sourcePort: e.target.value })}
                placeholder="e.g. 80 or 1024:65535"
                className="w-full bg-slate-800 border border-slate-700 rounded-xl py-2.5 px-3 text-white"
              />
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-slate-400 mb-2">Destination</label>
              <input
                type="text"
                value={formData.destNetwork}
                onChange={(e) => setFormData({ ...formData, destNetwork: e.target.value })}
                placeholder="any or IP/CIDR"
                className="w-full bg-slate-800 border border-slate-700 rounded-xl py-2.5 px-3 text-white"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-400 mb-2">Destination Port</label>
              <input
                type="text"
                value={formData.destPort}
                onChange={(e) => setFormData({ ...formData, destPort: e.target.value })}
                placeholder="e.g. 80, 443 or 1024:65535"
                className="w-full bg-slate-800 border border-slate-700 rounded-xl py-2.5 px-3 text-white"
              />
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-400 mb-2">Description</label>
            <input
              type="text"
              value={formData.description}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              placeholder="Rule description"
              className="w-full bg-slate-800 border border-slate-700 rounded-xl py-2.5 px-3 text-white"
            />
          </div>

          <div className="flex items-center gap-3">
            <input
              type="checkbox"
              id="log"
              checked={formData.log}
              onChange={(e) => setFormData({ ...formData, log: e.target.checked })}
              className="w-4 h-4 rounded border-slate-700 bg-slate-800 text-cyan-500"
            />
            <label htmlFor="log" className="text-slate-400 text-sm">Log packets matching this rule</label>
          </div>

          <div className="flex gap-3 pt-4">
            <button
              type="button"
              onClick={onClose}
              className="flex-1 px-4 py-2 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded-xl transition-all"
            >
              Cancel
            </button>
            <button
              type="submit"
              className="flex-1 px-4 py-2 bg-cyan-600 hover:bg-cyan-500 text-white font-bold rounded-xl transition-all"
            >
              {rule ? 'Update Rule' : 'Create Rule'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default FirewallRules;
