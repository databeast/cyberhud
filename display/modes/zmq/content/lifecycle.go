package content

// msgBuffer is the package-level message buffer shared between the receiver
// and the rendering/command layers. Its capacity tracks Policy.MaxLines.
var msgBuffer = newBuffer(DefaultPolicy().MaxLines)

// Activate starts the ZMQ receiver using the current policy and package-level
// buffer. If the receiver is already active, the call is a no-op.
func Activate() {
	policy := GetPolicy()
	defaultReceiver.Activate(policy, msgBuffer)
}

// Deactivate stops the ZMQ receiver. The message buffer contents are preserved.
// If the receiver is not active, the call is a no-op.
func Deactivate() {
	defaultReceiver.Deactivate()
}

// Clear empties the message buffer and increments the sequence counter.
func Clear() {
	msgBuffer.Clear()
}

// SnapshotLines returns a copy of the current message buffer.
func SnapshotLines() []string {
	return msgBuffer.Snapshot()
}

// BufferSeq returns the message buffer sequence counter.
func BufferSeq() uint64 {
	return msgBuffer.Seq()
}

// ClearMessagesForTest clears the package-level message buffer.
func ClearMessagesForTest() { msgBuffer.Clear() }

// PushMessageForTest appends one line to the package-level message buffer.
func PushMessageForTest(msg string) { msgBuffer.Push(msg) }
