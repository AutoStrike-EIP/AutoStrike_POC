package rest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"autostrike/internal/application"
	"autostrike/internal/domain/entity"
	"autostrike/internal/domain/service"
	"autostrike/internal/infrastructure/websocket"

	"go.uber.org/zap"
)

// Mock repository implementations for testing
type mockAgentRepo struct{}

func (m *mockAgentRepo) Create(ctx context.Context, agent *entity.Agent) error       { return nil }
func (m *mockAgentRepo) Update(ctx context.Context, agent *entity.Agent) error       { return nil }
func (m *mockAgentRepo) Delete(ctx context.Context, paw string) error                { return nil }
func (m *mockAgentRepo) FindByPaw(ctx context.Context, paw string) (*entity.Agent, error) {
	return &entity.Agent{Paw: paw}, nil
}
func (m *mockAgentRepo) FindByPaws(ctx context.Context, paws []string) ([]*entity.Agent, error) {
	return []*entity.Agent{}, nil
}
func (m *mockAgentRepo) FindAll(ctx context.Context) ([]*entity.Agent, error) {
	return []*entity.Agent{}, nil
}
func (m *mockAgentRepo) FindByStatus(ctx context.Context, status entity.AgentStatus) ([]*entity.Agent, error) {
	return []*entity.Agent{}, nil
}
func (m *mockAgentRepo) FindByPlatform(ctx context.Context, platform string) ([]*entity.Agent, error) {
	return []*entity.Agent{}, nil
}
func (m *mockAgentRepo) UpdateLastSeen(ctx context.Context, paw string) error { return nil }

type mockScenarioRepo struct{}

func (m *mockScenarioRepo) Create(ctx context.Context, scenario *entity.Scenario) error { return nil }
func (m *mockScenarioRepo) Update(ctx context.Context, scenario *entity.Scenario) error { return nil }
func (m *mockScenarioRepo) Delete(ctx context.Context, id string) error                 { return nil }
func (m *mockScenarioRepo) FindByID(ctx context.Context, id string) (*entity.Scenario, error) {
	return &entity.Scenario{ID: id}, nil
}
func (m *mockScenarioRepo) FindAll(ctx context.Context) ([]*entity.Scenario, error) {
	return []*entity.Scenario{}, nil
}
func (m *mockScenarioRepo) FindByTag(ctx context.Context, tag string) ([]*entity.Scenario, error) {
	return []*entity.Scenario{}, nil
}

type mockTechniqueRepo struct{}

func (m *mockTechniqueRepo) Create(ctx context.Context, technique *entity.Technique) error {
	return nil
}
func (m *mockTechniqueRepo) Update(ctx context.Context, technique *entity.Technique) error {
	return nil
}
func (m *mockTechniqueRepo) Delete(ctx context.Context, id string) error { return nil }
func (m *mockTechniqueRepo) FindByID(ctx context.Context, id string) (*entity.Technique, error) {
	return &entity.Technique{ID: id}, nil
}
func (m *mockTechniqueRepo) FindAll(ctx context.Context) ([]*entity.Technique, error) {
	return []*entity.Technique{}, nil
}
func (m *mockTechniqueRepo) FindByTactic(ctx context.Context, tactic entity.TacticType) ([]*entity.Technique, error) {
	return []*entity.Technique{}, nil
}
func (m *mockTechniqueRepo) FindByPlatform(ctx context.Context, platform string) ([]*entity.Technique, error) {
	return []*entity.Technique{}, nil
}
func (m *mockTechniqueRepo) ImportFromYAML(ctx context.Context, path string) error { return nil }

type mockResultRepo struct{}

func (m *mockResultRepo) CreateExecution(ctx context.Context, execution *entity.Execution) error {
	return nil
}
func (m *mockResultRepo) UpdateExecution(ctx context.Context, execution *entity.Execution) error {
	return nil
}
func (m *mockResultRepo) FindExecutionByID(ctx context.Context, id string) (*entity.Execution, error) {
	return &entity.Execution{ID: id}, nil
}
func (m *mockResultRepo) FindExecutionsByScenario(ctx context.Context, scenarioID string) ([]*entity.Execution, error) {
	return []*entity.Execution{}, nil
}
func (m *mockResultRepo) FindRecentExecutions(ctx context.Context, limit int) ([]*entity.Execution, error) {
	return []*entity.Execution{}, nil
}
func (m *mockResultRepo) CreateResult(ctx context.Context, result *entity.ExecutionResult) error {
	return nil
}
func (m *mockResultRepo) UpdateResult(ctx context.Context, result *entity.ExecutionResult) error {
	return nil
}
func (m *mockResultRepo) FindResultsByExecution(ctx context.Context, executionID string) ([]*entity.ExecutionResult, error) {
	return []*entity.ExecutionResult{}, nil
}
func (m *mockResultRepo) FindResultsByTechnique(ctx context.Context, techniqueID string) ([]*entity.ExecutionResult, error) {
	return []*entity.ExecutionResult{}, nil
}

func TestNewServerConfig_Default(t *testing.T) {
	// Clear environment variables
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("AGENT_SECRET")
	os.Unsetenv("ENABLE_AUTH")

	config := NewServerConfig()

	if config.JWTSecret != "" {
		t.Errorf("Expected empty JWT secret, got '%s'", config.JWTSecret)
	}
	if config.AgentSecret != "" {
		t.Errorf("Expected empty agent secret, got '%s'", config.AgentSecret)
	}
	if !config.EnableAuth {
		t.Error("Expected EnableAuth to be true by default")
	}
}

func TestNewServerConfig_WithEnv(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-jwt-secret")
	os.Setenv("AGENT_SECRET", "test-agent-secret")
	os.Setenv("ENABLE_AUTH", "true")
	defer func() {
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("AGENT_SECRET")
		os.Unsetenv("ENABLE_AUTH")
	}()

	config := NewServerConfig()

	if config.JWTSecret != "test-jwt-secret" {
		t.Errorf("Expected JWT secret 'test-jwt-secret', got '%s'", config.JWTSecret)
	}
	if config.AgentSecret != "test-agent-secret" {
		t.Errorf("Expected agent secret 'test-agent-secret', got '%s'", config.AgentSecret)
	}
	if !config.EnableAuth {
		t.Error("Expected EnableAuth to be true")
	}
}

func TestNewServerConfig_AuthDisabled(t *testing.T) {
	os.Setenv("ENABLE_AUTH", "false")
	defer os.Unsetenv("ENABLE_AUTH")

	config := NewServerConfig()

	if config.EnableAuth {
		t.Error("Expected EnableAuth to be false when ENABLE_AUTH=false")
	}
}

func TestNewServerWithConfig_AuthDisabled(t *testing.T) {
	logger := zap.NewNop()
	agentService := application.NewAgentService(&mockAgentRepo{})
	scenarioService := application.NewScenarioService(&mockScenarioRepo{}, &mockTechniqueRepo{}, service.NewTechniqueValidator())
	techniqueService := application.NewTechniqueService(&mockTechniqueRepo{})
	executionService := application.NewExecutionService(&mockResultRepo{}, &mockScenarioRepo{}, &mockTechniqueRepo{}, &mockAgentRepo{}, nil, nil)
	hub := websocket.NewHub(logger)

	config := &ServerConfig{
		JWTSecret:   "",
		AgentSecret: "",
		EnableAuth:  false,
	}

	server := NewServerWithConfig(agentService, scenarioService, executionService, techniqueService, hub, logger, config)

	if server == nil {
		t.Fatal("NewServerWithConfig returned nil")
	}

	if server.router == nil {
		t.Error("Server router is nil")
	}
}

func TestNewServerWithConfig_AuthEnabled(t *testing.T) {
	logger := zap.NewNop()
	agentService := application.NewAgentService(&mockAgentRepo{})
	scenarioService := application.NewScenarioService(&mockScenarioRepo{}, &mockTechniqueRepo{}, service.NewTechniqueValidator())
	techniqueService := application.NewTechniqueService(&mockTechniqueRepo{})
	executionService := application.NewExecutionService(&mockResultRepo{}, &mockScenarioRepo{}, &mockTechniqueRepo{}, &mockAgentRepo{}, nil, nil)
	hub := websocket.NewHub(logger)

	config := &ServerConfig{
		JWTSecret:   "test-secret",
		AgentSecret: "agent-secret",
		EnableAuth:  true,
	}

	server := NewServerWithConfig(agentService, scenarioService, executionService, techniqueService, hub, logger, config)

	if server == nil {
		t.Fatal("NewServerWithConfig returned nil")
	}
}

func TestNewServerWithConfig_NoHub(t *testing.T) {
	logger := zap.NewNop()
	agentService := application.NewAgentService(&mockAgentRepo{})
	scenarioService := application.NewScenarioService(&mockScenarioRepo{}, &mockTechniqueRepo{}, service.NewTechniqueValidator())
	techniqueService := application.NewTechniqueService(&mockTechniqueRepo{})
	executionService := application.NewExecutionService(&mockResultRepo{}, &mockScenarioRepo{}, &mockTechniqueRepo{}, &mockAgentRepo{}, nil, nil)

	config := &ServerConfig{EnableAuth: false}

	// Should not panic with nil hub
	server := NewServerWithConfig(agentService, scenarioService, executionService, techniqueService, nil, logger, config)

	if server == nil {
		t.Fatal("NewServerWithConfig returned nil")
	}
}

func TestNewServer(t *testing.T) {
	os.Setenv("ENABLE_AUTH", "false")
	defer os.Unsetenv("ENABLE_AUTH")

	logger := zap.NewNop()
	agentService := application.NewAgentService(&mockAgentRepo{})
	scenarioService := application.NewScenarioService(&mockScenarioRepo{}, &mockTechniqueRepo{}, service.NewTechniqueValidator())
	techniqueService := application.NewTechniqueService(&mockTechniqueRepo{})
	executionService := application.NewExecutionService(&mockResultRepo{}, &mockScenarioRepo{}, &mockTechniqueRepo{}, &mockAgentRepo{}, nil, nil)
	hub := websocket.NewHub(logger)

	server := NewServer(agentService, scenarioService, executionService, techniqueService, hub, logger)

	if server == nil {
		t.Fatal("NewServer returned nil")
	}
}

func TestServer_Router(t *testing.T) {
	logger := zap.NewNop()
	agentService := application.NewAgentService(&mockAgentRepo{})
	scenarioService := application.NewScenarioService(&mockScenarioRepo{}, &mockTechniqueRepo{}, service.NewTechniqueValidator())
	techniqueService := application.NewTechniqueService(&mockTechniqueRepo{})
	executionService := application.NewExecutionService(&mockResultRepo{}, &mockScenarioRepo{}, &mockTechniqueRepo{}, &mockAgentRepo{}, nil, nil)

	config := &ServerConfig{EnableAuth: false}
	server := NewServerWithConfig(agentService, scenarioService, executionService, techniqueService, nil, logger, config)

	router := server.Router()
	if router == nil {
		t.Error("Router() returned nil")
	}
}

func TestServer_HealthEndpoint(t *testing.T) {
	logger := zap.NewNop()
	agentService := application.NewAgentService(&mockAgentRepo{})
	scenarioService := application.NewScenarioService(&mockScenarioRepo{}, &mockTechniqueRepo{}, service.NewTechniqueValidator())
	techniqueService := application.NewTechniqueService(&mockTechniqueRepo{})
	executionService := application.NewExecutionService(&mockResultRepo{}, &mockScenarioRepo{}, &mockTechniqueRepo{}, &mockAgentRepo{}, nil, nil)

	config := &ServerConfig{EnableAuth: false}
	server := NewServerWithConfig(agentService, scenarioService, executionService, techniqueService, nil, logger, config)

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestServer_APIRoutes(t *testing.T) {
	logger := zap.NewNop()
	agentService := application.NewAgentService(&mockAgentRepo{})
	scenarioService := application.NewScenarioService(&mockScenarioRepo{}, &mockTechniqueRepo{}, service.NewTechniqueValidator())
	techniqueService := application.NewTechniqueService(&mockTechniqueRepo{})
	executionService := application.NewExecutionService(&mockResultRepo{}, &mockScenarioRepo{}, &mockTechniqueRepo{}, &mockAgentRepo{}, nil, nil)

	config := &ServerConfig{EnableAuth: false}
	server := NewServerWithConfig(agentService, scenarioService, executionService, techniqueService, nil, logger, config)

	routes := server.Router().Routes()
	expectedPaths := []string{"/health", "/api/v1/agents", "/api/v1/techniques"}

	for _, expected := range expectedPaths {
		found := false
		for _, route := range routes {
			if route.Path == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected route %s not found", expected)
		}
	}
}

func TestServer_Run_InvalidAddress(t *testing.T) {
	logger := zap.NewNop()
	agentService := application.NewAgentService(&mockAgentRepo{})
	scenarioService := application.NewScenarioService(&mockScenarioRepo{}, &mockTechniqueRepo{}, service.NewTechniqueValidator())
	techniqueService := application.NewTechniqueService(&mockTechniqueRepo{})
	executionService := application.NewExecutionService(&mockResultRepo{}, &mockScenarioRepo{}, &mockTechniqueRepo{}, &mockAgentRepo{}, nil, nil)

	config := &ServerConfig{EnableAuth: false}
	server := NewServerWithConfig(agentService, scenarioService, executionService, techniqueService, nil, logger, config)

	// Run in goroutine and expect it to fail with invalid address
	errCh := make(chan error, 1)
	go func() {
		errCh <- server.Run("invalid-address:99999")
	}()

	select {
	case err := <-errCh:
		if err == nil {
			t.Error("Expected error with invalid address")
		}
	case <-time.After(500 * time.Millisecond):
		// Server might be trying to bind, which is also acceptable
	}
}
