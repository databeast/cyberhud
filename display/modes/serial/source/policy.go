package source

import (
	"fmt"
	"strings"
	"time"

	"github.com/databeast/cyberhud/display/catalog"
)

// Policy controls runtime behavior for the serial monitor.
type Policy struct {
	Port       string
	Style      string // empty means auto-detect via style fitness
	Font       string // "auto" or a specific font ID
	Baud       int
	MaxLines   int
	ScanMS     int
	AutoSelect bool
	Raw        bool // true suppresses the instrument-panel chrome, showing data bytes only
}

// DefaultPolicy returns baseline serial monitor behavior.
func DefaultPolicy() Policy {
	return Policy{Port: "", Baud: DefaultBaud, MaxLines: DefaultMaxLines, AutoSelect: true, ScanMS: DefaultScanMS, Style: "", Font: "auto", Raw: false}
}

func (p Policy) Options() []catalog.OptionDefinition {
	return []catalog.OptionDefinition{
		{Key: "port", Type: "string", Summary: "Manual serial device path (disables auto-select when set).", Default: ""},
		{Key: "baud", Type: "int", Summary: "Serial baud rate.", Default: "115200"},
		{Key: "lines", Type: "int", Summary: "How many recent output lines to keep.", Default: "24"},
		{Key: "autoselect", Type: "bool", Summary: "Auto-pick the best available USB serial port.", Default: "true", Allowed: []string{"true", "false"}},
		{Key: "scan_ms", Type: "int", Summary: "Milliseconds between reconnect attempts.", Default: "500"},
		{Key: "font", Type: "string", Summary: "Font selection (auto or a registered font ID).", Default: "auto"},
		{Key: "raw", Type: "bool", Summary: "Show serial output bytes only, suppressing the connection/status chrome.", Default: "false", Allowed: []string{"true", "false"}},
	}
}

func (p Policy) Fingerprint() string {
	return fmt.Sprintf("%s|%d|%d|%t|%d|%s|%s|%t", p.Port, p.Baud, p.MaxLines, p.AutoSelect, p.ScanMS, p.Style, p.Font, p.Raw)
}

func (p Policy) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"port":       p.Port,
		"baud":       p.Baud,
		"lines":      p.MaxLines,
		"autoselect": p.AutoSelect,
		"scan_ms":    p.ScanMS,
		"style":      p.Style,
		"font":       p.Font,
		"raw":        p.Raw,
	}
}

// NormalizePolicy applies source-owned, registry-independent policy defaults.
func NormalizePolicy(p Policy) Policy {
	d := DefaultPolicy()
	p.Port = strings.TrimSpace(p.Port)
	if p.Baud <= 0 {
		p.Baud = d.Baud
	}
	if p.MaxLines <= 0 {
		p.MaxLines = d.MaxLines
	}
	if p.ScanMS <= 0 {
		p.ScanMS = d.ScanMS
	}
	if !p.AutoSelect && p.Port == "" && p.Baud == d.Baud && p.MaxLines == d.MaxLines && p.ScanMS == d.ScanMS {
		p.AutoSelect = d.AutoSelect
	}
	p.Style = strings.ToLower(strings.TrimSpace(p.Style))
	p.Font = strings.TrimSpace(p.Font)
	if p.Font == "" {
		p.Font = "auto"
	}
	return p
}

const (
	DefaultBaud        = 115200
	DefaultMaxLines    = 24
	DefaultScanMS      = 500
	DefaultReadTimeout = 500 * time.Millisecond
)
