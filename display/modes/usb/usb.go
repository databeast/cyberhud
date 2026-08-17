package usb

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/modes/usb/source"
	styles2 "github.com/databeast/cyberhud/display/modes/usb/styles"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/widgets/icons"

	"github.com/databeast/cyberhud/runtime/action"
)

const (
	defaultScanRoot     = "/sys/bus/usb/devices"
	defaultPollInterval = 500 * time.Millisecond
)

var usbRegistry = style.NewRegistry[source.Snapshot, source.Policy](
	// ── MonoSlow ──
	styles2.MonoSlow32x128Style,
	styles2.MonoSlow32x128MinimalStyle,
	styles2.MonoSlow64x128Style,
	styles2.MonoSlow64x128MinimalStyle,
	styles2.MonoSlow80x160Style,
	styles2.MonoSlow80x160MinimalStyle,
	styles2.MonoSlow104x212Style,
	styles2.MonoSlow104x212MinimalStyle,
	styles2.MonoSlow122x250Style,
	styles2.MonoSlow122x250MinimalStyle,
	styles2.MonoSlow128x32Style,
	styles2.MonoSlow128x32MinimalStyle,
	styles2.MonoSlow128x64Style,
	styles2.MonoSlow128x64MinimalStyle,
	styles2.MonoSlow128x128Style,
	styles2.MonoSlow128x128MinimalStyle,
	styles2.MonoSlow128x160Style,
	styles2.MonoSlow128x160MinimalStyle,
	styles2.MonoSlow128x296Style,
	styles2.MonoSlow128x296MinimalStyle,
	styles2.MonoSlow135x240Style,
	styles2.MonoSlow135x240MinimalStyle,
	styles2.MonoSlow160x80Style,
	styles2.MonoSlow160x80MinimalStyle,
	styles2.MonoSlow160x128Style,
	styles2.MonoSlow160x128MinimalStyle,
	styles2.MonoSlow176x264Style,
	styles2.MonoSlow176x264MinimalStyle,
	styles2.MonoSlow200x200Style,
	styles2.MonoSlow200x200MinimalStyle,
	styles2.MonoSlow212x104Style,
	styles2.MonoSlow212x104MinimalStyle,
	styles2.MonoSlow240x135Style,
	styles2.MonoSlow240x135MinimalStyle,
	styles2.MonoSlow240x240Style,
	styles2.MonoSlow240x240MinimalStyle,
	styles2.MonoSlow240x320Style,
	styles2.MonoSlow240x320MinimalStyle,
	styles2.MonoSlow250x122Style,
	styles2.MonoSlow250x122MinimalStyle,
	styles2.MonoSlow264x176Style,
	styles2.MonoSlow264x176MinimalStyle,
	styles2.MonoSlow296x128Style,
	styles2.MonoSlow296x128MinimalStyle,
	styles2.MonoSlow300x400Style,
	styles2.MonoSlow300x400MinimalStyle,
	styles2.MonoSlow320x240Style,
	styles2.MonoSlow320x240MinimalStyle,
	styles2.MonoSlow320x480Style,
	styles2.MonoSlow320x480MinimalStyle,
	styles2.MonoSlow400x300Style,
	styles2.MonoSlow400x300MinimalStyle,
	styles2.MonoSlow480x320Style,
	styles2.MonoSlow480x320MinimalStyle,
	styles2.MonoSlow480x800Style,
	styles2.MonoSlow480x800MinimalStyle,
	styles2.MonoSlow800x480Style,
	styles2.MonoSlow800x480MinimalStyle,

	// ── MonoFast ──
	styles2.MonoFast32x128Style,
	styles2.MonoFast32x128MinimalStyle,
	styles2.MonoFast64x128Style,
	styles2.MonoFast64x128MinimalStyle,
	styles2.MonoFast80x160Style,
	styles2.MonoFast80x160MinimalStyle,
	styles2.MonoFast104x212Style,
	styles2.MonoFast104x212MinimalStyle,
	styles2.MonoFast122x250Style,
	styles2.MonoFast122x250MinimalStyle,
	styles2.MonoFast128x32Style,
	styles2.MonoFast128x32MinimalStyle,
	styles2.MonoFast128x64Style,
	styles2.MonoFast128x64MinimalStyle,
	styles2.MonoFast128x128Style,
	styles2.MonoFast128x128MinimalStyle,
	styles2.MonoFast128x160Style,
	styles2.MonoFast128x160MinimalStyle,
	styles2.MonoFast128x296Style,
	styles2.MonoFast128x296MinimalStyle,
	styles2.MonoFast135x240Style,
	styles2.MonoFast135x240MinimalStyle,
	styles2.MonoFast160x80Style,
	styles2.MonoFast160x80MinimalStyle,
	styles2.MonoFast160x128Style,
	styles2.MonoFast160x128MinimalStyle,
	styles2.MonoFast176x264Style,
	styles2.MonoFast176x264MinimalStyle,
	styles2.MonoFast200x200Style,
	styles2.MonoFast200x200MinimalStyle,
	styles2.MonoFast212x104Style,
	styles2.MonoFast212x104MinimalStyle,
	styles2.MonoFast240x135Style,
	styles2.MonoFast240x135MinimalStyle,
	styles2.MonoFast240x240Style,
	styles2.MonoFast240x240MinimalStyle,
	styles2.MonoFast240x320Style,
	styles2.MonoFast240x320MinimalStyle,
	styles2.MonoFast250x122Style,
	styles2.MonoFast250x122MinimalStyle,
	styles2.MonoFast264x176Style,
	styles2.MonoFast264x176MinimalStyle,
	styles2.MonoFast296x128Style,
	styles2.MonoFast296x128MinimalStyle,
	styles2.MonoFast300x400Style,
	styles2.MonoFast300x400MinimalStyle,
	styles2.MonoFast320x240Style,
	styles2.MonoFast320x240MinimalStyle,
	styles2.MonoFast320x480Style,
	styles2.MonoFast320x480MinimalStyle,
	styles2.MonoFast400x300Style,
	styles2.MonoFast400x300MinimalStyle,
	styles2.MonoFast480x320Style,
	styles2.MonoFast480x320MinimalStyle,
	styles2.MonoFast480x800Style,
	styles2.MonoFast480x800MinimalStyle,
	styles2.MonoFast800x480Style,
	styles2.MonoFast800x480MinimalStyle,

	// ── GrayscaleSlow ──
	styles2.GrayscaleSlow32x128Style,
	styles2.GrayscaleSlow32x128MinimalStyle,
	styles2.GrayscaleSlow64x128Style,
	styles2.GrayscaleSlow64x128MinimalStyle,
	styles2.GrayscaleSlow80x160Style,
	styles2.GrayscaleSlow80x160MinimalStyle,
	styles2.GrayscaleSlow104x212Style,
	styles2.GrayscaleSlow104x212MinimalStyle,
	styles2.GrayscaleSlow122x250Style,
	styles2.GrayscaleSlow122x250MinimalStyle,
	styles2.GrayscaleSlow128x32Style,
	styles2.GrayscaleSlow128x32MinimalStyle,
	styles2.GrayscaleSlow128x64Style,
	styles2.GrayscaleSlow128x64MinimalStyle,
	styles2.GrayscaleSlow128x128Style,
	styles2.GrayscaleSlow128x128MinimalStyle,
	styles2.GrayscaleSlow128x160Style,
	styles2.GrayscaleSlow128x160MinimalStyle,
	styles2.GrayscaleSlow128x296Style,
	styles2.GrayscaleSlow128x296MinimalStyle,
	styles2.GrayscaleSlow135x240Style,
	styles2.GrayscaleSlow135x240MinimalStyle,
	styles2.GrayscaleSlow160x80Style,
	styles2.GrayscaleSlow160x80MinimalStyle,
	styles2.GrayscaleSlow160x128Style,
	styles2.GrayscaleSlow160x128MinimalStyle,
	styles2.GrayscaleSlow176x264Style,
	styles2.GrayscaleSlow176x264MinimalStyle,
	styles2.GrayscaleSlow200x200Style,
	styles2.GrayscaleSlow200x200MinimalStyle,
	styles2.GrayscaleSlow212x104Style,
	styles2.GrayscaleSlow212x104MinimalStyle,
	styles2.GrayscaleSlow240x135Style,
	styles2.GrayscaleSlow240x135MinimalStyle,
	styles2.GrayscaleSlow240x240Style,
	styles2.GrayscaleSlow240x240MinimalStyle,
	styles2.GrayscaleSlow240x320Style,
	styles2.GrayscaleSlow240x320MinimalStyle,
	styles2.GrayscaleSlow250x122Style,
	styles2.GrayscaleSlow250x122MinimalStyle,
	styles2.GrayscaleSlow264x176Style,
	styles2.GrayscaleSlow264x176MinimalStyle,
	styles2.GrayscaleSlow296x128Style,
	styles2.GrayscaleSlow296x128MinimalStyle,
	styles2.GrayscaleSlow300x400Style,
	styles2.GrayscaleSlow300x400MinimalStyle,
	styles2.GrayscaleSlow320x240Style,
	styles2.GrayscaleSlow320x240MinimalStyle,
	styles2.GrayscaleSlow320x480Style,
	styles2.GrayscaleSlow320x480MinimalStyle,
	styles2.GrayscaleSlow400x300Style,
	styles2.GrayscaleSlow400x300MinimalStyle,
	styles2.GrayscaleSlow480x320Style,
	styles2.GrayscaleSlow480x320MinimalStyle,
	styles2.GrayscaleSlow480x800Style,
	styles2.GrayscaleSlow480x800MinimalStyle,
	styles2.GrayscaleSlow800x480Style,
	styles2.GrayscaleSlow800x480MinimalStyle,

	// ── GrayscaleFast ──
	styles2.GrayscaleFast32x128Style,
	styles2.GrayscaleFast32x128MinimalStyle,
	styles2.GrayscaleFast64x128Style,
	styles2.GrayscaleFast64x128MinimalStyle,
	styles2.GrayscaleFast80x160Style,
	styles2.GrayscaleFast80x160MinimalStyle,
	styles2.GrayscaleFast104x212Style,
	styles2.GrayscaleFast104x212MinimalStyle,
	styles2.GrayscaleFast122x250Style,
	styles2.GrayscaleFast122x250MinimalStyle,
	styles2.GrayscaleFast128x32Style,
	styles2.GrayscaleFast128x32MinimalStyle,
	styles2.GrayscaleFast128x64Style,
	styles2.GrayscaleFast128x64MinimalStyle,
	styles2.GrayscaleFast128x128Style,
	styles2.GrayscaleFast128x128MinimalStyle,
	styles2.GrayscaleFast128x160Style,
	styles2.GrayscaleFast128x160MinimalStyle,
	styles2.GrayscaleFast128x296Style,
	styles2.GrayscaleFast128x296MinimalStyle,
	styles2.GrayscaleFast135x240Style,
	styles2.GrayscaleFast135x240MinimalStyle,
	styles2.GrayscaleFast160x80Style,
	styles2.GrayscaleFast160x80MinimalStyle,
	styles2.GrayscaleFast160x128Style,
	styles2.GrayscaleFast160x128MinimalStyle,
	styles2.GrayscaleFast176x264Style,
	styles2.GrayscaleFast176x264MinimalStyle,
	styles2.GrayscaleFast200x200Style,
	styles2.GrayscaleFast200x200MinimalStyle,
	styles2.GrayscaleFast212x104Style,
	styles2.GrayscaleFast212x104MinimalStyle,
	styles2.GrayscaleFast240x135Style,
	styles2.GrayscaleFast240x135MinimalStyle,
	styles2.GrayscaleFast240x240Style,
	styles2.GrayscaleFast240x240MinimalStyle,
	styles2.GrayscaleFast240x320Style,
	styles2.GrayscaleFast240x320MinimalStyle,
	styles2.GrayscaleFast250x122Style,
	styles2.GrayscaleFast250x122MinimalStyle,
	styles2.GrayscaleFast264x176Style,
	styles2.GrayscaleFast264x176MinimalStyle,
	styles2.GrayscaleFast296x128Style,
	styles2.GrayscaleFast296x128MinimalStyle,
	styles2.GrayscaleFast300x400Style,
	styles2.GrayscaleFast300x400MinimalStyle,
	styles2.GrayscaleFast320x240Style,
	styles2.GrayscaleFast320x240MinimalStyle,
	styles2.GrayscaleFast320x480Style,
	styles2.GrayscaleFast320x480MinimalStyle,
	styles2.GrayscaleFast400x300Style,
	styles2.GrayscaleFast400x300MinimalStyle,
	styles2.GrayscaleFast480x320Style,
	styles2.GrayscaleFast480x320MinimalStyle,
	styles2.GrayscaleFast480x800Style,
	styles2.GrayscaleFast480x800MinimalStyle,
	styles2.GrayscaleFast800x480Style,
	styles2.GrayscaleFast800x480MinimalStyle,

	// ── ColorSlow ──
	styles2.ColorSlow32x128Style,
	styles2.ColorSlow32x128MinimalStyle,
	styles2.ColorSlow64x128Style,
	styles2.ColorSlow64x128MinimalStyle,
	styles2.ColorSlow80x160Style,
	styles2.ColorSlow80x160MinimalStyle,
	styles2.ColorSlow104x212Style,
	styles2.ColorSlow104x212MinimalStyle,
	styles2.ColorSlow122x250Style,
	styles2.ColorSlow122x250MinimalStyle,
	styles2.ColorSlow128x32Style,
	styles2.ColorSlow128x32MinimalStyle,
	styles2.ColorSlow128x64Style,
	styles2.ColorSlow128x64MinimalStyle,
	styles2.ColorSlow128x128Style,
	styles2.ColorSlow128x128MinimalStyle,
	styles2.ColorSlow128x160Style,
	styles2.ColorSlow128x160MinimalStyle,
	styles2.ColorSlow128x296Style,
	styles2.ColorSlow128x296MinimalStyle,
	styles2.ColorSlow135x240Style,
	styles2.ColorSlow135x240MinimalStyle,
	styles2.ColorSlow160x80Style,
	styles2.ColorSlow160x80MinimalStyle,
	styles2.ColorSlow160x128Style,
	styles2.ColorSlow160x128MinimalStyle,
	styles2.ColorSlow176x264Style,
	styles2.ColorSlow176x264MinimalStyle,
	styles2.ColorSlow200x200Style,
	styles2.ColorSlow200x200MinimalStyle,
	styles2.ColorSlow212x104Style,
	styles2.ColorSlow212x104MinimalStyle,
	styles2.ColorSlow240x135Style,
	styles2.ColorSlow240x135MinimalStyle,
	styles2.ColorSlow240x240Style,
	styles2.ColorSlow240x240MinimalStyle,
	styles2.ColorSlow240x320Style,
	styles2.ColorSlow240x320MinimalStyle,
	styles2.ColorSlow250x122Style,
	styles2.ColorSlow250x122MinimalStyle,
	styles2.ColorSlow264x176Style,
	styles2.ColorSlow264x176MinimalStyle,
	styles2.ColorSlow296x128Style,
	styles2.ColorSlow296x128MinimalStyle,
	styles2.ColorSlow300x400Style,
	styles2.ColorSlow300x400MinimalStyle,
	styles2.ColorSlow320x240Style,
	styles2.ColorSlow320x240MinimalStyle,
	styles2.ColorSlow320x480Style,
	styles2.ColorSlow320x480MinimalStyle,
	styles2.ColorSlow400x300Style,
	styles2.ColorSlow400x300MinimalStyle,
	styles2.ColorSlow480x320Style,
	styles2.ColorSlow480x320MinimalStyle,
	styles2.ColorSlow480x800Style,
	styles2.ColorSlow480x800MinimalStyle,
	styles2.ColorSlow800x480Style,
	styles2.ColorSlow800x480MinimalStyle,

	// ── ColorFast ──
	styles2.ColorFast32x128Style,
	styles2.ColorFast32x128MinimalStyle,
	styles2.ColorFast64x128Style,
	styles2.ColorFast64x128MinimalStyle,
	styles2.ColorFast80x160Style,
	styles2.ColorFast80x160MinimalStyle,
	styles2.ColorFast104x212Style,
	styles2.ColorFast104x212MinimalStyle,
	styles2.ColorFast122x250Style,
	styles2.ColorFast122x250MinimalStyle,
	styles2.ColorFast128x32Style,
	styles2.ColorFast128x32MinimalStyle,
	styles2.ColorFast128x64Style,
	styles2.ColorFast128x64MinimalStyle,
	styles2.ColorFast128x128Style,
	styles2.ColorFast128x128MinimalStyle,
	styles2.ColorFast128x160Style,
	styles2.ColorFast128x160MinimalStyle,
	styles2.ColorFast128x296Style,
	styles2.ColorFast128x296MinimalStyle,
	styles2.ColorFast135x240Style,
	styles2.ColorFast135x240MinimalStyle,
	styles2.ColorFast160x80Style,
	styles2.ColorFast160x80MinimalStyle,
	styles2.ColorFast160x128Style,
	styles2.ColorFast160x128MinimalStyle,
	styles2.ColorFast176x264Style,
	styles2.ColorFast176x264MinimalStyle,
	styles2.ColorFast200x200Style,
	styles2.ColorFast200x200MinimalStyle,
	styles2.ColorFast212x104Style,
	styles2.ColorFast212x104MinimalStyle,
	styles2.ColorFast240x135Style,
	styles2.ColorFast240x135MinimalStyle,
	styles2.ColorFast240x240Style,
	styles2.ColorFast240x240MinimalStyle,
	styles2.ColorFast240x320Style,
	styles2.ColorFast240x320MinimalStyle,
	styles2.ColorFast250x122Style,
	styles2.ColorFast250x122MinimalStyle,
	styles2.ColorFast264x176Style,
	styles2.ColorFast264x176MinimalStyle,
	styles2.ColorFast296x128Style,
	styles2.ColorFast296x128MinimalStyle,
	styles2.ColorFast300x400Style,
	styles2.ColorFast300x400MinimalStyle,
	styles2.ColorFast320x240Style,
	styles2.ColorFast320x240MinimalStyle,
	styles2.ColorFast320x480Style,
	styles2.ColorFast320x480MinimalStyle,
	styles2.ColorFast400x300Style,
	styles2.ColorFast400x300MinimalStyle,
	styles2.ColorFast480x320Style,
	styles2.ColorFast480x320MinimalStyle,
	styles2.ColorFast480x800Style,
	styles2.ColorFast480x800MinimalStyle,
	styles2.ColorFast800x480Style,
	styles2.ColorFast800x480MinimalStyle,
)

// registeredStyleNames returns the list of style names from the registry.
// Used by catalog registration and cmdHandler for allowed-value validation.
func registeredStyleNames() []string {
	styles := usbRegistry.Enumerate()
	names := make([]string, len(styles))
	for i, s := range styles {
		names[i] = s.Name()
	}
	return names
}

// usbStyleNamePattern is the regex that all registered USB style names must match.
// Format: <capability_prefix><WxH>[-<variant>]
var usbStyleNamePattern = regexp.MustCompile(
	`^(mono-slow-|mono-fast-|mono-|grayscale-slow-|grayscale-fast-|color-slow-|color-fast-|color-)\d+x\d+(-[a-z0-9-]+)?$`,
)

func init() {
	for _, s := range usbRegistry.Enumerate() {
		if !usbStyleNamePattern.MatchString(s.Name()) {
			panic(fmt.Sprintf("usb: invalid style name %q does not match expected pattern %s", s.Name(), usbStyleNamePattern.String()))
		}
	}
}

func init() {
	catalog.Register(catalog.Definition{
		ID:      "usb",
		Title:   "USB Bench",
		Scope:   "any",
		Summary: "Highlights the most recently connected USB device for bench identification.",
		Order:   35,
		Options: []catalog.OptionDefinition{
			{Key: "poll_ms", Type: "int", Summary: "Milliseconds between fallback polling scans.", Default: "500"},
			{Key: "hold_unplugged_ms", Type: "int", Summary: "How long to keep a disconnected device visible (0 keeps indefinitely).", Default: "0"},
			{Key: "hide_root_hubs", Type: "bool", Summary: "Hide USB root hubs from bench identification results.", Default: "true", Allowed: []string{"true", "false"}},
			{Key: "style", Type: "string", Summary: "Visual presentation style.", Default: "", Allowed: registeredStyleNames()},
		},
	})
	catalog.RegisterCommand(catalog.CommandDefinition{
		Verb:    "usb",
		Summary: "Inspect and configure USB bench detection policy.",
		Usage:   "usb [poll_ms=<n>] [hold_unplugged_ms=<n>] [hide_root_hubs=<true|false>] [show_border=<true|false>] [style=<mono-slow-128x64|color-fast-240x320-minimal|...>]",
		Handle:  HandleConsoleCommand,
	})
}

// Handler implements action.Handler for USB bench mode.
type Handler struct{}

func (Handler) HandleAction(act action.Action, cursor, itemCount int) action.Result {
	switch act {
	case action.Primary:
		return action.Result{Navigate: "dashboard"}
	}
	return action.Result{}
}

var monitorState = struct {
	sync.Mutex
	known        map[string]source.DeviceInfo
	snapshot     Snapshot
	lastScanAt   time.Time
	scanRoot     string
	pollInterval time.Duration
	policy       Policy
	events       chan struct{}
	initOnce     sync.Once
	instant      bool
	now          func() time.Time
}{
	known:        map[string]source.DeviceInfo{},
	scanRoot:     defaultScanRoot,
	pollInterval: defaultPollInterval,
	policy:       DefaultPolicy(),
	events:       make(chan struct{}, 1),
	now:          time.Now,
}

// SnapshotNow returns the latest USB bench snapshot, refreshing from sysfs when stale.
func SnapshotNow() Snapshot {
	refresh()
	monitorState.Lock()
	defer monitorState.Unlock()
	return monitorState.snapshot
}

// BuildView returns the USB bench view data as a complete style.ViewData,
// including Title, Hint, and Static fields. All visual output (text, sprites,
// icons) is produced by the style's Build() method.
func BuildView() style.ViewData {
	snap := SnapshotNow()
	p := PolicySnapshot()
	p = normalizePolicy(p)
	hints, ok := getPanelHints()
	if !ok {
		return style.ViewData{Items: []string{"error"}}
	}

	s, reason := style.ResolveStyle(usbRegistry, hints, "usb", p.Style)

	ctx := style.NewStyleContext(hints)
	if consumer, ok := any(s).(styles2.IconConsumer); ok {
		consumer.SetIconGetter(icons.Get)
	}
	vd := s.Build(snap, p, ctx)

	vd.StyleReport = style.StyleReport{Name: s.Name(), Reason: reason}
	vd.Static = true
	return vd
}

// Signature returns a change token used by the UI refresh loop.
func Signature() uint32 {
	snap := SnapshotNow()
	p := PolicySnapshot()
	state := "0"
	if snap.Connected {
		state = "1"
	}
	return region.CalcRegionCacheKey("usb", snap.Sequence, state, snap.Device.Key, p.PollMS, p.HoldUnpluggedMS, p.HideRootHubs, p.Style)
}

func refresh() {
	initMonitor()

	monitorState.Lock()
	defer monitorState.Unlock()
	now := monitorState.now()
	force := drainEventsLocked()
	monitorState.pollInterval = time.Duration(monitorState.policy.PollMS) * time.Millisecond
	if !force && !monitorState.lastScanAt.IsZero() && now.Sub(monitorState.lastScanAt) < monitorState.pollInterval {
		return
	}
	curr, _ := source.ScanDevices(monitorState.scanRoot, monitorState.policy)
	nextKnown, nextSnapshot, changed := source.Transition(monitorState.known, curr, monitorState.snapshot, now)
	if applyHoldPolicy(&nextSnapshot, now, monitorState.policy) {
		changed = true
	}
	if changed {
		nextSnapshot.Sequence = monitorState.snapshot.Sequence + 1
	} else {
		nextSnapshot.Sequence = monitorState.snapshot.Sequence
	}
	monitorState.known = nextKnown
	monitorState.snapshot = nextSnapshot
	monitorState.lastScanAt = now
}

func initMonitor() {
	monitorState.initOnce.Do(func() {
		err := source.StartInstantMonitor(monitorState.scanRoot, monitorState.events)
		if err != nil {
			log.Printf("usb mode: instant monitoring unavailable (%v); using polling fallback", err)
			return
		}
		monitorState.Lock()
		monitorState.instant = true
		monitorState.Unlock()
		log.Printf("usb mode: instant monitoring enabled on %s", monitorState.scanRoot)
	})
}

func drainEventsLocked() bool {
	forced := false
	for {
		select {
		case <-monitorState.events:
			forced = true
		default:
			return forced
		}
	}
}

// boolStr returns "true" or "false" for a bool value.
func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// SetPolicy updates the USB mode policy.
func SetPolicy(policy source.Policy) {
	monitorState.Lock()
	monitorState.policy = normalizePolicy(policy)
	monitorState.lastScanAt = time.Time{}
	monitorState.Unlock()
	select {
	case monitorState.events <- struct{}{}:
	default:
	}
}

// PolicySnapshot returns the current USB mode policy.
func PolicySnapshot() source.Policy {
	monitorState.Lock()
	defer monitorState.Unlock()
	return monitorState.policy
}

func displayName(info source.DeviceInfo) string {
	name := strings.TrimSpace(info.Manufacturer + " " + info.Product)
	if name != "" {
		return name
	}
	if info.Product != "" {
		return info.Product
	}
	return "Unknown USB Device"
}

func safeValue(v, fallback string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return fallback
	}
	return v
}
