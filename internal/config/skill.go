package config

import "path/filepath"

// SkillConfig represents a GitHub Copilot skill configuration
// Skills are markdown files with YAML frontmatter containing instructions
// for specific tasks like debugging, code review, etc.
type SkillConfig struct {
	// Name is the skill identifier (e.g., "debugging", "code-review")
	Name string `yaml:"name" json:"name"`

	// Description describes when and how to use the skill
	Description string `yaml:"description" json:"description"`

	// Content is the markdown body with skill instructions
	Content string `yaml:"-" json:"content"`

	// SourcePath is the path to the source file in the repository
	SourcePath string `yaml:"-" json:"source_path"`
}

// FileName returns the expected filename for this skill
func (s *SkillConfig) FileName() string {
	return s.Name + ".md"
}

// TargetPath returns the deployment path relative to ~/.github/copilot
func (s *SkillConfig) TargetPath() string {
	return filepath.Join("skills", s.FileName())
}

// Validate checks if the skill configuration is valid
func (s *SkillConfig) Validate() error {
	if s.Name == "" {
		return &ValidationError{Field: "name", Message: "skill name is required"}
	}
	if s.Description == "" {
		return &ValidationError{Field: "description", Message: "skill description is required"}
	}
	return nil
}

// ValidationError represents a configuration validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}
