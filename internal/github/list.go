package github

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/tr4d3r/cp-config/internal/config"
	"github.com/tr4d3r/cp-config/internal/parser"
)

// InstalledConfig represents a configuration file installed in the target directory
type InstalledConfig struct {
	// Name is the configuration name (derived from filename)
	Name string

	// Path is the full path to the installed file
	Path string

	// Type is the configuration type (agent or skill)
	Type config.ConfigType
}

// ListInstalledAgents returns all agents installed in the target directory
func (t *Target) ListInstalledAgents() ([]InstalledConfig, error) {
	return t.listInstalled(t.AgentsDir, config.ConfigTypeAgent)
}

// ListInstalledSkills returns all skills installed in the target directory
func (t *Target) ListInstalledSkills() ([]InstalledConfig, error) {
	return t.listInstalled(t.SkillsDir, config.ConfigTypeSkill)
}

// listInstalled returns all configurations in a directory
func (t *Target) listInstalled(dir string, configType config.ConfigType) ([]InstalledConfig, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []InstalledConfig{}, nil
		}
		return nil, err
	}

	var configs []InstalledConfig
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), ".md")
		configs = append(configs, InstalledConfig{
			Name: name,
			Path: filepath.Join(dir, entry.Name()),
			Type: configType,
		})
	}

	return configs, nil
}

// LoadInstalledAgents loads and parses all installed agents
func (t *Target) LoadInstalledAgents() ([]*config.AgentConfig, error) {
	return parser.LoadAgentsFromDir(t.AgentsDir)
}

// LoadInstalledSkills loads and parses all installed skills
func (t *Target) LoadInstalledSkills() ([]*config.SkillConfig, error) {
	return parser.LoadSkillsFromDir(t.SkillsDir)
}

// HasAgent checks if an agent is installed
func (t *Target) HasAgent(name string) bool {
	_, err := os.Stat(t.AgentPath(name))
	return err == nil
}

// HasSkill checks if a skill is installed
func (t *Target) HasSkill(name string) bool {
	_, err := os.Stat(t.SkillPath(name))
	return err == nil
}

// ReadAgent reads and returns the content of an installed agent
func (t *Target) ReadAgent(name string) ([]byte, error) {
	return os.ReadFile(t.AgentPath(name))
}

// ReadSkill reads and returns the content of an installed skill
func (t *Target) ReadSkill(name string) ([]byte, error) {
	return os.ReadFile(t.SkillPath(name))
}
