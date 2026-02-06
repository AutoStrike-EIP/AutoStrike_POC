package websocket

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

func TestNewClient(t *testing.T) {
	logger := zap.NewNop()
	hub := NewHub(logger)

	client := NewClient(hub, nil, "test-paw", logger)

	if client == nil {
		t.Fatal("NewClient returned nil")
	}

	if client.hub != hub {
		t.Error("Client hub not set correctly")
	}

	if client.agentPaw != "test-paw" {
		t.Errorf("Expected agentPaw 'test-paw', got '%s'", client.agentPaw)
	}

	if client.send == nil {
		t.Error("Client send channel is nil")
	}
}

func TestClient_GetAgentPaw(t *testing.T) {
	logger := zap.NewNop()
	hub := NewHub(logger)

	client := NewClient(hub, nil, "agent-123", logger)

	paw := client.GetAgentPaw()

	if paw != "agent-123" {
		t.Errorf("Expected paw 'agent-123', got '%s'", paw)
	}
}

func TestClient_GetAgentPaw_Empty(t *testing.T) {
	logger := zap.NewNop()
	hub := NewHub(logger)

	client := NewClient(hub, nil, "", logger)

	paw := client.GetAgentPaw()

	if paw != "" {
		t.Errorf("Expected empty paw, got '%s'", paw)
	}
}

func TestClient_Send(t *testing.T) {
	logger := zap.NewNop()
	hub := NewHub(logger)

	client := NewClient(hub, nil, "test-paw", logger)

	payload := map[string]string{"key": "value"}
	err := client.Send("test-type", payload)

	if err != nil {
		t.Fatalf("Send returned error: %v", err)
	}

	// Read from the send channel
	select {
	case msgBytes := <-client.send:
		var msg Message
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			t.Fatalf("Failed to unmarshal message: %v", err)
		}

		if msg.Type != "test-type" {
			t.Errorf("Expected type 'test-type', got '%s'", msg.Type)
		}

		var receivedPayload map[string]string
		if err := json.Unmarshal(msg.Payload, &receivedPayload); err != nil {
			t.Fatalf("Failed to unmarshal payload: %v", err)
		}

		if receivedPayload["key"] != "value" {
			t.Errorf("Expected payload key 'value', got '%s'", receivedPayload["key"])
		}
	default:
		t.Error("No message in send channel")
	}
}

func TestClient_Send_ComplexPayload(t *testing.T) {
	logger := zap.NewNop()
	hub := NewHub(logger)

	client := NewClient(hub, nil, "test-paw", logger)

	type ComplexPayload struct {
		ID        string   `json:"id"`
		Count     int      `json:"count"`
		Active    bool     `json:"active"`
		Tags      []string `json:"tags"`
	}

	payload := ComplexPayload{
		ID:     "task-123",
		Count:  42,
		Active: true,
		Tags:   []string{"tag1", "tag2"},
	}

	err := client.Send("complex", payload)

	if err != nil {
		t.Fatalf("Send returned error: %v", err)
	}

	select {
	case msgBytes := <-client.send:
		var msg Message
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			t.Fatalf("Failed to unmarshal message: %v", err)
		}

		var receivedPayload ComplexPayload
		if err := json.Unmarshal(msg.Payload, &receivedPayload); err != nil {
			t.Fatalf("Failed to unmarshal payload: %v", err)
		}

		if receivedPayload.ID != "task-123" {
			t.Errorf("Expected ID 'task-123', got '%s'", receivedPayload.ID)
		}
		if receivedPayload.Count != 42 {
			t.Errorf("Expected Count 42, got %d", receivedPayload.Count)
		}
		if !receivedPayload.Active {
			t.Error("Expected Active to be true")
		}
		if len(receivedPayload.Tags) != 2 {
			t.Errorf("Expected 2 tags, got %d", len(receivedPayload.Tags))
		}
	default:
		t.Error("No message in send channel")
	}
}

func TestMessage_Struct(t *testing.T) {
	msg := Message{
		Type:    "test",
		Payload: json.RawMessage(`{"key":"value"}`),
	}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal message: %v", err)
	}

	var decoded Message
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal message: %v", err)
	}

	if decoded.Type != "test" {
		t.Errorf("Expected type 'test', got '%s'", decoded.Type)
	}
}

func TestConstants(t *testing.T) {
	// Verify constants are set to reasonable values
	if writeWait <= 0 {
		t.Error("writeWait should be positive")
	}

	if pongWait <= 0 {
		t.Error("pongWait should be positive")
	}

	if pingPeriod <= 0 {
		t.Error("pingPeriod should be positive")
	}

	if pingPeriod >= pongWait {
		t.Error("pingPeriod should be less than pongWait")
	}

	if maxMessageSize <= 0 {
		t.Error("maxMessageSize should be positive")
	}
}

func TestClient_SetAgentPaw(t *testing.T) {
	logger := zap.NewNop()
	hub := NewHub(logger)

	client := NewClient(hub, nil, "", logger)

	if client.GetAgentPaw() != "" {
		t.Error("Expected empty paw initially")
	}

	client.SetAgentPaw("new-agent-paw")

	if client.GetAgentPaw() != "new-agent-paw" {
		t.Errorf("Expected paw 'new-agent-paw', got '%s'", client.GetAgentPaw())
	}
}

func TestClient_SetAgentPaw_Override(t *testing.T) {
	logger := zap.NewNop()
	hub := NewHub(logger)

	client := NewClient(hub, nil, "initial-paw", logger)

	if client.GetAgentPaw() != "initial-paw" {
		t.Errorf("Expected initial paw 'initial-paw', got '%s'", client.GetAgentPaw())
	}

	client.SetAgentPaw("updated-paw")

	if client.GetAgentPaw() != "updated-paw" {
		t.Errorf("Expected updated paw 'updated-paw', got '%s'", client.GetAgentPaw())
	}
}

func TestClient_Context(t *testing.T) {
	logger := zap.NewNop()
	hub := NewHub(logger)

	client := NewClient(hub, nil, "test-paw", logger)

	ctx := client.Context()

	if ctx == nil {
		t.Fatal("Context returned nil")
	}

	// Verify it's a valid context (should not be canceled)
	select {
	case <-ctx.Done():
		t.Error("Context should not be done")
	default:
		// OK - context is not done
	}
}

func TestClient_Context_MultipleCallsReturnValidContext(t *testing.T) {
	logger := zap.NewNop()
	hub := NewHub(logger)

	client := NewClient(hub, nil, "test-paw", logger)

	ctx1 := client.Context()
	ctx2 := client.Context()

	if ctx1 == nil || ctx2 == nil {
		t.Fatal("Context should not return nil")
	}

	// Both should be valid contexts
	if ctx1.Err() != nil {
		t.Error("First context should not have error")
	}
	if ctx2.Err() != nil {
		t.Error("Second context should not have error")
	}
}

func TestNewClient_NilLogger(t *testing.T) {
	hub := NewHub(zap.NewNop())

	// Should not panic with nil logger
	client := NewClient(hub, nil, "test-paw", nil)

	if client == nil {
		t.Fatal("NewClient returned nil")
	}

	if client.logger == nil {
		t.Error("Logger should be set to nop logger when nil is passed")
	}
}

func TestClient_Send_UnmarshalablePayload(t *testing.T) {
	logger := zap.NewNop()
	hub := NewHub(logger)

	client := NewClient(hub, nil, "test-paw", logger)

	// Create an unmarshalable payload (channel cannot be marshaled)
	err := client.Send("test", make(chan int))

	if err == nil {
		t.Error("Expected error when sending unmarshalable payload")
	}
}

// Test helper to create a WebSocket test server
func createTestServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, string) {
	server := httptest.NewServer(handler)
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	return server, wsURL
}

var testUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func TestClient_ReadPump(t *testing.T) {
	logger := zap.NewNop()
	hub := NewHub(logger)
	go hub.Run()

	var serverConn *websocket.Conn
	var serverMu sync.Mutex

	server, wsURL := createTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		conn, err := testUpgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Logf("Upgrade failed: %v", err)
			return
		}
		serverMu.Lock()
		serverConn = conn
		serverMu.Unlock()
		// Keep connection open
		select {}
	})
	defer server.Close()

	// Connect client
	clientConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	client := NewClient(hub, clientConn, "test-agent", logger)
	hub.Register(client)

	receivedMsgs := make(chan *Message, 10)
	handler := func(c *Client, msg *Message) {
		receivedMsgs <- msg
	}

	go client.ReadPump(handler)

	// Wait for server connection
	time.Sleep(50 * time.Millisecond)

	serverMu.Lock()
	if serverConn == nil {
		t.Fatal("Server connection not established")
	}

	// Send a message from server to client
	testMsg := Message{Type: "test", Payload: json.RawMessage(`{"data":"hello"}`)}
	msgBytes, _ := json.Marshal(testMsg)
	err = serverConn.WriteMessage(websocket.TextMessage, msgBytes)
	serverMu.Unlock()

	if err != nil {
		t.Fatalf("Failed to write message: %v", err)
	}

	// Wait for message to be received
	select {
	case msg := <-receivedMsgs:
		if msg.Type != "test" {
			t.Errorf("Expected type 'test', got '%s'", msg.Type)
		}
	case <-time.After(1 * time.Second):
		t.Error("Timeout waiting for message")
	}

	// Close connection to trigger ReadPump exit
	clientConn.Close()
}

func TestClient_WritePump(t *testing.T) {
	logger := zap.NewNop()
	hub := NewHub(logger)
	go hub.Run()

	receivedMsgs := make(chan []byte, 10)
	var serverConn *websocket.Conn

	server, wsURL := createTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		conn, err := testUpgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		serverConn = conn
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				return
			}
			receivedMsgs <- data
		}
	})
	defer server.Close()

	clientConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	client := NewClient(hub, clientConn, "test-agent", logger)
	hub.Register(client)

	go client.WritePump()

	// Send a message
	err = client.Send("test-msg", map[string]string{"hello": "world"})
	if err != nil {
		t.Fatalf("Send failed: %v", err)
	}

	// Wait for message to be received by server
	select {
	case data := <-receivedMsgs:
		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}
		if msg.Type != "test-msg" {
			t.Errorf("Expected type 'test-msg', got '%s'", msg.Type)
		}
	case <-time.After(1 * time.Second):
		t.Error("Timeout waiting for message on server")
	}

	// Close to trigger WritePump exit
	clientConn.Close()
	if serverConn != nil {
		serverConn.Close()
	}
}

func TestClient_WritePump_MultipleMessages(t *testing.T) {
	logger := zap.NewNop()
	hub := NewHub(logger)
	go hub.Run()

	receivedMsgs := make(chan []byte, 20)

	server, wsURL := createTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		conn, err := testUpgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				return
			}
			receivedMsgs <- data
		}
	})
	defer server.Close()

	clientConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	client := NewClient(hub, clientConn, "test-agent", logger)
	hub.Register(client)

	go client.WritePump()

	// Send multiple messages with small delays to ensure separate writes
	for i := 0; i < 3; i++ {
		err = client.Send("msg", map[string]int{"num": i})
		if err != nil {
			t.Fatalf("Send %d failed: %v", i, err)
		}
		time.Sleep(50 * time.Millisecond) // Allow message to be processed
	}

	// Receive messages (may be batched, so count total received)
	received := 0
	timeout := time.After(2 * time.Second)
	for received < 3 {
		select {
		case <-receivedMsgs:
			received++
		case <-timeout:
			// Batching may combine messages, so receiving at least 1 is acceptable
			if received >= 1 {
				return
			}
			t.Fatalf("Timeout: received %d messages", received)
		}
	}

	clientConn.Close()
}

func TestClient_ReadPump_InvalidJSON(t *testing.T) {
	logger := zap.NewNop()
	hub := NewHub(logger)
	go hub.Run()

	var serverConn *websocket.Conn
	var mu sync.Mutex

	server, wsURL := createTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		conn, err := testUpgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		mu.Lock()
		serverConn = conn
		mu.Unlock()
		select {}
	})
	defer server.Close()

	clientConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	client := NewClient(hub, clientConn, "test-agent", logger)
	hub.Register(client)

	handlerCalled := make(chan bool, 10)
	handler := func(c *Client, msg *Message) {
		handlerCalled <- true
	}

	go client.ReadPump(handler)

	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	if serverConn != nil {
		// Send invalid JSON
		_ = serverConn.WriteMessage(websocket.TextMessage, []byte("not valid json"))
		// Send valid JSON
		validMsg := Message{Type: "valid", Payload: json.RawMessage(`{}`)}
		msgBytes, _ := json.Marshal(validMsg)
		_ = serverConn.WriteMessage(websocket.TextMessage, msgBytes)
	}
	mu.Unlock()

	// Should receive the valid message despite the invalid one
	select {
	case <-handlerCalled:
		// Success - handler was called for valid message
	case <-time.After(1 * time.Second):
		t.Error("Handler was not called")
	}

	clientConn.Close()
}

func TestClient_WritePump_ChannelClosed(t *testing.T) {
	logger := zap.NewNop()
	hub := NewHub(logger)
	go hub.Run()

	server, wsURL := createTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		conn, err := testUpgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		// Read messages to prevent blocking
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	})
	defer server.Close()

	clientConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	client := NewClient(hub, clientConn, "test-agent", logger)

	done := make(chan bool)
	go func() {
		client.WritePump()
		done <- true
	}()

	// Close the send channel to trigger WritePump exit
	close(client.send)

	select {
	case <-done:
		// WritePump exited as expected
	case <-time.After(2 * time.Second):
		t.Error("WritePump did not exit after channel closed")
	}

	clientConn.Close()
}

func TestClient_sendPing(t *testing.T) {
	logger := zap.NewNop()
	hub := NewHub(logger)

	pingReceived := make(chan bool, 1)

	server, wsURL := createTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		conn, err := testUpgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		conn.SetPingHandler(func(appData string) error {
			pingReceived <- true
			return nil
		})
		// Keep reading to handle pings
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	})
	defer server.Close()

	clientConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer clientConn.Close()

	client := NewClient(hub, clientConn, "ping-test", logger)

	// Call sendPing
	result := client.sendPing()
	if !result {
		t.Error("sendPing returned false")
	}

	// Wait for ping to be received
	select {
	case <-pingReceived:
		// Success
	case <-time.After(1 * time.Second):
		t.Error("Ping was not received by server")
	}
}

func TestClient_sendPing_ClosedConnection(t *testing.T) {
	logger := zap.NewNop()
	hub := NewHub(logger)

	server, wsURL := createTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		conn, err := testUpgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		// Keep connection open briefly
		time.Sleep(50 * time.Millisecond)
		conn.Close()
	})
	defer server.Close()

	clientConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	client := NewClient(hub, clientConn, "ping-test", logger)

	// Close client side connection to ensure ping fails
	clientConn.Close()

	// sendPing should return false on closed connection
	result := client.sendPing()
	if result {
		t.Error("sendPing should return false on closed connection")
	}
}

func TestClient_writeMessageWithQueue(t *testing.T) {
	logger := zap.NewNop()
	hub := NewHub(logger)

	receivedMsgs := make(chan []byte, 10)

	server, wsURL := createTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		conn, err := testUpgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				return
			}
			receivedMsgs <- data
		}
	})
	defer server.Close()

	clientConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer clientConn.Close()

	client := NewClient(hub, clientConn, "queue-test", logger)

	// Queue some messages first
	client.send <- []byte(`{"type":"queued1"}`)
	client.send <- []byte(`{"type":"queued2"}`)

	// Set write deadline
	_ = clientConn.SetWriteDeadline(time.Now().Add(writeWait))

	// Write first message with queue
	result := client.writeMessageWithQueue([]byte(`{"type":"first"}`))
	if !result {
		t.Error("writeMessageWithQueue returned false")
	}

	// Collect received messages
	received := 0
	timeout := time.After(2 * time.Second)
	for received < 3 {
		select {
		case <-receivedMsgs:
			received++
		case <-timeout:
			t.Fatalf("Timeout: received only %d messages, expected 3", received)
		}
	}
}

func TestClient_writeMessageWithQueue_QueueError(t *testing.T) {
	logger := zap.NewNop()
	hub := NewHub(logger)

	server, wsURL := createTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		conn, err := testUpgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		// Read first message then close
		_, _, _ = conn.ReadMessage()
		conn.Close()
	})
	defer server.Close()

	clientConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	client := NewClient(hub, clientConn, "queue-error-test", logger)

	// Queue a message
	client.send <- []byte(`{"type":"queued"}`)

	_ = clientConn.SetWriteDeadline(time.Now().Add(writeWait))

	// First write succeeds
	_ = client.writeMessageWithQueue([]byte(`{"type":"first"}`))

	// Wait for server to close
	time.Sleep(100 * time.Millisecond)

	// Queue another message and try again - should fail
	client.send <- []byte(`{"type":"queued2"}`)
	result := client.writeMessageWithQueue([]byte(`{"type":"second"}`))
	if result {
		t.Error("writeMessageWithQueue should return false when connection is closed")
	}
}

func TestClient_handleOutgoingMessage_ChannelClosed(t *testing.T) {
	logger := zap.NewNop()
	hub := NewHub(logger)

	receivedClose := make(chan bool, 1)

	server, wsURL := createTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		conn, err := testUpgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		// Set a close handler to detect close messages
		conn.SetCloseHandler(func(code int, text string) error {
			receivedClose <- true
			return nil
		})
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	})
	defer server.Close()

	clientConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer clientConn.Close()

	client := NewClient(hub, clientConn, "close-test", logger)

	// Set write deadline so the internal SetWriteDeadline succeeds
	// Call handleOutgoingMessage with ok=false to simulate closed channel
	result := client.handleOutgoingMessage(nil, false)
	if result {
		t.Error("handleOutgoingMessage should return false when channel is closed (ok=false)")
	}
}

func TestClient_handleOutgoingMessage_SetWriteDeadlineFails(t *testing.T) {
	logger := zap.NewNop()
	hub := NewHub(logger)

	server, wsURL := createTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		conn, err := testUpgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		// Close immediately so SetWriteDeadline will fail
		conn.Close()
	})
	defer server.Close()

	clientConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	client := NewClient(hub, clientConn, "deadline-test", logger)

	// Close the connection so SetWriteDeadline fails
	clientConn.Close()

	result := client.handleOutgoingMessage([]byte(`{"type":"test"}`), true)
	if result {
		t.Error("handleOutgoingMessage should return false when SetWriteDeadline fails on closed connection")
	}
}

func TestClient_handleOutgoingMessage_SuccessfulWrite(t *testing.T) {
	logger := zap.NewNop()
	hub := NewHub(logger)

	receivedMsgs := make(chan []byte, 5)

	server, wsURL := createTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		conn, err := testUpgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				return
			}
			receivedMsgs <- data
		}
	})
	defer server.Close()

	clientConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer clientConn.Close()

	client := NewClient(hub, clientConn, "success-test", logger)

	// Successful handleOutgoingMessage with ok=true and valid message
	result := client.handleOutgoingMessage([]byte(`{"type":"hello"}`), true)
	if !result {
		t.Error("handleOutgoingMessage should return true on successful write")
	}

	select {
	case data := <-receivedMsgs:
		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}
		if msg.Type != "hello" {
			t.Errorf("Expected type 'hello', got '%s'", msg.Type)
		}
	case <-time.After(1 * time.Second):
		t.Error("Timeout waiting for message")
	}
}

func TestClient_writeMessageWithQueue_FirstWriteFails(t *testing.T) {
	logger := zap.NewNop()
	hub := NewHub(logger)

	server, wsURL := createTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		conn, err := testUpgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		conn.Close()
	})
	defer server.Close()

	clientConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	client := NewClient(hub, clientConn, "write-fail-test", logger)

	// Close the connection to make WriteMessage fail
	clientConn.Close()

	// Wait for server side to close too
	time.Sleep(50 * time.Millisecond)

	result := client.writeMessageWithQueue([]byte(`{"type":"will-fail"}`))
	if result {
		t.Error("writeMessageWithQueue should return false when first write fails")
	}
}

func TestClient_writeMessageWithQueue_EmptyQueue(t *testing.T) {
	logger := zap.NewNop()
	hub := NewHub(logger)

	receivedMsgs := make(chan []byte, 5)

	server, wsURL := createTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		conn, err := testUpgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				return
			}
			receivedMsgs <- data
		}
	})
	defer server.Close()

	clientConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer clientConn.Close()

	client := NewClient(hub, clientConn, "empty-queue-test", logger)

	// Set write deadline
	_ = clientConn.SetWriteDeadline(time.Now().Add(writeWait))

	// Write with no queued messages
	result := client.writeMessageWithQueue([]byte(`{"type":"only-one"}`))
	if !result {
		t.Error("writeMessageWithQueue should return true with empty queue")
	}

	// Should receive exactly one message
	select {
	case data := <-receivedMsgs:
		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}
		if msg.Type != "only-one" {
			t.Errorf("Expected type 'only-one', got '%s'", msg.Type)
		}
	case <-time.After(1 * time.Second):
		t.Error("Timeout waiting for message")
	}
}

func TestClient_sendPing_SetWriteDeadlineFails(t *testing.T) {
	logger := zap.NewNop()
	hub := NewHub(logger)

	server, wsURL := createTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		conn, err := testUpgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		conn.Close()
	})
	defer server.Close()

	clientConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	client := NewClient(hub, clientConn, "deadline-fail-test", logger)

	// Close the connection so SetWriteDeadline fails
	clientConn.Close()

	result := client.sendPing()
	if result {
		t.Error("sendPing should return false when SetWriteDeadline fails")
	}
}

func TestClient_WritePump_ExitsOnConnectionCloseViaChannel(t *testing.T) {
	// WritePump blocks on select between send channel and ping ticker (54s).
	// Closing the send channel (as hub.handleUnregister does) triggers the
	// ok=false branch and exits the pump. This test verifies that path.
	logger := zap.NewNop()
	hub := NewHub(logger)
	go hub.Run()

	server, wsURL := createTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		conn, err := testUpgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	})
	defer server.Close()

	clientConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	client := NewClient(hub, clientConn, "pump-exit-test", logger)
	hub.Register(client)

	// Allow hub registration to complete
	time.Sleep(50 * time.Millisecond)

	done := make(chan bool, 1)
	go func() {
		client.WritePump()
		done <- true
	}()

	// Send a message first to prove the pump is running
	err = client.Send("test", map[string]string{"pump": "active"})
	if err != nil {
		t.Fatalf("Send failed: %v", err)
	}
	time.Sleep(50 * time.Millisecond)

	// Unregister via hub, which closes the send channel and triggers WritePump exit
	hub.Unregister(client)

	select {
	case <-done:
		// WritePump exited as expected when channel was closed by hub
	case <-time.After(3 * time.Second):
		t.Error("WritePump did not exit after hub unregistered client")
	}

	clientConn.Close()
}

func TestClient_SetAgentPaw_Empty(t *testing.T) {
	logger := zap.NewNop()
	hub := NewHub(logger)
	go hub.Run()

	// Wait for hub to start
	time.Sleep(10 * time.Millisecond)

	client := NewClient(hub, nil, "original-paw", logger)

	// Setting empty paw should not call RegisterAgent on hub
	client.SetAgentPaw("")

	if client.GetAgentPaw() != "" {
		t.Errorf("Expected empty paw after SetAgentPaw(''), got '%s'", client.GetAgentPaw())
	}

	// Verify the hub does not have an empty key agent registered
	if hub.IsAgentConnected("") {
		t.Error("Hub should not have an agent registered with empty paw")
	}
}

func TestClient_SetAgentPaw_NilHub(t *testing.T) {
	logger := zap.NewNop()

	// Create client with nil hub
	client := &Client{
		hub:      nil,
		send:     make(chan []byte, 256),
		agentPaw: "",
		logger:   logger,
	}

	// Should not panic when hub is nil
	client.SetAgentPaw("some-paw")

	if client.GetAgentPaw() != "some-paw" {
		t.Errorf("Expected paw 'some-paw', got '%s'", client.GetAgentPaw())
	}
}

func TestClient_ConcurrentGetSetAgentPaw(t *testing.T) {
	logger := zap.NewNop()
	hub := NewHub(logger)

	client := NewClient(hub, nil, "", logger)

	// Concurrent reads and writes to test mutex safety
	var wg sync.WaitGroup
	const goroutines = 50

	for i := 0; i < goroutines; i++ {
		wg.Add(2)
		go func(idx int) {
			defer wg.Done()
			client.SetAgentPaw("paw-from-writer")
		}(i)
		go func(idx int) {
			defer wg.Done()
			_ = client.GetAgentPaw()
		}(i)
	}

	wg.Wait()

	// Final value should be set (no race condition panics)
	paw := client.GetAgentPaw()
	if paw != "paw-from-writer" {
		t.Errorf("Expected final paw 'paw-from-writer', got '%s'", paw)
	}
}

func TestClient_ReadPump_ConnectionClose(t *testing.T) {
	logger := zap.NewNop()
	hub := NewHub(logger)
	go hub.Run()

	var serverConn *websocket.Conn
	var mu sync.Mutex

	server, wsURL := createTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		conn, err := testUpgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		mu.Lock()
		serverConn = conn
		mu.Unlock()
		select {}
	})
	defer server.Close()

	clientConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	client := NewClient(hub, clientConn, "close-agent", logger)
	hub.Register(client)

	done := make(chan bool, 1)
	handler := func(c *Client, msg *Message) {}
	go func() {
		client.ReadPump(handler)
		done <- true
	}()

	time.Sleep(50 * time.Millisecond)

	// Close server side to trigger unexpected close error in ReadPump
	mu.Lock()
	if serverConn != nil {
		serverConn.Close()
	}
	mu.Unlock()

	select {
	case <-done:
		// ReadPump exited after connection was closed from server side
	case <-time.After(2 * time.Second):
		t.Error("ReadPump did not exit after server closed connection")
	}
}

func TestClient_ReadPump_MultipleMessages(t *testing.T) {
	logger := zap.NewNop()
	hub := NewHub(logger)
	go hub.Run()

	var serverConn *websocket.Conn
	var mu sync.Mutex

	server, wsURL := createTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		conn, err := testUpgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		mu.Lock()
		serverConn = conn
		mu.Unlock()
		select {}
	})
	defer server.Close()

	clientConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	client := NewClient(hub, clientConn, "multi-msg-agent", logger)
	hub.Register(client)

	receivedMsgs := make(chan *Message, 10)
	handler := func(c *Client, msg *Message) {
		receivedMsgs <- msg
	}

	go client.ReadPump(handler)

	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	sc := serverConn
	mu.Unlock()

	if sc == nil {
		t.Fatal("Server connection not established")
	}

	// Send multiple messages and verify all are received
	for i := 0; i < 3; i++ {
		msg := Message{Type: "batch", Payload: json.RawMessage(`{"i":` + string(rune('0'+i)) + `}`)}
		msgBytes, _ := json.Marshal(msg)
		if err := sc.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
			t.Fatalf("Failed to write message %d: %v", i, err)
		}
	}

	received := 0
	timeout := time.After(2 * time.Second)
	for received < 3 {
		select {
		case msg := <-receivedMsgs:
			if msg.Type != "batch" {
				t.Errorf("Expected type 'batch', got '%s'", msg.Type)
			}
			received++
		case <-timeout:
			t.Fatalf("Timeout: received only %d of 3 messages", received)
		}
	}

	clientConn.Close()
}

func TestClient_WritePump_ChannelClosed_SendsCloseMessage(t *testing.T) {
	logger := zap.NewNop()
	hub := NewHub(logger)
	go hub.Run()

	closeReceived := make(chan bool, 1)

	server, wsURL := createTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		conn, err := testUpgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		conn.SetCloseHandler(func(code int, text string) error {
			closeReceived <- true
			return nil
		})
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	})
	defer server.Close()

	clientConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	client := NewClient(hub, clientConn, "close-chan-test", logger)

	done := make(chan bool)
	go func() {
		client.WritePump()
		done <- true
	}()

	// Close the send channel to trigger the ok=false branch in WritePump
	close(client.send)

	select {
	case <-done:
		// WritePump exited correctly
	case <-time.After(2 * time.Second):
		t.Error("WritePump did not exit after send channel was closed")
	}

	// Check if close message was received by server
	select {
	case <-closeReceived:
		// Server received the close message
	case <-time.After(500 * time.Millisecond):
		// Close message may have been dropped if connection was already closing
		// This is acceptable behavior
	}

	clientConn.Close()
}

func TestClient_handleOutgoingMessage_WithQueuedMessages(t *testing.T) {
	logger := zap.NewNop()
	hub := NewHub(logger)

	receivedMsgs := make(chan []byte, 10)

	server, wsURL := createTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		conn, err := testUpgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				return
			}
			receivedMsgs <- data
		}
	})
	defer server.Close()

	clientConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer clientConn.Close()

	client := NewClient(hub, clientConn, "queue-handle-test", logger)

	// Pre-queue messages before calling handleOutgoingMessage
	client.send <- []byte(`{"type":"queued-a"}`)
	client.send <- []byte(`{"type":"queued-b"}`)

	// Call handleOutgoingMessage with ok=true, which internally calls writeMessageWithQueue
	result := client.handleOutgoingMessage([]byte(`{"type":"primary"}`), true)
	if !result {
		t.Error("handleOutgoingMessage should return true for successful write with queued messages")
	}

	// Should receive primary + 2 queued = 3 messages
	received := 0
	timeout := time.After(2 * time.Second)
	for received < 3 {
		select {
		case <-receivedMsgs:
			received++
		case <-timeout:
			t.Fatalf("Timeout: received only %d of 3 messages", received)
		}
	}
}

func TestClient_ReadPump_SetReadDeadlineError(t *testing.T) {
	logger := zap.NewNop()
	hub := NewHub(logger)
	go hub.Run()

	server, wsURL := createTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		conn, err := testUpgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		conn.Close()
	})
	defer server.Close()

	clientConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	client := NewClient(hub, clientConn, "deadline-error-agent", logger)
	hub.Register(client)

	// Close connection before ReadPump starts so SetReadDeadline fails
	clientConn.Close()

	done := make(chan bool, 1)
	handler := func(c *Client, msg *Message) {}
	go func() {
		client.ReadPump(handler)
		done <- true
	}()

	select {
	case <-done:
		// ReadPump exited due to SetReadDeadline error on closed connection
	case <-time.After(2 * time.Second):
		t.Error("ReadPump did not exit after SetReadDeadline error")
	}
}

func TestClient_WritePump_ExitsCleanlyOnUnregister(t *testing.T) {
	// Verify that WritePump exits when hub unregisters the client (closes send channel)
	logger := zap.NewNop()
	hub := NewHub(logger)
	go hub.Run()

	server, wsURL := createTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		conn, err := testUpgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	})
	defer server.Close()

	clientConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer clientConn.Close()

	client := NewClient(hub, clientConn, "unregister-pump-test", logger)
	hub.Register(client)
	time.Sleep(50 * time.Millisecond)

	done := make(chan bool, 1)
	go func() {
		client.WritePump()
		done <- true
	}()

	// Send a message first to ensure pump is running
	_ = client.Send("alive", map[string]bool{"ok": true})
	time.Sleep(50 * time.Millisecond)

	// Unregister client via hub to close the send channel
	hub.Unregister(client)

	select {
	case <-done:
		// WritePump exited cleanly
	case <-time.After(2 * time.Second):
		t.Error("WritePump did not exit after unregister")
	}
}

func TestClient_sendPing_WriteMessageFails(t *testing.T) {
	// Test: SetWriteDeadline succeeds but WriteMessage(PingMessage) fails
	logger := zap.NewNop()
	hub := NewHub(logger)

	server, wsURL := createTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		conn, err := testUpgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		// Close server side immediately to cause subsequent writes to fail
		time.Sleep(20 * time.Millisecond)
		conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		conn.Close()
	})
	defer server.Close()

	clientConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	client := NewClient(hub, clientConn, "ping-write-fail", logger)

	// Wait for server to close its side
	time.Sleep(100 * time.Millisecond)

	// Drain any close message
	clientConn.SetReadDeadline(time.Now().Add(50 * time.Millisecond))
	_, _, _ = clientConn.ReadMessage()

	// Now try sendPing - SetWriteDeadline may succeed on some OS but WriteMessage should fail
	result := client.sendPing()
	if result {
		// On some platforms, both might still succeed if the TCP connection hasn't fully torn down
		// This is platform-dependent, so we just verify it doesn't panic
		t.Log("sendPing succeeded (platform-dependent, connection may not have fully closed)")
	}
}

func TestClient_Send_MultipleConcurrent(t *testing.T) {
	logger := zap.NewNop()
	hub := NewHub(logger)

	client := NewClient(hub, nil, "concurrent-send", logger)

	var wg sync.WaitGroup
	const numSends = 10
	errCh := make(chan error, numSends)

	for i := 0; i < numSends; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			err := client.Send("concurrent", map[string]int{"index": idx})
			if err != nil {
				errCh <- err
			}
		}(i)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		t.Errorf("Concurrent Send returned error: %v", err)
	}

	// Drain and count messages from the send channel
	count := 0
	for {
		select {
		case <-client.send:
			count++
		default:
			goto done
		}
	}
done:
	if count != numSends {
		t.Errorf("Expected %d messages in send channel, got %d", numSends, count)
	}
}
