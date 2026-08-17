package color

import (
	"image/color"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"

	"pgregory.net/rapid"
)

// --- From: color_prop_test.go ---

func TestProp_UnrecognizedNameFallback(t *testing.T) {
	registered := map[string]bool{
		"cyan":  true,
		"green": true,
		"amber": true,
		"red":   true,
		"white": true,
	}

	opaqueWhite := color.RGBA{255, 255, 255, 255}

	f := func(name string) bool {
		if registered[name] {
			// Skip registered names — they are not unrecognized
			return true
		}

		lookupResult := Lookup(name)
		if lookupResult != opaqueWhite {
			t.Logf("Lookup(%q) = %v, want %v", name, lookupResult, opaqueWhite)
			return false
		}

		resolveResult := ResolveAccent(name)
		if resolveResult != opaqueWhite {
			t.Logf("ResolveAccent(%q) = %v, want %v", name, resolveResult, opaqueWhite)
			return false
		}

		return true
	}

	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Errorf("Property 2 failed: %v", err)
	}
}

func TestProp_ResolveAccentLookupConsistency(t *testing.T) {
	names := Names()
	if len(names) == 0 {
		t.Fatal("palette has no registered names")
	}

	f := func(index uint) bool {
		// Pick a registered name and verify ResolveAccent matches Lookup.
		name := names[index%uint(len(names))]
		return ResolveAccent(name) == Lookup(name)
	}

	cfg := &quick.Config{
		MaxCount: 100,
		Values: func(args []reflect.Value, rng *rand.Rand) {
			args[0] = reflect.ValueOf(uint(rng.Intn(len(names) * 10)))
		},
	}

	if err := quick.Check(f, cfg); err != nil {
		t.Errorf("Property 3 failed: %v", err)
	}
}

func TestProp_DimFormulaCorrectness(t *testing.T) {
	f := func(r, g, b, a uint8) bool {
		c := color.RGBA{R: r, G: g, B: b, A: a}
		got := Dim(c)
		expected := color.RGBA{R: r / 2, G: g / 2, B: b / 2, A: 255}
		return got == expected
	}

	cfg := &quick.Config{MaxCount: 100}

	if err := quick.Check(f, cfg); err != nil {
		t.Errorf("Property 4 failed: %v", err)
	}
}

func TestProp_DimDoubleApplication(t *testing.T) {
	f := func(r, g, b, a uint8) bool {
		c := color.RGBA{R: r, G: g, B: b, A: a}
		got := Dim(Dim(c))
		expected := color.RGBA{R: r / 4, G: g / 4, B: b / 4, A: 255}
		return got == expected
	}

	cfg := &quick.Config{MaxCount: 100}

	if err := quick.Check(f, cfg); err != nil {
		t.Errorf("Property 5 failed: %v", err)
	}
}

func TestProp_BinaryPaletteSelect(t *testing.T) {
	f := func(aR, aG, aB, aA, iR, iG, iB, iA uint8) bool {
		active := color.RGBA{R: aR, G: aG, B: aB, A: aA}
		inactive := color.RGBA{R: iR, G: iG, B: iB, A: iA}
		bp := NewBinaryPalette(active, inactive)

		if bp.Select(true) != active {
			return false
		}
		if bp.Select(false) != inactive {
			return false
		}
		return true
	}

	cfg := &quick.Config{MaxCount: 100}

	if err := quick.Check(f, cfg); err != nil {
		t.Errorf("Property 6 failed: %v", err)
	}
}

func TestProp_BuildSliceDisabledReturnsNil(t *testing.T) {
	type buildSliceInput struct {
		States    []bool
		ActiveR   uint8
		ActiveG   uint8
		ActiveB   uint8
		ActiveA   uint8
		InactiveR uint8
		InactiveG uint8
		InactiveB uint8
		InactiveA uint8
	}

	f := func(input buildSliceInput) bool {
		active := color.RGBA{R: input.ActiveR, G: input.ActiveG, B: input.ActiveB, A: input.ActiveA}
		inactive := color.RGBA{R: input.InactiveR, G: input.InactiveG, B: input.InactiveB, A: input.InactiveA}
		bp := NewBinaryPalette(active, inactive)

		result := BuildSlice(input.States, bp, false)
		return result == nil
	}

	cfg := &quick.Config{MaxCount: 100}

	if err := quick.Check(f, cfg); err != nil {
		t.Errorf("Property 8 failed: %v", err)
	}
}

func TestProp_BuildSliceElementWiseCorrectness(t *testing.T) {
	f := func(states []bool, ar, ag, ab, aa, ir, ig, ib, ia uint8) bool {
		pal := BinaryPalette{
			Active:   color.RGBA{R: ar, G: ag, B: ab, A: aa},
			Inactive: color.RGBA{R: ir, G: ig, B: ib, A: ia},
		}

		result := BuildSlice(states, pal, true)

		if len(result) != len(states) {
			t.Logf("length mismatch: got %d, want %d", len(result), len(states))
			return false
		}

		for i, s := range states {
			expected := color.Color(pal.Select(s))
			if result[i] != expected {
				t.Logf("index %d: got %v, want %v", i, result[i], expected)
				return false
			}
		}

		return true
	}

	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Errorf("Property 7 failed: %v", err)
	}
}

// --- From: color_test.go ---

func TestNames(t *testing.T) {
	got := Names()
	expected := []string{"amber", "cyan", "emerald", "green", "red", "white"}

	if len(got) != len(expected) {
		t.Fatalf("Names() returned %d entries, want %d", len(got), len(expected))
	}
	for i, name := range expected {
		if got[i] != name {
			t.Errorf("Names()[%d] = %q, want %q", i, got[i], name)
		}
	}
}

func TestLookupKnownColors(t *testing.T) {
	tests := []struct {
		name string
		want color.RGBA
	}{
		{"cyan", color.RGBA{0, 255, 255, 255}},
		{"green", color.RGBA{0, 200, 0, 255}},
		{"amber", color.RGBA{255, 191, 0, 255}},
		{"red", color.RGBA{255, 0, 0, 255}},
		{"white", color.RGBA{255, 255, 255, 255}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Lookup(tc.name)
			if got != tc.want {
				t.Errorf("Lookup(%q) = %v, want %v", tc.name, got, tc.want)
			}
		})
	}
}

func TestResolveAccentNone(t *testing.T) {
	got := ResolveAccent("none")
	want := color.RGBA{255, 255, 255, 255}
	if got != want {
		t.Errorf("ResolveAccent(\"none\") = %v, want %v", got, want)
	}
}

func TestDimZeroValue(t *testing.T) {
	got := Dim(color.RGBA{0, 0, 0, 0})
	want := color.RGBA{0, 0, 0, 255}
	if got != want {
		t.Errorf("Dim({0,0,0,0}) = %v, want %v", got, want)
	}
}

func TestBuildSliceEmpty(t *testing.T) {
	pal := BinaryPalette{
		Active:   color.RGBA{0, 255, 0, 255},
		Inactive: color.RGBA{255, 0, 0, 255},
	}
	got := BuildSlice([]bool{}, pal, true)
	if got == nil {
		t.Fatal("BuildSlice([]bool{}, pal, true) returned nil, want non-nil empty slice")
	}
	if len(got) != 0 {
		t.Errorf("BuildSlice([]bool{}, pal, true) length = %d, want 0", len(got))
	}
}

func TestGPIOPaletteValues(t *testing.T) {
	wantActive := color.RGBA{0x00, 0xCC, 0x44, 0xFF}
	wantInactive := color.RGBA{0x66, 0x66, 0x66, 0xFF}

	if GPIOPalette.Active != wantActive {
		t.Errorf("GPIOPalette.Active = %v, want %v", GPIOPalette.Active, wantActive)
	}
	if GPIOPalette.Inactive != wantInactive {
		t.Errorf("GPIOPalette.Inactive = %v, want %v", GPIOPalette.Inactive, wantInactive)
	}
}

// --- From: scale_prop_test.go ---

// genColor generates a random color.RGBA value.
func genColor(t *rapid.T, label string) color.RGBA {
	return color.RGBA{
		R: uint8(rapid.IntRange(0, 255).Draw(t, label+"_r")),
		G: uint8(rapid.IntRange(0, 255).Draw(t, label+"_g")),
		B: uint8(rapid.IntRange(0, 255).Draw(t, label+"_b")),
		A: uint8(rapid.IntRange(0, 255).Draw(t, label+"_a")),
	}
}

// genFactor generates a random float64 in [0.0, 1.0].
func genFactor(t *rapid.T, label string) float64 {
	return rapid.Float64Range(0.0, 1.0).Draw(t, label)
}

func TestProp_ScaleMonotonicity(t *testing.T) {
	// For any color c and two factors f1 > f2, each RGB channel of Scale(c, f1) >= Scale(c, f2).
	rapid.Check(t, func(t *rapid.T) {
		c := genColor(t, "color")
		f1 := genFactor(t, "f1")
		f2 := genFactor(t, "f2")

		// Ensure f1 > f2 by swapping if needed; skip if equal.
		if f1 < f2 {
			f1, f2 = f2, f1
		}
		if f1 == f2 {
			return // trivially true when equal
		}

		high := Scale(c, f1)
		low := Scale(c, f2)

		if high.R < low.R {
			t.Fatalf("monotonicity violated for R: Scale(%v, %f).R=%d < Scale(%v, %f).R=%d",
				c, f1, high.R, c, f2, low.R)
		}
		if high.G < low.G {
			t.Fatalf("monotonicity violated for G: Scale(%v, %f).G=%d < Scale(%v, %f).G=%d",
				c, f1, high.G, c, f2, low.G)
		}
		if high.B < low.B {
			t.Fatalf("monotonicity violated for B: Scale(%v, %f).B=%d < Scale(%v, %f).B=%d",
				c, f1, high.B, c, f2, low.B)
		}
	})
}

func TestProp_ScaleIdentity(t *testing.T) {
	// Scale(c, 1.0) returns c with alpha forced to 255.
	rapid.Check(t, func(t *rapid.T) {
		c := genColor(t, "color")

		got := Scale(c, 1.0)
		want := color.RGBA{R: c.R, G: c.G, B: c.B, A: 255}

		if got != want {
			t.Fatalf("Scale(%v, 1.0) = %v, want %v", c, got, want)
		}
	})
}

func TestProp_ScaleZero(t *testing.T) {
	// Scale(c, 0.0) returns black {0, 0, 0, 255}.
	rapid.Check(t, func(t *rapid.T) {
		c := genColor(t, "color")

		got := Scale(c, 0.0)
		want := color.RGBA{R: 0, G: 0, B: 0, A: 255}

		if got != want {
			t.Fatalf("Scale(%v, 0.0) = %v, want %v", c, got, want)
		}
	})
}

func TestProp_ScaleAlphaAlways255(t *testing.T) {
	// Regardless of input alpha, output alpha is always 255.
	rapid.Check(t, func(t *rapid.T) {
		c := genColor(t, "color")
		f := genFactor(t, "factor")

		got := Scale(c, f)

		if got.A != 255 {
			t.Fatalf("Scale(%v, %f).A = %d, want 255", c, f, got.A)
		}
	})
}

// genGradientN generates a reasonable gradient length in [1, 50].
func genGradientN(t *rapid.T, label string) int {
	return rapid.IntRange(1, 50).Draw(t, label)
}

func TestProp_GradientLengthAndMonotonicity(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		c := genColor(t, "color")
		n := genGradientN(t, "n")

		grad := Gradient(c, n)

		// Property: len(Gradient(c, n)) == n for n >= 1
		if len(grad) != n {
			t.Fatalf("Gradient(%v, %d) has length %d, want %d", c, n, len(grad), n)
		}

		// Property: first element is full brightness (Scale(c, 1.0))
		wantFirst := Scale(c, 1.0)
		if grad[0] != wantFirst {
			t.Fatalf("Gradient(%v, %d)[0] = %v, want Scale(c, 1.0) = %v", c, n, grad[0], wantFirst)
		}

		// Property: last element is black (Scale(c, 0.0)) — only applies for n >= 2
		// When n == 1, the single element is both first and last; Gradient returns full brightness.
		if n >= 2 {
			wantLast := Scale(c, 0.0)
			if grad[n-1] != wantLast {
				t.Fatalf("Gradient(%v, %d)[%d] = %v, want Scale(c, 0.0) = %v", c, n, n-1, grad[n-1], wantLast)
			}
		}

		// Property: monotonic decrease — for all i < j, channels of grad[i] >= channels of grad[j]
		for i := 0; i < n-1; i++ {
			if grad[i].R < grad[i+1].R {
				t.Fatalf("monotonicity violated at R: Gradient[%d].R=%d < Gradient[%d].R=%d",
					i, grad[i].R, i+1, grad[i+1].R)
			}
			if grad[i].G < grad[i+1].G {
				t.Fatalf("monotonicity violated at G: Gradient[%d].G=%d < Gradient[%d].G=%d",
					i, grad[i].G, i+1, grad[i+1].G)
			}
			if grad[i].B < grad[i+1].B {
				t.Fatalf("monotonicity violated at B: Gradient[%d].B=%d < Gradient[%d].B=%d",
					i, grad[i].B, i+1, grad[i+1].B)
			}
		}
	})
}

// --- From: scale_test.go ---

func TestScale_Identity(t *testing.T) {
	c := color.RGBA{R: 100, G: 150, B: 200, A: 128}
	got := Scale(c, 1.0)
	want := color.RGBA{R: 100, G: 150, B: 200, A: 255}
	if got != want {
		t.Errorf("Scale(c, 1.0) = %v, want %v", got, want)
	}
}

func TestScale_Zero(t *testing.T) {
	c := color.RGBA{R: 100, G: 150, B: 200, A: 255}
	got := Scale(c, 0.0)
	want := color.RGBA{R: 0, G: 0, B: 0, A: 255}
	if got != want {
		t.Errorf("Scale(c, 0.0) = %v, want %v", got, want)
	}
}

func TestScale_Half(t *testing.T) {
	c := color.RGBA{R: 100, G: 200, B: 50, A: 255}
	got := Scale(c, 0.5)
	want := color.RGBA{R: 50, G: 100, B: 25, A: 255}
	if got != want {
		t.Errorf("Scale(c, 0.5) = %v, want %v", got, want)
	}
}

func TestScale_NegativeFactor(t *testing.T) {
	c := color.RGBA{R: 200, G: 200, B: 200, A: 255}
	got := Scale(c, -0.5)
	want := color.RGBA{R: 0, G: 0, B: 0, A: 255}
	if got != want {
		t.Errorf("Scale(c, -0.5) = %v, want %v", got, want)
	}
}

func TestScale_FactorAboveOne(t *testing.T) {
	c := color.RGBA{R: 100, G: 150, B: 200, A: 50}
	got := Scale(c, 1.5)
	want := color.RGBA{R: 100, G: 150, B: 200, A: 255}
	if got != want {
		t.Errorf("Scale(c, 1.5) = %v, want %v", got, want)
	}
}

func TestScale_AlphaAlways255(t *testing.T) {
	c := color.RGBA{R: 100, G: 100, B: 100, A: 0}
	got := Scale(c, 0.5)
	if got.A != 255 {
		t.Errorf("Scale alpha = %d, want 255", got.A)
	}
}

func TestGradient_NilForNLessThanOne(t *testing.T) {
	c := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	if got := Gradient(c, 0); got != nil {
		t.Errorf("Gradient(c, 0) = %v, want nil", got)
	}
	if got := Gradient(c, -1); got != nil {
		t.Errorf("Gradient(c, -1) = %v, want nil", got)
	}
}

func TestGradient_SingleElement(t *testing.T) {
	c := color.RGBA{R: 100, G: 200, B: 50, A: 128}
	got := Gradient(c, 1)
	if len(got) != 1 {
		t.Fatalf("Gradient(c, 1) len = %d, want 1", len(got))
	}
	want := color.RGBA{R: 100, G: 200, B: 50, A: 255}
	if got[0] != want {
		t.Errorf("Gradient(c, 1)[0] = %v, want %v", got[0], want)
	}
}

func TestGradient_Length(t *testing.T) {
	c := color.RGBA{R: 255, G: 128, B: 64, A: 255}
	for _, n := range []int{1, 2, 5, 10, 100} {
		got := Gradient(c, n)
		if len(got) != n {
			t.Errorf("Gradient(c, %d) len = %d, want %d", n, len(got), n)
		}
	}
}

func TestGradient_FirstIsFull(t *testing.T) {
	c := color.RGBA{R: 200, G: 100, B: 50, A: 255}
	got := Gradient(c, 5)
	want := color.RGBA{R: 200, G: 100, B: 50, A: 255}
	if got[0] != want {
		t.Errorf("Gradient first = %v, want %v", got[0], want)
	}
}

func TestGradient_LastIsBlack(t *testing.T) {
	c := color.RGBA{R: 200, G: 100, B: 50, A: 255}
	got := Gradient(c, 5)
	want := color.RGBA{R: 0, G: 0, B: 0, A: 255}
	if got[4] != want {
		t.Errorf("Gradient last = %v, want %v", got[4], want)
	}
}

func TestGradient_MonotonicDecrease(t *testing.T) {
	c := color.RGBA{R: 200, G: 100, B: 50, A: 255}
	got := Gradient(c, 10)
	for i := 1; i < len(got); i++ {
		if got[i].R > got[i-1].R || got[i].G > got[i-1].G || got[i].B > got[i-1].B {
			t.Errorf("Gradient not monotonic at index %d: %v > %v", i, got[i], got[i-1])
		}
	}
}
