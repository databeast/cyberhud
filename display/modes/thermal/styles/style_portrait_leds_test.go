package styles

import "testing"

func TestPortraitLEDDiameterRespectsGlowFootprint(t *testing.T) {
	tests := []struct {
		name       string
		width      int
		slotHeight int
		want       int
	}{
		{name: "large slot keeps LEDs smaller than full width", width: 64, slotHeight: 40, want: 30},
		{name: "smaller slot still keeps glow inside bounds", width: 32, slotHeight: 15, want: 11},
		{name: "tiny slot stays within available space", width: 16, slotHeight: 9, want: 5},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := portraitLEDDiameter(tc.width, tc.slotHeight)
			if got != tc.want {
				t.Fatalf("portraitLEDDiameter(%d, %d) = %d, want %d", tc.width, tc.slotHeight, got, tc.want)
			}
			if got > tc.width || got > tc.slotHeight {
				t.Fatalf("portraitLEDDiameter(%d, %d) = %d exceeded available width/slot", tc.width, tc.slotHeight, got)
			}
			if got > (tc.slotHeight*10)/13 {
				t.Fatalf("portraitLEDDiameter(%d, %d) = %d exceeded glow-safe cap of %d", tc.width, tc.slotHeight, got, (tc.slotHeight*10)/13)
			}
		})
	}
}
