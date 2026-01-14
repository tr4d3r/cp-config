package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tr4d3r/cp-config/internal/config"
)

// ParseAgentFile parses an agent markdown file with frontmatter
func ParseAgentFile(path string) (*config.AgentConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading agent file: %w", err)
	}

	return ParseAgentData(string(data), path)
}

// ParseAgentData parses agent data from a string
func ParseAgentData(data string, sourcePath string) (*config.AgentConfig, error) {
	parsed, err := ExtractFrontmatter(data)
	if err != nil {
		return nil, fmt.Errorf("extracting frontmatter: %w", err)
	}

	var fm config.Frontmatter
	if err := UnmarshalFrontmatter(parsed, &fm); err != nil {
		return nil, fmt.Errorf("parsing frontmatter: %w", err)
	}

	agent := &config.AgentConfig{
		Name:        fm.Name,
		Description: fm.Description,
		Content:     parsed.Content,
		SourcePath:  sourcePath,
	}

	if err := agent.Validate(); err != nil {
		return nil, fmt.Errorf("validating agent: %w", err)
	}

	return agent, nil
}

// LoadAgentsFromDir loads all agent configurations from a directory
// It expects files named {agent-name}.md with frontmatter
func LoadAgentsFromDir(dir string) ([]*config.AgentConfig, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []*config.AgentConfig{}, nil
		}
		return nil, fmt.Errorf("reading agents directory: %w", err)
	}

	var agents []*config.AgentConfig
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		agent, err := ParseAgentFile(path)
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %w", entry.Name(), err)
		}
		agents = append(agents, agent)
	}

	return agents, nil
}
