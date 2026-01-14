package config

import (
	"os"
	"path/filepath"
)

// AppConfig represents the application-level configuration
type AppConfig struct {
	// RepoRoot is the root directory of the cp-config repository
	RepoRoot string `yaml:"repo_root" json:"repo_root"`

	// TargetDir is the target directory for deployment (defaults to ~/.github/copilot)
	TargetDir string `yaml:"target_dir" json:"target_dir"`

	// ActiveProfile is the currently active profile name
	ActiveProfile string `yaml:"active_profile" json:"active_profile"`
}

// DefaultAppConfig returns the default application configuration
func DefaultAppConfig() (*AppConfig, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return &AppConfig{
		RepoRoot:      "",
		TargetDir:     filepath.Join(homeDir, ".github", "copilot"),
		ActiveProfile: "default",
	}, nil
}

// AgentsDir returns the path to the agents directory in the repository
func (c *AppConfig) AgentsDir() string {
	return filepath.Join(c.RepoRoot, "agents")
}

// SkillsDir returns the path to the skills directory in the repository
func (c *AppConfig) SkillsDir() string {
	return filepath.Join(c.RepoRoot, "skills")
}

// ProfilesDir returns the path to the profiles directory in the repository
func (c *AppConfig) ProfilesDir() string {
	return filepath.Join(c.RepoRoot, "profiles")
}

// TargetAgentsDir returns the target directory for agent deployment
func (c *AppConfig) TargetAgentsDir() string {
	return filepath.Join(c.TargetDir, "agents")
}

// TargetSkillsDir returns the target directory for skill deployment
func (c *AppConfig) TargetSkillsDir() string {
	return filepath.Join(c.TargetDir, "skills")
}

// Paths holds standard directory and file paths
type Paths struct {
	// RepoRoot is the repository root directory
	RepoRoot string

	// Agents is the agents directory in the repo
	Agents string

	// Skills is the skills directory in the repo
	Skills string

	// Profiles is the profiles directory in the repo
	Profiles string

	// Target is the deployment target directory (~/.github/copilot)
	Target string

	// TargetAgents is the agents directory in the target
	TargetAgents string

	// TargetSkills is the skills directory in the target
	TargetSkills string
}

// NewPaths creates a Paths instance from an AppConfig
func NewPaths(cfg *AppConfig) *Paths {
	return &Paths{
		RepoRoot:     cfg.RepoRoot,
		Agents:       cfg.AgentsDir(),
		Skills:       cfg.SkillsDir(),
		Profiles:     cfg.ProfilesDir(),
		Target:       cfg.TargetDir,
		TargetAgents: cfg.TargetAgentsDir(),
		TargetSkills: cfg.TargetSkillsDir(),
	}
}
