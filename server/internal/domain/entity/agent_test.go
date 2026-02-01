package entity

import (
	"testing"
	"time"
)

func TestAgent_IsOnline(t *testing.T) {
	tests := []struct {
		name     string
		lastSeen time.Time
		timeout  time.Duration
		want     bool
	}{
		{
			name:     "agent seen recently is online",
			lastSeen: time.Now().Add(-30 * time.Second),
			timeout:  2 * time.Minute,
			want:     true,
		},
		{
			name:     "agent seen long ago is offline",
			lastSeen: time.Now().Add(-5 * time.Minute),
			timeout:  2 * time.Minute,
			want:     false,
		},
		{
			name:     "agent at timeout boundary is offline",
			lastSeen: time.Now().Add(-2 * time.Minute),
			timeout:  2 * time.Minute,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent := &Agent{
				Paw:      "test-agent",
				LastSeen: tt.lastSeen,
			}
			if got := agent.IsOnline(tt.timeout); got != tt.want {
				t.Errorf("Agent.IsOnline() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAgent_IsCompatible(t *testing.T) {
	tests := []struct {
		name      string
		agent     *Agent
		technique *Technique
		want      bool
	}{
		{
			name: "compatible windows agent with powershell technique",
			agent: &Agent{
				Platform:  "windows",
				Executors: []string{"psh", "cmd"},
			},
			technique: &Technique{
				Platforms: []string{"windows"},
				Executors: []Executor{{Type: "psh", Command: "whoami"}},
			},
			want: true,
		},
		{
			name: "incompatible platform",
			agent: &Agent{
				Platform:  "linux",
				Executors: []string{"bash", "sh"},
			},
			technique: &Technique{
				Platforms: []string{"windows"},
				Executors: []Executor{{Type: "psh", Command: "whoami"}},
			},
			want: false,
		},
		{
			name: "incompatible executor",
			agent: &Agent{
				Platform:  "windows",
				Executors: []string{"cmd"},
			},
			technique: &Technique{
				Platforms: []string{"windows"},
				Executors: []Executor{{Type: "psh", Command: "whoami"}},
			},
			want: false,
		},
		{
			name: "compatible linux agent",
			agent: &Agent{
				Platform:  "linux",
				Executors: []string{"bash", "sh"},
			},
			technique: &Technique{
				Platforms: []string{"linux", "darwin"},
				Executors: []Executor{{Type: "bash", Command: "id"}},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.agent.IsCompatible(tt.technique); got != tt.want {
				t.Errorf("Agent.IsCompatible() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAgent_SupportsExecutor(t *testing.T) {
	agent := &Agent{
		Executors: []string{"psh", "cmd", "bash"},
	}

	tests := []struct {
		executor string
		want     bool
	}{
		{"psh", true},
		{"cmd", true},
		{"bash", true},
		{"zsh", false},
		{"python", false},
	}

	for _, tt := range tests {
		t.Run(tt.executor, func(t *testing.T) {
			if got := agent.SupportsExecutor(tt.executor); got != tt.want {
				t.Errorf("Agent.SupportsExecutor(%s) = %v, want %v", tt.executor, got, tt.want)
			}
		})
	}
}

func TestAgentStatus_Constants(t *testing.T) {
	// Verify status constants are defined correctly
	if AgentOnline != "online" {
		t.Errorf("AgentOnline = %s, want online", AgentOnline)
	}
	if AgentOffline != "offline" {
		t.Errorf("AgentOffline = %s, want offline", AgentOffline)
	}
	if AgentBusy != "busy" {
		t.Errorf("AgentBusy = %s, want busy", AgentBusy)
	}
	if AgentUntrusted != "untrusted" {
		t.Errorf("AgentUntrusted = %s, want untrusted", AgentUntrusted)
	}
}
