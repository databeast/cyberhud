package displaymode

import (
	"strings"
	"sync"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/catalog/cmdutil"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/runtime/action"
)

// Compile-time check: Handler implements action.Handler.
var _ action.Handler = Handler{}

// Policy captures the runtime-configurable parameters for the template mode.
type Policy struct {
	Style string
}

var _ catalog.ConfigPolicy = Policy{}

func (p Policy) Fingerprint() string           { return "template|" + p.Style }
func (p Policy) ToMap() map[string]interface{} { return map[string]interface{}{"style": p.Style} }
func (p Policy) Options() []catalog.OptionDefinition {
	return []catalog.OptionDefinition{
		{Key: "style", Type: "string", Summary: "Visual style variant.", Default: "", Allowed: registeredStyleNames()},
	}
}

// DefaultPolicy returns the template default policy.
func DefaultPolicy() Policy {
	return Policy{Style: ""}
}

// policyState holds the mutex-protected active policy.
var policyState = struct {
	sync.RWMutex
	policy Policy
}{policy: DefaultPolicy()}

// GetPolicy returns the current policy (thread-safe).
func GetPolicy() Policy {
	policyState.RLock()
	defer policyState.RUnlock()
	return policyState.policy
}

// SetPolicy updates the policy after normalization (thread-safe).
func SetPolicy(p Policy) {
	policyState.Lock()
	defer policyState.Unlock()
	policyState.policy = normalizePolicy(p)
}

// normalizePolicy trims, lowercases, and validates the Style field against the registry.
func normalizePolicy(p Policy) Policy {
	p.Style = strings.ToLower(strings.TrimSpace(p.Style))
	if templateRegistry.Lookup(p.Style) == nil {
		p.Style = ""
	}
	return p
}

// cmdHandler provides CLI verb processing for the displaymode template.
var cmdHandler = &cmdutil.CmdHandler{
	Verb: "displaymode",
	Keys: []cmdutil.KeyDef{
		{Name: "style", Validate: cmdutil.AllowedValidator(registeredStyleNames())},
	},
	Get: func(key string) string {
		p := GetPolicy()
		switch key {
		case "style":
			if p.Style != "" {
				return p.Style
			}
			hints, _ := getPanelHints()
			s, _ := style.ResolveStyle(templateRegistry, hints, "displaymode", "")
			return s.Name()
		}
		return ""
	},
	Apply: func(key, value string) {
		policyState.Lock()
		defer policyState.Unlock()
		switch key {
		case "style":
			policyState.policy.Style = strings.ToLower(strings.TrimSpace(value))
		}
	},
	PostApply: fitnessNotesPostApply,
}

// HandleCommand processes displaymode console commands.
func HandleCommand(args []string) string {
	return cmdHandler.Handle(args)
}

// RenderCacheKey returns a change-detection string derived from the current policy.
func RenderCacheKey() string {
	p := GetPolicy()
	return "template|" + p.Style
}

// BuildView performs registry-based style dispatch and returns mode-internal ViewData.
func BuildView(hints textlayout.TextHints) ViewData {
	pol := GetPolicy()
	s, _ := style.ResolveStyle(templateRegistry, hints, "displaymode", pol.Style)
	ctx := style.NewStyleContext(hints)
	svd := s.Build(Snapshot{}, pol, ctx)
	return convertFromStyleViewData(svd)
}

// Handler implements action.Handler for the template display mode.
type Handler struct{}

// HandleAction processes logical UI inputs. Left/Right return Dirty=true.
func (Handler) HandleAction(act action.Action, cursor, itemCount int) action.Result {
	switch act {
	case action.Left, action.Right:
		return action.Result{Dirty: true}
	default:
		return action.Result{}
	}
}

func init() {
	catalog.Register(catalog.Definition{
		ID:      "displaymode",
		Title:   "Display Mode Template",
		Summary: "Skeleton template demonstrating all display mode framework patterns.",
		Order:   999,
		Options: []catalog.OptionDefinition{
			{Key: "style", Type: "string", Summary: "Visual style variant.", Default: "", Allowed: registeredStyleNames()},
		},
	})
	catalog.RegisterCommand(catalog.CommandDefinition{
		Verb:    "displaymode",
		Summary: "Query or set template display mode options.",
		Usage:   "displaymode [style=<name>]",
		Handle:  HandleCommand,
	})
}
