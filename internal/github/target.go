package github

import (
	"fmt"
	"os"
	"path/filepath"
)

// Target represents the ~/.github/copilot target directory
type Target struct {
	// Root is the root target directory (typically ~/.github/copilot)
	Root string

	// AgentsDir is the agents subdirectory
	AgentsDir string

	// SkillsDir is the skills subdirectory
	SkillsDir string
}

// NewTarget creates a new Target with the default ~/.github/copilot path
func NewTarget() (*Target, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("getting home directory: %w", err)
	}

	return NewTargetAt(filepath.Join(homeDir, ".github", "copilot")), nil
}

// NewTargetAt creates a new Target at the specified path
func NewTargetAt(root string) *Target {
	return &Target{
		Root:      root,
		AgentsDir: filepath.Join(root, "agents"),
		SkillsDir: filepath.Join(root, "skills"),
	}
}

// Exists checks if the target directory exists
func (t *Target) Exists() bool {
	_, err := os.Stat(t.Root)
	return err == nil
}

// EnsureDirectories creates the target directory structure if it doesn't exist
func (t *Target) EnsureDirectories() error {
	if err := os.MkdirAll(t.AgentsDir, 0755); err != nil {
		return fmt.Errorf("creating agents directory: %w", err)
	}
	if err := os.MkdirAll(t.SkillsDir, 0755); err != nil {
		return fmt.Errorf("creating skills directory: %w", err)
	}
	return nil
}

// AgentPath returns the full path for an agent file
func (t *Target) AgentPath(name string) string {
	return filepath.Join(t.AgentsDir, name+".md")
}

// SkillPath returns the full path for a skill file
func (t *Target) SkillPath(name string) string {
	return filepath.Join(t.SkillsDir, name+".md")
}
