package main

import (
	"regexp"
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
	IsSafe            bool
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
		hasSafeExecutor := false
		for _, exec := range atomicTech.Executors {
			if !platformSet[exec.Platform] {
				continue
			}
			// Safety: no elevation AND no dangerous command patterns
			execSafe := !exec.ElevationRequired && !hasDangerousPattern(exec.Command) && !hasDangerousPattern(exec.Cleanup)
			if execSafe {
				hasSafeExecutor = true
			}
			executors = append(executors, MergedExecutor{
				Name:              exec.Name,
				Type:              exec.Type,
				Platform:          exec.Platform,
				Command:           exec.Command,
				Cleanup:           exec.Cleanup,
				Timeout:           defaultTimeout(exec.Type),
				ElevationRequired: exec.ElevationRequired,
				IsSafe:            execSafe,
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

		// Technique is safe if it has at least one safe executor
		isSafe := hasSafeExecutor

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

// defaultTimeout returns a default timeout based on executor type
func defaultTimeout(execType string) int {
	switch execType {
	case "powershell", "cmd":
		return 120
	default:
		return 60
	}
}

// truncateDescription truncates extremely long descriptions for YAML readability
func truncateDescription(desc string) string {
	const maxLen = 6000
	if len(desc) <= maxLen {
		return desc
	}
	return desc[:maxLen] + "..."
}

// dangerousPatterns matches command patterns that are destructive or dangerous
var dangerousPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\brm\s+-[^\s]*r[^\s]*f`),      // rm -rf, rm -fr
	regexp.MustCompile(`(?i)\brm\s+-[^\s]*f[^\s]*r`),      // rm -fr variant
	regexp.MustCompile(`(?i)\bdel\s+/[^\s]*f`),             // del /f (Windows)
	regexp.MustCompile(`(?i)\brd\s+/s`),                    // rd /s (Windows recursive delete)
	regexp.MustCompile(`(?i)\brmdir\s+/s`),                 // rmdir /s (Windows)
	regexp.MustCompile(`(?i)\bdd\s+.*\bof=/dev/`),          // dd of=/dev/sda etc.
	regexp.MustCompile(`(?i)\bmkfs\b`),                     // mkfs (format filesystem)
	regexp.MustCompile(`(?i)\bfdisk\b`),                    // fdisk (partition table)
	regexp.MustCompile(`(?i)\bformat\s+[a-z]:`),            // format C: (Windows)
	regexp.MustCompile(`(?i)\bshutdown\b`),                 // shutdown
	regexp.MustCompile(`(?i)\breboot\b`),                   // reboot
	regexp.MustCompile(`(?i)\binit\s+0\b`),                 // init 0 (halt)
	regexp.MustCompile(`(?i)\btaskkill\s+/f`),              // taskkill /f (force kill)
	regexp.MustCompile(`(?i)\bkill\s+-9\b`),                // kill -9
	regexp.MustCompile(`(?i)\bkillall\b`),                  // killall
	regexp.MustCompile(`(?i)\bpkill\b`),                    // pkill
	regexp.MustCompile(`(?i)>\s*/dev/sd[a-z]`),             // > /dev/sda (overwrite disk)
	regexp.MustCompile(`(?i)\bsystemctl\s+(stop|disable)`), // systemctl stop/disable
	regexp.MustCompile(`(?i)\bchmod\s+000\b`),              // chmod 000 (remove all perms)
	regexp.MustCompile(`(?i)\biptables\s+-F\b`),            // iptables -F (flush rules)
}

// hasDangerousPattern checks if a command contains any destructive patterns
func hasDangerousPattern(command string) bool {
	if strings.TrimSpace(command) == "" {
		return false
	}
	for _, pattern := range dangerousPatterns {
		if pattern.MatchString(command) {
			return true
		}
	}
	return false
}
