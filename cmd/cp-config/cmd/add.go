package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tr4d3r/cp-config/internal/repo"
)

var addCmd = &cobra.Command{
	Use:   "add <type> <name>",
	Short: "Create a new agent or skill configuration",
	Long: `Create a new agent or skill configuration from a template.

Type must be either 'agent' or 'skill'. The configuration file will be created
in the appropriate directory (agents/ or skills/) with a basic template.

Examples:
  cp-config add agent my-assistant
  cp-config add skill code-review`,
	Args:      cobra.ExactArgs(2),
	ValidArgs: []string{"agent", "skill"},
	Run:       runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func runAdd(cmd *cobra.Command, args []string) {
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

	var dir, template string
	if configType == "agent" {
		dir = filepath.Join(root, "agents")
		template = agentTemplate(name)
	} else {
		dir = filepath.Join(root, "skills")
		template = skillTemplate(name)
	}

	// Ensure directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating directory: %v\n", err)
		os.Exit(1)
	}

	// Check if file already exists
	filePath := filepath.Join(dir, name+".md")
	if _, err := os.Stat(filePath); err == nil {
		fmt.Fprintf(os.Stderr, "Error: %s '%s' already exists at %s\n", configType, name, filePath)
		os.Exit(1)
	}

	// Write template
	if err := os.WriteFile(filePath, []byte(template), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created %s: %s\n", configType, filePath)
	fmt.Printf("Edit the file to customize the %s configuration.\n", configType)
}

func agentTemplate(name string) string {
	return fmt.Sprintf(`---
name: %s
description: Description of what this agent does
---

# %s Agent

## Purpose

Describe the purpose and capabilities of this agent.

## Instructions

Provide instructions for how this agent should behave.

## Examples

Include examples of how to use this agent.
`, name, name)
}

func skillTemplate(name string) string {
	return fmt.Sprintf(`---
name: %s
description: Description of when and how to use this skill
---

# %s Skill

## Overview

Describe what this skill helps with.

## Guidelines

1. First guideline
2. Second guideline
3. Third guideline

## Best Practices

- Best practice 1
- Best practice 2
`, name, name)
}
