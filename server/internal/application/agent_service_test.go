package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"autostrike/internal/domain/entity"
)

// Mock AgentRepository for testing
type mockAgentRepo struct {
	agents      map[string]*entity.Agent
	findErr     error
	createErr   error
	updateErr   error
	deleteErr   error
	lastSeenErr error
}

func newMockAgentRepo() *mockAgentRepo {
	return &mockAgentRepo{
		agents: make(map[string]*entity.Agent),
	}
}

func (m *mockAgentRepo) Create(ctx context.Context, agent *entity.Agent) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.agents[agent.Paw] = agent
	return nil
}

func (m *mockAgentRepo) Update(ctx context.Context, agent *entity.Agent) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	m.agents[agent.Paw] = agent
	return nil
}

func (m *mockAgentRepo) Delete(ctx context.Context, paw string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	delete(m.agents, paw)
	return nil
}

func (m *mockAgentRepo) FindByPaw(ctx context.Context, paw string) (*entity.Agent, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	agent, ok := m.agents[paw]
	if !ok {
		return nil, errors.New("agent not found")
	}
	return agent, nil
}

func (m *mockAgentRepo) FindAll(ctx context.Context) ([]*entity.Agent, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	result := make([]*entity.Agent, 0, len(m.agents))
	for _, agent := range m.agents {
		result = append(result, agent)
	}
	return result, nil
}

func (m *mockAgentRepo) FindByStatus(ctx context.Context, status entity.AgentStatus) ([]*entity.Agent, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	var result []*entity.Agent
	for _, agent := range m.agents {
		if agent.Status == status {
			result = append(result, agent)
		}
	}
	return result, nil
}

func (m *mockAgentRepo) FindByPlatform(ctx context.Context, platform string) ([]*entity.Agent, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	var result []*entity.Agent
	for _, agent := range m.agents {
		if agent.Platform == platform {
			result = append(result, agent)
		}
	}
	return result, nil
}

func (m *mockAgentRepo) UpdateLastSeen(ctx context.Context, paw string) error {
	if m.lastSeenErr != nil {
		return m.lastSeenErr
	}
	agent, ok := m.agents[paw]
	if ok {
		agent.LastSeen = time.Now()
	}
	return nil
}

func TestNewAgentService(t *testing.T) {
	repo := newMockAgentRepo()
	service := NewAgentService(repo)

	if service == nil {
		t.Fatal("Expected non-nil service")
	}
}

func TestRegisterAgent_NewAgent(t *testing.T) {
	repo := newMockAgentRepo()
	service := NewAgentService(repo)
	ctx := context.Background()

	agent := &entity.Agent{
		Paw:       "test-paw",
		Hostname:  "test-host",
		Platform:  "linux",
		Executors: []string{"sh", "bash"},
	}

	err := service.RegisterAgent(ctx, agent)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if agent.Status != entity.AgentOnline {
		t.Errorf("Expected status Online, got %v", agent.Status)
	}

	if agent.LastSeen.IsZero() {
		t.Error("Expected LastSeen to be set")
	}
}

func TestRegisterAgent_ExistingAgent(t *testing.T) {
	repo := newMockAgentRepo()
	existing := &entity.Agent{
		Paw:       "existing-paw",
		Hostname:  "old-host",
		Platform:  "windows",
		Status:    entity.AgentOffline,
		Executors: []string{"powershell"},
	}
	repo.agents[existing.Paw] = existing

	service := NewAgentService(repo)
	ctx := context.Background()

	updated := &entity.Agent{
		Paw:       "existing-paw",
		Hostname:  "new-host",
		Platform:  "linux",
		Executors: []string{"sh"},
	}

	err := service.RegisterAgent(ctx, updated)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	agent := repo.agents["existing-paw"]
	if agent.Status != entity.AgentOnline {
		t.Errorf("Expected status Online, got %v", agent.Status)
	}
	if agent.Platform != "linux" {
		t.Errorf("Expected platform linux, got %s", agent.Platform)
	}
}

func TestHeartbeat(t *testing.T) {
	repo := newMockAgentRepo()
	repo.agents["test-paw"] = &entity.Agent{Paw: "test-paw"}

	service := NewAgentService(repo)
	ctx := context.Background()

	err := service.Heartbeat(ctx, "test-paw")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestHeartbeat_Error(t *testing.T) {
	repo := newMockAgentRepo()
	repo.lastSeenErr = errors.New("db error")

	service := NewAgentService(repo)
	ctx := context.Background()

	err := service.Heartbeat(ctx, "test-paw")
	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestGetAgent(t *testing.T) {
	repo := newMockAgentRepo()
	expected := &entity.Agent{Paw: "test-paw", Hostname: "test-host"}
	repo.agents["test-paw"] = expected

	service := NewAgentService(repo)
	ctx := context.Background()

	agent, err := service.GetAgent(ctx, "test-paw")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if agent.Hostname != expected.Hostname {
		t.Errorf("Expected hostname %s, got %s", expected.Hostname, agent.Hostname)
	}
}

func TestGetAgent_NotFound(t *testing.T) {
	repo := newMockAgentRepo()
	service := NewAgentService(repo)
	ctx := context.Background()

	_, err := service.GetAgent(ctx, "non-existent")
	if err == nil {
		t.Fatal("Expected error for non-existent agent")
	}
}

func TestGetAllAgents(t *testing.T) {
	repo := newMockAgentRepo()
	repo.agents["paw1"] = &entity.Agent{Paw: "paw1"}
	repo.agents["paw2"] = &entity.Agent{Paw: "paw2"}

	service := NewAgentService(repo)
	ctx := context.Background()

	agents, err := service.GetAllAgents(ctx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(agents) != 2 {
		t.Errorf("Expected 2 agents, got %d", len(agents))
	}
}

func TestGetOnlineAgents(t *testing.T) {
	repo := newMockAgentRepo()
	repo.agents["online1"] = &entity.Agent{Paw: "online1", Status: entity.AgentOnline}
	repo.agents["online2"] = &entity.Agent{Paw: "online2", Status: entity.AgentOnline}
	repo.agents["offline"] = &entity.Agent{Paw: "offline", Status: entity.AgentOffline}

	service := NewAgentService(repo)
	ctx := context.Background()

	agents, err := service.GetOnlineAgents(ctx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(agents) != 2 {
		t.Errorf("Expected 2 online agents, got %d", len(agents))
	}
}

func TestMarkAgentOffline(t *testing.T) {
	repo := newMockAgentRepo()
	repo.agents["test-paw"] = &entity.Agent{Paw: "test-paw", Status: entity.AgentOnline}

	service := NewAgentService(repo)
	ctx := context.Background()

	err := service.MarkAgentOffline(ctx, "test-paw")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if repo.agents["test-paw"].Status != entity.AgentOffline {
		t.Error("Expected agent to be offline")
	}
}

func TestMarkAgentOffline_NotFound(t *testing.T) {
	repo := newMockAgentRepo()
	service := NewAgentService(repo)
	ctx := context.Background()

	err := service.MarkAgentOffline(ctx, "non-existent")
	if err == nil {
		t.Fatal("Expected error for non-existent agent")
	}
}

func TestDeleteAgent(t *testing.T) {
	repo := newMockAgentRepo()
	repo.agents["test-paw"] = &entity.Agent{Paw: "test-paw"}

	service := NewAgentService(repo)
	ctx := context.Background()

	err := service.DeleteAgent(ctx, "test-paw")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if _, exists := repo.agents["test-paw"]; exists {
		t.Error("Expected agent to be deleted")
	}
}

func TestCheckStaleAgents(t *testing.T) {
	repo := newMockAgentRepo()
	staleTime := time.Now().Add(-5 * time.Minute)
	freshTime := time.Now()

	repo.agents["stale"] = &entity.Agent{
		Paw:      "stale",
		Status:   entity.AgentOnline,
		LastSeen: staleTime,
	}
	repo.agents["fresh"] = &entity.Agent{
		Paw:      "fresh",
		Status:   entity.AgentOnline,
		LastSeen: freshTime,
	}

	service := NewAgentService(repo)
	ctx := context.Background()

	err := service.CheckStaleAgents(ctx, 2*time.Minute)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if repo.agents["stale"].Status != entity.AgentOffline {
		t.Error("Expected stale agent to be offline")
	}
	if repo.agents["fresh"].Status != entity.AgentOnline {
		t.Error("Expected fresh agent to still be online")
	}
}

func TestCheckStaleAgents_FindError(t *testing.T) {
	repo := newMockAgentRepo()
	repo.findErr = errors.New("db error")

	service := NewAgentService(repo)
	ctx := context.Background()

	err := service.CheckStaleAgents(ctx, 2*time.Minute)
	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestCheckStaleAgents_UpdateError(t *testing.T) {
	repo := newMockAgentRepo()
	repo.agents["stale"] = &entity.Agent{
		Paw:      "stale",
		Status:   entity.AgentOnline,
		LastSeen: time.Now().Add(-5 * time.Minute),
	}
	repo.updateErr = errors.New("update error")

	service := NewAgentService(repo)
	ctx := context.Background()

	err := service.CheckStaleAgents(ctx, 2*time.Minute)
	if err == nil {
		t.Fatal("Expected error")
	}
}
