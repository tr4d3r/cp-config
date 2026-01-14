package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Deploy configurations to ~/.github/copilot",
	Long: `Deploy agent and skill configurations from the repository to ~/.github/copilot.

This command copies all agents and skills (filtered by the active profile) to your
local GitHub Copilot configuration directory.

Examples:
  cp-config install                    Install using default profile
  cp-config install --profile work     Install using 'work' profile`,
	Run: runInstall,
}

func init() {
	rootCmd.AddCommand(installCmd)
}

func runInstall(cmd *cobra.Command, args []string) {
	fmt.Printf("Installing configurations with profile: %s\n", GetProfile())
	fmt.Println("(not yet implemented)")
}
