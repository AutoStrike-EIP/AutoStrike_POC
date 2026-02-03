package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"autostrike/internal/domain/entity"

	"gopkg.in/yaml.v3"
)

// SQL column constants for scenarios
const (
	scenarioColumns = "id, name, description, phases, tags, created_at, updated_at"
)

// ScenarioRepository implements repository.ScenarioRepository using SQLite
type ScenarioRepository struct {
	db *sql.DB
}

// NewScenarioRepository creates a new SQLite scenario repository
func NewScenarioRepository(db *sql.DB) *ScenarioRepository {
	return &ScenarioRepository{db: db}
}

// Create creates a new scenario
func (r *ScenarioRepository) Create(ctx context.Context, scenario *entity.Scenario) error {
	phases, err := json.Marshal(scenario.Phases)
	if err != nil {
		return fmt.Errorf("failed to marshal phases: %w", err)
	}
	tags, err := json.Marshal(scenario.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO scenarios (id, name, description, phases, tags, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, scenario.ID, scenario.Name, scenario.Description, phases, tags, scenario.CreatedAt, scenario.UpdatedAt)

	return err
}

// Update updates an existing scenario
func (r *ScenarioRepository) Update(ctx context.Context, scenario *entity.Scenario) error {
	phases, err := json.Marshal(scenario.Phases)
	if err != nil {
		return fmt.Errorf("failed to marshal phases: %w", err)
	}
	tags, err := json.Marshal(scenario.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	_, err = r.db.ExecContext(ctx, `
		UPDATE scenarios SET name = ?, description = ?, phases = ?, tags = ?, updated_at = ?
		WHERE id = ?
	`, scenario.Name, scenario.Description, phases, tags, time.Now(), scenario.ID)

	return err
}

// Delete deletes a scenario
func (r *ScenarioRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM scenarios WHERE id = ?", id)
	return err
}

// FindByID finds a scenario by ID
func (r *ScenarioRepository) FindByID(ctx context.Context, id string) (*entity.Scenario, error) {
	scenario := &entity.Scenario{}
	var phases, tags string

	err := r.db.QueryRowContext(ctx,
		fmt.Sprintf("SELECT %s FROM scenarios WHERE id = ?", scenarioColumns),
		id).Scan(&scenario.ID, &scenario.Name, &scenario.Description, &phases, &tags, &scenario.CreatedAt, &scenario.UpdatedAt)

	if err != nil {
		return nil, err
	}

	// Parse JSON fields, default to empty on error
	if json.Unmarshal([]byte(phases), &scenario.Phases) != nil {
		scenario.Phases = []entity.Phase{}
	}
	if json.Unmarshal([]byte(tags), &scenario.Tags) != nil {
		scenario.Tags = []string{}
	}

	return scenario, nil
}

// FindAll finds all scenarios
func (r *ScenarioRepository) FindAll(ctx context.Context) ([]*entity.Scenario, error) {
	rows, err := r.db.QueryContext(ctx,
		fmt.Sprintf("SELECT %s FROM scenarios ORDER BY updated_at DESC", scenarioColumns))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanScenarios(rows)
}

// FindByTag finds scenarios by tag
func (r *ScenarioRepository) FindByTag(ctx context.Context, tag string) ([]*entity.Scenario, error) {
	rows, err := r.db.QueryContext(ctx,
		fmt.Sprintf("SELECT %s FROM scenarios WHERE tags LIKE ? ORDER BY updated_at DESC", scenarioColumns),
		"%"+tag+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanScenarios(rows)
}

func (r *ScenarioRepository) scanScenarios(rows *sql.Rows) ([]*entity.Scenario, error) {
	var scenarios []*entity.Scenario

	for rows.Next() {
		scenario := &entity.Scenario{}
		var phases, tags string

		err := rows.Scan(&scenario.ID, &scenario.Name, &scenario.Description, &phases, &tags, &scenario.CreatedAt, &scenario.UpdatedAt)
		if err != nil {
			return nil, err
		}

		// Parse JSON fields, default to empty on error
		if json.Unmarshal([]byte(phases), &scenario.Phases) != nil {
			scenario.Phases = []entity.Phase{}
		}
		if json.Unmarshal([]byte(tags), &scenario.Tags) != nil {
			scenario.Tags = []string{}
		}

		scenarios = append(scenarios, scenario)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return scenarios, nil
}

// ImportFromYAML imports scenarios from a YAML file
func (r *ScenarioRepository) ImportFromYAML(ctx context.Context, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var scenarios []*entity.Scenario
	if err := yaml.Unmarshal(data, &scenarios); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	now := time.Now()
	for _, s := range scenarios {
		if s.CreatedAt.IsZero() {
			s.CreatedAt = now
		}
		if s.UpdatedAt.IsZero() {
			s.UpdatedAt = now
		}
		if err := r.upsert(ctx, s); err != nil {
			return fmt.Errorf("failed to import scenario %s: %w", s.ID, err)
		}
	}

	return nil
}

// upsert inserts or updates a scenario
func (r *ScenarioRepository) upsert(ctx context.Context, scenario *entity.Scenario) error {
	phases, err := json.Marshal(scenario.Phases)
	if err != nil {
		return fmt.Errorf("failed to marshal phases: %w", err)
	}
	tags, err := json.Marshal(scenario.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO scenarios (id, name, description, phases, tags, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			name = excluded.name,
			description = excluded.description,
			phases = excluded.phases,
			tags = excluded.tags,
			updated_at = excluded.updated_at
	`, scenario.ID, scenario.Name, scenario.Description, phases, tags, scenario.CreatedAt, scenario.UpdatedAt)

	return err
}
