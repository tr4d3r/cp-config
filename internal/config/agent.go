package config

import "path/filepath"

// AgentConfig represents a GitHub Copilot agent configuration
// Agents are markdown files with YAML frontmatter that define
// specialized assistants with specific capabilities and behaviors.
type AgentConfig struct {
	// Name is the agent identifier (e.g., "code-reviewer", "documentation-helper")
	Name string `yaml:"name" json:"name"`

	// Description describes the agent's purpose and capabilities
	Description string `yaml:"description" json:"description"`

	// Content is the markdown body with agent instructions and behavior
	Content string `yaml:"-" json:"content"`

	// SourcePath is the path to the source file in the repository
	SourcePath string `yaml:"-" json:"source_path"`
}

// FileName returns the expected filename for this agent
func (a *AgentConfig) FileName() string {
	return a.Name + ".md"
}

// TargetPath returns the deployment path relative to ~/.github/copilot
func (a *AgentConfig) TargetPath() string {
	return filepath.Join("agents", a.FileName())
}

// Validate checks if the agent configuration is valid
func (a *AgentConfig) Validate() error {
	if a.Name == "" {
		return &ValidationError{Field: "name", Message: "agent name is required"}
	}
	if a.Description == "" {
		return &ValidationError{Field: "description", Message: "agent description is required"}
	}
	return nil
}
