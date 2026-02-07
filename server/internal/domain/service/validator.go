package service

import (
	"autostrike/internal/domain/entity"
)

// TechniqueValidator validates technique compatibility
type TechniqueValidator struct{}

// NewTechniqueValidator creates a new validator instance
func NewTechniqueValidator() *TechniqueValidator {
	return &TechniqueValidator{}
}

// ValidationResult contains the validation result
type ValidationResult struct {
	IsValid  bool
	Errors   []string
	Warnings []string
}

// ValidateAgentCompatibility checks if an agent can execute a technique
func (v *TechniqueValidator) ValidateAgentCompatibility(
	agent *entity.Agent,
	technique *entity.Technique,
) *ValidationResult {
	result := &ValidationResult{
		IsValid:  true,
		Errors:   make([]string, 0),
		Warnings: make([]string, 0),
	}

	// Check platform compatibility
	platformMatch := false
	for _, platform := range technique.Platforms {
		if platform == agent.Platform {
			platformMatch = true
			break
		}
	}

	if !platformMatch {
		result.IsValid = false
		result.Errors = append(result.Errors,
			"agent platform not supported by technique")
		return result
	}

	// Check executor compatibility
	executorMatch := false
	for _, executor := range technique.Executors {
		for _, agentExec := range agent.Executors {
			if executor.Type == agentExec {
				executorMatch = true
				break
			}
		}
		if executorMatch {
			break
		}
	}

	if !executorMatch {
		result.IsValid = false
		result.Errors = append(result.Errors,
			"no compatible executor found")
		return result
	}

	// Check agent status
	if agent.Status != entity.AgentOnline {
		result.Warnings = append(result.Warnings,
			"agent is not currently online")
	}

	// Check if technique is safe
	if !technique.IsSafe {
		result.Warnings = append(result.Warnings,
			"technique may cause system modifications")
	}

	return result
}

// ValidateScenario validates a complete scenario
func (v *TechniqueValidator) ValidateScenario(
	scenario *entity.Scenario,
	techniques []*entity.Technique,
) *ValidationResult {
	result := &ValidationResult{
		IsValid:  true,
		Errors:   make([]string, 0),
		Warnings: make([]string, 0),
	}

	if len(scenario.Phases) == 0 {
		result.IsValid = false
		result.Errors = append(result.Errors, "scenario has no phases")
		return result
	}

	techniqueMap := make(map[string]*entity.Technique)
	for _, t := range techniques {
		techniqueMap[t.ID] = t
	}

	for _, phase := range scenario.Phases {
		if len(phase.Techniques) == 0 {
			result.Warnings = append(result.Warnings,
				"phase '"+phase.Name+"' has no techniques")
		}

		for _, sel := range phase.Techniques {
			if _, exists := techniqueMap[sel.TechniqueID]; !exists {
				result.Errors = append(result.Errors,
					"technique '"+sel.TechniqueID+"' not found")
				result.IsValid = false
			}
		}
	}

	return result
}
