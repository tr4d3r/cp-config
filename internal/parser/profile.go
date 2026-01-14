package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tr4d3r/cp-config/internal/config"
	"gopkg.in/yaml.v3"
)

// ParseProfileFile parses a profile YAML file
func ParseProfileFile(path string) (*config.ProfileConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading profile file: %w", err)
	}

	return ParseProfileData(string(data), path)
}

// ParseProfileData parses profile data from a string
func ParseProfileData(data string, sourcePath string) (*config.ProfileConfig, error) {
	var profile config.ProfileConfig
	if err := yaml.Unmarshal([]byte(data), &profile); err != nil {
		return nil, fmt.Errorf("parsing profile YAML: %w", err)
	}

	profile.SourcePath = sourcePath

	if err := profile.Validate(); err != nil {
		return nil, fmt.Errorf("validating profile: %w", err)
	}

	return &profile, nil
}

// LoadProfilesFromDir loads all profile configurations from a directory
// It expects files named {profile-name}.yaml
func LoadProfilesFromDir(dir string) ([]*config.ProfileConfig, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []*config.ProfileConfig{}, nil
		}
		return nil, fmt.Errorf("reading profiles directory: %w", err)
	}

	var profiles []*config.ProfileConfig
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".yaml") && !strings.HasSuffix(entry.Name(), ".yml") {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		profile, err := ParseProfileFile(path)
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %w", entry.Name(), err)
		}
		profiles = append(profiles, profile)
	}

	return profiles, nil
}

// WriteProfileFile writes a profile configuration to a YAML file
func WriteProfileFile(profile *config.ProfileConfig, path string) error {
	data, err := yaml.Marshal(profile)
	if err != nil {
		return fmt.Errorf("marshaling profile: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("writing profile file: %w", err)
	}

	return nil
}
