package rest

import (
	"os"

	"autostrike/internal/application"
	"autostrike/internal/infrastructure/http/handlers"
	"autostrike/internal/infrastructure/http/middleware"
	"autostrike/internal/infrastructure/websocket"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Server represents the HTTP REST server
type Server struct {
	router *gin.Engine
	logger *zap.Logger
}

// ServerConfig contains server configuration options
type ServerConfig struct {
	JWTSecret      string
	AgentSecret    string
	EnableAuth     bool
}

// NewServerConfig creates a server config from environment variables
func NewServerConfig() *ServerConfig {
	enableAuth := os.Getenv("ENABLE_AUTH") != "false" // Enabled by default
	return &ServerConfig{
		JWTSecret:   os.Getenv("JWT_SECRET"),
		AgentSecret: os.Getenv("AGENT_SECRET"),
		EnableAuth:  enableAuth,
	}
}

// NewServer creates a new REST server with all routes configured
func NewServer(
	agentService *application.AgentService,
	scenarioService *application.ScenarioService,
	executionService *application.ExecutionService,
	techniqueService *application.TechniqueService,
	hub *websocket.Hub,
	logger *zap.Logger,
) *Server {
	return NewServerWithConfig(agentService, scenarioService, executionService, techniqueService, hub, logger, NewServerConfig())
}

// NewServerWithConfig creates a new REST server with explicit configuration
func NewServerWithConfig(
	agentService *application.AgentService,
	scenarioService *application.ScenarioService,
	executionService *application.ExecutionService,
	techniqueService *application.TechniqueService,
	hub *websocket.Hub,
	logger *zap.Logger,
	config *ServerConfig,
) *Server {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Global middleware
	router.Use(middleware.LoggingMiddleware(logger))
	router.Use(middleware.RecoveryMiddleware(logger))

	// Health check (always public)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// WebSocket routes (uses agent auth)
	if hub != nil {
		wsHandler := handlers.NewWebSocketHandler(hub, agentService, logger)
		wsHandler.RegisterRoutes(router)
	}

	// API v1 routes
	api := router.Group("/api/v1")

	// Apply authentication middleware if enabled
	if config.EnableAuth && config.JWTSecret != "" {
		authConfig := &middleware.AuthConfig{
			JWTSecret:   config.JWTSecret,
			AgentSecret: config.AgentSecret,
		}
		api.Use(middleware.AuthMiddleware(authConfig))
		logger.Info("Authentication middleware enabled for API routes")
	} else {
		logger.Warn("Authentication middleware DISABLED - set ENABLE_AUTH=true and JWT_SECRET in production")
	}

	{
		// Agents
		agentHandler := handlers.NewAgentHandler(agentService)
		agentHandler.RegisterRoutes(api)

		// Techniques
		techniqueHandler := handlers.NewTechniqueHandler(techniqueService)
		techniqueHandler.RegisterRoutes(api)

		// Executions
		executionHandler := handlers.NewExecutionHandler(executionService)
		executionHandler.RegisterRoutes(api)
	}

	return &Server{
		router: router,
		logger: logger,
	}
}

// Run starts the HTTP server
func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}

// Router returns the underlying gin router for testing
func (s *Server) Router() *gin.Engine {
	return s.router
}
