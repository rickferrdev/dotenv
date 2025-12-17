// Package dotenv provides simple utilities to load environment variables from local files.
package dotenv

import (
	"os"
	"strings"
)

// FilenameVariables defines the default files the package searches for.
var FilenameVariables = []string{".env", ".env.local"}

// Collect iterates through the predefined filenames in FilenameVariables,
// parses their content, and sets the resulting key-value pairs as
// environment variables in the current process.
//
// It supports:
//   - Standard KEY=VALUE pairs.
//   - Lines starting with "export ".
//   - Comments starting with "#".
//   - Basic handling of quoted values (via the internal quotes function
func Collect() {
	for _, filename := range FilenameVariables {
		content, err := os.ReadFile(filename)
		if err != nil {
			continue
		}

		if len(content) <= 1 {
			continue
		}

		for _, line := range strings.Split(string(content), "\n") {
			if strings.HasPrefix(line, "export ") {
				line = strings.TrimPrefix(line, "export")
				line = strings.TrimSpace(line)
			}

			if line == "" || strings.HasPrefix("#", line) {
				continue
			}

			key, value, found := strings.Cut(line, "=")
			if !found {
				continue
			}

			value = quotes(value)

			os.Setenv(key, value)
		}
	}
}
