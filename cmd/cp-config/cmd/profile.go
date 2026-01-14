package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/tr4d3r/cp-config/internal/config"
	"github.com/tr4d3r/cp-config/internal/profile"
	"github.com/tr4d3r/cp-config/internal/repo"
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
	Run:   runProfileList,
}

var profileCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new profile",
	Long: `Create a new profile with the given name.

The profile will be created with empty include/exclude lists, meaning it will
include all agents and skills by default. Edit the profile YAML file to customize.

Examples:
  cp-config profile create work
  cp-config profile create personal`,
	Args: cobra.ExactArgs(1),
	Run:  runProfileCreate,
}

var profileSwitchCmd = &cobra.Command{
	Use:   "switch <name>",
	Short: "Switch to a different profile",
	Long: `Switch the active profile to the specified one.

The active profile is used by default for install, diff, and list commands.

Examples:
  cp-config profile switch work
  cp-config profile switch default`,
	Args: cobra.ExactArgs(1),
	Run:  runProfileSwitch,
}

var profileDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a profile",
	Long: `Delete the specified profile.

The default profile cannot be deleted. If the deleted profile was active,
the active profile will be reset to 'default'.

Examples:
  cp-config profile delete work`,
	Args: cobra.ExactArgs(1),
	Run:  runProfileDelete,
}

var profileShowCmd = &cobra.Command{
	Use:   "show [name]",
	Short: "Show profile details",
	Long: `Show details of a profile including its include/exclude lists.

If no name is provided, shows the active profile.

Examples:
  cp-config profile show
  cp-config profile show work`,
	Args: cobra.MaximumNArgs(1),
	Run:  runProfileShow,
}

func init() {
	rootCmd.AddCommand(profileCmd)
	profileCmd.AddCommand(profileListCmd)
	profileCmd.AddCommand(profileCreateCmd)
	profileCmd.AddCommand(profileSwitchCmd)
	profileCmd.AddCommand(profileDeleteCmd)
	profileCmd.AddCommand(profileShowCmd)
}

func getProfileManager() *profile.Manager {
	root, err := repo.FindRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	return profile.NewManager(root)
}

func runProfileList(cmd *cobra.Command, args []string) {
	mgr := getProfileManager()

	profiles, err := mgr.List()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing profiles: %v\n", err)
		os.Exit(1)
	}

	active, _ := mgr.GetActiveProfile()

	fmt.Println("Profiles:")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "  NAME\tDESCRIPTION")
	fmt.Fprintln(w, "  ----\t-----------")

	for _, p := range profiles {
		marker := " "
		if p.Name == active {
			marker = "*"
		}
		desc := p.Description
		if desc == "" {
			desc = "(no description)"
		}
		fmt.Fprintf(w, " %s%s\t%s\n", marker, p.Name, truncate(desc, 50))
	}
	w.Flush()

	fmt.Printf("\nActive profile: %s\n", active)
}

func runProfileCreate(cmd *cobra.Command, args []string) {
	name := args[0]
	mgr := getProfileManager()

	newProfile := &config.ProfileConfig{
		Name:        name,
		Description: fmt.Sprintf("Profile: %s", name),
		Agents:      config.FilterConfig{Include: []string{}, Exclude: []string{}},
		Skills:      config.FilterConfig{Include: []string{}, Exclude: []string{}},
	}

	if err := mgr.Create(newProfile); err != nil {
		if err == profile.ErrProfileExists {
			fmt.Fprintf(os.Stderr, "Error: profile '%s' already exists\n", name)
		} else {
			fmt.Fprintf(os.Stderr, "Error creating profile: %v\n", err)
		}
		os.Exit(1)
	}

	fmt.Printf("Created profile: %s\n", name)
	fmt.Printf("Edit profiles/%s.yaml to customize include/exclude lists.\n", name)
}

func runProfileSwitch(cmd *cobra.Command, args []string) {
	name := args[0]
	mgr := getProfileManager()

	if err := mgr.SetActiveProfile(name); err != nil {
		if err == profile.ErrProfileNotFound {
			fmt.Fprintf(os.Stderr, "Error: profile '%s' not found\n", name)
		} else {
			fmt.Fprintf(os.Stderr, "Error switching profile: %v\n", err)
		}
		os.Exit(1)
	}

	fmt.Printf("Switched to profile: %s\n", name)
}

func runProfileDelete(cmd *cobra.Command, args []string) {
	name := args[0]
	mgr := getProfileManager()

	if err := mgr.Delete(name); err != nil {
		switch err {
		case profile.ErrProfileNotFound:
			fmt.Fprintf(os.Stderr, "Error: profile '%s' not found\n", name)
		case profile.ErrCannotDeleteDefault:
			fmt.Fprintf(os.Stderr, "Error: cannot delete the default profile\n")
		default:
			fmt.Fprintf(os.Stderr, "Error deleting profile: %v\n", err)
		}
		os.Exit(1)
	}

	fmt.Printf("Deleted profile: %s\n", name)
}

func runProfileShow(cmd *cobra.Command, args []string) {
	mgr := getProfileManager()

	var name string
	if len(args) > 0 {
		name = args[0]
	} else {
		var err error
		name, err = mgr.GetActiveProfile()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting active profile: %v\n", err)
			os.Exit(1)
		}
	}

	p, err := mgr.Get(name)
	if err != nil {
		if err == profile.ErrProfileNotFound {
			fmt.Fprintf(os.Stderr, "Error: profile '%s' not found\n", name)
		} else {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		os.Exit(1)
	}

	fmt.Printf("Profile: %s\n", p.Name)
	if p.Description != "" {
		fmt.Printf("Description: %s\n", p.Description)
	}
	fmt.Println()

	fmt.Println("Agents:")
	if len(p.Agents.Include) > 0 {
		fmt.Printf("  Include: %v\n", p.Agents.Include)
	} else {
		fmt.Println("  Include: (all)")
	}
	if len(p.Agents.Exclude) > 0 {
		fmt.Printf("  Exclude: %v\n", p.Agents.Exclude)
	}

	fmt.Println()
	fmt.Println("Skills:")
	if len(p.Skills.Include) > 0 {
		fmt.Printf("  Include: %v\n", p.Skills.Include)
	} else {
		fmt.Println("  Include: (all)")
	}
	if len(p.Skills.Exclude) > 0 {
		fmt.Printf("  Exclude: %v\n", p.Skills.Exclude)
	}
}
