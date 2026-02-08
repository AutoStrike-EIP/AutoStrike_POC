package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseAtomicData_BasicTest(t *testing.T) {
	data := []byte(`
attack_technique: T1082
display_name: System Information Discovery
atomic_tests:
  - name: "System Info via systeminfo"
    auto_generated_guid: "abc-123"
    description: "Get system info"
    supported_platforms:
      - windows
    executor:
      name: command_prompt
      command: |
        systeminfo
      elevation_required: false
`)

	tech, err := ParseAtomicData(data, "T1082")
	if err != nil {
		t.Fatalf("ParseAtomicData failed: %v", err)
	}

	if tech.ID != "T1082" {
		t.Errorf("ID = %s, want T1082", tech.ID)
	}
	if len(tech.Executors) != 1 {
		t.Fatalf("Expected 1 executor, got %d", len(tech.Executors))
	}

	exec := tech.Executors[0]
	if exec.Type != "cmd" {
		t.Errorf("Type = %s, want cmd", exec.Type)
	}
	if exec.Platform != "windows" {
		t.Errorf("Platform = %s, want windows", exec.Platform)
	}
	if exec.Name != "System Info via systeminfo" {
		t.Errorf("Name = %s, want System Info via systeminfo", exec.Name)
	}
}

func TestParseAtomicData_WithCleanupAndElevation(t *testing.T) {
	data := []byte(`
attack_technique: T1053.005
display_name: Scheduled Task
atomic_tests:
  - name: "Create Scheduled Task"
    auto_generated_guid: "def-456"
    description: "Create a scheduled task"
    supported_platforms:
      - windows
    executor:
      name: powershell
      command: |
        schtasks /create /tn "Test" /tr "cmd /c echo test" /sc once /st 00:00
      cleanup_command: |
        schtasks /delete /tn "Test" /f
      elevation_required: true
`)

	tech, err := ParseAtomicData(data, "T1053.005")
	if err != nil {
		t.Fatalf("ParseAtomicData failed: %v", err)
	}

	if len(tech.Executors) != 1 {
		t.Fatalf("Expected 1 executor, got %d", len(tech.Executors))
	}

	exec := tech.Executors[0]
	if exec.Type != "powershell" {
		t.Errorf("Type = %s, want powershell", exec.Type)
	}
	if exec.Cleanup == "" {
		t.Error("Expected non-empty cleanup command")
	}
	if !exec.ElevationRequired {
		t.Error("Expected elevation_required = true")
	}
}

func TestParseAtomicData_WithInputArguments(t *testing.T) {
	data := []byte(`
attack_technique: T1083
display_name: File and Directory Discovery
atomic_tests:
  - name: "File Discovery"
    auto_generated_guid: "ghi-789"
    description: "Discover files"
    supported_platforms:
      - linux
    input_arguments:
      output_file:
        description: "Output file"
        type: path
        default: /tmp/T1083.txt
      search_path:
        description: "Path to search"
        type: path
        default: /home
    executor:
      name: bash
      command: |
        find #{search_path} -type f > #{output_file}
      cleanup_command: |
        rm #{output_file}
      elevation_required: false
`)

	tech, err := ParseAtomicData(data, "T1083")
	if err != nil {
		t.Fatalf("ParseAtomicData failed: %v", err)
	}

	exec := tech.Executors[0]
	if exec.Command != "find /home -type f > /tmp/T1083.txt" {
		t.Errorf("Command template not resolved: %s", exec.Command)
	}
	if exec.Cleanup != "rm /tmp/T1083.txt" {
		t.Errorf("Cleanup template not resolved: %s", exec.Cleanup)
	}
}

func TestParseAtomicData_SkipManualExecutor(t *testing.T) {
	data := []byte(`
attack_technique: T1001
display_name: Manual Test
atomic_tests:
  - name: "Manual test"
    auto_generated_guid: "xyz-000"
    description: "Requires manual steps"
    supported_platforms:
      - windows
    executor:
      name: manual
      command: ""
`)

	tech, err := ParseAtomicData(data, "T1001")
	if err != nil {
		t.Fatalf("ParseAtomicData failed: %v", err)
	}

	if len(tech.Executors) != 0 {
		t.Errorf("Expected 0 executors (manual skipped), got %d", len(tech.Executors))
	}
}

func TestParseAtomicData_MultiPlatform(t *testing.T) {
	data := []byte(`
attack_technique: T1082
display_name: System Info
atomic_tests:
  - name: "Cross-platform info"
    auto_generated_guid: "abc-000"
    description: "test"
    supported_platforms:
      - windows
      - linux
    executor:
      name: bash
      command: uname -a
`)

	tech, err := ParseAtomicData(data, "T1082")
	if err != nil {
		t.Fatalf("ParseAtomicData failed: %v", err)
	}

	if len(tech.Executors) != 2 {
		t.Fatalf("Expected 2 executors (one per platform), got %d", len(tech.Executors))
	}

	platforms := make(map[string]bool)
	for _, exec := range tech.Executors {
		platforms[exec.Platform] = true
	}
	if !platforms["windows"] || !platforms["linux"] {
		t.Errorf("Expected windows and linux platforms, got %v", platforms)
	}
}

func TestParseAtomicData_DuplicateNames(t *testing.T) {
	data := []byte(`
attack_technique: T1082
display_name: System Info
atomic_tests:
  - name: "System Info"
    auto_generated_guid: "abc-001"
    description: "test1"
    supported_platforms:
      - windows
    executor:
      name: command_prompt
      command: systeminfo
  - name: "System Info"
    auto_generated_guid: "abc-002"
    description: "test2"
    supported_platforms:
      - windows
    executor:
      name: powershell
      command: Get-ComputerInfo
`)

	tech, err := ParseAtomicData(data, "T1082")
	if err != nil {
		t.Fatalf("ParseAtomicData failed: %v", err)
	}

	if len(tech.Executors) != 2 {
		t.Fatalf("Expected 2 executors, got %d", len(tech.Executors))
	}

	// Second should have deduplicated name
	if tech.Executors[0].Name != "System Info" {
		t.Errorf("First name = %s, want System Info", tech.Executors[0].Name)
	}
	if tech.Executors[1].Name != "System Info (2)" {
		t.Errorf("Second name = %s, want System Info (2)", tech.Executors[1].Name)
	}
}

func TestResolveTemplates(t *testing.T) {
	tests := []struct {
		name    string
		command string
		args    map[string]InputArgument
		want    string
	}{
		{
			name:    "no templates",
			command: "whoami",
			args:    nil,
			want:    "whoami",
		},
		{
			name:    "single template with default",
			command: "cat #{file}",
			args:    map[string]InputArgument{"file": {Default: "/etc/passwd"}},
			want:    "cat /etc/passwd",
		},
		{
			name:    "template without default",
			command: "cat #{file}",
			args:    map[string]InputArgument{},
			want:    "cat ",
		},
		{
			name:    "multiple templates",
			command: "cp #{src} #{dst}",
			args: map[string]InputArgument{
				"src": {Default: "/tmp/a"},
				"dst": {Default: "/tmp/b"},
			},
			want: "cp /tmp/a /tmp/b",
		},
		{
			name:    "empty command",
			command: "",
			args:    nil,
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveTemplates(tt.command, tt.args)
			if got != tt.want {
				t.Errorf("resolveTemplates() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseAtomicData_InvalidYAML(t *testing.T) {
	_, err := ParseAtomicData([]byte(`invalid: [yaml`), "T1000")
	if err == nil {
		t.Error("Expected error for invalid YAML")
	}
}

func TestParseAtomicData_UnsupportedPlatformFiltered(t *testing.T) {
	data := []byte(`
attack_technique: T1000
display_name: Test
atomic_tests:
  - name: "Cloud test"
    auto_generated_guid: "abc-999"
    description: "test"
    supported_platforms:
      - azure
    executor:
      name: bash
      command: echo test
`)

	tech, err := ParseAtomicData(data, "T1000")
	if err != nil {
		t.Fatalf("ParseAtomicData failed: %v", err)
	}

	if len(tech.Executors) != 0 {
		t.Errorf("Expected 0 executors (unsupported platform), got %d", len(tech.Executors))
	}
}

func TestParseAtomicFile_ValidFile(t *testing.T) {
	tmpDir := t.TempDir()
	yamlContent := []byte(`
attack_technique: T1082
display_name: System Information Discovery
atomic_tests:
  - name: "systeminfo"
    auto_generated_guid: "abc-123"
    description: "test"
    supported_platforms:
      - windows
    executor:
      name: command_prompt
      command: systeminfo
`)
	path := filepath.Join(tmpDir, "T1082.yaml")
	if err := os.WriteFile(path, yamlContent, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	tech, err := ParseAtomicFile(path, "T1082")
	if err != nil {
		t.Fatalf("ParseAtomicFile failed: %v", err)
	}

	if tech.ID != "T1082" {
		t.Errorf("ID = %s, want T1082", tech.ID)
	}
	if len(tech.Executors) != 1 {
		t.Errorf("Expected 1 executor, got %d", len(tech.Executors))
	}
}

func TestParseAtomicFile_FileNotFound(t *testing.T) {
	_, err := ParseAtomicFile("/nonexistent/path.yaml", "T9999")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}

func TestParseAtomics_DirectoryWithTechniques(t *testing.T) {
	tmpDir := t.TempDir()

	// Create T1082 directory with YAML file
	t1082Dir := filepath.Join(tmpDir, "T1082")
	os.Mkdir(t1082Dir, 0755)
	os.WriteFile(filepath.Join(t1082Dir, "T1082.yaml"), []byte(`
attack_technique: T1082
display_name: System Information Discovery
atomic_tests:
  - name: "systeminfo"
    auto_generated_guid: "abc-123"
    description: "test"
    supported_platforms:
      - windows
    executor:
      name: command_prompt
      command: systeminfo
`), 0644)

	// Create T1083 directory with .yml file
	t1083Dir := filepath.Join(tmpDir, "T1083")
	os.Mkdir(t1083Dir, 0755)
	os.WriteFile(filepath.Join(t1083Dir, "T1083.yml"), []byte(`
attack_technique: T1083
display_name: File and Directory Discovery
atomic_tests:
  - name: "ls"
    auto_generated_guid: "def-456"
    description: "test"
    supported_platforms:
      - linux
    executor:
      name: bash
      command: ls -la /home
`), 0644)

	// Create a non-technique directory (should be skipped)
	os.Mkdir(filepath.Join(tmpDir, "Indexes"), 0755)

	techniques, err := ParseAtomics(tmpDir)
	if err != nil {
		t.Fatalf("ParseAtomics failed: %v", err)
	}

	if len(techniques) != 2 {
		t.Errorf("Expected 2 techniques, got %d", len(techniques))
	}
	if techniques["T1082"] == nil {
		t.Error("T1082 should be parsed")
	}
	if techniques["T1083"] == nil {
		t.Error("T1083 should be parsed")
	}
}

func TestParseAtomics_EmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	techniques, err := ParseAtomics(tmpDir)
	if err != nil {
		t.Fatalf("ParseAtomics failed: %v", err)
	}
	if len(techniques) != 0 {
		t.Errorf("Expected 0 techniques for empty dir, got %d", len(techniques))
	}
}

func TestParseAtomics_NonexistentDirectory(t *testing.T) {
	_, err := ParseAtomics("/nonexistent/directory")
	if err == nil {
		t.Error("Expected error for nonexistent directory")
	}
}

func TestParseAtomics_SkipsInvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	t1000Dir := filepath.Join(tmpDir, "T1000")
	os.Mkdir(t1000Dir, 0755)
	os.WriteFile(filepath.Join(t1000Dir, "T1000.yaml"), []byte("invalid: [yaml"), 0644)

	techniques, err := ParseAtomics(tmpDir)
	if err != nil {
		t.Fatalf("ParseAtomics failed: %v", err)
	}
	if len(techniques) != 0 {
		t.Errorf("Expected 0 techniques (invalid YAML skipped), got %d", len(techniques))
	}
}

func TestParseAtomics_SkipsManualOnlyTechniques(t *testing.T) {
	tmpDir := t.TempDir()
	t1000Dir := filepath.Join(tmpDir, "T1000")
	os.Mkdir(t1000Dir, 0755)
	os.WriteFile(filepath.Join(t1000Dir, "T1000.yaml"), []byte(`
attack_technique: T1000
display_name: Manual Only
atomic_tests:
  - name: "Manual"
    auto_generated_guid: "xxx-000"
    description: "Manual only"
    supported_platforms:
      - windows
    executor:
      name: manual
      command: ""
`), 0644)

	techniques, err := ParseAtomics(tmpDir)
	if err != nil {
		t.Fatalf("ParseAtomics failed: %v", err)
	}
	if len(techniques) != 0 {
		t.Errorf("Expected 0 techniques (manual-only skipped), got %d", len(techniques))
	}
}
