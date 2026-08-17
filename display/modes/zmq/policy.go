package zmq

import (
	"strings"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/modes/zmq/content"
)

func GetPolicy() content.Policy {
	return content.GetPolicy()
}

func SetPolicy(p content.Policy) {
	content.SetPolicy(normalizePolicy(p))
}

func normalizePolicy(p content.Policy) content.Policy {
	p.Style = strings.ToLower(strings.TrimSpace(p.Style))
	if p.Style != "" {
		if resolved := zmqRegistry.Normalize(p.Style); resolved != "" {
			p.Style = resolved
		} else {
			p.Style = ""
		}
	}
	return p
}

// zmqSnapshotter implements catalog.PolicySnapshotter for the zmq mode.
type zmqSnapshotter struct{}

// SnapshotPolicy returns the current zmq policy as a JSON-serializable map
// using snake_case keys matching the protocol wire format.
func (zmqSnapshotter) SnapshotPolicy() map[string]interface{} {
	return GetPolicy().ToMap()
}

// RestorePolicy applies policy values from a JSON map, running them through
// the same normalization as SetPolicy. Keys use the protocol wire format.
func (zmqSnapshotter) RestorePolicy(data map[string]interface{}) error {
	p := content.DefaultPolicy()

	if v, ok := data["endpoint"]; ok {
		if s, ok := v.(string); ok {
			p.Endpoint = s
		}
	}
	if v, ok := data["socket_type"]; ok {
		if s, ok := v.(string); ok {
			p.SocketType = s
		}
	}
	if v, ok := data["topic"]; ok {
		if s, ok := v.(string); ok {
			p.Topic = s
		}
	}
	if v, ok := data["max_lines"]; ok {
		if n, ok := toInt(v); ok {
			p.MaxLines = n
		}
	}
	if v, ok := data["json_fields"]; ok {
		if s, ok := v.(string); ok {
			p.JSONFields = s
		}
	}
	if v, ok := data["style"]; ok {
		if s, ok := v.(string); ok {
			p.Style = s
		}
	}
	if v, ok := data["font"]; ok {
		if s, ok := v.(string); ok {
			p.Font = s
		}
	}

	SetPolicy(p)
	return nil
}

// toInt extracts an int from an interface value, handling both float64
// (JSON numbers decode as float64) and direct int types.
func toInt(v interface{}) (int, bool) {
	switch n := v.(type) {
	case float64:
		return int(n), true
	case int:
		return n, true
	case int64:
		return int(n), true
	default:
		return 0, false
	}
}

func init() {
	catalog.RegisterSnapshotter("zmq", zmqSnapshotter{})
}
