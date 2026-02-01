package entity

import (
	"testing"
)

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
					{Name: "Phase1", Techniques: []string{"T1059"}},
				},
			},
			want: []string{"T1059"},
		},
		{
			name: "multiple phases multiple techniques",
			scenario: &Scenario{
				Phases: []Phase{
					{Name: "Phase1", Techniques: []string{"T1059", "T1082"}},
					{Name: "Phase2", Techniques: []string{"T1055", "T1071"}},
				},
			},
			want: []string{"T1059", "T1082", "T1055", "T1071"},
		},
		{
			name: "duplicate techniques across phases",
			scenario: &Scenario{
				Phases: []Phase{
					{Name: "Phase1", Techniques: []string{"T1059", "T1082"}},
					{Name: "Phase2", Techniques: []string{"T1059", "T1071"}},
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
					{Name: "Phase1", Techniques: []string{}},
					{Name: "Phase2", Techniques: []string{"T1059"}},
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
					{Name: "Phase1", Techniques: []string{"T1059"}},
				},
			},
			want: 1,
		},
		{
			name: "multiple techniques no duplicates",
			scenario: &Scenario{
				Phases: []Phase{
					{Name: "Phase1", Techniques: []string{"T1059", "T1082"}},
					{Name: "Phase2", Techniques: []string{"T1055"}},
				},
			},
			want: 3,
		},
		{
			name: "multiple techniques with duplicates",
			scenario: &Scenario{
				Phases: []Phase{
					{Name: "Phase1", Techniques: []string{"T1059", "T1082"}},
					{Name: "Phase2", Techniques: []string{"T1059", "T1055"}},
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
		Techniques:  []string{"T1595", "T1592"},
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
