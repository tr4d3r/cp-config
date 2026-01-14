package github

import (
	"bytes"
	"os"

	"github.com/tr4d3r/cp-config/internal/config"
)

// DiffStatus represents the status of a configuration in a diff
type DiffStatus string

const (
	// DiffStatusAdded means the config exists in repo but not in target
	DiffStatusAdded DiffStatus = "added"

	// DiffStatusRemoved means the config exists in target but not in repo
	DiffStatusRemoved DiffStatus = "removed"

	// DiffStatusModified means the config exists in both but differs
	DiffStatusModified DiffStatus = "modified"

	// DiffStatusUnchanged means the config is identical in both
	DiffStatusUnchanged DiffStatus = "unchanged"
)

// DiffEntry represents a single configuration in a diff result
type DiffEntry struct {
	// Name is the configuration name
	Name string

	// Type is the configuration type (agent or skill)
	Type config.ConfigType

	// Status indicates the diff status
	Status DiffStatus

	// RepoPath is the path in the repository (if exists)
	RepoPath string

	// TargetPath is the path in the target directory (if exists)
	TargetPath string
}

// DiffResult contains the complete diff between repo and target
type DiffResult struct {
	// Agents contains diff entries for agents
	Agents []DiffEntry

	// Skills contains diff entries for skills
	Skills []DiffEntry
}

// HasChanges returns true if there are any changes
func (d *DiffResult) HasChanges() bool {
	for _, e := range d.Agents {
		if e.Status != DiffStatusUnchanged {
			return true
		}
	}
	for _, e := range d.Skills {
		if e.Status != DiffStatusUnchanged {
			return true
		}
	}
	return false
}

// Added returns all entries that need to be added to target
func (d *DiffResult) Added() []DiffEntry {
	var entries []DiffEntry
	for _, e := range d.Agents {
		if e.Status == DiffStatusAdded {
			entries = append(entries, e)
		}
	}
	for _, e := range d.Skills {
		if e.Status == DiffStatusAdded {
			entries = append(entries, e)
		}
	}
	return entries
}

// Removed returns all entries that exist in target but not in repo
func (d *DiffResult) Removed() []DiffEntry {
	var entries []DiffEntry
	for _, e := range d.Agents {
		if e.Status == DiffStatusRemoved {
			entries = append(entries, e)
		}
	}
	for _, e := range d.Skills {
		if e.Status == DiffStatusRemoved {
			entries = append(entries, e)
		}
	}
	return entries
}

// Modified returns all entries that differ between repo and target
func (d *DiffResult) Modified() []DiffEntry {
	var entries []DiffEntry
	for _, e := range d.Agents {
		if e.Status == DiffStatusModified {
			entries = append(entries, e)
		}
	}
	for _, e := range d.Skills {
		if e.Status == DiffStatusModified {
			entries = append(entries, e)
		}
	}
	return entries
}

// DiffAgents compares repo agents with installed agents
func (t *Target) DiffAgents(repoAgents []*config.AgentConfig) ([]DiffEntry, error) {
	installed, err := t.ListInstalledAgents()
	if err != nil {
		return nil, err
	}

	// Build map of installed agents
	installedMap := make(map[string]InstalledConfig)
	for _, cfg := range installed {
		installedMap[cfg.Name] = cfg
	}

	// Build map of repo agents
	repoMap := make(map[string]*config.AgentConfig)
	for _, agent := range repoAgents {
		repoMap[agent.Name] = agent
	}

	var entries []DiffEntry

	// Check repo agents against installed
	for _, agent := range repoAgents {
		entry := DiffEntry{
			Name:     agent.Name,
			Type:     config.ConfigTypeAgent,
			RepoPath: agent.SourcePath,
		}

		if inst, exists := installedMap[agent.Name]; exists {
			entry.TargetPath = inst.Path
			// Compare contents
			if t.filesMatch(agent.SourcePath, inst.Path) {
				entry.Status = DiffStatusUnchanged
			} else {
				entry.Status = DiffStatusModified
			}
		} else {
			entry.Status = DiffStatusAdded
		}

		entries = append(entries, entry)
	}

	// Check for installed agents not in repo
	for _, inst := range installed {
		if _, exists := repoMap[inst.Name]; !exists {
			entries = append(entries, DiffEntry{
				Name:       inst.Name,
				Type:       config.ConfigTypeAgent,
				Status:     DiffStatusRemoved,
				TargetPath: inst.Path,
			})
		}
	}

	return entries, nil
}

// DiffSkills compares repo skills with installed skills
func (t *Target) DiffSkills(repoSkills []*config.SkillConfig) ([]DiffEntry, error) {
	installed, err := t.ListInstalledSkills()
	if err != nil {
		return nil, err
	}

	// Build map of installed skills
	installedMap := make(map[string]InstalledConfig)
	for _, cfg := range installed {
		installedMap[cfg.Name] = cfg
	}

	// Build map of repo skills
	repoMap := make(map[string]*config.SkillConfig)
	for _, skill := range repoSkills {
		repoMap[skill.Name] = skill
	}

	var entries []DiffEntry

	// Check repo skills against installed
	for _, skill := range repoSkills {
		entry := DiffEntry{
			Name:     skill.Name,
			Type:     config.ConfigTypeSkill,
			RepoPath: skill.SourcePath,
		}

		if inst, exists := installedMap[skill.Name]; exists {
			entry.TargetPath = inst.Path
			// Compare contents
			if t.filesMatch(skill.SourcePath, inst.Path) {
				entry.Status = DiffStatusUnchanged
			} else {
				entry.Status = DiffStatusModified
			}
		} else {
			entry.Status = DiffStatusAdded
		}

		entries = append(entries, entry)
	}

	// Check for installed skills not in repo
	for _, inst := range installed {
		if _, exists := repoMap[inst.Name]; !exists {
			entries = append(entries, DiffEntry{
				Name:       inst.Name,
				Type:       config.ConfigTypeSkill,
				Status:     DiffStatusRemoved,
				TargetPath: inst.Path,
			})
		}
	}

	return entries, nil
}

// filesMatch compares two files for equality
func (t *Target) filesMatch(path1, path2 string) bool {
	data1, err := os.ReadFile(path1)
	if err != nil {
		return false
	}

	data2, err := os.ReadFile(path2)
	if err != nil {
		return false
	}

	return bytes.Equal(data1, data2)
}
