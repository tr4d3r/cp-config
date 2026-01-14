package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractFrontmatter(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantFM      string
		wantContent string
		wantErr     bool
	}{
		{
			name: "valid frontmatter",
			input: `---
name: test
description: A test skill
---

# Test Skill

This is the content.`,
			wantFM:      "name: test\ndescription: A test skill",
			wantContent: "\n# Test Skill\n\nThis is the content.",
			wantErr:     false,
		},
		{
			name: "frontmatter only",
			input: `---
name: test
---`,
			wantFM:      "name: test",
			wantContent: "",
			wantErr:     false,
		},
		{
			name:    "no frontmatter",
			input:   "# Just content\n\nNo frontmatter here.",
			wantErr: true,
		},
		{
			name:    "unclosed frontmatter",
			input:   "---\nname: test\nNo closing delimiter",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractFrontmatter(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractFrontmatter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if got.Frontmatter != tt.wantFM {
				t.Errorf("Frontmatter = %q, want %q", got.Frontmatter, tt.wantFM)
			}
			if got.Content != tt.wantContent {
				t.Errorf("Content = %q, want %q", got.Content, tt.wantContent)
			}
		})
	}
}

func TestParseSkillData(t *testing.T) {
	data := `---
name: debugging
description: Systematic debugging approach
---

# Debugging Skill

Debug all the things!`

	skill, err := ParseSkillData(data, "skills/debugging.md")
	if err != nil {
		t.Fatalf("ParseSkillData() error = %v", err)
	}

	if skill.Name != "debugging" {
		t.Errorf("Name = %q, want %q", skill.Name, "debugging")
	}
	if skill.Description != "Systematic debugging approach" {
		t.Errorf("Description = %q, want %q", skill.Description, "Systematic debugging approach")
	}
	if skill.SourcePath != "skills/debugging.md" {
		t.Errorf("SourcePath = %q, want %q", skill.SourcePath, "skills/debugging.md")
	}
}

func TestParseAgentData(t *testing.T) {
	data := `---
name: code-reviewer
description: Reviews code for quality
---

# Code Reviewer Agent

I review code.`

	agent, err := ParseAgentData(data, "agents/code-reviewer.md")
	if err != nil {
		t.Fatalf("ParseAgentData() error = %v", err)
	}

	if agent.Name != "code-reviewer" {
		t.Errorf("Name = %q, want %q", agent.Name, "code-reviewer")
	}
}

func TestParseProfileData(t *testing.T) {
	data := `name: work
description: Work profile
agents:
  include:
    - code-reviewer
  exclude: []
skills:
  include: []
  exclude:
    - gaming`

	profile, err := ParseProfileData(data, "profiles/work.yaml")
	if err != nil {
		t.Fatalf("ParseProfileData() error = %v", err)
	}

	if profile.Name != "work" {
		t.Errorf("Name = %q, want %q", profile.Name, "work")
	}
	if len(profile.Agents.Include) != 1 || profile.Agents.Include[0] != "code-reviewer" {
		t.Errorf("Agents.Include = %v, want [code-reviewer]", profile.Agents.Include)
	}
	if len(profile.Skills.Exclude) != 1 || profile.Skills.Exclude[0] != "gaming" {
		t.Errorf("Skills.Exclude = %v, want [gaming]", profile.Skills.Exclude)
	}
}

func TestLoadSkillsFromDir(t *testing.T) {
	// Create temp directory with test files
	dir := t.TempDir()

	skill1 := `---
name: skill-one
description: First skill
---

Content one.`

	skill2 := `---
name: skill-two
description: Second skill
---

Content two.`

	if err := os.WriteFile(filepath.Join(dir, "skill-one.md"), []byte(skill1), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "skill-two.md"), []byte(skill2), 0644); err != nil {
		t.Fatal(err)
	}
	// Write a non-md file that should be ignored
	if err := os.WriteFile(filepath.Join(dir, "readme.txt"), []byte("ignore me"), 0644); err != nil {
		t.Fatal(err)
	}

	skills, err := LoadSkillsFromDir(dir)
	if err != nil {
		t.Fatalf("LoadSkillsFromDir() error = %v", err)
	}

	if len(skills) != 2 {
		t.Errorf("got %d skills, want 2", len(skills))
	}
}

func TestLoadSkillsFromDir_NonExistent(t *testing.T) {
	skills, err := LoadSkillsFromDir("/nonexistent/path")
	if err != nil {
		t.Fatalf("LoadSkillsFromDir() should not error for non-existent dir, got %v", err)
	}
	if len(skills) != 0 {
		t.Errorf("got %d skills, want 0", len(skills))
	}
}

func TestRepository_GetMethods(t *testing.T) {
	// Create temp directory structure
	root := t.TempDir()
	agentsDir := filepath.Join(root, "agents")
	skillsDir := filepath.Join(root, "skills")
	profilesDir := filepath.Join(root, "profiles")

	os.MkdirAll(agentsDir, 0755)
	os.MkdirAll(skillsDir, 0755)
	os.MkdirAll(profilesDir, 0755)

	// Create test files
	agent := `---
name: test-agent
description: Test agent
---
Agent content.`

	skill := `---
name: test-skill
description: Test skill
---
Skill content.`

	profile := `name: test-profile
description: Test profile
agents:
  include: []
  exclude: []
skills:
  include: []
  exclude: []`

	os.WriteFile(filepath.Join(agentsDir, "test-agent.md"), []byte(agent), 0644)
	os.WriteFile(filepath.Join(skillsDir, "test-skill.md"), []byte(skill), 0644)
	os.WriteFile(filepath.Join(profilesDir, "test-profile.yaml"), []byte(profile), 0644)

	// Load repository
	agents, _ := LoadAgentsFromDir(agentsDir)
	skills, _ := LoadSkillsFromDir(skillsDir)
	profiles, _ := LoadProfilesFromDir(profilesDir)

	repo := &Repository{
		Root:     root,
		Agents:   agents,
		Skills:   skills,
		Profiles: profiles,
	}

	// Test Get methods
	if a := repo.GetAgent("test-agent"); a == nil {
		t.Error("GetAgent() returned nil for existing agent")
	}
	if a := repo.GetAgent("nonexistent"); a != nil {
		t.Error("GetAgent() should return nil for nonexistent agent")
	}

	if s := repo.GetSkill("test-skill"); s == nil {
		t.Error("GetSkill() returned nil for existing skill")
	}

	if p := repo.GetProfile("test-profile"); p == nil {
		t.Error("GetProfile() returned nil for existing profile")
	}
}
