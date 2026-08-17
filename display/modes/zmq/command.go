package zmq

import (
	"strconv"
	"strings"
)

// endpointValidator accepts any string (even empty means no connection).
func endpointValidator(value string) string {
	return ""
}

// styleValidator validates style names against the root style registry.
func styleValidator(value string) string {
	v := strings.TrimSpace(value)
	if v == "" {
		return "must be a non-empty style name"
	}
	if zmqRegistry.Lookup(v) == nil {
		return "must be one of: " + strings.Join(registeredStyleNames(), ", ")
	}
	return ""
}

// fontValidator accepts "auto" or any non-empty trimmed string.
func fontValidator(value string) string {
	v := strings.TrimSpace(value)
	if v == "" {
		return "must be \"auto\" or a non-empty font ID"
	}
	return ""
}

// getZmqPolicyValue reads the current value for a given key from the ZMQ policy.
func getZmqPolicyValue(key string) string {
	p := GetPolicy()
	switch key {
	case "endpoint":
		return p.Endpoint
	case "socket_type":
		return p.SocketType
	case "topic":
		return p.Topic
	case "max_lines":
		return strconv.Itoa(p.MaxLines)
	case "json_fields":
		return p.JSONFields
	case "style":
		return p.Style
	case "font":
		return p.Font
	default:
		return ""
	}
}

// applyZmqPolicyValue updates the policy for a single key.
// SetPolicy handles normalization and reconnection detection.
func applyZmqPolicyValue(key, value string) {
	p := GetPolicy()
	switch key {
	case "endpoint":
		p.Endpoint = strings.TrimSpace(value)
	case "socket_type":
		p.SocketType = strings.TrimSpace(value)
	case "topic":
		p.Topic = strings.TrimSpace(value)
	case "max_lines":
		if n, err := strconv.Atoi(strings.TrimSpace(value)); err == nil {
			p.MaxLines = n
		}
	case "json_fields":
		p.JSONFields = strings.TrimSpace(value)
	case "style":
		p.Style = strings.TrimSpace(value)
	case "font":
		p.Font = strings.TrimSpace(value)
	}
	SetPolicy(p)
}
