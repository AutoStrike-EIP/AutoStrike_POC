package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"autostrike/internal/application"
	"autostrike/internal/domain/entity"
	"autostrike/internal/domain/service"

	"github.com/gin-gonic/gin"
)

// testScenarioRepo is a mock repository for scenario handler tests
type testScenarioRepo struct {
	scenarios map[string]*entity.Scenario
	err       error
}

func newTestScenarioRepo() *testScenarioRepo {
	return &testScenarioRepo{scenarios: make(map[string]*entity.Scenario)}
}

func (m *testScenarioRepo) Create(ctx context.Context, s *entity.Scenario) error {
	if m.err != nil {
		return m.err
	}
	m.scenarios[s.ID] = s
	return nil
}

func (m *testScenarioRepo) Update(ctx context.Context, s *entity.Scenario) error {
	if m.err != nil {
		return m.err
	}
	m.scenarios[s.ID] = s
	return nil
}

func (m *testScenarioRepo) Delete(ctx context.Context, id string) error {
	if m.err != nil {
		return m.err
	}
	delete(m.scenarios, id)
	return nil
}

func (m *testScenarioRepo) FindByID(ctx context.Context, id string) (*entity.Scenario, error) {
	if m.err != nil {
		return nil, m.err
	}
	s, ok := m.scenarios[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return s, nil
}

func (m *testScenarioRepo) FindAll(ctx context.Context) ([]*entity.Scenario, error) {
	if m.err != nil {
		return nil, m.err
	}
	result := make([]*entity.Scenario, 0, len(m.scenarios))
	for _, s := range m.scenarios {
		result = append(result, s)
	}
	return result, nil
}

func (m *testScenarioRepo) FindByTag(ctx context.Context, tag string) ([]*entity.Scenario, error) {
	if m.err != nil {
		return nil, m.err
	}
	var result []*entity.Scenario
	for _, s := range m.scenarios {
		for _, t := range s.Tags {
			if t == tag {
				result = append(result, s)
				break
			}
		}
	}
	return result, nil
}

// testTechniqueRepo is a mock technique repository for validation
type testTechniqueRepo struct {
	techniques map[string]*entity.Technique
	err        error
}

func newTestTechniqueRepo() *testTechniqueRepo {
	return &testTechniqueRepo{techniques: make(map[string]*entity.Technique)}
}

func (m *testTechniqueRepo) Create(ctx context.Context, t *entity.Technique) error  { return m.err }
func (m *testTechniqueRepo) Update(ctx context.Context, t *entity.Technique) error  { return m.err }
func (m *testTechniqueRepo) Delete(ctx context.Context, id string) error            { return m.err }
func (m *testTechniqueRepo) FindByID(ctx context.Context, id string) (*entity.Technique, error) {
	if m.err != nil {
		return nil, m.err
	}
	t, ok := m.techniques[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return t, nil
}
func (m *testTechniqueRepo) FindAll(ctx context.Context) ([]*entity.Technique, error) {
	if m.err != nil {
		return nil, m.err
	}
	result := make([]*entity.Technique, 0, len(m.techniques))
	for _, t := range m.techniques {
		result = append(result, t)
	}
	return result, nil
}
func (m *testTechniqueRepo) FindByTactic(ctx context.Context, tactic entity.TacticType) ([]*entity.Technique, error) {
	return nil, nil
}
func (m *testTechniqueRepo) FindByPlatform(ctx context.Context, platform string) ([]*entity.Technique, error) {
	return nil, nil
}
func (m *testTechniqueRepo) ImportFromYAML(ctx context.Context, path string) error { return nil }

// createTestScenarioService creates a test scenario service
func createTestScenarioService(scenarioRepo *testScenarioRepo, techRepo *testTechniqueRepo) *application.ScenarioService {
	validator := service.NewTechniqueValidator()
	return application.NewScenarioService(scenarioRepo, techRepo, validator)
}

func TestNewScenarioHandler(t *testing.T) {
	scenarioRepo := newTestScenarioRepo()
	techRepo := newTestTechniqueRepo()
	svc := createTestScenarioService(scenarioRepo, techRepo)
	handler := NewScenarioHandler(svc)
	if handler == nil {
		t.Error("Expected non-nil handler")
	}
}

func TestScenarioHandler_RegisterRoutes(t *testing.T) {
	scenarioRepo := newTestScenarioRepo()
	techRepo := newTestTechniqueRepo()
	svc := createTestScenarioService(scenarioRepo, techRepo)
	handler := NewScenarioHandler(svc)

	router := gin.New()
	api := router.Group("/api")
	handler.RegisterRoutes(api)

	// Verify routes are registered
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/scenarios", nil)
	router.ServeHTTP(w, req)
	if w.Code == http.StatusNotFound {
		t.Error("Route /api/scenarios not registered")
	}
}

func TestScenarioHandler_ListScenarios(t *testing.T) {
	scenarioRepo := newTestScenarioRepo()
	scenarioRepo.scenarios["s1"] = &entity.Scenario{
		ID:   "s1",
		Name: "Test Scenario",
		Phases: []entity.Phase{
			{Name: "Phase 1", Techniques: []string{"T1082"}, Order: 1},
		},
	}
	techRepo := newTestTechniqueRepo()
	svc := createTestScenarioService(scenarioRepo, techRepo)
	handler := NewScenarioHandler(svc)

	router := gin.New()
	router.GET("/scenarios", handler.ListScenarios)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/scenarios", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestScenarioHandler_ListScenarios_Error(t *testing.T) {
	scenarioRepo := newTestScenarioRepo()
	scenarioRepo.err = errors.New("database error")
	techRepo := newTestTechniqueRepo()
	svc := createTestScenarioService(scenarioRepo, techRepo)
	handler := NewScenarioHandler(svc)

	router := gin.New()
	router.GET("/scenarios", handler.ListScenarios)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/scenarios", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestScenarioHandler_GetScenario(t *testing.T) {
	scenarioRepo := newTestScenarioRepo()
	scenarioRepo.scenarios["s1"] = &entity.Scenario{
		ID:   "s1",
		Name: "Test Scenario",
		Phases: []entity.Phase{
			{Name: "Phase 1", Techniques: []string{"T1082"}, Order: 1},
		},
	}
	techRepo := newTestTechniqueRepo()
	svc := createTestScenarioService(scenarioRepo, techRepo)
	handler := NewScenarioHandler(svc)

	router := gin.New()
	router.GET("/scenarios/:id", handler.GetScenario)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/scenarios/s1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var scenario entity.Scenario
	if err := json.NewDecoder(w.Body).Decode(&scenario); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}
	if scenario.ID != "s1" {
		t.Errorf("Expected scenario ID s1, got %s", scenario.ID)
	}
}

func TestScenarioHandler_GetScenario_NotFound(t *testing.T) {
	scenarioRepo := newTestScenarioRepo()
	techRepo := newTestTechniqueRepo()
	svc := createTestScenarioService(scenarioRepo, techRepo)
	handler := NewScenarioHandler(svc)

	router := gin.New()
	router.GET("/scenarios/:id", handler.GetScenario)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/scenarios/nonexistent", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestScenarioHandler_GetScenariosByTag(t *testing.T) {
	scenarioRepo := newTestScenarioRepo()
	scenarioRepo.scenarios["s1"] = &entity.Scenario{
		ID:   "s1",
		Name: "Test Scenario",
		Tags: []string{"apt29", "discovery"},
		Phases: []entity.Phase{
			{Name: "Phase 1", Techniques: []string{"T1082"}, Order: 1},
		},
	}
	scenarioRepo.scenarios["s2"] = &entity.Scenario{
		ID:   "s2",
		Name: "Another Scenario",
		Tags: []string{"execution"},
		Phases: []entity.Phase{
			{Name: "Phase 1", Techniques: []string{"T1059"}, Order: 1},
		},
	}
	techRepo := newTestTechniqueRepo()
	svc := createTestScenarioService(scenarioRepo, techRepo)
	handler := NewScenarioHandler(svc)

	router := gin.New()
	router.GET("/scenarios/tag/:tag", handler.GetScenariosByTag)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/scenarios/tag/apt29", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var scenarios []*entity.Scenario
	if err := json.NewDecoder(w.Body).Decode(&scenarios); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}
	if len(scenarios) != 1 {
		t.Errorf("Expected 1 scenario, got %d", len(scenarios))
	}
}

func TestScenarioHandler_CreateScenario(t *testing.T) {
	scenarioRepo := newTestScenarioRepo()
	techRepo := newTestTechniqueRepo()
	techRepo.techniques["T1082"] = &entity.Technique{
		ID:        "T1082",
		Name:      "System Information Discovery",
		Tactic:    entity.TacticDiscovery,
		Platforms: []string{"windows", "linux"},
	}
	svc := createTestScenarioService(scenarioRepo, techRepo)
	handler := NewScenarioHandler(svc)

	router := gin.New()
	router.POST("/scenarios", handler.CreateScenario)

	body := CreateScenarioRequest{
		Name:        "New Scenario",
		Description: "A test scenario",
		Phases: []entity.Phase{
			{Name: "Discovery", Techniques: []string{"T1082"}, Order: 1},
		},
		Tags:   []string{"test"},
		Author: "test-user",
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/scenarios", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestScenarioHandler_CreateScenario_BadRequest(t *testing.T) {
	scenarioRepo := newTestScenarioRepo()
	techRepo := newTestTechniqueRepo()
	svc := createTestScenarioService(scenarioRepo, techRepo)
	handler := NewScenarioHandler(svc)

	router := gin.New()
	router.POST("/scenarios", handler.CreateScenario)

	// Missing required fields
	body := `{"description": "missing name and phases"}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/scenarios", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestScenarioHandler_UpdateScenario(t *testing.T) {
	scenarioRepo := newTestScenarioRepo()
	scenarioRepo.scenarios["s1"] = &entity.Scenario{
		ID:        "s1",
		Name:      "Original Name",
		CreatedAt: time.Now(),
		Phases: []entity.Phase{
			{Name: "Phase 1", Techniques: []string{"T1082"}, Order: 1},
		},
	}
	techRepo := newTestTechniqueRepo()
	techRepo.techniques["T1082"] = &entity.Technique{
		ID:        "T1082",
		Name:      "System Information Discovery",
		Tactic:    entity.TacticDiscovery,
		Platforms: []string{"windows", "linux"},
	}
	svc := createTestScenarioService(scenarioRepo, techRepo)
	handler := NewScenarioHandler(svc)

	router := gin.New()
	router.PUT("/scenarios/:id", handler.UpdateScenario)

	body := UpdateScenarioRequest{
		Name:        "Updated Name",
		Description: "Updated description",
		Phases: []entity.Phase{
			{Name: "Phase 1", Techniques: []string{"T1082"}, Order: 1},
		},
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/scenarios/s1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestScenarioHandler_UpdateScenario_NotFound(t *testing.T) {
	scenarioRepo := newTestScenarioRepo()
	techRepo := newTestTechniqueRepo()
	svc := createTestScenarioService(scenarioRepo, techRepo)
	handler := NewScenarioHandler(svc)

	router := gin.New()
	router.PUT("/scenarios/:id", handler.UpdateScenario)

	body := UpdateScenarioRequest{
		Name: "Updated Name",
		Phases: []entity.Phase{
			{Name: "Phase 1", Techniques: []string{"T1082"}, Order: 1},
		},
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/scenarios/nonexistent", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestScenarioHandler_DeleteScenario(t *testing.T) {
	scenarioRepo := newTestScenarioRepo()
	scenarioRepo.scenarios["s1"] = &entity.Scenario{
		ID:   "s1",
		Name: "To Delete",
	}
	techRepo := newTestTechniqueRepo()
	svc := createTestScenarioService(scenarioRepo, techRepo)
	handler := NewScenarioHandler(svc)

	router := gin.New()
	router.DELETE("/scenarios/:id", handler.DeleteScenario)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/scenarios/s1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}
}

func TestScenarioHandler_DeleteScenario_Error(t *testing.T) {
	scenarioRepo := newTestScenarioRepo()
	scenarioRepo.err = errors.New("database error")
	techRepo := newTestTechniqueRepo()
	svc := createTestScenarioService(scenarioRepo, techRepo)
	handler := NewScenarioHandler(svc)

	router := gin.New()
	router.DELETE("/scenarios/:id", handler.DeleteScenario)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/scenarios/s1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}
