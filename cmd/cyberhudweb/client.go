package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

// Client talks to the cyberhud daemon's console Unix-domain socket. It is a thin
// wrapper that speaks the existing line protocol; the web UI is just another
// client of the same commands the CLI uses.
type Client struct {
	socketPath string
	timeout    time.Duration
}

// NewClient returns a Client that dials socketPath with the given per-command
// timeout. A zero timeout defaults to 2s.
func NewClient(socketPath string, timeout time.Duration) *Client {
	if timeout <= 0 {
		timeout = 2 * time.Second
	}
	return &Client{socketPath: socketPath, timeout: timeout}
}

// Send issues a single command and returns the full response as a slice of
// trimmed, non-empty lines (greeting and the trailing "OK bye" excluded).
//
// It uses the robust read strategy: send the command, immediately send "quit",
// then drain every line until "OK bye" or EOF. This captures multi-line "OK"
// blocks (display status/regions/modes/geometry) whose first line is a bare
// "OK" — a pattern that breaks naive readers that stop on the first OK line.
func (c *Client) Send(cmd string) ([]string, error) {
	conn, err := net.DialTimeout("unix", c.socketPath, c.timeout)
	if err != nil {
		return nil, fmt.Errorf("connect to daemon socket %q: %w", c.socketPath, err)
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(c.timeout))
	r := bufio.NewReader(conn)

	// Read and discard the greeting line.
	if _, err := r.ReadString('\n'); err != nil {
		return nil, fmt.Errorf("read greeting: %w", err)
	}

	if _, err := io.WriteString(conn, cmd+"\n"); err != nil {
		return nil, fmt.Errorf("write command: %w", err)
	}
	if _, err := io.WriteString(conn, "quit\n"); err != nil {
		return nil, fmt.Errorf("write quit: %w", err)
	}

	var lines []string
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				return lines, fmt.Errorf("timeout reading response for %q", cmd)
			}
			return lines, fmt.Errorf("read response for %q: %w", cmd, err)
		}
		line = strings.TrimRight(line, "\r\n")
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if trimmed == "OK bye" {
			break
		}
		// Preserve the raw (untrimmed) line so parsers keep indentation cues if
		// they want them; callers mostly TrimSpace per line.
		lines = append(lines, line)
	}
	return lines, nil
}

// IsErr reports whether a response (as returned by Send) is an error response,
// i.e. its first line begins with "ERR ".
func IsErr(lines []string) bool {
	if len(lines) == 0 {
		return true
	}
	return strings.HasPrefix(strings.TrimSpace(lines[0]), "ERR ")
}

// ErrText returns the error message from an ERR response's first line, or "" if
// the response is not an error.
func ErrText(lines []string) string {
	if len(lines) == 0 {
		return "empty response"
	}
	first := strings.TrimSpace(lines[0])
	if strings.HasPrefix(first, "ERR ") {
		return strings.TrimPrefix(first, "ERR ")
	}
	return ""
}
