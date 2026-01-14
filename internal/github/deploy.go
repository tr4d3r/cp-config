package github

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tr4d3r/cp-config/internal/config"
)

// DeployAgent writes an agent configuration to the target directory
func (t *Target) DeployAgent(agent *config.AgentConfig) error {
	if err := t.EnsureDirectories(); err != nil {
		return err
	}

	content := formatConfigFile(agent.Name, agent.Description, agent.Content)
	path := t.AgentPath(agent.Name)

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("writing agent file: %w", err)
	}

	return nil
}

// DeploySkill writes a skill configuration to the target directory
func (t *Target) DeploySkill(skill *config.SkillConfig) error {
	if err := t.EnsureDirectories(); err != nil {
		return err
	}

	content := formatConfigFile(skill.Name, skill.Description, skill.Content)
	path := t.SkillPath(skill.Name)

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("writing skill file: %w", err)
	}

	return nil
}

// DeployAgents deploys multiple agents to the target directory
func (t *Target) DeployAgents(agents []*config.AgentConfig) error {
	for _, agent := range agents {
		if err := t.DeployAgent(agent); err != nil {
			return fmt.Errorf("deploying agent %s: %w", agent.Name, err)
		}
	}
	return nil
}

// DeploySkills deploys multiple skills to the target directory
func (t *Target) DeploySkills(skills []*config.SkillConfig) error {
	for _, skill := range skills {
		if err := t.DeploySkill(skill); err != nil {
			return fmt.Errorf("deploying skill %s: %w", skill.Name, err)
		}
	}
	return nil
}

// RemoveAgent removes an agent from the target directory
func (t *Target) RemoveAgent(name string) error {
	path := t.AgentPath(name)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("removing agent file: %w", err)
	}
	return nil
}

// RemoveSkill removes a skill from the target directory
func (t *Target) RemoveSkill(name string) error {
	path := t.SkillPath(name)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("removing skill file: %w", err)
	}
	return nil
}

// DeployFromSource copies a file directly from a source path
func (t *Target) DeployFromSource(sourcePath, destSubdir, filename string) error {
	if err := t.EnsureDirectories(); err != nil {
		return err
	}

	data, err := os.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("reading source file: %w", err)
	}

	destPath := filepath.Join(t.Root, destSubdir, filename)
	if err := os.WriteFile(destPath, data, 0644); err != nil {
		return fmt.Errorf("writing destination file: %w", err)
	}

	return nil
}

// formatConfigFile creates a markdown file with frontmatter
func formatConfigFile(name, description, content string) string {
	return fmt.Sprintf(`---
name: %s
description: %s
---
%s`, name, description, content)
}
