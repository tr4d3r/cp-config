package repo

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/tr4d3r/cp-config/internal/config"
	"github.com/tr4d3r/cp-config/internal/parser"
)

var (
	// ErrNotInRepo indicates the current directory is not in a cp-config repository
	ErrNotInRepo = errors.New("not in a cp-config repository (no agents/ or skills/ directory found)")
)

// FindRoot finds the repository root by looking for agents/ or skills/ directories
// starting from the current directory and walking up
func FindRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return FindRootFrom(cwd)
}

// FindRootFrom finds the repository root starting from the given directory
func FindRootFrom(start string) (string, error) {
	dir := start
	for {
		// Check if this looks like a cp-config repo
		if isRepoRoot(dir) {
			return dir, nil
		}

		// Move up one directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root
			break
		}
		dir = parent
	}

	return "", ErrNotInRepo
}

// isRepoRoot checks if a directory looks like a cp-config repository root
func isRepoRoot(dir string) bool {
	// Check for agents/ directory
	agentsDir := filepath.Join(dir, "agents")
	if info, err := os.Stat(agentsDir); err == nil && info.IsDir() {
		return true
	}

	// Check for skills/ directory
	skillsDir := filepath.Join(dir, "skills")
	if info, err := os.Stat(skillsDir); err == nil && info.IsDir() {
		return true
	}

	// Check for go.mod with cp-config module
	goMod := filepath.Join(dir, "go.mod")
	if _, err := os.Stat(goMod); err == nil {
		return true
	}

	return false
}

// Load loads the repository from the given root directory
func Load(root string) (*parser.Repository, error) {
	cfg := &config.AppConfig{
		RepoRoot: root,
	}

	return parser.LoadRepository(cfg)
}

// LoadFromCwd loads the repository from the current working directory
func LoadFromCwd() (*parser.Repository, error) {
	root, err := FindRoot()
	if err != nil {
		return nil, err
	}
	return Load(root)
}
