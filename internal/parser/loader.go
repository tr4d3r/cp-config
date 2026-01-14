package parser

import (
	"fmt"

	"github.com/tr4d3r/cp-config/internal/config"
)

// Repository represents a loaded cp-config repository
type Repository struct {
	// Root is the repository root directory
	Root string

	// Agents contains all loaded agent configurations
	Agents []*config.AgentConfig

	// Skills contains all loaded skill configurations
	Skills []*config.SkillConfig

	// Profiles contains all loaded profile configurations
	Profiles []*config.ProfileConfig
}

// LoadRepository loads all configurations from a cp-config repository
func LoadRepository(cfg *config.AppConfig) (*Repository, error) {
	agents, err := LoadAgentsFromDir(cfg.AgentsDir())
	if err != nil {
		return nil, fmt.Errorf("loading agents: %w", err)
	}

	skills, err := LoadSkillsFromDir(cfg.SkillsDir())
	if err != nil {
		return nil, fmt.Errorf("loading skills: %w", err)
	}

	profiles, err := LoadProfilesFromDir(cfg.ProfilesDir())
	if err != nil {
		return nil, fmt.Errorf("loading profiles: %w", err)
	}

	return &Repository{
		Root:     cfg.RepoRoot,
		Agents:   agents,
		Skills:   skills,
		Profiles: profiles,
	}, nil
}

// GetAgent returns an agent by name, or nil if not found
func (r *Repository) GetAgent(name string) *config.AgentConfig {
	for _, a := range r.Agents {
		if a.Name == name {
			return a
		}
	}
	return nil
}

// GetSkill returns a skill by name, or nil if not found
func (r *Repository) GetSkill(name string) *config.SkillConfig {
	for _, s := range r.Skills {
		if s.Name == name {
			return s
		}
	}
	return nil
}

// GetProfile returns a profile by name, or nil if not found
func (r *Repository) GetProfile(name string) *config.ProfileConfig {
	for _, p := range r.Profiles {
		if p.Name == name {
			return p
		}
	}
	return nil
}

// FilteredAgents returns agents filtered by a profile
func (r *Repository) FilteredAgents(profile *config.ProfileConfig) []*config.AgentConfig {
	if profile == nil {
		return r.Agents
	}

	var filtered []*config.AgentConfig
	for _, agent := range r.Agents {
		if profile.ShouldIncludeAgent(agent.Name) {
			filtered = append(filtered, agent)
		}
	}
	return filtered
}

// FilteredSkills returns skills filtered by a profile
func (r *Repository) FilteredSkills(profile *config.ProfileConfig) []*config.SkillConfig {
	if profile == nil {
		return r.Skills
	}

	var filtered []*config.SkillConfig
	for _, skill := range r.Skills {
		if profile.ShouldIncludeSkill(skill.Name) {
			filtered = append(filtered, skill)
		}
	}
	return filtered
}
