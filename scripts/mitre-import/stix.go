package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// STIXBundle represents the top-level STIX bundle
type STIXBundle struct {
	Objects []json.RawMessage `json:"objects"`
}

// STIXAttackPattern represents a MITRE ATT&CK technique from STIX
type STIXAttackPattern struct {
	Type               string              `json:"type"`
	Name               string              `json:"name"`
	Description        string              `json:"description"`
	ExternalReferences []ExternalReference  `json:"external_references"`
	KillChainPhases    []KillChainPhase     `json:"kill_chain_phases"`
	Platforms          []string             `json:"x_mitre_platforms"`
	Deprecated         bool                 `json:"x_mitre_deprecated"`
	Revoked            bool                 `json:"revoked"`
	IsSubtechnique     bool                 `json:"x_mitre_is_subtechnique"`
}

// ExternalReference represents an external reference in STIX
type ExternalReference struct {
	SourceName string `json:"source_name"`
	ExternalID string `json:"external_id"`
	URL        string `json:"url"`
}

// KillChainPhase represents a kill chain phase in STIX
type KillChainPhase struct {
	KillChainName string `json:"kill_chain_name"`
	PhaseName     string `json:"phase_name"`
}

// STIXTechnique is the parsed/normalized result from a STIX attack pattern
type STIXTechnique struct {
	ID          string
	Name        string
	Description string
	Tactics     []string // Normalized tactic names (e.g., "discovery", "initial-access")
	Platforms   []string // Normalized: "windows", "linux", "macos"
	References  []string // URLs
}

// ParseSTIX parses STIX JSON from a file and returns normalized techniques
func ParseSTIX(path string) (map[string]*STIXTechnique, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read STIX file: %w", err)
	}

	return ParseSTIXData(data)
}

// ParseSTIXData parses STIX JSON from bytes
func ParseSTIXData(data []byte) (map[string]*STIXTechnique, error) {
	var bundle STIXBundle
	if err := json.Unmarshal(data, &bundle); err != nil {
		return nil, fmt.Errorf("failed to parse STIX JSON: %w", err)
	}

	techniques := make(map[string]*STIXTechnique)

	for _, raw := range bundle.Objects {
		// Quick check for type field
		var typeCheck struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(raw, &typeCheck); err != nil {
			continue
		}
		if typeCheck.Type != "attack-pattern" {
			continue
		}

		var ap STIXAttackPattern
		if err := json.Unmarshal(raw, &ap); err != nil {
			continue
		}

		// Skip revoked and deprecated
		if ap.Revoked || ap.Deprecated {
			continue
		}

		// Extract MITRE ATT&CK ID and build citation map
		id := ""
		var refs []string
		citationURLs := make(map[string]string)
		for _, ref := range ap.ExternalReferences {
			if ref.SourceName == "mitre-attack" {
				if ref.ExternalID != "" {
					id = ref.ExternalID
				}
				if ref.URL != "" {
					refs = append(refs, ref.URL)
				}
			} else if ref.URL != "" {
				citationURLs[ref.SourceName] = ref.URL
			}
		}
		if id == "" {
			continue
		}

		// Extract tactics
		var tactics []string
		for _, phase := range ap.KillChainPhases {
			if phase.KillChainName == "mitre-attack" {
				tactics = append(tactics, phase.PhaseName)
			}
		}

		// Normalize platforms
		platforms := normalizePlatforms(ap.Platforms)
		if len(platforms) == 0 {
			continue // Skip cloud-only techniques
		}

		// Resolve (Citation: Name) to markdown links
		description := resolveCitations(ap.Description, citationURLs)

		techniques[id] = &STIXTechnique{
			ID:          id,
			Name:        ap.Name,
			Description: description,
			Tactics:     tactics,
			Platforms:   platforms,
			References:  refs,
		}
	}

	return techniques, nil
}

// citationRegex matches (Citation: Source Name) patterns in STIX descriptions
var citationRegex = regexp.MustCompile(`\(Citation: ([^)]+)\)`)

// resolveCitations replaces (Citation: Name) with [Name](url) when a URL is available,
// or removes the citation markup if no URL exists.
func resolveCitations(description string, citationURLs map[string]string) string {
	return citationRegex.ReplaceAllStringFunc(description, func(match string) string {
		sub := citationRegex.FindStringSubmatch(match)
		if len(sub) < 2 {
			return match
		}
		name := sub[1]
		if url, ok := citationURLs[name]; ok {
			return "[" + name + "](" + url + ")"
		}
		// No URL: keep as plain text reference
		return "(Ref: " + name + ")"
	})
}

// normalizePlatforms normalizes STIX platform names to lowercase and filters non-supported ones
func normalizePlatforms(platforms []string) []string {
	supported := map[string]string{
		"windows": "windows",
		"linux":   "linux",
		"macos":   "macos",
	}

	seen := make(map[string]bool)
	var result []string
	for _, p := range platforms {
		normalized, ok := supported[strings.ToLower(p)]
		if ok && !seen[normalized] {
			seen[normalized] = true
			result = append(result, normalized)
		}
	}
	return result
}
