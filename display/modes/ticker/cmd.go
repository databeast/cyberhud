package ticker

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/modes/ticker/source"
	"github.com/databeast/cyberhud/display/style/color"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// HandleConsoleCommand handles the top-level "ticker" console verb.
func HandleConsoleCommand(args []string) string {
	if len(args) < 1 {
		return "ERR usage: ticker <set|get|policy> ..."
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
	case "policy":
		if len(args) == 1 {
			return source.FormatPolicyResponse(PolicySnapshot())
		}
		p := PolicySnapshot()
		var appliedKeys []string
		for _, token := range args[1:] {
			kv := strings.SplitN(token, "=", 2)
			if len(kv) != 2 {
				return "ERR usage: ticker policy [style=<plain|bordered>] [font=<auto|font-id>] [font_tier=<auto|small|normal|large|fullsize>] [line_mode=<truncate|clip>] [direction=<vertical|horizontal|none>] [auto_scroll_ms=<n>] [accent=<name>] [show_border=<true|false>] [show_glow=<true|false>]"
			}
			key := strings.ToLower(strings.TrimSpace(kv[0]))
			val := strings.TrimSpace(kv[1])
			switch key {
			case "style":
				lower := strings.ToLower(val)
				if tickerRegistry.Lookup(lower) == nil {
					return fmt.Sprintf("ERR style: must be one of [%s]", strings.Join(registeredStyleNames(), ", "))
				}
				p.Style = lower
			case "font":
				if strings.TrimSpace(val) == "" {
					return "ERR font: must be \"auto\" or a registered font ID"
				}
				p.Font = val
			case "font_tier":
				if !source.ValidFontTier(val) {
					return fmt.Sprintf("ERR font_tier: must be one of [%s]", strings.Join(source.AllowedFontTiers, ", "))
				}
				p.FontTier = val
			case "accent":
				if !source.ValidAccent(val) {
					return fmt.Sprintf("ERR accent: must be one of [%s]", strings.Join(color.Names(), ", "))
				}
				p.Accent = val
			case "line_mode":
				p.LineMode = val
			case "direction":
				p.Direction = val
			case "auto_scroll_ms":
				n, err := strconv.Atoi(val)
				if err != nil {
					return fmt.Sprintf("ERR auto_scroll_ms must be an integer, got %q", val)
				}
				p.AutoScrollMS = n
			case "show_border":
				b, err := strconv.ParseBool(val)
				if err != nil {
					return "ERR show_border: must be true or false"
				}
				p.ShowBorder = b
			case "show_glow":
				b, err := strconv.ParseBool(val)
				if err != nil {
					return "ERR show_glow: must be true or false"
				}
				p.ShowGlow = b
			default:
				return fmt.Sprintf("ERR unknown ticker policy key %q", key)
			}
			appliedKeys = append(appliedKeys, key)
		}
		SetPolicy(p)
		resp := source.FormatPolicyResponse(PolicySnapshot())
		if notes := fitnessNotesPostApply(appliedKeys); len(notes) > 0 {
			resp += "\n" + strings.Join(notes, "\n")
		}
		return resp
	default:
		return "ERR usage: ticker <set|get|policy> ..."
	}
}

func init() {
	catalog.Register(catalog.Definition{
		ID:      "ticker",
		Title:   "Ticker",
		Scope:   "any",
		Summary: "Multi-style externally supplied text stream with color TFT, mono OLED, grayscale, and e-ink panel support.",
		Order:   70,
		Options: []catalog.OptionDefinition{
			{Key: "style", Type: "string", Summary: "Visual presentation style.", Default: "", Allowed: registeredStyleNames()},
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
		Usage:   "ticker <set|get|policy> ...",
		Handle:  HandleConsoleCommand,
	})
}
