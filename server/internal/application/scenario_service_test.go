package application

import (
	"context"
	"errors"
	"testing"

	"autostrike/internal/domain/entity"
	"autostrike/internal/domain/service"
)

type mockScenarioRepo struct {
	scenarios map[string]*entity.Scenario
	err       error
}

func newMockScenarioRepo() *mockScenarioRepo {
	return &mockScenarioRepo{scenarios: make(map[string]*entity.Scenario)}
}

func (m *mockScenarioRepo) Create(ctx context.Context, s *entity.Scenario) error {
	if m.err != nil {
		return m.err
	}
	m.scenarios[s.ID] = s
	return nil
}

func (m *mockScenarioRepo) Update(ctx context.Context, s *entity.Scenario) error {
	if m.err != nil {
		return m.err
	}
	m.scenarios[s.ID] = s
	return nil
}

func (m *mockScenarioRepo) Delete(ctx context.Context, id string) error {
	if m.err != nil {
		return m.err
	}
	delete(m.scenarios, id)
	return nil
}

func (m *mockScenarioRepo) FindByID(ctx context.Context, id string) (*entity.Scenario, error) {
	if m.err != nil {
		return nil, m.err
	}
	s, ok := m.scenarios[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return s, nil
}

func (m *mockScenarioRepo) FindAll(ctx context.Context) ([]*entity.Scenario, error) {
	if m.err != nil {
		return nil, m.err
	}
	result := make([]*entity.Scenario, 0, len(m.scenarios))
	for _, s := range m.scenarios {
		result = append(result, s)
	}
	return result, nil
}

func (m *mockScenarioRepo) FindByTag(ctx context.Context, tag string) ([]*entity.Scenario, error) {
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

func (m *mockScenarioRepo) ImportFromYAML(ctx context.Context, path string) error {
	return m.err
}

func TestNewScenarioService(t *testing.T) {
	scenarioRepo := newMockScenarioRepo()
	techRepo := newMockTechniqueRepo()
	validator := service.NewTechniqueValidator()

	svc := NewScenarioService(scenarioRepo, techRepo, validator)
	if svc == nil {
		t.Fatal("Expected non-nil service")
	}
}

func TestGetScenario(t *testing.T) {
	scenarioRepo := newMockScenarioRepo()
	scenarioRepo.scenarios["s1"] = &entity.Scenario{ID: "s1", Name: "Test"}
	techRepo := newMockTechniqueRepo()
	validator := service.NewTechniqueValidator()

	svc := NewScenarioService(scenarioRepo, techRepo, validator)
	scenario, err := svc.GetScenario(context.Background(), "s1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if scenario.Name != "Test" {
		t.Errorf("Expected name Test, got %s", scenario.Name)
	}
}

func TestGetScenario_NotFound(t *testing.T) {
	scenarioRepo := newMockScenarioRepo()
	techRepo := newMockTechniqueRepo()
	validator := service.NewTechniqueValidator()

	svc := NewScenarioService(scenarioRepo, techRepo, validator)
	_, err := svc.GetScenario(context.Background(), "invalid")
	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestGetAllScenarios(t *testing.T) {
	scenarioRepo := newMockScenarioRepo()
	scenarioRepo.scenarios["s1"] = &entity.Scenario{ID: "s1"}
	scenarioRepo.scenarios["s2"] = &entity.Scenario{ID: "s2"}
	techRepo := newMockTechniqueRepo()
	validator := service.NewTechniqueValidator()

	svc := NewScenarioService(scenarioRepo, techRepo, validator)
	scenarios, err := svc.GetAllScenarios(context.Background())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(scenarios) != 2 {
		t.Errorf("Expected 2 scenarios, got %d", len(scenarios))
	}
}

func TestGetScenariosByTag(t *testing.T) {
	scenarioRepo := newMockScenarioRepo()
	scenarioRepo.scenarios["s1"] = &entity.Scenario{ID: "s1", Tags: []string{"apt"}}
	scenarioRepo.scenarios["s2"] = &entity.Scenario{ID: "s2", Tags: []string{"ransomware"}}
	techRepo := newMockTechniqueRepo()
	validator := service.NewTechniqueValidator()

	svc := NewScenarioService(scenarioRepo, techRepo, validator)
	scenarios, err := svc.GetScenariosByTag(context.Background(), "apt")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(scenarios) != 1 {
		t.Errorf("Expected 1 scenario, got %d", len(scenarios))
	}
}

func TestCreateScenario(t *testing.T) {
	scenarioRepo := newMockScenarioRepo()
	techRepo := newMockTechniqueRepo()
	techRepo.techniques["T1059"] = &entity.Technique{ID: "T1059"}
	validator := service.NewTechniqueValidator()

	svc := NewScenarioService(scenarioRepo, techRepo, validator)
	scenario := &entity.Scenario{
		Name: "Test Scenario",
		Phases: []entity.Phase{
			{Name: "Phase 1", Techniques: []entity.TechniqueSelection{{TechniqueID: "T1059"}}},
		},
	}

	err := svc.CreateScenario(context.Background(), scenario)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if scenario.ID == "" {
		t.Error("Expected ID to be set")
	}
	if scenario.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
}

func TestCreateScenario_TechniqueRepoError(t *testing.T) {
	scenarioRepo := newMockScenarioRepo()
	techRepo := newMockTechniqueRepo()
	techRepo.err = errors.New("db error")
	validator := service.NewTechniqueValidator()

	svc := NewScenarioService(scenarioRepo, techRepo, validator)
	scenario := &entity.Scenario{Name: "Test"}

	err := svc.CreateScenario(context.Background(), scenario)
	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestCreateScenario_ValidationError(t *testing.T) {
	scenarioRepo := newMockScenarioRepo()
	techRepo := newMockTechniqueRepo()
	validator := service.NewTechniqueValidator()

	svc := NewScenarioService(scenarioRepo, techRepo, validator)
	scenario := &entity.Scenario{
		Name:   "Test Scenario",
		Phases: []entity.Phase{}, // Empty phases should fail validation
	}

	err := svc.CreateScenario(context.Background(), scenario)
	if err == nil {
		t.Fatal("Expected validation error")
	}
}

func TestUpdateScenario(t *testing.T) {
	scenarioRepo := newMockScenarioRepo()
	scenarioRepo.scenarios["s1"] = &entity.Scenario{ID: "s1", Name: "Old"}
	techRepo := newMockTechniqueRepo()
	techRepo.techniques["T1059"] = &entity.Technique{ID: "T1059"}
	validator := service.NewTechniqueValidator()

	svc := NewScenarioService(scenarioRepo, techRepo, validator)
	scenario := &entity.Scenario{
		ID:   "s1",
		Name: "Updated",
		Phases: []entity.Phase{
			{Name: "Phase 1", Techniques: []entity.TechniqueSelection{{TechniqueID: "T1059"}}},
		},
	}

	err := svc.UpdateScenario(context.Background(), scenario)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if scenario.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}
}

func TestUpdateScenario_TechniqueRepoError(t *testing.T) {
	scenarioRepo := newMockScenarioRepo()
	techRepo := newMockTechniqueRepo()
	techRepo.err = errors.New("db error")
	validator := service.NewTechniqueValidator()

	svc := NewScenarioService(scenarioRepo, techRepo, validator)
	scenario := &entity.Scenario{ID: "s1", Name: "Test"}

	err := svc.UpdateScenario(context.Background(), scenario)
	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestUpdateScenario_ValidationError(t *testing.T) {
	scenarioRepo := newMockScenarioRepo()
	techRepo := newMockTechniqueRepo()
	validator := service.NewTechniqueValidator()

	svc := NewScenarioService(scenarioRepo, techRepo, validator)
	scenario := &entity.Scenario{
		ID:     "s1",
		Name:   "Test",
		Phases: []entity.Phase{},
	}

	err := svc.UpdateScenario(context.Background(), scenario)
	if err == nil {
		t.Fatal("Expected validation error")
	}
}

func TestDeleteScenario(t *testing.T) {
	scenarioRepo := newMockScenarioRepo()
	scenarioRepo.scenarios["s1"] = &entity.Scenario{ID: "s1"}
	techRepo := newMockTechniqueRepo()
	validator := service.NewTechniqueValidator()

	svc := NewScenarioService(scenarioRepo, techRepo, validator)
	err := svc.DeleteScenario(context.Background(), "s1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestValidationError_Error(t *testing.T) {
	err := &ValidationError{Errors: []string{"error1", "error2"}}
	if err.Error() != "error1" {
		t.Errorf("Expected 'error1', got '%s'", err.Error())
	}
}

func TestValidationError_Empty(t *testing.T) {
	err := &ValidationError{Errors: []string{}}
	if err.Error() != "validation failed" {
		t.Errorf("Expected 'validation failed', got '%s'", err.Error())
	}
}

func TestImportScenarios(t *testing.T) {
	scenarioRepo := newMockScenarioRepo()
	techRepo := newMockTechniqueRepo()
	validator := service.NewTechniqueValidator()

	svc := NewScenarioService(scenarioRepo, techRepo, validator)
	err := svc.ImportScenarios(context.Background(), "/path/to/scenarios.yaml")
	// Mock repo returns nil error
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestImportScenarios_Error(t *testing.T) {
	scenarioRepo := newMockScenarioRepo()
	scenarioRepo.err = errors.New("import error")
	techRepo := newMockTechniqueRepo()
	validator := service.NewTechniqueValidator()

	svc := NewScenarioService(scenarioRepo, techRepo, validator)
	err := svc.ImportScenarios(context.Background(), "/path/to/scenarios.yaml")
	if err == nil {
		t.Fatal("Expected error")
	}
}
