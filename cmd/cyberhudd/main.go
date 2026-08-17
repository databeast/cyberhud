// Command cyberhudd is a background daemon that drives SPI and I2C displays
// attached to a Raspberry Pi (any model with a 40-pin GPIO header).
//
// The default panel is the Waveshare 1.3-inch LCD HAT (240×240 ST7789) with
// buttons and a five-way joystick, but other panels can be selected at startup.
//
// A Unix domain socket at /run/cyberhudd/console.sock provides programmatic
// access using a simple line-oriented protocol (see internal/runtime/console).
//
// Usage:
//
//	cyberhudd [flags]
//
// Flags:
//
//	-socket   path to the Unix domain socket  (default /run/cyberhudd/console.sock)
//	-i2c      comma-separated I2C bus paths   (default /dev/i2c-1)
//	-scan     I2C re-scan interval            (default 2s)
//	-nodisplay        skip display/input init             (useful for headless testing)
//	-noinput          disable GPIO buttons/joystick       (useful for display-only panels)
//	-panel            select a built-in display panel
package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/coordinator"
	displaymodes "github.com/databeast/cyberhud/display/modes"
	cyclemode "github.com/databeast/cyberhud/display/modes/cycle"
	gpiomode "github.com/databeast/cyberhud/display/modes/gpio"
	gpioControlSource "github.com/databeast/cyberhud/display/modes/gpio_control/source"
	stemmapkg "github.com/databeast/cyberhud/display/modes/stemma"
	"github.com/databeast/cyberhud/display/modes/stemma/source"
	gpiomgr "github.com/databeast/cyberhud/hardware/gpio"
	"github.com/databeast/cyberhud/hardware/panels"
	_ "github.com/databeast/cyberhud/hardware/panels/all"
	"github.com/databeast/cyberhud/runtime/console"
	"periph.io/x/host/v3"
)

func main() {
	socketPath := flag.String("socket", "/run/cyberhudd/console.sock",
		"path to the Unix domain socket")
	i2cBuses := flag.String("i2c", "/dev/i2c-1",
		"comma-separated I2C bus paths to scan for STEMMA QT / QWIIC devices")
	scanInterval := flag.Duration("scan", 2*time.Second,
		"interval between STEMMA QT / QWIIC I2C bus re-scans")
	noDisplay := flag.Bool("nodisplay", false,
		"disable display and input (headless / testing mode)")
	configPath := flag.String("config", "",
		"optional path to JSON configuration file (CLI flags take precedence)")
	noInput := flag.Bool("noinput", false,
		"disable GPIO input even when the selected panel provides buttons")
	displayProfileName := flag.String("panel", "waveshare-1.3hat",
		"built-in display panel to use ("+panels.Supported()+")")
	spiDevice := flag.String("spi", "", "optional SPI device name for display (e.g. SPI0.0)")
	displayWidth := flag.Int("display-width", -1,
		"override display width in pixels for the selected panel")
	displayHeight := flag.Int("display-height", -1,
		"override display height in pixels for the selected panel")
	displayMADCTL := flag.String("display-madctl", "",
		"override MADCTL byte (decimal or hex, e.g. 0x20)")
	displayRotate := flag.String("display-rotate", "",
		"rotate display 180 degrees (true/false, overrides panel default)")
	displayXOffset := flag.Int("display-x-offset", -1,
		"override X offset for panels that require a shifted framebuffer")
	displayYOffset := flag.Int("display-y-offset", -1,
		"override Y offset for panels that require a shifted framebuffer")
	displayDC := flag.String("display-dc", "",
		"override the display data/command GPIO pin name")
	displayRST := flag.String("display-rst", "",
		"override the display reset GPIO pin name")
	displayBL := flag.String("display-bl", "",
		"override the display backlight GPIO pin name (use 'none' to disable backlight control)")
	displayBusy := flag.String("display-busy", "",
		"override the e-ink display busy GPIO pin name")
	displayPPI := flag.Float64("display-ppi", 0,
		"Override PPI for the display panel (0 = use panel/config default)")
	inputKey1 := flag.String("input-key1", "", "override GPIO pin name for KEY1")
	inputKey2 := flag.String("input-key2", "", "override GPIO pin name for KEY2")
	inputKey3 := flag.String("input-key3", "", "override GPIO pin name for KEY3")
	inputUp := flag.String("input-up", "", "override GPIO pin name for D-pad up")
	inputDown := flag.String("input-down", "", "override GPIO pin name for D-pad down")
	inputLeft := flag.String("input-left", "", "override GPIO pin name for D-pad left")
	inputRight := flag.String("input-right", "", "override GPIO pin name for D-pad right")
	inputPress := flag.String("input-press", "", "override GPIO pin name for D-pad press")
	flag.Parse()
	seen := visitedFlags()

	cfg, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	if err := mergeConfig(
		cfg,
		seen,
		socketPath,
		i2cBuses,
		scanInterval,
		noDisplay,
		noInput,
		displayProfileName,
		displayWidth,
		displayHeight,
		displayMADCTL,
		displayRotate,
		displayXOffset,
		displayYOffset,
		displayDC,
		displayRST,
		displayBL,
		displayBusy,
		inputKey1,
		inputKey2,
		inputKey3,
		inputUp,
		inputDown,
		inputLeft,
		inputRight,
		inputPress,
	); err != nil {
		log.Fatalf("config: %v", err)
	}

	// ── PPI override merge ────────────────────────────────────────────────
	// CLI flag takes precedence over config file. Zero means "no override."
	var configPPI float64
	if seen["display-ppi"] {
		configPPI = *displayPPI
	} else if cfg != nil && cfg.Display.PPI != nil && *cfg.Display.PPI > 0 {
		configPPI = *cfg.Display.PPI
	}

	// Build runtime config snapshot for the "config dump" protocol command.
	runtimeConfigFn := buildRuntimeConfig(
		*socketPath,
		*i2cBuses,
		*scanInterval,
		*noDisplay,
		*noInput,
		*displayProfileName,
		*displayWidth,
		*displayHeight,
		*displayMADCTL,
		*displayXOffset,
		*displayYOffset,
		*displayDC,
		*displayRST,
		*displayBL,
		*displayBusy,
		*inputKey1,
		*inputKey2,
		*inputKey3,
		*inputUp,
		*inputDown,
		*inputLeft,
		*inputRight,
		*inputPress,
	)
	runtimeConfigJSON := func() string {
		cfg := runtimeConfigFn()
		data, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			return ""
		}
		return string(data)
	}

	profile, err := panels.Resolve(
		*displayProfileName,
		panels.Overrides{
			Width:              *displayWidth,
			Height:             *displayHeight,
			MADCTL:             *displayMADCTL,
			Rotate:             *displayRotate,
			XOffset:            *displayXOffset,
			YOffset:            *displayYOffset,
			Orientation:        orientationFromConfig(cfg),
			ScreenOrientations: screenOrientationsFromConfig(cfg),
			DCPin:              *displayDC,
			RSTPin:             *displayRST,
			BLPin:              *displayBL,
			BusyPin:            *displayBusy,
			InputKey1:          *inputKey1,
			InputKey2:          *inputKey2,
			InputKey3:          *inputKey3,
			InputUp:            *inputUp,
			InputDown:          *inputDown,
			InputLeft:          *inputLeft,
			InputRight:         *inputRight,
			InputPress:         *inputPress,
		},
	)
	if err != nil {
		log.Fatalf("panel: %v", err)
	}
	modeState := coordinator.NewState()
	validateModeRegistrations()
	cyclemode.WireState(modeState)
	configureDisplayModes(modeState, profile, !*noInput)

	log.SetFlags(log.Ldate | log.Ltime | log.Lmsgprefix)
	log.SetPrefix("[cyberhudd] ")

	// Log startup configuration for diagnostics
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("CYBERHUD DAEMON STARTING")
	log.Printf("  Panel:       %s (%s, %dx%d)", profile.Name, profile.Controller, profile.Config.Width, profile.Config.Height)
	if !*noInput && profile.Inputs.Any() {
		log.Printf("  Input:       enabled")
	} else if !*noInput {
		log.Printf("  Input:       no pins configured (using passive mode)")
	} else {
		log.Printf("  Input:       disabled (-noinput flag)")
	}
	if *noDisplay {
		log.Printf("  Display:     disabled (headless mode)")
	}
	if *configPath != "" {
		log.Printf("  Config file: %s", *configPath)
	}
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// ── periph host initialisation ────────────────────────────────────────
	if _, err := host.Init(); err != nil {
		log.Fatalf("periph host init: %v", err)
	}

	// ── STEMMA QT / QWIIC scanner configuration ──────────────────────────
	// Provide I2C bus names and scan interval to the stemma package so that
	// Activate() can create the scanner when the stemma mode becomes active.
	buses := strings.Split(*i2cBuses, ",")
	stemmapkg.SetScannerConfig(buses, *scanInterval)
	log.Printf("STEMMA QT / QWIIC scanner configured (buses: %s, interval: %s)",
		*i2cBuses, *scanInterval)

	// ── GPIO manager ──────────────────────────────────────────────────────
	gpioManager := gpiomgr.New()
	gpiomode.SetGPIOManager(gpioManager)
	gpioControlSource.SetGPIOControlManager(gpioManager)
	log.Printf("GPIO manager started")

	// ── Console socket ────────────────────────────────────────────────────
	if err := os.MkdirAll(filepath.Dir(*socketPath), 0o755); err != nil {
		log.Fatalf("create socket directory: %v", err)
	}

	// ── PolicyStore (mode policy persistence) ─────────────────────────────
	policyStore := NewPolicyStore()
	if cfg != nil && len(cfg.Policies) > 0 {
		policyStore.LoadFromConfig(cfg)
		for modeID := range cfg.Policies {
			if !displaymodes.IsKnownInstance(modeID) {
				log.Printf("WARNING: config policies: unrecognized mode %q (ignored)", modeID)
			}
		}
		log.Printf("loaded %d mode policies from config", len(cfg.Policies))
	}

	// Register the PolicyStore as the centralized policy change listener.
	// Whenever any mode's CmdHandler applies key=value changes, the updated
	// policy snapshot is automatically captured in the PolicyStore.
	catalog.OnPolicyChange = policyStore.Set

	consoleServer := console.New(*socketPath, source.GlobalScanner, gpioManager, func() string {
		return panels.BuildPinReport(profile)
	}, modeState, runtimeConfigJSON, policyStore, *configPath)
	if err := consoleServer.Start(); err != nil {
		log.Fatalf("console server: %v", err)
	}
	defer consoleServer.Stop()
	log.Printf("console socket: %s", consoleServer.SocketPath())

	// ── Restore saved policies for initial active modes ───────────────────
	consoleServer.RestoreInitialPolicies()

	// ── Display + input (skipped in headless / testing mode) ──────────────
	stopUI := make(chan struct{})

	if !*noDisplay {
		go func() {
			if err := runDisplayWithRegions(stopUI, profile, !*noInput, *spiDevice, modeState, configPPI); err != nil {
				log.Printf("display error: %v", err)
			}
		}()
	} else {
		log.Printf("display disabled (-nodisplay flag)")
	}

	// ── Signal handling ───────────────────────────────────────────────────
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh
	log.Printf("received signal %s – shutting down", sig)

	close(stopUI)
	log.Printf("shutdown complete")
}
