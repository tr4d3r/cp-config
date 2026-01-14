package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show differences between repo and installed configs",
	Long: `Compare configurations in the repository with those installed in ~/.github/copilot.

Shows which configurations are:
  - Added: exist in repo but not installed
  - Removed: installed but not in repo
  - Modified: differ between repo and installed

Examples:
  cp-config diff                    Compare using default profile
  cp-config diff --profile work     Compare using 'work' profile`,
	Run: runDiff,
}

func init() {
	rootCmd.AddCommand(diffCmd)
}

func runDiff(cmd *cobra.Command, args []string) {
	fmt.Printf("Comparing configurations with profile: %s\n", GetProfile())
	fmt.Println("(not yet implemented)")
}
