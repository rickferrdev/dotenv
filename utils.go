// quotes processes a raw string value from an environment variable line.
package dotenv

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// It performs the following cleanup steps:
//  1. If the value starts with a single (') or double (") quote, it extracts
//     everything until the matching closing quote.
//  2. If no matching quote is found, it strips the leading quote.
//  3. It removes any trailing comments starting with "#" (only for unquoted
//     content or after the closing quote).
//  4. It trims leading and trailing whitespace from the final result.
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

// setField helps convert string values to basic Go types supported by the struct fields.
func setField(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(i)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(f)
	default:
		return fmt.Errorf("unsupported type: %s", field.Kind())
	}
	return nil
}
