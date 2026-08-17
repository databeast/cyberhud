package cmdutil

import (
	"fmt"
	"sort"
	"strings"
)

// KeyValidator validates a string value for a command option key.
// It returns an error description if the value is invalid, or "" if valid.
type KeyValidator func(value string) string

// KeyDef defines a recognized key for a command handler, including its validator
// and a getter for the current value.
type KeyDef struct {
	// Name is the key name (e.g., "style", "font").
	Name string
	// Validate checks a proposed value and returns "" if valid,
	// or a human-readable reason if invalid.
	Validate KeyValidator
}

// OnApplied is an optional package-level callback invoked after any CmdHandler
// successfully applies key=value pairs. It receives the handler's Verb (mode ID).
// The catalog package uses this to trigger PolicyStore updates centrally.
var OnApplied func(verb string)

// CmdHandler processes key=value command arguments following the uniform pattern:
//   - Zero args → return current state as "OK <verb> key1=value1 key2=value2 ..."
//   - Key=value args → validate all, apply atomically, return updated state or "ERR ..."
type CmdHandler struct {
	// Verb is the command verb (e.g., "clock", "dashboard").
	Verb string
	// Keys defines recognized option keys in display order.
	Keys []KeyDef
	// Get returns the current value for a given key name.
	// It is called to build the query response and the post-apply response.
	Get func(key string) string
	// Apply sets the value for a given key. It is called only after all
	// key=value pairs have passed validation.
	Apply func(key, value string)
	// PostApply is an optional hook called after all key=value pairs have been
	// applied successfully. It receives the list of applied key names and returns
	// additional response lines to append (e.g., fitness notes). May be nil.
	PostApply func(appliedKeys []string) []string
}

// Handle processes a set of arguments and returns the command response string.
// If args is empty, it returns the current state query. Otherwise it parses
// key=value pairs, validates them, and applies atomically.
func (h *CmdHandler) Handle(args []string) string {
	if len(args) == 0 {
		return h.queryResponse()
	}

	// Parse all key=value pairs first.
	type kv struct {
		key string
		val string
	}
	pairs := make([]kv, 0, len(args))
	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			// Treat bare tokens as keys with empty values — still a parse error.
			return fmt.Sprintf("ERR unknown key %s", arg)
		}
		pairs = append(pairs, kv{key: strings.TrimSpace(parts[0]), val: strings.TrimSpace(parts[1])})
	}

	// Validate all pairs before applying any.
	for _, p := range pairs {
		def, ok := h.findKey(p.key)
		if !ok {
			return fmt.Sprintf("ERR unknown key %s", p.key)
		}
		if def.Validate != nil {
			if reason := def.Validate(p.val); reason != "" {
				return fmt.Sprintf("ERR %s: %s", p.key, reason)
			}
		}
	}

	// All valid — apply atomically.
	for _, p := range pairs {
		h.Apply(p.key, p.val)
	}

	resp := h.queryResponse()

	// Invoke PostApply hook if present.
	if h.PostApply != nil {
		keys := make([]string, len(pairs))
		for i, p := range pairs {
			keys[i] = p.key
		}
		if extra := h.PostApply(keys); len(extra) > 0 {
			resp += "\n" + strings.Join(extra, "\n")
		}
	}

	// Notify the centralized policy change listener.
	if OnApplied != nil {
		OnApplied(h.Verb)
	}

	return resp
}

// queryResponse builds the "OK <verb> key1=value1 key2=value2 ..." string.
func (h *CmdHandler) queryResponse() string {
	var b strings.Builder
	b.WriteString("OK ")
	b.WriteString(h.Verb)
	for _, k := range h.Keys {
		b.WriteByte(' ')
		b.WriteString(k.Name)
		b.WriteByte('=')
		b.WriteString(h.Get(k.Name))
	}
	return b.String()
}

// findKey looks up a KeyDef by name (case-insensitive).
func (h *CmdHandler) findKey(name string) (KeyDef, bool) {
	lower := strings.ToLower(name)
	for _, k := range h.Keys {
		if strings.ToLower(k.Name) == lower {
			return k, true
		}
	}
	return KeyDef{}, false
}

// AllowedValidator returns a KeyValidator that accepts only values present in the
// allowed list (case-insensitive comparison). The error message lists allowed values.
func AllowedValidator(allowed []string) KeyValidator {
	return func(value string) string {
		lower := strings.ToLower(strings.TrimSpace(value))
		for _, a := range allowed {
			if strings.ToLower(a) == lower {
				return ""
			}
		}
		sorted := make([]string, len(allowed))
		copy(sorted, allowed)
		sort.Strings(sorted)
		return fmt.Sprintf("must be one of [%s]", strings.Join(sorted, ", "))
	}
}

// BoolValidator returns a KeyValidator that accepts boolean-like strings
// (true/false/yes/no/on/off/1/0).
func BoolValidator() KeyValidator {
	return func(value string) string {
		switch strings.ToLower(strings.TrimSpace(value)) {
		case "true", "false", "yes", "no", "on", "off", "1", "0":
			return ""
		default:
			return "must be true or false"
		}
	}
}

// ParseBool parses boolean-like strings. Returns the parsed value and
// whether the input was recognized (true/yes/on/1 → true; false/no/off/0 → false).
func ParseBool(s string) (bool, bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "true", "yes", "on", "1":
		return true, true
	case "false", "no", "off", "0":
		return false, true
	default:
		return false, false
	}
}

// IntValidator returns a KeyValidator that accepts integer strings,
// optionally constrained to a minimum value.
func IntValidator(min int) KeyValidator {
	return func(value string) string {
		value = strings.TrimSpace(value)
		n := 0
		negative := false
		if len(value) == 0 {
			return "must be an integer"
		}
		start := 0
		if value[0] == '-' {
			negative = true
			start = 1
		} else if value[0] == '+' {
			start = 1
		}
		if start >= len(value) {
			return "must be an integer"
		}
		for i := start; i < len(value); i++ {
			if value[i] < '0' || value[i] > '9' {
				return "must be an integer"
			}
			n = n*10 + int(value[i]-'0')
		}
		if negative {
			n = -n
		}
		if n < min {
			return fmt.Sprintf("must be >= %d", min)
		}
		return ""
	}
}
