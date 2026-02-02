package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"autostrike/internal/application"
	"autostrike/internal/domain/entity"
	"autostrike/internal/infrastructure/websocket"

	"github.com/gin-gonic/gin"
)

// ExecutionHandler handles execution-related HTTP requests
type ExecutionHandler struct {
	service *application.ExecutionService
	hub     *websocket.Hub
}

// NewExecutionHandler creates a new execution handler
func NewExecutionHandler(service *application.ExecutionService) *ExecutionHandler {
	return &ExecutionHandler{service: service}
}

// NewExecutionHandlerWithHub creates a new execution handler with WebSocket support
func NewExecutionHandlerWithHub(service *application.ExecutionService, hub *websocket.Hub) *ExecutionHandler {
	return &ExecutionHandler{service: service, hub: hub}
}

// broadcastExecutionEvent sends an execution event to all connected clients
func (h *ExecutionHandler) broadcastExecutionEvent(eventType string, executionID string, data interface{}) {
	if h.hub == nil {
		return
	}

	msg := map[string]interface{}{
		"type": eventType,
		"payload": map[string]interface{}{
			"execution_id": executionID,
			"data":         data,
		},
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return
	}

	h.hub.Broadcast(msgBytes)
}

// RegisterRoutes registers execution routes
func (h *ExecutionHandler) RegisterRoutes(r *gin.RouterGroup) {
	executions := r.Group("/executions")
	{
		executions.GET("", h.ListExecutions)
		executions.GET("/:id", h.GetExecution)
		executions.GET("/:id/results", h.GetResults)
		executions.POST("", h.StartExecution)
		executions.POST("/:id/complete", h.CompleteExecution)
		executions.POST("/:id/stop", h.StopExecution)
	}
}

// ListExecutions returns recent executions
func (h *ExecutionHandler) ListExecutions(c *gin.Context) {
	executions, err := h.service.GetRecentExecutions(c.Request.Context(), 50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return empty array instead of null
	if executions == nil {
		executions = []*entity.Execution{}
	}
	c.JSON(http.StatusOK, executions)
}

// GetExecution returns a specific execution
func (h *ExecutionHandler) GetExecution(c *gin.Context) {
	id := c.Param("id")

	execution, err := h.service.GetExecution(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "execution not found"})
		return
	}

	c.JSON(http.StatusOK, execution)
}

// GetResults returns results for an execution
func (h *ExecutionHandler) GetResults(c *gin.Context) {
	id := c.Param("id")

	results, err := h.service.GetExecutionResults(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return empty array instead of null
	if results == nil {
		results = []*entity.ExecutionResult{}
	}
	c.JSON(http.StatusOK, results)
}

// StartExecutionRequest represents the request body for starting an execution
type StartExecutionRequest struct {
	ScenarioID string   `json:"scenario_id" binding:"required"`
	AgentPaws  []string `json:"agent_paws" binding:"required"`
	SafeMode   bool     `json:"safe_mode"`
}

// StartExecution starts a new scenario execution
func (h *ExecutionHandler) StartExecution(c *gin.Context) {
	var req StartExecutionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	execution, err := h.service.StartExecution(c.Request.Context(), req.ScenarioID, req.AgentPaws, req.SafeMode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Broadcast execution started event to all connected clients
	h.broadcastExecutionEvent("execution_started", execution.ID, execution)

	c.JSON(http.StatusCreated, execution)
}

// CompleteExecution marks an execution as completed
func (h *ExecutionHandler) CompleteExecution(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.CompleteExecution(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Broadcast execution completed event to all connected clients
	h.broadcastExecutionEvent("execution_completed", id, map[string]string{
		"status": "completed",
	})

	c.JSON(http.StatusOK, gin.H{"status": "completed"})
}

// StopExecution cancels a running execution
func (h *ExecutionHandler) StopExecution(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.CancelExecution(c.Request.Context(), id); err != nil {
		errMsg := err.Error()
		// Return appropriate HTTP status based on error type
		if strings.Contains(errMsg, "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(errMsg, "cannot be cancelled") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Broadcast cancellation event to all connected clients
	h.broadcastExecutionEvent("execution_cancelled", id, map[string]string{
		"status": "cancelled",
	})

	c.JSON(http.StatusOK, gin.H{"status": "cancelled"})
}
