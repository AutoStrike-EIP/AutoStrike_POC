package application

import (
	"context"
	"errors"
	"testing"

	"autostrike/internal/domain/entity"
)

type mockTechniqueRepo struct {
	techniques map[string]*entity.Technique
	err        error
}

func newMockTechniqueRepo() *mockTechniqueRepo {
	return &mockTechniqueRepo{techniques: make(map[string]*entity.Technique)}
}

func (m *mockTechniqueRepo) Create(ctx context.Context, t *entity.Technique) error {
	if m.err != nil {
		return m.err
	}
	m.techniques[t.ID] = t
	return nil
}

func (m *mockTechniqueRepo) Update(ctx context.Context, t *entity.Technique) error {
	if m.err != nil {
		return m.err
	}
	m.techniques[t.ID] = t
	return nil
}

func (m *mockTechniqueRepo) Delete(ctx context.Context, id string) error {
	if m.err != nil {
		return m.err
	}
	delete(m.techniques, id)
	return nil
}

func (m *mockTechniqueRepo) FindByID(ctx context.Context, id string) (*entity.Technique, error) {
	if m.err != nil {
		return nil, m.err
	}
	t, ok := m.techniques[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return t, nil
}

func (m *mockTechniqueRepo) FindAll(ctx context.Context) ([]*entity.Technique, error) {
	if m.err != nil {
		return nil, m.err
	}
	result := make([]*entity.Technique, 0, len(m.techniques))
	for _, t := range m.techniques {
		result = append(result, t)
	}
	return result, nil
}

func (m *mockTechniqueRepo) FindByTactic(ctx context.Context, tactic entity.TacticType) ([]*entity.Technique, error) {
	if m.err != nil {
		return nil, m.err
	}
	var result []*entity.Technique
	for _, t := range m.techniques {
		if t.Tactic == tactic {
			result = append(result, t)
		}
	}
	return result, nil
}

func (m *mockTechniqueRepo) FindByPlatform(ctx context.Context, platform string) ([]*entity.Technique, error) {
	if m.err != nil {
		return nil, m.err
	}
	var result []*entity.Technique
	for _, t := range m.techniques {
		for _, p := range t.Platforms {
			if p == platform {
				result = append(result, t)
				break
			}
		}
	}
	return result, nil
}

func (m *mockTechniqueRepo) ImportFromYAML(ctx context.Context, path string) error {
	return m.err
}

func TestNewTechniqueService(t *testing.T) {
	repo := newMockTechniqueRepo()
	service := NewTechniqueService(repo)
	if service == nil {
		t.Fatal("Expected non-nil service")
	}
}

func TestGetTechnique(t *testing.T) {
	repo := newMockTechniqueRepo()
	repo.techniques["T1059"] = &entity.Technique{ID: "T1059", Name: "Command"}
	service := NewTechniqueService(repo)

	tech, err := service.GetTechnique(context.Background(), "T1059")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if tech.Name != "Command" {
		t.Errorf("Expected name Command, got %s", tech.Name)
	}
}

func TestGetTechnique_NotFound(t *testing.T) {
	repo := newMockTechniqueRepo()
	service := NewTechniqueService(repo)

	_, err := service.GetTechnique(context.Background(), "invalid")
	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestGetAllTechniques(t *testing.T) {
	repo := newMockTechniqueRepo()
	repo.techniques["T1"] = &entity.Technique{ID: "T1"}
	repo.techniques["T2"] = &entity.Technique{ID: "T2"}
	service := NewTechniqueService(repo)

	techniques, err := service.GetAllTechniques(context.Background())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(techniques) != 2 {
		t.Errorf("Expected 2 techniques, got %d", len(techniques))
	}
}

func TestGetTechniquesByTactic(t *testing.T) {
	repo := newMockTechniqueRepo()
	repo.techniques["T1"] = &entity.Technique{ID: "T1", Tactic: entity.TacticExecution}
	repo.techniques["T2"] = &entity.Technique{ID: "T2", Tactic: entity.TacticDiscovery}
	service := NewTechniqueService(repo)

	techniques, err := service.GetTechniquesByTactic(context.Background(), entity.TacticExecution)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(techniques) != 1 {
		t.Errorf("Expected 1 technique, got %d", len(techniques))
	}
}

func TestGetTechniquesByPlatform(t *testing.T) {
	repo := newMockTechniqueRepo()
	repo.techniques["T1"] = &entity.Technique{ID: "T1", Platforms: []string{"linux"}}
	repo.techniques["T2"] = &entity.Technique{ID: "T2", Platforms: []string{"windows"}}
	service := NewTechniqueService(repo)

	techniques, err := service.GetTechniquesByPlatform(context.Background(), "linux")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(techniques) != 1 {
		t.Errorf("Expected 1 technique, got %d", len(techniques))
	}
}

func TestImportTechniques(t *testing.T) {
	repo := newMockTechniqueRepo()
	service := NewTechniqueService(repo)

	err := service.ImportTechniques(context.Background(), "/path/to/file.yaml")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestImportTechniques_Error(t *testing.T) {
	repo := newMockTechniqueRepo()
	repo.err = errors.New("import error")
	service := NewTechniqueService(repo)

	err := service.ImportTechniques(context.Background(), "/path/to/file.yaml")
	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestCreateTechnique(t *testing.T) {
	repo := newMockTechniqueRepo()
	service := NewTechniqueService(repo)

	tech := &entity.Technique{ID: "T1059", Name: "Command"}
	err := service.CreateTechnique(context.Background(), tech)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestUpdateTechnique(t *testing.T) {
	repo := newMockTechniqueRepo()
	repo.techniques["T1059"] = &entity.Technique{ID: "T1059", Name: "Old"}
	service := NewTechniqueService(repo)

	tech := &entity.Technique{ID: "T1059", Name: "New"}
	err := service.UpdateTechnique(context.Background(), tech)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestDeleteTechnique(t *testing.T) {
	repo := newMockTechniqueRepo()
	repo.techniques["T1059"] = &entity.Technique{ID: "T1059"}
	service := NewTechniqueService(repo)

	err := service.DeleteTechnique(context.Background(), "T1059")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestGetCoverage(t *testing.T) {
	repo := newMockTechniqueRepo()
	repo.techniques["T1"] = &entity.Technique{ID: "T1", Tactic: entity.TacticExecution}
	repo.techniques["T2"] = &entity.Technique{ID: "T2", Tactic: entity.TacticExecution}
	repo.techniques["T3"] = &entity.Technique{ID: "T3", Tactic: entity.TacticDiscovery}
	service := NewTechniqueService(repo)

	coverage, err := service.GetCoverage(context.Background())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if coverage[entity.TacticExecution] != 2 {
		t.Errorf("Expected 2 execution techniques, got %d", coverage[entity.TacticExecution])
	}
	if coverage[entity.TacticDiscovery] != 1 {
		t.Errorf("Expected 1 discovery technique, got %d", coverage[entity.TacticDiscovery])
	}
}

func TestGetCoverage_Error(t *testing.T) {
	repo := newMockTechniqueRepo()
	repo.err = errors.New("db error")
	service := NewTechniqueService(repo)

	_, err := service.GetCoverage(context.Background())
	if err == nil {
		t.Fatal("Expected error")
	}
}
