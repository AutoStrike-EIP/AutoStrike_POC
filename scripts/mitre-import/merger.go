package main


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
			// Safety is determined solely by elevation_required from Atomic Red Team
			execSafe := !exec.ElevationRequired
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
