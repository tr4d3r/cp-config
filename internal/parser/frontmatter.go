package parser

import (
	"errors"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	// ErrNoFrontmatter indicates the file has no YAML frontmatter
	ErrNoFrontmatter = errors.New("no frontmatter found")

	// ErrInvalidFrontmatter indicates malformed frontmatter
	ErrInvalidFrontmatter = errors.New("invalid frontmatter format")
)

const frontmatterDelimiter = "---"

// ParsedMarkdown represents a markdown file with extracted frontmatter
type ParsedMarkdown struct {
	// Frontmatter is the raw YAML frontmatter string
	Frontmatter string

	// Content is the markdown body after the frontmatter
	Content string

	// Raw is the original complete file content
	Raw string
}

// ExtractFrontmatter parses a markdown file and separates frontmatter from content
// The file must start with "---" followed by YAML, then another "---"
func ExtractFrontmatter(data string) (*ParsedMarkdown, error) {
	data = strings.TrimSpace(data)

	// Must start with ---
	if !strings.HasPrefix(data, frontmatterDelimiter) {
		return nil, ErrNoFrontmatter
	}

	// Find the closing ---
	rest := data[len(frontmatterDelimiter):]
	rest = strings.TrimPrefix(rest, "\n")

	endIdx := strings.Index(rest, "\n"+frontmatterDelimiter)
	if endIdx == -1 {
		// Try without newline prefix (for files ending right after frontmatter)
		if strings.HasSuffix(rest, frontmatterDelimiter) {
			endIdx = len(rest) - len(frontmatterDelimiter)
		} else {
			return nil, ErrInvalidFrontmatter
		}
	}

	frontmatter := strings.TrimSpace(rest[:endIdx])
	content := ""

	// Extract content after the closing ---
	afterFrontmatter := rest[endIdx:]
	afterFrontmatter = strings.TrimPrefix(afterFrontmatter, "\n")
	afterFrontmatter = strings.TrimPrefix(afterFrontmatter, frontmatterDelimiter)
	content = strings.TrimPrefix(afterFrontmatter, "\n")

	return &ParsedMarkdown{
		Frontmatter: frontmatter,
		Content:     content,
		Raw:         data,
	}, nil
}

// UnmarshalFrontmatter parses the frontmatter YAML into the provided struct
func UnmarshalFrontmatter(parsed *ParsedMarkdown, v interface{}) error {
	if parsed.Frontmatter == "" {
		return ErrNoFrontmatter
	}
	return yaml.Unmarshal([]byte(parsed.Frontmatter), v)
}
