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
