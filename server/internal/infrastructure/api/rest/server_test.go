package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
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
func (m *mockScenarioRepo) ImportFromYAML(ctx context.Context, path string) error { return nil }

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
func (m *mockResultRepo) FindResultByID(ctx context.Context, id string) (*entity.ExecutionResult, error) {
	return &entity.ExecutionResult{ID: id}, nil
}
func (m *mockResultRepo) FindResultsByExecution(ctx context.Context, executionID string) ([]*entity.ExecutionResult, error) {
	return []*entity.ExecutionResult{}, nil
}
func (m *mockResultRepo) FindResultsByTechnique(ctx context.Context, techniqueID string) ([]*entity.ExecutionResult, error) {
	return []*entity.ExecutionResult{}, nil
}
func (m *mockResultRepo) FindExecutionsByDateRange(ctx context.Context, start, end time.Time) ([]*entity.Execution, error) {
	return []*entity.Execution{}, nil
}
func (m *mockResultRepo) FindCompletedExecutionsByDateRange(ctx context.Context, start, end time.Time) ([]*entity.Execution, error) {
	return []*entity.Execution{}, nil
}

func TestNewServerConfig_Default(t *testing.T) {
	// Clear environment variables
	_ = os.Unsetenv("JWT_SECRET")
	_ = os.Unsetenv("AGENT_SECRET")
	_ = os.Unsetenv("ENABLE_AUTH")

	config := NewServerConfig()

	if config.JWTSecret != "" {
		t.Errorf("Expected empty JWT secret, got '%s'", config.JWTSecret)
	}
	if config.AgentSecret != "" {
		t.Errorf("Expected empty agent secret, got '%s'", config.AgentSecret)
	}
	// Auth is disabled by default when JWT_SECRET is not set (dev mode)
	if config.EnableAuth {
		t.Error("Expected EnableAuth to be false when JWT_SECRET is not set")
	}
}

func TestNewServerConfig_EnabledWithSecret(t *testing.T) {
	_ = os.Setenv("JWT_SECRET", "my-secret-key")
	_ = os.Unsetenv("ENABLE_AUTH")
	defer func() { _ = os.Unsetenv("JWT_SECRET") }()

	config := NewServerConfig()

	// Auth is automatically enabled when JWT_SECRET is set
	if !config.EnableAuth {
		t.Error("Expected EnableAuth to be true when JWT_SECRET is set")
	}
}

func TestNewServerConfig_WithEnv(t *testing.T) {
	_ = os.Setenv("JWT_SECRET", "test-jwt-secret")
	_ = os.Setenv("AGENT_SECRET", "test-agent-secret")
	_ = os.Setenv("ENABLE_AUTH", "true")
	defer func() {
		_ = os.Unsetenv("JWT_SECRET")
		_ = os.Unsetenv("AGENT_SECRET")
		_ = os.Unsetenv("ENABLE_AUTH")
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
	services := &Services{
		Agent:     application.NewAgentService(&mockAgentRepo{}),
		Scenario:  application.NewScenarioService(&mockScenarioRepo{}, &mockTechniqueRepo{}, service.NewTechniqueValidator()),
		Technique: application.NewTechniqueService(&mockTechniqueRepo{}),
		Execution: application.NewExecutionService(&mockResultRepo{}, &mockScenarioRepo{}, &mockTechniqueRepo{}, &mockAgentRepo{}, nil, nil),
		Auth:      nil,
	}
	hub := websocket.NewHub(logger)

	config := &ServerConfig{
		JWTSecret:   "",
		AgentSecret: "",
		EnableAuth:  false,
	}

	server := NewServerWithConfig(services, hub, logger, config)

	if server == nil {
		t.Fatal("NewServerWithConfig returned nil")
	}

	if server.router == nil {
		t.Error("Server router is nil")
	}
}

func TestNewServerWithConfig_AuthEnabled(t *testing.T) {
	logger := zap.NewNop()
	services := &Services{
		Agent:     application.NewAgentService(&mockAgentRepo{}),
		Scenario:  application.NewScenarioService(&mockScenarioRepo{}, &mockTechniqueRepo{}, service.NewTechniqueValidator()),
		Technique: application.NewTechniqueService(&mockTechniqueRepo{}),
		Execution: application.NewExecutionService(&mockResultRepo{}, &mockScenarioRepo{}, &mockTechniqueRepo{}, &mockAgentRepo{}, nil, nil),
		Auth:      nil,
	}
	hub := websocket.NewHub(logger)

	config := &ServerConfig{
		JWTSecret:   "test-secret",
		AgentSecret: "agent-secret",
		EnableAuth:  true,
	}

	server := NewServerWithConfig(services, hub, logger, config)
	defer server.Close()

	if server == nil {
		t.Fatal("NewServerWithConfig returned nil")
	}
}

func TestNewServerWithConfig_NoHub(t *testing.T) {
	logger := zap.NewNop()
	services := &Services{
		Agent:     application.NewAgentService(&mockAgentRepo{}),
		Scenario:  application.NewScenarioService(&mockScenarioRepo{}, &mockTechniqueRepo{}, service.NewTechniqueValidator()),
		Technique: application.NewTechniqueService(&mockTechniqueRepo{}),
		Execution: application.NewExecutionService(&mockResultRepo{}, &mockScenarioRepo{}, &mockTechniqueRepo{}, &mockAgentRepo{}, nil, nil),
		Auth:      nil,
	}

	config := &ServerConfig{EnableAuth: false}

	// Should not panic with nil hub
	server := NewServerWithConfig(services, nil, logger, config)

	if server == nil {
		t.Fatal("NewServerWithConfig returned nil")
	}
}

func TestNewServer(t *testing.T) {
	os.Setenv("ENABLE_AUTH", "false")
	defer os.Unsetenv("ENABLE_AUTH")

	logger := zap.NewNop()
	services := &Services{
		Agent:     application.NewAgentService(&mockAgentRepo{}),
		Scenario:  application.NewScenarioService(&mockScenarioRepo{}, &mockTechniqueRepo{}, service.NewTechniqueValidator()),
		Technique: application.NewTechniqueService(&mockTechniqueRepo{}),
		Execution: application.NewExecutionService(&mockResultRepo{}, &mockScenarioRepo{}, &mockTechniqueRepo{}, &mockAgentRepo{}, nil, nil),
		Auth:      nil,
	}
	hub := websocket.NewHub(logger)

	server := NewServer(services, hub, logger)

	if server == nil {
		t.Fatal("NewServer returned nil")
	}
}

func TestServer_Router(t *testing.T) {
	logger := zap.NewNop()
	services := &Services{
		Agent:     application.NewAgentService(&mockAgentRepo{}),
		Scenario:  application.NewScenarioService(&mockScenarioRepo{}, &mockTechniqueRepo{}, service.NewTechniqueValidator()),
		Technique: application.NewTechniqueService(&mockTechniqueRepo{}),
		Execution: application.NewExecutionService(&mockResultRepo{}, &mockScenarioRepo{}, &mockTechniqueRepo{}, &mockAgentRepo{}, nil, nil),
		Auth:      nil,
	}

	config := &ServerConfig{EnableAuth: false}
	server := NewServerWithConfig(services, nil, logger, config)

	router := server.Router()
	if router == nil {
		t.Error("Router() returned nil")
	}
}

func TestServer_HealthEndpoint(t *testing.T) {
	logger := zap.NewNop()
	services := &Services{
		Agent:     application.NewAgentService(&mockAgentRepo{}),
		Scenario:  application.NewScenarioService(&mockScenarioRepo{}, &mockTechniqueRepo{}, service.NewTechniqueValidator()),
		Technique: application.NewTechniqueService(&mockTechniqueRepo{}),
		Execution: application.NewExecutionService(&mockResultRepo{}, &mockScenarioRepo{}, &mockTechniqueRepo{}, &mockAgentRepo{}, nil, nil),
		Auth:      nil,
	}

	config := &ServerConfig{EnableAuth: false}
	server := NewServerWithConfig(services, nil, logger, config)

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Verify response body contains auth_enabled
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}
	if response["status"] != "ok" {
		t.Errorf("Expected status 'ok', got %v", response["status"])
	}
	if response["auth_enabled"] != false {
		t.Errorf("Expected auth_enabled false, got %v", response["auth_enabled"])
	}
}

func TestServer_HealthEndpoint_AuthEnabled(t *testing.T) {
	logger := zap.NewNop()
	services := &Services{
		Agent:     application.NewAgentService(&mockAgentRepo{}),
		Scenario:  application.NewScenarioService(&mockScenarioRepo{}, &mockTechniqueRepo{}, service.NewTechniqueValidator()),
		Technique: application.NewTechniqueService(&mockTechniqueRepo{}),
		Execution: application.NewExecutionService(&mockResultRepo{}, &mockScenarioRepo{}, &mockTechniqueRepo{}, &mockAgentRepo{}, nil, nil),
		Auth:      nil,
	}

	config := &ServerConfig{
		EnableAuth: true,
		JWTSecret:  "test-secret",
	}
	server := NewServerWithConfig(services, nil, logger, config)
	defer server.Close()

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}
	if response["auth_enabled"] != true {
		t.Errorf("Expected auth_enabled true, got %v", response["auth_enabled"])
	}
}

func TestServer_APIRoutes(t *testing.T) {
	logger := zap.NewNop()
	services := &Services{
		Agent:     application.NewAgentService(&mockAgentRepo{}),
		Scenario:  application.NewScenarioService(&mockScenarioRepo{}, &mockTechniqueRepo{}, service.NewTechniqueValidator()),
		Technique: application.NewTechniqueService(&mockTechniqueRepo{}),
		Execution: application.NewExecutionService(&mockResultRepo{}, &mockScenarioRepo{}, &mockTechniqueRepo{}, &mockAgentRepo{}, nil, nil),
		Auth:      nil,
	}

	config := &ServerConfig{EnableAuth: false}
	server := NewServerWithConfig(services, nil, logger, config)

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
	services := &Services{
		Agent:     application.NewAgentService(&mockAgentRepo{}),
		Scenario:  application.NewScenarioService(&mockScenarioRepo{}, &mockTechniqueRepo{}, service.NewTechniqueValidator()),
		Technique: application.NewTechniqueService(&mockTechniqueRepo{}),
		Execution: application.NewExecutionService(&mockResultRepo{}, &mockScenarioRepo{}, &mockTechniqueRepo{}, &mockAgentRepo{}, nil, nil),
		Auth:      nil,
	}

	config := &ServerConfig{EnableAuth: false}
	server := NewServerWithConfig(services, nil, logger, config)

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

func TestNewServerConfig_EnableAuthTrueNoSecret(t *testing.T) {
	os.Unsetenv("JWT_SECRET")
	os.Setenv("ENABLE_AUTH", "true")
	defer os.Unsetenv("ENABLE_AUTH")

	config := NewServerConfig()

	// ENABLE_AUTH=true should enable auth even without secret
	if !config.EnableAuth {
		t.Error("Expected EnableAuth to be true when ENABLE_AUTH=true")
	}
}

func TestNewServerConfig_CustomDashboardPath(t *testing.T) {
	os.Setenv("DASHBOARD_PATH", "/custom/path")
	defer os.Unsetenv("DASHBOARD_PATH")

	config := NewServerConfig()

	if config.DashboardPath != "/custom/path" {
		t.Errorf("Expected dashboard path '/custom/path', got '%s'", config.DashboardPath)
	}
}

func TestServer_DashboardRoutes_InvalidPath(t *testing.T) {
	logger := zap.NewNop()
	services := &Services{
		Agent:     application.NewAgentService(&mockAgentRepo{}),
		Scenario:  application.NewScenarioService(&mockScenarioRepo{}, &mockTechniqueRepo{}, service.NewTechniqueValidator()),
		Technique: application.NewTechniqueService(&mockTechniqueRepo{}),
		Execution: application.NewExecutionService(&mockResultRepo{}, &mockScenarioRepo{}, &mockTechniqueRepo{}, &mockAgentRepo{}, nil, nil),
		Auth:      nil,
	}

	// Use a path that doesn't exist
	config := &ServerConfig{
		EnableAuth:    false,
		DashboardPath: "/nonexistent/path/that/does/not/exist",
	}

	// Should not panic, just log warning
	server := NewServerWithConfig(services, nil, logger, config)
	if server == nil {
		t.Fatal("Server should be created even with invalid dashboard path")
	}
}

func TestServer_NoRoute_APIPath(t *testing.T) {
	logger := zap.NewNop()
	services := &Services{
		Agent:     application.NewAgentService(&mockAgentRepo{}),
		Scenario:  application.NewScenarioService(&mockScenarioRepo{}, &mockTechniqueRepo{}, service.NewTechniqueValidator()),
		Technique: application.NewTechniqueService(&mockTechniqueRepo{}),
		Execution: application.NewExecutionService(&mockResultRepo{}, &mockScenarioRepo{}, &mockTechniqueRepo{}, &mockAgentRepo{}, nil, nil),
		Auth:      nil,
	}

	// Create temp directory with index.html for dashboard
	tmpDir := t.TempDir()
	indexFile := tmpDir + "/index.html"
	if err := os.WriteFile(indexFile, []byte("<html></html>"), 0644); err != nil {
		t.Fatalf("Failed to write index.html: %v", err)
	}
	if err := os.MkdirAll(tmpDir+"/assets", 0755); err != nil {
		t.Fatalf("Failed to create assets dir: %v", err)
	}

	config := &ServerConfig{
		EnableAuth:    false,
		DashboardPath: tmpDir,
	}

	server := NewServerWithConfig(services, nil, logger, config)

	// Request to non-existent API endpoint should return 404 JSON
	req, _ := http.NewRequest("GET", "/api/v1/nonexistent", nil)
	w := httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestServer_NoRoute_NonAPIPath(t *testing.T) {
	logger := zap.NewNop()
	services := &Services{
		Agent:     application.NewAgentService(&mockAgentRepo{}),
		Scenario:  application.NewScenarioService(&mockScenarioRepo{}, &mockTechniqueRepo{}, service.NewTechniqueValidator()),
		Technique: application.NewTechniqueService(&mockTechniqueRepo{}),
		Execution: application.NewExecutionService(&mockResultRepo{}, &mockScenarioRepo{}, &mockTechniqueRepo{}, &mockAgentRepo{}, nil, nil),
		Auth:      nil,
	}

	// Create temp directory with index.html for dashboard
	tmpDir := t.TempDir()
	indexFile := tmpDir + "/index.html"
	if err := os.WriteFile(indexFile, []byte("<html>Dashboard</html>"), 0644); err != nil {
		t.Fatalf("Failed to write index.html: %v", err)
	}
	if err := os.MkdirAll(tmpDir+"/assets", 0755); err != nil {
		t.Fatalf("Failed to create assets dir: %v", err)
	}

	config := &ServerConfig{
		EnableAuth:    false,
		DashboardPath: tmpDir,
	}

	server := NewServerWithConfig(services, nil, logger, config)

	// Request to non-API path should serve index.html (SPA fallback)
	req, _ := http.NewRequest("GET", "/some/spa/route", nil)
	w := httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if w.Body.String() != "<html>Dashboard</html>" {
		t.Errorf("Expected index.html content, got %s", w.Body.String())
	}
}

func TestServer_NoRoute_WSPath(t *testing.T) {
	logger := zap.NewNop()
	services := &Services{
		Agent:     application.NewAgentService(&mockAgentRepo{}),
		Scenario:  application.NewScenarioService(&mockScenarioRepo{}, &mockTechniqueRepo{}, service.NewTechniqueValidator()),
		Technique: application.NewTechniqueService(&mockTechniqueRepo{}),
		Execution: application.NewExecutionService(&mockResultRepo{}, &mockScenarioRepo{}, &mockTechniqueRepo{}, &mockAgentRepo{}, nil, nil),
		Auth:      nil,
	}

	tmpDir := t.TempDir()
	indexFile := tmpDir + "/index.html"
	if err := os.WriteFile(indexFile, []byte("<html>Dashboard</html>"), 0644); err != nil {
		t.Fatalf("Failed to write index.html: %v", err)
	}
	if err := os.MkdirAll(tmpDir+"/assets", 0755); err != nil {
		t.Fatalf("Failed to create assets dir: %v", err)
	}

	config := &ServerConfig{
		EnableAuth:    false,
		DashboardPath: tmpDir,
	}

	server := NewServerWithConfig(services, nil, logger, config)

	// Request to /ws/ path should return 404 JSON, not index.html
	req, _ := http.NewRequest("GET", "/ws/nonexistent", nil)
	w := httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestServer_WithAnalyticsService(t *testing.T) {
	logger := zap.NewNop()
	services := &Services{
		Agent:     application.NewAgentService(&mockAgentRepo{}),
		Scenario:  application.NewScenarioService(&mockScenarioRepo{}, &mockTechniqueRepo{}, service.NewTechniqueValidator()),
		Technique: application.NewTechniqueService(&mockTechniqueRepo{}),
		Execution: application.NewExecutionService(&mockResultRepo{}, &mockScenarioRepo{}, &mockTechniqueRepo{}, &mockAgentRepo{}, nil, nil),
		Analytics: application.NewAnalyticsService(&mockResultRepo{}),
		Auth:      nil,
	}

	config := &ServerConfig{EnableAuth: false}
	server := NewServerWithConfig(services, nil, logger, config)

	// Test that server was created with analytics service
	if server == nil {
		t.Fatal("Server should be created with analytics service")
	}

	// Verify routes exist
	routes := server.Router().Routes()
	analyticsRouteFound := false
	for _, route := range routes {
		if route.Path == "/api/v1/analytics/summary" {
			analyticsRouteFound = true
			break
		}
	}
	if !analyticsRouteFound {
		t.Error("Analytics routes should be registered when analytics service is provided")
	}
}

func TestServer_WithScheduleService(t *testing.T) {
	logger := zap.NewNop()
	scheduleService := application.NewScheduleService(
		&mockScheduleRepo{},
		nil,
		logger,
	)
	services := &Services{
		Agent:     application.NewAgentService(&mockAgentRepo{}),
		Scenario:  application.NewScenarioService(&mockScenarioRepo{}, &mockTechniqueRepo{}, service.NewTechniqueValidator()),
		Technique: application.NewTechniqueService(&mockTechniqueRepo{}),
		Execution: application.NewExecutionService(&mockResultRepo{}, &mockScenarioRepo{}, &mockTechniqueRepo{}, &mockAgentRepo{}, nil, nil),
		Schedule:  scheduleService,
		Auth:      nil,
	}

	config := &ServerConfig{EnableAuth: false}
	server := NewServerWithConfig(services, nil, logger, config)

	// Test that server was created with schedule service
	if server == nil {
		t.Fatal("Server should be created with schedule service")
	}

	// Verify routes exist
	routes := server.Router().Routes()
	scheduleRouteFound := false
	for _, route := range routes {
		if route.Path == "/api/v1/schedules" {
			scheduleRouteFound = true
			break
		}
	}
	if !scheduleRouteFound {
		t.Error("Schedule routes should be registered when schedule service is provided")
	}
}

func TestServer_WithNotificationService(t *testing.T) {
	logger := zap.NewNop()
	notificationService := application.NewNotificationService(
		&mockNotificationRepo{},
		&mockUserRepo{},
		nil,
		"https://localhost:8443",
		nil,
	)
	services := &Services{
		Agent:        application.NewAgentService(&mockAgentRepo{}),
		Scenario:     application.NewScenarioService(&mockScenarioRepo{}, &mockTechniqueRepo{}, service.NewTechniqueValidator()),
		Technique:    application.NewTechniqueService(&mockTechniqueRepo{}),
		Execution:    application.NewExecutionService(&mockResultRepo{}, &mockScenarioRepo{}, &mockTechniqueRepo{}, &mockAgentRepo{}, nil, nil),
		Notification: notificationService,
		Auth:         nil,
	}

	config := &ServerConfig{EnableAuth: false}
	server := NewServerWithConfig(services, nil, logger, config)

	// Test notifications endpoint exists
	req, _ := http.NewRequest("GET", "/api/v1/notifications", nil)
	w := httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	// Should return 200 or 401 depending on auth
	if w.Code != http.StatusOK && w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 200 or 401 for notifications endpoint, got %d", w.Code)
	}
}

func TestServer_WithWebSocketHub(t *testing.T) {
	logger := zap.NewNop()
	hub := websocket.NewHub(logger)
	services := &Services{
		Agent:     application.NewAgentService(&mockAgentRepo{}),
		Scenario:  application.NewScenarioService(&mockScenarioRepo{}, &mockTechniqueRepo{}, service.NewTechniqueValidator()),
		Technique: application.NewTechniqueService(&mockTechniqueRepo{}),
		Execution: application.NewExecutionService(&mockResultRepo{}, &mockScenarioRepo{}, &mockTechniqueRepo{}, &mockAgentRepo{}, nil, nil),
		Auth:      nil,
	}

	config := &ServerConfig{EnableAuth: false}
	server := NewServerWithConfig(services, hub, logger, config)

	// WebSocket routes should be registered
	routes := server.Router().Routes()
	wsRouteFound := false
	for _, route := range routes {
		if route.Path == "/ws/agent" || route.Path == "/ws/dashboard" {
			wsRouteFound = true
			break
		}
	}
	if !wsRouteFound {
		t.Error("WebSocket routes should be registered when hub is provided")
	}
}

func TestServer_EmptyDashboardPath(t *testing.T) {
	logger := zap.NewNop()
	services := &Services{
		Agent:     application.NewAgentService(&mockAgentRepo{}),
		Scenario:  application.NewScenarioService(&mockScenarioRepo{}, &mockTechniqueRepo{}, service.NewTechniqueValidator()),
		Technique: application.NewTechniqueService(&mockTechniqueRepo{}),
		Execution: application.NewExecutionService(&mockResultRepo{}, &mockScenarioRepo{}, &mockTechniqueRepo{}, &mockAgentRepo{}, nil, nil),
		Auth:      nil,
	}

	config := &ServerConfig{
		EnableAuth:    false,
		DashboardPath: "",
	}

	server := NewServerWithConfig(services, nil, logger, config)
	if server == nil {
		t.Fatal("Server should be created with empty dashboard path")
	}

	// Non-API routes should still return 404 (no SPA fallback)
	req, _ := http.NewRequest("GET", "/some/random/path", nil)
	w := httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected 404 for non-API path without dashboard, got %d", w.Code)
	}
}

func TestServer_DefaultDashboardPath(t *testing.T) {
	os.Unsetenv("DASHBOARD_PATH")

	config := NewServerConfig()
	if config.DashboardPath != "../dashboard/dist" {
		t.Errorf("Expected default dashboard path '../dashboard/dist', got '%s'", config.DashboardPath)
	}
}

func TestServer_AgentsEndpoint(t *testing.T) {
	logger := zap.NewNop()
	services := &Services{
		Agent:     application.NewAgentService(&mockAgentRepo{}),
		Scenario:  application.NewScenarioService(&mockScenarioRepo{}, &mockTechniqueRepo{}, service.NewTechniqueValidator()),
		Technique: application.NewTechniqueService(&mockTechniqueRepo{}),
		Execution: application.NewExecutionService(&mockResultRepo{}, &mockScenarioRepo{}, &mockTechniqueRepo{}, &mockAgentRepo{}, nil, nil),
		Auth:      nil,
	}

	config := &ServerConfig{EnableAuth: false}
	server := NewServerWithConfig(services, nil, logger, config)

	req, _ := http.NewRequest("GET", "/api/v1/agents", nil)
	w := httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for agents endpoint, got %d", w.Code)
	}
}

func TestServer_TechniquesEndpoint(t *testing.T) {
	logger := zap.NewNop()
	services := &Services{
		Agent:     application.NewAgentService(&mockAgentRepo{}),
		Scenario:  application.NewScenarioService(&mockScenarioRepo{}, &mockTechniqueRepo{}, service.NewTechniqueValidator()),
		Technique: application.NewTechniqueService(&mockTechniqueRepo{}),
		Execution: application.NewExecutionService(&mockResultRepo{}, &mockScenarioRepo{}, &mockTechniqueRepo{}, &mockAgentRepo{}, nil, nil),
		Auth:      nil,
	}

	config := &ServerConfig{EnableAuth: false}
	server := NewServerWithConfig(services, nil, logger, config)

	req, _ := http.NewRequest("GET", "/api/v1/techniques", nil)
	w := httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for techniques endpoint, got %d", w.Code)
	}
}

func TestServer_ScenariosEndpoint(t *testing.T) {
	logger := zap.NewNop()
	services := &Services{
		Agent:     application.NewAgentService(&mockAgentRepo{}),
		Scenario:  application.NewScenarioService(&mockScenarioRepo{}, &mockTechniqueRepo{}, service.NewTechniqueValidator()),
		Technique: application.NewTechniqueService(&mockTechniqueRepo{}),
		Execution: application.NewExecutionService(&mockResultRepo{}, &mockScenarioRepo{}, &mockTechniqueRepo{}, &mockAgentRepo{}, nil, nil),
		Auth:      nil,
	}

	config := &ServerConfig{EnableAuth: false}
	server := NewServerWithConfig(services, nil, logger, config)

	// Verify routes exist
	routes := server.Router().Routes()
	scenarioRouteFound := false
	for _, route := range routes {
		if route.Path == "/api/v1/scenarios" {
			scenarioRouteFound = true
			break
		}
	}
	if !scenarioRouteFound {
		t.Error("Scenario routes should be registered")
	}
}

func TestServer_ExecutionsEndpoint(t *testing.T) {
	logger := zap.NewNop()
	services := &Services{
		Agent:     application.NewAgentService(&mockAgentRepo{}),
		Scenario:  application.NewScenarioService(&mockScenarioRepo{}, &mockTechniqueRepo{}, service.NewTechniqueValidator()),
		Technique: application.NewTechniqueService(&mockTechniqueRepo{}),
		Execution: application.NewExecutionService(&mockResultRepo{}, &mockScenarioRepo{}, &mockTechniqueRepo{}, &mockAgentRepo{}, nil, nil),
		Auth:      nil,
	}

	config := &ServerConfig{EnableAuth: false}
	server := NewServerWithConfig(services, nil, logger, config)

	req, _ := http.NewRequest("GET", "/api/v1/executions", nil)
	w := httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for executions endpoint, got %d", w.Code)
	}
}

func TestServer_PermissionsEndpoint(t *testing.T) {
	logger := zap.NewNop()
	services := &Services{
		Agent:     application.NewAgentService(&mockAgentRepo{}),
		Scenario:  application.NewScenarioService(&mockScenarioRepo{}, &mockTechniqueRepo{}, service.NewTechniqueValidator()),
		Technique: application.NewTechniqueService(&mockTechniqueRepo{}),
		Execution: application.NewExecutionService(&mockResultRepo{}, &mockScenarioRepo{}, &mockTechniqueRepo{}, &mockAgentRepo{}, nil, nil),
		Auth:      nil,
	}

	config := &ServerConfig{EnableAuth: false}
	server := NewServerWithConfig(services, nil, logger, config)

	// Verify routes exist
	routes := server.Router().Routes()
	permissionRouteFound := false
	for _, route := range routes {
		if route.Path == "/api/v1/permissions/matrix" || route.Path == "/api/v1/permissions/me" {
			permissionRouteFound = true
			break
		}
	}
	if !permissionRouteFound {
		t.Error("Permission routes should be registered")
	}
}

func TestServer_StaticFiles(t *testing.T) {
	logger := zap.NewNop()
	services := &Services{
		Agent:     application.NewAgentService(&mockAgentRepo{}),
		Scenario:  application.NewScenarioService(&mockScenarioRepo{}, &mockTechniqueRepo{}, service.NewTechniqueValidator()),
		Technique: application.NewTechniqueService(&mockTechniqueRepo{}),
		Execution: application.NewExecutionService(&mockResultRepo{}, &mockScenarioRepo{}, &mockTechniqueRepo{}, &mockAgentRepo{}, nil, nil),
		Auth:      nil,
	}

	tmpDir := t.TempDir()
	indexFile := tmpDir + "/index.html"
	if err := os.WriteFile(indexFile, []byte("<html></html>"), 0644); err != nil {
		t.Fatalf("Failed to write index.html: %v", err)
	}
	assetsDir := tmpDir + "/assets"
	if err := os.MkdirAll(assetsDir, 0755); err != nil {
		t.Fatalf("Failed to create assets dir: %v", err)
	}
	if err := os.WriteFile(assetsDir+"/main.js", []byte("console.log('test');"), 0644); err != nil {
		t.Fatalf("Failed to write main.js: %v", err)
	}
	if err := os.WriteFile(tmpDir+"/favicon.ico", []byte("icon"), 0644); err != nil {
		t.Fatalf("Failed to write favicon.ico: %v", err)
	}
	if err := os.WriteFile(tmpDir+"/vite.svg", []byte("<svg></svg>"), 0644); err != nil {
		t.Fatalf("Failed to write vite.svg: %v", err)
	}

	config := &ServerConfig{
		EnableAuth:    false,
		DashboardPath: tmpDir,
	}

	server := NewServerWithConfig(services, nil, logger, config)

	// Test static asset
	req, _ := http.NewRequest("GET", "/assets/main.js", nil)
	w := httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for static asset, got %d", w.Code)
	}

	// Test favicon
	req, _ = http.NewRequest("GET", "/favicon.ico", nil)
	w = httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for favicon, got %d", w.Code)
	}

	// Test vite.svg
	req, _ = http.NewRequest("GET", "/vite.svg", nil)
	w = httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for vite.svg, got %d", w.Code)
	}
}

// mockScheduleRepo implements repository.ScheduleRepository for testing
type mockScheduleRepo struct{}

func (m *mockScheduleRepo) Create(ctx context.Context, schedule *entity.Schedule) error {
	return nil
}
func (m *mockScheduleRepo) Update(ctx context.Context, schedule *entity.Schedule) error {
	return nil
}
func (m *mockScheduleRepo) Delete(ctx context.Context, id string) error { return nil }
func (m *mockScheduleRepo) FindByID(ctx context.Context, id string) (*entity.Schedule, error) {
	return nil, nil
}
func (m *mockScheduleRepo) FindAll(ctx context.Context) ([]*entity.Schedule, error) {
	return []*entity.Schedule{}, nil
}
func (m *mockScheduleRepo) FindActiveSchedules(ctx context.Context) ([]*entity.Schedule, error) {
	return []*entity.Schedule{}, nil
}
func (m *mockScheduleRepo) FindByScenarioID(ctx context.Context, scenarioID string) ([]*entity.Schedule, error) {
	return []*entity.Schedule{}, nil
}
func (m *mockScheduleRepo) UpdateStatus(ctx context.Context, id string, status entity.ScheduleStatus) error {
	return nil
}
func (m *mockScheduleRepo) UpdateLastRun(ctx context.Context, id string, lastRunAt time.Time, lastRunID string) error {
	return nil
}
func (m *mockScheduleRepo) UpdateNextRun(ctx context.Context, id string, nextRunAt time.Time) error {
	return nil
}
func (m *mockScheduleRepo) CreateRun(ctx context.Context, run *entity.ScheduleRun) error {
	return nil
}
func (m *mockScheduleRepo) UpdateRun(ctx context.Context, run *entity.ScheduleRun) error {
	return nil
}
func (m *mockScheduleRepo) FindRunsByScheduleID(ctx context.Context, scheduleID string, limit int) ([]*entity.ScheduleRun, error) {
	return []*entity.ScheduleRun{}, nil
}
func (m *mockScheduleRepo) FindActiveSchedulesDue(ctx context.Context, before time.Time) ([]*entity.Schedule, error) {
	return []*entity.Schedule{}, nil
}
func (m *mockScheduleRepo) FindByStatus(ctx context.Context, status entity.ScheduleStatus) ([]*entity.Schedule, error) {
	return []*entity.Schedule{}, nil
}

// mockNotificationRepo implements repository.NotificationRepository for testing
type mockNotificationRepo struct{}

func (m *mockNotificationRepo) CreateSettings(ctx context.Context, settings *entity.NotificationSettings) error {
	return nil
}
func (m *mockNotificationRepo) UpdateSettings(ctx context.Context, settings *entity.NotificationSettings) error {
	return nil
}
func (m *mockNotificationRepo) DeleteSettings(ctx context.Context, id string) error { return nil }
func (m *mockNotificationRepo) FindSettingsByID(ctx context.Context, id string) (*entity.NotificationSettings, error) {
	return nil, nil
}
func (m *mockNotificationRepo) FindSettingsByUserID(ctx context.Context, userID string) (*entity.NotificationSettings, error) {
	return nil, nil
}
func (m *mockNotificationRepo) FindAllEnabledSettings(ctx context.Context) ([]*entity.NotificationSettings, error) {
	return []*entity.NotificationSettings{}, nil
}
func (m *mockNotificationRepo) CreateNotification(ctx context.Context, notification *entity.Notification) error {
	return nil
}
func (m *mockNotificationRepo) UpdateNotification(ctx context.Context, notification *entity.Notification) error {
	return nil
}
func (m *mockNotificationRepo) FindNotificationByID(ctx context.Context, id string) (*entity.Notification, error) {
	return nil, nil
}
func (m *mockNotificationRepo) FindNotificationsByUserID(ctx context.Context, userID string, limit int) ([]*entity.Notification, error) {
	return []*entity.Notification{}, nil
}
func (m *mockNotificationRepo) FindUnreadByUserID(ctx context.Context, userID string) ([]*entity.Notification, error) {
	return []*entity.Notification{}, nil
}
func (m *mockNotificationRepo) MarkAsRead(ctx context.Context, id string) error { return nil }
func (m *mockNotificationRepo) MarkAllAsRead(ctx context.Context, userID string) error {
	return nil
}

// mockUserRepo implements repository.UserRepository for testing
type mockUserRepo struct{}

func (m *mockUserRepo) Create(ctx context.Context, user *entity.User) error       { return nil }
func (m *mockUserRepo) Update(ctx context.Context, user *entity.User) error       { return nil }
func (m *mockUserRepo) Delete(ctx context.Context, id string) error               { return nil }
func (m *mockUserRepo) FindByID(ctx context.Context, id string) (*entity.User, error) { return nil, nil }
func (m *mockUserRepo) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	return nil, nil
}
func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	return nil, nil
}
func (m *mockUserRepo) FindAll(ctx context.Context) ([]*entity.User, error) { return nil, nil }
func (m *mockUserRepo) FindActive(ctx context.Context) ([]*entity.User, error) { return nil, nil }
func (m *mockUserRepo) UpdateLastLogin(ctx context.Context, id string) error { return nil }
func (m *mockUserRepo) Deactivate(ctx context.Context, id string) error      { return nil }
func (m *mockUserRepo) Reactivate(ctx context.Context, id string) error      { return nil }
func (m *mockUserRepo) CountByRole(ctx context.Context, role entity.UserRole) (int, error) {
	return 0, nil
}
func (m *mockUserRepo) DeactivateAdminIfNotLast(ctx context.Context, id string) error { return nil }

// --- Helper to create standard test services ---
func createTestServices(t *testing.T) *Services {
	t.Helper()
	return &Services{
		Agent:     application.NewAgentService(&mockAgentRepo{}),
		Scenario:  application.NewScenarioService(&mockScenarioRepo{}, &mockTechniqueRepo{}, service.NewTechniqueValidator()),
		Technique: application.NewTechniqueService(&mockTechniqueRepo{}),
		Execution: application.NewExecutionService(&mockResultRepo{}, &mockScenarioRepo{}, &mockTechniqueRepo{}, &mockAgentRepo{}, nil, nil),
		Auth:      nil,
	}
}

func createTestServicesWithAuth(t *testing.T) *Services {
	t.Helper()
	s := createTestServices(t)
	s.Auth = application.NewAuthService(&mockUserRepo{}, "test-jwt-secret-key")
	return s
}

// --- Middleware Order Tests ---

func TestServer_MiddlewareOrder_LoggingAndRecovery(t *testing.T) {
	logger := zap.NewNop()
	services := createTestServices(t)

	config := &ServerConfig{EnableAuth: false}
	server := NewServerWithConfig(services, nil, logger, config)

	// A normal request should go through logging + recovery middleware without issues
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 through middleware chain, got %d", w.Code)
	}
}

func TestServer_MiddlewareOrder_AuthBeforeHandlers(t *testing.T) {
	logger := zap.NewNop()
	services := createTestServicesWithAuth(t)

	config := &ServerConfig{
		EnableAuth: true,
		JWTSecret:  "test-jwt-secret-key",
	}
	server := NewServerWithConfig(services, nil, logger, config)
	defer server.Close()

	// Without auth header, API should return 401 (auth middleware runs before handler)
	req, _ := http.NewRequest("GET", "/api/v1/agents", nil)
	w := httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 when auth middleware runs before handler, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	if response["error"] != "authorization header required" {
		t.Errorf("Expected auth error message, got %v", response["error"])
	}
}

func TestServer_MiddlewareOrder_InvalidBearerToken(t *testing.T) {
	logger := zap.NewNop()
	services := createTestServicesWithAuth(t)

	config := &ServerConfig{
		EnableAuth: true,
		JWTSecret:  "test-jwt-secret-key",
	}
	server := NewServerWithConfig(services, nil, logger, config)
	defer server.Close()

	// With invalid bearer token format
	req, _ := http.NewRequest("GET", "/api/v1/agents", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	w := httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 for invalid bearer format, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	if response["error"] != "invalid authorization header format" {
		t.Errorf("Expected invalid format error, got %v", response["error"])
	}
}

func TestServer_MiddlewareOrder_NoAuthSetsDefaultContext(t *testing.T) {
	logger := zap.NewNop()
	services := createTestServices(t)

	config := &ServerConfig{EnableAuth: false}
	server := NewServerWithConfig(services, nil, logger, config)

	// Without auth, NoAuthMiddleware should set default user context (anonymous/admin)
	// so API requests should succeed
	req, _ := http.NewRequest("GET", "/api/v1/agents", nil)
	w := httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 with NoAuth middleware, got %d", w.Code)
	}
}

func TestServer_MiddlewareOrder_HealthBypassesAuth(t *testing.T) {
	logger := zap.NewNop()
	services := createTestServicesWithAuth(t)

	config := &ServerConfig{
		EnableAuth: true,
		JWTSecret:  "test-jwt-secret-key",
	}
	server := NewServerWithConfig(services, nil, logger, config)
	defer server.Close()

	// Health endpoint should be accessible without authentication
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected health endpoint to bypass auth with status 200, got %d", w.Code)
	}
}

// --- CORS Headers Tests ---

func TestServer_ResponseHeaders_ContentType(t *testing.T) {
	logger := zap.NewNop()
	services := createTestServices(t)

	config := &ServerConfig{EnableAuth: false}
	server := NewServerWithConfig(services, nil, logger, config)

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	contentType := w.Header().Get("Content-Type")
	if contentType == "" {
		t.Error("Expected Content-Type header to be set")
	}
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Expected Content-Type to contain 'application/json', got '%s'", contentType)
	}
}

func TestServer_ResponseHeaders_APIEndpoint(t *testing.T) {
	logger := zap.NewNop()
	services := createTestServices(t)

	config := &ServerConfig{EnableAuth: false}
	server := NewServerWithConfig(services, nil, logger, config)

	req, _ := http.NewRequest("GET", "/api/v1/agents", nil)
	w := httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	contentType := w.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Expected JSON Content-Type for API endpoint, got '%s'", contentType)
	}
}

// --- Auth Service Integration Tests ---

func TestServer_WithAuthService_RegistersAuthRoutes(t *testing.T) {
	logger := zap.NewNop()
	services := createTestServicesWithAuth(t)

	config := &ServerConfig{
		EnableAuth: true,
		JWTSecret:  "test-jwt-secret-key",
	}
	server := NewServerWithConfig(services, nil, logger, config)
	defer server.Close()

	routes := server.Router().Routes()
	expectedAuthRoutes := []string{
		"/api/v1/auth/login",
		"/api/v1/auth/refresh",
		"/api/v1/auth/logout",
	}

	for _, expected := range expectedAuthRoutes {
		found := false
		for _, route := range routes {
			if route.Path == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected auth route %s not found", expected)
		}
	}
}

func TestServer_WithAuthService_RegistersAdminRoutes(t *testing.T) {
	logger := zap.NewNop()
	services := createTestServicesWithAuth(t)

	config := &ServerConfig{
		EnableAuth: true,
		JWTSecret:  "test-jwt-secret-key",
	}
	server := NewServerWithConfig(services, nil, logger, config)
	defer server.Close()

	routes := server.Router().Routes()
	expectedAdminRoutes := []string{
		"/api/v1/admin/users",
		"/api/v1/admin/users/:id",
	}

	for _, expected := range expectedAdminRoutes {
		found := false
		for _, route := range routes {
			if route.Path == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected admin route %s not found", expected)
		}
	}
}

func TestServer_WithAuthService_RegistersProtectedMeRoute(t *testing.T) {
	logger := zap.NewNop()
	services := createTestServicesWithAuth(t)

	config := &ServerConfig{
		EnableAuth: true,
		JWTSecret:  "test-jwt-secret-key",
	}
	server := NewServerWithConfig(services, nil, logger, config)
	defer server.Close()

	routes := server.Router().Routes()
	found := false
	for _, route := range routes {
		if route.Path == "/api/v1/auth/me" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected protected /api/v1/auth/me route not found")
	}
}

func TestServer_WithAuthService_LoginEndpointAccessible(t *testing.T) {
	logger := zap.NewNop()
	services := createTestServicesWithAuth(t)

	config := &ServerConfig{
		EnableAuth: true,
		JWTSecret:  "test-jwt-secret-key",
	}
	server := NewServerWithConfig(services, nil, logger, config)
	defer server.Close()

	// Login endpoint should be accessible without auth (it's a public route)
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(`{"username":"test","password":"test"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	// Should not return 401 (the route itself may return 401 for bad credentials,
	// but NOT because of middleware blocking)
	if w.Code == http.StatusNotFound {
		t.Error("Login endpoint should exist and not return 404")
	}
}

func TestServer_WithAuthService_AdminRoutesRequireAuth(t *testing.T) {
	logger := zap.NewNop()
	services := createTestServicesWithAuth(t)

	config := &ServerConfig{
		EnableAuth: true,
		JWTSecret:  "test-jwt-secret-key",
	}
	server := NewServerWithConfig(services, nil, logger, config)
	defer server.Close()

	// Admin endpoint without auth should return 401
	req, _ := http.NewRequest("GET", "/api/v1/admin/users", nil)
	w := httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 for admin route without auth, got %d", w.Code)
	}
}

// --- NoRoute / SPA Fallback Edge Cases ---

func TestServer_NoRoute_MultiplePrefixes(t *testing.T) {
	logger := zap.NewNop()
	services := createTestServices(t)

	tmpDir := t.TempDir()
	if err := os.WriteFile(tmpDir+"/index.html", []byte("<html>SPA</html>"), 0644); err != nil {
		t.Fatalf("Failed to write index.html: %v", err)
	}
	if err := os.MkdirAll(tmpDir+"/assets", 0755); err != nil {
		t.Fatalf("Failed to create assets dir: %v", err)
	}

	config := &ServerConfig{
		EnableAuth:    false,
		DashboardPath: tmpDir,
	}
	server := NewServerWithConfig(services, nil, logger, config)

	tests := []struct {
		name           string
		path           string
		expectStatus   int
		expectSPABody  bool
	}{
		{"API path returns 404 JSON", "/api/v1/unknown", http.StatusNotFound, false},
		{"WS path returns 404 JSON", "/ws/unknown", http.StatusNotFound, false},
		{"SPA route serves index", "/dashboard", http.StatusOK, true},
		{"Nested SPA route serves index", "/settings/profile", http.StatusOK, true},
		{"Root-level SPA route", "/login", http.StatusOK, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()
			server.Router().ServeHTTP(w, req)

			if w.Code != tt.expectStatus {
				t.Errorf("Expected status %d for path %s, got %d", tt.expectStatus, tt.path, w.Code)
			}
			if tt.expectSPABody && w.Body.String() != "<html>SPA</html>" {
				t.Errorf("Expected SPA body for path %s, got '%s'", tt.path, w.Body.String())
			}
		})
	}
}

// --- Config Edge Cases ---

func TestNewServerConfig_EnableAuthFalseOverridesSecret(t *testing.T) {
	os.Setenv("JWT_SECRET", "some-secret")
	os.Setenv("ENABLE_AUTH", "false")
	defer func() {
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("ENABLE_AUTH")
	}()

	config := NewServerConfig()

	// ENABLE_AUTH=false should override even when JWT_SECRET is set
	if config.EnableAuth {
		t.Error("Expected EnableAuth to be false when ENABLE_AUTH=false, even with JWT_SECRET set")
	}
	if config.JWTSecret != "some-secret" {
		t.Errorf("JWTSecret should still be set to 'some-secret', got '%s'", config.JWTSecret)
	}
}

func TestNewServerConfig_EnableAuthInvalidValue(t *testing.T) {
	os.Unsetenv("JWT_SECRET")
	os.Setenv("ENABLE_AUTH", "maybe")
	defer os.Unsetenv("ENABLE_AUTH")

	config := NewServerConfig()

	// Invalid ENABLE_AUTH value should not change the default behavior
	// (default: auth disabled when JWT_SECRET is not set)
	if config.EnableAuth {
		t.Error("Expected EnableAuth to be false when ENABLE_AUTH has invalid value and no JWT_SECRET")
	}
}

func TestNewServerConfig_EnableAuthInvalidValueWithSecret(t *testing.T) {
	os.Setenv("JWT_SECRET", "my-secret")
	os.Setenv("ENABLE_AUTH", "maybe")
	defer func() {
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("ENABLE_AUTH")
	}()

	config := NewServerConfig()

	// Invalid ENABLE_AUTH should not override the JWT_SECRET-based default (true)
	if !config.EnableAuth {
		t.Error("Expected EnableAuth to be true when JWT_SECRET is set and ENABLE_AUTH is invalid")
	}
}

// --- Multiple Services Integration ---

func TestServer_AllServicesNil_ExceptRequired(t *testing.T) {
	logger := zap.NewNop()
	services := &Services{
		Agent:     application.NewAgentService(&mockAgentRepo{}),
		Scenario:  application.NewScenarioService(&mockScenarioRepo{}, &mockTechniqueRepo{}, service.NewTechniqueValidator()),
		Technique: application.NewTechniqueService(&mockTechniqueRepo{}),
		Execution: application.NewExecutionService(&mockResultRepo{}, &mockScenarioRepo{}, &mockTechniqueRepo{}, &mockAgentRepo{}, nil, nil),
		// All optional services nil
		Auth:         nil,
		Analytics:    nil,
		Notification: nil,
		Schedule:     nil,
	}

	config := &ServerConfig{EnableAuth: false}
	server := NewServerWithConfig(services, nil, logger, config)

	if server == nil {
		t.Fatal("Server should be created with optional services set to nil")
	}

	// Verify optional routes are not registered
	routes := server.Router().Routes()
	for _, route := range routes {
		if strings.Contains(route.Path, "/analytics/") {
			t.Error("Analytics routes should not be registered when analytics service is nil")
		}
		if strings.Contains(route.Path, "/schedules") && route.Path != "" {
			t.Error("Schedule routes should not be registered when schedule service is nil")
		}
		if strings.Contains(route.Path, "/notifications") {
			t.Error("Notification routes should not be registered when notification service is nil")
		}
		if strings.Contains(route.Path, "/admin/") {
			t.Error("Admin routes should not be registered when auth service is nil")
		}
	}
}

func TestServer_AllServicesProvided(t *testing.T) {
	logger := zap.NewNop()
	hub := websocket.NewHub(logger)

	scheduleService := application.NewScheduleService(&mockScheduleRepo{}, nil, logger)
	notificationService := application.NewNotificationService(&mockNotificationRepo{}, &mockUserRepo{}, nil, "https://localhost:8443", nil)

	services := &Services{
		Agent:        application.NewAgentService(&mockAgentRepo{}),
		Scenario:     application.NewScenarioService(&mockScenarioRepo{}, &mockTechniqueRepo{}, service.NewTechniqueValidator()),
		Technique:    application.NewTechniqueService(&mockTechniqueRepo{}),
		Execution:    application.NewExecutionService(&mockResultRepo{}, &mockScenarioRepo{}, &mockTechniqueRepo{}, &mockAgentRepo{}, nil, nil),
		Auth:         application.NewAuthService(&mockUserRepo{}, "full-service-secret"),
		Analytics:    application.NewAnalyticsService(&mockResultRepo{}),
		Notification: notificationService,
		Schedule:     scheduleService,
	}

	config := &ServerConfig{
		EnableAuth: true,
		JWTSecret:  "full-service-secret",
	}
	server := NewServerWithConfig(services, hub, logger, config)
	defer server.Close()

	if server == nil {
		t.Fatal("Server should be created with all services provided")
	}

	// Count total routes - should have all service routes
	routes := server.Router().Routes()
	if len(routes) < 20 {
		t.Errorf("Expected at least 20 routes with all services, got %d", len(routes))
	}
}

// --- Body Size Limit Integration Tests ---

func TestServer_BodySizeLimit_LargePostRejected(t *testing.T) {
	logger := zap.NewNop()
	services := createTestServices(t)

	config := &ServerConfig{EnableAuth: false}
	server := NewServerWithConfig(services, nil, logger, config)

	// Create a body larger than 10 MB
	largeBody := strings.Repeat("x", 11<<20) // 11 MB
	req, _ := http.NewRequest("POST", "/api/v1/agents", strings.NewReader(largeBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("Expected status 413 for body > 10MB, got %d", w.Code)
	}
}

func TestServer_BodySizeLimit_SmallPostAllowed(t *testing.T) {
	logger := zap.NewNop()
	services := createTestServices(t)

	config := &ServerConfig{EnableAuth: false}
	server := NewServerWithConfig(services, nil, logger, config)

	smallBody := `{"paw":"test-agent","hostname":"host","platform":"linux"}`
	req, _ := http.NewRequest("POST", "/api/v1/agents", strings.NewReader(smallBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	// Should not be rejected by body size limit (handler may return other status codes)
	if w.Code == http.StatusRequestEntityTooLarge {
		t.Error("Small body should not be rejected by body size limit")
	}
}

func TestServer_BodySizeLimit_HealthGetUnaffected(t *testing.T) {
	logger := zap.NewNop()
	services := createTestServices(t)

	config := &ServerConfig{EnableAuth: false}
	server := NewServerWithConfig(services, nil, logger, config)

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 for GET /health, got %d", w.Code)
	}
}

func TestServer_BodySizeLimit_AuthLoginLargeBodyRejected(t *testing.T) {
	logger := zap.NewNop()
	services := createTestServicesWithAuth(t)

	config := &ServerConfig{
		EnableAuth: true,
		JWTSecret:  "test-jwt-secret-key",
	}
	server := NewServerWithConfig(services, nil, logger, config)
	defer server.Close()

	// Large body to auth login endpoint
	largeBody := `{"username":"` + strings.Repeat("a", 11<<20) + `","password":"test"}`
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(largeBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("Expected status 413 for large auth login body, got %d", w.Code)
	}
}

// --- HTTP Method Verification ---

func TestServer_HealthEndpoint_MethodNotAllowed(t *testing.T) {
	logger := zap.NewNop()
	services := createTestServices(t)

	config := &ServerConfig{EnableAuth: false}
	server := NewServerWithConfig(services, nil, logger, config)

	// POST to health endpoint should fail (only GET is registered)
	req, _ := http.NewRequest("POST", "/health", nil)
	w := httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		t.Error("POST to /health should not return 200")
	}
}

func TestServer_APIEndpoints_MethodRouting(t *testing.T) {
	logger := zap.NewNop()
	services := createTestServices(t)

	config := &ServerConfig{EnableAuth: false}
	server := NewServerWithConfig(services, nil, logger, config)

	// GET /api/v1/agents should work
	req, _ := http.NewRequest("GET", "/api/v1/agents", nil)
	w := httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 for GET /api/v1/agents, got %d", w.Code)
	}

	// DELETE /api/v1/agents (without paw) should return 404 or 405
	req, _ = http.NewRequest("DELETE", "/api/v1/agents", nil)
	w = httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		t.Error("DELETE /api/v1/agents without paw should not return 200")
	}
}
