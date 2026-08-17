package source

import (
	"strings"
)

func ingestBytes(b []byte) {
	if len(b) == 0 {
		return
	}

	state.Lock()
	defer state.Unlock()

	data := string(b)

	// If the previous chunk ended with \r, check if this one starts with \n.
	// If so, the pair forms a single \r\n — drop the leading \n since the \r
	// was already normalized to \n in the previous call.
	if state.pendingCR && len(data) > 0 && data[0] == '\n' {
		data = data[1:]
	}
	state.pendingCR = len(b) > 0 && b[len(b)-1] == '\r'

	// Normalize line endings.
	data = strings.ReplaceAll(data, "\r\n", "\n")
	data = strings.ReplaceAll(data, "\r", "\n")

	if state.partial != "" {
		data = state.partial + data
		state.partial = ""
	}
	if data == "" {
		return
	}

	parts := strings.Split(data, "\n")
	if !strings.HasSuffix(data, "\n") {
		state.partial = parts[len(parts)-1]
		parts = parts[:len(parts)-1]
	} else {
		parts = parts[:len(parts)-1]
	}

	for _, line := range parts {
		appendLineLocked(line)
	}
	state.updatedAt = state.now()
}

func flushPartialLocked() {
	state.Lock()
	defer state.Unlock()
	if strings.TrimSpace(state.partial) != "" {
		appendLineLocked(state.partial)
	}
	state.partial = ""
	state.pendingCR = false
	state.updatedAt = state.now()
}

func appendLineLocked(line string) {
	state.lines = append(state.lines, line)
	max := state.policy.MaxLines
	if max <= 0 {
		max = DefaultMaxLines
	}
	if len(state.lines) > max {
		state.lines = append([]string(nil), state.lines[len(state.lines)-max:]...)
	}
	state.seq++
}
