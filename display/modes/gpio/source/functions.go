package source

import (
	"fmt"
	"sort"

	gpiomgr "github.com/databeast/cyberhud/hardware/gpio"
)

// FuncKind classifies a row in the live-functions view. The gpio mode shows
// "what the header is doing": each live function is either a standalone digital
// pin (output or input) whose level is the point, or a multi-pin bus/peripheral
// (SPI, I2C, …) collapsed to a single entry whose presence is the point.
type FuncKind uint8

const (
	// KindOutput is a standalone digital output pin. Its level (HIGH/LOW) is
	// the signal — this is the program's action made visible.
	KindOutput FuncKind = iota
	// KindInput is a standalone digital input pin. Its level is the signal —
	// what the outside world is presenting.
	KindInput
	// KindBus is a peripheral spanning one or more pins (SPI0, I2C1, UART0,
	// PWM0, …) collapsed to one row. Its presence ("up") is the signal; the
	// individual pin levels are noise and are not shown.
	KindBus
)

// Function is one row of the live GPIO view: a single thing the header is
// currently doing. Dead (unconfigured) pins never become a Function.
type Function struct {
	// Label is what the row shows on the left: for a bus, the bus name
	// ("SPI0", "I2C1"); for a standalone pin, its number ("17").
	Label string
	// PinCount is how many physical pins this function spans. 1 for standalone
	// pins; the member count for a bus.
	PinCount int
	Kind     FuncKind
	// Level is meaningful only for KindOutput / KindInput — true=HIGH.
	// Ignored for KindBus (a bus is shown as "up", not HI/LO).
	Level bool
}

// StateLabel renders the right-hand state token for the row, honoring the
// per-kind contract: standalone pins read HI/LO, a bus reads "up".
func (f Function) StateLabel() string {
	switch f.Kind {
	case KindBus:
		return "up"
	default:
		if f.Level {
			return "HI"
		}
		return "LO"
	}
}

// Row renders a full one-line "identity + state" string for roomy panels,
// e.g. "SPI0 up", "17 OUT HI", "23 IN LO".
func (f Function) Row() string {
	switch f.Kind {
	case KindBus:
		return fmt.Sprintf("%s %s", f.Label, f.StateLabel())
	case KindOutput:
		return fmt.Sprintf("%s OUT %s", f.Label, f.StateLabel())
	default:
		return fmt.Sprintf("%s IN %s", f.Label, f.StateLabel())
	}
}

// LiveFunctions reduces a raw pin snapshot to the ordered list of things the
// header is actively doing:
//
//   - Pins serving a recognized bus (SPI/I2C/UART/PWM/…) are collapsed so each
//     bus appears exactly once, spanning however many member pins it has.
//   - Standalone digital pins configured as output or input appear individually,
//     carrying their level.
//   - Unconfigured pins (ModeUnknown, no bus function) are dropped entirely —
//     "nothing" is not something worth a row on a headless box's only screen.
//
// Order (most actionable first): standalone outputs, then standalone inputs,
// then buses. Within each group, by pin number / bus name for stable output.
func LiveFunctions(pins []gpiomgr.PinState) []Function {
	type pinFn struct {
		fn  Function
		num int // pin number, for numeric ordering
	}
	var outputs, inputs []pinFn
	busMembers := make(map[string]int)
	var busOrder []string

	for _, p := range pins {
		if bus := p.BusGroup(); bus != "" {
			if _, seen := busMembers[bus]; !seen {
				busOrder = append(busOrder, bus)
			}
			busMembers[bus]++
			continue
		}
		switch p.Mode {
		case gpiomgr.ModeOutput:
			outputs = append(outputs, pinFn{
				fn:  Function{Kind: KindOutput, Label: fmt.Sprintf("%d", p.Number), Level: bool(p.Level), PinCount: 1},
				num: p.Number,
			})
		case gpiomgr.ModeInput:
			inputs = append(inputs, pinFn{
				fn:  Function{Kind: KindInput, Label: fmt.Sprintf("%d", p.Number), Level: bool(p.Level), PinCount: 1},
				num: p.Number,
			})
		default:
			// ModeUnknown / ModeAlt without a recognized bus family: not a live
			// standalone pin and not a known bus → drop it.
		}
	}

	sort.Slice(outputs, func(i, j int) bool { return outputs[i].num < outputs[j].num })
	sort.Slice(inputs, func(i, j int) bool { return inputs[i].num < inputs[j].num })
	sort.Strings(busOrder)

	out := make([]Function, 0, len(outputs)+len(inputs)+len(busOrder))
	for _, o := range outputs {
		out = append(out, o.fn)
	}
	for _, in := range inputs {
		out = append(out, in.fn)
	}
	for _, bus := range busOrder {
		out = append(out, Function{Kind: KindBus, Label: bus, PinCount: busMembers[bus]})
	}
	return out
}
