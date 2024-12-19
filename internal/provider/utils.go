package provider

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var indentRegex = regexp.MustCompile(`^[ \t]+`)

// heredoc removes exactly one leading and trailing newline and dedents the string based on
// the indentation of the first line.
func heredoc(s string) string {
	// Trim exactly one leading newline if it exists
	if strings.HasPrefix(s, "\n") {
		s = s[1:]
	}

	// Find the last newline and trim everything after it
	if lastNL := strings.LastIndex(s, "\n"); lastNL != -1 {
		s = s[:lastNL]
	}

	// Split into lines
	lines := strings.Split(s, "\n")
	if len(lines) == 0 {
		return ""
	}

	// Get indentation from first line using regex
	indentation := ""
	if match := indentRegex.FindString(lines[0]); match != "" {
		indentation = match
	}

	// Remove indentation from each line
	for i, line := range lines {
		lines[i] = strings.TrimPrefix(line, indentation)
	}

	return strings.Join(lines, "\n")
}

// expandPath expands the ~ in paths to the user's home directory.
// If expansion fails, it returns the original path.
func expandPath(path string) string {
	if len(path) == 0 || path[0] != '~' {
		return path
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}

	return filepath.Join(home, path[1:])
}
