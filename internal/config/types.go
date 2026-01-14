package config

// ConfigType represents the type of configuration (agent or skill)
type ConfigType string

const (
	ConfigTypeAgent ConfigType = "agent"
	ConfigTypeSkill ConfigType = "skill"
)

// Frontmatter represents the YAML frontmatter in markdown config files
type Frontmatter struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// ConfigFile represents a parsed configuration file (agent or skill)
type ConfigFile struct {
	// Path is the file path relative to the repository root
	Path string

	// Type indicates whether this is an agent or skill
	Type ConfigType

	// Frontmatter contains the parsed YAML frontmatter
	Frontmatter Frontmatter

	// Content is the markdown body (excluding frontmatter)
	Content string

	// Raw is the complete file contents
	Raw string
}
