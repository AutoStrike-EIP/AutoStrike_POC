package service

import (
	"context"
	"fmt"

	"autostrike/internal/domain/entity"
	"autostrike/internal/domain/repository"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// AttackOrchestrator coordinates attack execution across agents
type AttackOrchestrator struct {
	agentRepo     repository.AgentRepository
	techniqueRepo repository.TechniqueRepository
	validator     *TechniqueValidator
	logger        *zap.Logger
}

// NewAttackOrchestrator creates a new orchestrator instance
func NewAttackOrchestrator(
	agentRepo repository.AgentRepository,
	techniqueRepo repository.TechniqueRepository,
	validator *TechniqueValidator,
	logger *zap.Logger,
) *AttackOrchestrator {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &AttackOrchestrator{
		agentRepo:     agentRepo,
		techniqueRepo: techniqueRepo,
		validator:     validator,
		logger:        logger,
	}
}

// ExecutionPlan represents a planned execution
type ExecutionPlan struct {
	ID    string
	Tasks []PlannedTask
}

// PlannedTask represents a single task in the execution plan
type PlannedTask struct {
	TechniqueID  string
	AgentPaw     string
	Phase        string
	Order        int
	Command      string
	Cleanup      string
	Timeout      int
	ExecutorName string // Name of the executor (for display/debugging)
	ExecutorType string // Type of executor (sh, bash, cmd, powershell)
}

// PlanExecution creates an execution plan for a scenario
func (o *AttackOrchestrator) PlanExecution(
	ctx context.Context,
	scenario *entity.Scenario,
	targetAgents []*entity.Agent,
	safeMode bool,
) (*ExecutionPlan, error) {
	plan := &ExecutionPlan{
		ID:    uuid.New().String(),
		Tasks: make([]PlannedTask, 0),
	}

	taskOrder := 0
	for _, phase := range scenario.Phases {
		tasks := o.planPhase(ctx, phase, targetAgents, safeMode, taskOrder)
		plan.Tasks = append(plan.Tasks, tasks...)
		taskOrder += len(tasks)
	}

	if len(plan.Tasks) == 0 {
		return nil, fmt.Errorf("no executable tasks for the given scenario and agents")
	}

	return plan, nil
}

// planPhase creates tasks for a single phase
func (o *AttackOrchestrator) planPhase(
	ctx context.Context,
	phase entity.Phase,
	targetAgents []*entity.Agent,
	safeMode bool,
	startOrder int,
) []PlannedTask {
	var tasks []PlannedTask
	taskOrder := startOrder

	for _, selection := range phase.Techniques {
		technique := o.getTechnique(ctx, selection.TechniqueID, safeMode)
		if technique == nil {
			continue
		}

		for _, agent := range targetAgents {
			agentTasks := o.createTasksForAgent(agent, technique, selection.ExecutorName, phase.Name, taskOrder)
			tasks = append(tasks, agentTasks...)
			taskOrder += len(agentTasks)
		}
	}

	return tasks
}

// getTechnique retrieves and validates a technique.
// In safe mode, filters out unsafe executors instead of skipping the whole technique.
func (o *AttackOrchestrator) getTechnique(ctx context.Context, techID string, safeMode bool) *entity.Technique {
	technique, err := o.techniqueRepo.FindByID(ctx, techID)
	if err != nil {
		o.logger.Warn("Skipping technique: not found in repository", zap.String("technique_id", techID))
		return nil
	}

	if safeMode {
		if !technique.IsSafe {
			o.logger.Info("Skipping technique with no safe executors", zap.String("technique_id", techID))
			return nil
		}
		// Filter to only safe executors
		var safeExecutors []entity.Executor
		for _, exec := range technique.Executors {
			if exec.IsSafe {
				safeExecutors = append(safeExecutors, exec)
			}
		}
		if len(safeExecutors) == 0 {
			// Backward compatibility: if technique is marked safe but no executor has
			// individual is_safe set (legacy format), use all executors
			o.logger.Info("Safe technique has no per-executor is_safe flags, using all executors", zap.String("technique_id", techID))
			return technique
		}
		// Return a copy with only safe executors
		filtered := *technique
		filtered.Executors = safeExecutors
		return &filtered
	}

	return technique
}

// createTasksForAgent creates tasks for all compatible executors on the agent.
// If executorName is specified, only that executor is used (single task).
// If executorName is empty, ALL compatible executors are used (multiple tasks).
func (o *AttackOrchestrator) createTasksForAgent(
	agent *entity.Agent,
	technique *entity.Technique,
	executorName string,
	phaseName string,
	startOrder int,
) []PlannedTask {
	if !agent.IsCompatible(technique) {
		return nil
	}

	// If a specific executor is requested, use only that one
	if executorName != "" {
		executor := technique.GetExecutorByName(executorName, agent.Platform, agent.Executors)
		if executor == nil {
			// Fallback to auto-select first compatible executor
			executor = technique.GetExecutorForPlatform(agent.Platform, agent.Executors)
		}
		if executor == nil {
			return nil
		}
		return []PlannedTask{{
			TechniqueID:  technique.ID,
			AgentPaw:     agent.Paw,
			Phase:        phaseName,
			Order:        startOrder,
			Command:      executor.Command,
			Cleanup:      executor.Cleanup,
			Timeout:      executor.Timeout,
			ExecutorName: executor.Name,
			ExecutorType: executor.Type,
		}}
	}

	// No executor name: run ALL compatible executors
	executors := technique.GetExecutorsForPlatform(agent.Platform, agent.Executors)
	if len(executors) == 0 {
		return nil
	}

	tasks := make([]PlannedTask, 0, len(executors))
	for i, exec := range executors {
		tasks = append(tasks, PlannedTask{
			TechniqueID:  technique.ID,
			AgentPaw:     agent.Paw,
			Phase:        phaseName,
			Order:        startOrder + i,
			Command:      exec.Command,
			Cleanup:      exec.Cleanup,
			Timeout:      exec.Timeout,
			ExecutorName: exec.Name,
			ExecutorType: exec.Type,
		})
	}
	return tasks
}

// ValidatePlan validates an execution plan
func (o *AttackOrchestrator) ValidatePlan(ctx context.Context, plan *ExecutionPlan) error {
	for _, task := range plan.Tasks {
		agent, err := o.agentRepo.FindByPaw(ctx, task.AgentPaw)
		if err != nil {
			return fmt.Errorf("agent %s not found", task.AgentPaw)
		}

		if agent.Status != entity.AgentOnline {
			return fmt.Errorf("agent %s is not online", task.AgentPaw)
		}

		technique, err := o.techniqueRepo.FindByID(ctx, task.TechniqueID)
		if err != nil {
			return fmt.Errorf("technique %s not found", task.TechniqueID)
		}

		if !agent.IsCompatible(technique) {
			return fmt.Errorf("agent %s is not compatible with technique %s", task.AgentPaw, task.TechniqueID)
		}
	}

	return nil
}
