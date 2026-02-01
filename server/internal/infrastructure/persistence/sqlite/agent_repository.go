package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"autostrike/internal/domain/entity"
)

// AgentRepository implements repository.AgentRepository using SQLite
type AgentRepository struct {
	db *sql.DB
}

// NewAgentRepository creates a new SQLite agent repository
func NewAgentRepository(db *sql.DB) *AgentRepository {
	return &AgentRepository{db: db}
}

// Create creates a new agent
func (r *AgentRepository) Create(ctx context.Context, agent *entity.Agent) error {
	executors, err := json.Marshal(agent.Executors)
	if err != nil {
		return fmt.Errorf("failed to marshal executors: %w", err)
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO agents (paw, hostname, username, platform, executors, status, last_seen, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, agent.Paw, agent.Hostname, agent.Username, agent.Platform, executors, agent.Status, agent.LastSeen, agent.CreatedAt)

	return err
}

// Update updates an existing agent
func (r *AgentRepository) Update(ctx context.Context, agent *entity.Agent) error {
	executors, err := json.Marshal(agent.Executors)
	if err != nil {
		return fmt.Errorf("failed to marshal executors: %w", err)
	}

	_, err = r.db.ExecContext(ctx, `
		UPDATE agents SET hostname = ?, username = ?, platform = ?, executors = ?, status = ?, last_seen = ?
		WHERE paw = ?
	`, agent.Hostname, agent.Username, agent.Platform, executors, agent.Status, agent.LastSeen, agent.Paw)

	return err
}

// Delete deletes an agent
func (r *AgentRepository) Delete(ctx context.Context, paw string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM agents WHERE paw = ?", paw)
	return err
}

// FindByPaw finds an agent by paw
func (r *AgentRepository) FindByPaw(ctx context.Context, paw string) (*entity.Agent, error) {
	agent := &entity.Agent{}
	var executors string

	err := r.db.QueryRowContext(ctx, `
		SELECT paw, hostname, username, platform, executors, status, last_seen, created_at
		FROM agents WHERE paw = ?
	`, paw).Scan(&agent.Paw, &agent.Hostname, &agent.Username, &agent.Platform, &executors, &agent.Status, &agent.LastSeen, &agent.CreatedAt)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(executors), &agent.Executors); err != nil {
		agent.Executors = []string{} // Default to empty on parse error
	}
	return agent, nil
}

// FindAll finds all agents
func (r *AgentRepository) FindAll(ctx context.Context) ([]*entity.Agent, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT paw, hostname, username, platform, executors, status, last_seen, created_at
		FROM agents ORDER BY last_seen DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanAgents(rows)
}

// FindByStatus finds agents by status
func (r *AgentRepository) FindByStatus(ctx context.Context, status entity.AgentStatus) ([]*entity.Agent, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT paw, hostname, username, platform, executors, status, last_seen, created_at
		FROM agents WHERE status = ? ORDER BY last_seen DESC
	`, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanAgents(rows)
}

// FindByPlatform finds agents by platform
func (r *AgentRepository) FindByPlatform(ctx context.Context, platform string) ([]*entity.Agent, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT paw, hostname, username, platform, executors, status, last_seen, created_at
		FROM agents WHERE platform = ? ORDER BY last_seen DESC
	`, platform)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanAgents(rows)
}

// UpdateLastSeen updates the last seen timestamp
func (r *AgentRepository) UpdateLastSeen(ctx context.Context, paw string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE agents SET last_seen = ?, status = ? WHERE paw = ?
	`, time.Now(), entity.AgentOnline, paw)
	return err
}

func (r *AgentRepository) scanAgents(rows *sql.Rows) ([]*entity.Agent, error) {
	var agents []*entity.Agent

	for rows.Next() {
		agent := &entity.Agent{}
		var executors string

		err := rows.Scan(&agent.Paw, &agent.Hostname, &agent.Username, &agent.Platform, &executors, &agent.Status, &agent.LastSeen, &agent.CreatedAt)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(executors), &agent.Executors); err != nil {
			agent.Executors = []string{} // Default to empty on parse error
		}
		agents = append(agents, agent)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return agents, nil
}
