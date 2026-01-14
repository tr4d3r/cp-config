package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tr4d3r/cp-config/internal/repo"
)

var (
	forceRemove bool
)

var removeCmd = &cobra.Command{
	Use:     "remove <type> <name>",
	Aliases: []string{"rm"},
	Short:   "Remove an agent or skill configuration",
	Long: `Remove an agent or skill configuration from the repository.

Type must be either 'agent' or 'skill'. This deletes the configuration file
from the repository. Use --force to skip confirmation.

Examples:
  cp-config remove agent my-assistant
  cp-config remove skill code-review
  cp-config rm agent old-helper --force`,
	Args:      cobra.ExactArgs(2),
	ValidArgs: []string{"agent", "skill"},
	Run:       runRemove,
}

func init() {
	removeCmd.Flags().BoolVarP(&forceRemove, "force", "f", false, "Skip confirmation prompt")
	rootCmd.AddCommand(removeCmd)
}

func runRemove(cmd *cobra.Command, args []string) {
	configType := args[0]
	name := args[1]

	if configType != "agent" && configType != "skill" {
		fmt.Fprintf(os.Stderr, "Error: type must be 'agent' or 'skill', got '%s'\n", configType)
		os.Exit(1)
	}

	root, err := repo.FindRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	var dir string
	if configType == "agent" {
		dir = filepath.Join(root, "agents")
	} else {
		dir = filepath.Join(root, "skills")
	}

	filePath := filepath.Join(dir, name+".md")

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: %s '%s' not found at %s\n", configType, name, filePath)
		os.Exit(1)
	}

	// Confirm deletion unless --force
	if !forceRemove {
		fmt.Printf("Remove %s '%s' from %s? [y/N] ", configType, name, filePath)
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" && response != "yes" && response != "Yes" {
			fmt.Println("Aborted.")
			return
		}
	}

	// Remove the file
	if err := os.Remove(filePath); err != nil {
		fmt.Fprintf(os.Stderr, "Error removing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Removed %s: %s\n", configType, name)
}
