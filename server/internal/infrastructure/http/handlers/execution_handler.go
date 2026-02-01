package handlers

import (
	"net/http"

	"autostrike/internal/application"

	"github.com/gin-gonic/gin"
)

// ExecutionHandler handles execution-related HTTP requests
type ExecutionHandler struct {
	service *application.ExecutionService
}

// NewExecutionHandler creates a new execution handler
func NewExecutionHandler(service *application.ExecutionService) *ExecutionHandler {
	return &ExecutionHandler{service: service}
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
	}
}

// ListExecutions returns recent executions
func (h *ExecutionHandler) ListExecutions(c *gin.Context) {
	executions, err := h.service.GetRecentExecutions(c.Request.Context(), 50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
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

	c.JSON(http.StatusCreated, execution)
}

// CompleteExecution marks an execution as completed
func (h *ExecutionHandler) CompleteExecution(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.CompleteExecution(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "completed"})
}
