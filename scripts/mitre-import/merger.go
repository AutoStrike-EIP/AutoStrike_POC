package main

import (
	"strings"
)

// MergedTechnique represents a technique after merging STIX metadata with Atomic executors
type MergedTechnique struct {
	ID          string
	Name        string
	Description string
	Tactic      string   // Primary tactic (first in list)
	Tactics     []string // All tactics
	Platforms   []string
	Executors   []MergedExecutor
	Detection   string
	References  []string
	IsSafe      bool
}

// MergedExecutor represents a merged executor ready for YAML output
type MergedExecutor struct {
	Name              string
	Type              string
	Platform          string
	Command           string
	Cleanup           string
	Timeout           int
	ElevationRequired bool
}

// MergeStats holds statistics about the merge operation
type MergeStats struct {
	STIXTotal       int
	AtomicTotal     int
	Matched         int
	STIXOnly        int
	AtomicOnly      int
	ExecutorsTotal  int
	SafeCount       int
	UnsafeCount     int
	TacticBreakdown map[string]int
}

// dangerousPatterns are command patterns that force is_safe=false
var dangerousPatterns = []string{
	"rm -rf", "del /f", "format c:", "format d:", "shutdown",
	"encrypt", "mkfs", "dd if=", "fdisk", "wipefs",
	"cipher /w", "sdelete", "shred", "> /dev/sd",
	"Remove-Item -Recurse -Force", "Stop-Service",
	"Stop-Process", "net stop", "taskkill /f",
}

// safeTactics are tactics where techniques are safe by default
var safeTactics = map[string]bool{
	"discovery":      true,
	"reconnaissance": true,
}

// Merge performs an inner join between STIX techniques and Atomic executors
func Merge(stix map[string]*STIXTechnique, atomics map[string]*AtomicTechnique) ([]*MergedTechnique, MergeStats) {
	stats := MergeStats{
		STIXTotal:       len(stix),
		AtomicTotal:     len(atomics),
		TacticBreakdown: make(map[string]int),
	}

	// Track which IDs are in both
	stixOnly := 0
	atomicOnly := 0

	for id := range stix {
		if _, ok := atomics[id]; !ok {
			stixOnly++
		}
	}
	for id := range atomics {
		if _, ok := stix[id]; !ok {
			atomicOnly++
		}
	}
	stats.STIXOnly = stixOnly
	stats.AtomicOnly = atomicOnly

	var result []*MergedTechnique

	for id, stixTech := range stix {
		atomicTech, ok := atomics[id]
		if !ok {
			continue // Inner join: skip STIX-only
		}

		// Filter executors to only include supported platforms from STIX
		platformSet := make(map[string]bool)
		for _, p := range stixTech.Platforms {
			platformSet[p] = true
		}

		var executors []MergedExecutor
		for _, exec := range atomicTech.Executors {
			if !platformSet[exec.Platform] {
				continue
			}
			executors = append(executors, MergedExecutor{
				Name:              exec.Name,
				Type:              exec.Type,
				Platform:          exec.Platform,
				Command:           exec.Command,
				Cleanup:           exec.Cleanup,
				Timeout:           defaultTimeout(exec.Type),
				ElevationRequired: exec.ElevationRequired,
			})
		}

		if len(executors) == 0 {
			continue // No executable tests for supported platforms
		}

		stats.Matched++
		stats.ExecutorsTotal += len(executors)

		// Determine primary tactic and all tactics
		tactic := ""
		if len(stixTech.Tactics) > 0 {
			tactic = stixTech.Tactics[0]
		}

		// Determine is_safe (check ALL tactics, not just primary)
		isSafe := determineSafety(stixTech.Tactics, executors)

		if isSafe {
			stats.SafeCount++
		} else {
			stats.UnsafeCount++
		}

		// Update tactic breakdown using primary tactic
		if tactic != "" {
			stats.TacticBreakdown[tactic]++
		}

		merged := &MergedTechnique{
			ID:          id,
			Name:        stixTech.Name,
			Description: truncateDescription(stixTech.Description),
			Tactic:      tactic,
			Tactics:     stixTech.Tactics,
			Platforms:   stixTech.Platforms,
			Executors:   executors,
			References:  stixTech.References,
			IsSafe:      isSafe,
		}

		result = append(result, merged)
	}

	return result, stats
}

// determineSafety determines if a technique is safe based on ALL tactics and commands
func determineSafety(tactics []string, executors []MergedExecutor) bool {
	// Check if any executor requires elevation
	for _, exec := range executors {
		if exec.ElevationRequired {
			return false
		}
	}

	// Check for dangerous command patterns
	for _, exec := range executors {
		cmdLower := strings.ToLower(exec.Command)
		cleanupLower := strings.ToLower(exec.Cleanup)
		for _, pattern := range dangerousPatterns {
			patternLower := strings.ToLower(pattern)
			if strings.Contains(cmdLower, patternLower) || strings.Contains(cleanupLower, patternLower) {
				return false
			}
		}
	}

	// ALL tactics must be safe; if any tactic is unsafe, the technique is unsafe
	if len(tactics) == 0 {
		return false
	}
	for _, tactic := range tactics {
		if !safeTactics[tactic] {
			return false
		}
	}
	return true
}

// defaultTimeout returns a default timeout based on executor type
func defaultTimeout(execType string) int {
	switch execType {
	case "powershell", "cmd":
		return 120
	default:
		return 60
	}
}

// truncateDescription truncates long descriptions for YAML readability
func truncateDescription(desc string) string {
	const maxLen = 500
	if len(desc) <= maxLen {
		return desc
	}
	return desc[:maxLen] + "..."
}
