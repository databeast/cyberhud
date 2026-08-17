package source

import (
	"fmt"
	"strings"
)

// ErrorCategory classifies serial port failures for targeted user guidance.
type ErrorCategory int

const (
	ErrNone           ErrorCategory = iota
	ErrPermission                   // "permission denied" in error string
	ErrDeviceNotFound               // "no such file or directory", "device not found"
	ErrDisconnected                 // "input/output error", "device disconnected"
	ErrBaudMismatch                 // "framing error"
	ErrUnknown                      // Anything else
)

// Classify maps a raw error to an ErrorCategory based on string matching.
func Classify(err error) ErrorCategory {
	if err == nil {
		return ErrNone
	}
	msg := strings.ToLower(err.Error())

	if strings.Contains(msg, "permission denied") {
		return ErrPermission
	}
	if strings.Contains(msg, "no such file or directory") || strings.Contains(msg, "device not found") {
		return ErrDeviceNotFound
	}
	if strings.Contains(msg, "input/output error") || strings.Contains(msg, "device disconnected") {
		return ErrDisconnected
	}
	if strings.Contains(msg, "framing error") {
		return ErrBaudMismatch
	}
	return ErrUnknown
}

// Guidance returns the user-facing help text for a given ErrorCategory.
func Guidance(cat ErrorCategory) string {
	switch cat {
	case ErrNone:
		return ""
	case ErrPermission:
		return "Add your user to the dialout group: sudo usermod -aG dialout $USER"
	case ErrDeviceNotFound:
		return "Check cable connection and device path"
	case ErrDisconnected:
		return "Cable disconnected. Reconnecting..."
	case ErrBaudMismatch:
		return "Verify baud rate matches the connected device"
	case ErrUnknown:
		return fmt.Sprintf("Unexpected error")
	default:
		return ""
	}
}
