package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/tr4d3r/cp-config/internal/config"
	"github.com/tr4d3r/cp-config/internal/parser"
	"github.com/tr4d3r/cp-config/internal/repo"
)

var listCmd = &cobra.Command{
	Use:   "list [agents|skills]",
	Short: "List available configurations",
	Long: `List available agent and skill configurations from the repository.

Examples:
  cp-config list              List all agents and skills
  cp-config list agents       List only agents
  cp-config list skills       List only skills`,
	Args:      cobra.MaximumNArgs(1),
	ValidArgs: []string{"agents", "skills"},
	Run:       runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) {
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
			fmt.Fprintf(os.Stderr, "Warning: profile '%s' not found, using all configurations\n", profileName)
		}
	}

	configType := ""
	if len(args) > 0 {
		configType = args[0]
	}

	switch configType {
	case "":
		listAgents(repository, activeProfile)
		fmt.Println()
		listSkills(repository, activeProfile)
	case "agents":
		listAgents(repository, activeProfile)
	case "skills":
		listSkills(repository, activeProfile)
	default:
		fmt.Fprintf(os.Stderr, "Unknown type: %s (use 'agents' or 'skills')\n", configType)
		os.Exit(1)
	}
}

func listAgents(repository *parser.Repository, profile *config.ProfileConfig) {
	agents := repository.FilteredAgents(profile)

	fmt.Println("Agents:")
	if len(agents) == 0 {
		fmt.Println("  (no agents configured)")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "  NAME\tDESCRIPTION")
	fmt.Fprintln(w, "  ----\t-----------")
	for _, agent := range agents {
		desc := truncate(agent.Description, 60)
		fmt.Fprintf(w, "  %s\t%s\n", agent.Name, desc)
	}
	w.Flush()
}

func listSkills(repository *parser.Repository, profile *config.ProfileConfig) {
	skills := repository.FilteredSkills(profile)

	fmt.Println("Skills:")
	if len(skills) == 0 {
		fmt.Println("  (no skills configured)")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "  NAME\tDESCRIPTION")
	fmt.Fprintln(w, "  ----\t-----------")
	for _, skill := range skills {
		desc := truncate(skill.Description, 60)
		fmt.Fprintf(w, "  %s\t%s\n", skill.Name, desc)
	}
	w.Flush()
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
