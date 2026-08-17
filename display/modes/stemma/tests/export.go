package tests

import (
	displaymodes "github.com/databeast/cyberhud/display/modes"
	"github.com/databeast/cyberhud/display/modes/stemma"
	"github.com/databeast/cyberhud/display/modes/stemma/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

var StemmaRegistryEnumerate = func() []style.Style[source.StemmaSnapshot, source.Policy] {
	return stemma.StemmaRegistryEnumerate()
}

var StemmaRegistryDefaultName = stemma.StemmaRegistryDefaultName

func NewInstanceForTest(hints textlayout.TextHints) displaymodes.ModeInstance {
	return stemma.NewInstanceForTest()
}
