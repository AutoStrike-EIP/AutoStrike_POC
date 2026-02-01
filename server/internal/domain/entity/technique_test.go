package entity

import (
	"testing"
)

func TestTechnique_GetExecutorForPlatform(t *testing.T) {
	technique := &Technique{
		ID:        "T1059.001",
		Name:      "PowerShell",
		Platforms: []string{"windows"},
		Executors: []Executor{
			{Type: "psh", Command: "Get-Process", Timeout: 60},
			{Type: "cmd", Command: "tasklist", Timeout: 30},
		},
	}

	tests := []struct {
		name           string
		platform       string
		agentExecutors []string
		wantNil        bool
		wantType       string
	}{
		{
			name:           "matching platform and executor",
			platform:       "windows",
			agentExecutors: []string{"psh"},
			wantNil:        false,
			wantType:       "psh",
		},
		{
			name:           "matching platform, second executor",
			platform:       "windows",
			agentExecutors: []string{"cmd"},
			wantNil:        false,
			wantType:       "cmd",
		},
		{
			name:           "non-matching platform",
			platform:       "linux",
			agentExecutors: []string{"psh", "cmd"},
			wantNil:        true,
		},
		{
			name:           "no matching executor",
			platform:       "windows",
			agentExecutors: []string{"bash", "sh"},
			wantNil:        true,
		},
		{
			name:           "empty agent executors",
			platform:       "windows",
			agentExecutors: []string{},
			wantNil:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := technique.GetExecutorForPlatform(tt.platform, tt.agentExecutors)
			if tt.wantNil {
				if got != nil {
					t.Errorf("GetExecutorForPlatform() = %v, want nil", got)
				}
			} else {
				if got == nil {
					t.Errorf("GetExecutorForPlatform() = nil, want executor of type %s", tt.wantType)
				} else if got.Type != tt.wantType {
					t.Errorf("GetExecutorForPlatform().Type = %s, want %s", got.Type, tt.wantType)
				}
			}
		})
	}
}

func TestTechnique_MultiPlatform(t *testing.T) {
	technique := &Technique{
		ID:        "T1082",
		Name:      "System Information Discovery",
		Platforms: []string{"windows", "linux", "darwin"},
		Executors: []Executor{
			{Type: "psh", Command: "systeminfo", Timeout: 60},
			{Type: "bash", Command: "uname -a", Timeout: 30},
		},
	}

	// Test Windows
	exec := technique.GetExecutorForPlatform("windows", []string{"psh"})
	if exec == nil || exec.Type != "psh" {
		t.Error("Expected psh executor for windows")
	}

	// Test Linux
	exec = technique.GetExecutorForPlatform("linux", []string{"bash"})
	if exec == nil || exec.Type != "bash" {
		t.Error("Expected bash executor for linux")
	}

	// Test Darwin
	exec = technique.GetExecutorForPlatform("darwin", []string{"bash"})
	if exec == nil || exec.Type != "bash" {
		t.Error("Expected bash executor for darwin")
	}
}

func TestTacticType_Constants(t *testing.T) {
	// Verify all MITRE tactics are defined
	tactics := []TacticType{
		TacticReconnaissance,
		TacticResourceDevelopment,
		TacticInitialAccess,
		TacticExecution,
		TacticPersistence,
		TacticPrivilegeEscalation,
		TacticDefenseEvasion,
		TacticCredentialAccess,
		TacticDiscovery,
		TacticLateralMovement,
		TacticCollection,
		TacticCommandAndControl,
		TacticExfiltration,
		TacticImpact,
	}

	if len(tactics) != 14 {
		t.Errorf("Expected 14 MITRE tactics, got %d", len(tactics))
	}

	// Verify naming convention
	expectedValues := []string{
		"reconnaissance", "resource-development", "initial-access", "execution",
		"persistence", "privilege-escalation", "defense-evasion", "credential-access",
		"discovery", "lateral-movement", "collection", "command-and-control",
		"exfiltration", "impact",
	}

	for i, tactic := range tactics {
		if string(tactic) != expectedValues[i] {
			t.Errorf("Tactic %d = %s, want %s", i, tactic, expectedValues[i])
		}
	}
}
