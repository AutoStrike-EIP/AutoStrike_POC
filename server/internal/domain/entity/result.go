package entity

import (
	"time"
)

// ResultStatus represents the outcome of a technique execution
type ResultStatus string

const (
	StatusPending  ResultStatus = "pending"  // Waiting to execute
	StatusRunning  ResultStatus = "running"  // Currently executing
	StatusSuccess  ResultStatus = "success"  // Executed, not detected
	StatusBlocked  ResultStatus = "blocked"  // Blocked by defense
	StatusDetected ResultStatus = "detected" // Executed but alerted
	StatusFailed   ResultStatus = "failed"   // Technical error
	StatusSkipped  ResultStatus = "skipped"  // Not executed
	StatusTimeout  ResultStatus = "timeout"  // Execution timed out
)

// ExecutionResult represents the result of a single technique execution
type ExecutionResult struct {
	ID           string        `json:"id"`
	ExecutionID  string        `json:"execution_id"`
	TechniqueID  string        `json:"technique_id"`
	AgentPaw     string        `json:"agent_paw"`
	ExecutorName string        `json:"executor_name,omitempty"` // Which executor was used
	Command      string        `json:"command,omitempty"`       // The command that was executed
	Status       ResultStatus  `json:"status"`
	Output       string        `json:"output"`                  // Command output (always present)
	Stderr       string        `json:"stderr,omitempty"`        // Base64 encoded
	ExitCode     int           `json:"exit_code"`
	Detected     bool          `json:"detected"`                // Was the technique detected?
	DetectedBy   string        `json:"detected_by,omitempty"`   // "Windows Defender", "CrowdStrike"
	StartedAt    time.Time     `json:"started_at"`
	CompletedAt  *time.Time    `json:"completed_at,omitempty"`
	Duration     time.Duration `json:"duration_ms"`
}

// Execution represents a scenario execution session
type Execution struct {
	ID          string            `json:"id"`
	ScenarioID  string            `json:"scenario_id"`
	AgentPaws   []string          `json:"agent_paws"`
	Status      ExecutionStatus   `json:"status"`
	Progress    ExecutionProgress `json:"progress"`
	Results     []ExecutionResult `json:"results,omitempty"`
	Score       *SecurityScore    `json:"score,omitempty"`
	StartedAt   time.Time         `json:"started_at"`
	CompletedAt *time.Time        `json:"completed_at,omitempty"`
	StartedBy   string            `json:"started_by"`
	SafeMode    bool              `json:"safe_mode"`
}

// ExecutionStatus represents the status of an execution
type ExecutionStatus string

const (
	ExecutionPending   ExecutionStatus = "pending"
	ExecutionRunning   ExecutionStatus = "running"
	ExecutionCompleted ExecutionStatus = "completed"
	ExecutionFailed    ExecutionStatus = "failed"
	ExecutionCancelled ExecutionStatus = "cancelled"
)

// ExecutionProgress tracks execution progress
type ExecutionProgress struct {
	Total     int `json:"total"`
	Completed int `json:"completed"`
	Failed    int `json:"failed"`
	Skipped   int `json:"skipped"`
}

// SecurityScore represents the calculated security score
type SecurityScore struct {
	Overall    float64            `json:"overall"`    // 0-100
	ByTactic   map[string]float64 `json:"by_tactic"`  // Score per tactic
	Blocked    int                `json:"blocked"`    // Count
	Detected   int                `json:"detected"`   // Count
	Successful int                `json:"successful"` // Undetected executions
	Total      int                `json:"total"`      // Total techniques tested
}

// IsComplete returns true if the result has completed (success, failed, blocked, etc.)
func (r *ExecutionResult) IsComplete() bool {
	return r.Status != StatusPending && r.Status != StatusRunning
}

// IsSuccessful returns true if the technique executed without detection
func (r *ExecutionResult) IsSuccessful() bool {
	return r.Status == StatusSuccess
}
