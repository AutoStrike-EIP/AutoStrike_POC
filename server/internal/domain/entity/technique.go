package entity

// TacticType represents a MITRE ATT&CK tactic
type TacticType string

const (
	TacticReconnaissance      TacticType = "reconnaissance"
	TacticResourceDevelopment TacticType = "resource-development"
	TacticInitialAccess       TacticType = "initial-access"
	TacticExecution           TacticType = "execution"
	TacticPersistence         TacticType = "persistence"
	TacticPrivilegeEscalation TacticType = "privilege-escalation"
	TacticDefenseEvasion      TacticType = "defense-evasion"
	TacticCredentialAccess    TacticType = "credential-access"
	TacticDiscovery           TacticType = "discovery"
	TacticLateralMovement     TacticType = "lateral-movement"
	TacticCollection          TacticType = "collection"
	TacticCommandAndControl   TacticType = "command-and-control"
	TacticExfiltration        TacticType = "exfiltration"
	TacticImpact              TacticType = "impact"
)

// Technique represents a MITRE ATT&CK technique
type Technique struct {
	ID          string       `json:"id" yaml:"id"`                                    // "T1059.001"
	Name        string       `json:"name" yaml:"name"`                                // "PowerShell"
	Tactic      TacticType   `json:"tactic" yaml:"tactic"`                            // Primary tactic (retro-compat)
	Tactics     []TacticType `json:"tactics,omitempty" yaml:"tactics,omitempty"`       // All tactics (multi-tactic)
	Description string       `json:"description" yaml:"description"`                  // Detailed description
	Platforms   []string     `json:"platforms" yaml:"platforms"`                       // ["windows"]
	Executors   []Executor   `json:"executors" yaml:"executors"`
	Detection   []Detection  `json:"detection,omitempty" yaml:"detection,omitempty"`
	References  []string     `json:"references,omitempty" yaml:"references,omitempty"`
	IsSafe      bool         `json:"is_safe" yaml:"is_safe"` // Safe for production
}

// Executor defines how to execute the technique
type Executor struct {
	Name              string `json:"name,omitempty" yaml:"name,omitempty"`
	Type              string `json:"type" yaml:"type"`       // "psh", "cmd", "bash"
	Platform          string `json:"platform,omitempty" yaml:"platform,omitempty"`
	Command           string `json:"command" yaml:"command"` // The command to execute
	Cleanup           string `json:"cleanup,omitempty" yaml:"cleanup,omitempty"`
	Timeout           int    `json:"timeout" yaml:"timeout"` // Seconds
	ElevationRequired bool   `json:"elevation_required,omitempty" yaml:"elevation_required,omitempty"`
}

// Detection describes expected detection indicators
type Detection struct {
	Source    string `json:"source" yaml:"source"`       // "Process Creation", "File Creation"
	Indicator string `json:"indicator" yaml:"indicator"` // Pattern description
}

// GetExecutorForPlatform returns the first compatible executor for the given platform
func (t *Technique) GetExecutorForPlatform(platform string, agentExecutors []string) *Executor {
	// Check if platform is supported
	platformSupported := false
	for _, p := range t.Platforms {
		if p == platform {
			platformSupported = true
			break
		}
	}

	if !platformSupported {
		return nil
	}

	// Find compatible executor
	for i := range t.Executors {
		// If executor has a platform set, it must match
		if t.Executors[i].Platform != "" && t.Executors[i].Platform != platform {
			continue
		}
		for _, agentExec := range agentExecutors {
			if t.Executors[i].Type == agentExec {
				return &t.Executors[i]
			}
		}
	}
	return nil
}

// GetExecutorByName returns an executor by name, filtered by platform and agent capabilities
func (t *Technique) GetExecutorByName(name string, platform string, agentExecutors []string) *Executor {
	for i := range t.Executors {
		if t.Executors[i].Name != name {
			continue
		}
		// If executor has a platform set, it must match
		if t.Executors[i].Platform != "" && t.Executors[i].Platform != platform {
			continue
		}
		for _, agentExec := range agentExecutors {
			if t.Executors[i].Type == agentExec {
				return &t.Executors[i]
			}
		}
	}
	return nil
}

// GetExecutorsForPlatform returns all compatible executors for the given platform
func (t *Technique) GetExecutorsForPlatform(platform string, agentExecutors []string) []Executor {
	// Check if platform is supported
	platformSupported := false
	for _, p := range t.Platforms {
		if p == platform {
			platformSupported = true
			break
		}
	}

	if !platformSupported {
		return nil
	}

	var result []Executor
	for _, exec := range t.Executors {
		if exec.Platform != "" && exec.Platform != platform {
			continue
		}
		for _, agentExec := range agentExecutors {
			if exec.Type == agentExec {
				result = append(result, exec)
				break
			}
		}
	}
	return result
}

// GetTactics returns all tactics for the technique. Falls back to the primary tactic if Tactics is empty.
func (t *Technique) GetTactics() []TacticType {
	if len(t.Tactics) > 0 {
		return t.Tactics
	}
	if t.Tactic != "" {
		return []TacticType{t.Tactic}
	}
	return nil
}
