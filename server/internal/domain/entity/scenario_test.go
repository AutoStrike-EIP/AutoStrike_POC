package entity

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
)

// sel is a helper to create TechniqueSelection from an ID
func sel(id string) TechniqueSelection {
	return TechniqueSelection{TechniqueID: id}
}

// selExec is a helper to create TechniqueSelection with executor name
func selExec(id, executor string) TechniqueSelection {
	return TechniqueSelection{TechniqueID: id, ExecutorName: executor}
}

func TestScenario_GetAllTechniques(t *testing.T) {
	tests := []struct {
		name     string
		scenario *Scenario
		want     []string
	}{
		{
			name: "single phase single technique",
			scenario: &Scenario{
				Phases: []Phase{
					{Name: "Phase1", Techniques: []TechniqueSelection{sel("T1059")}},
				},
			},
			want: []string{"T1059"},
		},
		{
			name: "multiple phases multiple techniques",
			scenario: &Scenario{
				Phases: []Phase{
					{Name: "Phase1", Techniques: []TechniqueSelection{sel("T1059"), sel("T1082")}},
					{Name: "Phase2", Techniques: []TechniqueSelection{sel("T1055"), sel("T1071")}},
				},
			},
			want: []string{"T1059", "T1082", "T1055", "T1071"},
		},
		{
			name: "duplicate techniques across phases",
			scenario: &Scenario{
				Phases: []Phase{
					{Name: "Phase1", Techniques: []TechniqueSelection{sel("T1059"), sel("T1082")}},
					{Name: "Phase2", Techniques: []TechniqueSelection{sel("T1059"), sel("T1071")}},
				},
			},
			want: []string{"T1059", "T1082", "T1071"},
		},
		{
			name: "empty phases",
			scenario: &Scenario{
				Phases: []Phase{},
			},
			want: nil,
		},
		{
			name: "phases with empty techniques",
			scenario: &Scenario{
				Phases: []Phase{
					{Name: "Phase1", Techniques: []TechniqueSelection{}},
					{Name: "Phase2", Techniques: []TechniqueSelection{sel("T1059")}},
				},
			},
			want: []string{"T1059"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.scenario.GetAllTechniques()
			if len(got) != len(tt.want) {
				t.Errorf("GetAllTechniques() returned %d techniques, want %d", len(got), len(tt.want))
				return
			}
			for i, tech := range got {
				if tech != tt.want[i] {
					t.Errorf("GetAllTechniques()[%d] = %s, want %s", i, tech, tt.want[i])
				}
			}
		})
	}
}

func TestScenario_TechniqueCount(t *testing.T) {
	tests := []struct {
		name     string
		scenario *Scenario
		want     int
	}{
		{
			name: "no techniques",
			scenario: &Scenario{
				Phases: []Phase{},
			},
			want: 0,
		},
		{
			name: "single technique",
			scenario: &Scenario{
				Phases: []Phase{
					{Name: "Phase1", Techniques: []TechniqueSelection{sel("T1059")}},
				},
			},
			want: 1,
		},
		{
			name: "multiple techniques no duplicates",
			scenario: &Scenario{
				Phases: []Phase{
					{Name: "Phase1", Techniques: []TechniqueSelection{sel("T1059"), sel("T1082")}},
					{Name: "Phase2", Techniques: []TechniqueSelection{sel("T1055")}},
				},
			},
			want: 3,
		},
		{
			name: "multiple techniques with duplicates",
			scenario: &Scenario{
				Phases: []Phase{
					{Name: "Phase1", Techniques: []TechniqueSelection{sel("T1059"), sel("T1082")}},
					{Name: "Phase2", Techniques: []TechniqueSelection{sel("T1059"), sel("T1055")}},
				},
			},
			want: 3, // T1059 counted only once
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.scenario.TechniqueCount(); got != tt.want {
				t.Errorf("TechniqueCount() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestPhase_Struct(t *testing.T) {
	phase := Phase{
		Name:        "Reconnaissance",
		Description: "Initial reconnaissance phase",
		Techniques:  []TechniqueSelection{sel("T1595"), sel("T1592")},
		Order:       1,
	}

	if phase.Name != "Reconnaissance" {
		t.Errorf("Name = %s, want Reconnaissance", phase.Name)
	}
	if phase.Order != 1 {
		t.Errorf("Order = %d, want 1", phase.Order)
	}
	if len(phase.Techniques) != 2 {
		t.Errorf("Techniques length = %d, want 2", len(phase.Techniques))
	}
}

func TestPhase_UnmarshalJSON_OldFormat(t *testing.T) {
	// Old format: techniques is a []string
	data := `{"name":"Phase1","description":"test","techniques":["T1059","T1082"],"order":1}`
	var phase Phase
	if err := json.Unmarshal([]byte(data), &phase); err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	if phase.Name != "Phase1" {
		t.Errorf("Name = %s, want Phase1", phase.Name)
	}
	if phase.Order != 1 {
		t.Errorf("Order = %d, want 1", phase.Order)
	}
	if len(phase.Techniques) != 2 {
		t.Fatalf("Techniques length = %d, want 2", len(phase.Techniques))
	}
	if phase.Techniques[0].TechniqueID != "T1059" {
		t.Errorf("Techniques[0].TechniqueID = %s, want T1059", phase.Techniques[0].TechniqueID)
	}
	if phase.Techniques[0].ExecutorName != "" {
		t.Errorf("Techniques[0].ExecutorName = %s, want empty", phase.Techniques[0].ExecutorName)
	}
	if phase.Techniques[1].TechniqueID != "T1082" {
		t.Errorf("Techniques[1].TechniqueID = %s, want T1082", phase.Techniques[1].TechniqueID)
	}
}

func TestPhase_UnmarshalJSON_NewFormat(t *testing.T) {
	// New format: techniques is a []TechniqueSelection
	data := `{"name":"Phase1","techniques":[{"technique_id":"T1059","executor_name":"whoami"},{"technique_id":"T1082"}],"order":2}`
	var phase Phase
	if err := json.Unmarshal([]byte(data), &phase); err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	if phase.Name != "Phase1" {
		t.Errorf("Name = %s, want Phase1", phase.Name)
	}
	if phase.Order != 2 {
		t.Errorf("Order = %d, want 2", phase.Order)
	}
	if len(phase.Techniques) != 2 {
		t.Fatalf("Techniques length = %d, want 2", len(phase.Techniques))
	}
	if phase.Techniques[0].TechniqueID != "T1059" {
		t.Errorf("Techniques[0].TechniqueID = %s, want T1059", phase.Techniques[0].TechniqueID)
	}
	if phase.Techniques[0].ExecutorName != "whoami" {
		t.Errorf("Techniques[0].ExecutorName = %s, want whoami", phase.Techniques[0].ExecutorName)
	}
	if phase.Techniques[1].TechniqueID != "T1082" {
		t.Errorf("Techniques[1].TechniqueID = %s, want T1082", phase.Techniques[1].TechniqueID)
	}
	if phase.Techniques[1].ExecutorName != "" {
		t.Errorf("Techniques[1].ExecutorName = %s, want empty", phase.Techniques[1].ExecutorName)
	}
}

func TestPhase_UnmarshalJSON_EmptyTechniques(t *testing.T) {
	data := `{"name":"Phase1","techniques":[],"order":0}`
	var phase Phase
	if err := json.Unmarshal([]byte(data), &phase); err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	if len(phase.Techniques) != 0 {
		t.Errorf("Techniques length = %d, want 0", len(phase.Techniques))
	}
}

func TestPhase_MarshalJSON(t *testing.T) {
	phase := Phase{
		Name:       "Phase1",
		Techniques: []TechniqueSelection{selExec("T1059", "whoami"), sel("T1082")},
		Order:      1,
	}

	data, err := json.Marshal(phase)
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	// Unmarshal back to verify round-trip
	var result Phase
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("UnmarshalJSON round-trip failed: %v", err)
	}

	if len(result.Techniques) != 2 {
		t.Fatalf("Round-trip Techniques length = %d, want 2", len(result.Techniques))
	}
	if result.Techniques[0].TechniqueID != "T1059" {
		t.Errorf("Round-trip Techniques[0].TechniqueID = %s, want T1059", result.Techniques[0].TechniqueID)
	}
	if result.Techniques[0].ExecutorName != "whoami" {
		t.Errorf("Round-trip Techniques[0].ExecutorName = %s, want whoami", result.Techniques[0].ExecutorName)
	}
}

func TestTechniqueSelection_Struct(t *testing.T) {
	s := TechniqueSelection{TechniqueID: "T1059", ExecutorName: "test-exec"}
	if s.TechniqueID != "T1059" {
		t.Errorf("TechniqueID = %s, want T1059", s.TechniqueID)
	}
	if s.ExecutorName != "test-exec" {
		t.Errorf("ExecutorName = %s, want test-exec", s.ExecutorName)
	}
}

func TestTechniqueSelection_UnmarshalYAML_ScalarNode(t *testing.T) {
	yamlData := `- T1059
- T1082`
	var selections []TechniqueSelection
	if err := yaml.Unmarshal([]byte(yamlData), &selections); err != nil {
		t.Fatalf("UnmarshalYAML failed: %v", err)
	}
	if len(selections) != 2 {
		t.Fatalf("Expected 2 selections, got %d", len(selections))
	}
	if selections[0].TechniqueID != "T1059" {
		t.Errorf("selections[0].TechniqueID = %s, want T1059", selections[0].TechniqueID)
	}
	if selections[0].ExecutorName != "" {
		t.Errorf("selections[0].ExecutorName = %s, want empty", selections[0].ExecutorName)
	}
	if selections[1].TechniqueID != "T1082" {
		t.Errorf("selections[1].TechniqueID = %s, want T1082", selections[1].TechniqueID)
	}
}

func TestTechniqueSelection_UnmarshalYAML_MappingNode(t *testing.T) {
	yamlData := `- technique_id: T1059
  executor_name: whoami
- technique_id: T1082`
	var selections []TechniqueSelection
	if err := yaml.Unmarshal([]byte(yamlData), &selections); err != nil {
		t.Fatalf("UnmarshalYAML failed: %v", err)
	}
	if len(selections) != 2 {
		t.Fatalf("Expected 2 selections, got %d", len(selections))
	}
	if selections[0].TechniqueID != "T1059" {
		t.Errorf("TechniqueID = %s, want T1059", selections[0].TechniqueID)
	}
	if selections[0].ExecutorName != "whoami" {
		t.Errorf("ExecutorName = %s, want whoami", selections[0].ExecutorName)
	}
	if selections[1].TechniqueID != "T1082" {
		t.Errorf("TechniqueID = %s, want T1082", selections[1].TechniqueID)
	}
	if selections[1].ExecutorName != "" {
		t.Errorf("ExecutorName = %s, want empty", selections[1].ExecutorName)
	}
}

func TestTechniqueSelection_UnmarshalYAML_InvalidMapping(t *testing.T) {
	yamlData := `- [invalid, yaml, sequence]`
	var selections []TechniqueSelection
	err := yaml.Unmarshal([]byte(yamlData), &selections)
	if err == nil {
		t.Error("Expected error for invalid YAML mapping, got nil")
	}
}

func TestPhase_UnmarshalJSON_InvalidJSON(t *testing.T) {
	data := `{invalid json}`
	var phase Phase
	err := json.Unmarshal([]byte(data), &phase)
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestPhase_UnmarshalJSON_InvalidTechniques(t *testing.T) {
	data := `{"name":"Phase1","techniques":123,"order":1}`
	var phase Phase
	err := json.Unmarshal([]byte(data), &phase)
	if err == nil {
		t.Error("Expected error for invalid techniques field, got nil")
	}
}

func TestPhase_UnmarshalJSON_FallbackWhenTechniqueIDEmpty(t *testing.T) {
	// This tests the case where JSON unmarshal into []TechniqueSelection succeeds
	// but TechniqueID is empty (it's actually a string array parsed as objects)
	data := `{"name":"Phase1","techniques":["T1059","T1082"],"order":1}`
	var phase Phase
	if err := json.Unmarshal([]byte(data), &phase); err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}
	// Should fall back to old format parsing
	if len(phase.Techniques) != 2 {
		t.Fatalf("Expected 2 techniques, got %d", len(phase.Techniques))
	}
	if phase.Techniques[0].TechniqueID != "T1059" {
		t.Errorf("TechniqueID = %s, want T1059", phase.Techniques[0].TechniqueID)
	}
}
