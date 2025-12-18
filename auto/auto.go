// Package auto provides a side-effect import to automatically load
// environment variables from .env files during application initialization.
//
// Usage:
//
//	import _ "github.com/rickferrdev/dotenv/auto"
package auto

import "github.com/rickferrdev/dotenv"

// The init function runs automatically when the package is imported,
// ensuring that environment variables are available before main() starts.
func init() {
	dotenv.Collect()
}
