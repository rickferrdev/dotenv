// Package dotenv provides simple utilities to load environment variables from local files.
package dotenv

import (
	"errors"
	"fmt"
	"os"
	"reflect"
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
//   - Basic handling of quoted values (via the internal quotes function).
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

			if line == "" || strings.HasPrefix(line, "#") {
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

// Unmarshal parses environment variables into the provided struct.
// The struct must have 'env' tags defining which variables to map.
func Unmarshal(dest interface{}) error {
	rv := reflect.ValueOf(dest)

	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("dest must be a non-nil pointer")
	}

	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return errors.New("dest must be a pointer to a struct")
	}

	t := rv.Type()

	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		fieldType := t.Field(i)

		if !field.CanSet() {
			continue
		}

		key := fieldType.Tag.Get("env")
		if key == "" {
			continue
		}

		value := os.Getenv(key)
		if value == "" {
			continue
		}

		if err := setField(field, value); err != nil {
			return fmt.Errorf("error setting field %s: %w", fieldType.Name, err)
		}
	}

	return nil
}

// Marshal converts a struct into a .env formatted byte slice.
// It uses 'env' tags to define the keys.
func Marshal(dest interface{}) ([]byte, error) {
	rv := reflect.ValueOf(dest)

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return nil, errors.New("dest must be a struct or a pointer to a struct")
	}

	var builder strings.Builder
	t := rv.Type()

	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		fieldType := t.Field(i)

		key := fieldType.Tag.Get("env")
		if key == "" {
			continue
		}

		value := fmt.Sprintf("%v", field.Interface())

		if strings.Contains(value, " ") {
			value = fmt.Sprintf(`"%s"`, value)
		}

		builder.WriteString(fmt.Sprintf("%s=%s\n", key, value))
	}

	return []byte(builder.String()), nil
}
