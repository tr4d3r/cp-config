package config

// ProfileConfig represents a named configuration profile
// Profiles allow different sets of agents and skills for different contexts
// (e.g., work, personal, project-specific)
type ProfileConfig struct {
	// Name is the profile identifier (e.g., "work", "personal", "default")
	Name string `yaml:"name" json:"name"`

	// Description describes the profile's purpose
	Description string `yaml:"description" json:"description"`

	// Agents specifies which agents to include or exclude
	Agents FilterConfig `yaml:"agents" json:"agents"`

	// Skills specifies which skills to include or exclude
	Skills FilterConfig `yaml:"skills" json:"skills"`

	// SourcePath is the path to the source file in the repository
	SourcePath string `yaml:"-" json:"source_path"`
}

// FilterConfig specifies include/exclude lists for agents or skills
type FilterConfig struct {
	// Include lists the names to include (if empty, include all)
	Include []string `yaml:"include" json:"include"`

	// Exclude lists the names to exclude (applied after include)
	Exclude []string `yaml:"exclude" json:"exclude"`
}

// FileName returns the expected filename for this profile
func (p *ProfileConfig) FileName() string {
	return p.Name + ".yaml"
}

// Validate checks if the profile configuration is valid
func (p *ProfileConfig) Validate() error {
	if p.Name == "" {
		return &ValidationError{Field: "name", Message: "profile name is required"}
	}
	return nil
}

// ShouldIncludeAgent determines if an agent should be included based on filters
func (p *ProfileConfig) ShouldIncludeAgent(name string) bool {
	return shouldInclude(name, p.Agents)
}

// ShouldIncludeSkill determines if a skill should be included based on filters
func (p *ProfileConfig) ShouldIncludeSkill(name string) bool {
	return shouldInclude(name, p.Skills)
}

// shouldInclude checks if a name passes the include/exclude filters
func shouldInclude(name string, filter FilterConfig) bool {
	// If include list is specified and non-empty, name must be in it
	if len(filter.Include) > 0 {
		found := false
		for _, inc := range filter.Include {
			if inc == name {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Check exclude list
	for _, exc := range filter.Exclude {
		if exc == name {
			return false
		}
	}

	return true
}

// DefaultProfile returns the default profile configuration
func DefaultProfile() *ProfileConfig {
	return &ProfileConfig{
		Name:        "default",
		Description: "Default profile - includes all agents and skills",
		Agents:      FilterConfig{Include: []string{}, Exclude: []string{}},
		Skills:      FilterConfig{Include: []string{}, Exclude: []string{}},
	}
}
