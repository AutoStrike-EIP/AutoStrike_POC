package application

import (
	"context"

	"autostrike/internal/domain/entity"
	"autostrike/internal/domain/repository"
)

// TechniqueService handles technique-related business logic
type TechniqueService struct {
	repo repository.TechniqueRepository
}

// NewTechniqueService creates a new technique service
func NewTechniqueService(repo repository.TechniqueRepository) *TechniqueService {
	return &TechniqueService{repo: repo}
}

// GetTechnique retrieves a technique by ID
func (s *TechniqueService) GetTechnique(ctx context.Context, id string) (*entity.Technique, error) {
	return s.repo.FindByID(ctx, id)
}

// GetAllTechniques retrieves all techniques
func (s *TechniqueService) GetAllTechniques(ctx context.Context) ([]*entity.Technique, error) {
	return s.repo.FindAll(ctx)
}

// GetTechniquesByTactic retrieves techniques by MITRE tactic
func (s *TechniqueService) GetTechniquesByTactic(ctx context.Context, tactic entity.TacticType) ([]*entity.Technique, error) {
	return s.repo.FindByTactic(ctx, tactic)
}

// GetTechniquesByPlatform retrieves techniques by platform
func (s *TechniqueService) GetTechniquesByPlatform(ctx context.Context, platform string) ([]*entity.Technique, error) {
	return s.repo.FindByPlatform(ctx, platform)
}

// ImportTechniques imports techniques from YAML file
func (s *TechniqueService) ImportTechniques(ctx context.Context, path string) error {
	return s.repo.ImportFromYAML(ctx, path)
}

// CreateTechnique creates a new technique
func (s *TechniqueService) CreateTechnique(ctx context.Context, technique *entity.Technique) error {
	return s.repo.Create(ctx, technique)
}

// UpdateTechnique updates an existing technique
func (s *TechniqueService) UpdateTechnique(ctx context.Context, technique *entity.Technique) error {
	return s.repo.Update(ctx, technique)
}

// DeleteTechnique deletes a technique
func (s *TechniqueService) DeleteTechnique(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// GetCoverage returns MITRE ATT&CK coverage statistics
func (s *TechniqueService) GetCoverage(ctx context.Context) (map[entity.TacticType]int, error) {
	techniques, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	coverage := make(map[entity.TacticType]int)
	for _, t := range techniques {
		for _, tactic := range t.GetTactics() {
			coverage[tactic]++
		}
	}

	return coverage, nil
}
