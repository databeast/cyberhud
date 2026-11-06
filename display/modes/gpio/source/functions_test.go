package source

import (
	"testing"

	gpiomgr "github.com/databeast/cyberhud/hardware/gpio"
)

func TestLiveFunctionsGroupsBusesAndDropsDead(t *testing.T) {
	pins := []gpiomgr.PinState{
		// Two I2C1 pins → one bus row.
		{Number: 2, Mode: gpiomgr.ModeAlt, Func: "I2C1_SDA"},
		{Number: 3, Mode: gpiomgr.ModeAlt, Func: "I2C1_SCL"},
		// Five SPI0 pins → one bus row.
		{Number: 7, Mode: gpiomgr.ModeAlt, Func: "SPI0_CE1"},
		{Number: 8, Mode: gpiomgr.ModeAlt, Func: "SPI0_CE0"},
		{Number: 9, Mode: gpiomgr.ModeAlt, Func: "SPI0_MISO"},
		{Number: 10, Mode: gpiomgr.ModeAlt, Func: "SPI0_MOSI"},
		{Number: 11, Mode: gpiomgr.ModeAlt, Func: "SPI0_SCLK"},
		// Standalone output + input.
		{Number: 17, Mode: gpiomgr.ModeOutput, Level: true, Func: "Out"},
		{Number: 27, Mode: gpiomgr.ModeInput, Level: false, Func: "In"},
		// Dead pins → dropped.
		{Number: 22, Mode: gpiomgr.ModeUnknown},
		{Number: 23, Mode: gpiomgr.ModeUnknown, Func: ""},
	}

	got := LiveFunctions(pins)

	// Expect: output(17), input(27), bus(I2C1), bus(SPI0) in that order.
	if len(got) != 4 {
		t.Fatalf("expected 4 live functions, got %d: %+v", len(got), got)
	}
	if got[0].Kind != KindOutput || got[0].Label != "17" || !got[0].Level {
		t.Errorf("row0 want output 17 HI, got %+v", got[0])
	}
	if got[1].Kind != KindInput || got[1].Label != "27" || got[1].Level {
		t.Errorf("row1 want input 27 LO, got %+v", got[1])
	}
	if got[2].Kind != KindBus || got[2].Label != "I2C1" || got[2].PinCount != 2 {
		t.Errorf("row2 want bus I2C1 count 2, got %+v", got[2])
	}
	if got[3].Kind != KindBus || got[3].Label != "SPI0" || got[3].PinCount != 5 {
		t.Errorf("row3 want bus SPI0 count 5, got %+v", got[3])
	}
}

func TestLiveFunctionsNumericPinOrder(t *testing.T) {
	pins := []gpiomgr.PinState{
		{Number: 2, Mode: gpiomgr.ModeOutput, Func: "Out"},
		{Number: 17, Mode: gpiomgr.ModeOutput, Func: "Out"},
		{Number: 3, Mode: gpiomgr.ModeOutput, Func: "Out"},
	}
	got := LiveFunctions(pins)
	want := []string{"2", "3", "17"}
	if len(got) != 3 {
		t.Fatalf("expected 3, got %d", len(got))
	}
	for i, w := range want {
		if got[i].Label != w {
			t.Errorf("position %d: want %q, got %q (numeric order broken)", i, w, got[i].Label)
		}
	}
}

func TestFunctionRowRendering(t *testing.T) {
	cases := []struct {
		f    Function
		want string
	}{
		{Function{Kind: KindOutput, Label: "17", Level: true}, "17 OUT HI"},
		{Function{Kind: KindInput, Label: "27", Level: false}, "27 IN LO"},
		{Function{Kind: KindBus, Label: "SPI0"}, "SPI0 up"},
	}
	for _, c := range cases {
		if got := c.f.Row(); got != c.want {
			t.Errorf("Row() = %q, want %q", got, c.want)
		}
	}
}

func TestEmptyWhenAllDead(t *testing.T) {
	pins := []gpiomgr.PinState{
		{Number: 4, Mode: gpiomgr.ModeUnknown},
		{Number: 22, Mode: gpiomgr.ModeUnknown},
	}
	if got := LiveFunctions(pins); len(got) != 0 {
		t.Errorf("all-dead snapshot should yield no functions, got %+v", got)
	}
}
