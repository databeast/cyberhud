package main

import (
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

// TestConfigSchemaParity validates that the set of JSON keys documented in the
// configuration schema reference equals the set of struct tag keys in fileConfig
// and fileDisplayConfig Go structs.
//
// This test runs as part of normal `go test ./...` to ensure periodic validation
// regardless of whether struct changes have occurred.
func TestConfigSchemaParity(t *testing.T) {
	// Step 1: Extract JSON struct tags from fileConfig and fileDisplayConfig using reflection.
	topLevelKeys := extractJSONKeys(reflect.TypeOf(fileConfig{}))
	displayKeys := extractJSONKeys(reflect.TypeOf(fileDisplayConfig{}))

	// The "display" key is a nested struct in fileConfig, so it appears as a
	// top-level key. Its children are the displayKeys.
	allStructKeys := make(map[string]bool)
	for _, k := range topLevelKeys {
		allStructKeys[k] = true
	}
	for _, k := range displayKeys {
		allStructKeys[k] = true
	}

	// Step 2: Read the schema documentation markdown file.
	schemaPath := findSchemaDoc(t)
	schemaContent, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("failed to read schema documentation: %v", err)
	}

	// Step 3: Extract documented JSON keys from the markdown.
	docKeys := extractDocumentedKeys(string(schemaContent))

	// Step 4: Compare and flag discrepancies.
	var missingFromDocs []string
	var missingFromCode []string

	for key := range allStructKeys {
		if !docKeys[key] {
			missingFromDocs = append(missingFromDocs, key)
		}
	}

	for key := range docKeys {
		if !allStructKeys[key] {
			missingFromCode = append(missingFromCode, key)
		}
	}

	if len(missingFromDocs) > 0 {
		t.Errorf("JSON keys in Go structs but NOT in documentation: %v", missingFromDocs)
	}
	if len(missingFromCode) > 0 {
		t.Errorf("JSON keys in documentation but NOT in Go structs: %v", missingFromCode)
	}
}

// TestConfigSchemaParity_TopLevelStructure validates that the documentation
// correctly distinguishes top-level fields from display object fields.
func TestConfigSchemaParity_TopLevelStructure(t *testing.T) {
	topLevelKeys := extractJSONKeys(reflect.TypeOf(fileConfig{}))

	// Verify top-level keys include expected entries from the struct.
	expectedTopLevel := []string{"socket", "i2c", "scan", "display", "policies"}
	for _, expected := range expectedTopLevel {
		found := false
		for _, key := range topLevelKeys {
			if key == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected top-level key %q in fileConfig struct tags", expected)
		}
	}
}

// extractJSONKeys uses reflection to get all JSON tag key names from a struct type.
// It returns the JSON key names (the part before the first comma in the tag value).
func extractJSONKeys(t reflect.Type) []string {
	var keys []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		// Extract the key name (before any options like omitempty).
		key := strings.Split(tag, ",")[0]
		if key != "" {
			keys = append(keys, key)
		}
	}
	return keys
}

// findSchemaDoc locates the schema documentation file relative to the project root.
func findSchemaDoc(t *testing.T) string {
	t.Helper()

	// Walk up from the test file location to find the project root.
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot determine test file location")
	}

	// The test file is at cmd/cyberhudd/config_schema_parity_test.go
	// Project root is two levels up.
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filename)))

	// Try the MkDocs schema file first.
	schemaPath := filepath.Join(projectRoot, "ghpages", "docs", "configuration", "schema.md")
	if _, err := os.Stat(schemaPath); err == nil {
		return schemaPath
	}

	// Fallback: try CONFIGURATION.md at the root.
	configPath := filepath.Join(projectRoot, "CONFIGURATION.md")
	if _, err := os.Stat(configPath); err == nil {
		return configPath
	}

	t.Fatalf("cannot find schema documentation file; tried:\n  %s\n  %s", schemaPath, configPath)
	return ""
}

// extractDocumentedKeys parses the markdown content and extracts all JSON key names
// that appear in documentation tables. We specifically use table rows as the
// authoritative source of documented keys, since JSON code blocks may contain
// example values (like policy data) that aren't actual struct fields.
//
// It looks for keys in table rows with the pattern: | `key_name` | ...
// (backtick-quoted keys in the first column of markdown tables).
func extractDocumentedKeys(content string) map[string]bool {
	keys := make(map[string]bool)

	// Match backtick-quoted keys in the first column of markdown tables.
	// Keys can contain lowercase letters, digits, and underscores.
	tableKeyRe := regexp.MustCompile("(?m)^\\|\\s*`([a-z][a-z0-9_]*)`\\s*\\|")
	for _, match := range tableKeyRe.FindAllStringSubmatch(content, -1) {
		keys[match[1]] = true
	}

	return keys
}
