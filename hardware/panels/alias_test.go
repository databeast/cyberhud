package panels_test

import (
	"testing"

	stylealias "github.com/databeast/cyberhud/display/style"
	_ "github.com/databeast/cyberhud/hardware/panels/all"
)

func TestWaveshare13OLEDHatWifiAlias(t *testing.T) {
	got, ok := stylealias.Lookup("waveshare-1.3-oled-hat", "wifi")
	if !ok {
		t.Fatal("expected wifi style alias for waveshare-1.3-oled-hat")
	}
	if got != "mono-fast-128x64" {
		t.Fatalf("wifi alias = %q, want %q", got, "mono-fast-128x64")
	}
}
