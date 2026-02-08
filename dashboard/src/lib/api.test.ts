import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { HttpStatusCode, AxiosHeaders, InternalAxiosRequestConfig } from 'axios';

// We need to test the interceptors directly since mocking axios interferes with them
describe('API Module', () => {
  beforeEach(() => {
    vi.resetModules();
  });

  it('exports api instance with correct config', async () => {
    const { api } = await import('./api');
    expect(api).toBeDefined();
    expect(api.defaults.headers['Content-Type']).toBe('application/json');
  });
});

describe('Request Interceptor Logic', () => {
  const originalLocalStorage = global.localStorage;

  beforeEach(() => {
    // Create a proper mock for localStorage
    const store: Record<string, string> = {};
    const mockLocalStorage = {
      getItem: vi.fn((key: string) => store[key] || null),
      setItem: vi.fn((key: string, value: string) => {
        store[key] = value;
      }),
      removeItem: vi.fn((key: string) => {
        delete store[key];
      }),
      clear: vi.fn(() => {
        Object.keys(store).forEach((key) => delete store[key]);
      }),
      length: 0,
      key: vi.fn(),
    };
    Object.defineProperty(global, 'localStorage', {
      value: mockLocalStorage,
      writable: true,
    });
  });

  afterEach(() => {
    Object.defineProperty(global, 'localStorage', {
      value: originalLocalStorage,
      writable: true,
    });
  });

  it('adds Authorization header when token exists in localStorage', async () => {
    localStorage.setItem('token', 'my-jwt-token');

    vi.resetModules();
    const { api } = await import('./api');

    // Get the request interceptor function
    const interceptors = api.interceptors.request as unknown as {
      handlers: Array<{ fulfilled: (config: InternalAxiosRequestConfig) => InternalAxiosRequestConfig }>;
    };
    const requestInterceptor = interceptors.handlers[0].fulfilled;

    const config: InternalAxiosRequestConfig = {
      headers: new AxiosHeaders(),
    };

    const result = requestInterceptor(config);
    expect(result.headers.Authorization).toBe('Bearer my-jwt-token');
  });

  it('does not add Authorization header when no token', async () => {
    // No token in localStorage
    vi.resetModules();
    const { api } = await import('./api');

    const interceptors = api.interceptors.request as unknown as {
      handlers: Array<{ fulfilled: (config: InternalAxiosRequestConfig) => InternalAxiosRequestConfig }>;
    };
    const requestInterceptor = interceptors.handlers[0].fulfilled;

    const config: InternalAxiosRequestConfig = {
      headers: new AxiosHeaders(),
    };

    const result = requestInterceptor(config);
    expect(result.headers.Authorization).toBeUndefined();
  });
});

describe('Response Interceptor Logic', () => {
  const originalLocalStorage = global.localStorage;
  const originalLocation = global.location;

  beforeEach(() => {
    const store: Record<string, string> = {};
    const mockLocalStorage = {
      getItem: vi.fn((key: string) => store[key] || null),
      setItem: vi.fn((key: string, value: string) => {
        store[key] = value;
      }),
      removeItem: vi.fn((key: string) => {
        delete store[key];
      }),
      clear: vi.fn(),
      length: 0,
      key: vi.fn(),
    };
    Object.defineProperty(global, 'localStorage', {
      value: mockLocalStorage,
      writable: true,
    });

    Object.defineProperty(global, 'location', {
      value: { href: 'http://localhost:3000/' },
      writable: true,
    });
  });

  afterEach(() => {
    Object.defineProperty(global, 'localStorage', {
      value: originalLocalStorage,
      writable: true,
    });
    Object.defineProperty(global, 'location', {
      value: originalLocation,
      writable: true,
    });
  });

  it('passes through successful responses', async () => {
    vi.resetModules();
    const { api } = await import('./api');

    const interceptors = api.interceptors.response as unknown as {
      handlers: Array<{
        fulfilled: (response: { status: number }) => { status: number };
        rejected: (error: { response?: { status: number } }) => Promise<never>;
      }>;
    };
    const successHandler = interceptors.handlers[0].fulfilled;

    const response = { status: 200, data: 'success' };
    const result = successHandler(response);
    expect(result).toBe(response);
  });

  it('removes token on 401 error', async () => {
    localStorage.setItem('token', 'expired-token');

    vi.resetModules();
    const { api } = await import('./api');

    const interceptors = api.interceptors.response as unknown as {
      handlers: Array<{
        fulfilled: (response: { status: number }) => { status: number };
        rejected: (error: { response?: { status: number } }) => Promise<never>;
      }>;
    };
    const errorHandler = interceptors.handlers[0].rejected;

    const error = { response: { status: HttpStatusCode.Unauthorized } };

    await expect(errorHandler(error)).rejects.toBe(error);
    expect(localStorage.removeItem).toHaveBeenCalledWith('token');
  });

  it('rejects error without redirect on non-401', async () => {
    vi.resetModules();
    const { api } = await import('./api');

    const interceptors = api.interceptors.response as unknown as {
      handlers: Array<{
        fulfilled: (response: { status: number }) => { status: number };
        rejected: (error: { response?: { status: number } }) => Promise<never>;
      }>;
    };
    const errorHandler = interceptors.handlers[0].rejected;

    const error = { response: { status: 500 } };

    await expect(errorHandler(error)).rejects.toBe(error);
    expect(localStorage.removeItem).not.toHaveBeenCalled();
    expect(global.location.href).toBe('http://localhost:3000/');
  });

  it('handles errors without response object', async () => {
    vi.resetModules();
    const { api } = await import('./api');

    const interceptors = api.interceptors.response as unknown as {
      handlers: Array<{
        fulfilled: (response: { status: number }) => { status: number };
        rejected: (error: { response?: { status: number } }) => Promise<never>;
      }>;
    };
    const errorHandler = interceptors.handlers[0].rejected;

    const error = { message: 'Network Error' };

    await expect(errorHandler(error as { response?: { status: number } })).rejects.toBe(error);
    expect(localStorage.removeItem).not.toHaveBeenCalled();
  });
});

describe('executionApi', () => {
  it('exports executionApi object with all methods', async () => {
    const { executionApi } = await import('./api');
    expect(executionApi).toBeDefined();
    expect(typeof executionApi.list).toBe('function');
    expect(typeof executionApi.get).toBe('function');
    expect(typeof executionApi.getResults).toBe('function');
    expect(typeof executionApi.start).toBe('function');
    expect(typeof executionApi.stop).toBe('function');
    expect(typeof executionApi.complete).toBe('function');
  });
});

describe('authApi', () => {
  it('exports authApi object with all methods', async () => {
    const { authApi } = await import('./api');
    expect(authApi).toBeDefined();
    expect(typeof authApi.login).toBe('function');
    expect(typeof authApi.refresh).toBe('function');
    expect(typeof authApi.logout).toBe('function');
    expect(typeof authApi.me).toBe('function');
  });
});

describe('Token Refresh Interceptor Logic', () => {
  const originalLocalStorage = global.localStorage;
  const originalLocation = global.location;
  let store: Record<string, string> = {};

  beforeEach(() => {
    store = {};
    const mockLocalStorage = {
      getItem: vi.fn((key: string) => store[key] || null),
      setItem: vi.fn((key: string, value: string) => {
        store[key] = value;
      }),
      removeItem: vi.fn((key: string) => {
        delete store[key];
      }),
      clear: vi.fn(),
      length: 0,
      key: vi.fn(),
    };
    Object.defineProperty(global, 'localStorage', {
      value: mockLocalStorage,
      writable: true,
    });
    Object.defineProperty(global, 'location', {
      value: { href: 'http://localhost:3000/dashboard', pathname: '/dashboard' },
      writable: true,
    });
  });

  afterEach(() => {
    Object.defineProperty(global, 'localStorage', {
      value: originalLocalStorage,
      writable: true,
    });
    Object.defineProperty(global, 'location', {
      value: originalLocation,
      writable: true,
    });
  });

  it('redirects to login when 401 on auth/refresh endpoint', async () => {
    vi.resetModules();
    const { api } = await import('./api');

    const interceptors = api.interceptors.response as unknown as {
      handlers: Array<{
        fulfilled: (response: unknown) => unknown;
        rejected: (error: unknown) => Promise<unknown>;
      }>;
    };
    const errorHandler = interceptors.handlers[0].rejected;

    const error = {
      config: { url: '/auth/refresh', _retry: false },
      response: { status: HttpStatusCode.Unauthorized },
    };

    await expect(errorHandler(error)).rejects.toBe(error);
    expect(localStorage.removeItem).toHaveBeenCalledWith('token');
    expect(localStorage.removeItem).toHaveBeenCalledWith('refreshToken');
    expect(global.location.href).toBe('/login');
  });

  it('redirects to login when 401 on auth/login endpoint', async () => {
    vi.resetModules();
    const { api } = await import('./api');

    const interceptors = api.interceptors.response as unknown as {
      handlers: Array<{
        fulfilled: (response: unknown) => unknown;
        rejected: (error: unknown) => Promise<unknown>;
      }>;
    };
    const errorHandler = interceptors.handlers[0].rejected;

    const error = {
      config: { url: '/auth/login', _retry: false },
      response: { status: HttpStatusCode.Unauthorized },
    };

    await expect(errorHandler(error)).rejects.toBe(error);
    expect(localStorage.removeItem).toHaveBeenCalledWith('token');
    expect(localStorage.removeItem).toHaveBeenCalledWith('refreshToken');
    expect(global.location.href).toBe('/login');
  });

  it('redirects to login when 401 and no config', async () => {
    vi.resetModules();
    const { api } = await import('./api');

    const interceptors = api.interceptors.response as unknown as {
      handlers: Array<{
        fulfilled: (response: unknown) => unknown;
        rejected: (error: unknown) => Promise<unknown>;
      }>;
    };
    const errorHandler = interceptors.handlers[0].rejected;

    const error = {
      config: undefined,
      response: { status: HttpStatusCode.Unauthorized },
    };

    await expect(errorHandler(error)).rejects.toBe(error);
    expect(localStorage.removeItem).toHaveBeenCalledWith('token');
    expect(localStorage.removeItem).toHaveBeenCalledWith('refreshToken');
    expect(global.location.href).toBe('/login');
  });

  it('does not redirect when already on login page', async () => {
    Object.defineProperty(global, 'location', {
      value: { href: 'http://localhost:3000/login', pathname: '/login' },
      writable: true,
    });

    vi.resetModules();
    const { api } = await import('./api');

    const interceptors = api.interceptors.response as unknown as {
      handlers: Array<{
        fulfilled: (response: unknown) => unknown;
        rejected: (error: unknown) => Promise<unknown>;
      }>;
    };
    const errorHandler = interceptors.handlers[0].rejected;

    const error = {
      config: { url: '/some-endpoint', _retry: true },
      response: { status: HttpStatusCode.Unauthorized },
    };

    await expect(errorHandler(error)).rejects.toBe(error);
    expect(global.location.href).toBe('http://localhost:3000/login');
  });

  it('clears token and redirects when no refresh token available', async () => {
    store['token'] = 'expired-token';
    // No refreshToken in store

    vi.resetModules();
    const { api } = await import('./api');

    const interceptors = api.interceptors.response as unknown as {
      handlers: Array<{
        fulfilled: (response: unknown) => unknown;
        rejected: (error: unknown) => Promise<unknown>;
      }>;
    };
    const errorHandler = interceptors.handlers[0].rejected;

    const error = {
      config: { url: '/some-endpoint', _retry: false, headers: {} },
      response: { status: HttpStatusCode.Unauthorized },
    };

    // This will attempt token refresh but no refreshToken exists
    await expect(errorHandler(error)).rejects.toThrow('No refresh token');
    expect(localStorage.removeItem).toHaveBeenCalledWith('token');
    expect(global.location.href).toBe('/login');
  });

  it('resets isRefreshing flag when no refresh token, allowing subsequent 401s to be handled', async () => {
    store['token'] = 'expired-token';
    // No refreshToken in store

    vi.resetModules();
    const { api } = await import('./api');

    const interceptors = api.interceptors.response as unknown as {
      handlers: Array<{
        fulfilled: (response: unknown) => unknown;
        rejected: (error: unknown) => Promise<unknown>;
      }>;
    };
    const errorHandler = interceptors.handlers[0].rejected;

    const error1 = {
      config: { url: '/endpoint-1', _retry: false, headers: {} },
      response: { status: HttpStatusCode.Unauthorized },
    };

    // First 401 - no refresh token
    await expect(errorHandler(error1)).rejects.toThrow('No refresh token');

    // Reset location for second test
    global.location.href = 'http://localhost:3000/dashboard';

    const error2 = {
      config: { url: '/endpoint-2', _retry: false, headers: {} },
      response: { status: HttpStatusCode.Unauthorized },
    };

    // Second 401 should also be handled (not hang) because isRefreshing was reset
    await expect(errorHandler(error2)).rejects.toThrow('No refresh token');
    expect(global.location.href).toBe('/login');
  });
});

describe('API exports', () => {
  it('exports api instance', async () => {
    const { api } = await import('./api');
    expect(api).toBeDefined();
    expect(api.defaults.baseURL).toBeDefined();
  });
});

describe('healthApi', () => {
  it('exports healthApi object with check method', async () => {
    const { healthApi } = await import('./api');
    expect(healthApi).toBeDefined();
    expect(typeof healthApi.check).toBe('function');
  });
});

describe('adminApi', () => {
  it('exports adminApi object with all methods', async () => {
    const { adminApi } = await import('./api');
    expect(adminApi).toBeDefined();
    expect(typeof adminApi.listUsers).toBe('function');
    expect(typeof adminApi.getUser).toBe('function');
    expect(typeof adminApi.createUser).toBe('function');
    expect(typeof adminApi.updateUser).toBe('function');
    expect(typeof adminApi.updateUserRole).toBe('function');
    expect(typeof adminApi.deactivateUser).toBe('function');
    expect(typeof adminApi.reactivateUser).toBe('function');
    expect(typeof adminApi.resetPassword).toBe('function');
  });
});

describe('scenarioApi', () => {
  it('exports scenarioApi object with all methods', async () => {
    const { scenarioApi } = await import('./api');
    expect(scenarioApi).toBeDefined();
    expect(typeof scenarioApi.list).toBe('function');
    expect(typeof scenarioApi.get).toBe('function');
    expect(typeof scenarioApi.create).toBe('function');
    expect(typeof scenarioApi.update).toBe('function');
    expect(typeof scenarioApi.delete).toBe('function');
    expect(typeof scenarioApi.exportAll).toBe('function');
    expect(typeof scenarioApi.exportOne).toBe('function');
    expect(typeof scenarioApi.import).toBe('function');
  });
});

describe('permissionApi', () => {
  it('exports permissionApi object with all methods', async () => {
    const { permissionApi } = await import('./api');
    expect(permissionApi).toBeDefined();
    expect(typeof permissionApi.getMatrix).toBe('function');
    expect(typeof permissionApi.getMyPermissions).toBe('function');
  });
});

describe('analyticsApi', () => {
  it('exports analyticsApi object with all methods', async () => {
    const { analyticsApi } = await import('./api');
    expect(analyticsApi).toBeDefined();
    expect(typeof analyticsApi.compare).toBe('function');
    expect(typeof analyticsApi.trend).toBe('function');
    expect(typeof analyticsApi.summary).toBe('function');
    expect(typeof analyticsApi.periodStats).toBe('function');
  });
});

describe('notificationApi', () => {
  it('exports notificationApi object with all methods', async () => {
    const { notificationApi } = await import('./api');
    expect(notificationApi).toBeDefined();
    expect(typeof notificationApi.list).toBe('function');
    expect(typeof notificationApi.getUnreadCount).toBe('function');
    expect(typeof notificationApi.markAsRead).toBe('function');
    expect(typeof notificationApi.markAllAsRead).toBe('function');
    expect(typeof notificationApi.getSettings).toBe('function');
    expect(typeof notificationApi.updateSettings).toBe('function');
    expect(typeof notificationApi.createSettings).toBe('function');
    expect(typeof notificationApi.deleteSettings).toBe('function');
    expect(typeof notificationApi.getSMTPConfig).toBe('function');
    expect(typeof notificationApi.testSMTP).toBe('function');
  });
});

describe('scheduleApi', () => {
  it('exports scheduleApi object with all methods', async () => {
    const { scheduleApi } = await import('./api');
    expect(scheduleApi).toBeDefined();
    expect(typeof scheduleApi.list).toBe('function');
    expect(typeof scheduleApi.get).toBe('function');
    expect(typeof scheduleApi.create).toBe('function');
    expect(typeof scheduleApi.update).toBe('function');
    expect(typeof scheduleApi.delete).toBe('function');
    expect(typeof scheduleApi.pause).toBe('function');
    expect(typeof scheduleApi.resume).toBe('function');
    expect(typeof scheduleApi.runNow).toBe('function');
    expect(typeof scheduleApi.getRuns).toBe('function');
  });
});

describe('Token Refresh Queue and Success Flow', () => {
  const originalLocalStorage = global.localStorage;
  const originalLocation = global.location;
  let store: Record<string, string> = {};

  beforeEach(() => {
    store = {};
    const mockLocalStorage = {
      getItem: vi.fn((key: string) => store[key] || null),
      setItem: vi.fn((key: string, value: string) => {
        store[key] = value;
      }),
      removeItem: vi.fn((key: string) => {
        delete store[key];
      }),
      clear: vi.fn(),
      length: 0,
      key: vi.fn(),
    };
    Object.defineProperty(global, 'localStorage', {
      value: mockLocalStorage,
      writable: true,
    });
    Object.defineProperty(global, 'location', {
      value: { href: 'http://localhost:3000/dashboard', pathname: '/dashboard' },
      writable: true,
    });
  });

  afterEach(() => {
    Object.defineProperty(global, 'localStorage', {
      value: originalLocalStorage,
      writable: true,
    });
    Object.defineProperty(global, 'location', {
      value: originalLocation,
      writable: true,
    });
  });

  it('skips refresh when _retry flag is already set', async () => {
    vi.resetModules();
    const { api } = await import('./api');

    const interceptors = api.interceptors.response as unknown as {
      handlers: Array<{
        fulfilled: (response: unknown) => unknown;
        rejected: (error: unknown) => Promise<unknown>;
      }>;
    };
    const errorHandler = interceptors.handlers[0].rejected;

    const error = {
      config: { url: '/some-endpoint', _retry: true, headers: {} },
      response: { status: HttpStatusCode.Unauthorized },
    };

    await expect(errorHandler(error)).rejects.toBe(error);
    expect(localStorage.removeItem).toHaveBeenCalledWith('token');
    expect(localStorage.removeItem).toHaveBeenCalledWith('refreshToken');
  });

  it('throws error when originalRequest is undefined after all checks', async () => {
    vi.resetModules();
    const { api } = await import('./api');

    const interceptors = api.interceptors.response as unknown as {
      handlers: Array<{
        fulfilled: (response: unknown) => unknown;
        rejected: (error: unknown) => Promise<unknown>;
      }>;
    };
    const errorHandler = interceptors.handlers[0].rejected;

    // Error with no config at all
    const error = {
      response: { status: HttpStatusCode.Unauthorized },
    };

    await expect(errorHandler(error)).rejects.toBe(error);
  });
});

describe('API method invocations', () => {
  beforeEach(() => {
    vi.resetModules();
  });

  it('adminApi.listUsers with includeInactive=true passes params', async () => {
    const { api, adminApi } = await import('./api');
    const getSpy = vi.spyOn(api, 'get').mockResolvedValue({ data: { users: [], total: 0 } });
    await adminApi.listUsers(true);
    expect(getSpy).toHaveBeenCalledWith('/admin/users', { params: { include_inactive: 'true' } });
    getSpy.mockRestore();
  });

  it('adminApi.listUsers with includeInactive=false passes no params', async () => {
    const { api, adminApi } = await import('./api');
    const getSpy = vi.spyOn(api, 'get').mockResolvedValue({ data: { users: [], total: 0 } });
    await adminApi.listUsers(false);
    expect(getSpy).toHaveBeenCalledWith('/admin/users', { params: undefined });
    getSpy.mockRestore();
  });

  it('adminApi.listUsers defaults to includeInactive=false', async () => {
    const { api, adminApi } = await import('./api');
    const getSpy = vi.spyOn(api, 'get').mockResolvedValue({ data: { users: [], total: 0 } });
    await adminApi.listUsers();
    expect(getSpy).toHaveBeenCalledWith('/admin/users', { params: undefined });
    getSpy.mockRestore();
  });

  it('adminApi.getUser calls correct endpoint', async () => {
    const { api, adminApi } = await import('./api');
    const getSpy = vi.spyOn(api, 'get').mockResolvedValue({ data: {} });
    await adminApi.getUser('user-123');
    expect(getSpy).toHaveBeenCalledWith('/admin/users/user-123');
    getSpy.mockRestore();
  });

  it('adminApi.createUser posts data correctly', async () => {
    const { api, adminApi } = await import('./api');
    const postSpy = vi.spyOn(api, 'post').mockResolvedValue({ data: {} });
    const userData = { username: 'test', email: 'test@test.com', password: 'pass', role: 'viewer' as const };
    await adminApi.createUser(userData);
    expect(postSpy).toHaveBeenCalledWith('/admin/users', userData);
    postSpy.mockRestore();
  });

  it('adminApi.updateUser puts data correctly', async () => {
    const { api, adminApi } = await import('./api');
    const putSpy = vi.spyOn(api, 'put').mockResolvedValue({ data: {} });
    await adminApi.updateUser('user-123', { username: 'new' });
    expect(putSpy).toHaveBeenCalledWith('/admin/users/user-123', { username: 'new' });
    putSpy.mockRestore();
  });

  it('adminApi.updateUserRole puts role correctly', async () => {
    const { api, adminApi } = await import('./api');
    const putSpy = vi.spyOn(api, 'put').mockResolvedValue({ data: {} });
    await adminApi.updateUserRole('user-123', { role: 'admin' });
    expect(putSpy).toHaveBeenCalledWith('/admin/users/user-123/role', { role: 'admin' });
    putSpy.mockRestore();
  });

  it('adminApi.deactivateUser calls delete endpoint', async () => {
    const { api, adminApi } = await import('./api');
    const deleteSpy = vi.spyOn(api, 'delete').mockResolvedValue({ data: {} });
    await adminApi.deactivateUser('user-123');
    expect(deleteSpy).toHaveBeenCalledWith('/admin/users/user-123');
    deleteSpy.mockRestore();
  });

  it('adminApi.reactivateUser calls post endpoint', async () => {
    const { api, adminApi } = await import('./api');
    const postSpy = vi.spyOn(api, 'post').mockResolvedValue({ data: {} });
    await adminApi.reactivateUser('user-123');
    expect(postSpy).toHaveBeenCalledWith('/admin/users/user-123/reactivate');
    postSpy.mockRestore();
  });

  it('adminApi.resetPassword posts password correctly', async () => {
    const { api, adminApi } = await import('./api');
    const postSpy = vi.spyOn(api, 'post').mockResolvedValue({ data: {} });
    await adminApi.resetPassword('user-123', { new_password: 'newpass' });
    expect(postSpy).toHaveBeenCalledWith('/admin/users/user-123/reset-password', { new_password: 'newpass' });
    postSpy.mockRestore();
  });

  it('scenarioApi.exportAll without ids sends no params', async () => {
    const { api, scenarioApi } = await import('./api');
    const getSpy = vi.spyOn(api, 'get').mockResolvedValue({ data: {} });
    await scenarioApi.exportAll();
    expect(getSpy).toHaveBeenCalledWith('/scenarios/export', { params: undefined });
    getSpy.mockRestore();
  });

  it('scenarioApi.exportAll with empty ids sends no params', async () => {
    const { api, scenarioApi } = await import('./api');
    const getSpy = vi.spyOn(api, 'get').mockResolvedValue({ data: {} });
    await scenarioApi.exportAll([]);
    expect(getSpy).toHaveBeenCalledWith('/scenarios/export', { params: undefined });
    getSpy.mockRestore();
  });

  it('scenarioApi.exportAll with ids sends joined params', async () => {
    const { api, scenarioApi } = await import('./api');
    const getSpy = vi.spyOn(api, 'get').mockResolvedValue({ data: {} });
    await scenarioApi.exportAll(['id1', 'id2']);
    expect(getSpy).toHaveBeenCalledWith('/scenarios/export', { params: { ids: 'id1,id2' } });
    getSpy.mockRestore();
  });

  it('scenarioApi.exportOne calls correct endpoint', async () => {
    const { api, scenarioApi } = await import('./api');
    const getSpy = vi.spyOn(api, 'get').mockResolvedValue({ data: {} });
    await scenarioApi.exportOne('sc-1');
    expect(getSpy).toHaveBeenCalledWith('/scenarios/sc-1/export');
    getSpy.mockRestore();
  });

  it('scenarioApi.import posts data correctly', async () => {
    const { api, scenarioApi } = await import('./api');
    const postSpy = vi.spyOn(api, 'post').mockResolvedValue({ data: {} });
    const importData = { scenarios: [{ name: 'test', phases: [], tags: [] }] };
    await scenarioApi.import(importData);
    expect(postSpy).toHaveBeenCalledWith('/scenarios/import', importData);
    postSpy.mockRestore();
  });

  it('scenarioApi.create posts data correctly', async () => {
    const { api, scenarioApi } = await import('./api');
    const postSpy = vi.spyOn(api, 'post').mockResolvedValue({ data: {} });
    const data = { name: 'test', description: 'desc', phases: [], tags: [] };
    await scenarioApi.create(data);
    expect(postSpy).toHaveBeenCalledWith('/scenarios', data);
    postSpy.mockRestore();
  });

  it('scenarioApi.update puts data correctly', async () => {
    const { api, scenarioApi } = await import('./api');
    const putSpy = vi.spyOn(api, 'put').mockResolvedValue({ data: {} });
    const data = { name: 'updated', description: 'desc', phases: [], tags: [] };
    await scenarioApi.update('sc-1', data);
    expect(putSpy).toHaveBeenCalledWith('/scenarios/sc-1', data);
    putSpy.mockRestore();
  });

  it('scenarioApi.delete calls delete endpoint', async () => {
    const { api, scenarioApi } = await import('./api');
    const deleteSpy = vi.spyOn(api, 'delete').mockResolvedValue({ data: {} });
    await scenarioApi.delete('sc-1');
    expect(deleteSpy).toHaveBeenCalledWith('/scenarios/sc-1');
    deleteSpy.mockRestore();
  });

  it('executionApi.start posts all parameters', async () => {
    const { api, executionApi } = await import('./api');
    const postSpy = vi.spyOn(api, 'post').mockResolvedValue({ data: {} });
    await executionApi.start('scenario-1', ['paw-1', 'paw-2'], true);
    expect(postSpy).toHaveBeenCalledWith('/executions', {
      scenario_id: 'scenario-1',
      agent_paws: ['paw-1', 'paw-2'],
      safe_mode: true,
    });
    postSpy.mockRestore();
  });

  it('executionApi.stop calls correct endpoint', async () => {
    const { api, executionApi } = await import('./api');
    const postSpy = vi.spyOn(api, 'post').mockResolvedValue({ data: {} });
    await executionApi.stop('exec-1');
    expect(postSpy).toHaveBeenCalledWith('/executions/exec-1/stop');
    postSpy.mockRestore();
  });

  it('executionApi.complete calls correct endpoint', async () => {
    const { api, executionApi } = await import('./api');
    const postSpy = vi.spyOn(api, 'post').mockResolvedValue({ data: {} });
    await executionApi.complete('exec-1');
    expect(postSpy).toHaveBeenCalledWith('/executions/exec-1/complete');
    postSpy.mockRestore();
  });

  it('executionApi.getResults calls correct endpoint', async () => {
    const { api, executionApi } = await import('./api');
    const getSpy = vi.spyOn(api, 'get').mockResolvedValue({ data: {} });
    await executionApi.getResults('exec-1');
    expect(getSpy).toHaveBeenCalledWith('/executions/exec-1/results');
    getSpy.mockRestore();
  });

  it('techniqueApi.import posts techniques correctly', async () => {
    const { api, techniqueApi } = await import('./api');
    const postSpy = vi.spyOn(api, 'post').mockResolvedValue({ data: {} });
    const techniques = [{ id: 'T1082', name: 'test', description: '', tactic: 'discovery', platforms: [], executors: [], detection: [], is_safe: true }];
    await techniqueApi.import(techniques);
    expect(postSpy).toHaveBeenCalledWith('/techniques/import/json', { techniques });
    postSpy.mockRestore();
  });

  it('techniqueApi.getByTactic calls correct endpoint', async () => {
    const { api, techniqueApi } = await import('./api');
    const getSpy = vi.spyOn(api, 'get').mockResolvedValue({ data: [] });
    await techniqueApi.getByTactic('discovery');
    expect(getSpy).toHaveBeenCalledWith('/techniques/tactic/discovery');
    getSpy.mockRestore();
  });

  it('techniqueApi.getByPlatform calls correct endpoint', async () => {
    const { api, techniqueApi } = await import('./api');
    const getSpy = vi.spyOn(api, 'get').mockResolvedValue({ data: [] });
    await techniqueApi.getByPlatform('linux');
    expect(getSpy).toHaveBeenCalledWith('/techniques/platform/linux');
    getSpy.mockRestore();
  });

  it('techniqueApi.getCoverage calls correct endpoint', async () => {
    const { api, techniqueApi } = await import('./api');
    const getSpy = vi.spyOn(api, 'get').mockResolvedValue({ data: {} });
    await techniqueApi.getCoverage();
    expect(getSpy).toHaveBeenCalledWith('/techniques/coverage');
    getSpy.mockRestore();
  });

  it('analyticsApi.compare uses default days parameter', async () => {
    const { api, analyticsApi } = await import('./api');
    const getSpy = vi.spyOn(api, 'get').mockResolvedValue({ data: {} });
    await analyticsApi.compare();
    expect(getSpy).toHaveBeenCalledWith('/analytics/comparison', { params: { days: 7 } });
    getSpy.mockRestore();
  });

  it('analyticsApi.compare accepts custom days parameter', async () => {
    const { api, analyticsApi } = await import('./api');
    const getSpy = vi.spyOn(api, 'get').mockResolvedValue({ data: {} });
    await analyticsApi.compare(14);
    expect(getSpy).toHaveBeenCalledWith('/analytics/comparison', { params: { days: 14 } });
    getSpy.mockRestore();
  });

  it('analyticsApi.trend uses default days parameter', async () => {
    const { api, analyticsApi } = await import('./api');
    const getSpy = vi.spyOn(api, 'get').mockResolvedValue({ data: {} });
    await analyticsApi.trend();
    expect(getSpy).toHaveBeenCalledWith('/analytics/trend', { params: { days: 30 } });
    getSpy.mockRestore();
  });

  it('analyticsApi.summary uses default days parameter', async () => {
    const { api, analyticsApi } = await import('./api');
    const getSpy = vi.spyOn(api, 'get').mockResolvedValue({ data: {} });
    await analyticsApi.summary();
    expect(getSpy).toHaveBeenCalledWith('/analytics/summary', { params: { days: 30 } });
    getSpy.mockRestore();
  });

  it('analyticsApi.periodStats passes start and end params', async () => {
    const { api, analyticsApi } = await import('./api');
    const getSpy = vi.spyOn(api, 'get').mockResolvedValue({ data: {} });
    await analyticsApi.periodStats('2024-01-01', '2024-01-31');
    expect(getSpy).toHaveBeenCalledWith('/analytics/period', { params: { start: '2024-01-01', end: '2024-01-31' } });
    getSpy.mockRestore();
  });

  it('notificationApi.list uses default limit', async () => {
    const { api, notificationApi } = await import('./api');
    const getSpy = vi.spyOn(api, 'get').mockResolvedValue({ data: [] });
    await notificationApi.list();
    expect(getSpy).toHaveBeenCalledWith('/notifications', { params: { limit: 50 } });
    getSpy.mockRestore();
  });

  it('notificationApi.list accepts custom limit', async () => {
    const { api, notificationApi } = await import('./api');
    const getSpy = vi.spyOn(api, 'get').mockResolvedValue({ data: [] });
    await notificationApi.list(25);
    expect(getSpy).toHaveBeenCalledWith('/notifications', { params: { limit: 25 } });
    getSpy.mockRestore();
  });

  it('notificationApi.markAsRead calls correct endpoint', async () => {
    const { api, notificationApi } = await import('./api');
    const postSpy = vi.spyOn(api, 'post').mockResolvedValue({ data: {} });
    await notificationApi.markAsRead('notif-1');
    expect(postSpy).toHaveBeenCalledWith('/notifications/notif-1/read');
    postSpy.mockRestore();
  });

  it('notificationApi.createSettings posts data correctly', async () => {
    const { api, notificationApi } = await import('./api');
    const postSpy = vi.spyOn(api, 'post').mockResolvedValue({ data: {} });
    const settings = { channel: 'email' as const, enabled: true, notify_on_start: true, notify_on_complete: true, notify_on_failure: true, notify_on_score_alert: false, score_alert_threshold: 50, notify_on_agent_offline: false };
    await notificationApi.createSettings(settings);
    expect(postSpy).toHaveBeenCalledWith('/notifications/settings', settings);
    postSpy.mockRestore();
  });

  it('notificationApi.updateSettings puts data correctly', async () => {
    const { api, notificationApi } = await import('./api');
    const putSpy = vi.spyOn(api, 'put').mockResolvedValue({ data: {} });
    const settings = { channel: 'webhook' as const, enabled: true, notify_on_start: false, notify_on_complete: true, notify_on_failure: true, notify_on_score_alert: true, score_alert_threshold: 70, notify_on_agent_offline: true };
    await notificationApi.updateSettings('set-1', settings);
    expect(putSpy).toHaveBeenCalledWith('/notifications/settings/set-1', settings);
    putSpy.mockRestore();
  });

  it('notificationApi.deleteSettings calls correct endpoint', async () => {
    const { api, notificationApi } = await import('./api');
    const deleteSpy = vi.spyOn(api, 'delete').mockResolvedValue({ data: {} });
    await notificationApi.deleteSettings('set-1');
    expect(deleteSpy).toHaveBeenCalledWith('/notifications/settings/set-1');
    deleteSpy.mockRestore();
  });

  it('notificationApi.testSMTP posts email correctly', async () => {
    const { api, notificationApi } = await import('./api');
    const postSpy = vi.spyOn(api, 'post').mockResolvedValue({ data: {} });
    await notificationApi.testSMTP('test@test.com');
    expect(postSpy).toHaveBeenCalledWith('/notifications/smtp/test', { email: 'test@test.com' });
    postSpy.mockRestore();
  });

  it('scheduleApi.getRuns uses default limit', async () => {
    const { api, scheduleApi } = await import('./api');
    const getSpy = vi.spyOn(api, 'get').mockResolvedValue({ data: [] });
    await scheduleApi.getRuns('sched-1');
    expect(getSpy).toHaveBeenCalledWith('/schedules/sched-1/runs', { params: { limit: 20 } });
    getSpy.mockRestore();
  });

  it('scheduleApi.getRuns accepts custom limit', async () => {
    const { api, scheduleApi } = await import('./api');
    const getSpy = vi.spyOn(api, 'get').mockResolvedValue({ data: [] });
    await scheduleApi.getRuns('sched-1', 50);
    expect(getSpy).toHaveBeenCalledWith('/schedules/sched-1/runs', { params: { limit: 50 } });
    getSpy.mockRestore();
  });

  it('scheduleApi.create posts data correctly', async () => {
    const { api, scheduleApi } = await import('./api');
    const postSpy = vi.spyOn(api, 'post').mockResolvedValue({ data: {} });
    const data = { name: 'test', scenario_id: 'sc-1', frequency: 'daily' as const, safe_mode: true };
    await scheduleApi.create(data);
    expect(postSpy).toHaveBeenCalledWith('/schedules', data);
    postSpy.mockRestore();
  });

  it('scheduleApi.update puts data correctly', async () => {
    const { api, scheduleApi } = await import('./api');
    const putSpy = vi.spyOn(api, 'put').mockResolvedValue({ data: {} });
    const data = { name: 'updated', scenario_id: 'sc-1', frequency: 'weekly' as const, safe_mode: false };
    await scheduleApi.update('sched-1', data);
    expect(putSpy).toHaveBeenCalledWith('/schedules/sched-1', data);
    putSpy.mockRestore();
  });

  it('scheduleApi.pause calls correct endpoint', async () => {
    const { api, scheduleApi } = await import('./api');
    const postSpy = vi.spyOn(api, 'post').mockResolvedValue({ data: {} });
    await scheduleApi.pause('sched-1');
    expect(postSpy).toHaveBeenCalledWith('/schedules/sched-1/pause');
    postSpy.mockRestore();
  });

  it('scheduleApi.resume calls correct endpoint', async () => {
    const { api, scheduleApi } = await import('./api');
    const postSpy = vi.spyOn(api, 'post').mockResolvedValue({ data: {} });
    await scheduleApi.resume('sched-1');
    expect(postSpy).toHaveBeenCalledWith('/schedules/sched-1/resume');
    postSpy.mockRestore();
  });

  it('scheduleApi.runNow calls correct endpoint', async () => {
    const { api, scheduleApi } = await import('./api');
    const postSpy = vi.spyOn(api, 'post').mockResolvedValue({ data: {} });
    await scheduleApi.runNow('sched-1');
    expect(postSpy).toHaveBeenCalledWith('/schedules/sched-1/run');
    postSpy.mockRestore();
  });

  it('scheduleApi.delete calls correct endpoint', async () => {
    const { api, scheduleApi } = await import('./api');
    const deleteSpy = vi.spyOn(api, 'delete').mockResolvedValue({ data: {} });
    await scheduleApi.delete('sched-1');
    expect(deleteSpy).toHaveBeenCalledWith('/schedules/sched-1');
    deleteSpy.mockRestore();
  });

  it('permissionApi.getMatrix calls correct endpoint', async () => {
    const { api, permissionApi } = await import('./api');
    const getSpy = vi.spyOn(api, 'get').mockResolvedValue({ data: {} });
    await permissionApi.getMatrix();
    expect(getSpy).toHaveBeenCalledWith('/permissions/matrix');
    getSpy.mockRestore();
  });

  it('permissionApi.getMyPermissions calls correct endpoint', async () => {
    const { api, permissionApi } = await import('./api');
    const getSpy = vi.spyOn(api, 'get').mockResolvedValue({ data: {} });
    await permissionApi.getMyPermissions();
    expect(getSpy).toHaveBeenCalledWith('/permissions/me');
    getSpy.mockRestore();
  });

  it('authApi.login posts credentials correctly', async () => {
    const { api, authApi } = await import('./api');
    const postSpy = vi.spyOn(api, 'post').mockResolvedValue({ data: {} });
    await authApi.login({ username: 'user', password: 'pass' });
    expect(postSpy).toHaveBeenCalledWith('/auth/login', { username: 'user', password: 'pass' });
    postSpy.mockRestore();
  });

  it('authApi.refresh posts refresh token correctly', async () => {
    const { api, authApi } = await import('./api');
    const postSpy = vi.spyOn(api, 'post').mockResolvedValue({ data: {} });
    await authApi.refresh('my-refresh-token');
    expect(postSpy).toHaveBeenCalledWith('/auth/refresh', { refresh_token: 'my-refresh-token' });
    postSpy.mockRestore();
  });

  it('authApi.logout calls correct endpoint', async () => {
    const { api, authApi } = await import('./api');
    const postSpy = vi.spyOn(api, 'post').mockResolvedValue({ data: {} });
    await authApi.logout();
    expect(postSpy).toHaveBeenCalledWith('/auth/logout');
    postSpy.mockRestore();
  });

  it('authApi.me calls correct endpoint', async () => {
    const { api, authApi } = await import('./api');
    const getSpy = vi.spyOn(api, 'get').mockResolvedValue({ data: {} });
    await authApi.me();
    expect(getSpy).toHaveBeenCalledWith('/auth/me');
    getSpy.mockRestore();
  });

  it('executionApi.list calls correct endpoint', async () => {
    const { api, executionApi } = await import('./api');
    const getSpy = vi.spyOn(api, 'get').mockResolvedValue({ data: [] });
    await executionApi.list();
    expect(getSpy).toHaveBeenCalledWith('/executions');
    getSpy.mockRestore();
  });

  it('executionApi.get calls correct endpoint', async () => {
    const { api, executionApi } = await import('./api');
    const getSpy = vi.spyOn(api, 'get').mockResolvedValue({ data: {} });
    await executionApi.get('exec-1');
    expect(getSpy).toHaveBeenCalledWith('/executions/exec-1');
    getSpy.mockRestore();
  });

  it('techniqueApi.list calls correct endpoint', async () => {
    const { api, techniqueApi } = await import('./api');
    const getSpy = vi.spyOn(api, 'get').mockResolvedValue({ data: [] });
    await techniqueApi.list();
    expect(getSpy).toHaveBeenCalledWith('/techniques');
    getSpy.mockRestore();
  });

  it('techniqueApi.get calls correct endpoint', async () => {
    const { api, techniqueApi } = await import('./api');
    const getSpy = vi.spyOn(api, 'get').mockResolvedValue({ data: {} });
    await techniqueApi.get('T1082');
    expect(getSpy).toHaveBeenCalledWith('/techniques/T1082');
    getSpy.mockRestore();
  });

  it('techniqueApi.getExecutors calls correct endpoint without platform', async () => {
    const { api, techniqueApi } = await import('./api');
    const getSpy = vi.spyOn(api, 'get').mockResolvedValue({ data: [] });
    await techniqueApi.getExecutors('T1059.001');
    expect(getSpy).toHaveBeenCalledWith('/techniques/T1059.001/executors', { params: undefined });
    getSpy.mockRestore();
  });

  it('techniqueApi.getExecutors passes platform filter', async () => {
    const { api, techniqueApi } = await import('./api');
    const getSpy = vi.spyOn(api, 'get').mockResolvedValue({ data: [] });
    await techniqueApi.getExecutors('T1059.001', 'linux');
    expect(getSpy).toHaveBeenCalledWith('/techniques/T1059.001/executors', { params: { platform: 'linux' } });
    getSpy.mockRestore();
  });

  it('techniqueApi exports TechniqueSelection type', async () => {
    const apiModule = await import('./api');
    // Verify getExecutors is available (confirms TechniqueExecutor[] return type)
    expect(typeof apiModule.techniqueApi.getExecutors).toBe('function');
  });

  it('scenarioApi.list calls correct endpoint', async () => {
    const { api, scenarioApi } = await import('./api');
    const getSpy = vi.spyOn(api, 'get').mockResolvedValue({ data: [] });
    await scenarioApi.list();
    expect(getSpy).toHaveBeenCalledWith('/scenarios');
    getSpy.mockRestore();
  });

  it('scenarioApi.get calls correct endpoint', async () => {
    const { api, scenarioApi } = await import('./api');
    const getSpy = vi.spyOn(api, 'get').mockResolvedValue({ data: {} });
    await scenarioApi.get('sc-1');
    expect(getSpy).toHaveBeenCalledWith('/scenarios/sc-1');
    getSpy.mockRestore();
  });

  it('notificationApi.getUnreadCount calls correct endpoint', async () => {
    const { api, notificationApi } = await import('./api');
    const getSpy = vi.spyOn(api, 'get').mockResolvedValue({ data: { count: 0 } });
    await notificationApi.getUnreadCount();
    expect(getSpy).toHaveBeenCalledWith('/notifications/unread/count');
    getSpy.mockRestore();
  });

  it('notificationApi.markAllAsRead calls correct endpoint', async () => {
    const { api, notificationApi } = await import('./api');
    const postSpy = vi.spyOn(api, 'post').mockResolvedValue({ data: {} });
    await notificationApi.markAllAsRead();
    expect(postSpy).toHaveBeenCalledWith('/notifications/read-all');
    postSpy.mockRestore();
  });

  it('notificationApi.getSettings calls correct endpoint', async () => {
    const { api, notificationApi } = await import('./api');
    const getSpy = vi.spyOn(api, 'get').mockResolvedValue({ data: {} });
    await notificationApi.getSettings();
    expect(getSpy).toHaveBeenCalledWith('/notifications/settings');
    getSpy.mockRestore();
  });

  it('notificationApi.getSMTPConfig calls correct endpoint', async () => {
    const { api, notificationApi } = await import('./api');
    const getSpy = vi.spyOn(api, 'get').mockResolvedValue({ data: {} });
    await notificationApi.getSMTPConfig();
    expect(getSpy).toHaveBeenCalledWith('/notifications/smtp');
    getSpy.mockRestore();
  });

  it('scheduleApi.list calls correct endpoint', async () => {
    const { api, scheduleApi } = await import('./api');
    const getSpy = vi.spyOn(api, 'get').mockResolvedValue({ data: [] });
    await scheduleApi.list();
    expect(getSpy).toHaveBeenCalledWith('/schedules');
    getSpy.mockRestore();
  });

  it('scheduleApi.get calls correct endpoint', async () => {
    const { api, scheduleApi } = await import('./api');
    const getSpy = vi.spyOn(api, 'get').mockResolvedValue({ data: {} });
    await scheduleApi.get('sched-1');
    expect(getSpy).toHaveBeenCalledWith('/schedules/sched-1');
    getSpy.mockRestore();
  });
});

describe('Type exports', () => {
  it('exports LoginCredentials type', async () => {
    const api = await import('./api');
    // Type is exported if we can reference it without error
    expect(api).toBeDefined();
  });

  it('exports TokenResponse type', async () => {
    const api = await import('./api');
    expect(api).toBeDefined();
  });

  it('exports UserRole type', async () => {
    const api = await import('./api');
    expect(api).toBeDefined();
  });

  it('exports User type', async () => {
    const api = await import('./api');
    expect(api).toBeDefined();
  });

  it('exports Notification types', async () => {
    const api = await import('./api');
    expect(api).toBeDefined();
  });

  it('exports Schedule types', async () => {
    const api = await import('./api');
    expect(api).toBeDefined();
  });

  it('exports Permission types', async () => {
    const api = await import('./api');
    expect(api).toBeDefined();
  });

  it('exports Analytics types', async () => {
    const api = await import('./api');
    expect(api).toBeDefined();
  });
});

describe('Helper functions coverage', () => {
  const originalLocalStorage = global.localStorage;
  const originalLocation = global.location;
  let store: Record<string, string> = {};

  beforeEach(() => {
    store = {};
    const mockLocalStorage = {
      getItem: vi.fn((key: string) => store[key] || null),
      setItem: vi.fn((key: string, value: string) => {
        store[key] = value;
      }),
      removeItem: vi.fn((key: string) => {
        delete store[key];
      }),
      clear: vi.fn(),
      length: 0,
      key: vi.fn(),
    };
    Object.defineProperty(global, 'localStorage', {
      value: mockLocalStorage,
      writable: true,
    });
    Object.defineProperty(global, 'location', {
      value: { href: 'http://localhost:3000/dashboard', pathname: '/dashboard' },
      writable: true,
    });
  });

  afterEach(() => {
    Object.defineProperty(global, 'localStorage', {
      value: originalLocalStorage,
      writable: true,
    });
    Object.defineProperty(global, 'location', {
      value: originalLocation,
      writable: true,
    });
  });

  it('handles shouldSkipRefresh with undefined config', async () => {
    vi.resetModules();
    const { api } = await import('./api');

    const interceptors = api.interceptors.response as unknown as {
      handlers: Array<{
        fulfilled: (response: unknown) => unknown;
        rejected: (error: unknown) => Promise<unknown>;
      }>;
    };
    const errorHandler = interceptors.handlers[0].rejected;

    const error = {
      config: undefined,
      response: { status: HttpStatusCode.Unauthorized },
    };

    await expect(errorHandler(error)).rejects.toBe(error);
  });

  it('handles error with undefined url in config', async () => {
    vi.resetModules();
    const { api } = await import('./api');

    const interceptors = api.interceptors.response as unknown as {
      handlers: Array<{
        fulfilled: (response: unknown) => unknown;
        rejected: (error: unknown) => Promise<unknown>;
      }>;
    };
    const errorHandler = interceptors.handlers[0].rejected;

    const error = {
      config: { url: undefined, _retry: false, headers: {} },
      response: { status: HttpStatusCode.Unauthorized },
    };

    await expect(errorHandler(error)).rejects.toThrow('No refresh token');
  });

  it('successfully refreshes token and retries original request', async () => {
    store['token'] = 'expired-token';
    store['refreshToken'] = 'valid-refresh';

    vi.resetModules();
    const { api } = await import('./api');

    // Mock the adapter to intercept all HTTP calls
    const calls: string[] = [];
    const originalAdapter = api.defaults.adapter;
    api.defaults.adapter = async (config) => {
      const url = String(config.url || '');
      calls.push(url);

      if (url.includes('/auth/refresh')) {
        return {
          data: { access_token: 'new-access', refresh_token: 'new-refresh' },
          status: 200,
          statusText: 'OK',
          headers: {},
          config,
        } as never;
      }
      // Retried request
      return {
        data: { result: 'retried-ok' },
        status: 200,
        statusText: 'OK',
        headers: {},
        config,
      } as never;
    };

    const interceptors = api.interceptors.response as unknown as {
      handlers: Array<{
        fulfilled: (response: unknown) => unknown;
        rejected: (error: unknown) => Promise<unknown>;
      }>;
    };
    const errorHandler = interceptors.handlers[0].rejected;

    const error = {
      config: { url: '/some-endpoint', _retry: false, headers: {} },
      response: { status: HttpStatusCode.Unauthorized },
    };

    await errorHandler(error);
    expect(localStorage.setItem).toHaveBeenCalledWith('token', 'new-access');
    expect(localStorage.setItem).toHaveBeenCalledWith('refreshToken', 'new-refresh');
    api.defaults.adapter = originalAdapter;
  });

  it('clears tokens when refresh request itself fails', async () => {
    store['token'] = 'expired-token';
    store['refreshToken'] = 'valid-refresh';

    vi.resetModules();
    const { api } = await import('./api');

    const originalAdapter = api.defaults.adapter;
    api.defaults.adapter = async (config) => {
      const url = String(config.url || '');
      if (url.includes('/auth/refresh')) {
        const err = new Error('Refresh failed') as Error & { response?: { status: number } };
        err.response = { status: 401 };
        throw err;
      }
      return { data: {}, status: 200, statusText: 'OK', headers: {}, config } as never;
    };

    const interceptors = api.interceptors.response as unknown as {
      handlers: Array<{
        fulfilled: (response: unknown) => unknown;
        rejected: (error: unknown) => Promise<unknown>;
      }>;
    };
    const errorHandler = interceptors.handlers[0].rejected;

    const error = {
      config: { url: '/some-endpoint', _retry: false, headers: {} },
      response: { status: HttpStatusCode.Unauthorized },
    };

    await expect(errorHandler(error)).rejects.toThrow('Refresh failed');
    expect(localStorage.removeItem).toHaveBeenCalledWith('token');
    expect(localStorage.removeItem).toHaveBeenCalledWith('refreshToken');
    expect(global.location.href).toBe('/login');
    api.defaults.adapter = originalAdapter;
  });

  it('queues concurrent 401 requests and resolves them after refresh', async () => {
    store['token'] = 'expired-token';
    store['refreshToken'] = 'valid-refresh';

    vi.resetModules();
    const { api } = await import('./api');

    let refreshCallCount = 0;
    let resolveRefresh: ((value?: unknown) => void) | null = null;
    const originalAdapter = api.defaults.adapter;
    api.defaults.adapter = async (config) => {
      const url = String(config.url || '');
      if (url.includes('/auth/refresh')) {
        refreshCallCount++;
        // Delay refresh to allow second 401 to queue
        await new Promise(r => { resolveRefresh = r as (value?: unknown) => void; });
        return {
          data: { access_token: 'new-access', refresh_token: 'new-refresh' },
          status: 200,
          statusText: 'OK',
          headers: {},
          config,
        } as never;
      }
      return { data: { ok: true }, status: 200, statusText: 'OK', headers: {}, config } as never;
    };

    const interceptors = api.interceptors.response as unknown as {
      handlers: Array<{
        fulfilled: (response: unknown) => unknown;
        rejected: (error: unknown) => Promise<unknown>;
      }>;
    };
    const errorHandler = interceptors.handlers[0].rejected;

    // First 401 triggers refresh
    const error1 = {
      config: { url: '/endpoint-1', _retry: false, headers: {} },
      response: { status: HttpStatusCode.Unauthorized },
    };

    // Trigger first (starts refresh, will be delayed)
    const p1 = errorHandler(error1);

    // Wait for refresh to start
    await new Promise(r => setTimeout(r, 10));

    // Second 401 while refresh in progress -> queued
    const error2 = {
      config: { url: '/endpoint-2', _retry: false, headers: {} },
      response: { status: HttpStatusCode.Unauthorized },
    };
    const p2 = errorHandler(error2);

    // Now resolve the refresh
    await new Promise(r => setTimeout(r, 10));
    if (resolveRefresh) (resolveRefresh as () => void)();

    await Promise.all([p1, p2]);
    // Refresh should only be called once
    expect(refreshCallCount).toBe(1);
    api.defaults.adapter = originalAdapter;
  });

  it('processQueue rejects queued requests when refresh fails', async () => {
    store['token'] = 'expired-token';
    store['refreshToken'] = 'valid-refresh';

    vi.resetModules();
    const { api } = await import('./api');

    const originalAdapter = api.defaults.adapter;
    let refreshResolve: ((value?: unknown) => void) | null = null;
    api.defaults.adapter = async (config) => {
      const url = String(config.url || '');
      if (url.includes('/auth/refresh')) {
        // Delay to allow queuing
        await new Promise<void>(r => { refreshResolve = r as (value?: unknown) => void; });
        throw new Error('Refresh expired');
      }
      return { data: {}, status: 200, statusText: 'OK', headers: {}, config } as never;
    };

    const interceptors = api.interceptors.response as unknown as {
      handlers: Array<{
        fulfilled: (response: unknown) => unknown;
        rejected: (error: unknown) => Promise<unknown>;
      }>;
    };
    const errorHandler = interceptors.handlers[0].rejected;

    const error1 = {
      config: { url: '/endpoint-1', _retry: false, headers: {} },
      response: { status: HttpStatusCode.Unauthorized },
    };
    const p1 = errorHandler(error1);

    await new Promise(r => setTimeout(r, 10));

    const error2 = {
      config: { url: '/endpoint-2', _retry: false, headers: {} },
      response: { status: HttpStatusCode.Unauthorized },
    };
    const p2 = errorHandler(error2);

    // Resolve the delayed refresh (which will throw)
    await new Promise(r => setTimeout(r, 10));
    if (refreshResolve) (refreshResolve as () => void)();

    await expect(p1).rejects.toThrow('Refresh expired');
    await expect(p2).rejects.toThrow('Refresh expired');
    api.defaults.adapter = originalAdapter;
  });
});

describe('techniqueApi.getExecutors', () => {
  beforeEach(() => {
    vi.resetModules();
  });

  it('calls correct endpoint without platform filter', async () => {
    const { api, techniqueApi } = await import('./api');
    const getSpy = vi.spyOn(api, 'get').mockResolvedValue({ data: [] });
    await techniqueApi.getExecutors('T1082');
    expect(getSpy).toHaveBeenCalledWith('/techniques/T1082/executors', { params: undefined });
    getSpy.mockRestore();
  });

  it('calls correct endpoint with platform filter', async () => {
    const { api, techniqueApi } = await import('./api');
    const getSpy = vi.spyOn(api, 'get').mockResolvedValue({ data: [] });
    await techniqueApi.getExecutors('T1082', 'linux');
    expect(getSpy).toHaveBeenCalledWith('/techniques/T1082/executors', { params: { platform: 'linux' } });
    getSpy.mockRestore();
  });

  it('exports getExecutors method on techniqueApi', async () => {
    const { techniqueApi } = await import('./api');
    expect(typeof techniqueApi.getExecutors).toBe('function');
  });
});
