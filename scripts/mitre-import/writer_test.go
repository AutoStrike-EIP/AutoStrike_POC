package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestWriteYAMLFiles_BasicWrite(t *testing.T) {
	tmpDir := t.TempDir()

	techniques := []*MergedTechnique{
		{
			ID:        "T1082",
			Name:      "System Information Discovery",
			Tactic:    "discovery",
			Tactics:   []string{"discovery"},
			Platforms: []string{"windows", "linux"},
			Executors: []MergedExecutor{
				{Name: "systeminfo", Type: "cmd", Platform: "windows", Command: "systeminfo", Timeout: 60},
				{Name: "uname", Type: "bash", Platform: "linux", Command: "uname -a", Timeout: 60},
			},
			IsSafe: true,
		},
	}

	result, err := WriteYAMLFiles(techniques, tmpDir)
	if err != nil {
		t.Fatalf("WriteYAMLFiles failed: %v", err)
	}

	if result.FilesWritten != 1 {
		t.Errorf("FilesWritten = %d, want 1", result.FilesWritten)
	}
	if result.TotalTechniques != 1 {
		t.Errorf("TotalTechniques = %d, want 1", result.TotalTechniques)
	}

	// Check file exists
	path := filepath.Join(tmpDir, "discovery.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	content := string(data)

	// Should have header comment
	if !strings.HasPrefix(content, "# AutoStrike MITRE ATT&CK Techniques") {
		t.Error("Missing header comment")
	}

	// Should contain technique ID
	if !strings.Contains(content, "T1082") {
		t.Error("Missing technique ID T1082")
	}
}

func TestWriteYAMLFiles_MultipleFiles(t *testing.T) {
	tmpDir := t.TempDir()

	techniques := []*MergedTechnique{
		{
			ID:        "T1082",
			Name:      "System Info",
			Tactic:    "discovery",
			Platforms: []string{"windows"},
			Executors: []MergedExecutor{
				{Name: "test", Type: "cmd", Platform: "windows", Command: "systeminfo", Timeout: 60},
			},
			IsSafe: true,
		},
		{
			ID:        "T1059",
			Name:      "Command Execution",
			Tactic:    "execution",
			Platforms: []string{"windows"},
			Executors: []MergedExecutor{
				{Name: "test", Type: "cmd", Platform: "windows", Command: "cmd /c whoami", Timeout: 60},
			},
			IsSafe: false,
		},
	}

	result, err := WriteYAMLFiles(techniques, tmpDir)
	if err != nil {
		t.Fatalf("WriteYAMLFiles failed: %v", err)
	}

	if result.FilesWritten != 2 {
		t.Errorf("FilesWritten = %d, want 2", result.FilesWritten)
	}

	// Check both files exist
	if _, err := os.Stat(filepath.Join(tmpDir, "discovery.yaml")); err != nil {
		t.Error("discovery.yaml not created")
	}
	if _, err := os.Stat(filepath.Join(tmpDir, "execution.yaml")); err != nil {
		t.Error("execution.yaml not created")
	}
}

func TestWriteYAMLFiles_OutputFormat(t *testing.T) {
	tmpDir := t.TempDir()

	techniques := []*MergedTechnique{
		{
			ID:          "T1082",
			Name:        "System Info",
			Description: "Test description",
			Tactic:      "discovery",
			Tactics:     []string{"discovery"},
			Platforms:   []string{"windows"},
			Executors: []MergedExecutor{
				{Name: "systeminfo", Type: "cmd", Platform: "windows", Command: "systeminfo", Timeout: 120},
			},
			References: []string{"https://attack.mitre.org/techniques/T1082"},
			IsSafe:     true,
		},
	}

	_, err := WriteYAMLFiles(techniques, tmpDir)
	if err != nil {
		t.Fatalf("WriteYAMLFiles failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmpDir, "discovery.yaml"))
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	// Skip header comments and parse YAML
	content := string(data)
	yamlStart := strings.Index(content, "- id:")
	if yamlStart == -1 {
		// Try parsing from after the last comment line
		lines := strings.Split(content, "\n")
		var yamlLines []string
		for _, line := range lines {
			if !strings.HasPrefix(line, "#") && line != "" {
				yamlLines = append(yamlLines, line)
			}
		}
		content = strings.Join(yamlLines, "\n")
	} else {
		content = content[yamlStart:]
	}

	var parsed []YAMLTechnique
	if err := yaml.Unmarshal([]byte(content), &parsed); err != nil {
		t.Fatalf("Failed to parse output YAML: %v", err)
	}

	if len(parsed) != 1 {
		t.Fatalf("Expected 1 technique in YAML, got %d", len(parsed))
	}

	tech := parsed[0]
	if tech.ID != "T1082" {
		t.Errorf("ID = %s, want T1082", tech.ID)
	}
	if tech.Name != "System Info" {
		t.Errorf("Name = %s, want System Info", tech.Name)
	}
	if tech.Tactic != "discovery" {
		t.Errorf("Tactic = %s, want discovery", tech.Tactic)
	}
	if !tech.IsSafe {
		t.Error("IsSafe should be true")
	}
	if len(tech.Executors) != 1 {
		t.Errorf("Expected 1 executor, got %d", len(tech.Executors))
	}
	if tech.Executors[0].Type != "cmd" {
		t.Errorf("Executor type = %s, want cmd", tech.Executors[0].Type)
	}
}

func TestWriteYAMLFiles_MultiTacticIncludesTactics(t *testing.T) {
	tmpDir := t.TempDir()

	techniques := []*MergedTechnique{
		{
			ID:        "T1078",
			Name:      "Valid Accounts",
			Tactic:    "defense-evasion",
			Tactics:   []string{"defense-evasion", "persistence", "initial-access"},
			Platforms: []string{"windows"},
			Executors: []MergedExecutor{
				{Name: "test", Type: "cmd", Platform: "windows", Command: "net user", Timeout: 60},
			},
			IsSafe: false,
		},
	}

	_, err := WriteYAMLFiles(techniques, tmpDir)
	if err != nil {
		t.Fatalf("WriteYAMLFiles failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmpDir, "defense-evasion.yaml"))
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	content := string(data)
	// Multi-tactic technique should include tactics array
	if !strings.Contains(content, "tactics:") {
		t.Error("Multi-tactic technique should include tactics array")
	}
	if !strings.Contains(content, "persistence") {
		t.Error("tactics should include persistence")
	}
	if !strings.Contains(content, "initial-access") {
		t.Error("tactics should include initial-access")
	}
}

func TestWriteYAMLFiles_SingleTacticOmitsTactics(t *testing.T) {
	tmpDir := t.TempDir()

	techniques := []*MergedTechnique{
		{
			ID:        "T1082",
			Name:      "System Info",
			Tactic:    "discovery",
			Tactics:   []string{"discovery"},
			Platforms: []string{"windows"},
			Executors: []MergedExecutor{
				{Name: "test", Type: "cmd", Platform: "windows", Command: "systeminfo", Timeout: 60},
			},
			IsSafe: true,
		},
	}

	_, err := WriteYAMLFiles(techniques, tmpDir)
	if err != nil {
		t.Fatalf("WriteYAMLFiles failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmpDir, "discovery.yaml"))
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	content := string(data)
	// Single-tactic should NOT include tactics array (omitempty)
	if strings.Contains(content, "tactics:") {
		t.Error("Single-tactic technique should not include tactics array")
	}
}

func TestWriteYAMLFiles_SortedByID(t *testing.T) {
	tmpDir := t.TempDir()

	techniques := []*MergedTechnique{
		{
			ID: "T1083", Name: "File Discovery", Tactic: "discovery",
			Platforms: []string{"linux"},
			Executors: []MergedExecutor{{Name: "test", Type: "bash", Platform: "linux", Command: "ls", Timeout: 60}},
			IsSafe:    true,
		},
		{
			ID: "T1057", Name: "Process Discovery", Tactic: "discovery",
			Platforms: []string{"linux"},
			Executors: []MergedExecutor{{Name: "test", Type: "bash", Platform: "linux", Command: "ps", Timeout: 60}},
			IsSafe:    true,
		},
		{
			ID: "T1082", Name: "System Info", Tactic: "discovery",
			Platforms: []string{"linux"},
			Executors: []MergedExecutor{{Name: "test", Type: "bash", Platform: "linux", Command: "uname", Timeout: 60}},
			IsSafe:    true,
		},
	}

	_, err := WriteYAMLFiles(techniques, tmpDir)
	if err != nil {
		t.Fatalf("WriteYAMLFiles failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmpDir, "discovery.yaml"))
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	content := string(data)
	posT1057 := strings.Index(content, "T1057")
	posT1082 := strings.Index(content, "T1082")
	posT1083 := strings.Index(content, "T1083")

	if posT1057 > posT1082 || posT1082 > posT1083 {
		t.Error("Techniques should be sorted by ID (T1057 < T1082 < T1083)")
	}
}

func TestWriteYAMLFiles_EmptyTechniques(t *testing.T) {
	tmpDir := t.TempDir()

	result, err := WriteYAMLFiles(nil, tmpDir)
	if err != nil {
		t.Fatalf("WriteYAMLFiles failed: %v", err)
	}

	if result.FilesWritten != 0 {
		t.Errorf("FilesWritten = %d, want 0", result.FilesWritten)
	}
}

func TestWriteYAMLFiles_CreatesOutputDir(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "nested", "output", "dir")

	techniques := []*MergedTechnique{
		{
			ID: "T1082", Name: "System Info", Tactic: "discovery",
			Platforms: []string{"windows"},
			Executors: []MergedExecutor{{Name: "test", Type: "cmd", Platform: "windows", Command: "systeminfo", Timeout: 60}},
			IsSafe:    true,
		},
	}

	_, err := WriteYAMLFiles(techniques, outputDir)
	if err != nil {
		t.Fatalf("WriteYAMLFiles failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(outputDir, "discovery.yaml")); err != nil {
		t.Error("Should create nested output directory")
	}
}

func TestTacticFilename(t *testing.T) {
	tests := []struct {
		tactic   string
		expected string
	}{
		{"discovery", "discovery.yaml"},
		{"initial-access", "initial-access.yaml"},
		{"command-and-control", "command-and-control.yaml"},
		{"privilege-escalation", "privilege-escalation.yaml"},
		{"defense-evasion", "defense-evasion.yaml"},
		{"unknown", "unknown.yaml"},
	}

	for _, tt := range tests {
		t.Run(tt.tactic, func(t *testing.T) {
			got := tacticFilename(tt.tactic)
			if got != tt.expected {
				t.Errorf("tacticFilename(%q) = %q, want %q", tt.tactic, got, tt.expected)
			}
		})
	}
}

func TestToYAMLTechnique(t *testing.T) {
	tech := &MergedTechnique{
		ID:          "T1082",
		Name:        "System Info",
		Description: "Discovers system information",
		Tactic:      "discovery",
		Tactics:     []string{"discovery"},
		Platforms:   []string{"windows", "linux"},
		Executors: []MergedExecutor{
			{Name: "systeminfo", Type: "cmd", Platform: "windows", Command: "systeminfo", Cleanup: "del output.txt", Timeout: 120, ElevationRequired: false},
		},
		References: []string{"https://attack.mitre.org/techniques/T1082"},
		IsSafe:     true,
	}

	yt := toYAMLTechnique(tech)

	if yt.ID != "T1082" {
		t.Errorf("ID = %s, want T1082", yt.ID)
	}
	if yt.Tactic != "discovery" {
		t.Errorf("Tactic = %s, want discovery", yt.Tactic)
	}
	if yt.Tactics != nil {
		t.Error("Single tactic should have nil Tactics (omitempty)")
	}
	if len(yt.Executors) != 1 {
		t.Errorf("Expected 1 executor, got %d", len(yt.Executors))
	}
	if yt.Executors[0].Name != "systeminfo" {
		t.Errorf("Executor name = %s, want systeminfo", yt.Executors[0].Name)
	}
	if yt.Executors[0].Cleanup != "del output.txt" {
		t.Errorf("Executor cleanup = %s, want del output.txt", yt.Executors[0].Cleanup)
	}
}

func TestWriteYAMLFiles_WithElevationRequired(t *testing.T) {
	tmpDir := t.TempDir()

	techniques := []*MergedTechnique{
		{
			ID: "T1548", Name: "Abuse Elevation", Tactic: "privilege-escalation",
			Platforms: []string{"linux"},
			Executors: []MergedExecutor{
				{Name: "sudo test", Type: "bash", Platform: "linux", Command: "sudo whoami", Timeout: 60, ElevationRequired: true},
			},
			IsSafe: false,
		},
	}

	_, err := WriteYAMLFiles(techniques, tmpDir)
	if err != nil {
		t.Fatalf("WriteYAMLFiles failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmpDir, "privilege-escalation.yaml"))
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "elevation_required: true") {
		t.Error("Should contain elevation_required: true")
	}
}

func TestWriteYAMLFiles_FileBreakdown(t *testing.T) {
	tmpDir := t.TempDir()

	techniques := []*MergedTechnique{
		{ID: "T1082", Name: "A", Tactic: "discovery", Platforms: []string{"linux"},
			Executors: []MergedExecutor{{Name: "a", Type: "bash", Platform: "linux", Command: "a", Timeout: 60}}, IsSafe: true},
		{ID: "T1083", Name: "B", Tactic: "discovery", Platforms: []string{"linux"},
			Executors: []MergedExecutor{{Name: "b", Type: "bash", Platform: "linux", Command: "b", Timeout: 60}}, IsSafe: true},
		{ID: "T1059", Name: "C", Tactic: "execution", Platforms: []string{"linux"},
			Executors: []MergedExecutor{{Name: "c", Type: "bash", Platform: "linux", Command: "c", Timeout: 60}}, IsSafe: false},
	}

	result, err := WriteYAMLFiles(techniques, tmpDir)
	if err != nil {
		t.Fatalf("WriteYAMLFiles failed: %v", err)
	}

	if result.FileBreakdown["discovery.yaml"] != 2 {
		t.Errorf("discovery.yaml count = %d, want 2", result.FileBreakdown["discovery.yaml"])
	}
	if result.FileBreakdown["execution.yaml"] != 1 {
		t.Errorf("execution.yaml count = %d, want 1", result.FileBreakdown["execution.yaml"])
	}
}

func TestWriteYAMLFiles_UnknownTactic(t *testing.T) {
	tmpDir := t.TempDir()

	techniques := []*MergedTechnique{
		{
			ID: "T9999", Name: "Unknown", Tactic: "",
			Platforms: []string{"linux"},
			Executors: []MergedExecutor{{Name: "test", Type: "bash", Platform: "linux", Command: "echo", Timeout: 60}},
			IsSafe:    true,
		},
	}

	result, err := WriteYAMLFiles(techniques, tmpDir)
	if err != nil {
		t.Fatalf("WriteYAMLFiles failed: %v", err)
	}

	if result.FileBreakdown["unknown.yaml"] != 1 {
		t.Errorf("unknown.yaml count = %d, want 1", result.FileBreakdown["unknown.yaml"])
	}

	if _, err := os.Stat(filepath.Join(tmpDir, "unknown.yaml")); err != nil {
		t.Error("unknown.yaml should be created for empty tactic")
	}
}

func TestWriteYAMLFiles_ReferencesIncluded(t *testing.T) {
	tmpDir := t.TempDir()

	techniques := []*MergedTechnique{
		{
			ID: "T1082", Name: "System Info", Tactic: "discovery",
			Platforms:  []string{"linux"},
			Executors:  []MergedExecutor{{Name: "test", Type: "bash", Platform: "linux", Command: "uname", Timeout: 60}},
			References: []string{"https://attack.mitre.org/techniques/T1082"},
			IsSafe:     true,
		},
	}

	_, err := WriteYAMLFiles(techniques, tmpDir)
	if err != nil {
		t.Fatalf("WriteYAMLFiles failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmpDir, "discovery.yaml"))
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "references:") {
		t.Error("Should contain references")
	}
	if !strings.Contains(content, "attack.mitre.org") {
		t.Error("Should contain the reference URL")
	}
}

func TestPrintDryRunStats_DoesNotPanic(t *testing.T) {
	stats := MergeStats{
		STIXTotal:       100,
		AtomicTotal:     200,
		Matched:         50,
		STIXOnly:        50,
		AtomicOnly:      150,
		ExecutorsTotal:  120,
		SafeCount:       30,
		UnsafeCount:     20,
		TacticBreakdown: map[string]int{"discovery": 10, "execution": 15},
	}

	techniques := []*MergedTechnique{
		{
			ID: "T1082", Name: "System Info", Tactic: "discovery",
			Executors: []MergedExecutor{{Name: "test", Type: "cmd", Platform: "windows", Command: "systeminfo", Timeout: 60}},
			IsSafe:    true,
		},
	}

	// Redirect stdout to avoid test output noise
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintDryRunStats(stats, techniques)

	w.Close()
	os.Stdout = oldStdout

	var buf [8192]byte
	n, _ := r.Read(buf[:])
	output := string(buf[:n])

	if !strings.Contains(output, "STIX techniques total") {
		t.Error("Output should contain STIX techniques total")
	}
	if !strings.Contains(output, "discovery") {
		t.Error("Output should contain tactic name")
	}
}

func TestPrintDryRunStats_EmptyTacticTechnique(t *testing.T) {
	stats := MergeStats{
		TacticBreakdown: map[string]int{},
	}
	techniques := []*MergedTechnique{
		{
			ID: "T9999", Name: "Unknown", Tactic: "",
			Executors: []MergedExecutor{{Name: "test", Type: "bash", Platform: "linux", Command: "echo", Timeout: 60}},
			IsSafe:    true,
		},
	}

	// Redirect stdout
	oldStdout := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w

	PrintDryRunStats(stats, techniques)

	w.Close()
	os.Stdout = oldStdout
}
