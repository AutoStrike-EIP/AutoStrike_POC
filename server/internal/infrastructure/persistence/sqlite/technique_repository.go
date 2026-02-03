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

// SQL column constants for techniques
const (
	techniqueColumns = "id, name, description, tactic, platforms, executors, detection, is_safe"
)

// TechniqueRepository implements repository.TechniqueRepository using SQLite
type TechniqueRepository struct {
	db *sql.DB
}

// NewTechniqueRepository creates a new SQLite technique repository
func NewTechniqueRepository(db *sql.DB) *TechniqueRepository {
	return &TechniqueRepository{db: db}
}

// Create creates a new technique
func (r *TechniqueRepository) Create(ctx context.Context, technique *entity.Technique) error {
	platforms, err := json.Marshal(technique.Platforms)
	if err != nil {
		return fmt.Errorf("failed to marshal platforms: %w", err)
	}
	executors, err := json.Marshal(technique.Executors)
	if err != nil {
		return fmt.Errorf("failed to marshal executors: %w", err)
	}
	detection, err := json.Marshal(technique.Detection)
	if err != nil {
		return fmt.Errorf("failed to marshal detection: %w", err)
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO techniques (id, name, description, tactic, platforms, executors, detection, is_safe, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, technique.ID, technique.Name, technique.Description, technique.Tactic, platforms, executors, detection, technique.IsSafe, time.Now())

	return err
}

// Update updates an existing technique
func (r *TechniqueRepository) Update(ctx context.Context, technique *entity.Technique) error {
	platforms, err := json.Marshal(technique.Platforms)
	if err != nil {
		return fmt.Errorf("failed to marshal platforms: %w", err)
	}
	executors, err := json.Marshal(technique.Executors)
	if err != nil {
		return fmt.Errorf("failed to marshal executors: %w", err)
	}
	detection, err := json.Marshal(technique.Detection)
	if err != nil {
		return fmt.Errorf("failed to marshal detection: %w", err)
	}

	_, err = r.db.ExecContext(ctx, `
		UPDATE techniques SET name = ?, description = ?, tactic = ?, platforms = ?, executors = ?, detection = ?, is_safe = ?
		WHERE id = ?
	`, technique.Name, technique.Description, technique.Tactic, platforms, executors, detection, technique.IsSafe, technique.ID)

	return err
}

// Delete deletes a technique
func (r *TechniqueRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM techniques WHERE id = ?", id)
	return err
}

// FindByID finds a technique by ID
func (r *TechniqueRepository) FindByID(ctx context.Context, id string) (*entity.Technique, error) {
	technique := &entity.Technique{}
	var platforms, executors, detection string

	err := r.db.QueryRowContext(ctx,
		fmt.Sprintf("SELECT %s FROM techniques WHERE id = ?", techniqueColumns),
		id).Scan(&technique.ID, &technique.Name, &technique.Description, &technique.Tactic, &platforms, &executors, &detection, &technique.IsSafe)

	if err != nil {
		return nil, err
	}

	// Parse JSON fields, default to empty on error
	if json.Unmarshal([]byte(platforms), &technique.Platforms) != nil {
		technique.Platforms = []string{}
	}
	if json.Unmarshal([]byte(executors), &technique.Executors) != nil {
		technique.Executors = []entity.Executor{}
	}
	if json.Unmarshal([]byte(detection), &technique.Detection) != nil {
		technique.Detection = []entity.Detection{}
	}

	return technique, nil
}

// FindAll finds all techniques
func (r *TechniqueRepository) FindAll(ctx context.Context) ([]*entity.Technique, error) {
	rows, err := r.db.QueryContext(ctx,
		fmt.Sprintf("SELECT %s FROM techniques ORDER BY id", techniqueColumns))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTechniques(rows)
}

// FindByTactic finds techniques by tactic
func (r *TechniqueRepository) FindByTactic(ctx context.Context, tactic entity.TacticType) ([]*entity.Technique, error) {
	rows, err := r.db.QueryContext(ctx,
		fmt.Sprintf("SELECT %s FROM techniques WHERE tactic = ? ORDER BY id", techniqueColumns),
		tactic)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTechniques(rows)
}

// FindByPlatform finds techniques by platform
func (r *TechniqueRepository) FindByPlatform(ctx context.Context, platform string) ([]*entity.Technique, error) {
	rows, err := r.db.QueryContext(ctx,
		fmt.Sprintf("SELECT %s FROM techniques WHERE platforms LIKE ? ORDER BY id", techniqueColumns),
		"%"+platform+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTechniques(rows)
}

// ImportFromYAML imports techniques from a YAML file
func (r *TechniqueRepository) ImportFromYAML(ctx context.Context, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var techniques []*entity.Technique
	if err := yaml.Unmarshal(data, &techniques); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	for _, t := range techniques {
		if err := r.upsert(ctx, t); err != nil {
			return fmt.Errorf("failed to import technique %s: %w", t.ID, err)
		}
	}

	return nil
}

// upsert inserts or updates a technique
func (r *TechniqueRepository) upsert(ctx context.Context, technique *entity.Technique) error {
	platforms, err := json.Marshal(technique.Platforms)
	if err != nil {
		return fmt.Errorf("failed to marshal platforms: %w", err)
	}
	executors, err := json.Marshal(technique.Executors)
	if err != nil {
		return fmt.Errorf("failed to marshal executors: %w", err)
	}
	detection, err := json.Marshal(technique.Detection)
	if err != nil {
		return fmt.Errorf("failed to marshal detection: %w", err)
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO techniques (id, name, description, tactic, platforms, executors, detection, is_safe, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			name = excluded.name,
			description = excluded.description,
			tactic = excluded.tactic,
			platforms = excluded.platforms,
			executors = excluded.executors,
			detection = excluded.detection,
			is_safe = excluded.is_safe
	`, technique.ID, technique.Name, technique.Description, technique.Tactic, platforms, executors, detection, technique.IsSafe, time.Now())

	return err
}

func (r *TechniqueRepository) scanTechniques(rows *sql.Rows) ([]*entity.Technique, error) {
	var techniques []*entity.Technique

	for rows.Next() {
		technique := &entity.Technique{}
		var platforms, executors, detection string

		err := rows.Scan(&technique.ID, &technique.Name, &technique.Description, &technique.Tactic, &platforms, &executors, &detection, &technique.IsSafe)
		if err != nil {
			return nil, err
		}

		// Parse JSON fields, default to empty on error
		if json.Unmarshal([]byte(platforms), &technique.Platforms) != nil {
			technique.Platforms = []string{}
		}
		if json.Unmarshal([]byte(executors), &technique.Executors) != nil {
			technique.Executors = []entity.Executor{}
		}
		if json.Unmarshal([]byte(detection), &technique.Detection) != nil {
			technique.Detection = []entity.Detection{}
		}

		techniques = append(techniques, technique)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return techniques, nil
}
