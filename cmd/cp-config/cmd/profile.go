package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage configuration profiles",
	Long: `Manage configuration profiles for different contexts (work, personal, etc).

Profiles allow you to maintain different sets of agents and skills for different
purposes and switch between them easily.`,
}

var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available profiles",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Available profiles:")
		fmt.Println("  * default")
		fmt.Println("(not yet implemented)")
	},
}

var profileCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new profile",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Creating profile: %s\n", args[0])
		fmt.Println("(not yet implemented)")
	},
}

var profileSwitchCmd = &cobra.Command{
	Use:   "switch <name>",
	Short: "Switch to a different profile",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Switching to profile: %s\n", args[0])
		fmt.Println("(not yet implemented)")
	},
}

var profileDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a profile",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Deleting profile: %s\n", args[0])
		fmt.Println("(not yet implemented)")
	},
}

func init() {
	rootCmd.AddCommand(profileCmd)
	profileCmd.AddCommand(profileListCmd)
	profileCmd.AddCommand(profileCreateCmd)
	profileCmd.AddCommand(profileSwitchCmd)
	profileCmd.AddCommand(profileDeleteCmd)
}
