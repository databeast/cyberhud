package attract_matrix

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/catalog/cmdutil"
	"github.com/databeast/cyberhud/display/modes/attract_matrix/source"
)

// policyState holds the mutex-protected active policy.
var policyState = struct {
	sync.RWMutex
	policy source.Policy
}{
	policy: source.DefaultPolicy(),
}

// GetPolicy returns the current matrix policy (thread-safe read under RWMutex).
//
// to the shared policy state, enabling render loops and CLI queries to observe
// consistent snapshots without blocking writers.
func GetPolicy() source.Policy {
	policyState.RLock()
	defer policyState.RUnlock()
	return policyState.policy
}

// SetPolicy updates the matrix policy after normalization (thread-safe write under Mutex).
//
// incoming policy values before persisting them, ensuring the active policy always
// contains well-formed data regardless of caller input.
func SetPolicy(p source.Policy) {
	policyState.Lock()
	defer policyState.Unlock()
	policyState.policy = normalizePolicy(p)
}

func DefaultPolicy() source.Policy {
	return source.DefaultPolicy()
}

func normalizePolicy(p source.Policy) source.Policy {
	return source.NormalizePolicy(p)
}

func policyFingerprint(p source.Policy) string {
	return source.PolicyFingerprint(p)
}

// matrixSnapshotter implements catalog.PolicySnapshotter for the attract_matrix mode.
type matrixSnapshotter struct{}

// SnapshotPolicy returns the current matrix policy as a JSON-serializable map
// using snake_case keys matching the protocol wire format.
func (matrixSnapshotter) SnapshotPolicy() map[string]interface{} {
	p := GetPolicy()
	return p.ToMap()
}

// RestorePolicy applies policy values from a JSON map, running them through
// the same normalization as SetPolicy. Keys use the protocol wire format.
func (matrixSnapshotter) RestorePolicy(data map[string]interface{}) error {
	p := DefaultPolicy()

	if v, ok := data["min_speed"]; ok {
		if f, ok := toFloat64(v); ok {
			p.MinSpeed = f
		}
	}
	if v, ok := data["max_speed"]; ok {
		if f, ok := toFloat64(v); ok {
			p.MaxSpeed = f
		}
	}
	if v, ok := data["trail_length"]; ok {
		if n, ok := toInt(v); ok {
			p.TrailLength = n
		}
	}
	if v, ok := data["density"]; ok {
		if f, ok := toFloat64(v); ok {
			p.Density = f
		}
	}
	if v, ok := data["show_background"]; ok {
		if b, ok := toBool(v); ok {
			p.ShowBackground = b
		}
	}

	SetPolicy(p)
	return nil
}

// toFloat64 extracts a float64 from an interface value, handling both
// float64 (native JSON number) and int/int64 conversions.
func toFloat64(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	default:
		return 0, false
	}
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

// toBool extracts a bool from an interface value.
func toBool(v interface{}) (bool, bool) {
	b, ok := v.(bool)
	return b, ok
}

func init() {
	catalog.RegisterSnapshotter("attract_matrix", matrixSnapshotter{})
}

// floatPositiveValidator returns a KeyValidator that accepts float values > 0.
func floatPositiveValidator() cmdutil.KeyValidator {
	return func(value string) string {
		f, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
		if err != nil {
			return "must be a number"
		}
		if f <= 0 {
			return "must be > 0"
		}
		return ""
	}
}

// floatRangeValidator returns a KeyValidator that accepts float values in [min, max].
func floatRangeValidator(min, max float64) cmdutil.KeyValidator {
	return func(value string) string {
		f, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
		if err != nil {
			return "must be a number"
		}
		if f < min || f > max {
			return fmt.Sprintf("must be in [%.1f, %.1f]", min, max)
		}
		return ""
	}
}

// intRangeValidator returns a KeyValidator that accepts integer values in [min, max].
func intRangeValidator(min, max int) cmdutil.KeyValidator {
	return func(value string) string {
		v := strings.TrimSpace(value)
		n, err := strconv.Atoi(v)
		if err != nil {
			return "must be an integer"
		}
		if n < min || n > max {
			return fmt.Sprintf("must be in [%d, %d]", min, max)
		}
		return ""
	}
}

// getPolicy returns the current value for a given policy key.
func getPolicy(key string) string {
	p := GetPolicy()
	switch key {
	case "min_speed":
		return fmt.Sprintf("%g", p.MinSpeed)
	case "max_speed":
		return fmt.Sprintf("%g", p.MaxSpeed)
	case "trail_length":
		return fmt.Sprintf("%d", p.TrailLength)
	case "density":
		return fmt.Sprintf("%g", p.Density)
	case "show_background":
		if p.ShowBackground {
			return "true"
		}
		return "false"
	default:
		return ""
	}
}

// applyPolicy updates a single policy key with its validated value.
func applyPolicy(key, value string) {
	policyState.Lock()
	defer policyState.Unlock()
	switch key {
	case "min_speed":
		if f, err := strconv.ParseFloat(strings.TrimSpace(value), 64); err == nil {
			policyState.policy.MinSpeed = f
		}
	case "max_speed":
		if f, err := strconv.ParseFloat(strings.TrimSpace(value), 64); err == nil {
			policyState.policy.MaxSpeed = f
		}
	case "trail_length":
		if n, err := strconv.Atoi(strings.TrimSpace(value)); err == nil {
			policyState.policy.TrailLength = n
		}
	case "density":
		if f, err := strconv.ParseFloat(strings.TrimSpace(value), 64); err == nil {
			policyState.policy.Density = f
		}
	case "show_background":
		if v, ok := cmdutil.ParseBool(value); ok {
			policyState.policy.ShowBackground = v
		}
	}
	policyState.policy = normalizePolicy(policyState.policy)
}
