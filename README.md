# cp-config

A CLI tool for managing GitHub Copilot agent and skill configurations.

cp-config enables source-controlled management of Copilot configurations with deployment to your `~/.github/copilot` directory, supporting sharing and versioning of agents and skills across machines and teams.

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/tr4d3r/cp-config.git
cd cp-config

# Build and install
go install ./cmd/cp-config
```

### Pre-built Binaries

Download the latest release from the [releases page](https://github.com/tr4d3r/cp-config/releases).

## Quick Start

```bash
# List available agents and skills
cp-config list

# Deploy configurations to ~/.github/copilot
cp-config install

# Show differences between repo and installed configs
cp-config diff
```

## Commands

### list

List available agent and skill configurations from the repository.

```bash
cp-config list              # List all agents and skills
cp-config list agents       # List only agents
cp-config list skills       # List only skills
cp-config list -p work      # List using 'work' profile
```

### install

Deploy configurations to `~/.github/copilot`.

```bash
cp-config install                    # Install using default profile
cp-config install --profile work     # Install using 'work' profile
```

### diff

Show differences between repository and installed configurations.

```bash
cp-config diff                    # Compare using default profile
cp-config diff --profile work     # Compare using 'work' profile
```

### profile

Manage configuration profiles for different contexts.

```bash
cp-config profile list            # List available profiles
cp-config profile create <name>   # Create a new profile
cp-config profile switch <name>   # Switch to a different profile
cp-config profile delete <name>   # Delete a profile
```

### version

Print version information.

```bash
cp-config version
```

## Configuration Format

### Agents

Agents are markdown files with YAML frontmatter stored in the `agents/` directory:

```markdown
---
name: code-reviewer
description: Reviews code for quality, security, and best practices
---

# Code Reviewer Agent

Instructions for the agent...
```

### Skills

Skills are markdown files with YAML frontmatter stored in the `skills/` directory:

```markdown
---
name: debugging
description: Systematic debugging approach with root cause documentation
---

# Debugging Skill

Instructions for the skill...
```

### Profiles

Profiles are YAML files stored in the `profiles/` directory:

```yaml
name: work
description: Work-related configurations
agents:
  include:
    - code-reviewer
    - documentation-helper
  exclude: []
skills:
  include:
    - api-design
    - testing
  exclude: []
```

## Repository Structure

```
your-config-repo/
├── agents/                    # Agent configuration files
│   └── {agent-name}.md
├── skills/                    # Skill configuration files
│   └── {skill-name}.md
├── profiles/                  # Profile configurations
│   └── {profile-name}.yaml
└── README.md
```

## Target Directory

Configurations are deployed to:

```
~/.github/
└── copilot/
    ├── agents/
    │   └── {agent-name}.md
    └── skills/
        └── {skill-name}.md
```

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--config` | | Config file path (default: `.cp-config.yaml`) |
| `--profile` | `-p` | Profile to use (default: `default`) |
| `--help` | `-h` | Help for any command |

## License

MIT License - see [LICENSE](LICENSE) for details.
