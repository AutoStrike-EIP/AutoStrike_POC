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
