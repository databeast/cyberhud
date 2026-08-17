package source

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	tserial "github.com/tarm/serial"
)

var state = struct {
	sync.RWMutex
	policy       Policy
	lines        []string
	partial      string
	pendingCR    bool
	seq          uint64
	connected    bool
	port         string
	baud         int
	lastError    string
	errCategory  ErrorCategory
	throughput   *ThroughputTracker
	scrollOffset int
	updatedAt    time.Time
	now          func() time.Time
	stopCh       chan struct{} // closed to signal the monitor goroutine to exit
}{
	policy: DefaultPolicy(),
	now:    time.Now,
}

// PolicySnapshot returns the current serial monitor policy.
func PolicySnapshot() Policy {
	state.RLock()
	defer state.RUnlock()
	return state.policy
}

// SetPolicy updates the serial monitor policy and resets dependent state when needed.
func SetPolicy(policy Policy) {
	state.Lock()
	defer state.Unlock()
	policy = NormalizePolicy(policy)
	prev := state.policy
	state.policy = policy
	if prev.Port != policy.Port || prev.Baud != policy.Baud {
		state.lines = nil
		state.partial = ""
		state.pendingCR = false
		state.lastError = ""
		if state.throughput != nil {
			state.throughput.Reset()
		}
	}
	if len(state.lines) > policy.MaxLines {
		state.lines = append([]string(nil), state.lines[len(state.lines)-policy.MaxLines:]...)
	}
	state.seq++
	state.updatedAt = state.now()
}

// SnapshotNow returns the current serial monitor snapshot.
func SnapshotNow() Snapshot {
	state.RLock()
	defer state.RUnlock()
	return snapshotLocked()
}

func SetSnapshotForTest(s Snapshot, p Policy) {
	state.Lock()
	defer state.Unlock()
	state.connected = s.Connected
	state.port = s.Port
	state.baud = s.Baud
	state.lines = append([]string(nil), s.Lines...)
	state.partial = ""
	state.pendingCR = false
	state.lastError = s.LastError
	state.errCategory = s.ErrorCategory
	state.seq = s.Sequence
	state.scrollOffset = s.ScrollOffset
	state.policy = NormalizePolicy(p)
	state.updatedAt = state.now()
}

func ResetForTest() {
	state.Lock()
	defer state.Unlock()
	state.connected = false
	state.port = ""
	state.baud = 0
	state.lines = nil
	state.partial = ""
	state.pendingCR = false
	state.lastError = ""
	state.errCategory = ErrNone
	state.seq = 0
	state.scrollOffset = 0
	state.policy = DefaultPolicy()
	state.updatedAt = state.now()
}

func snapshotLocked() Snapshot {
	lines := append([]string(nil), state.lines...)

	// Generate per-line ANSI color segments.
	var lineColors [][]ColorSegment
	if len(lines) > 0 {
		lineColors = make([][]ColorSegment, len(lines))
		for i, line := range lines {
			_, segments := ParseLine(line)
			lineColors[i] = segments
		}
	}

	// Get throughput history.
	var throughput [32]int
	if state.throughput != nil {
		throughput = state.throughput.History()
	}

	return Snapshot{
		Sequence:      state.seq,
		Connected:     state.connected,
		Port:          state.port,
		Baud:          state.baud,
		AutoSelect:    state.policy.AutoSelect,
		LastError:     state.lastError,
		ErrorCategory: state.errCategory,
		Lines:         lines,
		LineColors:    lineColors,
		Throughput:    throughput,
		ScrollOffset:  state.scrollOffset,
		UpdatedAt:     state.updatedAt,
	}
}

// Clear clears the retained output buffer.
func Clear() {
	state.Lock()
	defer state.Unlock()
	state.lines = nil
	state.partial = ""
	state.pendingCR = false
	state.lastError = ""
	if state.throughput != nil {
		state.throughput.Reset()
	}
	state.seq++
	state.updatedAt = state.now()
}

// Activate starts the serial monitor background goroutine.
// Safe to call multiple times; only one goroutine runs at a time.
func Activate() {
	state.Lock()
	defer state.Unlock()
	if state.stopCh != nil {
		return // already running
	}
	if state.throughput == nil {
		state.throughput = &ThroughputTracker{}
	}
	state.stopCh = make(chan struct{})
	go monitorLoop(state.stopCh)
}

// Deactivate stops the serial monitor background goroutine.
func Deactivate() {
	state.Lock()
	ch := state.stopCh
	state.stopCh = nil
	state.Unlock()
	if ch != nil {
		close(ch)
	}
}

func monitorLoop(stop <-chan struct{}) {
	for {
		select {
		case <-stop:
			return
		default:
		}

		policy := PolicySnapshot()
		portName := resolvePort(policy)
		if portName == "" {
			setStatus(false, "", policy.Baud, "waiting for serial port")
			select {
			case <-stop:
				return
			case <-time.After(scanDelay(policy)):
			}
			continue
		}

		port, err := openPort(portName, policy)
		if err != nil {
			cat := Classify(err)
			setStatusWithCategory(false, portName, policy.Baud, err.Error(), cat)
			select {
			case <-stop:
				return
			case <-time.After(scanDelay(policy)):
			}
			continue
		}

		log.Printf("[cyberhudd] serial mode: connected to %s @%d", portName, policy.Baud)
		setStatusWithCategory(true, portName, policy.Baud, "", ErrNone)

		reopen, readErr := readLoop(port, policy, stop)
		_ = port.Close()

		if reopen {
			continue
		}
		if readErr != nil {
			cat := Classify(readErr)
			setStatusWithCategory(false, portName, policy.Baud, readErr.Error(), cat)
		} else {
			setStatusWithCategory(false, portName, policy.Baud, "", ErrNone)
		}

		select {
		case <-stop:
			return
		case <-time.After(200 * time.Millisecond):
		}
	}
}

func readLoop(port *tserial.Port, initial Policy, stop <-chan struct{}) (reopen bool, err error) {
	startSig := policySignature(initial)
	buf := make([]byte, 256)

	for {
		select {
		case <-stop:
			flushPartialLocked()
			return false, nil
		default:
		}
		if policySignature(PolicySnapshot()) != startSig {
			flushPartialLocked()
			return true, nil
		}

		n, readErr := port.Read(buf)
		if n > 0 {
			ingestBytes(buf[:n])
			if state.throughput != nil {
				state.throughput.Add(n)
			}
		}
		// Tick the throughput tracker on every iteration to advance the window.
		if state.throughput != nil {
			state.throughput.Tick(time.Now())
		}
		if readErr != nil {
			if isTimeoutErr(readErr) {
				continue
			}
			if errors.Is(readErr, io.EOF) {
				flushPartialLocked()
				return false, nil
			}
			flushPartialLocked()
			return false, readErr
		}
	}
}

func isTimeoutErr(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, os.ErrDeadlineExceeded) {
		return true
	}
	s := strings.ToLower(err.Error())
	return strings.Contains(s, "timeout") || strings.Contains(s, "timed out")
}

func setStatus(connected bool, port string, baud int, errMsg string) {
	setStatusWithCategory(connected, port, baud, errMsg, ErrNone)
}

func setStatusWithCategory(connected bool, port string, baud int, errMsg string, cat ErrorCategory) {
	state.Lock()
	defer state.Unlock()
	state.connected = connected
	if port != "" {
		state.port = port
	}
	state.baud = baud
	state.lastError = strings.TrimSpace(errMsg)
	state.errCategory = cat
	state.updatedAt = state.now()
	state.seq++
}

func resolvePort(p Policy) string {
	state.RLock()
	current := state.port
	connected := state.connected
	state.RUnlock()

	p = NormalizePolicy(p)
	if p.Port != "" {
		return p.Port
	}
	if connected && current != "" {
		return current
	}
	if !p.AutoSelect {
		return ""
	}
	ports := candidatePorts([]string{"/dev/serial/by-id/*", "/dev/ttyUSB*", "/dev/ttyACM*"})
	if len(ports) == 0 {
		return ""
	}
	return ports[0]
}

func candidatePorts(patterns []string) []string {
	type cand struct {
		name string
		mod  time.Time
	}
	seen := map[string]bool{}
	list := make([]cand, 0)
	for _, pattern := range patterns {
		matches, _ := filepath.Glob(pattern)
		for _, name := range matches {
			if seen[name] {
				continue
			}
			seen[name] = true
			fi, err := os.Stat(name)
			mod := time.Time{}
			if err == nil {
				mod = fi.ModTime()
			}
			list = append(list, cand{name: name, mod: mod})
		}
	}
	sort.Slice(list, func(i, j int) bool {
		if list[i].mod.Equal(list[j].mod) {
			return list[i].name < list[j].name
		}
		return list[i].mod.After(list[j].mod)
	})
	out := make([]string, 0, len(list))
	for _, c := range list {
		out = append(out, c.name)
	}
	return out
}

func openPort(name string, p Policy) (*tserial.Port, error) {
	cfg := &tserial.Config{Name: name, Baud: p.Baud, ReadTimeout: scanDelay(p)}
	return tserial.OpenPort(cfg)
}

func scanDelay(p Policy) time.Duration {
	if p.ScanMS <= 0 {
		return DefaultScanMS * time.Millisecond
	}
	return time.Duration(p.ScanMS) * time.Millisecond
}

func policySignature(p Policy) string {
	p = NormalizePolicy(p)
	return fmt.Sprintf("%s|%d|%d|%t|%d", p.Port, p.Baud, p.MaxLines, p.AutoSelect, p.ScanMS)
}
