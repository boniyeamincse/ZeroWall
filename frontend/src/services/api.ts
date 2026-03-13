import axios, { AxiosError } from 'axios';
import type { LoginResponse, FirewallRule, NATRule, Alias, FirewallStats, APIResponse } from '../types';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (response) => response,
  (error: AxiosError) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export const authService = {
  login: async (username: string, password: string): Promise<LoginResponse> => {
    const response = await api.post<LoginResponse>('/login', { username, password });
    if (response.data.success && response.data.token) {
      localStorage.setItem('token', response.data.token);
      localStorage.setItem('user', JSON.stringify(response.data.user));
    }
    return response.data;
  },

  logout: async (): Promise<void> => {
    try {
      await api.post('/logout');
    } finally {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
    }
  },

  getStoredUser: () => {
    const userStr = localStorage.getItem('user');
    return userStr ? JSON.parse(userStr) : null;
  },

  isAuthenticated: () => {
    return !!localStorage.getItem('token');
  },
};

export const firewallService = {
  getRules: async (): Promise<FirewallRule[]> => {
    const response = await api.get<APIResponse<FirewallRule[]>>('/firewall/rules');
    return response.data.data || [];
  },

  getRule: async (id: string): Promise<FirewallRule> => {
    const response = await api.get<APIResponse<FirewallRule>>(`/firewall/rules/${id}`);
    return response.data.data!;
  },

  createRule: async (rule: Omit<FirewallRule, 'id' | 'uuid' | 'sequence' | 'created' | 'modified'>): Promise<FirewallRule> => {
    const response = await api.post<APIResponse<FirewallRule>>('/firewall/rule', rule);
    return response.data.data!;
  },

  updateRule: async (id: string, rule: Partial<FirewallRule>): Promise<FirewallRule> => {
    const response = await api.put<APIResponse<FirewallRule>>(`/firewall/rule/${id}`, rule);
    return response.data.data!;
  },

  deleteRule: async (id: string): Promise<void> => {
    await api.delete(`/firewall/rule/${id}`);
  },

  toggleRule: async (id: string): Promise<void> => {
    await api.patch(`/firewall/rule/toggle/${id}`);
  },

  reorderRules: async (rules: { id: string; sequence: number }[]): Promise<void> => {
    await api.post('/firewall/rules/reorder', { rules });
  },

  getNATRules: async (): Promise<NATRule[]> => {
    const response = await api.get<APIResponse<NATRule[]>>('/firewall/nat');
    return response.data.data || [];
  },

  createNATRule: async (rule: Omit<NATRule, 'id' | 'uuid'>): Promise<NATRule> => {
    const response = await api.post<APIResponse<NATRule>>('/firewall/nat/', rule);
    return response.data.data!;
  },

  deleteNATRule: async (id: string): Promise<void> => {
    await api.delete(`/firewall/nat/${id}`);
  },

  getAliases: async (): Promise<Alias[]> => {
    const response = await api.get<APIResponse<Alias[]>>('/firewall/aliases');
    return response.data.data || [];
  },

  getStats: async (): Promise<FirewallStats> => {
    const response = await api.get<APIResponse<FirewallStats>>('/firewall/stats');
    return response.data.data || { states: 0, rules: 0, blocked_hour: 0 };
  },

  applyFirewall: async (): Promise<void> => {
    await api.post('/firewall/apply');
  },

  flushStates: async (): Promise<void> => {
    await api.delete('/firewall/flush');
  },
};

export default api;
