package config

import "testing"

func TestSkillConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		skill   SkillConfig
		wantErr bool
	}{
		{
			name: "valid skill",
			skill: SkillConfig{
				Name:        "debugging",
				Description: "Systematic debugging approach",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			skill: SkillConfig{
				Description: "Some description",
			},
			wantErr: true,
		},
		{
			name: "missing description",
			skill: SkillConfig{
				Name: "debugging",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.skill.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSkillConfig_TargetPath(t *testing.T) {
	s := SkillConfig{Name: "debugging"}
	want := "skills/debugging.md"
	if got := s.TargetPath(); got != want {
		t.Errorf("TargetPath() = %v, want %v", got, want)
	}
}

func TestAgentConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		agent   AgentConfig
		wantErr bool
	}{
		{
			name: "valid agent",
			agent: AgentConfig{
				Name:        "code-reviewer",
				Description: "Reviews code for quality",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			agent: AgentConfig{
				Description: "Some description",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.agent.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProfileConfig_ShouldInclude(t *testing.T) {
	profile := &ProfileConfig{
		Name: "test",
		Agents: FilterConfig{
			Include: []string{"agent-a", "agent-b"},
			Exclude: []string{"agent-b"},
		},
		Skills: FilterConfig{
			Include: []string{}, // empty means include all
			Exclude: []string{"skill-x"},
		},
	}

	tests := []struct {
		name     string
		itemName string
		isAgent  bool
		want     bool
	}{
		{"included agent", "agent-a", true, true},
		{"excluded agent", "agent-b", true, false},
		{"not in include list", "agent-c", true, false},
		{"skill not excluded", "skill-y", false, true},
		{"skill excluded", "skill-x", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got bool
			if tt.isAgent {
				got = profile.ShouldIncludeAgent(tt.itemName)
			} else {
				got = profile.ShouldIncludeSkill(tt.itemName)
			}
			if got != tt.want {
				t.Errorf("ShouldInclude(%s) = %v, want %v", tt.itemName, got, tt.want)
			}
		})
	}
}

func TestDefaultProfile(t *testing.T) {
	p := DefaultProfile()
	if p.Name != "default" {
		t.Errorf("DefaultProfile().Name = %v, want default", p.Name)
	}
	// Default profile should include everything
	if !p.ShouldIncludeAgent("any-agent") {
		t.Error("DefaultProfile should include all agents")
	}
	if !p.ShouldIncludeSkill("any-skill") {
		t.Error("DefaultProfile should include all skills")
	}
}
