package entity

import (
	"time"
)

type AgentStatus string

const (
	AgentOnline    AgentStatus = "online"
	AgentOffline   AgentStatus = "offline"
	AgentBusy      AgentStatus = "busy"
	AgentUntrusted AgentStatus = "untrusted"
)

// Agent represents a deployed AutoStrike agent
type Agent struct {
	Paw       string            `json:"paw"`
	Hostname  string            `json:"hostname"`
	Platform  string            `json:"platform"` // "windows", "linux", "darwin"
	Username  string            `json:"username"`
	Executors []string          `json:"executors"` // ["psh", "cmd", "bash"]
	Status    AgentStatus       `json:"status"`
	LastSeen  time.Time         `json:"last_seen"`
	IPAddress string            `json:"ip_address"`
	OSVersion string            `json:"os_version"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
}

// IsOnline returns true if the agent is considered online
func (a *Agent) IsOnline(timeout time.Duration) bool {
	return time.Since(a.LastSeen) < timeout
}

// IsCompatible checks if the agent can execute the given technique
func (a *Agent) IsCompatible(technique *Technique) bool {
	for _, platform := range technique.Platforms {
		if platform == a.Platform {
			for _, executor := range technique.Executors {
				for _, agentExec := range a.Executors {
					if executor.Type == agentExec {
						return true
					}
				}
			}
		}
	}
	return false
}

// SupportsExecutor checks if the agent supports a specific executor type
func (a *Agent) SupportsExecutor(executorType string) bool {
	for _, exec := range a.Executors {
		if exec == executorType {
			return true
		}
	}
	return false
}
