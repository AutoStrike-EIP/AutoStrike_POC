import axios, { HttpStatusCode } from 'axios';

export const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor for auth token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Response interceptor for error handling
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === HttpStatusCode.Unauthorized) {
      localStorage.removeItem('token');
      localStorage.removeItem('refreshToken');
      // Redirect to login if not already there
      if (globalThis.location.pathname !== '/login') {
        globalThis.location.href = '/login';
      }
    }
    return Promise.reject(error);
  }
);

// Auth types
export interface LoginCredentials {
  username: string;
  password: string;
}

export interface TokenResponse {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  token_type: string;
}

export interface User {
  id: string;
  username: string;
  email: string;
  role: 'admin' | 'operator' | 'viewer';
  created_at: string;
  updated_at: string;
}

// Auth API methods
export const authApi = {
  /**
   * Login with username and password
   */
  login: (credentials: LoginCredentials) =>
    api.post<TokenResponse>('/auth/login', credentials),

  /**
   * Refresh access token using refresh token
   */
  refresh: (refreshToken: string) =>
    api.post<TokenResponse>('/auth/refresh', { refresh_token: refreshToken }),

  /**
   * Logout (client-side token removal)
   */
  logout: () => api.post('/auth/logout'),

  /**
   * Get current authenticated user
   */
  me: () => api.get<User>('/auth/me'),
};

// Execution API methods
export const executionApi = {
  /**
   * List all executions
   */
  list: () => api.get('/executions'),

  /**
   * Get execution by ID
   */
  get: (id: string) => api.get(`/executions/${id}`),

  /**
   * Get execution results
   */
  getResults: (id: string) => api.get(`/executions/${id}/results`),

  /**
   * Start a new execution
   */
  start: (scenarioId: string, agentPaws: string[], safeMode: boolean) =>
    api.post('/executions', {
      scenario_id: scenarioId,
      agent_paws: agentPaws,
      safe_mode: safeMode,
    }),

  /**
   * Stop a running execution
   */
  stop: (id: string) => api.post(`/executions/${id}/stop`),

  /**
   * Complete an execution
   */
  complete: (id: string) => api.post(`/executions/${id}/complete`),
};
