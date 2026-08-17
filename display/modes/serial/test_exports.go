package serial

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
)

func Registry() *style.StyleRegistry[source.Snapshot, source.Policy] { return serialRegistry }

func SetSnapshotForTest(s Snapshot, p Policy) { source.SetSnapshotForTest(s, p) }

func ResetForTest() { source.ResetForTest() }
