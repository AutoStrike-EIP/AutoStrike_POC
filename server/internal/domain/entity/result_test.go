package entity

import (
	"testing"
)

func TestExecutionResult_IsComplete(t *testing.T) {
	tests := []struct {
		status ResultStatus
		want   bool
	}{
		{StatusPending, false},
		{StatusRunning, false},
		{StatusSuccess, true},
		{StatusBlocked, true},
		{StatusDetected, true},
		{StatusFailed, true},
		{StatusSkipped, true},
		{StatusTimeout, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			result := &ExecutionResult{Status: tt.status}
			if got := result.IsComplete(); got != tt.want {
				t.Errorf("IsComplete() for status %s = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}

func TestExecutionResult_IsSuccessful(t *testing.T) {
	tests := []struct {
		status ResultStatus
		want   bool
	}{
		{StatusSuccess, true},
		{StatusBlocked, false},
		{StatusDetected, false},
		{StatusFailed, false},
		{StatusPending, false},
		{StatusRunning, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			result := &ExecutionResult{Status: tt.status}
			if got := result.IsSuccessful(); got != tt.want {
				t.Errorf("IsSuccessful() for status %s = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}

func TestResultStatus_Constants(t *testing.T) {
	// Verify all status constants
	statuses := map[ResultStatus]string{
		StatusPending:  "pending",
		StatusRunning:  "running",
		StatusSuccess:  "success",
		StatusBlocked:  "blocked",
		StatusDetected: "detected",
		StatusFailed:   "failed",
		StatusSkipped:  "skipped",
		StatusTimeout:  "timeout",
	}

	for status, expected := range statuses {
		if string(status) != expected {
			t.Errorf("Status %v = %s, want %s", status, status, expected)
		}
	}
}

func TestExecutionStatus_Constants(t *testing.T) {
	statuses := map[ExecutionStatus]string{
		ExecutionPending:   "pending",
		ExecutionRunning:   "running",
		ExecutionCompleted: "completed",
		ExecutionFailed:    "failed",
		ExecutionCancelled: "cancelled",
	}

	for status, expected := range statuses {
		if string(status) != expected {
			t.Errorf("ExecutionStatus %v = %s, want %s", status, status, expected)
		}
	}
}

func TestSecurityScore_Initialization(t *testing.T) {
	score := &SecurityScore{
		Overall:    75.5,
		ByTactic:   map[string]float64{"discovery": 80.0, "execution": 70.0},
		Blocked:    5,
		Detected:   3,
		Successful: 2,
		Total:      10,
	}

	if score.Overall != 75.5 {
		t.Errorf("Overall = %f, want 75.5", score.Overall)
	}
	if score.Total != 10 {
		t.Errorf("Total = %d, want 10", score.Total)
	}
	if len(score.ByTactic) != 2 {
		t.Errorf("ByTactic length = %d, want 2", len(score.ByTactic))
	}
}
