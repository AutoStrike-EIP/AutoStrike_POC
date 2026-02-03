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
      // Auth redirect will be implemented in Phase 3
      console.warn('Unauthorized - token cleared');
    }
    return Promise.reject(error);
  }
);

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
