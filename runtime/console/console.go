// Package console provides a Unix-domain-socket server that exposes a simple
// line-oriented protocol for querying and controlling the cyberhud daemon.
//
// Each client connection is served in its own goroutine.  The protocol is
// line-based UTF-8:
//
//	Request  → <verb> [<sub-command>] [<args>]\n
//	Response → OK <data>\n  or  ERR <message>\n
//
// Commands:
//
//	status                                  – one-line daemon summary
//	gpio status                             – list all GPIO pin states
//	gpio pins                               – report pin usage and connector conflicts
//	gpio set <n> <0|1>                      – drive GPIO pin n to level 0 or 1
//	gpio in <n>                             – configure pin n as input and read its level
//	stemma status                           – list all detected STEMMA QT / QWIIC devices
//	display status                          – list regions with current mode
//	display modes                           – list available modes per region
//	display set <region> <mode> [key=value ...] – set region mode with optional config
//	display config <region> [key=value ...]  – query or configure the active mode
//	display next <region>                    – switch region to next mode
//	display prev <region>                    – switch region to previous mode
//	display <mode> [args...]                – mode-specific data commands (ticker, image, etc.)
//	config dump                             – dump runtime configuration as JSON
//	help [modes]                            – show command help
//	quit / exit                             – close the connection
package console

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	stemmapkg "github.com/databeast/cyberhud/display/modes/stemma/source"
	"periph.io/x/conn/v3/gpio"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/coordinator"
	_ "github.com/databeast/cyberhud/display/modes"
	"github.com/databeast/cyberhud/display/regionid"
	gpiopkg "github.com/databeast/cyberhud/hardware/gpio"
)

// PolicyStoreAccess is an interface for accessing the daemon's PolicyStore.
// The console server uses this to serve policy dump and freeze commands.
type PolicyStoreAccess interface {
	Get(modeID string) json.RawMessage
	Set(modeID string, data json.RawMessage)
	All() map[string]json.RawMessage
}

// Main command verbs
const (
	cmdQuit    = "quit"
	cmdExit    = "exit"
	cmdStatus  = "status"
	cmdGPIO    = "gpio"
	cmdStemma  = "stemma"
	cmdDisplay = "display"
	cmdHelp    = "help"
	cmdConfig  = "config"
	cmdPolicy  = "policy"
	cmdFreeze  = "freeze"
)

// GPIO subcommands
const (
	gpioStatus = "status"
	gpioSet    = "set"
	gpioIn     = "in"
	gpioPins   = "pins"
)

// Stemma subcommands
const (
	stemmaStatus = "status"
)

// Help subcommands
const (
	helpModes = "modes"
)

// Config subcommands
const (
	configDump = "dump"
)

// Policy subcommands
const (
	policyDump = "dump"
)

// Freeze subcommands
const (
	freezePolicy = "policy"
)

// Response prefixes
const (
	respOK  = "OK"
	respErr = "ERR"
)

// Response messages
const (
	msgGreeting             = "OK cyberhud daemon ready"
	msgBye                  = "OK bye"
	msgNone                 = "OK (none)"
	msgEmptyCommand         = "ERR empty command"
	msgUnknownVerb          = "ERR unknown verb"
	msgPinReportUnavailable = "ERR pin report unavailable"
	msgDisplayUnavailable   = "ERR display control unavailable"
	msgLevelInvalid         = "ERR level must be 0 or 1"
	msgConfigUnavailable    = "ERR config snapshot unavailable"
	msgConfigUsage          = "ERR usage: config dump"
	msgHelpUsage            = "ERR usage: help [modes]"
	msgGPIOUsage            = "ERR usage: gpio <status|set|in|pins> ..."
	msgGPIOSetUsage         = "ERR usage: gpio set <pin> <0|1>"
	msgGPIOInUsage          = "ERR usage: gpio in <pin>"
	msgStemmaUsage          = "ERR usage: stemma status"
	msgPolicyUsage          = "ERR usage: policy dump"
	msgPolicyStoreNil       = "ERR policy store not configured"
	msgDisplayUsage         = "ERR usage: display <status|modes|set|config|next|prev|<mode>> ..."
	msgDisplayModesUsage    = "ERR usage: display modes"
	msgDisplaySetUsage      = "ERR usage: display set <region> <mode> [key=value ...]"
	msgDisplayNextUsage     = "ERR usage: display next <region>"
	msgDisplayPrevUsage     = "ERR usage: display prev <region>"
	msgPresent              = "present"
	msgAbsent               = "absent"
)

// Server listens on a Unix domain socket and answers queries about the
// cyberhud hardware state.
type Server struct {
	sockPath      string
	scanner       func() *stemmapkg.Scanner
	gpiomgr       *gpiopkg.Manager
	pinReport     func() string
	modes         *coordinator.State
	runtimeConfig func() string
	policyStore   PolicyStoreAccess
	configPath    string
	listener      net.Listener
	wg            sync.WaitGroup
	stopOnce      sync.Once
}

// New creates a Server.  sockPath is the path of the Unix socket to create.
// scanner is an accessor function that returns the current STEMMA scanner
// instance, or nil when the stemma mode is inactive.
// runtimeConfig is an optional closure that returns the current runtime
// configuration as a pre-serialized JSON string; it may be nil.
// policyStore provides access to the daemon's PolicyStore for policy dump and
// freeze commands; it may be nil if policy persistence is not configured.
// configPath is the path to the daemon's config file, used for freeze writes.
func New(sockPath string, scanner func() *stemmapkg.Scanner, gpiomgr *gpiopkg.Manager, pinReport func() string, modes *coordinator.State, runtimeConfig func() string, policyStore PolicyStoreAccess, configPath string) *Server {
	return &Server{
		sockPath:      sockPath,
		scanner:       scanner,
		gpiomgr:       gpiomgr,
		pinReport:     pinReport,
		modes:         modes,
		runtimeConfig: runtimeConfig,
		policyStore:   policyStore,
		configPath:    configPath,
	}
}

// Start begins listening on the Unix socket.
func (s *Server) Start() error {
	// Remove any stale socket file.
	_ = os.Remove(s.sockPath)

	l, err := net.Listen("unix", s.sockPath)
	if err != nil {
		return fmt.Errorf("console: listen %q: %w", s.sockPath, err)
	}
	if err := configureSocketPermissions(s.sockPath); err != nil {
		_ = l.Close()
		_ = os.Remove(s.sockPath)
		return err
	}
	s.listener = l

	s.wg.Add(1)
	go s.accept()
	return nil
}

func configureSocketPermissions(path string) error {
	if err := os.Chmod(path, 0o600); err != nil {
		return fmt.Errorf("console: chmod socket: %w", err)
	}
	grp, err := user.LookupGroup("cyberhud")
	if err != nil {
		return nil
	}
	gid, err := strconv.Atoi(strings.TrimSpace(grp.Gid))
	if err != nil {
		return nil
	}
	if err := os.Chmod(path, 0o660); err != nil {
		return fmt.Errorf("console: chmod socket group mode: %w", err)
	}
	if err := os.Chown(path, -1, gid); err != nil {
		return fmt.Errorf("console: chown socket group: %w", err)
	}
	dir := filepath.Dir(path)
	if err := os.Chmod(dir, 0o750); err != nil {
		return fmt.Errorf("console: chmod socket dir: %w", err)
	}
	if err := os.Chown(dir, -1, gid); err != nil {
		return fmt.Errorf("console: chown socket dir group: %w", err)
	}
	return nil
}

// Stop shuts down the server.
func (s *Server) Stop() {
	s.stopOnce.Do(func() {
		if s.listener != nil {
			_ = s.listener.Close()
		}
	})
	s.wg.Wait()
}

// SocketPath returns the path of the Unix socket.
func (s *Server) SocketPath() string {
	return s.sockPath
}

func (s *Server) accept() {
	defer s.wg.Done()
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			continue
		}
		s.wg.Add(1)
		go s.serve(conn)
	}
}

func (s *Server) serve(conn net.Conn) {
	defer s.wg.Done()
	defer conn.Close()

	fmt.Fprintf(conn, "OK cyberhud daemon ready\n")

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		resp := s.dispatch(line)
		fmt.Fprintf(conn, "%s\n", resp)
	}
}

func (s *Server) dispatch(line string) string {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return msgEmptyCommand
	}
	verb := strings.ToLower(parts[0])

	switch verb {
	case cmdQuit, cmdExit:
		return msgBye

	case cmdStatus:
		scanner := s.scanner()
		var devCount int
		if scanner != nil {
			devCount = len(scanner.PresentDevices())
		}
		pins := s.gpiomgr.Snapshot()
		if s.modes != nil {
			return fmt.Sprintf("OK stemma_devices=%d gpio_pins=%d display_regions=%d", devCount, len(pins), len(s.modes.Status()))
		}
		return fmt.Sprintf("OK stemma_devices=%d gpio_pins=%d", devCount, len(pins))

	case cmdGPIO:
		return s.handleGPIOCmd(parts[1:])

	case cmdStemma:
		return s.handleStemmaCmd(parts[1:])

	case cmdDisplay:
		return s.handleDisplay(parts[1:])

	case cmdPolicy:
		return s.handlePolicy(parts[1:])

	case cmdFreeze:
		return s.handleFreeze(parts[1:])

	case cmdHelp:
		return s.handleHelp(parts[1:])

	case cmdConfig:
		return s.handleConfig(parts[1:])

	default:
		return msgUnknownVerb
	}
}

func (s *Server) handleHelp(args []string) string {
	if len(args) == 0 {
		return msgHelpUsage
	}
	sub := strings.ToLower(strings.TrimSpace(args[0]))
	switch sub {
	case helpModes:
		return s.listModeCommands()
	default:
		return msgHelpUsage
	}
}

func (s *Server) handleConfig(args []string) string {
	if len(args) < 1 {
		return msgConfigUsage
	}
	sub := strings.ToLower(args[0])
	switch sub {
	case configDump:
		return s.handleConfigDump()
	default:
		return msgConfigUsage
	}
}

func (s *Server) handleConfigDump() string {
	if s.runtimeConfig == nil {
		return msgConfigUnavailable
	}
	result := s.runtimeConfig()
	if result == "" {
		return msgConfigUnavailable
	}
	return "OK\n" + result
}

func (s *Server) handlePolicy(args []string) string {
	if len(args) < 1 {
		return msgPolicyUsage
	}
	sub := strings.ToLower(args[0])
	switch sub {
	case policyDump:
		return s.handlePolicyDump()
	default:
		return msgPolicyUsage
	}
}

func (s *Server) handlePolicyDump() string {
	if s.policyStore == nil {
		return msgPolicyStoreNil
	}
	all := s.policyStore.All()
	data, err := json.Marshal(all)
	if err != nil {
		return fmt.Sprintf("ERR failed to serialize policies: %v", err)
	}
	return "OK " + string(data)
}

func (s *Server) handleFreeze(args []string) string {
	// No args or first arg is "policy" → freeze policy
	if len(args) == 0 || strings.ToLower(args[0]) == freezePolicy {
		return s.handleFreezePolicy()
	}
	return "ERR usage: freeze [policy]"
}

func (s *Server) handleFreezePolicy() string {
	if s.configPath == "" {
		return "ERR config path not configured"
	}
	if s.policyStore == nil {
		return msgPolicyStoreNil
	}

	// Step 1: Snapshot all policies from the PolicyStore.
	policies := s.policyStore.All()

	// Step 2: Read existing config file (if it exists).
	var existing map[string]json.RawMessage
	data, err := os.ReadFile(s.configPath)
	if err != nil {
		if !os.IsNotExist(err) {
			if os.IsPermission(err) {
				return "ERR cannot read config: permission denied"
			}
			return fmt.Sprintf("ERR cannot read config: %v", err)
		}
		// File doesn't exist — start with an empty object.
		existing = make(map[string]json.RawMessage)
	} else {
		if err := json.Unmarshal(data, &existing); err != nil {
			return fmt.Sprintf("ERR cannot parse config: %v", err)
		}
	}

	// Step 3: Merge policies into the existing config.
	policiesJSON, err := json.Marshal(policies)
	if err != nil {
		return fmt.Sprintf("ERR failed to serialize policies: %v", err)
	}
	existing["policies"] = json.RawMessage(policiesJSON)

	// Step 4: Marshal back to JSON with indentation.
	output, err := json.MarshalIndent(existing, "", "  ")
	if err != nil {
		return fmt.Sprintf("ERR failed to serialize config: %v", err)
	}
	output = append(output, '\n')

	// Step 5: Create .bak backup of old file (if it exists).
	if _, statErr := os.Stat(s.configPath); statErr == nil {
		backupPath := s.configPath + ".bak"
		// Copy existing content to backup (don't rename, so atomic write is safe).
		if copyErr := copyFile(s.configPath, backupPath); copyErr != nil {
			// Non-fatal: log warning but continue with the write.
			_ = copyErr
		}
	}

	// Step 6: Atomic write — write to temp file, then rename.
	dir := filepath.Dir(s.configPath)
	tmp, err := os.CreateTemp(dir, ".cyberhud-cfg-*.tmp")
	if err != nil {
		if os.IsPermission(err) {
			return "ERR cannot write config: permission denied"
		}
		return fmt.Sprintf("ERR cannot write config: %v", err)
	}
	tmpName := tmp.Name()

	if _, err := tmp.Write(output); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		if os.IsPermission(err) {
			return "ERR cannot write config: permission denied"
		}
		return fmt.Sprintf("ERR cannot write config: %v", err)
	}
	if err := tmp.Chmod(0644); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		return fmt.Sprintf("ERR cannot write config: %v", err)
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpName)
		return fmt.Sprintf("ERR cannot write config: %v", err)
	}
	if err := os.Rename(tmpName, s.configPath); err != nil {
		os.Remove(tmpName)
		if os.IsPermission(err) {
			return "ERR cannot write config: permission denied"
		}
		return fmt.Sprintf("ERR cannot write config: %v", err)
	}

	return "OK"
}

// copyFile copies src to dst, creating or overwriting dst.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}

func (s *Server) handleDisplayPolicy(args []string) string {
	if len(args) < 1 {
		return "ERR usage: display policy <mode|region_id>"
	}
	arg := args[0]

	// First, check if the argument is a known mode ID.
	snap := catalog.Snapshotter(arg)
	if snap != nil {
		// It's a known mode — return its global policy snapshot.
		policy := snap.SnapshotPolicy()
		data, err := json.Marshal(policy)
		if err != nil {
			return fmt.Sprintf("ERR failed to serialize policy: %v", err)
		}
		return "OK " + string(data)
	}

	// Not a known mode — try to resolve as a region_id.
	region, errMsg := s.parseRegion(arg)
	if errMsg != "" {
		// Neither a known mode nor a valid region.
		return fmt.Sprintf("ERR unknown mode or region %q; available regions: %s", arg, s.availableRegions())
	}

	// Look up the current mode running on this region.
	currentMode := s.modes.CurrentMode(region)
	if currentMode == "" {
		return fmt.Sprintf("ERR no active mode on region %s", s.regionID(region))
	}

	// Get the mode's snapshotter and return mode name + policy.
	regionSnap := catalog.Snapshotter(currentMode)
	if regionSnap == nil {
		return fmt.Sprintf("OK mode=%s {}", currentMode)
	}
	policy := regionSnap.SnapshotPolicy()
	data, err := json.Marshal(policy)
	if err != nil {
		return fmt.Sprintf("ERR failed to serialize policy: %v", err)
	}
	return fmt.Sprintf("OK mode=%s %s", currentMode, string(data))
}

func (s *Server) listModeCommands() string {
	cmds := catalog.Commands()
	if len(cmds) == 0 {
		return msgNone
	}
	defs := catalog.Definitions()
	meta := make(map[string]catalog.Definition, len(defs))
	for _, def := range defs {
		meta[def.ID] = def
	}
	var sb strings.Builder
	sb.WriteString(respOK)
	for _, cmd := range cmds {
		fmt.Fprintf(&sb, "\n  verb=%s", cmd.Verb)
		if cmd.Usage != "" {
			fmt.Fprintf(&sb, " usage=%q", cmd.Usage)
		}
		if cmd.Summary != "" {
			fmt.Fprintf(&sb, " summary=%q", cmd.Summary)
		}
		if def, ok := meta[cmd.Verb]; ok {
			fmt.Fprintf(&sb, " mode_title=%q scope=%s", def.Title, def.Scope)
			for _, opt := range def.Options {
				fmt.Fprintf(&sb, "\n    option=%s type=%s default=%q summary=%q", opt.Key, opt.Type, opt.Default, opt.Summary)
			}
		}
	}
	return sb.String()
}

func (s *Server) handleDisplay(args []string) string {
	if s.modes == nil {
		return msgDisplayUnavailable
	}
	if len(args) < 1 {
		return msgDisplayUsage
	}
	sub := strings.ToLower(args[0])
	switch sub {
	case "status", "list":
		return s.listDisplay()
	case "regions":
		return s.listDisplayRegions()
	case "modes":
		if len(args) != 1 {
			return msgDisplayModesUsage
		}
		return s.listDisplayModes()
	case "set":
		if len(args) < 3 {
			return msgDisplaySetUsage
		}
		region, errMsg := s.parseRegion(args[1])
		if errMsg != "" {
			return errMsg
		}
		mode, err := s.modes.Set(region, args[2])
		if err != nil {
			return fmt.Sprintf("ERR unsupported mode %q for region %s", args[2], s.regionID(region))
		}
		// Restore saved policy for the newly activated mode.
		s.restoreModePolicy(mode)
		// If extra key=value args provided, forward them to the mode's catalog command handler.
		if len(args) > 3 {
			if def, ok := catalog.Command(mode); ok {
				result := def.Handle(args[3:])
				if strings.HasPrefix(result, "ERR ") {
					return result
				}
			}
		}
		return fmt.Sprintf("OK region=%s mode=%s", s.regionID(region), mode)
	case "config":
		if len(args) < 2 {
			return "ERR usage: display config <region> [key=value ...]"
		}
		region, errMsg := s.parseRegion(args[1])
		if errMsg != "" {
			return errMsg
		}
		// Find the current mode on this region and forward args to its command handler.
		current := s.modes.CurrentMode(region)
		if current == "" {
			return "ERR no active mode on region"
		}
		def, ok := catalog.Command(current)
		if !ok {
			return fmt.Sprintf("ERR mode %q has no configurable options", current)
		}
		return def.Handle(args[2:])
	case "next":
		if len(args) != 2 {
			return msgDisplayNextUsage
		}
		region, errMsg := s.parseRegion(args[1])
		if errMsg != "" {
			return errMsg
		}
		mode, err := s.modes.Next(region)
		if err != nil {
			return fmt.Sprintf("ERR %v", err)
		}
		// Restore saved policy for the newly activated mode.
		s.restoreModePolicy(mode)
		return fmt.Sprintf("OK region=%s mode=%s", s.regionID(region), mode)
	case "prev":
		if len(args) != 2 {
			return msgDisplayPrevUsage
		}
		region, errMsg := s.parseRegion(args[1])
		if errMsg != "" {
			return errMsg
		}
		mode, err := s.modes.Prev(region)
		if err != nil {
			return fmt.Sprintf("ERR %v", err)
		}
		// Restore saved policy for the newly activated mode.
		s.restoreModePolicy(mode)
		return fmt.Sprintf("OK region=%s mode=%s", s.regionID(region), mode)
	case "policy":
		return s.handleDisplayPolicy(args[1:])
	default:
		// Mode-specific commands: "display <mode-id> [args...]"
		// e.g., "display ticker set hello", "display image clear"
		if def, ok := catalog.Command(sub); ok {
			return def.Handle(args[1:])
		}
		return msgDisplayUsage
	}
}

func (s *Server) parseRegion(arg string) (int, string) {
	arg = strings.TrimSpace(arg)

	// Try bare integer first — resolves to coordinator region at that index.
	if idx, ok := regionid.ParseBareInt(arg); ok {
		if s.modes.HasRegion(idx) {
			return idx, ""
		}
		return 0, fmt.Sprintf("ERR unknown region %q; available: %s", arg, s.availableRegions())
	}

	// Try surface.index notation.
	id, err := regionid.Parse(strings.ToLower(arg))
	if err == nil {
		// Look up the region whose Name matches the surface name.
		for _, p := range s.modes.Status() {
			if strings.EqualFold(p.Name, id.Surface) {
				return p.Index, ""
			}
		}
		return 0, fmt.Sprintf("ERR unknown region %q; available: %s", arg, s.availableRegions())
	}

	// Neither bare integer nor valid region ID format.
	return 0, fmt.Sprintf("ERR unknown region %q; available: %s", arg, s.availableRegions())
}

// regionID returns the region identifier for a region index in <surface>.0 format.
func (s *Server) regionID(regionIndex int) string {
	ps, ok := s.modes.Region(regionIndex)
	if !ok {
		return fmt.Sprintf("%d.0", regionIndex)
	}
	if ps.Name == "" {
		return fmt.Sprintf("%d.0", regionIndex)
	}
	return ps.Name + ".0"
}

// availableRegions returns a comma-separated list of available regions in
// <surface>.<index> format (e.g., "main.0, left-aux.0, right-aux.0").
func (s *Server) availableRegions() string {
	regions := s.modes.Status()
	if len(regions) == 0 {
		return "(none)"
	}
	names := make([]string, 0, len(regions))
	for _, r := range regions {
		name := r.Name
		if name == "" {
			name = fmt.Sprintf("%d", r.Index)
		}
		names = append(names, fmt.Sprintf("%s.0", name))
	}
	return strings.Join(names, ", ")
}

// restoreModePolicy restores saved policy for a mode from the PolicyStore.
// It looks up the PolicySnapshotter for the mode and, if saved policy data
// exists, deserializes it and calls RestorePolicy (which applies normalization).
func (s *Server) restoreModePolicy(modeID string) {
	if s.policyStore == nil {
		return
	}
	data := s.policyStore.Get(modeID)
	if data == nil {
		return
	}
	snap := catalog.Snapshotter(modeID)
	if snap == nil {
		return
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		log.Printf("policy restore: failed to parse saved policy for %q: %v", modeID, err)
		return
	}
	if err := snap.RestorePolicy(parsed); err != nil {
		log.Printf("policy restore: failed to apply saved policy for %q: %v", modeID, err)
		return
	}
	log.Printf("policy restore: applied saved policy for mode %q", modeID)
}

// RestoreInitialPolicies restores saved policies for all currently active modes.
// Call this after the coordinator is configured and the PolicyStore is loaded.
func (s *Server) RestoreInitialPolicies() {
	if s.policyStore == nil || s.modes == nil {
		return
	}
	regions := s.modes.Status()
	restored := 0
	for _, r := range regions {
		if r.Current == "" {
			continue
		}
		data := s.policyStore.Get(r.Current)
		if data == nil {
			continue
		}
		snap := catalog.Snapshotter(r.Current)
		if snap == nil {
			continue
		}
		var parsed map[string]interface{}
		if err := json.Unmarshal(data, &parsed); err != nil {
			log.Printf("policy restore: failed to parse saved policy for %q: %v", r.Current, err)
			continue
		}
		if err := snap.RestorePolicy(parsed); err != nil {
			log.Printf("policy restore: failed to apply saved policy for %q: %v", r.Current, err)
			continue
		}
		restored++
	}
	if restored > 0 {
		log.Printf("policy restore: applied saved policies for %d initial mode(s)", restored)
	}
}

func (s *Server) listDisplay() string {
	regions := s.modes.Status()
	if len(regions) == 0 {
		return "OK (none)"
	}
	var sb strings.Builder
	sb.WriteString("OK")
	for _, r := range regions {
		name := r.Name
		if name == "" {
			name = fmt.Sprintf("%d", r.Index)
		}
		fmt.Fprintf(&sb, "\n  region=%s.0", name)
		if r.Name != "" {
			fmt.Fprintf(&sb, " name=%q", r.Name)
		}
		if r.Controller != "" {
			fmt.Fprintf(&sb, " controller=%s", r.Controller)
		}
		fmt.Fprintf(&sb, " mode=%s modes=%s", r.Current, strings.Join(r.Modes, ","))
	}
	return sb.String()
}

func (s *Server) listDisplayRegions() string {
	regions := s.modes.Status()
	if len(regions) == 0 {
		return "OK no regions configured"
	}
	var sb strings.Builder
	sb.WriteString("OK")
	for _, r := range regions {
		name := r.Name
		if name == "" {
			name = fmt.Sprintf("region%d", r.Index)
		}
		fmt.Fprintf(&sb, "\n%s.0 mode=%s modes=%s", name, r.Current, strings.Join(r.Modes, ","))
	}
	return sb.String()
}

func (s *Server) listDisplayModes() string {
	regions := s.modes.Definitions()
	if len(regions) == 0 {
		return "OK (none)"
	}
	var sb strings.Builder
	sb.WriteString("OK")
	for _, r := range regions {
		name := r.Name
		if name == "" {
			name = fmt.Sprintf("%d", r.Index)
		}
		fmt.Fprintf(&sb, "\n  region=%s.0", name)
		if r.Name != "" {
			fmt.Fprintf(&sb, " name=%q", r.Name)
		}
		if r.Controller != "" {
			fmt.Fprintf(&sb, " controller=%s", r.Controller)
		}
		fmt.Fprintf(&sb, " current=%s", r.Current)
		for _, mode := range r.Modes {
			fmt.Fprintf(&sb, "\n    mode=%s title=%q", mode.ID, mode.Title)
			if mode.Scope != "" {
				fmt.Fprintf(&sb, " scope=%s", mode.Scope)
			}
			fmt.Fprintf(&sb, " summary=%q", mode.Summary)
		}
	}
	return sb.String()
}

func (s *Server) listStemma() string {
	scanner := s.scanner()
	if scanner == nil {
		return "OK (none)"
	}
	devs := scanner.Devices()
	if len(devs) == 0 {
		return "OK (none)"
	}
	var sb strings.Builder
	sb.WriteString("OK")
	for _, d := range devs {
		statusStr := "absent"
		if d.Present {
			statusStr = "present"
		}
		fmt.Fprintf(&sb, "\n  bus=%s addr=0x%02X status=%s name=%q",
			d.Bus, d.Addr, statusStr, d.Name)
	}
	return sb.String()
}

func (s *Server) listGPIO() string {
	pins := s.gpiomgr.Snapshot()
	if len(pins) == 0 {
		return "OK (none)"
	}
	var sb strings.Builder
	sb.WriteString("OK")
	for _, p := range pins {
		fmt.Fprintf(&sb, "\n  %s", p.String())
	}
	return sb.String()
}

func (s *Server) handleGPIOCmd(args []string) string {
	if len(args) < 1 {
		return msgGPIOUsage
	}
	sub := strings.ToLower(args[0])

	switch sub {
	case gpioStatus:
		return s.listGPIO()
	case gpioPins:
		if s.pinReport == nil {
			return msgPinReportUnavailable
		}
		return s.pinReport()
	case gpioSet:
		if len(args) < 3 {
			return msgGPIOSetUsage
		}
		pinNum, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Sprintf("ERR invalid pin number %q", args[1])
		}
		lvlInt, err := strconv.Atoi(args[2])
		if err != nil || (lvlInt != 0 && lvlInt != 1) {
			return msgLevelInvalid
		}
		lvl := gpio.Level(lvlInt == 1)
		if err := s.gpiomgr.SetOutput(pinNum, lvl); err != nil {
			return fmt.Sprintf("ERR %v", err)
		}
		return fmt.Sprintf("OK GPIO%d set to %d", pinNum, lvlInt)
	case gpioIn:
		if len(args) < 2 {
			return msgGPIOInUsage
		}
		pinNum, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Sprintf("ERR invalid pin number %q", args[1])
		}
		if err := s.gpiomgr.SetInput(pinNum, gpio.PullNoChange); err != nil {
			return fmt.Sprintf("ERR %v", err)
		}
		lvl, err := s.gpiomgr.Read(pinNum)
		if err != nil {
			return fmt.Sprintf("ERR %v", err)
		}
		val := 0
		if lvl {
			val = 1
		}
		return fmt.Sprintf("OK GPIO%d input level=%d", pinNum, val)
	default:
		return msgGPIOUsage
	}
}

func (s *Server) handleStemmaCmd(args []string) string {
	if len(args) < 1 {
		return msgStemmaUsage
	}
	sub := strings.ToLower(args[0])
	switch sub {
	case stemmaStatus:
		return s.listStemma()
	default:
		return msgStemmaUsage
	}
}

// Dial is a convenience helper for tests/tools to connect to a running Server.
func Dial(sockPath string) (net.Conn, error) {
	return net.Dial("unix", sockPath)
}

// SendCommand sends a single command to the server at sockPath and returns
// the response line(s) up to and including the first "OK …" or "ERR …" line.
func SendCommand(sockPath, cmd string) (string, error) {
	conn, err := Dial(sockPath)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	sc := bufio.NewScanner(conn)
	// Read the greeting.
	if !sc.Scan() {
		return "", io.EOF
	}

	fmt.Fprintf(conn, "%s\n", cmd)

	var lines []string
	for sc.Scan() {
		lines = append(lines, sc.Text())
		if strings.HasPrefix(sc.Text(), "OK ") || strings.HasPrefix(sc.Text(), "ERR ") {
			break
		}
	}
	return strings.Join(lines, "\n"), nil
}
