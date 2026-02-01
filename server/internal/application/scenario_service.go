package application

import (
	"context"
	"time"

	"autostrike/internal/domain/entity"
	"autostrike/internal/domain/repository"
	"autostrike/internal/domain/service"

	"github.com/google/uuid"
)

// ScenarioService handles scenario-related business logic
type ScenarioService struct {
	repo      repository.ScenarioRepository
	techRepo  repository.TechniqueRepository
	validator *service.TechniqueValidator
}

// NewScenarioService creates a new scenario service
func NewScenarioService(
	repo repository.ScenarioRepository,
	techRepo repository.TechniqueRepository,
	validator *service.TechniqueValidator,
) *ScenarioService {
	return &ScenarioService{
		repo:      repo,
		techRepo:  techRepo,
		validator: validator,
	}
}

// GetScenario retrieves a scenario by ID
func (s *ScenarioService) GetScenario(ctx context.Context, id string) (*entity.Scenario, error) {
	return s.repo.FindByID(ctx, id)
}

// GetAllScenarios retrieves all scenarios
func (s *ScenarioService) GetAllScenarios(ctx context.Context) ([]*entity.Scenario, error) {
	return s.repo.FindAll(ctx)
}

// GetScenariosByTag retrieves scenarios by tag
func (s *ScenarioService) GetScenariosByTag(ctx context.Context, tag string) ([]*entity.Scenario, error) {
	return s.repo.FindByTag(ctx, tag)
}

// CreateScenario creates a new scenario
func (s *ScenarioService) CreateScenario(ctx context.Context, scenario *entity.Scenario) error {
	scenario.ID = uuid.New().String()
	scenario.CreatedAt = time.Now()
	scenario.UpdatedAt = time.Now()

	// Validate scenario
	techniques, err := s.techRepo.FindAll(ctx)
	if err != nil {
		return err
	}

	result := s.validator.ValidateScenario(scenario, techniques)
	if !result.IsValid {
		return &ValidationError{Errors: result.Errors}
	}

	return s.repo.Create(ctx, scenario)
}

// UpdateScenario updates an existing scenario
func (s *ScenarioService) UpdateScenario(ctx context.Context, scenario *entity.Scenario) error {
	scenario.UpdatedAt = time.Now()

	// Validate scenario
	techniques, err := s.techRepo.FindAll(ctx)
	if err != nil {
		return err
	}

	result := s.validator.ValidateScenario(scenario, techniques)
	if !result.IsValid {
		return &ValidationError{Errors: result.Errors}
	}

	return s.repo.Update(ctx, scenario)
}

// DeleteScenario deletes a scenario
func (s *ScenarioService) DeleteScenario(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// ValidationError represents validation errors
type ValidationError struct {
	Errors []string
}

func (e *ValidationError) Error() string {
	if len(e.Errors) > 0 {
		return e.Errors[0]
	}
	return "validation failed"
}
