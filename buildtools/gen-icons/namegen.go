package main

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// SnakeToPascal converts a snake_case icon name to PascalCase with "Icon" prefix.
// e.g. "signal_wifi" → "IconSignalWifi"
func SnakeToPascal(name string) string {
	var b strings.Builder
	b.WriteString("Icon")
	for _, seg := range strings.Split(name, "_") {
		if seg == "" {
			continue
		}
		r, size := utf8.DecodeRuneInString(seg)
		b.WriteRune(unicode.ToUpper(r))
		b.WriteString(seg[size:])
	}
	return b.String()
}

// CheckCollisions verifies that no two icon names produce the same Go identifier.
func CheckCollisions(entries []IconEntry) error {
	seen := make(map[string]string) // converted name → original name
	for _, e := range entries {
		converted := SnakeToPascal(e.Name)
		if prev, ok := seen[converted]; ok {
			return fmt.Errorf("constant name collision: %q and %q both produce %q", prev, e.Name, converted)
		}
		seen[converted] = e.Name
	}
	return nil
}
