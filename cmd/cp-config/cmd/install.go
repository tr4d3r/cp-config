package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tr4d3r/cp-config/internal/config"
	"github.com/tr4d3r/cp-config/internal/github"
	"github.com/tr4d3r/cp-config/internal/repo"
)

var (
	dryRun bool
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Deploy configurations to ~/.github/copilot",
	Long: `Deploy agent and skill configurations from the repository to ~/.github/copilot.

This command copies all agents and skills (filtered by the active profile) to your
local GitHub Copilot configuration directory.

Examples:
  cp-config install                    Install using default profile
  cp-config install --profile work     Install using 'work' profile
  cp-config install --dry-run          Show what would be installed`,
	Run: runInstall,
}

func init() {
	installCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be installed without making changes")
	rootCmd.AddCommand(installCmd)
}

func runInstall(cmd *cobra.Command, args []string) {
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
			fmt.Fprintf(os.Stderr, "Warning: profile '%s' not found, installing all configurations\n", profileName)
		}
	}

	// Get filtered agents and skills
	agents := repository.FilteredAgents(activeProfile)
	skills := repository.FilteredSkills(activeProfile)

	if len(agents) == 0 && len(skills) == 0 {
		fmt.Println("No configurations to install.")
		return
	}

	if dryRun {
		fmt.Println("Dry run - would install:")
		fmt.Println()
	}

	// Install agents
	if len(agents) > 0 {
		fmt.Println("Agents:")
		for _, agent := range agents {
			if dryRun {
				fmt.Printf("  %s -> %s\n", agent.Name, target.AgentPath(agent.Name))
			} else {
				if err := target.DeployAgent(agent); err != nil {
					fmt.Fprintf(os.Stderr, "Error installing agent %s: %v\n", agent.Name, err)
					os.Exit(1)
				}
				fmt.Printf("  ✓ %s\n", agent.Name)
			}
		}
	}

	// Install skills
	if len(skills) > 0 {
		if len(agents) > 0 {
			fmt.Println()
		}
		fmt.Println("Skills:")
		for _, skill := range skills {
			if dryRun {
				fmt.Printf("  %s -> %s\n", skill.Name, target.SkillPath(skill.Name))
			} else {
				if err := target.DeploySkill(skill); err != nil {
					fmt.Fprintf(os.Stderr, "Error installing skill %s: %v\n", skill.Name, err)
					os.Exit(1)
				}
				fmt.Printf("  ✓ %s\n", skill.Name)
			}
		}
	}

	// Summary
	fmt.Println()
	if dryRun {
		fmt.Printf("Would install %d agent(s) and %d skill(s) to %s\n", len(agents), len(skills), target.Root)
	} else {
		fmt.Printf("Installed %d agent(s) and %d skill(s) to %s\n", len(agents), len(skills), target.Root)
	}
}
