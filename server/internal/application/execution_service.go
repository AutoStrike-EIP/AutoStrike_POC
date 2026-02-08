package application

import (
	"context"
	"fmt"
	"time"

	"autostrike/internal/domain/entity"
	"autostrike/internal/domain/repository"
	"autostrike/internal/domain/service"

	"github.com/google/uuid"
)

// ExecutionService handles execution-related business logic
type ExecutionService struct {
	resultRepo    repository.ResultRepository
	scenarioRepo  repository.ScenarioRepository
	techniqueRepo repository.TechniqueRepository
	agentRepo     repository.AgentRepository
	orchestrator  *service.AttackOrchestrator
	calculator    *service.ScoreCalculator
}

// NewExecutionService creates a new execution service
func NewExecutionService(
	resultRepo repository.ResultRepository,
	scenarioRepo repository.ScenarioRepository,
	techniqueRepo repository.TechniqueRepository,
	agentRepo repository.AgentRepository,
	orchestrator *service.AttackOrchestrator,
	calculator *service.ScoreCalculator,
) *ExecutionService {
	return &ExecutionService{
		resultRepo:    resultRepo,
		scenarioRepo:  scenarioRepo,
		techniqueRepo: techniqueRepo,
		agentRepo:     agentRepo,
		orchestrator:  orchestrator,
		calculator:    calculator,
	}
}

// TaskDispatchInfo contains information needed to dispatch a task to an agent
type TaskDispatchInfo struct {
	ResultID     string
	AgentPaw     string
	TechniqueID  string
	Command      string
	Executor     string
	ExecutorName string
	Timeout      int
	Cleanup      string
}

// ExecutionWithTasks contains the execution and tasks to dispatch
type ExecutionWithTasks struct {
	Execution *entity.Execution
	Tasks     []TaskDispatchInfo
}

// StartExecution starts a new scenario execution
func (s *ExecutionService) StartExecution(
	ctx context.Context,
	scenarioID string,
	agentPaws []string,
	safeMode bool,
) (*ExecutionWithTasks, error) {
	scenario, err := s.scenarioRepo.FindByID(ctx, scenarioID)
	if err != nil {
		return nil, fmt.Errorf("scenario not found: %w", err)
	}

	agentMap, agents, err := s.loadAndValidateAgents(ctx, agentPaws)
	if err != nil {
		return nil, err
	}

	plan, err := s.orchestrator.PlanExecution(ctx, scenario, agents, safeMode)
	if err != nil {
		return nil, fmt.Errorf("failed to plan execution: %w", err)
	}

	execution := &entity.Execution{
		ID:         uuid.New().String(),
		ScenarioID: scenarioID,
		AgentPaws:  agentPaws,
		Status:     entity.ExecutionRunning,
		StartedAt:  time.Now(),
		SafeMode:   safeMode,
	}

	if err := s.resultRepo.CreateExecution(ctx, execution); err != nil {
		return nil, fmt.Errorf("failed to create execution: %w", err)
	}

	tasks, err := s.createTasksForExecution(ctx, execution.ID, plan.Tasks, agentMap)
	if err != nil {
		return nil, err
	}

	return &ExecutionWithTasks{
		Execution: execution,
		Tasks:     tasks,
	}, nil
}

// loadAndValidateAgents loads agents and validates they exist and are online
func (s *ExecutionService) loadAndValidateAgents(
	ctx context.Context,
	agentPaws []string,
) (map[string]*entity.Agent, []*entity.Agent, error) {
	agents, err := s.agentRepo.FindByPaws(ctx, agentPaws)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load agents: %w", err)
	}

	agentMap := make(map[string]*entity.Agent, len(agents))
	for _, agent := range agents {
		agentMap[agent.Paw] = agent
	}

	for _, paw := range agentPaws {
		agent, found := agentMap[paw]
		if !found {
			return nil, nil, fmt.Errorf("agent %s not found", paw)
		}
		if agent.Status != entity.AgentOnline {
			return nil, nil, fmt.Errorf("agent %s is not online", paw)
		}
	}

	return agentMap, agents, nil
}

// createTasksForExecution creates task results and dispatch info for each planned task
func (s *ExecutionService) createTasksForExecution(
	ctx context.Context,
	executionID string,
	planTasks []service.PlannedTask,
	agentMap map[string]*entity.Agent,
) ([]TaskDispatchInfo, error) {
	tasks := make([]TaskDispatchInfo, 0, len(planTasks))

	for _, task := range planTasks {
		result := &entity.ExecutionResult{
			ID:           uuid.New().String(),
			ExecutionID:  executionID,
			TechniqueID:  task.TechniqueID,
			AgentPaw:     task.AgentPaw,
			ExecutorName: task.ExecutorName,
			Command:      task.Command,
			Status:       entity.StatusPending,
			StartedAt:    time.Now(),
		}

		if err := s.resultRepo.CreateResult(ctx, result); err != nil {
			return nil, fmt.Errorf("failed to create result: %w", err)
		}

		// Use executor type from planning phase; fallback to determineExecutor
		executor := task.ExecutorType
		if executor == "" {
			executor = s.determineExecutor(ctx, task.TechniqueID, agentMap[task.AgentPaw])
		}

		tasks = append(tasks, TaskDispatchInfo{
			ResultID:     result.ID,
			AgentPaw:     task.AgentPaw,
			TechniqueID:  task.TechniqueID,
			Command:      task.Command,
			Executor:     executor,
			ExecutorName: task.ExecutorName,
			Timeout:      task.Timeout,
			Cleanup:      task.Cleanup,
		})
	}

	return tasks, nil
}

// determineExecutor finds the appropriate executor for a technique on an agent
func (s *ExecutionService) determineExecutor(ctx context.Context, techniqueID string, agent *entity.Agent) string {
	if agent == nil {
		return "sh"
	}

	technique, err := s.techniqueRepo.FindByID(ctx, techniqueID)
	if err != nil || technique == nil {
		// Return platform-appropriate fallback when technique lookup fails
		return defaultExecutorForPlatform(agent.Platform)
	}

	if exec := technique.GetExecutorForPlatform(agent.Platform, agent.Executors); exec != nil {
		return exec.Type
	}

	// No compatible executor found - return platform-appropriate fallback
	return defaultExecutorForPlatform(agent.Platform)
}

// defaultExecutorForPlatform returns the default shell executor for a platform
func defaultExecutorForPlatform(platform string) string {
	switch platform {
	case "windows":
		return "cmd"
	case "darwin":
		return "bash"
	default:
		return "sh"
	}
}

// UpdateResult updates an execution result
func (s *ExecutionService) UpdateResult(
	ctx context.Context,
	resultID string,
	status entity.ResultStatus,
	output string,
	detected bool,
) error {
	result, err := s.resultRepo.FindResultByID(ctx, resultID)
	if err != nil {
		return err
	}

	now := time.Now()
	result.Status = status
	result.Output = output
	result.Detected = detected
	result.CompletedAt = &now
	return s.resultRepo.UpdateResult(ctx, result)
}

// UpdateResultByID updates a result by its ID with exit code
// If agentPaw is provided, it validates that the result belongs to the specified agent
func (s *ExecutionService) UpdateResultByID(
	ctx context.Context,
	resultID string,
	status entity.ResultStatus,
	output string,
	exitCode int,
	agentPaw string,
) error {
	result, err := s.resultRepo.FindResultByID(ctx, resultID)
	if err != nil {
		return fmt.Errorf("result not found: %w", err)
	}

	// Validate that the result belongs to the requesting agent
	if agentPaw != "" && result.AgentPaw != agentPaw {
		return fmt.Errorf("agent %s is not authorized to update result %s (belongs to %s)", agentPaw, resultID, result.AgentPaw)
	}

	executionID := result.ExecutionID

	now := time.Now()
	result.Status = status
	result.Output = output
	result.ExitCode = exitCode
	result.CompletedAt = &now

	if err := s.resultRepo.UpdateResult(ctx, result); err != nil {
		return err
	}

	// Check if all results are completed and auto-complete execution
	return s.checkAndCompleteExecution(ctx, executionID)
}

// checkAndCompleteExecution checks if all results are done and completes the execution
func (s *ExecutionService) checkAndCompleteExecution(ctx context.Context, executionID string) error {
	results, err := s.resultRepo.FindResultsByExecution(ctx, executionID)
	if err != nil {
		return nil // Don't fail the result update if we can't check
	}

	// Check if all results are completed (not pending or running)
	allDone := true
	for _, r := range results {
		if r.Status == entity.StatusPending || r.Status == entity.StatusRunning {
			allDone = false
			break
		}
	}

	if allDone && len(results) > 0 {
		// All tasks completed, mark execution as completed
		return s.CompleteExecution(ctx, executionID)
	}

	return nil
}

// CompleteExecution marks an execution as completed and calculates score
func (s *ExecutionService) CompleteExecution(ctx context.Context, executionID string) error {
	if s.calculator == nil {
		return fmt.Errorf("cannot complete execution: score calculator is not configured")
	}

	execution, err := s.resultRepo.FindExecutionByID(ctx, executionID)
	if err != nil {
		return err
	}

	results, err := s.resultRepo.FindResultsByExecution(ctx, executionID)
	if err != nil {
		return err
	}

	// Calculate score
	now := time.Now()
	score := s.calculator.CalculateScore(results)
	execution.Score = score
	execution.Status = entity.ExecutionCompleted
	execution.CompletedAt = &now

	return s.resultRepo.UpdateExecution(ctx, execution)
}

// GetExecution retrieves an execution by ID
func (s *ExecutionService) GetExecution(ctx context.Context, id string) (*entity.Execution, error) {
	return s.resultRepo.FindExecutionByID(ctx, id)
}

// GetExecutionResults retrieves results for an execution
func (s *ExecutionService) GetExecutionResults(ctx context.Context, executionID string) ([]*entity.ExecutionResult, error) {
	return s.resultRepo.FindResultsByExecution(ctx, executionID)
}

// GetRecentExecutions retrieves recent executions
func (s *ExecutionService) GetRecentExecutions(ctx context.Context, limit int) ([]*entity.Execution, error) {
	return s.resultRepo.FindRecentExecutions(ctx, limit)
}

// CancelExecution stops a running execution
func (s *ExecutionService) CancelExecution(ctx context.Context, executionID string) error {
	execution, err := s.resultRepo.FindExecutionByID(ctx, executionID)
	if err != nil {
		return fmt.Errorf("execution not found: %w", err)
	}

	// Only running or pending executions can be cancelled
	if execution.Status != entity.ExecutionRunning && execution.Status != entity.ExecutionPending {
		return fmt.Errorf("execution cannot be cancelled: status is %s", execution.Status)
	}

	// Update all pending results to skipped
	results, err := s.resultRepo.FindResultsByExecution(ctx, executionID)
	if err != nil {
		return fmt.Errorf("failed to get results: %w", err)
	}

	now := time.Now()
	for _, result := range results {
		if result.Status == entity.StatusPending || result.Status == entity.StatusRunning {
			result.Status = entity.StatusSkipped
			result.CompletedAt = &now
			if err := s.resultRepo.UpdateResult(ctx, result); err != nil {
				return fmt.Errorf("failed to update result %s: %w", result.ID, err)
			}
		}
	}

	// Mark execution as cancelled
	execution.Status = entity.ExecutionCancelled
	execution.CompletedAt = &now

	return s.resultRepo.UpdateExecution(ctx, execution)
}
