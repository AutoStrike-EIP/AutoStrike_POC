package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	stixURL    = "https://raw.githubusercontent.com/mitre/cti/master/enterprise-attack/enterprise-attack.json"
	atomicRepo = "https://github.com/redcanaryco/atomic-red-team.git"
)

func main() {
	stixPath := flag.String("stix-path", "", "Local path to enterprise-attack.json (downloads if not set)")
	atomicsPath := flag.String("atomics-path", "", "Local path to atomic-red-team repo (clones if not set)")
	outputDir := flag.String("output-dir", "../../server/configs/techniques", "Output directory for YAML files")
	cacheDir := flag.String("cache-dir", "", "Cache directory (default: ~/.cache/autostrike)")
	dryRun := flag.Bool("dry-run", false, "Print stats without writing files")
	safeOnly := flag.Bool("safe-only", false, "Import only safe techniques")
	forceDownload := flag.Bool("force-download", false, "Re-download even if cache exists")
	verbose := flag.Bool("verbose", false, "Verbose logging")
	flag.Parse()

	if *cacheDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Failed to get home directory: %v", err)
		}
		*cacheDir = filepath.Join(home, ".cache", "autostrike")
	}

	if err := os.MkdirAll(*cacheDir, 0755); err != nil {
		log.Fatalf("Failed to create cache directory: %v", err)
	}

	// Resolve STIX path
	resolvedSTIXPath := *stixPath
	if resolvedSTIXPath == "" {
		resolvedSTIXPath = filepath.Join(*cacheDir, "enterprise-attack.json")
		if *forceDownload || !fileExists(resolvedSTIXPath) {
			fmt.Println("Downloading MITRE ATT&CK STIX data...")
			if err := downloadFile(stixURL, resolvedSTIXPath); err != nil {
				log.Fatalf("Failed to download STIX data: %v", err)
			}
			fmt.Println("Download complete.")
		} else {
			fmt.Println("Using cached STIX data.")
		}
	}

	// Resolve Atomics path
	resolvedAtomicsPath := *atomicsPath
	if resolvedAtomicsPath == "" {
		repoDir := filepath.Join(*cacheDir, "atomic-red-team")
		atomicsDir := filepath.Join(repoDir, "atomics")
		if *forceDownload || !fileExists(atomicsDir) {
			fmt.Println("Cloning Atomic Red Team repository (shallow)...")
			if err := cloneAtomics(atomicRepo, repoDir); err != nil {
				log.Fatalf("Failed to clone Atomic Red Team: %v", err)
			}
			fmt.Println("Clone complete.")
		} else {
			fmt.Println("Using cached Atomic Red Team data.")
		}
		resolvedAtomicsPath = atomicsDir
	} else {
		// User-provided path points to the repo root; append /atomics
		atomicsDir := filepath.Join(resolvedAtomicsPath, "atomics")
		if fileExists(atomicsDir) {
			resolvedAtomicsPath = atomicsDir
		}
	}

	// Parse STIX
	fmt.Println("Parsing STIX data...")
	stixTechniques, err := ParseSTIX(resolvedSTIXPath)
	if err != nil {
		log.Fatalf("Failed to parse STIX: %v", err)
	}
	if *verbose {
		fmt.Printf("  Parsed %d STIX techniques\n", len(stixTechniques))
	}

	// Parse Atomics
	fmt.Println("Parsing Atomic Red Team data...")
	atomicTechniques, err := ParseAtomics(resolvedAtomicsPath)
	if err != nil {
		log.Fatalf("Failed to parse Atomics: %v", err)
	}
	if *verbose {
		fmt.Printf("  Parsed %d Atomic techniques\n", len(atomicTechniques))
	}

	// Merge
	fmt.Println("Merging techniques (inner join)...")
	merged, stats := Merge(stixTechniques, atomicTechniques)

	// Filter safe-only if requested
	if *safeOnly {
		var filtered []*MergedTechnique
		filteredExecutors := 0
		for _, tech := range merged {
			if tech.IsSafe {
				filtered = append(filtered, tech)
				filteredExecutors += len(tech.Executors)
			}
		}
		fmt.Printf("Filtered to %d safe techniques (from %d total)\n", len(filtered), len(merged))
		merged = filtered
		stats.ExecutorsTotal = filteredExecutors
		stats.SafeCount = len(filtered)
		stats.UnsafeCount = 0
	}

	if *dryRun {
		PrintDryRunStats(stats, merged)
		return
	}

	// Write output
	absOutput, err := filepath.Abs(*outputDir)
	if err != nil {
		log.Fatalf("Failed to resolve output path: %v", err)
	}
	fmt.Printf("Writing %d techniques to %s...\n", len(merged), absOutput)
	result, err := WriteYAMLFiles(merged, absOutput)
	if err != nil {
		log.Fatalf("Failed to write YAML files: %v", err)
	}

	fmt.Println()
	fmt.Printf("Import complete: %d techniques with %d executors across %d files\n",
		result.TotalTechniques, stats.ExecutorsTotal, result.FilesWritten)
	fmt.Printf("Safe: %d, Unsafe: %d\n", stats.SafeCount, stats.UnsafeCount)
}

// fileExists checks if a path exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// downloadFile downloads a URL to a local file
func downloadFile(url, dest string) (retErr error) {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("HTTP GET failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := out.Close(); retErr == nil {
			retErr = cerr
		}
		// Remove partial file on error to avoid cache poisoning
		if retErr != nil {
			os.Remove(dest)
		}
	}()

	if _, retErr = io.Copy(out, resp.Body); retErr != nil {
		return retErr
	}
	return retErr
}

// cloneAtomics performs a shallow git clone of the Atomic Red Team repository
func cloneAtomics(repoURL, dest string) error {
	// Remove existing directory if force download
	if fileExists(dest) {
		if err := os.RemoveAll(dest); err != nil {
			return fmt.Errorf("failed to remove existing directory: %w", err)
		}
	}

	cmd := exec.Command("git", "clone", "--depth", "1", repoURL, dest)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
