package service

import (
	"testing"

	"autostrike/internal/domain/entity"
)

func TestScoreCalculator_CalculateScore(t *testing.T) {
	calc := NewScoreCalculator()

	tests := []struct {
		name        string
		results     []*entity.ExecutionResult
		wantOverall float64
		wantBlocked int
		wantDetect  int
		wantSuccess int
		wantTotal   int
	}{
		{
			name:        "empty results returns 100",
			results:     []*entity.ExecutionResult{},
			wantOverall: 100.0,
			wantTotal:   0,
		},
		{
			name: "all blocked returns 100",
			results: []*entity.ExecutionResult{
				{Status: entity.StatusBlocked},
				{Status: entity.StatusBlocked},
				{Status: entity.StatusBlocked},
			},
			wantOverall: 100.0,
			wantBlocked: 3,
			wantTotal:   3,
		},
		{
			name: "all successful returns 0",
			results: []*entity.ExecutionResult{
				{Status: entity.StatusSuccess},
				{Status: entity.StatusSuccess},
			},
			wantOverall: 0.0,
			wantSuccess: 2,
			wantTotal:   2,
		},
		{
			name: "all detected returns 50",
			results: []*entity.ExecutionResult{
				{Status: entity.StatusDetected},
				{Status: entity.StatusDetected},
			},
			wantOverall: 50.0,
			wantDetect:  2,
			wantTotal:   2,
		},
		{
			name: "mixed results",
			results: []*entity.ExecutionResult{
				{Status: entity.StatusBlocked},  // 100 points
				{Status: entity.StatusDetected}, // 50 points
				{Status: entity.StatusSuccess},  // 0 points
				{Status: entity.StatusSuccess},  // 0 points
			},
			wantOverall: 37.5, // (100 + 50) / 400 * 100 = 37.5
			wantBlocked: 1,
			wantDetect:  1,
			wantSuccess: 2,
			wantTotal:   4,
		},
		{
			name: "skipped and pending are ignored",
			results: []*entity.ExecutionResult{
				{Status: entity.StatusBlocked},
				{Status: entity.StatusSkipped},
				{Status: entity.StatusPending},
			},
			wantOverall: 100.0,
			wantBlocked: 1,
			wantTotal:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := calc.CalculateScore(tt.results)

			if score.Overall != tt.wantOverall {
				t.Errorf("Overall = %f, want %f", score.Overall, tt.wantOverall)
			}
			if score.Blocked != tt.wantBlocked {
				t.Errorf("Blocked = %d, want %d", score.Blocked, tt.wantBlocked)
			}
			if score.Detected != tt.wantDetect {
				t.Errorf("Detected = %d, want %d", score.Detected, tt.wantDetect)
			}
			if score.Successful != tt.wantSuccess {
				t.Errorf("Successful = %d, want %d", score.Successful, tt.wantSuccess)
			}
			if score.Total != tt.wantTotal {
				t.Errorf("Total = %d, want %d", score.Total, tt.wantTotal)
			}
		})
	}
}

func TestScoreCalculator_CalculateTrend(t *testing.T) {
	calc := NewScoreCalculator()

	tests := []struct {
		name     string
		current  *entity.SecurityScore
		previous *entity.SecurityScore
		want     float64
	}{
		{
			name:     "nil previous returns 0",
			current:  &entity.SecurityScore{Overall: 80.0},
			previous: nil,
			want:     0,
		},
		{
			name:     "zero previous returns 0",
			current:  &entity.SecurityScore{Overall: 80.0},
			previous: &entity.SecurityScore{Overall: 0},
			want:     0,
		},
		{
			name:     "positive trend",
			current:  &entity.SecurityScore{Overall: 80.0},
			previous: &entity.SecurityScore{Overall: 60.0},
			want:     20.0,
		},
		{
			name:     "negative trend",
			current:  &entity.SecurityScore{Overall: 50.0},
			previous: &entity.SecurityScore{Overall: 70.0},
			want:     -20.0,
		},
		{
			name:     "no change",
			current:  &entity.SecurityScore{Overall: 75.0},
			previous: &entity.SecurityScore{Overall: 75.0},
			want:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calc.CalculateTrend(tt.current, tt.previous)
			if got != tt.want {
				t.Errorf("CalculateTrend() = %f, want %f", got, tt.want)
			}
		})
	}
}

func TestScoreCalculator_GetCoveragePercentage(t *testing.T) {
	calc := NewScoreCalculator()

	tests := []struct {
		tested  int
		total   int
		wantMin float64
		wantMax float64
	}{
		{0, 0, 0, 0},
		{0, 100, 0, 0},
		{50, 100, 50.0, 50.0},
		{100, 100, 100.0, 100.0},
		{25, 50, 50.0, 50.0},
		{1, 3, 33.33, 33.34}, // ~33.333...%
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := calc.GetCoveragePercentage(tt.tested, tt.total)
			if got < tt.wantMin || got > tt.wantMax {
				t.Errorf("GetCoveragePercentage(%d, %d) = %f, want between %f and %f", tt.tested, tt.total, got, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestScoreCalculator_CalculateScoreByTactic(t *testing.T) {
	calc := NewScoreCalculator()

	techniques := map[string]*entity.Technique{
		"T1082": {ID: "T1082", Tactic: entity.TacticDiscovery},
		"T1059": {ID: "T1059", Tactic: entity.TacticExecution},
	}

	results := []*entity.ExecutionResult{
		{TechniqueID: "T1082", Status: entity.StatusBlocked},
		{TechniqueID: "T1082", Status: entity.StatusDetected},
		{TechniqueID: "T1059", Status: entity.StatusSuccess},
	}

	scores := calc.CalculateScoreByTactic(results, techniques)

	if len(scores) != 2 {
		t.Errorf("Expected 2 tactic scores, got %d", len(scores))
	}

	discoveryScore := scores[entity.TacticDiscovery]
	if discoveryScore == nil {
		t.Fatal("Discovery score is nil")
	}
	// Discovery: 1 blocked (100) + 1 detected (50) = 150/200 = 75%
	if discoveryScore.Overall != 75.0 {
		t.Errorf("Discovery score = %f, want 75.0", discoveryScore.Overall)
	}

	executionScore := scores[entity.TacticExecution]
	if executionScore == nil {
		t.Fatal("Execution score is nil")
	}
	// Execution: 1 success (0) = 0/100 = 0%
	if executionScore.Overall != 0.0 {
		t.Errorf("Execution score = %f, want 0.0", executionScore.Overall)
	}
}
