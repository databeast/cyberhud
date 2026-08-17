package main

import (
	"fmt"
	"image"
	"io"
	"log"
	"sort"
	"strings"
	"time"

	displaymodes "github.com/databeast/cyberhud/display/modes"
	"github.com/databeast/cyberhud/display/region/modehints"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/coordinator"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/hardware/driver"
	"github.com/databeast/cyberhud/hardware/input"
	"github.com/databeast/cyberhud/hardware/panels"
	"github.com/databeast/cyberhud/runtime/ui"
)

// isI2COnly reports whether a driver definition supports only I2C (no SPI factory).
func isI2COnly(drv driver.Definition) bool {
	return drv.NewI2C != nil && drv.NewSPI == nil
}

// indexedTarget pairs a display index with its draw target.
type indexedTarget struct {
	index  int
	target driver.DrawTarget
}

func initDisplayTargets(profile panels.Definition, spiDevice string) ([]region.ScreenPosition, func(), string, error) {
	// Determine communication mode for logging.
	commMode := "SPI"
	if profile.I2CBus != "" {
		commMode = "I2C"
	}

	controller := strings.ToLower(strings.TrimSpace(profile.Controller))
	if controller == "" {
		controller = "st7789"
	}
	log.Printf("display: init panel=%q controller=%s mode=%s", profile.Name, controller, commMode)

	// Log available buses before any open attempt.
	i2cBuses := i2creg.All()
	spiBuses := spireg.All()
	i2cNames := make([]string, len(i2cBuses))
	for i, ref := range i2cBuses {
		i2cNames[i] = ref.Name
	}
	spiNames := make([]string, len(spiBuses))
	for i, ref := range spiBuses {
		spiNames[i] = ref.Name
	}
	log.Printf("display: available buses i2c=[%s] spi=[%s]", strings.Join(i2cNames, " "), strings.Join(spiNames, " "))

	// Log output pin assignment.
	pinOrUnassigned := func(name string) string {
		if name == "" {
			return "unassigned"
		}
		return name
	}
	log.Printf("display: pin assignment DC=%s RST=%s BL=%s BUSY=%s",
		pinOrUnassigned(profile.DCPin),
		pinOrUnassigned(profile.RSTPin),
		pinOrUnassigned(profile.BLPin),
		pinOrUnassigned(profile.BusyPin))

	// targets collects (index, DrawTarget) pairs in insertion order;
	// we sort by index later to produce left-to-right positioning.
	targets := make([]indexedTarget, 0, 4)

	closers := make([]io.Closer, 0, 4)
	cleanup := func() {
		for i := len(closers) - 1; i >= 0; i-- {
			_ = closers[i].Close()
		}
	}

	if len(profile.Virtual) == 0 {
		drv, ok := driver.Get(controller)
		if !ok {
			log.Printf("display: unrecognized controller %q; registered drivers: %v", profile.Controller, driver.Names())
			return nil, func() {}, "", fmt.Errorf("unsupported display controller %q", profile.Controller)
		}

		if isI2COnly(drv) || profile.I2CBus != "" {
			// I2C path: skip GPIO pin resolution entirely.
			log.Printf("display: i2c bus=%q addr=0x%02X", profile.I2CBus, profile.Config.I2CAddr)
			primary, bus, controller, err := initI2CDisplayTarget(profile.I2CBus, profile.Controller, profile.Config)
			if err != nil {
				return nil, func() {}, "", err
			}
			closers = append(closers, bus)
			targets = append(targets, indexedTarget{index: 0, target: primary})
			b := primary.Bounds()
			w, h := b.Dx(), b.Dy()
			log.Printf("display: %s ready (%dx%d %s)", profile.Name, w, h, strings.ToUpper(controller))

			positions, err := buildPositions(targets, profile)
			if err != nil {
				cleanup()
				return nil, func() {}, "", err
			}
			return positions, cleanup, controller, nil
		}

		// SPI path: resolve GPIO pins and open the SPI port for a single-screen panel.
		primary, port, controller, err := initDisplayTarget(panels.Screen{
			SPI:        spiDevice,
			Controller: profile.Controller,
			Config:     profile.Config,
			DCPin:      profile.DCPin,
			RSTPin:     profile.RSTPin,
			BusyPin:    profile.BusyPin,
			BLPin:      profile.BLPin,
		})
		if err != nil {
			return nil, func() {}, "", err
		}
		closers = append(closers, port)
		targets = append(targets, indexedTarget{index: 0, target: primary})
		b := primary.Bounds()
		w, h := b.Dx(), b.Dy()
		log.Printf("display: %s ready (%dx%d %s)", profile.Name, w, h, strings.ToUpper(controller))

		positions, err := buildPositions(targets, profile)
		if err != nil {
			cleanup()
			return nil, func() {}, "", err
		}
		return positions, cleanup, controller, nil
	}

	// Pre-flight: verify all required SPI buses can be opened before touching any GPIO pins.
	// This prevents driving DC/RST/BL pins for partially-initialized screens when a bus is missing.
	requiredSPI := map[string]bool{}
	for _, vp := range profile.Virtual {
		if name := strings.TrimSpace(vp.SPI); name != "" {
			requiredSPI[name] = true
		}
	}
	var missingSPI []string
	for name := range requiredSPI {
		port, err := spireg.Open(name)
		if err != nil {
			missingSPI = append(missingSPI, name)
		} else {
			_ = port.Close()
		}
	}
	if len(missingSPI) > 0 {
		sort.Strings(missingSPI)
		log.Printf("display: panel %q requires SPI bus(es) %v but they could not be opened", profile.Name, missingSPI)
		return nil, func() {}, "", fmt.Errorf("display: panel %q requires SPI bus(es) %v but they could not be opened; check dtoverlay settings and reboot", profile.Name, missingSPI)
	}

	primaryController := ""
	seen := map[int]bool{}
	for i, vp := range profile.Virtual {
		if seen[vp.Index] {
			cleanup()
			return nil, func() {}, "", fmt.Errorf("display: duplicate virtual panel index %d", vp.Index)
		}
		seen[vp.Index] = true

		vpController := strings.ToLower(strings.TrimSpace(vp.Controller))
		if vpController == "" {
			vpController = "st7789"
		}
		drv, ok := driver.Get(vpController)
		if !ok {
			log.Printf("display: unrecognized controller %q; registered drivers: %v", vp.Controller, driver.Names())
			cleanup()
			return nil, func() {}, "", fmt.Errorf("unsupported display controller %q", vp.Controller)
		}

		if isI2COnly(drv) {
			// I2C path for this virtual screen.
			target, bus, controller, err := initI2CDisplayTarget(vp.I2CBus, vp.Controller, vp.Config)
			if err != nil {
				cleanup()
				return nil, func() {}, "", err
			}
			closers = append(closers, bus)
			targets = append(targets, indexedTarget{index: vp.Index, target: target})
			b := target.Bounds()
			w, h := b.Dx(), b.Dy()
			log.Printf("display[%d] %s ready (%dx%d %s via %s)", vp.Index, vp.Name, w, h, strings.ToUpper(controller), vp.I2CBus)
			if vp.Index == 0 {
				primaryController = controller
			}
			continue
		}

		// SPI path for this virtual screen.
		spiName := strings.TrimSpace(vp.SPI)
		if i == 0 && strings.TrimSpace(spiDevice) != "" {
			spiName = strings.TrimSpace(spiDevice)
		}
		if spiName == "" {
			cleanup()
			return nil, func() {}, "", fmt.Errorf("display: virtual panel %d has empty SPI device", vp.Index)
		}

		spec := panels.Screen{
			SPI:        spiName,
			Controller: vp.Controller,
			Config:     vp.Config,
			DCPin:      vp.DCPin,
			RSTPin:     vp.RSTPin,
			BusyPin:    vp.BusyPin,
			BLPin:      vp.BLPin,
		}

		target, port, controller, err := initDisplayTarget(spec)
		if err != nil {
			cleanup()
			return nil, func() {}, "", err
		}
		closers = append(closers, port)
		targets = append(targets, indexedTarget{index: vp.Index, target: target})
		b := target.Bounds()
		w, h := b.Dx(), b.Dy()
		log.Printf("display[%d] %s ready (%dx%d %s via %s)", vp.Index, vp.Name, w, h, strings.ToUpper(controller), spiName)
		if vp.Index == 0 {
			primaryController = controller
		}
	}
	if primaryController == "" {
		cleanup()
		return nil, func() {}, "", fmt.Errorf("display: virtual panel %q must define panel index 0", profile.Name)
	}

	positions, err := buildPositions(targets, profile)
	if err != nil {
		cleanup()
		return nil, func() {}, "", err
	}
	return positions, cleanup, primaryController, nil
}

// initDisplayTarget initialises a single SPI-connected display from a panels.Screen.
func initDisplayTarget(screen panels.Screen) (driver.DrawTarget, spi.PortCloser, string, error) {
	spiPort, err := spireg.Open(strings.TrimSpace(screen.SPI))
	if err != nil {
		log.Printf("display: spi device=%q open=error err=%q", screen.SPI, err.Error())
		return nil, nil, "", err
	}
	log.Printf("display: spi device=%q open=ok", screen.SPI)

	dc := pinOut(screen.DCPin)
	log.Printf("display: pin DC=%q type=DC found=%t", screen.DCPin, dc != nil)

	rst := pinOut(screen.RSTPin)
	log.Printf("display: pin RST=%q type=RST found=%t", screen.RSTPin, rst != nil)

	bl := pinOut(screen.BLPin)
	if screen.BLPin == "" {
		log.Printf("display: pin BL=%q type=BL found=false (unassigned)", screen.BLPin)
	} else {
		log.Printf("display: pin BL=%q type=BL found=%t", screen.BLPin, bl != nil)
	}

	busy := pinIn(screen.BusyPin)
	if screen.BusyPin == "" {
		log.Printf("display: pin BUSY=%q type=BUSY found=false (unassigned)", screen.BusyPin)
	} else {
		log.Printf("display: pin BUSY=%q type=BUSY found=%t", screen.BusyPin, busy != nil)
	}

	if dc == nil {
		log.Printf("display: required pin missing name=%q type=DC", screen.DCPin)
		_ = spiPort.Close()
		return nil, nil, "", errGPIO(screen.DCPin + " (DC) not found")
	}
	if rst == nil {
		log.Printf("display: required pin missing name=%q type=RST", screen.RSTPin)
		_ = spiPort.Close()
		return nil, nil, "", errGPIO(screen.RSTPin + " (RST) not found")
	}

	controller := strings.ToLower(strings.TrimSpace(screen.Controller))
	if controller == "" {
		controller = "st7789"
	}
	drv, ok := driver.Get(controller)
	if !ok {
		_ = spiPort.Close()
		return nil, nil, "", fmt.Errorf("unsupported display controller %q", screen.Controller)
	}

	log.Printf("display: driver factory id=%s type=SPI config=%dx%d madctl=0x%02X", controller, screen.Config.Width, screen.Config.Height, screen.Config.MADCTL)
	dev, err := drv.NewSPI(spiPort, dc, rst, bl, busy, screen.Config)
	if err != nil {
		log.Printf("display: driver factory error id=%s type=SPI err=%q", controller, err.Error())
		_ = spiPort.Close()
		return nil, nil, "", err
	}
	return dev, spiPort, controller, nil
}

// pinOut looks up a GPIO pin by name and returns it as gpio.PinOut.
// Returns nil when the pin is not found on this host.
func pinOut(name string) gpio.PinOut {
	p := gpioreg.ByName(name)
	if p == nil {
		return nil
	}
	if po, ok := p.(gpio.PinOut); ok {
		return po
	}
	return nil
}

// pinIn looks up a GPIO pin by name and returns it as gpio.PinIn.
// Returns nil when the pin is not found on this host (silently skipped).
func pinIn(name string) gpio.PinIn {
	p := gpioreg.ByName(name)
	if p == nil {
		return nil
	}
	if pi, ok := p.(gpio.PinIn); ok {
		return pi
	}
	return nil
}

type gpioError string

func errGPIO(msg string) error {
	return gpioError(msg)
}

func (e gpioError) Error() string {
	return "gpio: " + string(e)
}

// initI2CDisplayTarget initialises a single I2C-connected display.
func initI2CDisplayTarget(i2cBus string, controller string, config driver.DriverConfig) (driver.DrawTarget, i2c.BusCloser, string, error) {
	bus, err := i2creg.Open(i2cBus)
	if err != nil {
		log.Printf("display: i2c bus=%q open=error err=%q", i2cBus, err.Error())
		return nil, nil, "", fmt.Errorf("i2c bus %q: %w", i2cBus, err)
	}
	log.Printf("display: i2c bus=%q addr=0x%02X open=ok", i2cBus, config.I2CAddr)

	ctrl := strings.ToLower(strings.TrimSpace(controller))
	drv, ok := driver.Get(ctrl)
	if !ok {
		_ = bus.Close()
		return nil, nil, "", fmt.Errorf("unsupported controller %q", ctrl)
	}

	log.Printf("display: driver factory id=%s type=I2C config=%dx%d", ctrl, config.Width, config.Height)
	dev, err := drv.NewI2C(bus, config)
	if err != nil {
		log.Printf("display: driver factory error id=%s type=I2C err=%q", ctrl, err.Error())
		_ = bus.Close()
		return nil, nil, "", fmt.Errorf("controller %q: %w", ctrl, err)
	}
	return dev, bus, ctrl, nil
}

func resolveInputConfig(pins panels.InputPins) (input.Config, int) {
	cfg := input.Config{}
	detected := 0
	if pin := pinIn(pins.Key1); pin != nil {
		cfg.Key1 = pin
		detected++
	}
	if pin := pinIn(pins.Key2); pin != nil {
		cfg.Key2 = pin
		detected++
	}
	if pin := pinIn(pins.Key3); pin != nil {
		cfg.Key3 = pin
		detected++
	}
	if pin := pinIn(pins.JoyUp); pin != nil {
		cfg.JoyUp = pin
		detected++
	}
	if pin := pinIn(pins.JoyDown); pin != nil {
		cfg.JoyDown = pin
		detected++
	}
	if pin := pinIn(pins.JoyLeft); pin != nil {
		cfg.JoyLeft = pin
		detected++
	}
	if pin := pinIn(pins.JoyRight); pin != nil {
		cfg.JoyRight = pin
		detected++
	}
	if pin := pinIn(pins.JoyPressed); pin != nil {
		cfg.JoyPressed = pin
		detected++
	}

	// Log each declared input pin and whether it was resolved.
	type namedPin struct {
		name  string
		value string
		found bool
	}
	declared := []namedPin{
		{"Key1", pins.Key1, cfg.Key1 != nil},
		{"Key2", pins.Key2, cfg.Key2 != nil},
		{"Key3", pins.Key3, cfg.Key3 != nil},
		{"JoyUp", pins.JoyUp, cfg.JoyUp != nil},
		{"JoyDown", pins.JoyDown, cfg.JoyDown != nil},
		{"JoyLeft", pins.JoyLeft, cfg.JoyLeft != nil},
		{"JoyRight", pins.JoyRight, cfg.JoyRight != nil},
		{"JoyPressed", pins.JoyPressed, cfg.JoyPressed != nil},
	}
	var parts []string
	for _, d := range declared {
		if d.value != "" {
			parts = append(parts, fmt.Sprintf("%s=%s resolved=%t", d.name, d.value, d.found))
		}
	}
	if len(parts) > 0 {
		log.Printf("display: input pins %s", strings.Join(parts, " "))
	}

	return cfg, detected
}

// buildPositions converts the collected indexed targets into []region.ScreenPosition.
// Screens are placed left-to-right by ascending index unless the panel definition
// provides explicit XPosition/YPosition values. Screen names resolve from
// the panel's Virtual entries, falling back to the panel name for single-screen panels
// or "screen0", "screen1" etc. for multi-screen panels with unnamed entries.
//
// Returns an error if no screen has positive-dimension bounds.
func buildPositions(targets []indexedTarget, profile panels.Definition) ([]region.ScreenPosition, error) {
	if len(targets) == 0 {
		return nil, fmt.Errorf("display: no screens available for position layout")
	}

	// Sort targets by ascending index for deterministic ordering.
	sort.Slice(targets, func(i, j int) bool {
		return targets[i].index < targets[j].index
	})

	// Build a lookup from the panel's Virtual screen definitions.
	type screenMeta struct {
		name     string
		xPos     int
		yPos     int
		rotation int
		mirrorX  bool
		ppi      float64
	}
	metaByIndex := make(map[int]screenMeta, len(profile.Virtual))
	hasExplicitPositions := false
	for _, scr := range profile.Virtual {
		meta := screenMeta{name: scr.Name, xPos: scr.XPosition, yPos: scr.YPosition, rotation: scr.Rotation, mirrorX: scr.MirrorX, ppi: scr.PPI}
		metaByIndex[scr.Index] = meta
		if scr.XPosition > 0 || scr.YPosition > 0 {
			hasExplicitPositions = true
		}
	}

	positions := make([]region.ScreenPosition, 0, len(targets))
	hasPositiveBounds := false

	// Track cumulative X offset for auto-positioning (when no explicit positions).
	xOffset := 0
	for _, it := range targets {
		b := it.target.Bounds()
		w, h := b.Dx(), b.Dy()
		meta := metaByIndex[it.index]

		// When the screen is rotated 90° or 270°, the logical (VD) dimensions
		// are the transpose of the hardware dimensions.
		logW, logH := w, h
		if meta.rotation == 90 || meta.rotation == 270 {
			logW, logH = h, w
		}

		// Compute the screen's position within the virtual display coordinate space.
		var bounds image.Rectangle
		if hasExplicitPositions {
			bounds = image.Rect(meta.xPos, meta.yPos, meta.xPos+logW, meta.yPos+logH)
		} else {
			bounds = image.Rect(xOffset, 0, xOffset+logW, logH)
			xOffset += logW
		}

		if logW > 0 && logH > 0 {
			hasPositiveBounds = true
		}

		// Determine screen name: use virtual screen name if available,
		// otherwise fall back to the panel name or a generated name.
		name := meta.name
		if name == "" {
			if len(targets) == 1 {
				name = profile.Name
			} else {
				name = fmt.Sprintf("screen%d", it.index)
			}
		}

		sp := region.ScreenPosition{
			Index:    it.index,
			Name:     name,
			Bounds:   bounds,
			Target:   it.target,
			Rotation: meta.rotation,
			MirrorX:  meta.mirrorX,
			PPI:      meta.ppi,
		}
		if hp, ok := it.target.(textlayout.TextHintProvider); ok {
			sp.HintProvider = hp.TextHints
		}
		positions = append(positions, sp)
	}

	if !hasPositiveBounds {
		return nil, fmt.Errorf("display: no screen has positive-dimension bounds")
	}

	return positions, nil
}

// isEPaperController reports whether the named controller is an e-paper/e-ink display.
func isEPaperController(controller string) bool {
	drv, ok := driver.Get(controller)
	if !ok {
		return false
	}
	return drv.IsEPaper
}

// clearScreensOnExit blanks all non-e-paper screens to prevent burn-in when
// the daemon exits. E-paper displays are skipped because they retain their
// image without power and a refresh cycle would be unnecessary.
func clearScreensOnExit(screens []region.ScreenPosition, profile panels.Definition) {
	// Build a map from screen index to controller name.
	controllerByIndex := make(map[int]string, len(screens))
	if len(profile.Virtual) == 0 {
		// Single-screen panel: all screens use the primary controller.
		ctrl := strings.ToLower(strings.TrimSpace(profile.Controller))
		if ctrl == "" {
			ctrl = "st7789"
		}
		for _, s := range screens {
			controllerByIndex[s.Index] = ctrl
		}
	} else {
		for _, vp := range profile.Virtual {
			ctrl := strings.ToLower(strings.TrimSpace(vp.Controller))
			if ctrl == "" {
				ctrl = "st7789"
			}
			controllerByIndex[vp.Index] = ctrl
		}
	}

	for _, s := range screens {
		ctrl := controllerByIndex[s.Index]
		if isEPaperController(ctrl) {
			log.Printf("display: skipping clear for e-paper screen %q (index %d, controller %s)", s.Name, s.Index, ctrl)
			continue
		}
		b := s.Target.Bounds()
		black := image.NewRGBA(b)
		// image.NewRGBA zero-initializes all pixels to transparent black (0,0,0,0).
		// For display purposes this renders as black on all supported controllers.
		if err := s.Target.DrawImage(black); err != nil {
			log.Printf("display: failed to clear screen %q on exit: %v", s.Name, err)
		} else {
			log.Printf("display: cleared screen %q (index %d) on exit", s.Name, s.Index)
		}
	}
}

// runDisplayWithRegions initialises the display using the region-based render
// system. It constructs PanelActivationConfig, calls
// ActivatePanel, wires the RegionRenderer, InputDispatcher, TickRateResolver,
// and blocks on the RenderLoop until stopCh is closed.
func runDisplayWithRegions(
	stopCh <-chan struct{},
	profile panels.Definition,
	enableInput bool,
	spiDevice string,
	modeState *coordinator.State,
	configPPI float64,
) error {
	// ── Initialise display targets ────────────────────────────────────────
	screens, cleanupDisplays, _, err := initDisplayTargets(profile, spiDevice)
	if err != nil {
		return err
	}
	defer cleanupDisplays()

	// ── Resolve available modes ───────────────────────────────────────────
	policy := catalog.StandardPolicy{}
	profile.InputEnabled = enableInput
	for i := range profile.Virtual {
		profile.Virtual[i].InputEnabled = enableInput
	}
	displays := panels.Displays(profile, policy)
	var availModes []string
	var defaultMode string
	screenModes := map[string]string{}
	for _, d := range displays {
		if d.Index == 0 {
			availModes = d.Modes
			defaultMode = d.Default
		}
		screenModes[d.Name] = d.Default
	}
	if defaultMode == "" && len(availModes) > 0 {
		defaultMode = availModes[0]
	}

	// ── Resolve input ─────────────────────────────────────────────────────
	var events <-chan input.Event
	inputEnabled := false
	if enableInput && profile.Inputs.Any() {
		cfg, detected := resolveInputConfig(profile.Inputs)
		if detected == 0 {
			log.Printf("input: no input pins detected for panel %q; switching to passive mode", profile.Name)
		} else {
			inputHandler, err := input.New(cfg, 20*time.Millisecond)
			if err != nil {
				return err
			}
			inputHandler.Start()
			defer inputHandler.Stop()
			events = inputHandler.Events()
			inputEnabled = true
			log.Printf("input: %d pin(s) initialized (KEY1, KEY2, KEY3, joystick)", detected)
		}
	} else if enableInput {
		log.Printf("input: disabled (no pins configured for this panel)")
	} else {
		log.Printf("input: disabled by user flag (-noinput)")
	}

	// ── Construct PanelActivationConfig ───────────────────────────────────
	// Log config PPI override if active.
	if configPPI > 0 {
		log.Printf("display: config PPI override active: %.1f", configPPI)
	}

	config := region.PanelActivationConfig{
		Screens:       screens,
		Layout:        nil, // nil = use default generation (single region covering full bounds)
		DefaultMode:   defaultMode,
		InputEnabled:  inputEnabled,
		AvailModes:    availModes,
		ScreenModes:   screenModes,
		ModeValidator: displaymodes.IsKnown,
		Events:        events,
		PanelProduct:  strings.TrimSpace(profile.Name),
		PanelPPI:      profile.PPI,
		ConfigPPI:     configPPI,
	}

	// ── Activate panel ────────────────────────────────────────────────────
	activation, err := region.ActivatePanel(config)
	if err != nil {
		return fmt.Errorf("region activation: %w", err)
	}

	// ── Wire ModeFactory on all allocated regions ────────────────────────
	// This bridges the region package (which defines ModeFactory) with the
	// displaymodes package (which provides GetInstance). Without this,
	// Region.SetMode takes the legacy path and never constructs ModeInstances,
	// resulting in nil Instance() and black screens.
	for _, r := range activation.RegionManager.Regions() {
		r.SetModeFactory(func(id string, hints textlayout.TextHints) (region.ModeInstance, bool) {
			inst, ok := displaymodes.GetInstance(id)
			if !ok {
				return nil, false
			}
			return inst, true
		})
	}

	// ── Log display readiness ─────────────────────────────────────────────
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("DISPLAY INITIALIZED (region system)")
	log.Printf("  Panel:       %s", profile.Name)
	log.Printf("  Controller:  %s (%dx%d)", profile.Controller, profile.Config.Width, profile.Config.Height)
	log.Printf("  Screens:     %d", len(screens))
	if inputEnabled {
		log.Printf("  Input:       active (interactive mode)")
	} else {
		log.Printf("  Input:       disabled (passive mode)")
	}
	log.Printf("  Default mode: %s", defaultMode)
	for _, r := range activation.RegionManager.Regions() {
		hints := r.TextHints()
		log.Printf("  Surface %q: PPI=%.1f (%dx%d)", r.Name(), hints.PPI, hints.PixelWidth, hints.PixelHeight)
	}
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// ── Propagate panel TextHints to mode packages for fitness notes ──────
	// Mode command handlers (clock, thermal) use TextHints in their PostApply
	// hooks to generate fitness notes when a style is set via cyberhudctl.
	// Resolve hints from the primary screen (index 0) and distribute.
	if len(screens) > 0 {
		primaryHints := textlayout.ResolveTextHints(screens[0].Target, screens[0].HintProvider)
		modehints.PropagateHints(primaryHints)
		log.Printf("display: panel hints propagated (%dx%d) for fitness notes",
			primaryHints.PixelWidth, primaryHints.PixelHeight)
	}

	// ── Create RegionRenderer ─────────────────────────────────────────────
	var warnings []string
	if notices := panels.PinNotices(profile); len(notices) > 0 {
		for _, n := range notices {
			log.Printf("pin info: %s", n)
		}
		warnings = notices
	}
	displaymodes.WireWarnings(warnings)

	inputMapper := ui.BuildInputMapper(
		strings.TrimSpace(profile.Inputs.Key1) != "",
		strings.TrimSpace(profile.Inputs.Key2) != "",
		strings.TrimSpace(profile.Inputs.Key3) != "",
		strings.TrimSpace(profile.Inputs.JoyPressed) != "",
	)

	renderer := ui.NewRegionRenderer(
		profile.Monochrome,
		warnings,
		inputMapper,
		ui.WithRendererModeState(modeState),
	)

	// ── Create RegionInputDispatcher ──────────────────────────────────────
	dispatcher := ui.NewRegionInputDispatcher(renderer, inputMapper)

	// ── Create TickRateResolver ───────────────────────────────────────────
	tickResolver := &region.DefaultTickRateResolver{}

	// ── Configure RenderLoop ──────────────────────────────────────────────
	activation.RenderLoop.Apply(
		region.WithRenderer(renderer),
		region.WithInputDispatcher(dispatcher),
		region.WithTickRateResolver(tickResolver),
	)

	// Wire ModeSwitch for console mode-change commands.
	// The catalog.State is shared with the console server and handles mode
	// tracking. When catalog.State.Set is called by the console, the
	// RegionRenderer.syncMode() detects the change on the next render tick
	// and calls Region.SetMode to construct a fresh mode instance.

	// ── Run render loop (blocking until stop) ─────────────────────────────
	go func() {
		<-stopCh
		activation.RenderLoop.Stop()
	}()

	activation.RenderLoop.Run()

	// ── Clear non-e-paper screens on exit (burn-in prevention) ────────────
	clearScreensOnExit(screens, profile)

	return nil
}
