package service

import (
	"context"
	"errors"
	"testing"

	"autostrike/internal/domain/entity"
)

// Mock repositories for testing
type mockAgentRepo struct {
	agents map[string]*entity.Agent
}

func (m *mockAgentRepo) Create(ctx context.Context, agent *entity.Agent) error {
	return nil
}
func (m *mockAgentRepo) Update(ctx context.Context, agent *entity.Agent) error {
	return nil
}
func (m *mockAgentRepo) Delete(ctx context.Context, paw string) error { return nil }
func (m *mockAgentRepo) FindByPaw(ctx context.Context, paw string) (*entity.Agent, error) {
	if agent, ok := m.agents[paw]; ok {
		return agent, nil
	}
	return nil, errors.New("agent not found")
}
func (m *mockAgentRepo) FindByPaws(ctx context.Context, paws []string) ([]*entity.Agent, error) {
	var result []*entity.Agent
	for _, paw := range paws {
		if agent, ok := m.agents[paw]; ok {
			result = append(result, agent)
		}
	}
	return result, nil
}
func (m *mockAgentRepo) FindAll(ctx context.Context) ([]*entity.Agent, error) {
	return nil, nil
}
func (m *mockAgentRepo) FindByStatus(ctx context.Context, status entity.AgentStatus) ([]*entity.Agent, error) {
	return nil, nil
}
func (m *mockAgentRepo) FindByPlatform(ctx context.Context, platform string) ([]*entity.Agent, error) {
	return nil, nil
}
func (m *mockAgentRepo) UpdateLastSeen(ctx context.Context, paw string) error {
	return nil
}

type mockTechniqueRepo struct {
	techniques map[string]*entity.Technique
}

func (m *mockTechniqueRepo) Create(ctx context.Context, technique *entity.Technique) error {
	return nil
}
func (m *mockTechniqueRepo) Update(ctx context.Context, technique *entity.Technique) error {
	return nil
}
func (m *mockTechniqueRepo) Delete(ctx context.Context, id string) error { return nil }
func (m *mockTechniqueRepo) FindByID(ctx context.Context, id string) (*entity.Technique, error) {
	if tech, ok := m.techniques[id]; ok {
		return tech, nil
	}
	return nil, errors.New("technique not found")
}
func (m *mockTechniqueRepo) FindAll(ctx context.Context) ([]*entity.Technique, error) {
	return nil, nil
}
func (m *mockTechniqueRepo) FindByTactic(ctx context.Context, tactic entity.TacticType) ([]*entity.Technique, error) {
	return nil, nil
}
func (m *mockTechniqueRepo) FindByPlatform(ctx context.Context, platform string) ([]*entity.Technique, error) {
	return nil, nil
}
func (m *mockTechniqueRepo) ImportFromYAML(ctx context.Context, path string) error {
	return nil
}

func TestNewAttackOrchestrator(t *testing.T) {
	agentRepo := &mockAgentRepo{}
	techRepo := &mockTechniqueRepo{}
	validator := NewTechniqueValidator()

	orchestrator := NewAttackOrchestrator(agentRepo, techRepo, validator, nil)

	if orchestrator == nil {
		t.Error("NewAttackOrchestrator returned nil")
	}
}

func TestAttackOrchestrator_PlanExecution(t *testing.T) {
	technique := &entity.Technique{
		ID:        "T1059",
		Name:      "Command Execution",
		Platforms: []string{"windows"},
		Executors: []entity.Executor{
			{Type: "psh", Command: "whoami", Timeout: 30},
		},
		IsSafe: true,
	}

	techRepo := &mockTechniqueRepo{
		techniques: map[string]*entity.Technique{"T1059": technique},
	}

	agentRepo := &mockAgentRepo{}
	validator := NewTechniqueValidator()
	orchestrator := NewAttackOrchestrator(agentRepo, techRepo, validator, nil)

	agent := &entity.Agent{
		Paw:       "test-agent",
		Platform:  "windows",
		Executors: []string{"psh", "cmd"},
		Status:    entity.AgentOnline,
	}

	scenario := &entity.Scenario{
		ID:   "test-scenario",
		Name: "Test Scenario",
		Phases: []entity.Phase{
			{Name: "Phase1", Techniques: []entity.TechniqueSelection{{TechniqueID: "T1059"}}},
		},
	}

	plan, err := orchestrator.PlanExecution(context.Background(), scenario, []*entity.Agent{agent}, false)

	if err != nil {
		t.Errorf("PlanExecution returned error: %v", err)
	}
	if plan == nil {
		t.Fatal("PlanExecution returned nil plan")
	}
	if len(plan.Tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(plan.Tasks))
	}
	if plan.Tasks[0].TechniqueID != "T1059" {
		t.Errorf("Expected technique T1059, got %s", plan.Tasks[0].TechniqueID)
	}
}

func TestAttackOrchestrator_PlanExecution_SafeMode(t *testing.T) {
	safeTech := &entity.Technique{
		ID:        "T1082",
		Name:      "System Info",
		Platforms: []string{"windows"},
		Executors: []entity.Executor{
			{Type: "psh", Command: "systeminfo", Timeout: 30, IsSafe: true},
		},
		IsSafe: true,
	}

	unsafeTech := &entity.Technique{
		ID:        "T1055",
		Name:      "Process Injection",
		Platforms: []string{"windows"},
		Executors: []entity.Executor{
			{Type: "psh", Command: "dangerous", Timeout: 30, ElevationRequired: true, IsSafe: false},
		},
		IsSafe: false,
	}

	techRepo := &mockTechniqueRepo{
		techniques: map[string]*entity.Technique{
			"T1082": safeTech,
			"T1055": unsafeTech,
		},
	}

	agentRepo := &mockAgentRepo{}
	validator := NewTechniqueValidator()
	orchestrator := NewAttackOrchestrator(agentRepo, techRepo, validator, nil)

	agent := &entity.Agent{
		Paw:       "test-agent",
		Platform:  "windows",
		Executors: []string{"psh"},
		Status:    entity.AgentOnline,
	}

	scenario := &entity.Scenario{
		ID:   "test-scenario",
		Name: "Test Scenario",
		Phases: []entity.Phase{
			{Name: "Phase1", Techniques: []entity.TechniqueSelection{{TechniqueID: "T1082"}, {TechniqueID: "T1055"}}},
		},
	}

	// Safe mode should skip unsafe techniques
	plan, err := orchestrator.PlanExecution(context.Background(), scenario, []*entity.Agent{agent}, true)

	if err != nil {
		t.Errorf("PlanExecution returned error: %v", err)
	}
	if plan == nil {
		t.Fatal("PlanExecution returned nil plan")
	}
	if len(plan.Tasks) != 1 {
		t.Errorf("Expected 1 task (safe only), got %d", len(plan.Tasks))
	}
	if plan.Tasks[0].TechniqueID != "T1082" {
		t.Errorf("Expected safe technique T1082, got %s", plan.Tasks[0].TechniqueID)
	}
}

func TestAttackOrchestrator_PlanExecution_NoCompatibleAgents(t *testing.T) {
	technique := &entity.Technique{
		ID:        "T1059",
		Platforms: []string{"windows"},
		Executors: []entity.Executor{
			{Type: "psh", Command: "whoami", Timeout: 30},
		},
		IsSafe: true,
	}

	techRepo := &mockTechniqueRepo{
		techniques: map[string]*entity.Technique{"T1059": technique},
	}

	agentRepo := &mockAgentRepo{}
	validator := NewTechniqueValidator()
	orchestrator := NewAttackOrchestrator(agentRepo, techRepo, validator, nil)

	// Linux agent cannot run Windows technique
	agent := &entity.Agent{
		Paw:       "linux-agent",
		Platform:  "linux",
		Executors: []string{"bash"},
		Status:    entity.AgentOnline,
	}

	scenario := &entity.Scenario{
		ID:   "test-scenario",
		Name: "Test Scenario",
		Phases: []entity.Phase{
			{Name: "Phase1", Techniques: []entity.TechniqueSelection{{TechniqueID: "T1059"}}},
		},
	}

	_, err := orchestrator.PlanExecution(context.Background(), scenario, []*entity.Agent{agent}, false)

	if err == nil {
		t.Error("Expected error for no compatible agents, got nil")
	}
}

func TestAttackOrchestrator_PlanExecution_TechniqueNotFound(t *testing.T) {
	techRepo := &mockTechniqueRepo{
		techniques: map[string]*entity.Technique{}, // empty
	}

	agentRepo := &mockAgentRepo{}
	validator := NewTechniqueValidator()
	orchestrator := NewAttackOrchestrator(agentRepo, techRepo, validator, nil)

	agent := &entity.Agent{
		Paw:       "test-agent",
		Platform:  "windows",
		Executors: []string{"psh"},
		Status:    entity.AgentOnline,
	}

	scenario := &entity.Scenario{
		ID:   "test-scenario",
		Name: "Test Scenario",
		Phases: []entity.Phase{
			{Name: "Phase1", Techniques: []entity.TechniqueSelection{{TechniqueID: "T9999"}}}, // doesn't exist
		},
	}

	_, err := orchestrator.PlanExecution(context.Background(), scenario, []*entity.Agent{agent}, false)

	if err == nil {
		t.Error("Expected error for technique not found")
	}
}

func TestAttackOrchestrator_ValidatePlan(t *testing.T) {
	agent := &entity.Agent{
		Paw:       "test-agent",
		Platform:  "windows",
		Executors: []string{"psh"},
		Status:    entity.AgentOnline,
	}

	technique := &entity.Technique{
		ID:        "T1059",
		Platforms: []string{"windows"},
		Executors: []entity.Executor{
			{Type: "psh", Command: "whoami"},
		},
	}

	agentRepo := &mockAgentRepo{
		agents: map[string]*entity.Agent{"test-agent": agent},
	}

	techRepo := &mockTechniqueRepo{
		techniques: map[string]*entity.Technique{"T1059": technique},
	}

	validator := NewTechniqueValidator()
	orchestrator := NewAttackOrchestrator(agentRepo, techRepo, validator, nil)

	plan := &ExecutionPlan{
		ID: "test-plan",
		Tasks: []PlannedTask{
			{TechniqueID: "T1059", AgentPaw: "test-agent"},
		},
	}

	err := orchestrator.ValidatePlan(context.Background(), plan)

	if err != nil {
		t.Errorf("ValidatePlan returned error: %v", err)
	}
}

func TestAttackOrchestrator_ValidatePlan_AgentNotFound(t *testing.T) {
	agentRepo := &mockAgentRepo{
		agents: map[string]*entity.Agent{}, // empty
	}

	techRepo := &mockTechniqueRepo{}
	validator := NewTechniqueValidator()
	orchestrator := NewAttackOrchestrator(agentRepo, techRepo, validator, nil)

	plan := &ExecutionPlan{
		ID: "test-plan",
		Tasks: []PlannedTask{
			{TechniqueID: "T1059", AgentPaw: "missing-agent"},
		},
	}

	err := orchestrator.ValidatePlan(context.Background(), plan)

	if err == nil {
		t.Error("Expected error for missing agent")
	}
}

func TestAttackOrchestrator_ValidatePlan_AgentNotOnline(t *testing.T) {
	agent := &entity.Agent{
		Paw:       "offline-agent",
		Platform:  "windows",
		Executors: []string{"psh"},
		Status:    entity.AgentOffline, // Not online
	}

	agentRepo := &mockAgentRepo{
		agents: map[string]*entity.Agent{"offline-agent": agent},
	}

	techRepo := &mockTechniqueRepo{}
	validator := NewTechniqueValidator()
	orchestrator := NewAttackOrchestrator(agentRepo, techRepo, validator, nil)

	plan := &ExecutionPlan{
		ID: "test-plan",
		Tasks: []PlannedTask{
			{TechniqueID: "T1059", AgentPaw: "offline-agent"},
		},
	}

	err := orchestrator.ValidatePlan(context.Background(), plan)

	if err == nil {
		t.Error("Expected error for offline agent")
	}
}

func TestAttackOrchestrator_ValidatePlan_TechniqueNotFound(t *testing.T) {
	agent := &entity.Agent{
		Paw:       "test-agent",
		Platform:  "windows",
		Executors: []string{"psh"},
		Status:    entity.AgentOnline,
	}

	agentRepo := &mockAgentRepo{
		agents: map[string]*entity.Agent{"test-agent": agent},
	}

	techRepo := &mockTechniqueRepo{
		techniques: map[string]*entity.Technique{}, // empty
	}

	validator := NewTechniqueValidator()
	orchestrator := NewAttackOrchestrator(agentRepo, techRepo, validator, nil)

	plan := &ExecutionPlan{
		ID: "test-plan",
		Tasks: []PlannedTask{
			{TechniqueID: "T9999", AgentPaw: "test-agent"}, // doesn't exist
		},
	}

	err := orchestrator.ValidatePlan(context.Background(), plan)

	if err == nil {
		t.Error("Expected error for missing technique")
	}
}

func TestAttackOrchestrator_ValidatePlan_IncompatibleAgent(t *testing.T) {
	agent := &entity.Agent{
		Paw:       "linux-agent",
		Platform:  "linux",
		Executors: []string{"bash"},
		Status:    entity.AgentOnline,
	}

	technique := &entity.Technique{
		ID:        "T1059",
		Platforms: []string{"windows"}, // Windows only
		Executors: []entity.Executor{
			{Type: "psh", Command: "whoami"},
		},
	}

	agentRepo := &mockAgentRepo{
		agents: map[string]*entity.Agent{"linux-agent": agent},
	}

	techRepo := &mockTechniqueRepo{
		techniques: map[string]*entity.Technique{"T1059": technique},
	}

	validator := NewTechniqueValidator()
	orchestrator := NewAttackOrchestrator(agentRepo, techRepo, validator, nil)

	plan := &ExecutionPlan{
		ID: "test-plan",
		Tasks: []PlannedTask{
			{TechniqueID: "T1059", AgentPaw: "linux-agent"},
		},
	}

	err := orchestrator.ValidatePlan(context.Background(), plan)

	if err == nil {
		t.Error("Expected error for incompatible agent")
	}
}

func TestAttackOrchestrator_PlanExecution_NoCompatibleExecutor(t *testing.T) {
	// Technique is compatible with platform but agent has incompatible executor
	technique := &entity.Technique{
		ID:        "T1059",
		Platforms: []string{"windows"},
		Executors: []entity.Executor{
			{Type: "psh", Command: "whoami", Timeout: 30}, // Requires psh
		},
		IsSafe: true,
	}

	techRepo := &mockTechniqueRepo{
		techniques: map[string]*entity.Technique{"T1059": technique},
	}

	agentRepo := &mockAgentRepo{}
	validator := NewTechniqueValidator()
	orchestrator := NewAttackOrchestrator(agentRepo, techRepo, validator, nil)

	// Agent is Windows but only has cmd, not psh
	agent := &entity.Agent{
		Paw:       "cmd-only-agent",
		Platform:  "windows",
		Executors: []string{"cmd"}, // Only cmd, not psh
		Status:    entity.AgentOnline,
	}

	scenario := &entity.Scenario{
		ID:   "test-scenario",
		Name: "Test Scenario",
		Phases: []entity.Phase{
			{Name: "Phase1", Techniques: []entity.TechniqueSelection{{TechniqueID: "T1059"}}},
		},
	}

	_, err := orchestrator.PlanExecution(context.Background(), scenario, []*entity.Agent{agent}, false)

	if err == nil {
		t.Error("Expected error when no compatible executor found")
	}
}

func TestAttackOrchestrator_PlanExecution_WithExecutorName(t *testing.T) {
	technique := &entity.Technique{
		ID:        "T1059",
		Name:      "Command Execution",
		Platforms: []string{"windows"},
		Executors: []entity.Executor{
			{Name: "whoami-exec", Type: "psh", Platform: "windows", Command: "whoami", Timeout: 30},
			{Name: "hostname-exec", Type: "cmd", Platform: "windows", Command: "hostname", Timeout: 30},
		},
		IsSafe: true,
	}

	techRepo := &mockTechniqueRepo{
		techniques: map[string]*entity.Technique{"T1059": technique},
	}
	agentRepo := &mockAgentRepo{}
	validator := NewTechniqueValidator()
	orchestrator := NewAttackOrchestrator(agentRepo, techRepo, validator, nil)

	agent := &entity.Agent{
		Paw:       "test-agent",
		Platform:  "windows",
		Executors: []string{"psh", "cmd"},
		Status:    entity.AgentOnline,
	}

	scenario := &entity.Scenario{
		ID:   "test-scenario",
		Name: "Test Scenario",
		Phases: []entity.Phase{
			{Name: "Phase1", Techniques: []entity.TechniqueSelection{
				{TechniqueID: "T1059", ExecutorName: "hostname-exec"},
			}},
		},
	}

	plan, err := orchestrator.PlanExecution(context.Background(), scenario, []*entity.Agent{agent}, false)
	if err != nil {
		t.Fatalf("PlanExecution returned error: %v", err)
	}
	if len(plan.Tasks) != 1 {
		t.Fatalf("Expected 1 task, got %d", len(plan.Tasks))
	}
	if plan.Tasks[0].Command != "hostname" {
		t.Errorf("Expected command 'hostname' from named executor, got '%s'", plan.Tasks[0].Command)
	}
}

func TestAttackOrchestrator_PlanExecution_WithInvalidExecutorName(t *testing.T) {
	technique := &entity.Technique{
		ID:        "T1059",
		Name:      "Command Execution",
		Platforms: []string{"windows"},
		Executors: []entity.Executor{
			{Name: "whoami-exec", Type: "psh", Platform: "windows", Command: "whoami", Timeout: 30},
		},
		IsSafe: true,
	}

	techRepo := &mockTechniqueRepo{
		techniques: map[string]*entity.Technique{"T1059": technique},
	}
	agentRepo := &mockAgentRepo{}
	validator := NewTechniqueValidator()
	orchestrator := NewAttackOrchestrator(agentRepo, techRepo, validator, nil)

	agent := &entity.Agent{
		Paw:       "test-agent",
		Platform:  "windows",
		Executors: []string{"psh"},
		Status:    entity.AgentOnline,
	}

	scenario := &entity.Scenario{
		ID:   "test-scenario",
		Name: "Test Scenario",
		Phases: []entity.Phase{
			{Name: "Phase1", Techniques: []entity.TechniqueSelection{
				{TechniqueID: "T1059", ExecutorName: "nonexistent-exec"},
			}},
		},
	}

	// Should fallback to auto-select since named executor doesn't exist
	plan, err := orchestrator.PlanExecution(context.Background(), scenario, []*entity.Agent{agent}, false)
	if err != nil {
		t.Fatalf("PlanExecution returned error: %v", err)
	}
	if len(plan.Tasks) != 1 {
		t.Fatalf("Expected 1 task (fallback), got %d", len(plan.Tasks))
	}
	if plan.Tasks[0].Command != "whoami" {
		t.Errorf("Expected fallback command 'whoami', got '%s'", plan.Tasks[0].Command)
	}
}

func TestExecutionPlan_Struct(t *testing.T) {
	plan := &ExecutionPlan{
		ID: "plan-123",
		Tasks: []PlannedTask{
			{
				TechniqueID: "T1059",
				AgentPaw:    "agent-1",
				Phase:       "Execution",
				Order:       0,
				Command:     "whoami",
				Cleanup:     "rm log.txt",
				Timeout:     30,
			},
		},
	}

	if plan.ID != "plan-123" {
		t.Errorf("ID = %s, want plan-123", plan.ID)
	}
	if len(plan.Tasks) != 1 {
		t.Errorf("Tasks length = %d, want 1", len(plan.Tasks))
	}
	if plan.Tasks[0].Timeout != 30 {
		t.Errorf("Timeout = %d, want 30", plan.Tasks[0].Timeout)
	}
}

func TestAttackOrchestrator_PlanExecution_AllExecutors(t *testing.T) {
	// Technique with 3 linux executors — all should produce tasks when no executor name is specified
	technique := &entity.Technique{
		ID:        "T1552.001",
		Name:      "Credentials In Files",
		Platforms: []string{"linux"},
		Executors: []entity.Executor{
			{Name: "Find AWS credentials", Type: "sh", Platform: "linux", Command: "find //.aws -name credentials 2>/dev/null", Timeout: 60, IsSafe: true},
			{Name: "Find Github Credentials", Type: "bash", Platform: "linux", Command: "find /home -name .netrc 2>/dev/null", Timeout: 60, IsSafe: true},
			{Name: "Find Azure credentials", Type: "sh", Platform: "linux", Command: "find //.azure -name msal_token_cache.json 2>/dev/null", Timeout: 60, IsSafe: true},
		},
		IsSafe: true,
	}

	techRepo := &mockTechniqueRepo{
		techniques: map[string]*entity.Technique{"T1552.001": technique},
	}
	agentRepo := &mockAgentRepo{}
	validator := NewTechniqueValidator()
	orchestrator := NewAttackOrchestrator(agentRepo, techRepo, validator, nil)

	agent := &entity.Agent{
		Paw:       "linux-agent",
		Platform:  "linux",
		Executors: []string{"sh", "bash"},
		Status:    entity.AgentOnline,
	}

	scenario := &entity.Scenario{
		ID:   "test-scenario",
		Name: "Credential Test",
		Phases: []entity.Phase{
			{Name: "Credentials", Techniques: []entity.TechniqueSelection{{TechniqueID: "T1552.001"}}},
		},
	}

	plan, err := orchestrator.PlanExecution(context.Background(), scenario, []*entity.Agent{agent}, false)
	if err != nil {
		t.Fatalf("PlanExecution returned error: %v", err)
	}

	// All 3 executors should produce tasks
	if len(plan.Tasks) != 3 {
		t.Fatalf("Expected 3 tasks (one per executor), got %d", len(plan.Tasks))
	}

	// Verify each task has different executor names
	names := map[string]bool{}
	for _, task := range plan.Tasks {
		names[task.ExecutorName] = true
		if task.TechniqueID != "T1552.001" {
			t.Errorf("Expected technique T1552.001, got %s", task.TechniqueID)
		}
		if task.ExecutorType == "" {
			t.Error("ExecutorType should not be empty")
		}
	}
	if len(names) != 3 {
		t.Errorf("Expected 3 unique executor names, got %d", len(names))
	}
}

func TestAttackOrchestrator_PlanExecution_NamedExecutorStillSingle(t *testing.T) {
	// When executor_name is specified, only that executor should run (even with multiple available)
	technique := &entity.Technique{
		ID:        "T1552.001",
		Name:      "Credentials In Files",
		Platforms: []string{"linux"},
		Executors: []entity.Executor{
			{Name: "Find AWS credentials", Type: "sh", Platform: "linux", Command: "find //.aws", Timeout: 60, IsSafe: true},
			{Name: "Find Github Credentials", Type: "bash", Platform: "linux", Command: "find /home -name .netrc", Timeout: 60, IsSafe: true},
			{Name: "Find Azure credentials", Type: "sh", Platform: "linux", Command: "find //.azure", Timeout: 60, IsSafe: true},
		},
		IsSafe: true,
	}

	techRepo := &mockTechniqueRepo{
		techniques: map[string]*entity.Technique{"T1552.001": technique},
	}
	agentRepo := &mockAgentRepo{}
	validator := NewTechniqueValidator()
	orchestrator := NewAttackOrchestrator(agentRepo, techRepo, validator, nil)

	agent := &entity.Agent{
		Paw:       "linux-agent",
		Platform:  "linux",
		Executors: []string{"sh", "bash"},
		Status:    entity.AgentOnline,
	}

	scenario := &entity.Scenario{
		ID:   "test-scenario",
		Name: "Credential Test",
		Phases: []entity.Phase{
			{Name: "Credentials", Techniques: []entity.TechniqueSelection{
				{TechniqueID: "T1552.001", ExecutorName: "Find Github Credentials"},
			}},
		},
	}

	plan, err := orchestrator.PlanExecution(context.Background(), scenario, []*entity.Agent{agent}, false)
	if err != nil {
		t.Fatalf("PlanExecution returned error: %v", err)
	}

	// Only 1 task — the named executor
	if len(plan.Tasks) != 1 {
		t.Fatalf("Expected 1 task (named executor only), got %d", len(plan.Tasks))
	}
	if plan.Tasks[0].Command != "find /home -name .netrc" {
		t.Errorf("Expected Github credentials command, got '%s'", plan.Tasks[0].Command)
	}
	if plan.Tasks[0].ExecutorName != "Find Github Credentials" {
		t.Errorf("Expected executor name 'Find Github Credentials', got '%s'", plan.Tasks[0].ExecutorName)
	}
}

func TestAttackOrchestrator_PlanExecution_SafeMode_BackwardCompat(t *testing.T) {
	// Technique is marked safe but ALL executors have IsSafe: false (legacy format).
	// This triggers the backward compatibility path at lines 133-137 where
	// len(safeExecutors) == 0 but technique.IsSafe is true, so it returns
	// the technique with ALL executors.
	technique := &entity.Technique{
		ID:        "T1082",
		Name:      "System Info",
		Platforms: []string{"windows"},
		Executors: []entity.Executor{
			{Type: "psh", Command: "systeminfo", Timeout: 30, IsSafe: false},
		},
		IsSafe: true, // Technique is safe, but no executor has is_safe=true
	}

	techRepo := &mockTechniqueRepo{
		techniques: map[string]*entity.Technique{"T1082": technique},
	}

	agentRepo := &mockAgentRepo{}
	validator := NewTechniqueValidator()
	orchestrator := NewAttackOrchestrator(agentRepo, techRepo, validator, nil)

	agent := &entity.Agent{
		Paw:       "test-agent",
		Platform:  "windows",
		Executors: []string{"psh"},
		Status:    entity.AgentOnline,
	}

	scenario := &entity.Scenario{
		ID:   "test-scenario",
		Name: "Test Scenario",
		Phases: []entity.Phase{
			{Name: "Phase1", Techniques: []entity.TechniqueSelection{{TechniqueID: "T1082"}}},
		},
	}

	// Safe mode with backward compat: technique is safe but no executor has is_safe=true
	plan, err := orchestrator.PlanExecution(context.Background(), scenario, []*entity.Agent{agent}, true)

	if err != nil {
		t.Fatalf("PlanExecution returned error: %v", err)
	}
	if plan == nil {
		t.Fatal("PlanExecution returned nil plan")
	}
	if len(plan.Tasks) != 1 {
		t.Errorf("Expected 1 task (backward compat returns all executors), got %d", len(plan.Tasks))
	}
	if plan.Tasks[0].Command != "systeminfo" {
		t.Errorf("Expected command 'systeminfo', got '%s'", plan.Tasks[0].Command)
	}
}
