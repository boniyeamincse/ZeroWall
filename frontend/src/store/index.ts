import { create } from 'zustand';
import type { User, FirewallRule, NATRule, FirewallStats } from '../types';

interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  login: (user: User, token: string) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  user: JSON.parse(localStorage.getItem('user') || 'null'),
  token: localStorage.getItem('token'),
  isAuthenticated: !!localStorage.getItem('token'),
  login: (user, token) => {
    localStorage.setItem('user', JSON.stringify(user));
    localStorage.setItem('token', token);
    set({ user, token, isAuthenticated: true });
  },
  logout: () => {
    localStorage.removeItem('user');
    localStorage.removeItem('token');
    set({ user: null, token: null, isAuthenticated: false });
  },
}));

interface FirewallState {
  rules: FirewallRule[];
  natRules: NATRule[];
  stats: FirewallStats | null;
  loading: boolean;
  error: string | null;
  setRules: (rules: FirewallRule[]) => void;
  setNATRules: (rules: NATRule[]) => void;
  setStats: (stats: FirewallStats) => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
}

export const useFirewallStore = create<FirewallState>((set) => ({
  rules: [],
  natRules: [],
  stats: null,
  loading: false,
  error: null,
  setRules: (rules) => set({ rules }),
  setNATRules: (rules) => set({ natRules: rules }),
  setStats: (stats) => set({ stats }),
  setLoading: (loading) => set({ loading }),
  setError: (error) => set({ error }),
}));

interface UIState {
  sidebarOpen: boolean;
  toggleSidebar: () => void;
}

export const useUIStore = create<UIState>((set) => ({
  sidebarOpen: true,
  toggleSidebar: () => set((state) => ({ sidebarOpen: !state.sidebarOpen })),
}));
