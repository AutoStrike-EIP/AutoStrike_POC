package sqlite

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"autostrike/internal/domain/entity"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	if err := InitSchema(db); err != nil {
		t.Fatalf("Failed to init schema: %v", err)
	}
	return db
}

// Schema tests
func TestInitSchema(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	err = InitSchema(db)
	if err != nil {
		t.Fatalf("InitSchema failed: %v", err)
	}

	// Verify tables exist
	tables := []string{"agents", "techniques", "scenarios", "executions", "execution_results"}
	for _, table := range tables {
		var name string
		err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&name)
		if err != nil {
			t.Errorf("Table %s not created: %v", table, err)
		}
	}
}

// Agent Repository tests
func TestNewAgentRepository(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewAgentRepository(db)
	if repo == nil {
		t.Error("Expected non-nil repository")
	}
}

func TestAgentRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewAgentRepository(db)
	ctx := context.Background()

	agent := &entity.Agent{
		Paw:       "test-paw",
		Hostname:  "test-host",
		Username:  "test-user",
		Platform:  "linux",
		Executors: []string{"sh", "bash"},
		Status:    entity.AgentOnline,
		LastSeen:  time.Now(),
		}

	err := repo.Create(ctx, agent)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
}

func TestAgentRepository_FindByPaw(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewAgentRepository(db)
	ctx := context.Background()

	agent := &entity.Agent{
		Paw:       "test-paw",
		Hostname:  "test-host",
		Username:  "test-user",
		Platform:  "linux",
		Executors: []string{"sh"},
		Status:    entity.AgentOnline,
		LastSeen:  time.Now(),
		}
	_ = repo.Create(ctx, agent)

	found, err := repo.FindByPaw(ctx, "test-paw")
	if err != nil {
		t.Fatalf("FindByPaw failed: %v", err)
	}
	if found.Hostname != "test-host" {
		t.Errorf("Expected hostname test-host, got %s", found.Hostname)
	}
}

func TestAgentRepository_FindByPaw_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewAgentRepository(db)
	ctx := context.Background()

	_, err := repo.FindByPaw(ctx, "nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent agent")
	}
}

func TestAgentRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewAgentRepository(db)
	ctx := context.Background()

	agent := &entity.Agent{
		Paw:       "test-paw",
		Hostname:  "old-host",
		Username:  "test-user",
		Platform:  "linux",
		Executors: []string{"sh"},
		Status:    entity.AgentOnline,
		LastSeen:  time.Now(),
		}
	_ = repo.Create(ctx, agent)

	agent.Hostname = "new-host"
	err := repo.Update(ctx, agent)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	found, _ := repo.FindByPaw(ctx, "test-paw")
	if found.Hostname != "new-host" {
		t.Errorf("Expected hostname new-host, got %s", found.Hostname)
	}
}

func TestAgentRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewAgentRepository(db)
	ctx := context.Background()

	agent := &entity.Agent{
		Paw:       "test-paw",
		Hostname:  "test-host",
		Username:  "test-user",
		Platform:  "linux",
		Executors: []string{"sh"},
		Status:    entity.AgentOnline,
		LastSeen:  time.Now(),
		}
	_ = repo.Create(ctx, agent)

	err := repo.Delete(ctx, "test-paw")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = repo.FindByPaw(ctx, "test-paw")
	if err == nil {
		t.Error("Expected error after delete")
	}
}

func TestAgentRepository_FindAll(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewAgentRepository(db)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		agent := &entity.Agent{
			Paw:       "paw-" + string(rune('a'+i)),
			Hostname:  "host",
			Username:  "user",
			Platform:  "linux",
			Executors: []string{"sh"},
			Status:    entity.AgentOnline,
			LastSeen:  time.Now(),
				}
	_ = repo.Create(ctx, agent)
	}

	agents, err := repo.FindAll(ctx)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	if len(agents) != 3 {
		t.Errorf("Expected 3 agents, got %d", len(agents))
	}
}

func TestAgentRepository_FindByStatus(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewAgentRepository(db)
	ctx := context.Background()

	online := &entity.Agent{
		Paw: "online", Hostname: "h", Username: "u", Platform: "linux",
		Executors: []string{"sh"}, Status: entity.AgentOnline,
		LastSeen: time.Now(), CreatedAt: time.Now(),
	}
	offline := &entity.Agent{
		Paw: "offline", Hostname: "h", Username: "u", Platform: "linux",
		Executors: []string{"sh"}, Status: entity.AgentOffline,
		LastSeen: time.Now(), CreatedAt: time.Now(),
	}
	_ = repo.Create(ctx, online)
	_ = repo.Create(ctx, offline)

	agents, err := repo.FindByStatus(ctx, entity.AgentOnline)
	if err != nil {
		t.Fatalf("FindByStatus failed: %v", err)
	}
	if len(agents) != 1 {
		t.Errorf("Expected 1 online agent, got %d", len(agents))
	}
}

func TestAgentRepository_FindByPlatform(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewAgentRepository(db)
	ctx := context.Background()

	linux := &entity.Agent{
		Paw: "linux", Hostname: "h", Username: "u", Platform: "linux",
		Executors: []string{"sh"}, Status: entity.AgentOnline,
		LastSeen: time.Now(), CreatedAt: time.Now(),
	}
	windows := &entity.Agent{
		Paw: "windows", Hostname: "h", Username: "u", Platform: "windows",
		Executors: []string{"cmd"}, Status: entity.AgentOnline,
		LastSeen: time.Now(), CreatedAt: time.Now(),
	}
	_ = repo.Create(ctx, linux)
	_ = repo.Create(ctx, windows)

	agents, err := repo.FindByPlatform(ctx, "linux")
	if err != nil {
		t.Fatalf("FindByPlatform failed: %v", err)
	}
	if len(agents) != 1 {
		t.Errorf("Expected 1 linux agent, got %d", len(agents))
	}
}

func TestAgentRepository_UpdateLastSeen(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewAgentRepository(db)
	ctx := context.Background()

	agent := &entity.Agent{
		Paw: "test", Hostname: "h", Username: "u", Platform: "linux",
		Executors: []string{"sh"}, Status: entity.AgentOffline,
		LastSeen: time.Now().Add(-time.Hour), CreatedAt: time.Now(),
	}
	_ = repo.Create(ctx, agent)

	err := repo.UpdateLastSeen(ctx, "test")
	if err != nil {
		t.Fatalf("UpdateLastSeen failed: %v", err)
	}

	found, _ := repo.FindByPaw(ctx, "test")
	if found.Status != entity.AgentOnline {
		t.Error("Expected agent to be online after UpdateLastSeen")
	}
}

// Technique Repository tests
func TestNewTechniqueRepository(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewTechniqueRepository(db)
	if repo == nil {
		t.Error("Expected non-nil repository")
	}
}

func TestTechniqueRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewTechniqueRepository(db)
	ctx := context.Background()

	tech := &entity.Technique{
		ID:          "T1059",
		Name:        "Command Execution",
		Description: "Execute commands",
		Tactic:      entity.TacticExecution,
		Platforms:   []string{"windows", "linux"},
		Executors: []entity.Executor{
			{Type: "sh", Command: "whoami"},
		},
		IsSafe:    true,
		}

	err := repo.Create(ctx, tech)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
}

func TestTechniqueRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewTechniqueRepository(db)
	ctx := context.Background()

	tech := &entity.Technique{
		ID:        "T1059",
		Name:      "Command Execution",
		Tactic:    entity.TacticExecution,
		Platforms: []string{"linux"},
		Executors: []entity.Executor{{Type: "sh", Command: "whoami"}},
		IsSafe:    true,
		}
	_ = repo.Create(ctx, tech)

	found, err := repo.FindByID(ctx, "T1059")
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if found.Name != "Command Execution" {
		t.Errorf("Expected name Command Execution, got %s", found.Name)
	}
}

func TestTechniqueRepository_FindByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewTechniqueRepository(db)
	ctx := context.Background()

	_, err := repo.FindByID(ctx, "nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent technique")
	}
}

func TestTechniqueRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewTechniqueRepository(db)
	ctx := context.Background()

	tech := &entity.Technique{
		ID:        "T1059",
		Name:      "Old Name",
		Tactic:    entity.TacticExecution,
		Platforms: []string{"linux"},
		Executors: []entity.Executor{{Type: "sh", Command: "whoami"}},
		IsSafe:    true,
		}
	_ = repo.Create(ctx, tech)

	tech.Name = "New Name"
	err := repo.Update(ctx, tech)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	found, _ := repo.FindByID(ctx, "T1059")
	if found.Name != "New Name" {
		t.Errorf("Expected name New Name, got %s", found.Name)
	}
}

func TestTechniqueRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewTechniqueRepository(db)
	ctx := context.Background()

	tech := &entity.Technique{
		ID:        "T1059",
		Name:      "Test",
		Tactic:    entity.TacticExecution,
		Platforms: []string{"linux"},
		Executors: []entity.Executor{{Type: "sh", Command: "whoami"}},
		}
	_ = repo.Create(ctx, tech)

	err := repo.Delete(ctx, "T1059")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = repo.FindByID(ctx, "T1059")
	if err == nil {
		t.Error("Expected error after delete")
	}
}

func TestTechniqueRepository_FindAll(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewTechniqueRepository(db)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		tech := &entity.Technique{
			ID:        "T" + string(rune('0'+i)),
			Name:      "Tech",
			Tactic:    entity.TacticExecution,
			Platforms: []string{"linux"},
			Executors: []entity.Executor{{Type: "sh", Command: "cmd"}},
				}
	_ = repo.Create(ctx, tech)
	}

	techniques, err := repo.FindAll(ctx)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	if len(techniques) != 3 {
		t.Errorf("Expected 3 techniques, got %d", len(techniques))
	}
}

func TestTechniqueRepository_FindByTactic(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewTechniqueRepository(db)
	ctx := context.Background()

	exec := &entity.Technique{
		ID: "T1", Name: "Exec", Tactic: entity.TacticExecution,
		Platforms: []string{"linux"}, Executors: []entity.Executor{{Type: "sh", Command: "c"}},
		}
	persist := &entity.Technique{
		ID: "T2", Name: "Persist", Tactic: entity.TacticPersistence,
		Platforms: []string{"linux"}, Executors: []entity.Executor{{Type: "sh", Command: "c"}},
		}
	_ = repo.Create(ctx, exec)
	_ = repo.Create(ctx, persist)

	techniques, err := repo.FindByTactic(ctx, entity.TacticExecution)
	if err != nil {
		t.Fatalf("FindByTactic failed: %v", err)
	}
	if len(techniques) != 1 {
		t.Errorf("Expected 1 execution technique, got %d", len(techniques))
	}
}

func TestTechniqueRepository_FindByPlatform(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewTechniqueRepository(db)
	ctx := context.Background()

	linux := &entity.Technique{
		ID: "T1", Name: "Linux", Tactic: entity.TacticExecution,
		Platforms: []string{"linux"}, Executors: []entity.Executor{{Type: "sh", Command: "c"}},
		}
	windows := &entity.Technique{
		ID: "T2", Name: "Windows", Tactic: entity.TacticExecution,
		Platforms: []string{"windows"}, Executors: []entity.Executor{{Type: "cmd", Command: "c"}},
		}
	_ = repo.Create(ctx, linux)
	_ = repo.Create(ctx, windows)

	techniques, err := repo.FindByPlatform(ctx, "linux")
	if err != nil {
		t.Fatalf("FindByPlatform failed: %v", err)
	}
	if len(techniques) != 1 {
		t.Errorf("Expected 1 linux technique, got %d", len(techniques))
	}
}

// Scenario Repository tests
func TestNewScenarioRepository(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewScenarioRepository(db)
	if repo == nil {
		t.Error("Expected non-nil repository")
	}
}

func TestScenarioRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewScenarioRepository(db)
	ctx := context.Background()

	scenario := &entity.Scenario{
		ID:          "s1",
		Name:        "Test Scenario",
		Description: "A test",
		Phases: []entity.Phase{
			{Name: "Phase1", Techniques: []string{"T1059"}},
		},
		Tags:      []string{"test"},
			UpdatedAt: time.Now(),
	}

	err := repo.Create(ctx, scenario)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
}

func TestScenarioRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewScenarioRepository(db)
	ctx := context.Background()

	scenario := &entity.Scenario{
		ID:   "s1",
		Name: "Test",
		Phases: []entity.Phase{
			{Name: "Phase1", Techniques: []string{"T1059"}},
		},
		Tags:      []string{"test"},
			UpdatedAt: time.Now(),
	}
	_ = repo.Create(ctx, scenario)

	found, err := repo.FindByID(ctx, "s1")
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if found.Name != "Test" {
		t.Errorf("Expected name Test, got %s", found.Name)
	}
}

func TestScenarioRepository_FindByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewScenarioRepository(db)
	ctx := context.Background()

	_, err := repo.FindByID(ctx, "nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent scenario")
	}
}

func TestScenarioRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewScenarioRepository(db)
	ctx := context.Background()

	scenario := &entity.Scenario{
		ID:        "s1",
		Name:      "Old Name",
		Phases:    []entity.Phase{{Name: "P1", Techniques: []string{"T1"}}},
		Tags:      []string{},
			UpdatedAt: time.Now(),
	}
	_ = repo.Create(ctx, scenario)

	scenario.Name = "New Name"
	err := repo.Update(ctx, scenario)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	found, _ := repo.FindByID(ctx, "s1")
	if found.Name != "New Name" {
		t.Errorf("Expected name New Name, got %s", found.Name)
	}
}

func TestScenarioRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewScenarioRepository(db)
	ctx := context.Background()

	scenario := &entity.Scenario{
		ID:        "s1",
		Name:      "Test",
		Phases:    []entity.Phase{{Name: "P1", Techniques: []string{"T1"}}},
		Tags:      []string{},
			UpdatedAt: time.Now(),
	}
	_ = repo.Create(ctx, scenario)

	err := repo.Delete(ctx, "s1")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = repo.FindByID(ctx, "s1")
	if err == nil {
		t.Error("Expected error after delete")
	}
}

func TestScenarioRepository_FindAll(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewScenarioRepository(db)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		scenario := &entity.Scenario{
			ID:        "s" + string(rune('0'+i)),
			Name:      "Scenario",
			Phases:    []entity.Phase{{Name: "P", Techniques: []string{"T1"}}},
			Tags:      []string{},
					UpdatedAt: time.Now(),
		}
	_ = repo.Create(ctx, scenario)
	}

	scenarios, err := repo.FindAll(ctx)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	if len(scenarios) != 3 {
		t.Errorf("Expected 3 scenarios, got %d", len(scenarios))
	}
}

func TestScenarioRepository_FindByTag(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewScenarioRepository(db)
	ctx := context.Background()

	tagged := &entity.Scenario{
		ID:        "s1",
		Name:      "Tagged",
		Phases:    []entity.Phase{{Name: "P", Techniques: []string{"T1"}}},
		Tags:      []string{"important"},
			UpdatedAt: time.Now(),
	}
	untagged := &entity.Scenario{
		ID:        "s2",
		Name:      "Untagged",
		Phases:    []entity.Phase{{Name: "P", Techniques: []string{"T1"}}},
		Tags:      []string{"other"},
			UpdatedAt: time.Now(),
	}
	_ = repo.Create(ctx, tagged)
	_ = repo.Create(ctx, untagged)

	scenarios, err := repo.FindByTag(ctx, "important")
	if err != nil {
		t.Fatalf("FindByTag failed: %v", err)
	}
	if len(scenarios) != 1 {
		t.Errorf("Expected 1 tagged scenario, got %d", len(scenarios))
	}
}

// Result Repository tests
func TestNewResultRepository(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewResultRepository(db)
	if repo == nil {
		t.Error("Expected non-nil repository")
	}
}

func TestResultRepository_CreateExecution(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewResultRepository(db)
	ctx := context.Background()

	exec := &entity.Execution{
		ID:         "e1",
		ScenarioID: "s1",
		Status:     entity.ExecutionRunning,
		StartedAt:  time.Now(),
		SafeMode:   true,
	}

	err := repo.CreateExecution(ctx, exec)
	if err != nil {
		t.Fatalf("CreateExecution failed: %v", err)
	}
}

func TestResultRepository_FindExecutionByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewResultRepository(db)
	ctx := context.Background()

	exec := &entity.Execution{
		ID:         "e1",
		ScenarioID: "s1",
		Status:     entity.ExecutionRunning,
		StartedAt:  time.Now(),
	}
	_ = repo.CreateExecution(ctx, exec)

	found, err := repo.FindExecutionByID(ctx, "e1")
	if err != nil {
		t.Fatalf("FindExecutionByID failed: %v", err)
	}
	if found.ScenarioID != "s1" {
		t.Errorf("Expected scenario s1, got %s", found.ScenarioID)
	}
}

func TestResultRepository_FindExecutionByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewResultRepository(db)
	ctx := context.Background()

	_, err := repo.FindExecutionByID(ctx, "nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent execution")
	}
}

func TestResultRepository_UpdateExecution(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewResultRepository(db)
	ctx := context.Background()

	exec := &entity.Execution{
		ID:         "e1",
		ScenarioID: "s1",
		Status:     entity.ExecutionRunning,
		StartedAt:  time.Now(),
	}
	_ = repo.CreateExecution(ctx, exec)

	now := time.Now()
	exec.Status = entity.ExecutionCompleted
	exec.CompletedAt = &now
	exec.Score = &entity.SecurityScore{Overall: 0.8, Blocked: 1, Detected: 1, Successful: 8, Total: 10}

	err := repo.UpdateExecution(ctx, exec)
	if err != nil {
		t.Fatalf("UpdateExecution failed: %v", err)
	}

	found, _ := repo.FindExecutionByID(ctx, "e1")
	if found.Status != entity.ExecutionCompleted {
		t.Error("Expected status Completed")
	}
}

func TestResultRepository_FindExecutionsByScenario(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewResultRepository(db)
	ctx := context.Background()

	for i := 0; i < 2; i++ {
		exec := &entity.Execution{
			ID:         "e" + string(rune('0'+i)),
			ScenarioID: "s1",
			Status:     entity.ExecutionCompleted,
			StartedAt:  time.Now(),
		}
	_ = repo.CreateExecution(ctx, exec)
	}
	exec3 := &entity.Execution{
		ID:         "e3",
		ScenarioID: "s2",
		Status:     entity.ExecutionCompleted,
		StartedAt:  time.Now(),
	}
	_ = repo.CreateExecution(ctx, exec3)

	executions, err := repo.FindExecutionsByScenario(ctx, "s1")
	if err != nil {
		t.Fatalf("FindExecutionsByScenario failed: %v", err)
	}
	if len(executions) != 2 {
		t.Errorf("Expected 2 executions for s1, got %d", len(executions))
	}
}

func TestResultRepository_FindRecentExecutions(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewResultRepository(db)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		exec := &entity.Execution{
			ID:         "e" + string(rune('0'+i)),
			ScenarioID: "s1",
			Status:     entity.ExecutionCompleted,
			StartedAt:  time.Now(),
		}
	_ = repo.CreateExecution(ctx, exec)
	}

	executions, err := repo.FindRecentExecutions(ctx, 3)
	if err != nil {
		t.Fatalf("FindRecentExecutions failed: %v", err)
	}
	if len(executions) != 3 {
		t.Errorf("Expected 3 recent executions, got %d", len(executions))
	}
}

func TestResultRepository_CreateResult(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewResultRepository(db)
	ctx := context.Background()

	result := &entity.ExecutionResult{
		ID:          "r1",
		ExecutionID: "e1",
		TechniqueID: "T1059",
		AgentPaw:    "paw1",
		Status:      entity.StatusPending,
		StartedAt:   time.Now(),
	}

	err := repo.CreateResult(ctx, result)
	if err != nil {
		t.Fatalf("CreateResult failed: %v", err)
	}
}

func TestResultRepository_UpdateResult(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewResultRepository(db)
	ctx := context.Background()

	result := &entity.ExecutionResult{
		ID:          "r1",
		ExecutionID: "e1",
		TechniqueID: "T1059",
		AgentPaw:    "paw1",
		Status:      entity.StatusPending,
		StartedAt:   time.Now(),
	}
	_ = repo.CreateResult(ctx, result)

	now := time.Now()
	result.Status = entity.StatusSuccess
	result.Output = "test output"
	result.CompletedAt = &now

	err := repo.UpdateResult(ctx, result)
	if err != nil {
		t.Fatalf("UpdateResult failed: %v", err)
	}
}

func TestResultRepository_FindResultsByExecution(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewResultRepository(db)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		result := &entity.ExecutionResult{
			ID:          "r" + string(rune('0'+i)),
			ExecutionID: "e1",
			TechniqueID: "T1059",
			AgentPaw:    "paw1",
			Status:      entity.StatusSuccess,
			StartedAt:   time.Now(),
		}
	_ = repo.CreateResult(ctx, result)
	}

	results, err := repo.FindResultsByExecution(ctx, "e1")
	if err != nil {
		t.Fatalf("FindResultsByExecution failed: %v", err)
	}
	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}
}

func TestResultRepository_FindResultsByTechnique(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewResultRepository(db)
	ctx := context.Background()

	r1 := &entity.ExecutionResult{
		ID: "r1", ExecutionID: "e1", TechniqueID: "T1059",
		AgentPaw: "paw1", Status: entity.StatusSuccess, StartedAt: time.Now(),
	}
	r2 := &entity.ExecutionResult{
		ID: "r2", ExecutionID: "e1", TechniqueID: "T1055",
		AgentPaw: "paw1", Status: entity.StatusSuccess, StartedAt: time.Now(),
	}
	_ = repo.CreateResult(ctx, r1)
	_ = repo.CreateResult(ctx, r2)

	results, err := repo.FindResultsByTechnique(ctx, "T1059")
	if err != nil {
		t.Fatalf("FindResultsByTechnique failed: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("Expected 1 result for T1059, got %d", len(results))
	}
}
