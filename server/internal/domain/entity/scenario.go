package entity

import (
	"encoding/json"
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

// Scenario represents an attack scenario with multiple phases
type Scenario struct {
	ID          string    `json:"id" yaml:"id"`
	Name        string    `json:"name" yaml:"name"`
	Description string    `json:"description" yaml:"description"`
	Phases      []Phase   `json:"phases" yaml:"phases"`
	Tags        []string  `json:"tags,omitempty" yaml:"tags,omitempty"`
	Author      string    `json:"author,omitempty" yaml:"author,omitempty"`
	CreatedAt   time.Time `json:"created_at" yaml:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" yaml:"updated_at"`
}

// TechniqueSelection represents a technique with an optional executor preference
type TechniqueSelection struct {
	TechniqueID  string `json:"technique_id" yaml:"technique_id"`
	ExecutorName string `json:"executor_name,omitempty" yaml:"executor_name,omitempty"`
}

// UnmarshalYAML supports both old format (plain string) and new format (map with technique_id)
func (ts *TechniqueSelection) UnmarshalYAML(value *yaml.Node) error {
	// Try plain string first (old format)
	if value.Kind == yaml.ScalarNode {
		ts.TechniqueID = value.Value
		ts.ExecutorName = ""
		return nil
	}

	// New format: mapping node with technique_id and executor_name
	type techniqueSelectionAlias TechniqueSelection
	var alias techniqueSelectionAlias
	if err := value.Decode(&alias); err != nil {
		return fmt.Errorf("cannot unmarshal technique selection: %w", err)
	}
	ts.TechniqueID = alias.TechniqueID
	ts.ExecutorName = alias.ExecutorName
	return nil
}

// Phase represents a phase in a scenario
type Phase struct {
	Name        string               `json:"name" yaml:"name"`
	Description string               `json:"description,omitempty" yaml:"description,omitempty"`
	Techniques  []TechniqueSelection `json:"techniques" yaml:"techniques"`
	Order       int                  `json:"order" yaml:"order"`
}

// phaseJSON is used for custom unmarshalling to support both old and new formats
type phaseJSON struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Techniques  json.RawMessage `json:"techniques"`
	Order       int             `json:"order"`
}

// UnmarshalJSON supports both old format ([]string) and new format ([]TechniqueSelection)
func (p *Phase) UnmarshalJSON(data []byte) error {
	var raw phaseJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	p.Name = raw.Name
	p.Description = raw.Description
	p.Order = raw.Order

	// Try new format first: []TechniqueSelection
	var selections []TechniqueSelection
	if json.Unmarshal(raw.Techniques, &selections) == nil {
		// Verify it's actually the new format (objects with technique_id field)
		// by checking if the first element has a non-empty TechniqueID
		if len(selections) == 0 || selections[0].TechniqueID != "" {
			p.Techniques = selections
			return nil
		}
	}

	// Fallback: old format []string
	var techIDs []string
	if err := json.Unmarshal(raw.Techniques, &techIDs); err != nil {
		return err
	}

	p.Techniques = make([]TechniqueSelection, len(techIDs))
	for i, id := range techIDs {
		p.Techniques[i] = TechniqueSelection{TechniqueID: id}
	}
	return nil
}

// GetAllTechniques returns all unique technique IDs from all phases
func (s *Scenario) GetAllTechniques() []string {
	seen := make(map[string]bool)
	var techniques []string

	for _, phase := range s.Phases {
		for _, sel := range phase.Techniques {
			if !seen[sel.TechniqueID] {
				seen[sel.TechniqueID] = true
				techniques = append(techniques, sel.TechniqueID)
			}
		}
	}
	return techniques
}

// TechniqueCount returns the total number of techniques in the scenario
func (s *Scenario) TechniqueCount() int {
	return len(s.GetAllTechniques())
}
