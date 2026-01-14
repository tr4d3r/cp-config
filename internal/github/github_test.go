package github

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tr4d3r/cp-config/internal/config"
)

func TestNewTargetAt(t *testing.T) {
	target := NewTargetAt("/test/path")
	if target.Root != "/test/path" {
		t.Errorf("Root = %q, want %q", target.Root, "/test/path")
	}
	if target.AgentsDir != "/test/path/agents" {
		t.Errorf("AgentsDir = %q, want %q", target.AgentsDir, "/test/path/agents")
	}
	if target.SkillsDir != "/test/path/skills" {
		t.Errorf("SkillsDir = %q, want %q", target.SkillsDir, "/test/path/skills")
	}
}

func TestTarget_EnsureDirectories(t *testing.T) {
	dir := t.TempDir()
	target := NewTargetAt(dir)

	if err := target.EnsureDirectories(); err != nil {
		t.Fatalf("EnsureDirectories() error = %v", err)
	}

	// Check directories exist
	if _, err := os.Stat(target.AgentsDir); os.IsNotExist(err) {
		t.Error("AgentsDir was not created")
	}
	if _, err := os.Stat(target.SkillsDir); os.IsNotExist(err) {
		t.Error("SkillsDir was not created")
	}
}

func TestTarget_DeployAgent(t *testing.T) {
	dir := t.TempDir()
	target := NewTargetAt(dir)

	agent := &config.AgentConfig{
		Name:        "test-agent",
		Description: "A test agent",
		Content:     "\n# Test Agent\n\nThis is content.",
	}

	if err := target.DeployAgent(agent); err != nil {
		t.Fatalf("DeployAgent() error = %v", err)
	}

	// Verify file was created
	path := target.AgentPath("test-agent")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading deployed file: %v", err)
	}

	content := string(data)
	if !contains(content, "name: test-agent") {
		t.Error("deployed file missing name in frontmatter")
	}
	if !contains(content, "description: A test agent") {
		t.Error("deployed file missing description in frontmatter")
	}
	if !contains(content, "# Test Agent") {
		t.Error("deployed file missing content")
	}
}

func TestTarget_DeploySkill(t *testing.T) {
	dir := t.TempDir()
	target := NewTargetAt(dir)

	skill := &config.SkillConfig{
		Name:        "debugging",
		Description: "Debug skill",
		Content:     "\n# Debugging\n\nDebug content.",
	}

	if err := target.DeploySkill(skill); err != nil {
		t.Fatalf("DeploySkill() error = %v", err)
	}

	// Verify file was created
	if !target.HasSkill("debugging") {
		t.Error("skill was not deployed")
	}
}

func TestTarget_RemoveAgent(t *testing.T) {
	dir := t.TempDir()
	target := NewTargetAt(dir)

	// Deploy then remove
	agent := &config.AgentConfig{
		Name:        "to-remove",
		Description: "Will be removed",
		Content:     "Content",
	}
	target.DeployAgent(agent)

	if !target.HasAgent("to-remove") {
		t.Fatal("agent was not deployed")
	}

	if err := target.RemoveAgent("to-remove"); err != nil {
		t.Fatalf("RemoveAgent() error = %v", err)
	}

	if target.HasAgent("to-remove") {
		t.Error("agent was not removed")
	}
}

func TestTarget_ListInstalledAgents(t *testing.T) {
	dir := t.TempDir()
	target := NewTargetAt(dir)

	// Deploy some agents
	agents := []*config.AgentConfig{
		{Name: "agent-1", Description: "First", Content: "Content 1"},
		{Name: "agent-2", Description: "Second", Content: "Content 2"},
	}
	target.DeployAgents(agents)

	installed, err := target.ListInstalledAgents()
	if err != nil {
		t.Fatalf("ListInstalledAgents() error = %v", err)
	}

	if len(installed) != 2 {
		t.Errorf("got %d agents, want 2", len(installed))
	}
}

func TestTarget_ListInstalledAgents_Empty(t *testing.T) {
	dir := t.TempDir()
	target := NewTargetAt(dir)

	// Don't create directories
	installed, err := target.ListInstalledAgents()
	if err != nil {
		t.Fatalf("ListInstalledAgents() error = %v", err)
	}

	if len(installed) != 0 {
		t.Errorf("got %d agents, want 0", len(installed))
	}
}

func TestTarget_DiffAgents(t *testing.T) {
	dir := t.TempDir()
	target := NewTargetAt(dir)

	// Create a repo agent file
	repoDir := t.TempDir()
	repoAgentPath := filepath.Join(repoDir, "existing.md")
	os.WriteFile(repoAgentPath, []byte(`---
name: existing
description: Existing agent
---

Content`), 0644)

	// Deploy one agent that's also in repo (same content)
	installedPath := target.AgentPath("existing")
	target.EnsureDirectories()
	os.WriteFile(installedPath, []byte(`---
name: existing
description: Existing agent
---

Content`), 0644)

	// Deploy one agent that's not in repo
	target.DeployAgent(&config.AgentConfig{
		Name:        "orphan",
		Description: "Not in repo",
		Content:     "Orphan content",
	})

	// Repo agents (one matches, one is new)
	repoAgents := []*config.AgentConfig{
		{Name: "existing", Description: "Existing agent", Content: "\nContent", SourcePath: repoAgentPath},
		{Name: "new-agent", Description: "New agent", Content: "\nNew content", SourcePath: filepath.Join(repoDir, "new.md")},
	}

	entries, err := target.DiffAgents(repoAgents)
	if err != nil {
		t.Fatalf("DiffAgents() error = %v", err)
	}

	// Should have 3 entries: existing (unchanged), new-agent (added), orphan (removed)
	if len(entries) != 3 {
		t.Errorf("got %d entries, want 3", len(entries))
	}

	statusMap := make(map[string]DiffStatus)
	for _, e := range entries {
		statusMap[e.Name] = e.Status
	}

	if statusMap["existing"] != DiffStatusUnchanged {
		t.Errorf("existing status = %v, want %v", statusMap["existing"], DiffStatusUnchanged)
	}
	if statusMap["new-agent"] != DiffStatusAdded {
		t.Errorf("new-agent status = %v, want %v", statusMap["new-agent"], DiffStatusAdded)
	}
	if statusMap["orphan"] != DiffStatusRemoved {
		t.Errorf("orphan status = %v, want %v", statusMap["orphan"], DiffStatusRemoved)
	}
}

func TestDiffResult_HasChanges(t *testing.T) {
	tests := []struct {
		name   string
		result DiffResult
		want   bool
	}{
		{
			name: "no changes",
			result: DiffResult{
				Agents: []DiffEntry{{Status: DiffStatusUnchanged}},
				Skills: []DiffEntry{{Status: DiffStatusUnchanged}},
			},
			want: false,
		},
		{
			name: "has added",
			result: DiffResult{
				Agents: []DiffEntry{{Status: DiffStatusAdded}},
			},
			want: true,
		},
		{
			name: "has modified skill",
			result: DiffResult{
				Skills: []DiffEntry{{Status: DiffStatusModified}},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.result.HasChanges(); got != tt.want {
				t.Errorf("HasChanges() = %v, want %v", got, tt.want)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
