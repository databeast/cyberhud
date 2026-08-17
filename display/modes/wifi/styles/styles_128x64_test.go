package styles

import (
	"image"
	"testing"

	"github.com/databeast/cyberhud/display/modes/wifi/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

func TestMonoFast128x64Build_TrimsToThreeRowsWhenConnected(t *testing.T) {
	hints := textlayout.DefaultTextHints(image.Rect(0, 0, 128, 64))
	ctx := style.NewStyleContext(hints)

	vd := monoFast128x64Build(source.WifiData{
		SSID:            "very-long-network-name",
		SignalStrength:  -48,
		LinkQuality:     92,
		IPAddress:       "192.168.1.42",
		Frequency:       2.437,
		Channel:         6,
		InterfaceName:   "wlan0",
		LinkSpeed:       72,
		ConnectionState: source.Connected,
	}, source.DefaultPolicy(), ctx, def{})

	if got, want := len(vd.Items), 3; got != want {
		t.Fatalf("len(vd.Items) = %d, want %d", got, want)
	}
	if got, want := vd.VisibleCount, 3; got != want {
		t.Fatalf("VisibleCount = %d, want %d", got, want)
	}
	if vd.Items[0] == "" {
		t.Fatal("top row is empty, want SSID text")
	}
	if vd.Items[1] == "" {
		t.Fatal("signal row is empty, want signal text")
	}
	if vd.Items[2] != "192.168.1.42" {
		t.Fatalf("third row = %q, want IP address", vd.Items[2])
	}
}
