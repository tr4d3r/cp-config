package profile

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tr4d3r/cp-config/internal/config"
	"github.com/tr4d3r/cp-config/internal/parser"
	"gopkg.in/yaml.v3"
)

var (
	// ErrProfileNotFound indicates the requested profile does not exist
	ErrProfileNotFound = errors.New("profile not found")

	// ErrProfileExists indicates a profile with that name already exists
	ErrProfileExists = errors.New("profile already exists")

	// ErrCannotDeleteDefault indicates an attempt to delete the default profile
	ErrCannotDeleteDefault = errors.New("cannot delete the default profile")
)

// Manager handles profile operations
type Manager struct {
	// ProfilesDir is the directory containing profile YAML files
	ProfilesDir string

	// StateFile is the path to the state file storing the active profile
	StateFile string
}

// State represents the persistent state (active profile)
type State struct {
	ActiveProfile string `yaml:"active_profile"`
}

// NewManager creates a new profile manager
func NewManager(repoRoot string) *Manager {
	return &Manager{
		ProfilesDir: filepath.Join(repoRoot, "profiles"),
		StateFile:   filepath.Join(repoRoot, ".cp-config-state.yaml"),
	}
}

// List returns all available profiles
func (m *Manager) List() ([]*config.ProfileConfig, error) {
	profiles, err := parser.LoadProfilesFromDir(m.ProfilesDir)
	if err != nil {
		return nil, err
	}

	// Always include a default profile if none exists
	hasDefault := false
	for _, p := range profiles {
		if p.Name == "default" {
			hasDefault = true
			break
		}
	}

	if !hasDefault {
		profiles = append([]*config.ProfileConfig{config.DefaultProfile()}, profiles...)
	}

	return profiles, nil
}

// Get returns a profile by name
func (m *Manager) Get(name string) (*config.ProfileConfig, error) {
	if name == "default" {
		// Check if there's a default.yaml file first
		path := filepath.Join(m.ProfilesDir, "default.yaml")
		if _, err := os.Stat(path); err == nil {
			return parser.ParseProfileFile(path)
		}
		// Return the built-in default
		return config.DefaultProfile(), nil
	}

	path := filepath.Join(m.ProfilesDir, name+".yaml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, ErrProfileNotFound
	}

	return parser.ParseProfileFile(path)
}

// Create creates a new profile
func (m *Manager) Create(profile *config.ProfileConfig) error {
	if err := profile.Validate(); err != nil {
		return err
	}

	path := filepath.Join(m.ProfilesDir, profile.Name+".yaml")
	if _, err := os.Stat(path); err == nil {
		return ErrProfileExists
	}

	// Ensure profiles directory exists
	if err := os.MkdirAll(m.ProfilesDir, 0755); err != nil {
		return fmt.Errorf("creating profiles directory: %w", err)
	}

	return parser.WriteProfileFile(profile, path)
}

// Delete removes a profile
func (m *Manager) Delete(name string) error {
	if name == "default" {
		return ErrCannotDeleteDefault
	}

	path := filepath.Join(m.ProfilesDir, name+".yaml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return ErrProfileNotFound
	}

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("removing profile: %w", err)
	}

	// If this was the active profile, switch back to default
	active, _ := m.GetActiveProfile()
	if active == name {
		m.SetActiveProfile("default")
	}

	return nil
}

// GetActiveProfile returns the currently active profile name
func (m *Manager) GetActiveProfile() (string, error) {
	data, err := os.ReadFile(m.StateFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "default", nil
		}
		return "", err
	}

	var state State
	if err := yaml.Unmarshal(data, &state); err != nil {
		return "default", nil
	}

	if state.ActiveProfile == "" {
		return "default", nil
	}

	return state.ActiveProfile, nil
}

// SetActiveProfile sets the active profile
func (m *Manager) SetActiveProfile(name string) error {
	// Verify profile exists (unless it's default)
	if name != "default" {
		if _, err := m.Get(name); err != nil {
			return err
		}
	}

	state := State{ActiveProfile: name}
	data, err := yaml.Marshal(&state)
	if err != nil {
		return err
	}

	return os.WriteFile(m.StateFile, data, 0644)
}

// Exists checks if a profile exists
func (m *Manager) Exists(name string) bool {
	if name == "default" {
		return true
	}
	path := filepath.Join(m.ProfilesDir, name+".yaml")
	_, err := os.Stat(path)
	return err == nil
}
