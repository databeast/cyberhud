package stemma

import (
	"strconv"
	"strings"
	"time"

	displaymodes "github.com/databeast/cyberhud/display/modes"
	"github.com/databeast/cyberhud/display/modes/stemma/source"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/widgets/icons"
	"github.com/databeast/cyberhud/runtime/action"
)

func init() {
	displaymodes.RegisterFactory("stemma", newInstance)
}

// instance implements displaymodes.ModeInstance for the stemma mode.
// It holds a scanner reference obtained from the package-level singleton at
// construction time. If the scanner is nil (no I2C hardware available),
// BuildView gracefully returns valid ViewData with zero data items.
type instance struct {
	scanner         *source.Scanner // fetched from package-level singleton; may be nil
	scannerBuses    []string
	scannerInterval time.Duration
}

func newInstance() displaymodes.ModeInstance {
	buses, interval := source.GetScannerConfig()
	return &instance{
		scannerBuses:    buses,
		scannerInterval: interval,
	}
}

func (i *instance) ID() string { return "stemma" }

func (i *instance) Activate() {
	if s := source.GlobalScanner(); s != nil {
		i.scanner = s
		return
	}

	// Create and start scanner with configured bus names and interval.
	s := source.New(i.scannerBuses, i.scannerInterval)
	s.Start()
	i.scanner = s
	source.SetGlobalScanner(s)
}

func (i *instance) Deactivate() {
	if i.scanner != nil {
		i.scanner.Stop()
		i.scanner = nil
	}
	source.SetGlobalScanner(nil)
}

func (i *instance) ActionHandler() action.Handler { return Handler{} }

func (i *instance) BuildView() style.ViewData {
	if hints, ok := getPanelHints(); ok {
		return BuildView(i.devices(), hints, icons.Get, displaymodes.Warnings)
	}
	return style.ViewData{Items: []string{"error"}}
}

func (i *instance) RenderCacheKey() uint32 {
	return region.CalcRegionCacheKey(signatureDevices(i.devices()))
}

// devices returns the current device list from the scanner.
// Returns nil if the scanner is nil (graceful degradation — no panic).
func (i *instance) devices() []*source.Device {
	if i.scanner == nil {
		return nil
	}
	return i.scanner.Devices()
}

// signatureDevices builds a data-change signature string from stemma device states.
// The format encodes presence, bus, and address for each device.
func signatureDevices(devs []*source.Device) string {
	if len(devs) == 0 {
		return "dev:none"
	}

	var b strings.Builder
	for _, d := range devs {
		if d.Present {
			b.WriteByte('1')
		} else {
			b.WriteByte('0')
		}
		b.WriteString(":")
		b.WriteString(d.Bus)
		b.WriteString(":")
		b.WriteString(strconv.FormatUint(uint64(d.Addr), 16))
		b.WriteByte(';')
	}
	return b.String()
}

// SetScannerConfig stores I2C scanner configuration for future stemma instances.
func SetScannerConfig(buses []string, interval time.Duration) {
	source.SetScannerConfig(buses, interval)
}
