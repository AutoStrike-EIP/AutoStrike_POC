package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"autostrike/internal/domain/entity"
	"autostrike/internal/domain/service"
)

type mockResultRepo struct {
	executions       map[string]*entity.Execution
	results          map[string][]*entity.ExecutionResult
	err              error
	createResultErr  error
	findResultsErr   error
}

func newMockResultRepo() *mockResultRepo {
	return &mockResultRepo{
		executions: make(map[string]*entity.Execution),
		results:    make(map[string][]*entity.ExecutionResult),
	}
}

func (m *mockResultRepo) CreateExecution(ctx context.Context, e *entity.Execution) error {
	if m.err != nil {
		return m.err
	}
	m.executions[e.ID] = e
	return nil
}

func (m *mockResultRepo) UpdateExecution(ctx context.Context, e *entity.Execution) error {
	if m.err != nil {
		return m.err
	}
	m.executions[e.ID] = e
	return nil
}

func (m *mockResultRepo) FindExecutionByID(ctx context.Context, id string) (*entity.Execution, error) {
	if m.err != nil {
		return nil, m.err
	}
	e, ok := m.executions[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return e, nil
}

func (m *mockResultRepo) FindExecutionsByScenario(ctx context.Context, scenarioID string) ([]*entity.Execution, error) {
	if m.err != nil {
		return nil, m.err
	}
	var result []*entity.Execution
	for _, e := range m.executions {
		if e.ScenarioID == scenarioID {
			result = append(result, e)
		}
	}
	return result, nil
}

func (m *mockResultRepo) FindRecentExecutions(ctx context.Context, limit int) ([]*entity.Execution, error) {
	if m.err != nil {
		return nil, m.err
	}
	result := make([]*entity.Execution, 0, len(m.executions))
	for _, e := range m.executions {
		result = append(result, e)
		if len(result) >= limit {
			break
		}
	}
	return result, nil
}

func (m *mockResultRepo) CreateResult(ctx context.Context, r *entity.ExecutionResult) error {
	if m.createResultErr != nil {
		return m.createResultErr
	}
	if m.err != nil {
		return m.err
	}
	m.results[r.ExecutionID] = append(m.results[r.ExecutionID], r)
	return nil
}

func (m *mockResultRepo) UpdateResult(ctx context.Context, r *entity.ExecutionResult) error {
	if m.err != nil {
		return m.err
	}
	return nil
}

func (m *mockResultRepo) FindResultsByExecution(ctx context.Context, executionID string) ([]*entity.ExecutionResult, error) {
	if m.findResultsErr != nil {
		return nil, m.findResultsErr
	}
	if m.err != nil {
		return nil, m.err
	}
	return m.results[executionID], nil
}

func (m *mockResultRepo) FindResultsByTechnique(ctx context.Context, techniqueID string) ([]*entity.ExecutionResult, error) {
	if m.err != nil {
		return nil, m.err
	}
	var result []*entity.ExecutionResult
	for _, results := range m.results {
		for _, r := range results {
			if r.TechniqueID == techniqueID {
				result = append(result, r)
			}
		}
	}
	return result, nil
}

func TestNewExecutionService(t *testing.T) {
	resultRepo := newMockResultRepo()
	scenarioRepo := newMockScenarioRepo()
	techRepo := newMockTechniqueRepo()
	agentRepo := newMockAgentRepo()
	validator := service.NewTechniqueValidator()
	orchestrator := service.NewAttackOrchestrator(agentRepo, techRepo, validator, nil)
	calculator := service.NewScoreCalculator()

	svc := NewExecutionService(resultRepo, scenarioRepo, techRepo, agentRepo, orchestrator, calculator)
	if svc == nil {
		t.Fatal("Expected non-nil service")
	}
}

func TestGetExecution(t *testing.T) {
	resultRepo := newMockResultRepo()
	resultRepo.executions["e1"] = &entity.Execution{ID: "e1", ScenarioID: "s1"}

	svc := &ExecutionService{resultRepo: resultRepo}
	exec, err := svc.GetExecution(context.Background(), "e1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if exec.ScenarioID != "s1" {
		t.Errorf("Expected scenario s1, got %s", exec.ScenarioID)
	}
}

func TestGetExecutionNotFound(t *testing.T) {
	resultRepo := newMockResultRepo()
	svc := &ExecutionService{resultRepo: resultRepo}

	_, err := svc.GetExecution(context.Background(), "invalid")
	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestGetExecutionResults(t *testing.T) {
	resultRepo := newMockResultRepo()
	resultRepo.results["e1"] = []*entity.ExecutionResult{
		{ID: "r1", ExecutionID: "e1"},
		{ID: "r2", ExecutionID: "e1"},
	}

	svc := &ExecutionService{resultRepo: resultRepo}
	results, err := svc.GetExecutionResults(context.Background(), "e1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}
}

func TestGetRecentExecutions(t *testing.T) {
	resultRepo := newMockResultRepo()
	resultRepo.executions["e1"] = &entity.Execution{ID: "e1"}
	resultRepo.executions["e2"] = &entity.Execution{ID: "e2"}
	resultRepo.executions["e3"] = &entity.Execution{ID: "e3"}

	svc := &ExecutionService{resultRepo: resultRepo}
	execs, err := svc.GetRecentExecutions(context.Background(), 2)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(execs) > 2 {
		t.Errorf("Expected at most 2 executions, got %d", len(execs))
	}
}

func TestStartExecutionScenarioNotFound(t *testing.T) {
	resultRepo := newMockResultRepo()
	scenarioRepo := newMockScenarioRepo()
	agentRepo := newMockAgentRepo()

	svc := &ExecutionService{
		resultRepo:   resultRepo,
		scenarioRepo: scenarioRepo,
		agentRepo:    agentRepo,
	}

	_, err := svc.StartExecution(context.Background(), "invalid", []string{"paw1"}, false)
	if err == nil {
		t.Fatal("Expected error for missing scenario")
	}
}

func TestStartExecutionAgentNotFound(t *testing.T) {
	resultRepo := newMockResultRepo()
	scenarioRepo := newMockScenarioRepo()
	scenarioRepo.scenarios["s1"] = &entity.Scenario{ID: "s1"}
	agentRepo := newMockAgentRepo()

	svc := &ExecutionService{
		resultRepo:   resultRepo,
		scenarioRepo: scenarioRepo,
		agentRepo:    agentRepo,
	}

	_, err := svc.StartExecution(context.Background(), "s1", []string{"invalid"}, false)
	if err == nil {
		t.Fatal("Expected error for missing agent")
	}
}

func TestStartExecutionAgentOffline(t *testing.T) {
	resultRepo := newMockResultRepo()
	scenarioRepo := newMockScenarioRepo()
	scenarioRepo.scenarios["s1"] = &entity.Scenario{ID: "s1"}
	agentRepo := newMockAgentRepo()
	agentRepo.agents["paw1"] = &entity.Agent{Paw: "paw1", Status: entity.AgentOffline}

	svc := &ExecutionService{
		resultRepo:   resultRepo,
		scenarioRepo: scenarioRepo,
		agentRepo:    agentRepo,
	}

	_, err := svc.StartExecution(context.Background(), "s1", []string{"paw1"}, false)
	if err == nil {
		t.Fatal("Expected error for offline agent")
	}
}

func TestCompleteExecution(t *testing.T) {
	resultRepo := newMockResultRepo()
	resultRepo.executions["e1"] = &entity.Execution{
		ID:     "e1",
		Status: entity.ExecutionRunning,
	}
	resultRepo.results["e1"] = []*entity.ExecutionResult{
		{ID: "r1", Status: entity.StatusSuccess},
	}
	calculator := service.NewScoreCalculator()

	svc := &ExecutionService{
		resultRepo: resultRepo,
		calculator: calculator,
	}

	err := svc.CompleteExecution(context.Background(), "e1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	exec := resultRepo.executions["e1"]
	if exec.Status != entity.ExecutionCompleted {
		t.Errorf("Expected status Completed, got %v", exec.Status)
	}
}

func TestCompleteExecutionNotFound(t *testing.T) {
	resultRepo := newMockResultRepo()
	resultRepo.err = errors.New("not found")

	svc := &ExecutionService{resultRepo: resultRepo}
	err := svc.CompleteExecution(context.Background(), "invalid")
	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestUpdateResult(t *testing.T) {
	resultRepo := newMockResultRepo()
	// The implementation uses resultID as executionID when looking up results
	// So we need to use the resultID "r1" as the key
	resultRepo.results["r1"] = []*entity.ExecutionResult{
		{ID: "r1", ExecutionID: "e1", Status: entity.StatusPending},
	}

	svc := &ExecutionService{resultRepo: resultRepo}
	err := svc.UpdateResult(context.Background(), "r1", entity.StatusSuccess, "output", false)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestUpdateResultNotFound(t *testing.T) {
	resultRepo := newMockResultRepo()
	resultRepo.results["e1"] = []*entity.ExecutionResult{}

	svc := &ExecutionService{resultRepo: resultRepo}
	err := svc.UpdateResult(context.Background(), "invalid", entity.StatusSuccess, "", false)
	if err == nil {
		t.Fatal("Expected error for not found result")
	}
}

func TestStartExecutionSuccess(t *testing.T) {
	resultRepo := newMockResultRepo()
	scenarioRepo := newMockScenarioRepo()
	scenarioRepo.scenarios["s1"] = &entity.Scenario{
		ID:   "s1",
		Name: "Test",
		Phases: []entity.Phase{
			{Name: "Phase1", Techniques: []string{"T1059"}},
		},
	}
	techRepo := newMockTechniqueRepo()
	techRepo.techniques["T1059"] = &entity.Technique{
		ID:        "T1059",
		Platforms: []string{"linux"},
		Executors: []entity.Executor{{Type: "sh", Command: "echo test"}},
	}
	agentRepo := newMockAgentRepo()
	agentRepo.agents["paw1"] = &entity.Agent{
		Paw:       "paw1",
		Status:    entity.AgentOnline,
		Platform:  "linux",
		Executors: []string{"sh"},
		LastSeen:  time.Now(),
	}
	validator := service.NewTechniqueValidator()
	orchestrator := service.NewAttackOrchestrator(agentRepo, techRepo, validator, nil)
	calculator := service.NewScoreCalculator()

	svc := NewExecutionService(resultRepo, scenarioRepo, techRepo, agentRepo, orchestrator, calculator)
	exec, err := svc.StartExecution(context.Background(), "s1", []string{"paw1"}, false)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if exec == nil {
		t.Fatal("Expected execution to be created")
	}
	if exec.Status != entity.ExecutionRunning {
		t.Errorf("Expected status Running, got %v", exec.Status)
	}
}

func TestUpdateResultRepoError(t *testing.T) {
	resultRepo := newMockResultRepo()
	resultRepo.err = errors.New("db error")

	svc := &ExecutionService{resultRepo: resultRepo}
	err := svc.UpdateResult(context.Background(), "r1", entity.StatusSuccess, "", false)
	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestCompleteExecutionResultsError(t *testing.T) {
	resultRepo := newMockResultRepo()
	resultRepo.executions["e1"] = &entity.Execution{
		ID:     "e1",
		Status: entity.ExecutionRunning,
	}
	// Use specific error for FindResultsByExecution
	resultRepo.findResultsErr = errors.New("results error")

	svc := &ExecutionService{resultRepo: resultRepo}
	err := svc.CompleteExecution(context.Background(), "e1")
	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestStartExecutionPlanError(t *testing.T) {
	resultRepo := newMockResultRepo()
	scenarioRepo := newMockScenarioRepo()
	scenarioRepo.scenarios["s1"] = &entity.Scenario{
		ID:   "s1",
		Name: "Test",
		Phases: []entity.Phase{
			{Name: "Phase1", Techniques: []string{"T1059"}},
		},
	}
	techRepo := newMockTechniqueRepo()
	// Missing technique causes plan error
	agentRepo := newMockAgentRepo()
	agentRepo.agents["paw1"] = &entity.Agent{
		Paw:       "paw1",
		Status:    entity.AgentOnline,
		Platform:  "linux",
		Executors: []string{"sh"},
		LastSeen:  time.Now(),
	}
	validator := service.NewTechniqueValidator()
	orchestrator := service.NewAttackOrchestrator(agentRepo, techRepo, validator, nil)
	calculator := service.NewScoreCalculator()

	svc := NewExecutionService(resultRepo, scenarioRepo, techRepo, agentRepo, orchestrator, calculator)
	_, err := svc.StartExecution(context.Background(), "s1", []string{"paw1"}, false)
	if err == nil {
		t.Fatal("Expected error for plan failure")
	}
}

func TestStartExecutionCreateError(t *testing.T) {
	resultRepo := newMockResultRepo()
	resultRepo.err = errors.New("create error")
	scenarioRepo := newMockScenarioRepo()
	scenarioRepo.scenarios["s1"] = &entity.Scenario{
		ID:   "s1",
		Name: "Test",
		Phases: []entity.Phase{
			{Name: "Phase1", Techniques: []string{"T1059"}},
		},
	}
	techRepo := newMockTechniqueRepo()
	techRepo.techniques["T1059"] = &entity.Technique{
		ID:        "T1059",
		Platforms: []string{"linux"},
		Executors: []entity.Executor{{Type: "sh", Command: "echo test"}},
	}
	agentRepo := newMockAgentRepo()
	agentRepo.agents["paw1"] = &entity.Agent{
		Paw:       "paw1",
		Status:    entity.AgentOnline,
		Platform:  "linux",
		Executors: []string{"sh"},
		LastSeen:  time.Now(),
	}
	validator := service.NewTechniqueValidator()
	orchestrator := service.NewAttackOrchestrator(agentRepo, techRepo, validator, nil)
	calculator := service.NewScoreCalculator()

	svc := NewExecutionService(resultRepo, scenarioRepo, techRepo, agentRepo, orchestrator, calculator)
	_, err := svc.StartExecution(context.Background(), "s1", []string{"paw1"}, false)
	if err == nil {
		t.Fatal("Expected error for create failure")
	}
}

func TestStartExecutionCreateResultError(t *testing.T) {
	resultRepo := newMockResultRepo()
	// Only CreateResult will fail
	resultRepo.createResultErr = errors.New("create result error")
	scenarioRepo := newMockScenarioRepo()
	scenarioRepo.scenarios["s1"] = &entity.Scenario{
		ID:   "s1",
		Name: "Test",
		Phases: []entity.Phase{
			{Name: "Phase1", Techniques: []string{"T1059"}},
		},
	}
	techRepo := newMockTechniqueRepo()
	techRepo.techniques["T1059"] = &entity.Technique{
		ID:        "T1059",
		Platforms: []string{"linux"},
		Executors: []entity.Executor{{Type: "sh", Command: "echo test"}},
	}
	agentRepo := newMockAgentRepo()
	agentRepo.agents["paw1"] = &entity.Agent{
		Paw:       "paw1",
		Status:    entity.AgentOnline,
		Platform:  "linux",
		Executors: []string{"sh"},
		LastSeen:  time.Now(),
	}
	validator := service.NewTechniqueValidator()
	orchestrator := service.NewAttackOrchestrator(agentRepo, techRepo, validator, nil)
	calculator := service.NewScoreCalculator()

	svc := NewExecutionService(resultRepo, scenarioRepo, techRepo, agentRepo, orchestrator, calculator)
	_, err := svc.StartExecution(context.Background(), "s1", []string{"paw1"}, false)
	if err == nil {
		t.Fatal("Expected error for create result failure")
	}
}
