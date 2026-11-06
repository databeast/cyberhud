package ticker

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/catalog/cmdutil"
	"github.com/databeast/cyberhud/display/modes/ticker/source"
	"github.com/databeast/cyberhud/display/style/color"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// cmdHandler is the uniform CmdHandler for the ticker mode's policy options.
// Policy edits flow through it so applying key=value pairs fires
// cmdutil.OnApplied, driving centralized PolicyStore capture (freeze policy
// persistence) — the same path every conformant mode uses. The data sub-verbs
// (set/get) are handled by HandleConsoleCommand before delegating here.
var cmdHandler = &cmdutil.CmdHandler{
	Verb: "ticker",
	Keys: []cmdutil.KeyDef{
		{Name: "style", Validate: tickerStyleValidator},
		{Name: "font", Validate: tickerFontValidator},
		{Name: "font_tier", Validate: func(v string) string {
			if !source.ValidFontTier(v) {
				return fmt.Sprintf("must be one of [%s]", strings.Join(source.AllowedFontTiers, ", "))
			}
			return ""
		}},
		{Name: "line_mode"},
		{Name: "direction"},
		{Name: "auto_scroll_ms", Validate: func(v string) string {
			if _, err := strconv.Atoi(strings.TrimSpace(v)); err != nil {
				return fmt.Sprintf("must be an integer, got %q", v)
			}
			return ""
		}},
		{Name: "accent", Validate: func(v string) string {
			if !source.ValidAccent(v) {
				return fmt.Sprintf("must be one of [%s]", strings.Join(append(color.Names(), "none"), ", "))
			}
			return ""
		}},
		{Name: "show_border", Validate: cmdutil.BoolValidator()},
		{Name: "show_glow", Validate: cmdutil.BoolValidator()},
	},
	PostApply: fitnessNotesPostApply,
	Get:       getPolicyValue,
	Apply:     applyPolicyValue,
}

// tickerStyleValidator accepts only registered ticker style names.
func tickerStyleValidator(value string) string {
	if tickerRegistry.Lookup(strings.ToLower(strings.TrimSpace(value))) == nil {
		return fmt.Sprintf("must be one of [%s]", strings.Join(tickerRegistry.RegisteredNames(), ", "))
	}
	return ""
}

// tickerFontValidator accepts "auto" or any non-empty font ID.
func tickerFontValidator(value string) string {
	if strings.TrimSpace(value) == "" {
		return "must be \"auto\" or a registered font ID"
	}
	return ""
}

// getPolicyValue returns the current value for a given policy key.
func getPolicyValue(key string) string {
	p := PolicySnapshot()
	switch key {
	case "style":
		return p.Style
	case "font":
		return p.Font
	case "font_tier":
		return p.FontTier
	case "line_mode":
		return p.LineMode
	case "direction":
		return p.Direction
	case "auto_scroll_ms":
		return strconv.Itoa(p.AutoScrollMS)
	case "accent":
		return p.Accent
	case "show_border":
		return strconv.FormatBool(p.ShowBorder)
	case "show_glow":
		return strconv.FormatBool(p.ShowGlow)
	default:
		return ""
	}
}

// applyPolicyValue sets a single policy key on the current policy and persists
// it via SetPolicy (which normalizes). Values have already passed validation.
func applyPolicyValue(key, value string) {
	p := PolicySnapshot()
	switch key {
	case "style":
		p.Style = strings.ToLower(strings.TrimSpace(value))
	case "font":
		p.Font = value
	case "font_tier":
		p.FontTier = value
	case "line_mode":
		p.LineMode = value
	case "direction":
		p.Direction = value
	case "auto_scroll_ms":
		if n, err := strconv.Atoi(strings.TrimSpace(value)); err == nil {
			p.AutoScrollMS = n
		}
	case "accent":
		p.Accent = value
	case "show_border":
		if b, ok := cmdutil.ParseBool(value); ok {
			p.ShowBorder = b
		}
	case "show_glow":
		if b, ok := cmdutil.ParseBool(value); ok {
			p.ShowGlow = b
		}
	}
	SetPolicy(p)
}

// HandleConsoleCommand handles the top-level "ticker" console verb. The data
// sub-verbs "set" and "get" manage the externally-supplied ticker text/feed;
// everything else is treated as policy key=value edits and delegated to the
// uniform cmdutil.CmdHandler. A bare "ticker" with no args returns the current
// policy query line.
func HandleConsoleCommand(args []string) string {
	if len(args) == 0 {
		return cmdHandler.Handle(args)
	}
	sub := strings.ToLower(strings.TrimSpace(args[0]))
	switch sub {
	case "set":
		if len(args) < 2 {
			return "ERR usage: ticker set <text>"
		}
		text := strings.TrimSpace(strings.Join(args[1:], " "))
		if text == "" {
			source.SetText(nil)
			return "OK ticker lines=1"
		}

		if strings.HasPrefix(text, "[") {
			// JSON path
			directives, err := source.ParseJSONFeed(text)
			if err != nil {
				return fmt.Sprintf("ERR %s", err.Error())
			}
			source.SetFeed(directives)
			return fmt.Sprintf("OK ticker lines=%d", len(source.FeedSnapshot()))
		}

		// Legacy plain-text path
		source.SetText(strings.Split(text, "|"))
		return fmt.Sprintf("OK ticker lines=%d", len(source.Snapshot()))
	case "get":
		serialized, err := source.SerializeFeed()
		if err != nil {
			return fmt.Sprintf("ERR %s", err.Error())
		}
		return "OK\n" + serialized
	default:
		// Policy edits: style=... font=... etc. — uniform key=value path.
		return cmdHandler.Handle(args)
	}
}

func init() {
	catalog.Register(catalog.Definition{
		ID:       "ticker",
		Title:    "Ticker",
		Category: "Data",
		Scope:    "any",
		Summary:  "Multi-style externally supplied text stream with color TFT, mono OLED, grayscale, and e-ink panel support.",
		Order:    70,
		Options: []catalog.OptionDefinition{
			{Key: "style", Type: "string", Summary: "Visual presentation style.", Default: "", Allowed: tickerRegistry.RegisteredNames()},
			{Key: "font", Type: "string", Summary: "Font selection (auto or a registered font ID).", Default: "auto"},
			{Key: "font_tier", Type: "string", Summary: "Font size tier for panel-appropriate text rendering.", Default: "auto", Allowed: []string{"auto", "small", "normal", "large", "fullsize"}},
			{Key: "line_mode", Type: "string", Summary: "How long lines are constrained to panel width.", Default: textlayout.LineModeTruncate, Allowed: []string{textlayout.LineModeTruncate, textlayout.LineModeClip}},
			{Key: "direction", Type: "string", Summary: "Ticker auto-scroll direction when supported by the panel.", Default: textlayout.TickerDirectionVertical, Allowed: []string{textlayout.TickerDirectionVertical, "horizontal", textlayout.TickerDirectionNone}},
			{Key: "auto_scroll_ms", Type: "int", Summary: "Milliseconds between automatic ticker advances.", Default: "0"},
			{Key: "accent", Type: "string", Summary: "Named accent color from the shared palette.", Default: "cyan", Allowed: color.Names()},
			{Key: "show_border", Type: "bool", Summary: "Show decorative border frame.", Default: "false"},
			{Key: "show_glow", Type: "bool", Summary: "Show accent-colored glow background (ColorFast panels only).", Default: "false"},
		},
	})
	catalog.RegisterCommand(catalog.CommandDefinition{
		Verb:    "ticker",
		Summary: "Manage externally supplied ticker text and policy.",
		Usage:   "ticker set <text> | ticker get | ticker [style=<name>] [font=<auto|font-id>] [font_tier=<auto|small|normal|large|fullsize>] [line_mode=<truncate|clip>] [direction=<vertical|horizontal|none>] [auto_scroll_ms=<n>] [accent=<name>] [show_border=<true|false>] [show_glow=<true|false>]",
		Handle:  HandleConsoleCommand,
	})
}
