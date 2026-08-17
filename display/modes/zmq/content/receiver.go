package content

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/go-zeromq/zmq4"
)

// receiver manages the ZMQ socket lifecycle in a background goroutine.
// It is unexported; a package-level instance is used by the public API.
type receiver struct {
	mu     sync.Mutex
	cancel context.CancelFunc
	active bool
}

// defaultReceiver is the package-level receiver instance.
var defaultReceiver receiver

// Activate starts the ZMQ receiver goroutine with the given policy and buffer.
// If the receiver is already active, the call is a no-op (idempotent).
// If the policy Endpoint is empty, the receiver stays idle.
func (r *receiver) Activate(policy Policy, buf *buffer) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.active {
		return
	}

	if strings.TrimSpace(policy.Endpoint) == "" {
		// No endpoint configured — stay idle without attempting connections.
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	r.cancel = cancel
	r.active = true

	go r.run(ctx, policy, buf)
}

// Deactivate stops the receiver goroutine and closes the socket.
// The message buffer contents are left unchanged.
// If the receiver is not active, the call is a no-op.
func (r *receiver) Deactivate() {
	r.mu.Lock()
	cancel := r.cancel
	r.cancel = nil
	r.active = false
	r.mu.Unlock()

	if cancel != nil {
		cancel()
	}
}

// IsActive returns whether the receiver goroutine is currently running.
func (r *receiver) IsActive() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.active
}

// run is the background goroutine that manages the ZMQ socket connection
// and message ingestion loop. It uses context cancellation for lifecycle
// management, which integrates cleanly with go-zeromq/zmq4's context-based
// socket creation.
func (r *receiver) run(ctx context.Context, policy Policy, buf *buffer) {
	defer func() {
		r.mu.Lock()
		r.active = false
		r.cancel = nil
		r.mu.Unlock()
	}()

	// Parse JSON fields from the comma-separated policy string.
	fields := parseJSONFields(policy.JSONFields)

	for {
		// Check for cancellation before attempting connection.
		select {
		case <-ctx.Done():
			return
		default:
		}

		sock := r.createSocket(ctx, policy)

		// Attempt to connect to the endpoint.
		err := sock.Dial(policy.Endpoint)
		if err != nil {
			log.Printf("[cyberhudd] zmq receiver: connection to %s failed: %v", policy.Endpoint, err)
			_ = sock.Close()
			// Wait 5 seconds before retrying, or exit if cancelled.
			select {
			case <-ctx.Done():
				return
			case <-time.After(5 * time.Second):
			}
			continue
		}

		log.Printf("[cyberhudd] zmq receiver: connected to %s (type=%s)", policy.Endpoint, policy.SocketType)

		// Run the message read loop. Returns when the context is cancelled
		// or an unrecoverable error occurs.
		r.readLoop(ctx, sock, buf, fields)

		// Close the socket after the read loop exits.
		_ = sock.Close()

		// If context is done, exit the goroutine entirely.
		select {
		case <-ctx.Done():
			return
		default:
			// Connection was lost; retry after delay.
			log.Printf("[cyberhudd] zmq receiver: connection lost, retrying in 5s")
			select {
			case <-ctx.Done():
				return
			case <-time.After(5 * time.Second):
			}
		}
	}
}

// createSocket creates the appropriate ZMQ socket (SUB or PULL) based on
// the policy's SocketType. For SUB sockets, it subscribes to the configured topic.
func (r *receiver) createSocket(ctx context.Context, policy Policy) zmq4.Socket {
	switch policy.SocketType {
	case "pull":
		return zmq4.NewPull(ctx)
	default:
		// Default to SUB socket.
		sub := zmq4.NewSub(ctx)
		// Subscribe to topic. Empty topic subscribes to all messages.
		if err := sub.SetOption(zmq4.OptionSubscribe, policy.Topic); err != nil {
			log.Printf("[cyberhudd] zmq receiver: subscribe to topic %q failed: %v", policy.Topic, err)
		}
		return sub
	}
}

// readLoop continuously reads messages from the socket and pushes them to
// the buffer, applying JSON field filtering when configured. It returns
// when the context is cancelled or a read error occurs.
func (r *receiver) readLoop(ctx context.Context, sock zmq4.Socket, buf *buffer, fields []string) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		msg, err := sock.Recv()
		if err != nil {
			// Check if the context was cancelled (graceful shutdown).
			select {
			case <-ctx.Done():
				return
			default:
			}
			// Log unexpected errors and break out for reconnection.
			log.Printf("[cyberhudd] zmq receiver: recv error: %v", err)
			return
		}

		raw := msg.String()
		if raw == "" {
			continue
		}

		// Apply JSON field filtering.
		if filtered := filterJSON(raw, fields); filtered != nil {
			for _, line := range filtered {
				buf.Push(line)
			}
		} else {
			buf.Push(raw)
		}
	}
}

// parseJSONFields splits the comma-separated JSONFields string into a
// trimmed slice of field names. Returns nil if the input is empty.
func parseJSONFields(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	var fields []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			fields = append(fields, p)
		}
	}
	if len(fields) == 0 {
		return nil
	}
	return fields
}
