package handlers

import (
	"errors"
	"net/http"
	"path/filepath"
	"strings"

	"autostrike/internal/application"
	"autostrike/internal/domain/entity"

	"github.com/gin-gonic/gin"
)

// allowedImportDirs is the whitelist of directories allowed for technique imports
var allowedImportDirs = []string{"./configs", "configs", "./config", "config"}

// TechniqueHandler handles technique-related HTTP requests
type TechniqueHandler struct {
	service *application.TechniqueService
}

// NewTechniqueHandler creates a new technique handler
func NewTechniqueHandler(service *application.TechniqueService) *TechniqueHandler {
	return &TechniqueHandler{service: service}
}

// RegisterRoutes registers technique routes
func (h *TechniqueHandler) RegisterRoutes(r *gin.RouterGroup) {
	techniques := r.Group("/techniques")
	{
		techniques.GET("", h.ListTechniques)
		techniques.GET("/:id", h.GetTechnique)
		techniques.GET("/tactic/:tactic", h.GetByTactic)
		techniques.GET("/platform/:platform", h.GetByPlatform)
		techniques.GET("/coverage", h.GetCoverage)
		techniques.POST("/import", h.ImportTechniques)
		techniques.POST("/import/json", h.ImportTechniquesJSON)
	}
}

// ListTechniques returns all techniques
func (h *TechniqueHandler) ListTechniques(c *gin.Context) {
	techniques, err := h.service.GetAllTechniques(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return empty array instead of null
	if techniques == nil {
		techniques = []*entity.Technique{}
	}
	c.JSON(http.StatusOK, techniques)
}

// GetTechnique returns a specific technique
func (h *TechniqueHandler) GetTechnique(c *gin.Context) {
	id := c.Param("id")

	technique, err := h.service.GetTechnique(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "technique not found"})
		return
	}

	c.JSON(http.StatusOK, technique)
}

// GetByTactic returns techniques by MITRE tactic
func (h *TechniqueHandler) GetByTactic(c *gin.Context) {
	tactic := entity.TacticType(c.Param("tactic"))

	techniques, err := h.service.GetTechniquesByTactic(c.Request.Context(), tactic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return empty array instead of null
	if techniques == nil {
		techniques = []*entity.Technique{}
	}
	c.JSON(http.StatusOK, techniques)
}

// GetByPlatform returns techniques by platform
func (h *TechniqueHandler) GetByPlatform(c *gin.Context) {
	platform := c.Param("platform")

	techniques, err := h.service.GetTechniquesByPlatform(c.Request.Context(), platform)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return empty array instead of null
	if techniques == nil {
		techniques = []*entity.Technique{}
	}
	c.JSON(http.StatusOK, techniques)
}

// GetCoverage returns MITRE ATT&CK coverage statistics
func (h *TechniqueHandler) GetCoverage(c *gin.Context) {
	coverage, err := h.service.GetCoverage(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, coverage)
}

// filterExecutorsByPlatform returns executors matching the given platform,
// or empty if the technique does not support that platform.
func filterExecutorsByPlatform(technique *entity.Technique, platform string) []entity.Executor {
	platformSupported := false
	for _, p := range technique.Platforms {
		if p == platform {
			platformSupported = true
			break
		}
	}
	if !platformSupported {
		return []entity.Executor{}
	}

	var executors []entity.Executor
	for _, exec := range technique.Executors {
		if exec.Platform == "" || exec.Platform == platform {
			executors = append(executors, exec)
		}
	}
	return executors
}

// GetExecutors returns executors for a technique, optionally filtered by platform
func (h *TechniqueHandler) GetExecutors(c *gin.Context) {
	id := c.Param("id")

	technique, err := h.service.GetTechnique(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "technique not found"})
		return
	}

	platform := c.Query("platform")
	var executors []entity.Executor

	if platform != "" {
		executors = filterExecutorsByPlatform(technique, platform)
	} else {
		executors = technique.Executors
	}

	if executors == nil {
		executors = []entity.Executor{}
	}
	c.JSON(http.StatusOK, executors)
}

// ImportRequest represents the request body for importing techniques
type ImportRequest struct {
	Path string `json:"path" binding:"required"`
}

// validateImportPath checks that the path is within allowed directories to prevent path traversal.
// Uses filepath.Abs and filepath.Rel for reliable containment checks instead of string matching.
func validateImportPath(requestPath string) error {
	// Resolve the request path to an absolute path
	absPath, err := filepath.Abs(requestPath)
	if err != nil {
		return errors.New("invalid path")
	}

	// Check against each allowed directory using filepath.Rel
	for _, dir := range allowedImportDirs {
		absDir, err := filepath.Abs(dir)
		if err != nil {
			continue
		}
		// filepath.Rel computes a relative path from absDir to absPath.
		// If the result starts with "..", absPath is outside absDir.
		rel, err := filepath.Rel(absDir, absPath)
		if err != nil {
			continue
		}
		if rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return nil
		}
	}

	return errors.New("import path must be within the configs directory")
}

// ImportTechniques imports techniques from YAML file
func (h *TechniqueHandler) ImportTechniques(c *gin.Context) {
	var req ImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate path to prevent path traversal attacks
	if err := validateImportPath(req.Path); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ImportTechniques(c.Request.Context(), req.Path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "imported"})
}

// ImportJSONRequest represents the request body for importing techniques from JSON
type ImportJSONRequest struct {
	Techniques []*entity.Technique `json:"techniques" binding:"required"`
}

// ImportJSONResponse represents the response for JSON import
type ImportJSONResponse struct {
	Imported int      `json:"imported"`
	Failed   int      `json:"failed"`
	Errors   []string `json:"errors,omitempty"`
}

// ImportTechniquesJSON imports techniques directly from JSON request body
func (h *TechniqueHandler) ImportTechniquesJSON(c *gin.Context) {
	var req ImportJSONRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	imported := 0
	failed := 0
	var errors []string

	for _, t := range req.Techniques {
		if err := h.service.CreateTechnique(c.Request.Context(), t); err != nil {
			// Try update if create fails (technique already exists)
			if err := h.service.UpdateTechnique(c.Request.Context(), t); err != nil {
				failed++
				errors = append(errors, "Failed to import "+t.ID+": "+err.Error())
			} else {
				imported++
			}
		} else {
			imported++
		}
	}

	c.JSON(http.StatusOK, ImportJSONResponse{
		Imported: imported,
		Failed:   failed,
		Errors:   errors,
	})
}
