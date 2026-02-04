package application

import (
	"context"
	"time"

	"autostrike/internal/domain/entity"
	"autostrike/internal/domain/repository"
)

// AnalyticsService provides analytics and reporting functionality
type AnalyticsService struct {
	resultRepo repository.ResultRepository
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(resultRepo repository.ResultRepository) *AnalyticsService {
	return &AnalyticsService{
		resultRepo: resultRepo,
	}
}

// PeriodStats represents statistics for a time period
type PeriodStats struct {
	Period          string          `json:"period"`
	StartDate       time.Time       `json:"start_date"`
	EndDate         time.Time       `json:"end_date"`
	ExecutionCount  int             `json:"execution_count"`
	AverageScore    float64         `json:"average_score"`
	TotalBlocked    int             `json:"total_blocked"`
	TotalDetected   int             `json:"total_detected"`
	TotalSuccessful int             `json:"total_successful"`
	TotalTechniques int             `json:"total_techniques"`
	ScoreByTactic   map[string]float64 `json:"score_by_tactic,omitempty"`
}

// ScoreComparison represents comparison between two periods
type ScoreComparison struct {
	Current       PeriodStats `json:"current"`
	Previous      PeriodStats `json:"previous"`
	ScoreChange   float64     `json:"score_change"`
	ScoreTrend    string      `json:"score_trend"` // "improving", "declining", "stable"
	BlockedChange int         `json:"blocked_change"`
	DetectedChange int        `json:"detected_change"`
}

// TrendDataPoint represents a single data point in a trend
type TrendDataPoint struct {
	Date           string  `json:"date"`
	AverageScore   float64 `json:"average_score"`
	ExecutionCount int     `json:"execution_count"`
	Blocked        int     `json:"blocked"`
	Detected       int     `json:"detected"`
	Successful     int     `json:"successful"`
}

// ScoreTrend represents score trend over time
type ScoreTrend struct {
	Period     string           `json:"period"` // "7d", "30d", "90d"
	DataPoints []TrendDataPoint `json:"data_points"`
	Summary    TrendSummary     `json:"summary"`
}

// TrendSummary provides summary statistics for a trend
type TrendSummary struct {
	StartScore       float64 `json:"start_score"`
	EndScore         float64 `json:"end_score"`
	AverageScore     float64 `json:"average_score"`
	MaxScore         float64 `json:"max_score"`
	MinScore         float64 `json:"min_score"`
	TotalExecutions  int     `json:"total_executions"`
	OverallTrend     string  `json:"overall_trend"` // "improving", "declining", "stable"
	PercentageChange float64 `json:"percentage_change"`
}

// ExecutionSummary provides overall execution analytics
type ExecutionSummary struct {
	TotalExecutions     int                `json:"total_executions"`
	CompletedExecutions int                `json:"completed_executions"`
	AverageScore        float64            `json:"average_score"`
	BestScore           float64            `json:"best_score"`
	WorstScore          float64            `json:"worst_score"`
	ScoresByScenario    map[string]float64 `json:"scores_by_scenario"`
	ExecutionsByStatus  map[string]int     `json:"executions_by_status"`
}

// GetPeriodStats calculates statistics for a given time period
func (s *AnalyticsService) GetPeriodStats(ctx context.Context, start, end time.Time, periodLabel string) (*PeriodStats, error) {
	executions, err := s.resultRepo.FindCompletedExecutionsByDateRange(ctx, start, end)
	if err != nil {
		return nil, err
	}

	stats := &PeriodStats{
		Period:        periodLabel,
		StartDate:     start,
		EndDate:       end,
		ScoreByTactic: make(map[string]float64),
	}

	if len(executions) == 0 {
		return stats, nil
	}

	stats.ExecutionCount = len(executions)

	var totalScore float64
	tacticScores := make(map[string][]float64)

	for _, exec := range executions {
		if exec.Score != nil {
			totalScore += exec.Score.Overall
			stats.TotalBlocked += exec.Score.Blocked
			stats.TotalDetected += exec.Score.Detected
			stats.TotalSuccessful += exec.Score.Successful
			stats.TotalTechniques += exec.Score.Total

			// Aggregate tactic scores
			for tactic, score := range exec.Score.ByTactic {
				tacticScores[tactic] = append(tacticScores[tactic], score)
			}
		}
	}

	stats.AverageScore = totalScore / float64(len(executions))

	// Calculate average per tactic
	for tactic, scores := range tacticScores {
		var sum float64
		for _, score := range scores {
			sum += score
		}
		stats.ScoreByTactic[tactic] = sum / float64(len(scores))
	}

	return stats, nil
}

// CompareScores compares scores between two consecutive periods
func (s *AnalyticsService) CompareScores(ctx context.Context, periodDays int) (*ScoreComparison, error) {
	now := time.Now()

	// Current period
	currentEnd := now
	currentStart := now.AddDate(0, 0, -periodDays)

	// Previous period
	previousEnd := currentStart
	previousStart := previousEnd.AddDate(0, 0, -periodDays)

	currentStats, err := s.GetPeriodStats(ctx, currentStart, currentEnd, "current")
	if err != nil {
		return nil, err
	}

	previousStats, err := s.GetPeriodStats(ctx, previousStart, previousEnd, "previous")
	if err != nil {
		return nil, err
	}

	comparison := &ScoreComparison{
		Current:        *currentStats,
		Previous:       *previousStats,
		ScoreChange:   currentStats.AverageScore - previousStats.AverageScore,
		BlockedChange:  currentStats.TotalBlocked - previousStats.TotalBlocked,
		DetectedChange: currentStats.TotalDetected - previousStats.TotalDetected,
	}

	// Determine trend
	if comparison.ScoreChange > 5 {
		comparison.ScoreTrend = "improving"
	} else if comparison.ScoreChange < -5 {
		comparison.ScoreTrend = "declining"
	} else {
		comparison.ScoreTrend = "stable"
	}

	return comparison, nil
}

// GetScoreTrend gets score trend over time
func (s *AnalyticsService) GetScoreTrend(ctx context.Context, days int) (*ScoreTrend, error) {
	now := time.Now()
	start := now.AddDate(0, 0, -days)

	executions, err := s.resultRepo.FindCompletedExecutionsByDateRange(ctx, start, now)
	if err != nil {
		return nil, err
	}

	// Group executions by date
	dateGroups := make(map[string][]*entity.Execution)
	for _, exec := range executions {
		dateKey := exec.StartedAt.Format("2006-01-02")
		dateGroups[dateKey] = append(dateGroups[dateKey], exec)
	}

	// Generate data points for each day
	dataPoints := make([]TrendDataPoint, 0)
	var allScores []float64
	var maxScore, minScore float64
	firstSet := true

	for d := start; !d.After(now); d = d.AddDate(0, 0, 1) {
		dateKey := d.Format("2006-01-02")
		dayExecs := dateGroups[dateKey]

		point := TrendDataPoint{
			Date:           dateKey,
			ExecutionCount: len(dayExecs),
		}

		if len(dayExecs) > 0 {
			var dayTotal float64
			for _, exec := range dayExecs {
				if exec.Score != nil {
					dayTotal += exec.Score.Overall
					point.Blocked += exec.Score.Blocked
					point.Detected += exec.Score.Detected
					point.Successful += exec.Score.Successful

					if firstSet {
						maxScore = exec.Score.Overall
						minScore = exec.Score.Overall
						firstSet = false
					} else {
						if exec.Score.Overall > maxScore {
							maxScore = exec.Score.Overall
						}
						if exec.Score.Overall < minScore {
							minScore = exec.Score.Overall
						}
					}
				}
			}
			point.AverageScore = dayTotal / float64(len(dayExecs))
			allScores = append(allScores, point.AverageScore)
		}

		dataPoints = append(dataPoints, point)
	}

	// Calculate summary
	var totalScore float64
	for _, score := range allScores {
		totalScore += score
	}

	summary := TrendSummary{
		TotalExecutions: len(executions),
		MaxScore:        maxScore,
		MinScore:        minScore,
	}

	if len(allScores) > 0 {
		summary.AverageScore = totalScore / float64(len(allScores))
		summary.StartScore = allScores[0]
		summary.EndScore = allScores[len(allScores)-1]

		if summary.StartScore > 0 {
			summary.PercentageChange = ((summary.EndScore - summary.StartScore) / summary.StartScore) * 100
		}

		if summary.PercentageChange > 5 {
			summary.OverallTrend = "improving"
		} else if summary.PercentageChange < -5 {
			summary.OverallTrend = "declining"
		} else {
			summary.OverallTrend = "stable"
		}
	}

	periodLabel := "7d"
	if days == 30 {
		periodLabel = "30d"
	} else if days == 90 {
		periodLabel = "90d"
	}

	return &ScoreTrend{
		Period:     periodLabel,
		DataPoints: dataPoints,
		Summary:    summary,
	}, nil
}

// GetExecutionSummary provides overall execution analytics
func (s *AnalyticsService) GetExecutionSummary(ctx context.Context, days int) (*ExecutionSummary, error) {
	now := time.Now()
	start := now.AddDate(0, 0, -days)

	executions, err := s.resultRepo.FindExecutionsByDateRange(ctx, start, now)
	if err != nil {
		return nil, err
	}

	summary := &ExecutionSummary{
		TotalExecutions:    len(executions),
		ScoresByScenario:   make(map[string]float64),
		ExecutionsByStatus: make(map[string]int),
	}

	scenarioScores := make(map[string][]float64)
	var completedScores []float64
	firstCompleted := true

	for _, exec := range executions {
		summary.ExecutionsByStatus[string(exec.Status)]++

		if exec.Status == entity.ExecutionCompleted && exec.Score != nil {
			summary.CompletedExecutions++
			completedScores = append(completedScores, exec.Score.Overall)
			scenarioScores[exec.ScenarioID] = append(scenarioScores[exec.ScenarioID], exec.Score.Overall)

			if firstCompleted {
				summary.BestScore = exec.Score.Overall
				summary.WorstScore = exec.Score.Overall
				firstCompleted = false
			} else {
				if exec.Score.Overall > summary.BestScore {
					summary.BestScore = exec.Score.Overall
				}
				if exec.Score.Overall < summary.WorstScore {
					summary.WorstScore = exec.Score.Overall
				}
			}
		}
	}

	// Calculate averages
	if len(completedScores) > 0 {
		var total float64
		for _, score := range completedScores {
			total += score
		}
		summary.AverageScore = total / float64(len(completedScores))
	}

	// Calculate per-scenario averages
	for scenarioID, scores := range scenarioScores {
		var total float64
		for _, score := range scores {
			total += score
		}
		summary.ScoresByScenario[scenarioID] = total / float64(len(scores))
	}

	return summary, nil
}
