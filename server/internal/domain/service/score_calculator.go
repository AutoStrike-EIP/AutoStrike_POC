package service

import (
	"autostrike/internal/domain/entity"
)

// ScoreCalculator calculates security scores from execution results
type ScoreCalculator struct{}

// NewScoreCalculator creates a new score calculator
func NewScoreCalculator() *ScoreCalculator {
	return &ScoreCalculator{}
}

// CalculateScore calculates the security score from execution results
// Score formula: (blocked*100 + detected*50) / (total*100) * 100
func (s *ScoreCalculator) CalculateScore(results []*entity.ExecutionResult) *entity.SecurityScore {
	score := &entity.SecurityScore{
		ByTactic: make(map[string]float64),
	}

	if len(results) == 0 {
		score.Overall = 100.0 // No attacks = perfect score (or undefined?)
		return score
	}

	var blocked, detected, successful, total int

	for _, result := range results {
		if result.Status == entity.StatusSkipped || result.Status == entity.StatusPending {
			continue
		}

		total++

		switch result.Status {
		case entity.StatusBlocked:
			blocked++
		case entity.StatusDetected:
			detected++
		case entity.StatusSuccess:
			successful++
		}
	}

	score.Blocked = blocked
	score.Detected = detected
	score.Successful = successful
	score.Total = total

	if total > 0 {
		// Blocked = 100 points, Detected = 50 points, Success = 0 points
		maxPoints := float64(total * 100)
		earnedPoints := float64(blocked*100 + detected*50)
		score.Overall = (earnedPoints / maxPoints) * 100
	}

	return score
}

// CalculateScoreByTactic calculates scores grouped by MITRE tactic
func (s *ScoreCalculator) CalculateScoreByTactic(
	results []*entity.ExecutionResult,
	techniques map[string]*entity.Technique,
) map[entity.TacticType]*entity.SecurityScore {
	tacticResults := make(map[entity.TacticType][]*entity.ExecutionResult)

	for _, result := range results {
		if tech, ok := techniques[result.TechniqueID]; ok {
			tacticResults[tech.Tactic] = append(tacticResults[tech.Tactic], result)
		}
	}

	scores := make(map[entity.TacticType]*entity.SecurityScore)
	for tactic, tacticRes := range tacticResults {
		scores[tactic] = s.CalculateScore(tacticRes)
	}

	return scores
}

// CalculateTrend compares two score sets and returns the difference
func (s *ScoreCalculator) CalculateTrend(current, previous *entity.SecurityScore) float64 {
	if previous == nil || previous.Overall == 0 {
		return 0
	}
	return current.Overall - previous.Overall
}

// GetCoveragePercentage returns the percentage of MITRE techniques tested
func (s *ScoreCalculator) GetCoveragePercentage(testedTechniques, totalTechniques int) float64 {
	if totalTechniques == 0 {
		return 0
	}
	return (float64(testedTechniques) / float64(totalTechniques)) * 100
}
