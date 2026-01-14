package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tr4d3r/cp-config/internal/config"
	"github.com/tr4d3r/cp-config/internal/github"
	"github.com/tr4d3r/cp-config/internal/repo"
)

var (
	syncDryRun bool
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Bidirectional sync between repo and installed configs",
	Long: `Sync configurations between the repository and installed configs.

This command performs a bidirectional sync:
- Deploys repo configs that are not installed
- Offers to import installed configs that are not in the repo
- Prompts for conflict resolution when both versions differ

Use --dry-run to preview changes without applying them.

Examples:
  cp-config sync                  Sync all configurations
  cp-config sync --dry-run        Preview sync changes
  cp-config sync -p work          Sync using the work profile`,
	Run: runSync,
}

func init() {
	syncCmd.Flags().BoolVar(&syncDryRun, "dry-run", false, "Preview changes without applying them")
	rootCmd.AddCommand(syncCmd)
}

var stdinReader *bufio.Reader

func runSync(cmd *cobra.Command, args []string) {
	stdinReader = bufio.NewReader(os.Stdin)

	// Load repository
	repository, err := repo.LoadFromCwd()
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
			fmt.Fprintf(os.Stderr, "Warning: profile '%s' not found, syncing all configurations\n", profileName)
		}
	}

	// Get filtered agents and skills
	agents := repository.FilteredAgents(activeProfile)
	skills := repository.FilteredSkills(activeProfile)

	// Set up target
	target, err := github.NewTarget()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Get diffs
	agentDiffs, err := target.DiffAgents(agents)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error diffing agents: %v\n", err)
		os.Exit(1)
	}

	skillDiffs, err := target.DiffSkills(skills)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error diffing skills: %v\n", err)
		os.Exit(1)
	}

	result := &github.DiffResult{
		Agents: agentDiffs,
		Skills: skillDiffs,
	}

	if !result.HasChanges() {
		fmt.Println("Everything is in sync.")
		return
	}

	repoRoot, _ := repo.FindRoot()

	// Process each category
	var actions []syncAction

	// Handle items that need to be deployed (in repo, not installed)
	for _, entry := range result.Added() {
		actions = append(actions, syncAction{
			entry:     entry,
			operation: "deploy",
			desc:      fmt.Sprintf("Deploy %s '%s' to installed", entry.Type, entry.Name),
		})
	}

	// Handle items that are installed but not in repo
	for _, entry := range result.Removed() {
		action := promptRemovedAction(entry, syncDryRun)
		if action.operation != "skip" {
			actions = append(actions, action)
		}
	}

	// Handle conflicts (both modified)
	for _, entry := range result.Modified() {
		action := promptConflictAction(entry, syncDryRun)
		if action.operation != "skip" {
			actions = append(actions, action)
		}
	}

	if len(actions) == 0 {
		fmt.Println("No actions to perform.")
		return
	}

	// Show summary
	fmt.Println("\nActions to perform:")
	for _, a := range actions {
		fmt.Printf("  • %s\n", a.desc)
	}

	if syncDryRun {
		fmt.Println("\n(dry-run) No changes made.")
		return
	}

	// Execute actions
	fmt.Println()
	for _, a := range actions {
		if err := executeAction(a, target, repoRoot, repository); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		} else {
			fmt.Printf("✓ %s\n", a.desc)
		}
	}
}

type syncAction struct {
	entry     github.DiffEntry
	operation string // "deploy", "import", "remove", "skip"
	desc      string
}

func promptRemovedAction(entry github.DiffEntry, dryRun bool) syncAction {
	if dryRun {
		return syncAction{
			entry:     entry,
			operation: "import",
			desc:      fmt.Sprintf("Import %s '%s' to repo (would prompt)", entry.Type, entry.Name),
		}
	}

	fmt.Printf("\n%s '%s' is installed but not in repo.\n", strings.Title(string(entry.Type)), entry.Name)
	fmt.Printf("  [i]mport to repo, [r]emove from installed, [s]kip? ")

	response := readResponse()
	switch strings.ToLower(response) {
	case "i", "import":
		return syncAction{
			entry:     entry,
			operation: "import",
			desc:      fmt.Sprintf("Import %s '%s' to repo", entry.Type, entry.Name),
		}
	case "r", "remove":
		return syncAction{
			entry:     entry,
			operation: "remove",
			desc:      fmt.Sprintf("Remove %s '%s' from installed", entry.Type, entry.Name),
		}
	default:
		return syncAction{entry: entry, operation: "skip"}
	}
}

func promptConflictAction(entry github.DiffEntry, dryRun bool) syncAction {
	if dryRun {
		return syncAction{
			entry:     entry,
			operation: "deploy",
			desc:      fmt.Sprintf("Overwrite installed %s '%s' with repo version (would prompt)", entry.Type, entry.Name),
		}
	}

	fmt.Printf("\n%s '%s' differs between repo and installed.\n", strings.Title(string(entry.Type)), entry.Name)
	fmt.Printf("  [r]epo wins (deploy), [i]nstalled wins (import), [s]kip? ")

	response := readResponse()
	switch strings.ToLower(response) {
	case "r", "repo":
		return syncAction{
			entry:     entry,
			operation: "deploy",
			desc:      fmt.Sprintf("Overwrite installed %s '%s' with repo version", entry.Type, entry.Name),
		}
	case "i", "installed":
		return syncAction{
			entry:     entry,
			operation: "import",
			desc:      fmt.Sprintf("Import installed %s '%s' to repo", entry.Type, entry.Name),
		}
	default:
		return syncAction{entry: entry, operation: "skip"}
	}
}

func readResponse() string {
	response, _ := stdinReader.ReadString('\n')
	return strings.TrimSpace(response)
}

func executeAction(action syncAction, target *github.Target, repoRoot string, repository interface{}) error {
	switch action.operation {
	case "deploy":
		return deployEntry(action.entry, target)
	case "import":
		return importEntry(action.entry, repoRoot)
	case "remove":
		return removeEntry(action.entry, target)
	}
	return nil
}

func deployEntry(entry github.DiffEntry, target *github.Target) error {
	if entry.RepoPath == "" {
		return fmt.Errorf("no repo path for %s", entry.Name)
	}

	subdir := "agents"
	if entry.Type == config.ConfigTypeSkill {
		subdir = "skills"
	}

	return target.DeployFromSource(entry.RepoPath, subdir, entry.Name+".md")
}

func importEntry(entry github.DiffEntry, repoRoot string) error {
	if entry.TargetPath == "" {
		return fmt.Errorf("no installed path for %s", entry.Name)
	}

	// Read installed file
	data, err := os.ReadFile(entry.TargetPath)
	if err != nil {
		return fmt.Errorf("reading installed file: %w", err)
	}

	// Determine destination
	subdir := "agents"
	if entry.Type == config.ConfigTypeSkill {
		subdir = "skills"
	}

	destDir := filepath.Join(repoRoot, subdir)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	destPath := filepath.Join(destDir, entry.Name+".md")
	if err := os.WriteFile(destPath, data, 0644); err != nil {
		return fmt.Errorf("writing to repo: %w", err)
	}

	return nil
}

func removeEntry(entry github.DiffEntry, target *github.Target) error {
	if entry.Type == config.ConfigTypeAgent {
		return target.RemoveAgent(entry.Name)
	}
	return target.RemoveSkill(entry.Name)
}
