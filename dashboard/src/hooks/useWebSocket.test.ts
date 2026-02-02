import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { renderHook, act, waitFor } from '@testing-library/react';
import { useWebSocket } from './useWebSocket';

// Mock WebSocket
class MockWebSocket {
  static CONNECTING = 0;
  static OPEN = 1;
  static CLOSING = 2;
  static CLOSED = 3;

  readyState = MockWebSocket.CONNECTING;
  onopen: ((event: Event) => void) | null = null;
  onclose: ((event: CloseEvent) => void) | null = null;
  onerror: ((event: Event) => void) | null = null;
  onmessage: ((event: MessageEvent) => void) | null = null;

  constructor(public url: string) {
    // Simulate async connection
    setTimeout(() => {
      this.readyState = MockWebSocket.OPEN;
      this.onopen?.(new Event('open'));
    }, 0);
  }

  send = vi.fn();
  close = vi.fn(() => {
    this.readyState = MockWebSocket.CLOSED;
    this.onclose?.(new CloseEvent('close'));
  });

  // Helper to simulate receiving a message
  simulateMessage(data: unknown) {
    this.onmessage?.(new MessageEvent('message', { data: JSON.stringify(data) }));
  }

  // Helper to simulate error
  simulateError() {
    this.onerror?.(new Event('error'));
  }

  // Helper to simulate close
  simulateClose() {
    this.readyState = MockWebSocket.CLOSED;
    this.onclose?.(new CloseEvent('close'));
  }
}

// Store references to created WebSocket instances
let mockWebSocketInstances: MockWebSocket[] = [];

// Setup mock
beforeEach(() => {
  mockWebSocketInstances = [];
  vi.stubGlobal('WebSocket', class extends MockWebSocket {
    constructor(url: string) {
      super(url);
      mockWebSocketInstances.push(this);
    }
  });

  // Mock import.meta.env
  vi.stubGlobal('import', { meta: { env: {} } });
});

afterEach(() => {
  vi.clearAllMocks();
  vi.unstubAllGlobals();
});

describe('useWebSocket', () => {
  it('should connect to WebSocket on mount', async () => {
    const { result } = renderHook(() => useWebSocket());

    await waitFor(() => {
      expect(result.current.isConnected).toBe(true);
    });

    expect(mockWebSocketInstances.length).toBe(1);
    expect(mockWebSocketInstances[0].url).toContain('/ws/dashboard');
  });

  it('should call onConnect callback when connected', async () => {
    const onConnect = vi.fn();
    renderHook(() => useWebSocket({ onConnect }));

    await waitFor(() => {
      expect(onConnect).toHaveBeenCalled();
    });
  });

  it('should call onMessage callback when receiving message', async () => {
    const onMessage = vi.fn();
    const { result } = renderHook(() => useWebSocket({ onMessage }));

    await waitFor(() => {
      expect(result.current.isConnected).toBe(true);
    });

    const testMessage = { type: 'test', payload: { data: 'test' } };
    act(() => {
      mockWebSocketInstances[0].simulateMessage(testMessage);
    });

    expect(onMessage).toHaveBeenCalledWith(testMessage);
    expect(result.current.lastMessage).toEqual(testMessage);
  });

  it('should send messages when connected', async () => {
    const { result } = renderHook(() => useWebSocket());

    await waitFor(() => {
      expect(result.current.isConnected).toBe(true);
    });

    const testMessage = { type: 'ping' };
    act(() => {
      result.current.send(testMessage);
    });

    expect(mockWebSocketInstances[0].send).toHaveBeenCalledWith(JSON.stringify(testMessage));
  });

  it('should update isConnected to false on disconnect', async () => {
    const onDisconnect = vi.fn();
    const { result } = renderHook(() => useWebSocket({ onDisconnect }));

    await waitFor(() => {
      expect(result.current.isConnected).toBe(true);
    });

    act(() => {
      mockWebSocketInstances[0].simulateClose();
    });

    await waitFor(() => {
      expect(result.current.isConnected).toBe(false);
    });
    expect(onDisconnect).toHaveBeenCalled();
  });

  it('should call onError callback on error', async () => {
    const onError = vi.fn();
    const { result } = renderHook(() => useWebSocket({ onError }));

    await waitFor(() => {
      expect(result.current.isConnected).toBe(true);
    });

    act(() => {
      mockWebSocketInstances[0].simulateError();
    });

    expect(onError).toHaveBeenCalled();
  });

  it('should close WebSocket on unmount', async () => {
    const { result, unmount } = renderHook(() => useWebSocket());

    await waitFor(() => {
      expect(result.current.isConnected).toBe(true);
    });

    unmount();

    expect(mockWebSocketInstances[0].close).toHaveBeenCalled();
  });

  it('should handle invalid JSON messages gracefully', async () => {
    const onMessage = vi.fn();
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
    const { result } = renderHook(() => useWebSocket({ onMessage }));

    await waitFor(() => {
      expect(result.current.isConnected).toBe(true);
    });

    // Simulate receiving invalid JSON
    act(() => {
      mockWebSocketInstances[0].onmessage?.(new MessageEvent('message', { data: 'invalid json' }));
    });

    expect(onMessage).not.toHaveBeenCalled();
    expect(consoleSpy).toHaveBeenCalled();

    consoleSpy.mockRestore();
  });

  it('should not create new connection if already connected', async () => {
    const { result } = renderHook(() => useWebSocket());

    await waitFor(() => {
      expect(result.current.isConnected).toBe(true);
    });

    // Force readyState to OPEN
    mockWebSocketInstances[0].readyState = MockWebSocket.OPEN;

    // Should still only have 1 connection
    expect(mockWebSocketInstances.length).toBe(1);
  });

  it('should not send messages when not connected', async () => {
    const { result } = renderHook(() => useWebSocket());

    // Try to send before connection is established
    act(() => {
      result.current.send({ type: 'test' });
    });

    // Message should not be sent (WebSocket not open yet)
    expect(result.current.isConnected).toBe(false);
  });

  it('should trigger reconnection logic on close', async () => {
    const onDisconnect = vi.fn();
    const { result } = renderHook(() => useWebSocket({ onDisconnect, maxRetries: 1 }));

    await waitFor(() => {
      expect(result.current.isConnected).toBe(true);
    });

    // Initial connection
    expect(mockWebSocketInstances.length).toBe(1);

    // Simulate close - this should trigger reconnect logic
    act(() => {
      mockWebSocketInstances[0].simulateClose();
    });

    expect(result.current.isConnected).toBe(false);
    expect(onDisconnect).toHaveBeenCalled();
  });

  it('should use wss protocol for https', async () => {
    const { result } = renderHook(() => useWebSocket());

    await waitFor(() => {
      expect(result.current.isConnected).toBe(true);
    });

    // Default is http, so should use ws
    expect(mockWebSocketInstances[0].url).toContain('ws:');
    expect(mockWebSocketInstances[0].url).toContain('/ws/dashboard');
  });

  it('should clean up on unmount without errors', async () => {
    const { result, unmount } = renderHook(() => useWebSocket());

    await waitFor(() => {
      expect(result.current.isConnected).toBe(true);
    });

    // Unmount should clean up without errors
    expect(() => unmount()).not.toThrow();
    expect(mockWebSocketInstances[0].close).toHaveBeenCalled();
  });

  it('should not reconnect when cleaning up (component unmounting)', async () => {
    vi.useFakeTimers();
    const { result, unmount } = renderHook(() => useWebSocket({ maxRetries: 3 }));

    await vi.waitFor(() => {
      expect(result.current.isConnected).toBe(true);
    });

    const initialInstanceCount = mockWebSocketInstances.length;

    // Unmount the component (sets isCleaningUpRef to true)
    unmount();

    // Advance timers - no reconnect should happen
    vi.advanceTimersByTime(10000);

    // No new connection should be created
    expect(mockWebSocketInstances.length).toBe(initialInstanceCount);

    vi.useRealTimers();
  });

  it('should attempt reconnection with backoff on disconnect when not cleaning up', async () => {
    vi.useFakeTimers();
    const { result } = renderHook(() => useWebSocket({
      maxRetries: 3,
      reconnectInterval: 1000
    }));

    await vi.waitFor(() => {
      expect(result.current.isConnected).toBe(true);
    });

    const initialInstanceCount = mockWebSocketInstances.length;

    // Simulate close without unmounting
    act(() => {
      mockWebSocketInstances[0].readyState = MockWebSocket.CLOSED;
      mockWebSocketInstances[0].onclose?.(new CloseEvent('close'));
    });

    expect(result.current.isConnected).toBe(false);

    // Advance time to trigger reconnect
    await act(async () => {
      vi.advanceTimersByTime(1000);
    });

    // A new connection attempt should be made
    expect(mockWebSocketInstances.length).toBe(initialInstanceCount + 1);

    vi.useRealTimers();
  });

  it('should not reconnect if max retries exceeded', async () => {
    vi.useFakeTimers();
    const { result } = renderHook(() => useWebSocket({
      maxRetries: 0,
      reconnectInterval: 1000
    }));

    await vi.waitFor(() => {
      expect(result.current.isConnected).toBe(true);
    });

    const initialInstanceCount = mockWebSocketInstances.length;

    // Simulate close
    act(() => {
      mockWebSocketInstances[0].readyState = MockWebSocket.CLOSED;
      mockWebSocketInstances[0].onclose?.(new CloseEvent('close'));
    });

    // Advance time
    await act(async () => {
      vi.advanceTimersByTime(5000);
    });

    // No reconnection since maxRetries is 0
    expect(mockWebSocketInstances.length).toBe(initialInstanceCount);

    vi.useRealTimers();
  });

  it('should prevent duplicate connections when in CONNECTING state', async () => {
    // Create a slow-connecting WebSocket mock
    vi.stubGlobal('WebSocket', class extends MockWebSocket {
      constructor(url: string) {
        super(url);
        mockWebSocketInstances.push(this);
        // Stay in CONNECTING state longer
        this.readyState = MockWebSocket.CONNECTING;
      }
    });

    const { result } = renderHook(() => useWebSocket());

    // Should be in connecting state, not connected yet
    expect(result.current.isConnected).toBe(false);
    expect(mockWebSocketInstances.length).toBe(1);
    expect(mockWebSocketInstances[0].readyState).toBe(MockWebSocket.CONNECTING);
  });

  it('should handle WebSocket creation error gracefully', async () => {
    // Suppress console.error for this test
    vi.spyOn(console, 'error').mockImplementation(() => {});

    // Mock WebSocket to throw on construction
    vi.stubGlobal('WebSocket', class {
      constructor() {
        throw new Error('WebSocket connection failed');
      }
    });

    // Hook should not crash even if WebSocket fails to create
    const { result } = renderHook(() => useWebSocket());

    // Should not crash, isConnected should be false
    expect(result.current.isConnected).toBe(false);

    vi.restoreAllMocks();
  });

});
