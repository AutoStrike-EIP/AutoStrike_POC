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

// StartExecution starts a new scenario execution
func (s *ExecutionService) StartExecution(
	ctx context.Context,
	scenarioID string,
	agentPaws []string,
	safeMode bool,
) (*entity.Execution, error) {
	// Load scenario
	scenario, err := s.scenarioRepo.FindByID(ctx, scenarioID)
	if err != nil {
		return nil, fmt.Errorf("scenario not found: %w", err)
	}

	// Load agents
	agents := make([]*entity.Agent, 0, len(agentPaws))
	for _, paw := range agentPaws {
		agent, err := s.agentRepo.FindByPaw(ctx, paw)
		if err != nil {
			return nil, fmt.Errorf("agent %s not found: %w", paw, err)
		}
		if agent.Status != entity.AgentOnline {
			return nil, fmt.Errorf("agent %s is not online", paw)
		}
		agents = append(agents, agent)
	}

	// Create execution plan
	plan, err := s.orchestrator.PlanExecution(ctx, scenario, agents, safeMode)
	if err != nil {
		return nil, fmt.Errorf("failed to plan execution: %w", err)
	}

	// Create execution record
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

	// Create pending results for each task
	for _, task := range plan.Tasks {
		result := &entity.ExecutionResult{
			ID:          uuid.New().String(),
			ExecutionID: execution.ID,
			TechniqueID: task.TechniqueID,
			AgentPaw:    task.AgentPaw,
			Status:      entity.StatusPending,
			StartedAt:   time.Now(),
		}

		if err := s.resultRepo.CreateResult(ctx, result); err != nil {
			return nil, fmt.Errorf("failed to create result: %w", err)
		}
	}

	return execution, nil
}

// UpdateResult updates an execution result
func (s *ExecutionService) UpdateResult(
	ctx context.Context,
	resultID string,
	status entity.ResultStatus,
	output string,
	detected bool,
) error {
	results, err := s.resultRepo.FindResultsByExecution(ctx, resultID)
	if err != nil {
		return err
	}

	for _, result := range results {
		if result.ID == resultID {
			now := time.Now()
			result.Status = status
			result.Output = output
			result.Detected = detected
			result.CompletedAt = &now
			return s.resultRepo.UpdateResult(ctx, result)
		}
	}

	return fmt.Errorf("result not found")
}

// CompleteExecution marks an execution as completed and calculates score
func (s *ExecutionService) CompleteExecution(ctx context.Context, executionID string) error {
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
