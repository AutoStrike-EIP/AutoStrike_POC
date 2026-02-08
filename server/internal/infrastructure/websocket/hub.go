package websocket

import (
	"sync"

	"go.uber.org/zap"
)

// AgentDisconnectCallback is called when an agent disconnects
type AgentDisconnectCallback func(paw string)

// Hub manages WebSocket connections
type Hub struct {
	clients           map[*Client]bool
	agents            map[string]*Client
	broadcast         chan []byte
	register          chan *Client
	unregister        chan *Client
	mu                sync.RWMutex
	logger            *zap.Logger
	onAgentDisconnect AgentDisconnectCallback
}

// NewHub creates a new WebSocket hub
func NewHub(logger *zap.Logger) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		agents:     make(map[string]*Client),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		logger:     logger,
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.handleRegister(client)

		case client := <-h.unregister:
			h.handleUnregister(client)

		case message := <-h.broadcast:
			h.handleBroadcast(message)
		}
	}
}

// handleRegister handles client registration with proper mutex handling
func (h *Hub) handleRegister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[client] = true
	paw := client.GetAgentPaw()
	if paw != "" {
		h.agents[paw] = client
		h.logger.Info("Agent connected", zap.String("paw", paw))
	}
}

// handleUnregister handles client unregistration with proper mutex handling
func (h *Hub) handleUnregister(client *Client) {
	h.mu.Lock()
	paw := ""
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		paw = client.GetAgentPaw()
		if paw != "" {
			// Only remove from agents map if THIS client is the one registered.
			// A newer connection with the same paw may have replaced it.
			if h.agents[paw] == client {
				delete(h.agents, paw)
				h.logger.Info("Agent disconnected", zap.String("paw", paw))
			} else {
				h.logger.Info("Stale agent connection closed", zap.String("paw", paw))
				paw = "" // Don't fire disconnect callback for stale connection
			}
		}
		close(client.send)
	}
	h.mu.Unlock()

	// Call disconnect callback outside of lock to avoid deadlock
	if paw != "" && h.onAgentDisconnect != nil {
		h.onAgentDisconnect(paw)
	}
}

// handleBroadcast handles message broadcasting with proper mutex handling
func (h *Hub) handleBroadcast(message []byte) {
	h.mu.Lock()
	for client := range h.clients {
		select {
		case client.send <- message:
		default:
			// Channel full — skip this client instead of disconnecting it.
			// A temporarily full buffer is not a fatal condition.
			h.logger.Warn("Broadcast: send channel full, skipping client",
				zap.String("paw", client.GetAgentPaw()),
			)
		}
	}
	h.mu.Unlock()
}

// SendToAgent sends a message to a specific agent
func (h *Hub) SendToAgent(paw string, message []byte) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if client, ok := h.agents[paw]; ok {
		select {
		case client.send <- message:
			return true
		default:
			return false
		}
	}
	return false
}

// Broadcast sends a message to all connected clients
func (h *Hub) Broadcast(message []byte) {
	h.broadcast <- message
}

// GetConnectedAgents returns list of connected agent paws
func (h *Hub) GetConnectedAgents() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	paws := make([]string, 0, len(h.agents))
	for paw := range h.agents {
		paws = append(paws, paw)
	}
	return paws
}

// IsAgentConnected checks if an agent is connected
func (h *Hub) IsAgentConnected(paw string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	_, ok := h.agents[paw]
	return ok
}

// RegisterAgent registers an agent with its paw after initial connection
func (h *Hub) RegisterAgent(paw string, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.agents[paw] = client
	h.logger.Info("Agent registered with paw", zap.String("paw", paw))
}

// Register adds a client to the hub
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Unregister removes a client from the hub
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

// SetOnAgentDisconnect sets the callback for agent disconnection
func (h *Hub) SetOnAgentDisconnect(callback AgentDisconnectCallback) {
	h.onAgentDisconnect = callback
}
