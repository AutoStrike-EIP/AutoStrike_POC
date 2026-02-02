package handlers

import (
	"net/http"

	"autostrike/internal/application"
	"autostrike/internal/domain/entity"

	"github.com/gin-gonic/gin"
)

// ScenarioHandler handles scenario-related HTTP requests
type ScenarioHandler struct {
	service *application.ScenarioService
}

// NewScenarioHandler creates a new scenario handler
func NewScenarioHandler(service *application.ScenarioService) *ScenarioHandler {
	return &ScenarioHandler{service: service}
}

// RegisterRoutes registers scenario routes
func (h *ScenarioHandler) RegisterRoutes(r *gin.RouterGroup) {
	scenarios := r.Group("/scenarios")
	{
		scenarios.GET("", h.ListScenarios)
		scenarios.GET("/tag/:tag", h.GetScenariosByTag) // Must be before /:id
		scenarios.GET("/:id", h.GetScenario)
		scenarios.POST("", h.CreateScenario)
		scenarios.PUT("/:id", h.UpdateScenario)
		scenarios.DELETE("/:id", h.DeleteScenario)
	}
}

// ListScenarios returns all scenarios
func (h *ScenarioHandler) ListScenarios(c *gin.Context) {
	scenarios, err := h.service.GetAllScenarios(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return empty array instead of null
	if scenarios == nil {
		scenarios = []*entity.Scenario{}
	}
	c.JSON(http.StatusOK, scenarios)
}

// GetScenario returns a specific scenario
func (h *ScenarioHandler) GetScenario(c *gin.Context) {
	id := c.Param("id")

	scenario, err := h.service.GetScenario(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "scenario not found"})
		return
	}

	c.JSON(http.StatusOK, scenario)
}

// GetScenariosByTag returns scenarios filtered by tag
func (h *ScenarioHandler) GetScenariosByTag(c *gin.Context) {
	tag := c.Param("tag")

	scenarios, err := h.service.GetScenariosByTag(c.Request.Context(), tag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return empty array instead of null
	if scenarios == nil {
		scenarios = []*entity.Scenario{}
	}
	c.JSON(http.StatusOK, scenarios)
}

// CreateScenarioRequest represents the request body for scenario creation
type CreateScenarioRequest struct {
	Name        string         `json:"name" binding:"required"`
	Description string         `json:"description"`
	Phases      []entity.Phase `json:"phases" binding:"required"`
	Tags        []string       `json:"tags,omitempty"`
}

// CreateScenario creates a new scenario
func (h *ScenarioHandler) CreateScenario(c *gin.Context) {
	var req CreateScenarioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	scenario := &entity.Scenario{
		Name:        req.Name,
		Description: req.Description,
		Phases:      req.Phases,
		Tags:        req.Tags,
	}

	if err := h.service.CreateScenario(c.Request.Context(), scenario); err != nil {
		// Check if it's a validation error
		if _, ok := err.(*application.ValidationError); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, scenario)
}

// UpdateScenarioRequest represents the request body for scenario update
type UpdateScenarioRequest struct {
	Name        string         `json:"name" binding:"required"`
	Description string         `json:"description"`
	Phases      []entity.Phase `json:"phases" binding:"required"`
	Tags        []string       `json:"tags,omitempty"`
}

// UpdateScenario updates an existing scenario
func (h *ScenarioHandler) UpdateScenario(c *gin.Context) {
	id := c.Param("id")

	// Check if scenario exists
	existing, err := h.service.GetScenario(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "scenario not found"})
		return
	}

	var req UpdateScenarioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	scenario := &entity.Scenario{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Phases:      req.Phases,
		Tags:        req.Tags,
		CreatedAt:   existing.CreatedAt,
	}

	if err := h.service.UpdateScenario(c.Request.Context(), scenario); err != nil {
		// Check if it's a validation error
		if _, ok := err.(*application.ValidationError); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, scenario)
}

// DeleteScenario deletes a scenario
func (h *ScenarioHandler) DeleteScenario(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteScenario(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
