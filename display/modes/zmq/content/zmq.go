// Package zmq implements the ZMQ display mode for cyberHUD.
// It connects to a ZeroMQ socket (SUB or PULL) and renders incoming
// messages on TFT panels, supporting optional JSON field filtering.
package content

import (
	"fmt"
	"strings"
	"sync"

	"github.com/databeast/cyberhud/display/catalog"
)

// Allowed socket type values.
var allowedSocketTypes = []string{"sub", "pull"}

// Policy captures all runtime-configurable parameters for the ZMQ display mode.
type Policy struct {
	Endpoint   string // ZMQ connection string, e.g. "tcp://localhost:5556" (max 256 chars)
	SocketType string // "sub" or "pull"
	Topic      string // SUB topic filter (ignored for PULL)
	MaxLines   int    // Ring buffer capacity [1, 1000]
	Style      string // Style registry name
	Font       string // Font ID or "auto"
	JSONFields string // Comma-separated field names for JSON extraction
}

// DefaultPolicy returns baseline ZMQ mode behavior with sensible defaults.
func DefaultPolicy() Policy {
	return Policy{
		Endpoint:   "",
		SocketType: "sub",
		Topic:      "",
		MaxLines:   24,
		Style:      defaultStyleName,
		Font:       "auto",
		JSONFields: "",
	}
}

func (p Policy) Options() []catalog.OptionDefinition {
	return []catalog.OptionDefinition{
		{Key: "endpoint", Type: "string", Summary: "ZMQ connection endpoint (e.g. tcp://localhost:5556).", Default: ""},
		{Key: "socket_type", Type: "string", Summary: "ZMQ socket type.", Default: "sub", Allowed: allowedSocketTypes},
		{Key: "topic", Type: "string", Summary: "SUB topic filter (ignored for PULL).", Default: ""},
		{Key: "max_lines", Type: "int", Summary: "Ring buffer capacity (message lines to retain).", Default: "24"},
		{Key: "json_fields", Type: "string", Summary: "Comma-separated JSON field names to extract and display.", Default: ""},
		{Key: "font", Type: "string", Summary: "Font selection (auto or a registered font ID).", Default: "auto"},
	}
}

func (p Policy) Fingerprint() string {
	return fmt.Sprintf("%s|%s|%s|%d|%s|%s|%s",
		p.Endpoint, p.SocketType, p.Topic, p.MaxLines, p.Style, p.Font, p.JSONFields)
}

func (p Policy) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"endpoint":    p.Endpoint,
		"socket_type": p.SocketType,
		"topic":       p.Topic,
		"max_lines":   p.MaxLines,
		"style":       p.Style,
		"font":        p.Font,
		"json_fields": p.JSONFields,
	}
}

// defaultStyleName is the registry default style name.
const defaultStyleName = "color-medium-240x240"

// policyState holds the mutex-protected active policy.
var policyState = struct {
	sync.RWMutex
	policy Policy
}{
	policy: DefaultPolicy(),
}

// GetPolicy returns the current ZMQ policy (thread-safe read under RWMutex).
func GetPolicy() Policy {
	policyState.RLock()
	defer policyState.RUnlock()
	return policyState.policy
}

// SetPolicy updates the ZMQ policy after normalization (thread-safe write under Mutex).
// If connection-relevant fields (Endpoint, SocketType, Topic) change while the
// receiver is active, the receiver is deactivated and reactivated with the new policy.
func SetPolicy(p Policy) {
	newPolicy := normalizePolicy(p)

	policyState.Lock()
	oldPolicy := policyState.policy
	policyState.policy = newPolicy
	policyState.Unlock()
	msgBuffer.SetMaxLines(newPolicy.MaxLines)

	// Detect connection-relevant changes and trigger reconnection if needed.
	if defaultReceiver.IsActive() && connectionFieldsChanged(oldPolicy, newPolicy) {
		defaultReceiver.Deactivate()
		defaultReceiver.Activate(newPolicy, msgBuffer)
	}
}

// normalizePolicy ensures all policy fields contain valid, canonical values.
// Normalization is idempotent: normalizePolicy(normalizePolicy(p)) == normalizePolicy(p).
func normalizePolicy(p Policy) Policy {
	// Trim all string fields.
	p.Endpoint = strings.TrimSpace(p.Endpoint)
	p.SocketType = strings.TrimSpace(p.SocketType)
	p.Topic = strings.TrimSpace(p.Topic)
	p.Style = strings.TrimSpace(p.Style)
	p.Font = strings.TrimSpace(p.Font)
	p.JSONFields = strings.TrimSpace(p.JSONFields)

	// Truncate Endpoint to 256 characters.
	if len(p.Endpoint) > 256 {
		p.Endpoint = p.Endpoint[:256]
	}

	// Lowercase SocketType, validate against allowed values.
	p.SocketType = strings.ToLower(p.SocketType)
	if !isAllowedSocketType(p.SocketType) {
		p.SocketType = "sub"
	}

	// Clamp MaxLines to [1, 1000].
	if p.MaxLines < 1 {
		p.MaxLines = 1
	}
	if p.MaxLines > 1000 {
		p.MaxLines = 1000
	}

	// Lowercase Style. Registry-aware validation is performed by the parent controller.
	p.Style = strings.ToLower(p.Style)
	if p.Style == "" {
		p.Style = defaultStyleName
	}

	// Fallback Font to "auto" if empty.
	if p.Font == "" {
		p.Font = "auto"
	}

	return p
}

// connectionFieldsChanged reports whether any connection-relevant policy fields
// (Endpoint, SocketType, Topic) differ between old and new. Both policies should
// already be normalized before comparison.
func connectionFieldsChanged(old, new Policy) bool {
	return old.Endpoint != new.Endpoint ||
		old.SocketType != new.SocketType ||
		old.Topic != new.Topic
}

// isAllowedSocketType checks if the value is a valid socket type.
func isAllowedSocketType(value string) bool {
	for _, a := range allowedSocketTypes {
		if value == a {
			return true
		}
	}
	return false
}
