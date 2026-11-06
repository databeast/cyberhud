package attract_cube

import (
	"fmt"
	"strings"
	"sync"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/catalog/cmdutil"
	"github.com/databeast/cyberhud/display/modes/attract_cube/source"
)

// Policy controls the cube's color.
type Policy struct {
	BaseColor string
}

func DefaultPolicy() Policy          { return Policy{BaseColor: "cyan"} }
func (p Policy) Fingerprint() string { return p.BaseColor }
func (p Policy) ToMap() map[string]interface{} {
	return map[string]interface{}{"base_color": p.BaseColor}
}
func (Policy) Options() []catalog.OptionDefinition {
	return []catalog.OptionDefinition{{
		Key: "base_color", Type: "string", Default: "cyan",
		Summary: "Named accent color or #RRGGBB cube color.",
	}}
}

var policyState = struct {
	sync.RWMutex
	policy Policy
}{policy: DefaultPolicy()}

func GetPolicy() Policy {
	policyState.RLock()
	defer policyState.RUnlock()
	return policyState.policy
}

// SetPolicy validates the policy before replacing the current value.
func SetPolicy(p Policy) error {
	normalized, _, err := source.ParseBaseColor(p.BaseColor)
	if err != nil {
		return err
	}
	p.BaseColor = normalized
	policyState.Lock()
	policyState.policy = p
	policyState.Unlock()
	return nil
}

var command = cmdutil.CmdHandler{
	Verb: "attract_cube",
	Keys: []cmdutil.KeyDef{{Name: "base_color", Validate: func(value string) string {
		_, _, err := source.ParseBaseColor(value)
		if err != nil {
			return err.Error()
		}
		return ""
	}}},
	Get: func(string) string { return GetPolicy().BaseColor },
	Apply: func(_ string, value string) {
		policyState.Lock()
		policyState.policy.BaseColor = strings.ToLower(strings.TrimSpace(value))
		policyState.Unlock()
	},
}

type snapshotter struct{}

func (snapshotter) SnapshotPolicy() map[string]interface{} { return GetPolicy().ToMap() }
func (snapshotter) RestorePolicy(data map[string]interface{}) error {
	p := DefaultPolicy()
	if value, ok := data["base_color"]; ok {
		s, ok := value.(string)
		if !ok {
			return fmt.Errorf("base_color must be a string")
		}
		p.BaseColor = s
	}
	return SetPolicy(p)
}

func init() {
	catalog.Register(catalog.Definition{
		ID: "attract_cube", Title: "Cube", Category: "Attract", Order: 200,
		Summary: "A translucent cube tumbling in 3D.", Options: DefaultPolicy().Options(),
	})
	catalog.RegisterCommand(catalog.CommandDefinition{
		Verb: "attract_cube", Summary: "Query or set cube display options.",
		Usage: "attract_cube [base_color=<name|#RRGGBB>]", Handle: command.Handle,
	})
	catalog.RegisterSnapshotter("attract_cube", snapshotter{})
}
