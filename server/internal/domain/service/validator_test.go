package service

import (
	"testing"

	"autostrike/internal/domain/entity"
)

func TestTechniqueValidator_ValidateAgentCompatibility(t *testing.T) {
	validator := NewTechniqueValidator()

	tests := []struct {
		name         string
		agent        *entity.Agent
		technique    *entity.Technique
		wantValid    bool
		wantErrors   int
		wantWarnings int
	}{
		{
			name: "compatible agent and technique",
			agent: &entity.Agent{
				Platform:  "windows",
				Executors: []string{"psh", "cmd"},
				Status:    entity.AgentOnline,
			},
			technique: &entity.Technique{
				Platforms: []string{"windows"},
				Executors: []entity.Executor{{Type: "psh"}},
				IsSafe:    true,
			},
			wantValid:    true,
			wantErrors:   0,
			wantWarnings: 0,
		},
		{
			name: "incompatible platform",
			agent: &entity.Agent{
				Platform:  "linux",
				Executors: []string{"bash"},
				Status:    entity.AgentOnline,
			},
			technique: &entity.Technique{
				Platforms: []string{"windows"},
				Executors: []entity.Executor{{Type: "psh"}},
				IsSafe:    true,
			},
			wantValid:  false,
			wantErrors: 1,
		},
		{
			name: "incompatible executor",
			agent: &entity.Agent{
				Platform:  "windows",
				Executors: []string{"cmd"},
				Status:    entity.AgentOnline,
			},
			technique: &entity.Technique{
				Platforms: []string{"windows"},
				Executors: []entity.Executor{{Type: "psh"}},
				IsSafe:    true,
			},
			wantValid:  false,
			wantErrors: 1,
		},
		{
			name: "agent offline warning",
			agent: &entity.Agent{
				Platform:  "windows",
				Executors: []string{"psh"},
				Status:    entity.AgentOffline,
			},
			technique: &entity.Technique{
				Platforms: []string{"windows"},
				Executors: []entity.Executor{{Type: "psh"}},
				IsSafe:    true,
			},
			wantValid:    true,
			wantWarnings: 1,
		},
		{
			name: "unsafe technique warning",
			agent: &entity.Agent{
				Platform:  "windows",
				Executors: []string{"psh"},
				Status:    entity.AgentOnline,
			},
			technique: &entity.Technique{
				Platforms: []string{"windows"},
				Executors: []entity.Executor{{Type: "psh"}},
				IsSafe:    false,
			},
			wantValid:    true,
			wantWarnings: 1,
		},
		{
			name: "offline agent with unsafe technique",
			agent: &entity.Agent{
				Platform:  "windows",
				Executors: []string{"psh"},
				Status:    entity.AgentBusy,
			},
			technique: &entity.Technique{
				Platforms: []string{"windows"},
				Executors: []entity.Executor{{Type: "psh"}},
				IsSafe:    false,
			},
			wantValid:    true,
			wantWarnings: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateAgentCompatibility(tt.agent, tt.technique)

			if result.IsValid != tt.wantValid {
				t.Errorf("IsValid = %v, want %v", result.IsValid, tt.wantValid)
			}
			if len(result.Errors) != tt.wantErrors {
				t.Errorf("Errors count = %d, want %d: %v", len(result.Errors), tt.wantErrors, result.Errors)
			}
			if len(result.Warnings) != tt.wantWarnings {
				t.Errorf("Warnings count = %d, want %d: %v", len(result.Warnings), tt.wantWarnings, result.Warnings)
			}
		})
	}
}

func TestTechniqueValidator_ValidateScenario(t *testing.T) {
	validator := NewTechniqueValidator()

	techniques := []*entity.Technique{
		{ID: "T1082", Name: "System Info"},
		{ID: "T1059", Name: "Command Execution"},
	}

	tests := []struct {
		name         string
		scenario     *entity.Scenario
		wantValid    bool
		wantErrors   int
		wantWarnings int
	}{
		{
			name: "valid scenario",
			scenario: &entity.Scenario{
				Name: "Test Scenario",
				Phases: []entity.Phase{
					{Name: "Recon", Techniques: []entity.TechniqueSelection{{TechniqueID: "T1082"}}},
					{Name: "Execution", Techniques: []entity.TechniqueSelection{{TechniqueID: "T1059"}}},
				},
			},
			wantValid:  true,
			wantErrors: 0,
		},
		{
			name: "empty phases",
			scenario: &entity.Scenario{
				Name:   "Empty Scenario",
				Phases: []entity.Phase{},
			},
			wantValid:  false,
			wantErrors: 1,
		},
		{
			name: "missing technique",
			scenario: &entity.Scenario{
				Name: "Bad Scenario",
				Phases: []entity.Phase{
					{Name: "Phase1", Techniques: []entity.TechniqueSelection{{TechniqueID: "T9999"}}},
				},
			},
			wantValid:  false,
			wantErrors: 1,
		},
		{
			name: "empty phase warning",
			scenario: &entity.Scenario{
				Name: "Warning Scenario",
				Phases: []entity.Phase{
					{Name: "Empty Phase", Techniques: []entity.TechniqueSelection{}},
					{Name: "Valid Phase", Techniques: []entity.TechniqueSelection{{TechniqueID: "T1082"}}},
				},
			},
			wantValid:    true,
			wantWarnings: 1,
		},
		{
			name: "multiple missing techniques",
			scenario: &entity.Scenario{
				Name: "Multi Error",
				Phases: []entity.Phase{
					{Name: "Phase1", Techniques: []entity.TechniqueSelection{{TechniqueID: "T9999"}, {TechniqueID: "T8888"}}},
				},
			},
			wantValid:  false,
			wantErrors: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateScenario(tt.scenario, techniques)

			if result.IsValid != tt.wantValid {
				t.Errorf("IsValid = %v, want %v", result.IsValid, tt.wantValid)
			}
			if len(result.Errors) != tt.wantErrors {
				t.Errorf("Errors count = %d, want %d: %v", len(result.Errors), tt.wantErrors, result.Errors)
			}
			if len(result.Warnings) != tt.wantWarnings {
				t.Errorf("Warnings count = %d, want %d: %v", len(result.Warnings), tt.wantWarnings, result.Warnings)
			}
		})
	}
}

func TestNewTechniqueValidator(t *testing.T) {
	validator := NewTechniqueValidator()
	if validator == nil {
		t.Error("NewTechniqueValidator() returned nil")
	}
}

func TestValidationResult_Initialization(t *testing.T) {
	result := &ValidationResult{
		IsValid:  true,
		Errors:   []string{},
		Warnings: []string{},
	}

	if !result.IsValid {
		t.Error("Expected IsValid to be true")
	}
	if len(result.Errors) != 0 {
		t.Error("Expected empty Errors slice")
	}
	if len(result.Warnings) != 0 {
		t.Error("Expected empty Warnings slice")
	}
}
