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
	ID          string      `json:"id" yaml:"id"`                   // "T1059.001"
	Name        string      `json:"name" yaml:"name"`               // "PowerShell"
	Tactic      TacticType  `json:"tactic" yaml:"tactic"`           // "execution"
	Description string      `json:"description" yaml:"description"` // Detailed description
	Platforms   []string    `json:"platforms" yaml:"platforms"`     // ["windows"]
	Executors   []Executor  `json:"executors" yaml:"executors"`
	Detection   []Detection `json:"detection,omitempty" yaml:"detection,omitempty"`
	References  []string    `json:"references,omitempty" yaml:"references,omitempty"`
	IsSafe      bool        `json:"is_safe" yaml:"is_safe"` // Safe for production
}

// Executor defines how to execute the technique
type Executor struct {
	Type    string `json:"type" yaml:"type"`       // "psh", "cmd", "bash"
	Command string `json:"command" yaml:"command"` // The command to execute
	Cleanup string `json:"cleanup,omitempty" yaml:"cleanup,omitempty"`
	Timeout int    `json:"timeout" yaml:"timeout"` // Seconds
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
		for _, agentExec := range agentExecutors {
			if t.Executors[i].Type == agentExec {
				return &t.Executors[i]
			}
		}
	}
	return nil
}
