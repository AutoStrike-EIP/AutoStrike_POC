package handlers

import (
	"net/http"

	"autostrike/internal/application"
	"autostrike/internal/domain/entity"

	"github.com/gin-gonic/gin"
)

// AgentHandler handles agent-related HTTP requests
type AgentHandler struct {
	service *application.AgentService
}

// NewAgentHandler creates a new agent handler
func NewAgentHandler(service *application.AgentService) *AgentHandler {
	return &AgentHandler{service: service}
}

// RegisterRoutes registers agent routes
func (h *AgentHandler) RegisterRoutes(r *gin.RouterGroup) {
	agents := r.Group("/agents")
	{
		agents.GET("", h.ListAgents)
		agents.GET("/:paw", h.GetAgent)
		agents.POST("", h.RegisterAgent)
		agents.DELETE("/:paw", h.DeleteAgent)
		agents.POST("/:paw/heartbeat", h.Heartbeat)
	}
}

// ListAgents returns all agents
func (h *AgentHandler) ListAgents(c *gin.Context) {
	agents, err := h.service.GetAllAgents(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, agents)
}

// GetAgent returns a specific agent
func (h *AgentHandler) GetAgent(c *gin.Context) {
	paw := c.Param("paw")

	agent, err := h.service.GetAgent(c.Request.Context(), paw)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "agent not found"})
		return
	}

	c.JSON(http.StatusOK, agent)
}

// RegisterAgentRequest represents the request body for agent registration
type RegisterAgentRequest struct {
	Paw       string   `json:"paw" binding:"required"`
	Hostname  string   `json:"hostname" binding:"required"`
	Username  string   `json:"username" binding:"required"`
	Platform  string   `json:"platform" binding:"required"`
	Executors []string `json:"executors" binding:"required"`
}

// RegisterAgent registers a new agent
func (h *AgentHandler) RegisterAgent(c *gin.Context) {
	var req RegisterAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	agent := &entity.Agent{
		Paw:       req.Paw,
		Hostname:  req.Hostname,
		Username:  req.Username,
		Platform:  req.Platform,
		Executors: req.Executors,
	}

	if err := h.service.RegisterAgent(c.Request.Context(), agent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, agent)
}

// DeleteAgent deletes an agent
func (h *AgentHandler) DeleteAgent(c *gin.Context) {
	paw := c.Param("paw")

	if err := h.service.DeleteAgent(c.Request.Context(), paw); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// Heartbeat updates agent's last seen
func (h *AgentHandler) Heartbeat(c *gin.Context) {
	paw := c.Param("paw")

	if err := h.service.Heartbeat(c.Request.Context(), paw); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
