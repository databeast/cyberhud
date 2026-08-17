package usb

import "github.com/databeast/cyberhud/display/style"

func RegisteredStyleNames() []string                        { return registeredStyleNames() }
func BoolStr(b bool) string                                 { return boolStr(b) }
func UsbRegistryEnumerate() []style.Style[Snapshot, Policy] { return usbRegistry.Enumerate() }
func LookupStyleForTest(name string) style.Style[Snapshot, Policy] {
	return usbRegistry.Lookup(name)
}
func SetTestMonitorState(snap Snapshot, p Policy) {
	monitorState.Lock()
	monitorState.snapshot = snap
	monitorState.policy = normalizePolicy(p)
	monitorState.lastScanAt = monitorState.now()
	monitorState.Unlock()
}
