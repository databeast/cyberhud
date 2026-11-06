package displaymodes

import (
	"github.com/databeast/cyberhud/display/style"

	"testing"
)

// TestCapGapReportsUnderUtilization verifies lever 2's core: a StyleReport whose
// chosen style sits below the panel's capability tier reports a positive gap, and
// one that matches the panel reports zero. This is the signal LogStyleSelection
// turns into an "UNDER-UTILIZES" warning.
func TestCapGapReportsUnderUtilization(t *testing.T) {
	cases := []struct {
		name    string
		report  style.StyleReport
		wantGap int
	}{
		{
			name:    "mono style on color-fast panel wastes 5 tiers",
			report:  style.StyleReport{HasCapInfo: true, PanelCap: style.ColorFast, StyleCap: style.MonoSlow},
			wantGap: 5,
		},
		{
			name:    "color style on color-fast panel wastes nothing",
			report:  style.StyleReport{HasCapInfo: true, PanelCap: style.ColorFast, StyleCap: style.ColorFast},
			wantGap: 0,
		},
		{
			name:    "grayscale style on grayscale-fast panel is a full match",
			report:  style.StyleReport{HasCapInfo: true, PanelCap: style.GrayscaleFast, StyleCap: style.GrayscaleFast},
			wantGap: 0,
		},
		{
			name:    "no capability info reports unknown (zero) gap",
			report:  style.StyleReport{HasCapInfo: false, PanelCap: style.ColorFast, StyleCap: style.MonoSlow},
			wantGap: 0,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.report.CapGap(); got != tc.wantGap {
				t.Fatalf("CapGap() = %d, want %d", got, tc.wantGap)
			}
		})
	}
}
