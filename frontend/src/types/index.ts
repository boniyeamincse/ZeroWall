export interface User {
  username: string;
  role: string;
}

export interface LoginResponse {
  success: boolean;
  token: string;
  user: User;
}

export interface FirewallRule {
  id: string;
  uuid?: string;
  sequence?: number;
  action: 'pass' | 'block' | 'reject';
  interface: string;
  direction: 'in' | 'out';
  protocol: string;
  source?: string;
  destination?: string;
  dst_port?: string;
  src_port?: string;
  state?: string;
  log?: boolean;
  description?: string;
  enabled: boolean;
  created?: string;
  modified?: string;
}

export interface NATRule {
  id: string;
  uuid?: string;
  interface: string;
  protocol: string;
  source?: string;
  source_port?: string;
  destination?: string;
  destination_port?: string;
  target?: string;
  target_port?: string;
  description?: string;
  enabled: boolean;
}

export interface Alias {
  id: string;
  name: string;
  type: 'host' | 'network' | 'port';
  address?: string;
  detail?: string;
}

export interface FirewallStats {
  states: number;
  rules: number;
  blocked_hour: number;
}

export interface APIResponse<T> {
  success: boolean;
  message?: string;
  data?: T;
  error?: string;
}
