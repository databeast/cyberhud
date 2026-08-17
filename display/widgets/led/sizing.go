package led

// DiameterForRow returns the LED diameter that fits within a text row of the
// given height. The LED is inset by 1px top and bottom so it aligns with
// adjacent text without touching the row edges.
//
// Use this to size LEDs that sit alongside text labels in list-style layouts.
// The returned value is always at least MinDiameter (3).
func DiameterForRow(rowHeight int) int {
	d := rowHeight - 2
	if d < MinDiameter {
		d = MinDiameter
	}
	return d
}

// MinDiameter is the smallest valid LED diameter (matches Config validation).
const MinDiameter = 3
