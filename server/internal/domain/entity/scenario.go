package entity

import (
	"time"
)

// Scenario represents an attack scenario with multiple phases
type Scenario struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Phases      []Phase   `json:"phases"`
	Tags        []string  `json:"tags,omitempty"`
	Author      string    `json:"author,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Phase represents a phase in a scenario
type Phase struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Techniques  []string `json:"techniques"` // Technique IDs
	Order       int      `json:"order"`
}

// GetAllTechniques returns all unique technique IDs from all phases
func (s *Scenario) GetAllTechniques() []string {
	seen := make(map[string]bool)
	var techniques []string

	for _, phase := range s.Phases {
		for _, techID := range phase.Techniques {
			if !seen[techID] {
				seen[techID] = true
				techniques = append(techniques, techID)
			}
		}
	}
	return techniques
}

// TechniqueCount returns the total number of techniques in the scenario
func (s *Scenario) TechniqueCount() int {
	return len(s.GetAllTechniques())
}
