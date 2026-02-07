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

func TestTechnique_GetExecutorForPlatform_WithPlatformField(t *testing.T) {
	technique := &Technique{
		ID:        "T1082",
		Name:      "System Info",
		Platforms: []string{"windows", "linux"},
		Executors: []Executor{
			{Name: "win-systeminfo", Type: "cmd", Platform: "windows", Command: "systeminfo", Timeout: 60},
			{Name: "linux-uname", Type: "bash", Platform: "linux", Command: "uname -a", Timeout: 30},
		},
	}

	// Windows agent should get the windows executor only
	exec := technique.GetExecutorForPlatform("windows", []string{"cmd", "bash"})
	if exec == nil || exec.Name != "win-systeminfo" {
		t.Errorf("Expected win-systeminfo, got %v", exec)
	}

	// Linux agent should get the linux executor only
	exec = technique.GetExecutorForPlatform("linux", []string{"cmd", "bash"})
	if exec == nil || exec.Name != "linux-uname" {
		t.Errorf("Expected linux-uname, got %v", exec)
	}
}

func TestTechnique_GetExecutorByName(t *testing.T) {
	technique := &Technique{
		ID:        "T1082",
		Name:      "System Info",
		Platforms: []string{"windows", "linux"},
		Executors: []Executor{
			{Name: "basic-info", Type: "cmd", Platform: "windows", Command: "systeminfo", Timeout: 60},
			{Name: "detailed-info", Type: "cmd", Platform: "windows", Command: "systeminfo /fo csv", Timeout: 60},
			{Name: "linux-info", Type: "bash", Platform: "linux", Command: "uname -a", Timeout: 30},
		},
	}

	tests := []struct {
		name           string
		execName       string
		platform       string
		agentExecutors []string
		wantNil        bool
		wantCommand    string
	}{
		{
			name:           "match by name",
			execName:       "detailed-info",
			platform:       "windows",
			agentExecutors: []string{"cmd"},
			wantNil:        false,
			wantCommand:    "systeminfo /fo csv",
		},
		{
			name:           "name not found",
			execName:       "nonexistent",
			platform:       "windows",
			agentExecutors: []string{"cmd"},
			wantNil:        true,
		},
		{
			name:           "name found but wrong platform",
			execName:       "linux-info",
			platform:       "windows",
			agentExecutors: []string{"bash"},
			wantNil:        true,
		},
		{
			name:           "name found but agent lacks executor",
			execName:       "basic-info",
			platform:       "windows",
			agentExecutors: []string{"bash"},
			wantNil:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := technique.GetExecutorByName(tt.execName, tt.platform, tt.agentExecutors)
			if tt.wantNil {
				if got != nil {
					t.Errorf("GetExecutorByName() = %v, want nil", got)
				}
			} else {
				if got == nil {
					t.Errorf("GetExecutorByName() = nil, want executor with command %s", tt.wantCommand)
				} else if got.Command != tt.wantCommand {
					t.Errorf("GetExecutorByName().Command = %s, want %s", got.Command, tt.wantCommand)
				}
			}
		})
	}
}

func TestTechnique_GetExecutorsForPlatform(t *testing.T) {
	technique := &Technique{
		ID:        "T1082",
		Name:      "System Info",
		Platforms: []string{"windows", "linux"},
		Executors: []Executor{
			{Name: "basic-info", Type: "cmd", Platform: "windows", Command: "systeminfo", Timeout: 60},
			{Name: "detailed-info", Type: "cmd", Platform: "windows", Command: "systeminfo /fo csv", Timeout: 60},
			{Name: "linux-info", Type: "bash", Platform: "linux", Command: "uname -a", Timeout: 30},
			{Name: "generic", Type: "bash", Command: "echo info", Timeout: 10},
		},
	}

	// Windows with cmd: should get basic-info and detailed-info only
	execs := technique.GetExecutorsForPlatform("windows", []string{"cmd"})
	if len(execs) != 2 {
		t.Fatalf("Expected 2 executors for windows/cmd, got %d", len(execs))
	}

	// Linux with bash: should get linux-info and generic
	execs = technique.GetExecutorsForPlatform("linux", []string{"bash"})
	if len(execs) != 2 {
		t.Fatalf("Expected 2 executors for linux/bash, got %d", len(execs))
	}

	// Unsupported platform: should get nil
	execs = technique.GetExecutorsForPlatform("macos", []string{"bash"})
	if execs != nil {
		t.Errorf("Expected nil for unsupported platform, got %d executors", len(execs))
	}
}

func TestTechnique_GetTactics(t *testing.T) {
	tests := []struct {
		name      string
		technique Technique
		want      []TacticType
	}{
		{
			name: "tactics field populated",
			technique: Technique{
				Tactic:  TacticInitialAccess,
				Tactics: []TacticType{TacticInitialAccess, TacticPersistence, TacticDefenseEvasion},
			},
			want: []TacticType{TacticInitialAccess, TacticPersistence, TacticDefenseEvasion},
		},
		{
			name: "tactics field empty, fallback to tactic",
			technique: Technique{
				Tactic: TacticDiscovery,
			},
			want: []TacticType{TacticDiscovery},
		},
		{
			name: "both empty",
			technique: Technique{
				Tactic:  "",
				Tactics: nil,
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.technique.GetTactics()
			if len(got) != len(tt.want) {
				t.Errorf("GetTactics() length = %d, want %d", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("GetTactics()[%d] = %s, want %s", i, got[i], tt.want[i])
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

func TestExecutor_NewFields(t *testing.T) {
	exec := Executor{
		Name:              "test-exec",
		Type:              "bash",
		Platform:          "linux",
		Command:           "id",
		Cleanup:           "rm /tmp/test",
		Timeout:           30,
		ElevationRequired: true,
	}

	if exec.Name != "test-exec" {
		t.Errorf("Name = %s, want test-exec", exec.Name)
	}
	if exec.Platform != "linux" {
		t.Errorf("Platform = %s, want linux", exec.Platform)
	}
	if !exec.ElevationRequired {
		t.Error("ElevationRequired = false, want true")
	}
}
