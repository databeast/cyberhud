package attract_bokeh

import "github.com/databeast/cyberhud/display/modes/attract_bokeh/source"

type Policy = source.Policy

func DefaultPolicy() source.Policy { return source.DefaultPolicy() }

func NormalizePolicy(p source.Policy) source.Policy { return source.NormalizePolicy(p) }

func SnapshotPolicyForTest() map[string]interface{} { return bokehSnapshotter{}.SnapshotPolicy() }

func RestorePolicyForTest(data map[string]interface{}) error {
	return bokehSnapshotter{}.RestorePolicy(data)
}
