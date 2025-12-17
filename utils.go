// quotes processes a raw string value from an environment variable line.
package dotenv

import "strings"

// It performs the following cleanup steps:
// 1. If the value starts with a single (') or double (") quote, it extracts
//    everything until the matching closing quote.
// 2. If no matching quote is found, it strips the leading quote.
// 3. It removes any trailing comments starting with "#" (only for unquoted
//    content or after the closing quote).
// 4. It trims leading and trailing whitespace from the final result.
func quotes(value string) string {
	if len(value) == 0 {
		return ""
	}

	quote := value[0]
	if quote == '"' || quote == '\'' {
		content, _, found := strings.Cut(value[1:], string(quote))
		if found {
			return content
		}

		value = value[1:]
	}
	value, _, _ = strings.Cut(value, "#")
	return strings.TrimSpace(value)
}
