package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tr4d3r/cp-config/internal/config"
)

// ParseSkillFile parses a skill markdown file with frontmatter
func ParseSkillFile(path string) (*config.SkillConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading skill file: %w", err)
	}

	return ParseSkillData(string(data), path)
}

// ParseSkillData parses skill data from a string
func ParseSkillData(data string, sourcePath string) (*config.SkillConfig, error) {
	parsed, err := ExtractFrontmatter(data)
	if err != nil {
		return nil, fmt.Errorf("extracting frontmatter: %w", err)
	}

	var fm config.Frontmatter
	if err := UnmarshalFrontmatter(parsed, &fm); err != nil {
		return nil, fmt.Errorf("parsing frontmatter: %w", err)
	}

	skill := &config.SkillConfig{
		Name:        fm.Name,
		Description: fm.Description,
		Content:     parsed.Content,
		SourcePath:  sourcePath,
	}

	if err := skill.Validate(); err != nil {
		return nil, fmt.Errorf("validating skill: %w", err)
	}

	return skill, nil
}

// LoadSkillsFromDir loads all skill configurations from a directory
// It expects files named {skill-name}.md with frontmatter
func LoadSkillsFromDir(dir string) ([]*config.SkillConfig, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []*config.SkillConfig{}, nil
		}
		return nil, fmt.Errorf("reading skills directory: %w", err)
	}

	var skills []*config.SkillConfig
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		skill, err := ParseSkillFile(path)
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %w", entry.Name(), err)
		}
		skills = append(skills, skill)
	}

	return skills, nil
}
