package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tr4d3r/cp-config/internal/config"
	"github.com/tr4d3r/cp-config/internal/github"
	"github.com/tr4d3r/cp-config/internal/repo"
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show differences between repo and installed configs",
	Long: `Compare configurations in the repository with those installed in ~/.github/copilot.

Shows which configurations are:
  + Added: exist in repo but not installed
  - Removed: installed but not in repo
  ~ Modified: differ between repo and installed

Examples:
  cp-config diff                    Compare using default profile
  cp-config diff --profile work     Compare using 'work' profile`,
	Run: runDiff,
}

func init() {
	rootCmd.AddCommand(diffCmd)
}

func runDiff(cmd *cobra.Command, args []string) {
	// Load repository
	repository, err := repo.LoadFromCwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Get target directory
	target, err := github.NewTarget()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Get the active profile if specified
	var activeProfile *config.ProfileConfig
	profileName := GetProfile()
	if profileName != "" && profileName != "default" {
		activeProfile = repository.GetProfile(profileName)
		if activeProfile == nil {
			fmt.Fprintf(os.Stderr, "Warning: profile '%s' not found, comparing all configurations\n", profileName)
		}
	}

	// Get filtered agents and skills
	agents := repository.FilteredAgents(activeProfile)
	skills := repository.FilteredSkills(activeProfile)

	// Compute diffs
	agentDiff, err := target.DiffAgents(agents)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error comparing agents: %v\n", err)
		os.Exit(1)
	}

	skillDiff, err := target.DiffSkills(skills)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error comparing skills: %v\n", err)
		os.Exit(1)
	}

	// Display results
	hasChanges := false

	// Show agents diff
	fmt.Println("Agents:")
	agentChanges := printDiffEntries(agentDiff)
	if !agentChanges {
		fmt.Println("  (no changes)")
	}
	hasChanges = hasChanges || agentChanges

	fmt.Println()

	// Show skills diff
	fmt.Println("Skills:")
	skillChanges := printDiffEntries(skillDiff)
	if !skillChanges {
		fmt.Println("  (no changes)")
	}
	hasChanges = hasChanges || skillChanges

	// Summary
	fmt.Println()
	if hasChanges {
		added := countByStatus(agentDiff, github.DiffStatusAdded) + countByStatus(skillDiff, github.DiffStatusAdded)
		modified := countByStatus(agentDiff, github.DiffStatusModified) + countByStatus(skillDiff, github.DiffStatusModified)
		removed := countByStatus(agentDiff, github.DiffStatusRemoved) + countByStatus(skillDiff, github.DiffStatusRemoved)
		fmt.Printf("Summary: %d added, %d modified, %d removed\n", added, modified, removed)
		fmt.Println("Run 'cp-config install' to apply changes.")
	} else {
		fmt.Println("Everything is up to date.")
	}
}

func printDiffEntries(entries []github.DiffEntry) bool {
	hasChanges := false
	for _, entry := range entries {
		switch entry.Status {
		case github.DiffStatusAdded:
			fmt.Printf("  + %s\n", entry.Name)
			hasChanges = true
		case github.DiffStatusModified:
			fmt.Printf("  ~ %s\n", entry.Name)
			hasChanges = true
		case github.DiffStatusRemoved:
			fmt.Printf("  - %s\n", entry.Name)
			hasChanges = true
		}
	}
	return hasChanges
}

func countByStatus(entries []github.DiffEntry, status github.DiffStatus) int {
	count := 0
	for _, e := range entries {
		if e.Status == status {
			count++
		}
	}
	return count
}
