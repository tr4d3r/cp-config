package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version is set at build time
	Version = "dev"

	// Commit is set at build time
	Commit = "none"

	// cfgFile stores the config file path
	cfgFile string

	// profileFlag stores the active profile name from command line
	profileFlag string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cp-config",
	Short: "Manage GitHub Copilot configurations",
	Long: `cp-config is a CLI tool for managing GitHub Copilot agent and skill configurations.

It enables source-controlled management of Copilot configurations with deployment
to your ~/.github/copilot directory, supporting sharing and versioning of agents
and skills across machines and teams.

Examples:
  cp-config list agents           List available agents
  cp-config list skills           List available skills
  cp-config install               Deploy configurations to ~/.github/copilot
  cp-config diff                  Show differences between repo and installed configs
  cp-config profile list          List available profiles`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .cp-config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&profileFlag, "profile", "p", "default", "profile to use")

	// Add version command
	rootCmd.AddCommand(versionCmd)
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("cp-config %s (commit: %s)\n", Version, Commit)
	},
}

// GetProfile returns the currently selected profile
func GetProfile() string {
	return profileFlag
}
