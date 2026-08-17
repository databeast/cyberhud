package is31fl3731

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math/bits"
	"strings"
	"testing"

	"periph.io/x/conn/v3/physic"
	"pgregory.net/rapid"

	"github.com/databeast/cyberhud/hardware/driver"
)

// mockBus is a minimal i2c.Bus implementation that succeeds on all operations.
type mockBus struct{}

func (m *mockBus) String() string                    { return "mock" }
func (m *mockBus) Tx(_ uint16, _, _ []byte) error    { return nil }
func (m *mockBus) SetSpeed(_ physic.Frequency) error { return nil }

// recordingBus records all Tx writes for later inspection.
type recordingBus struct {
	writes [][]byte // each entry is a copy of the w slice from Tx
}

func (r *recordingBus) String() string                    { return "recording" }
func (r *recordingBus) SetSpeed(_ physic.Frequency) error { return nil }
func (r *recordingBus) Tx(_ uint16, w, _ []byte) error {
	cp := make([]byte, len(w))
	copy(cp, w)
	r.writes = append(r.writes, cp)
	return nil
}

// randomImage is a draw.Image backed by a slice of NRGBA pixels for use in tests.
type randomImage struct {
	w, h   int
	pixels []color.NRGBA
}

func (img *randomImage) ColorModel() color.Model { return color.NRGBAModel }
func (img *randomImage) Bounds() image.Rectangle { return image.Rect(0, 0, img.w, img.h) }
func (img *randomImage) At(x, y int) color.Color { return img.pixels[y*img.w+x] }
func (img *randomImage) Set(x, y int, c color.Color) {
	img.pixels[y*img.w+x] = color.NRGBAModel.Convert(c).(color.NRGBA)
}

// Verify randomImage implements draw.Image at compile time.
var _ draw.Image = (*randomImage)(nil)

func TestProperty_BT601GrayscaleConversion(t *testing.T) {

	rapid.Check(t, func(t *rapid.T) {
		r := rapid.Uint32Range(0, 255).Draw(t, "r")
		g := rapid.Uint32Range(0, 255).Draw(t, "g")
		b := rapid.Uint32Range(0, 255).Draw(t, "b")

		got := luminance(r, g, b)

		// Expected: min(255, (299*r + 587*g + 114*b) / 1000)
		y := (299*r + 587*g + 114*b) / 1000
		if y > 255 {
			y = 255
		}
		expected := byte(y)

		if got != expected {
			t.Fatalf("luminance(%d, %d, %d) = %d, want %d", r, g, b, got, expected)
		}
	})
}

func TestProperty_CustomI2CAddressUsage(t *testing.T) {

	rapid.Check(t, func(t *rapid.T) {
		// Generate random 7-bit address in range 0x01–0x7F (avoid 0 since 0 means "use default").
		addr := rapid.Uint16Range(0x01, 0x7F).Draw(t, "addr")

		bus := &mockBus{}
		cfg := driver.DriverConfig{
			Width:   15,
			Height:  7,
			I2CAddr: addr,
		}

		d, err := New(bus, cfg)
		if err != nil {
			t.Fatalf("New() returned unexpected error: %v", err)
		}

		if d.dev.Addr != addr {
			t.Fatalf("device address = 0x%02X, want 0x%02X", d.dev.Addr, addr)
		}
	})
}

func TestProperty_DrawImagePixelMappingAndClipping(t *testing.T) {

	rapid.Check(t, func(t *rapid.T) {
		// Generate image dimensions >= 15x7 (up to 50x50 to keep tests fast).
		imgW := rapid.IntRange(15, 50).Draw(t, "imgWidth")
		imgH := rapid.IntRange(7, 50).Draw(t, "imgHeight")

		// Generate random pixel colors for the entire image.
		pixels := make([]color.NRGBA, imgW*imgH)
		for i := range pixels {
			pixels[i] = color.NRGBA{
				R: uint8(rapid.Uint32Range(0, 255).Draw(t, "r")),
				G: uint8(rapid.Uint32Range(0, 255).Draw(t, "g")),
				B: uint8(rapid.Uint32Range(0, 255).Draw(t, "b")),
				A: 255,
			}
		}
		img := &randomImage{w: imgW, h: imgH, pixels: pixels}

		// Create a recording bus and construct the driver (init writes are recorded too).
		bus := &recordingBus{}
		cfg := driver.DriverConfig{Width: 15, Height: 7}
		d, err := New(bus, cfg)
		if err != nil {
			t.Fatalf("New() failed: %v", err)
		}

		// Clear recorded init writes so we only inspect DrawImage output.
		bus.writes = nil

		// Call DrawImage.
		if err := d.DrawImage(img); err != nil {
			t.Fatalf("DrawImage() failed: %v", err)
		}

		// DrawImage should produce exactly one Tx write.
		if len(bus.writes) != 1 {
			t.Fatalf("expected 1 Tx write from DrawImage, got %d", len(bus.writes))
		}

		write := bus.writes[0]

		// The write should be 1 register byte (0x24) + 105 brightness bytes = 106 total.
		if len(write) != 106 {
			t.Fatalf("expected write length 106, got %d", len(write))
		}

		// First byte must be PWM start register 0x24.
		if write[0] != 0x24 {
			t.Fatalf("expected register byte 0x24, got 0x%02X", write[0])
		}

		// Verify each of the 105 brightness values matches BT.601 luminance
		// of the corresponding pixel in the top-left 15x7 sub-rectangle.
		data := write[1:] // 105 brightness bytes
		for row := 0; row < 7; row++ {
			for col := 0; col < 15; col++ {
				idx := row*15 + col
				px := pixels[row*imgW+col] // top-left sub-rectangle pixel

				// Compute expected BT.601 luminance.
				y := (299*uint32(px.R) + 587*uint32(px.G) + 114*uint32(px.B)) / 1000
				if y > 255 {
					y = 255
				}
				expected := byte(y)

				if data[idx] != expected {
					t.Fatalf("pixel (%d,%d): brightness = %d, want %d (R=%d G=%d B=%d)",
						col, row, data[idx], expected, px.R, px.G, px.B)
				}
			}
		}
	})
}

// failAfterNBus is a mock I2C bus that succeeds for the first N Tx calls,
// then returns a configured error on all subsequent calls. Used to test
// error propagation from DrawImage after a successful init sequence.
type failAfterNBus struct {
	allowCount int // number of Tx calls to allow before failing
	callCount  int
	err        error
}

func (b *failAfterNBus) String() string                    { return "failAfterN" }
func (b *failAfterNBus) SetSpeed(_ physic.Frequency) error { return nil }
func (b *failAfterNBus) Tx(_ uint16, _, _ []byte) error {
	b.callCount++
	if b.callCount > b.allowCount {
		return b.err
	}
	return nil
}

// initTxCount is the number of Tx calls performed during IS31FL3731 init:
//   - writeReg(regPageSelect, pageFunctionReg)
//   - writeReg(regShutdown, 0x01)
//   - writeReg(regMode, 0x00)
//   - writeReg(regDisplayFrame, 0x00)
//   - writeReg(regPageSelect, pageFrame0)
//   - writeBulk(ledEnableStart, enableBuf)
//   - writeBulk(pwmStart, pwmBuf)
const initTxCount = 7

func TestProperty_DrawImageErrorPropagation(t *testing.T) {

	rapid.Check(t, func(t *rapid.T) {
		// Generate a random error message to inject as an I2C write failure.
		errMsg := rapid.StringMatching(`[a-z]{3,20}`).Draw(t, "err_msg")
		injectedErr := fmt.Errorf("i2c write: %s", errMsg)

		// Create a mock bus that succeeds during init (7 Tx calls) but fails
		// on the next Tx call, which is the DrawImage PWM write.
		bus := &failAfterNBus{
			allowCount: initTxCount,
			err:        injectedErr,
		}

		cfg := driver.DriverConfig{Width: 15, Height: 7}
		d, err := New(bus, cfg)
		if err != nil {
			t.Fatalf("New() failed during init (unexpected): %v", err)
		}

		// Generate a random 15x7+ image to draw.
		imgW := rapid.IntRange(15, 40).Draw(t, "imgWidth")
		imgH := rapid.IntRange(7, 20).Draw(t, "imgHeight")
		pixels := make([]color.NRGBA, imgW*imgH)
		for i := range pixels {
			pixels[i] = color.NRGBA{
				R: uint8(rapid.Uint32Range(0, 255).Draw(t, "r")),
				G: uint8(rapid.Uint32Range(0, 255).Draw(t, "g")),
				B: uint8(rapid.Uint32Range(0, 255).Draw(t, "b")),
				A: 255,
			}
		}
		img := &randomImage{w: imgW, h: imgH, pixels: pixels}

		// Call DrawImage — should fail with the injected error.
		drawErr := d.DrawImage(img)
		if drawErr == nil {
			t.Fatalf("DrawImage returned nil error, expected error containing %q", errMsg)
		}

		// Verify the returned error contains the injected error message.
		if !strings.Contains(drawErr.Error(), errMsg) {
			t.Fatalf("DrawImage error %q does not contain injected message %q", drawErr.Error(), errMsg)
		}
	})
}

// errBus is a mock I2C bus that always returns a configured error on Tx.
// Used to test initialization error format.
type errBus struct {
	err error
}

func (b *errBus) String() string                    { return "errBus" }
func (b *errBus) SetSpeed(_ physic.Frequency) error { return nil }
func (b *errBus) Tx(_ uint16, _, _ []byte) error    { return b.err }

func TestProperty_InitializationErrorIncludesAddressAndReason(t *testing.T) {

	rapid.Check(t, func(t *rapid.T) {
		// Generate a random 7-bit I2C address (0x01–0x7F).
		addr := rapid.Uint16Range(0x01, 0x7F).Draw(t, "addr")

		// Generate a random error reason string.
		reason := rapid.StringMatching(`[a-z]{3,30}`).Draw(t, "reason")
		busErr := fmt.Errorf("%s", reason)

		// Create a bus that always fails.
		bus := &errBus{err: busErr}

		cfg := driver.DriverConfig{
			Width:   15,
			Height:  7,
			I2CAddr: addr,
		}

		// Attempt to construct the driver — should fail during init.
		_, err := New(bus, cfg)
		if err == nil {
			t.Fatalf("New() returned nil error, expected initialization failure")
		}

		errStr := err.Error()

		// Verify the error contains the address in hex format (e.g., "0x74").
		expectedAddrStr := fmt.Sprintf("0x%02X", addr)
		if !strings.Contains(errStr, expectedAddrStr) {
			t.Fatalf("error %q does not contain address %q", errStr, expectedAddrStr)
		}

		// Verify the error contains the underlying reason.
		if !strings.Contains(errStr, reason) {
			t.Fatalf("error %q does not contain reason %q", errStr, reason)
		}
	})
}

func TestProperty_BoundsFromConfigDimensions(t *testing.T) {

	rapid.Check(t, func(t *rapid.T) {
		w := rapid.IntRange(1, 16).Draw(t, "width")
		h := rapid.IntRange(1, 9).Draw(t, "height")

		d, err := New(&mockBus{}, driver.DriverConfig{Width: w, Height: h})
		if err != nil {
			t.Fatalf("New() returned unexpected error: %v", err)
		}

		expected := image.Rect(0, 0, w, h)
		got := d.Bounds()
		if got != expected {
			t.Fatalf("Bounds() = %v, want %v (width=%d, height=%d)", got, expected, w, h)
		}
	})
}

func TestProperty_DrawImagePixelCountMatchesConfig(t *testing.T) {

	rapid.Check(t, func(t *rapid.T) {
		w := rapid.IntRange(1, 16).Draw(t, "width")
		h := rapid.IntRange(1, 9).Draw(t, "height")

		bus := &recordingBus{}
		cfg := driver.DriverConfig{Width: w, Height: h}

		d, err := New(bus, cfg)
		if err != nil {
			t.Fatalf("New() failed: %v", err)
		}

		// Clear init writes so we only inspect DrawImage output.
		bus.writes = nil

		// Create a random image of exactly w×h pixels.
		pixels := make([]color.NRGBA, w*h)
		for i := range pixels {
			pixels[i] = color.NRGBA{
				R: uint8(rapid.Uint32Range(0, 255).Draw(t, "r")),
				G: uint8(rapid.Uint32Range(0, 255).Draw(t, "g")),
				B: uint8(rapid.Uint32Range(0, 255).Draw(t, "b")),
				A: 255,
			}
		}
		img := &randomImage{w: w, h: h, pixels: pixels}

		if err := d.DrawImage(img); err != nil {
			t.Fatalf("DrawImage() failed: %v", err)
		}

		// DrawImage should produce exactly one Tx write.
		if len(bus.writes) != 1 {
			t.Fatalf("expected 1 Tx write from DrawImage, got %d", len(bus.writes))
		}

		// The write should be 1 register byte + w*h brightness bytes.
		expectedLen := w*h + 1
		if len(bus.writes[0]) != expectedLen {
			t.Fatalf("expected write length %d (1 register + %d pixels), got %d",
				expectedLen, w*h, len(bus.writes[0]))
		}
	})
}

func TestProperty_LEDEnableMaskMatchesTotalLEDs(t *testing.T) {

	rapid.Check(t, func(t *rapid.T) {
		w := rapid.IntRange(1, 16).Draw(t, "width")
		h := rapid.IntRange(1, 9).Draw(t, "height")

		// Filter: only valid dimensions where total LEDs fit in 144 bits (18 bytes).
		if w*h > 144 {
			t.Skip("w*h exceeds 144")
		}

		bus := &recordingBus{}
		_, err := New(bus, driver.DriverConfig{Width: w, Height: h})
		if err != nil {
			t.Fatalf("New() returned unexpected error: %v", err)
		}

		// The LED-enable write is the 6th Tx call during init (index 5).
		// It starts with ledEnableStart (0x00) followed by 18 data bytes (total length 19).
		if len(bus.writes) < 6 {
			t.Fatalf("expected at least 6 Tx writes during init, got %d", len(bus.writes))
		}

		enableWrite := bus.writes[5]
		if len(enableWrite) != 19 {
			t.Fatalf("expected LED-enable write length 19, got %d", len(enableWrite))
		}
		if enableWrite[0] != 0x00 {
			t.Fatalf("expected LED-enable register byte 0x00, got 0x%02X", enableWrite[0])
		}

		// Count set bits in the 18 data bytes (skip register byte).
		setBits := 0
		for _, b := range enableWrite[1:] {
			setBits += bits.OnesCount8(b)
		}

		expected := w * h
		if setBits != expected {
			t.Fatalf("LED enable set bits = %d, want %d (width=%d, height=%d)", setBits, expected, w, h)
		}
	})
}

// --- Named pixel-address layout tests (Adafruit CharliePlex boards) ---

// charlie-wing reference mapping, from Adafruit's CircuitPython driver:
//
//	if x > 7: x = 15 - x; y += 8  else: y = 7 - y
//	addr = x*16 + y
func wingRefAddr(x, y int) int {
	if x > 7 {
		x = 15 - x
		y += 8
	} else {
		y = 7 - y
	}
	return x*16 + y
}

// charlie-bonnet reference mapping, from Adafruit's CircuitPython driver:
//
//	if x >= 8: addr = (x-6)*16 - (y+1)  else: addr = (x+1)*16 + (7-y)
func bonnetRefAddr(x, y int) int {
	if x >= 8 {
		return (x-6)*16 - (y + 1)
	}
	return (x+1)*16 + (7 - y)
}

func TestLayoutCharlieWingMapping(t *testing.T) {
	fn, err := resolvePixelAddr(driver.DriverConfig{Width: 15, Height: 7, Layout: "charlie-wing"})
	if err != nil {
		t.Fatalf("resolvePixelAddr failed: %v", err)
	}
	seen := map[int]bool{}
	for y := 0; y < 7; y++ {
		for x := 0; x < 15; x++ {
			got := fn(x, y)
			want := wingRefAddr(x, y)
			if got != want {
				t.Fatalf("charlie-wing addr(%d,%d) = %d, want %d", x, y, got, want)
			}
			if got < 0 || got >= 144 {
				t.Fatalf("charlie-wing addr(%d,%d) = %d out of range [0,144)", x, y, got)
			}
			if seen[got] {
				t.Fatalf("charlie-wing addr(%d,%d) = %d collides with another pixel", x, y, got)
			}
			seen[got] = true
		}
	}
}

func TestLayoutCharlieBonnetMapping(t *testing.T) {
	fn, err := resolvePixelAddr(driver.DriverConfig{Width: 16, Height: 8, Layout: "charlie-bonnet"})
	if err != nil {
		t.Fatalf("resolvePixelAddr failed: %v", err)
	}
	seen := map[int]bool{}
	for y := 0; y < 8; y++ {
		for x := 0; x < 16; x++ {
			got := fn(x, y)
			want := bonnetRefAddr(x, y)
			if got != want {
				t.Fatalf("charlie-bonnet addr(%d,%d) = %d, want %d", x, y, got, want)
			}
			if got < 0 || got >= 144 {
				t.Fatalf("charlie-bonnet addr(%d,%d) = %d out of range [0,144)", x, y, got)
			}
			if seen[got] {
				t.Fatalf("charlie-bonnet addr(%d,%d) = %d collides with another pixel", x, y, got)
			}
			seen[got] = true
		}
	}
}

func TestLayoutCharlieBonnetTransposedMapping(t *testing.T) {
	// Bonnet mounted rotated 90 degrees: configured 8x16 portrait. Every logical
	// coordinate must map to a unique in-range address via the rotated formula.
	fn, err := resolvePixelAddr(driver.DriverConfig{Width: 8, Height: 16, Layout: "charlie-bonnet"})
	if err != nil {
		t.Fatalf("resolvePixelAddr failed: %v", err)
	}
	seen := map[int]bool{}
	for y := 0; y < 16; y++ {
		for x := 0; x < 8; x++ {
			got := fn(x, y)
			want := bonnetRefAddr(y, 8-1-x)
			if got != want {
				t.Fatalf("rotated addr(%d,%d) = %d, want %d", x, y, got, want)
			}
			if got < 0 || got >= 144 {
				t.Fatalf("rotated addr(%d,%d) = %d out of range [0,144)", x, y, got)
			}
			if seen[got] {
				t.Fatalf("rotated addr(%d,%d) = %d collides with another pixel", x, y, got)
			}
			seen[got] = true
		}
	}
}

func TestLayoutUnknownNameFails(t *testing.T) {
	_, err := New(&mockBus{}, driver.DriverConfig{Width: 15, Height: 7, Layout: "nonsense"})
	if err == nil || !strings.Contains(err.Error(), "unknown layout") {
		t.Fatalf("New() error = %v, want unknown layout error", err)
	}
}

func TestLayoutDimensionMismatchFails(t *testing.T) {
	_, err := New(&mockBus{}, driver.DriverConfig{Width: 15, Height: 7, Layout: "charlie-bonnet"})
	if err == nil || !strings.Contains(err.Error(), "requires") {
		t.Fatalf("New() error = %v, want dimension mismatch error", err)
	}
}

func TestBonnetDrawImageWritesMappedPWM(t *testing.T) {
	bus := &recordingBus{}
	d, err := New(bus, driver.DriverConfig{Width: 16, Height: 8, Layout: "charlie-bonnet"})
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	bus.writes = nil

	img := image.NewRGBA(image.Rect(0, 0, 16, 8))
	img.SetRGBA(0, 0, color.RGBA{255, 255, 255, 255})
	img.SetRGBA(15, 7, color.RGBA{255, 255, 255, 255})
	img.SetRGBA(9, 3, color.RGBA{255, 255, 255, 255})

	if err := d.DrawImage(img); err != nil {
		t.Fatalf("DrawImage() failed: %v", err)
	}
	if len(bus.writes) != 1 {
		t.Fatalf("expected 1 write, got %d", len(bus.writes))
	}
	w := bus.writes[0]
	if len(w) != 1+144 {
		t.Fatalf("expected full 144-byte PWM write (145 with register), got %d", len(w))
	}
	if w[0] != pwmStart {
		t.Fatalf("expected register byte 0x%02X, got 0x%02X", pwmStart, w[0])
	}
	data := w[1:]
	lit := map[int]bool{
		bonnetRefAddr(0, 0):  true,
		bonnetRefAddr(15, 7): true,
		bonnetRefAddr(9, 3):  true,
	}
	for addr, v := range data {
		if lit[addr] && v != 255 {
			t.Fatalf("addr %d: brightness %d, want 255", addr, v)
		}
		if !lit[addr] && v != 0 {
			t.Fatalf("addr %d: brightness %d, want 0", addr, v)
		}
	}
}

func TestBonnetEnableBitmapMatchesMapping(t *testing.T) {
	bus := &recordingBus{}
	_, err := New(bus, driver.DriverConfig{Width: 16, Height: 8, Layout: "charlie-bonnet"})
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	if len(bus.writes) < 6 {
		t.Fatalf("expected at least 6 init writes, got %d", len(bus.writes))
	}
	enable := bus.writes[5]
	if len(enable) != 19 || enable[0] != 0x00 {
		t.Fatalf("unexpected enable write: len=%d reg=0x%02X", len(enable), enable[0])
	}
	want := make([]byte, 18)
	for y := 0; y < 8; y++ {
		for x := 0; x < 16; x++ {
			a := bonnetRefAddr(x, y)
			want[a/8] |= 1 << (a % 8)
		}
	}
	for i := range want {
		if enable[1+i] != want[i] {
			t.Fatalf("enable byte %d = 0x%02X, want 0x%02X", i, enable[1+i], want[i])
		}
	}
}
