package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, out io.Writer, errOut io.Writer) int {
	fs := flag.NewFlagSet("cyberhudctl", flag.ContinueOnError)
	fs.SetOutput(errOut)
	socketPath := fs.String("socket", "/run/cyberhudd/console.sock", "path to cyberhudd Unix socket")
	timeout := fs.Duration("timeout", 2*time.Second, "socket read/write timeout")

	if err := fs.Parse(args); err != nil {
		return 2
	}

	remaining := fs.Args()
	if len(remaining) == 0 {
		fmt.Fprintf(errOut, "ERR missing command\n\n")
		printUsage(errOut)
		return 2
	}

	// Check for multi-command path: if args contain a standalone ";"
	if containsSemicolon(remaining) {
		return runMultiCommand(remaining, *socketPath, *timeout, out, errOut)
	}

	// Single command path
	cmd, local, err := protocolCommand(remaining)
	if err != nil {
		fmt.Fprintf(errOut, "ERR %v\n\n", err)
		printUsage(errOut)
		return 2
	}
	if local {
		// Command was handled locally (e.g., help)
		return 0
	}

	resp, err := sendCommand(*socketPath, cmd, *timeout)
	if err != nil {
		fmt.Fprintf(errOut, "ERR %v\n", err)
		return 1
	}
	fmt.Fprintln(out, formatResponseForDisplay(cmd, resp))
	if isErrorResponse(resp) {
		return 1
	}
	return 0
}

// containsSemicolon returns true if the args slice contains a standalone ";" element.
func containsSemicolon(args []string) bool {
	for _, a := range args {
		if a == ";" {
			return true
		}
	}
	return false
}

// runMultiCommand handles a multi-command invocation by passing each command
// through the same explicit protocol parser as a single command. There is no
// client-side region context or command rewriting.
func runMultiCommand(args []string, sockPath string, timeout time.Duration, out io.Writer, errOut io.Writer) int {
	cmdGroups := SplitCommands(args)
	if len(cmdGroups) == 0 {
		fmt.Fprintf(errOut, "ERR no commands found\n")
		return 2
	}

	var protoCmds []string
	for _, group := range cmdGroups {
		if len(group) == 0 {
			continue
		}

		cmd, local, err := protocolCommand(group)
		if err != nil {
			fmt.Fprintf(errOut, "ERR %v\n", err)
			return 2
		}
		if local {
			continue
		}
		protoCmds = append(protoCmds, cmd)
	}

	if len(protoCmds) == 0 {
		return 0
	}

	responses, err := sendMultipleCommands(sockPath, protoCmds, timeout)
	if err != nil {
		for _, r := range responses {
			if !isErrorResponse(r) {
				fmt.Fprintln(out, r)
			}
		}
		fmt.Fprintf(errOut, "ERR %v\n", err)
		return 1
	}

	for _, r := range responses {
		fmt.Fprintln(out, r)
	}
	return 0
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  cyberhudctl [flags] <command>")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "System:")
	fmt.Fprintln(w, "  cyberhudctl [flags] status")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "GPIO:")
	fmt.Fprintln(w, "  cyberhudctl [flags] gpio status")
	fmt.Fprintln(w, "  cyberhudctl [flags] gpio pins")
	fmt.Fprintln(w, "  cyberhudctl [flags] gpio set <pin> <0|1>")
	fmt.Fprintln(w, "  cyberhudctl [flags] gpio in <pin>")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "STEMMA/QWIIC:")
	fmt.Fprintln(w, "  cyberhudctl [flags] stemma status")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Display regions:")
	fmt.Fprintln(w, "  cyberhudctl [flags] display regions                  # list all regions")
	fmt.Fprintln(w, "  cyberhudctl [flags] display status")
	fmt.Fprintln(w, "  cyberhudctl [flags] display modes")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Display control (region addressing):")
	fmt.Fprintln(w, "  cyberhudctl [flags] display set <region> <mode> [key=value ...]")
	fmt.Fprintln(w, "  cyberhudctl [flags] display config <region> [key=value ...]")
	fmt.Fprintln(w, "  cyberhudctl [flags] display next <region>")
	fmt.Fprintln(w, "  cyberhudctl [flags] display prev <region>")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Policy:")
	fmt.Fprintln(w, "  cyberhudctl [flags] policy dump                      # dump all mode policies")
	fmt.Fprintln(w, "  cyberhudctl [flags] display policy <mode|region>     # query single mode/region policy")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Persistence:")
	fmt.Fprintln(w, "  cyberhudctl [flags] freeze                           # save hardware config (daemon-side)")
	fmt.Fprintln(w, "  cyberhudctl [flags] freeze policy                    # save all mode policies (daemon-side)")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Mode data commands:")
	fmt.Fprintln(w, "  cyberhudctl [flags] display ticker set <text>")
	fmt.Fprintln(w, "  cyberhudctl [flags] display ticker policy [key=value ...]")
	fmt.Fprintln(w, "  cyberhudctl [flags] display image set file <path>")
	fmt.Fprintln(w, "  cyberhudctl [flags] display image set base64 <data>")
	fmt.Fprintln(w, "  cyberhudctl [flags] display image clear")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Multi-command:")
	fmt.Fprintln(w, "  cyberhudctl [flags] display set main.0 clock ';' status")
	fmt.Fprintln(w, "  Each command is explicit; there is no client-side region context.")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Advanced:")
	fmt.Fprintln(w, "  cyberhudctl [flags] raw <line...>                    # send raw protocol line")
	fmt.Fprintln(w, "  cyberhudctl [flags] help modes                       # query daemon mode help")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Region reference:")
	fmt.Fprintln(w, "  Use <surface>.<index> notation (e.g., main.0, left-aux.0)")
	fmt.Fprintln(w, "  Or bare integer for coordinator index (e.g., 0, 1, 2)")
	fmt.Fprintln(w, "  Run 'cyberhudctl display regions' to see available regions")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Flags:")
	fmt.Fprintln(w, "  -socket   Unix socket path (default /run/cyberhudd/console.sock)")
	fmt.Fprintln(w, "  -timeout  socket timeout (default 2s)")
}

// protocolCommand routes CLI args to a protocol command string.
// Returns the command string, whether it was handled locally (no send needed), and any error.
func protocolCommand(args []string) (string, bool, error) {
	if len(args) == 0 {
		return "", false, errors.New("missing command")
	}
	switch strings.ToLower(args[0]) {
	case "status":
		if len(args) != 1 {
			return "", false, errors.New("usage: status")
		}
		return "status", false, nil

	case "stemma":
		if len(args) < 2 {
			return "", false, errors.New("usage: stemma status")
		}
		sub := strings.ToLower(args[1])
		switch sub {
		case "status":
			if len(args) != 2 {
				return "", false, errors.New("usage: stemma status")
			}
			return "stemma status", false, nil
		default:
			return "", false, errors.New("usage: stemma status")
		}

	case "gpio":
		if len(args) < 2 {
			return "", false, errors.New("usage: gpio <status|set|in|pins> ...")
		}
		sub := strings.ToLower(args[1])
		switch sub {
		case "status":
			if len(args) != 2 {
				return "", false, errors.New("usage: gpio status")
			}
			return "gpio status", false, nil
		case "pins":
			if len(args) != 2 {
				return "", false, errors.New("usage: gpio pins")
			}
			return "gpio pins", false, nil
		case "in":
			if len(args) != 3 {
				return "", false, errors.New("usage: gpio in <pin>")
			}
			if _, err := strconv.Atoi(args[2]); err != nil {
				return "", false, errors.New("pin must be an integer")
			}
			return "gpio in " + args[2], false, nil
		case "set":
			if len(args) != 4 {
				return "", false, errors.New("usage: gpio set <pin> <0|1>")
			}
			if _, err := strconv.Atoi(args[2]); err != nil {
				return "", false, errors.New("pin must be an integer")
			}
			if args[3] != "0" && args[3] != "1" {
				return "", false, errors.New("level must be 0 or 1")
			}
			return "gpio set " + args[2] + " " + args[3], false, nil
		default:
			return "", false, errors.New("usage: gpio <status|set|in|pins> ...")
		}

	case "display":
		if len(args) < 2 {
			return "", false, errors.New("usage: display <regions|status|modes|set|config|next|prev|policy|...> ...")
		}
		sub := strings.ToLower(args[1])
		switch sub {
		case "regions":
			if len(args) != 2 {
				return "", false, errors.New("usage: display regions")
			}
			return "display regions", false, nil
		case "status", "list", "modes":
			if len(args) != 2 {
				return "", false, fmt.Errorf("usage: display %s", sub)
			}
			return "display " + sub, false, nil
		case "set":
			if len(args) < 4 {
				return "", false, errors.New("usage: display set <region> <mode> [key=value ...]")
			}
			if strings.TrimSpace(args[3]) == "" {
				return "", false, errors.New("mode must not be empty")
			}
			parts := []string{"display set", strings.ToLower(strings.TrimSpace(args[2])), strings.ToLower(strings.TrimSpace(args[3]))}
			for _, kv := range args[4:] {
				parts = append(parts, strings.TrimSpace(kv))
			}
			return strings.Join(parts, " "), false, nil
		case "config":
			if len(args) < 3 {
				return "", false, errors.New("usage: display config <region> [key=value ...]")
			}
			parts := []string{"display config", strings.ToLower(strings.TrimSpace(args[2]))}
			for _, kv := range args[3:] {
				parts = append(parts, strings.TrimSpace(kv))
			}
			return strings.Join(parts, " "), false, nil
		case "next", "prev":
			if len(args) != 3 {
				return "", false, fmt.Errorf("usage: display %s <region>", sub)
			}
			return "display " + sub + " " + strings.ToLower(strings.TrimSpace(args[2])), false, nil
		case "policy":
			// display policy <mode|region>
			if len(args) < 3 {
				return "", false, errors.New("usage: display policy <mode|region>")
			}
			parts := []string{"display policy"}
			for _, a := range args[2:] {
				parts = append(parts, strings.TrimSpace(a))
			}
			return strings.Join(parts, " "), false, nil
		default:
			// Mode-specific sub-commands: "display <mode> [args...]"
			// e.g., "display ticker set hello", "display image clear"
			parts := make([]string, 0, len(args))
			parts = append(parts, "display")
			for _, a := range args[1:] {
				parts = append(parts, strings.TrimSpace(a))
			}
			return strings.Join(parts, " "), false, nil
		}

	case "freeze":
		// "freeze" or "freeze policy" — both are protocol commands sent to daemon
		if len(args) == 1 {
			return "freeze", false, nil
		}
		if len(args) == 2 && strings.ToLower(args[1]) == "policy" {
			return "freeze policy", false, nil
		}
		return "", false, errors.New("usage: freeze | freeze policy")

	case "policy":
		if len(args) < 2 {
			return "", false, errors.New("usage: policy dump")
		}
		sub := strings.ToLower(args[1])
		switch sub {
		case "dump":
			if len(args) != 2 {
				return "", false, errors.New("usage: policy dump")
			}
			return "policy dump", false, nil
		default:
			return "", false, errors.New("usage: policy dump")
		}

	case "region":
		// In single-command mode, "region" is treated as a protocol command
		// that the daemon resolves (or used as context in multi-command mode).
		// For single commands, we pass it through to the daemon.
		if len(args) < 2 {
			return "", false, errors.New("usage: region <region_id>")
		}
		parts := make([]string, 0, len(args))
		for _, a := range args {
			parts = append(parts, strings.TrimSpace(a))
		}
		return strings.Join(parts, " "), false, nil

	case "raw":
		if len(args) < 2 {
			return "", false, errors.New("usage: raw <line...>")
		}
		return strings.Join(args[1:], " "), false, nil

	case "help":
		if len(args) < 2 {
			return "", false, errors.New("usage: help modes")
		}
		return strings.Join(args, " "), false, nil

	default:
		return "", false, fmt.Errorf("unknown command %q; see usage", args[0])
	}
}

// sendMultipleCommands opens a single socket connection, sends each command
// sequentially, reads each response, and only sends "quit" at the end.
// If any command returns an error response, execution stops immediately.
func sendMultipleCommands(sockPath string, cmds []string, timeout time.Duration) ([]string, error) {
	if timeout <= 0 {
		timeout = 2 * time.Second
	}
	if len(cmds) == 0 {
		return nil, errors.New("no commands to send")
	}

	conn, err := net.DialTimeout("unix", sockPath, timeout)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(timeout))

	r := bufio.NewReader(conn)

	// Read greeting line — must start with "OK"
	greeting, err := r.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("reading greeting: %w", err)
	}
	greeting = strings.TrimSpace(greeting)
	if !strings.HasPrefix(greeting, "OK") {
		return nil, fmt.Errorf("unexpected greeting: %s", greeting)
	}

	var responses []string

	for _, cmd := range cmds {
		// Extend deadline for each command
		_ = conn.SetDeadline(time.Now().Add(timeout))

		// Send command
		if _, err := io.WriteString(conn, cmd+"\n"); err != nil {
			return responses, fmt.Errorf("writing command %q: %w", cmd, err)
		}

		// Read response lines until we get a line starting with "OK" or "ERR"
		var respLines []string
		for {
			line, err := r.ReadString('\n')
			if err != nil {
				if ne, ok := err.(net.Error); ok && ne.Timeout() {
					return responses, fmt.Errorf("timeout reading response for %q", cmd)
				}
				if errors.Is(err, io.EOF) {
					return responses, fmt.Errorf("connection closed reading response for %q", cmd)
				}
				return responses, fmt.Errorf("reading response for %q: %w", cmd, err)
			}
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			respLines = append(respLines, line)
			if strings.HasPrefix(line, "OK") || strings.HasPrefix(line, "ERR") {
				break
			}
		}

		resp := strings.Join(respLines, "\n")
		responses = append(responses, resp)

		// Stop on first error response
		if isErrorResponse(resp) {
			// Send quit before returning
			_, _ = io.WriteString(conn, "quit\n")
			return responses, fmt.Errorf("command %q failed: %s", cmd, resp)
		}
	}

	// All commands succeeded — send quit
	_ = conn.SetDeadline(time.Now().Add(timeout))
	_, _ = io.WriteString(conn, "quit\n")

	return responses, nil
}

func sendCommand(sockPath, cmd string, timeout time.Duration) (string, error) {
	if timeout <= 0 {
		timeout = 2 * time.Second
	}
	conn, err := net.DialTimeout("unix", sockPath, timeout)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(timeout))

	r := bufio.NewReader(conn)
	if _, err := r.ReadString('\n'); err != nil {
		return "", err
	}

	if _, err := io.WriteString(conn, cmd+"\n"); err != nil {
		return "", err
	}
	if _, err := io.WriteString(conn, "quit\n"); err != nil {
		return "", err
	}

	lines := make([]string, 0, 8)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				break
			}
			if errors.Is(err, io.EOF) {
				break
			}
			return "", err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line == "OK bye" {
			break
		}
		lines = append(lines, line)
	}
	if len(lines) == 0 {
		return "", errors.New("empty response")
	}
	return strings.Join(lines, "\n"), nil
}

func isErrorResponse(resp string) bool {
	resp = strings.TrimSpace(resp)
	if resp == "" {
		return true
	}
	firstLine := resp
	if idx := strings.IndexByte(resp, '\n'); idx >= 0 {
		firstLine = resp[:idx]
	}
	return strings.HasPrefix(firstLine, "ERR ")
}

func formatResponseForDisplay(cmd, resp string) string {
	if strings.EqualFold(strings.TrimSpace(cmd), "help modes") {
		if pretty := formatHelpModesResponse(resp); pretty != "" {
			return pretty
		}
	}
	return resp
}

func formatHelpModesResponse(resp string) string {
	lines := strings.Split(strings.TrimSpace(resp), "\n")
	if len(lines) == 0 {
		return ""
	}
	if strings.TrimSpace(lines[0]) != "OK" {
		return ""
	}

	type modeOption struct {
		key     string
		typ     string
		def     string
		summary string
	}
	type modeCommand struct {
		verb    string
		usage   string
		summary string
		title   string
		scope   string
		options []modeOption
	}

	var out []modeCommand
	var current *modeCommand
	for _, raw := range lines[1:] {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "verb=") {
			kv := parseKeyValueLine(line)
			if kv["verb"] == "" {
				continue
			}
			cmd := modeCommand{
				verb:    kv["verb"],
				usage:   kv["usage"],
				summary: kv["summary"],
				title:   kv["mode_title"],
				scope:   kv["scope"],
			}
			out = append(out, cmd)
			current = &out[len(out)-1]
			continue
		}
		if strings.HasPrefix(line, "option=") && current != nil {
			kv := parseKeyValueLine(line)
			if kv["option"] == "" {
				continue
			}
			current.options = append(current.options, modeOption{
				key:     kv["option"],
				typ:     kv["type"],
				def:     kv["default"],
				summary: kv["summary"],
			})
		}
	}
	if len(out) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString("Mode Commands:\n")
	for i, cmd := range out {
		if i > 0 {
			b.WriteByte('\n')
		}
		fmt.Fprintf(&b, "- %s", cmd.verb)
		if cmd.title != "" || cmd.scope != "" {
			fmt.Fprintf(&b, " (%s", cmd.title)
			if cmd.scope != "" {
				fmt.Fprintf(&b, ", scope=%s", cmd.scope)
			}
			b.WriteString(")")
		}
		b.WriteByte('\n')
		if cmd.summary != "" {
			fmt.Fprintf(&b, "  %s\n", cmd.summary)
		}
		if cmd.usage != "" {
			fmt.Fprintf(&b, "  usage: %s\n", cmd.usage)
		}
		if len(cmd.options) > 0 {
			b.WriteString("  options:\n")
			for _, opt := range cmd.options {
				fmt.Fprintf(&b, "    - %s (%s, default=%s): %s\n", opt.key, opt.typ, opt.def, opt.summary)
			}
		}
	}
	return strings.TrimRight(b.String(), "\n")
}

var keyValuePattern = regexp.MustCompile(`([a-zA-Z_]+)=("[^"]*"|[^ ]+)`)

func parseKeyValueLine(line string) map[string]string {
	out := map[string]string{}
	matches := keyValuePattern.FindAllStringSubmatch(line, -1)
	for _, m := range matches {
		if len(m) != 3 {
			continue
		}
		val := m[2]
		if strings.HasPrefix(val, `"`) && strings.HasSuffix(val, `"`) {
			if u, err := strconv.Unquote(val); err == nil {
				val = u
			}
		}
		out[m[1]] = val
	}
	return out
}
