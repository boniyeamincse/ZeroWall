const API_BASE = '/api/v1';

interface LoginRequest {
  username: string;
  password: string;
}

interface LoginResponse {
  token: string;
  user: {
    username: string;
    role: string;
  };
  expiresIn: number;
}

interface ApiError {
  error: string;
  code?: string;
}

class ApiService {
  private token: string | null = null;

  constructor() {
    this.token = localStorage.getItem('auth_token');
  }

  setToken(token: string) {
    this.token = token;
    localStorage.setItem('auth_token', token);
  }

  clearToken() {
    this.token = null;
    localStorage.removeItem('auth_token');
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      ...options.headers,
    };

    if (this.token) {
      (headers as Record<string, string>)['Authorization'] = `Bearer ${this.token}`;
    }

    const response = await fetch(`${API_BASE}${endpoint}`, {
      ...options,
      headers,
    });

    if (!response.ok) {
      const error: ApiError = await response.json().catch(() => ({
        error: 'An unexpected error occurred',
      }));
      throw new Error(error.error || `HTTP ${response.status}`);
    }

    return response.json();
  }

  async login(credentials: LoginRequest): Promise<LoginResponse> {
    const response = await this.request<LoginResponse>('/login', {
      method: 'POST',
      body: JSON.stringify(credentials),
    });
    this.setToken(response.token);
    return response;
  }

  async logout(): Promise<void> {
    try {
      await this.request('/logout', { method: 'POST' });
    } finally {
      this.clearToken();
    }
  }

  async refreshToken(): Promise<LoginResponse> {
    const response = await this.request<LoginResponse>('/refresh', {
      method: 'POST',
    });
    this.setToken(response.token);
    return response;
  }

  async getStatus() {
    return this.request<{ status: string; version: string }>('/status');
  }

  async getHealth() {
    return this.request<{ healthy: boolean; timestamp: number }>('/health');
  }

  async getFirewallRules() {
    return this.request<{ rules: unknown[] }>('/firewall/rules');
  }

  async createFirewallRule(rule: unknown) {
    return this.request('/firewall/rules', {
      method: 'POST',
      body: JSON.stringify(rule),
    });
  }

  async updateFirewallRule(id: string, rule: unknown) {
    return this.request(`/firewall/rules/${id}`, {
      method: 'PUT',
      body: JSON.stringify(rule),
    });
  }

  async deleteFirewallRule(id: string) {
    return this.request(`/firewall/rules/${id}`, {
      method: 'DELETE',
    });
  }

  async getVPNStatus() {
    return this.request<{ connections: unknown[] }>('/vpn/status');
  }
}

export const api = new ApiService();
export default api;
