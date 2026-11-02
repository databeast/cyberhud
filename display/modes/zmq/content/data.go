package content

// ZMQData is the mode-specific snapshot consumed by style Build methods.
//
// It carries the current message buffer (Lines, newest last) plus the stream
// identity derived from the active policy, so styles can render a meaningful
// header (which endpoint/socket/topic the panel is monitoring) without reaching
// back into package state themselves.
type ZMQData struct {
	// Lines is the message buffer snapshot, oldest first and newest last.
	Lines []string

	// Endpoint is the configured ZMQ connection string (may be empty when the
	// mode is idle/unconfigured).
	Endpoint string

	// SocketType is the ZMQ socket pattern, "sub" or "pull".
	SocketType string

	// Topic is the SUB topic filter; empty means "all messages" (and is ignored
	// for PULL sockets).
	Topic string
}
