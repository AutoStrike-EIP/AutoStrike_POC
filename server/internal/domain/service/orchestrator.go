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
	TechniqueID string
	AgentPaw    string
	Phase       string
	Order       int
	Command     string
	Cleanup     string
	Timeout     int
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
			task := o.createTaskForAgent(agent, technique, selection.ExecutorName, phase.Name, taskOrder)
			if task != nil {
				tasks = append(tasks, *task)
				taskOrder++
			}
		}
	}

	return tasks
}

// getTechnique retrieves and validates a technique
func (o *AttackOrchestrator) getTechnique(ctx context.Context, techID string, safeMode bool) *entity.Technique {
	technique, err := o.techniqueRepo.FindByID(ctx, techID)
	if err != nil {
		o.logger.Warn("Skipping technique: not found in repository", zap.String("technique_id", techID))
		return nil
	}

	if safeMode && !technique.IsSafe {
		o.logger.Info("Skipping unsafe technique in safe mode", zap.String("technique_id", techID))
		return nil
	}

	return technique
}

// createTaskForAgent creates a task if the agent is compatible
func (o *AttackOrchestrator) createTaskForAgent(
	agent *entity.Agent,
	technique *entity.Technique,
	executorName string,
	phaseName string,
	order int,
) *PlannedTask {
	if !agent.IsCompatible(technique) {
		return nil
	}

	var executor *entity.Executor
	if executorName != "" {
		executor = technique.GetExecutorByName(executorName, agent.Platform, agent.Executors)
	}
	// Fallback to auto-select if no executor name or not found
	if executor == nil {
		executor = technique.GetExecutorForPlatform(agent.Platform, agent.Executors)
	}
	if executor == nil {
		return nil
	}

	return &PlannedTask{
		TechniqueID: technique.ID,
		AgentPaw:    agent.Paw,
		Phase:       phaseName,
		Order:       order,
		Command:     executor.Command,
		Cleanup:     executor.Cleanup,
		Timeout:     executor.Timeout,
	}
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
