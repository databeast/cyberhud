package systemd

import (
	"github.com/databeast/cyberhud/display/modes/systemd/source"
	"github.com/databeast/cyberhud/display/modes/systemd/styles"
)

type Policy = source.Policy
type Snapshot = source.Snapshot

func DefaultPolicy() Policy           { return source.DefaultPolicy() }
func BuildItems() []string            { return source.BuildDefaultItems(source.PollSnapshot()) }
func Signature() string               { return source.Signature(GetPolicy()) }
func SanitizeTarget(v string) string  { return source.SanitizeTarget(v) }
func BootFraction(s Snapshot) float64 { return styles.BootFraction(s) }

var SystemdRegistryExported = systemdRegistry

func ResetStateForTest() { source.ResetStateForTest() }
func SetProbesForTest(desired func() string, loading func() string, reached func() bool, targets func() (int, int)) {
	source.SetProbesForTest(desired, loading, reached, targets)
}
