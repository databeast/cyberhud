package ssd1305_test

import (
	"strings"
	"testing"

	conn "periph.io/x/conn/v3"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"

	"github.com/databeast/cyberhud/hardware/driver"
	_ "github.com/databeast/cyberhud/hardware/driver/ssd1305"
)

// --- Mock types ---

// mockPinOut implements gpio.PinOut for testing.
type mockPinOut struct {
	name string
	last gpio.Level
}

func (p *mockPinOut) String() string                        { return p.name }
func (p *mockPinOut) Halt() error                           { return nil }
func (p *mockPinOut) Name() string                          { return p.name }
func (p *mockPinOut) Number() int                           { return -1 }
func (p *mockPinOut) Function() string                      { return "Out" }
func (p *mockPinOut) Out(l gpio.Level) error                { p.last = l; return nil }
func (p *mockPinOut) PWM(gpio.Duty, physic.Frequency) error { return nil }

// mockSPIConn implements spi.Conn and captures all transmitted bytes.
type mockSPIConn struct {
	txBuf []byte
}

func (c *mockSPIConn) String() string      { return "mock-spi-conn" }
func (c *mockSPIConn) Duplex() conn.Duplex { return conn.Half }
func (c *mockSPIConn) Tx(w, r []byte) error {
	c.txBuf = append(c.txBuf, w...)
	return nil
}
func (c *mockSPIConn) TxPackets(p []spi.Packet) error { return nil }

// mockSPIPort implements spi.Port and returns the embedded mockSPIConn.
type mockSPIPort struct {
	conn *mockSPIConn
}

func (p *mockSPIPort) String() string { return "mock-spi-port" }
func (p *mockSPIPort) Connect(f physic.Frequency, mode spi.Mode, bits int) (spi.Conn, error) {
	return p.conn, nil
}

// mockI2CBus implements i2c.Bus for testing.
type mockI2CBus struct {
	txBuf []byte
}

func (b *mockI2CBus) String() string { return "mock-i2c-bus" }
func (b *mockI2CBus) Tx(addr uint16, w, r []byte) error {
	b.txBuf = append(b.txBuf, w...)
	return nil
}
func (b *mockI2CBus) SetSpeed(f physic.Frequency) error { return nil }

// --- Tests ---

func TestDriverRegistration(t *testing.T) {
	def, ok := driver.Get("ssd1305")
	if !ok {
		t.Fatal("driver.Get(\"ssd1305\") returned false; expected driver to be registered")
	}
	if def.ID != "ssd1305" {
		t.Fatalf("expected ID %q, got %q", "ssd1305", def.ID)
	}
	if def.NewSPI == nil {
		t.Fatal("expected NewSPI factory to be non-nil")
	}
	if def.NewI2C == nil {
		t.Fatal("expected NewI2C factory to be non-nil")
	}
	if !def.Monochrome {
		t.Fatal("expected Monochrome=true")
	}
}

func TestNewSPI_NilPort(t *testing.T) {
	dc := &mockPinOut{name: "DC"}
	rst := &mockPinOut{name: "RST"}
	bl := &mockPinOut{name: "BL"}
	cfg := driver.DriverConfig{Width: 128, Height: 32}

	def, ok := driver.Get("ssd1305")
	if !ok {
		t.Fatal("driver not registered")
	}

	_, err := def.NewSPI(nil, dc, rst, bl, nil, cfg)
	if err == nil {
		t.Fatal("expected error for nil SPI port")
	}
	if !strings.Contains(err.Error(), "spi port") {
		t.Fatalf("error should mention 'spi port', got: %q", err.Error())
	}
}

func TestNewSPI_NilDCPin(t *testing.T) {
	port := &mockSPIPort{conn: &mockSPIConn{}}
	rst := &mockPinOut{name: "RST"}
	bl := &mockPinOut{name: "BL"}
	cfg := driver.DriverConfig{Width: 128, Height: 32}

	def, ok := driver.Get("ssd1305")
	if !ok {
		t.Fatal("driver not registered")
	}

	_, err := def.NewSPI(port, nil, rst, bl, nil, cfg)
	if err == nil {
		t.Fatal("expected error for nil DC pin")
	}
	if !strings.Contains(err.Error(), "dc pin") {
		t.Fatalf("error should mention 'dc pin', got: %q", err.Error())
	}
}

func TestNewSPI_NilRSTPin(t *testing.T) {
	port := &mockSPIPort{conn: &mockSPIConn{}}
	dc := &mockPinOut{name: "DC"}
	bl := &mockPinOut{name: "BL"}
	cfg := driver.DriverConfig{Width: 128, Height: 32}

	def, ok := driver.Get("ssd1305")
	if !ok {
		t.Fatal("driver not registered")
	}

	_, err := def.NewSPI(port, dc, nil, bl, nil, cfg)
	if err == nil {
		t.Fatal("expected error for nil RST pin")
	}
	if !strings.Contains(err.Error(), "rst pin") {
		t.Fatalf("error should mention 'rst pin', got: %q", err.Error())
	}
}

func TestNewI2C_NilBus(t *testing.T) {
	cfg := driver.DriverConfig{Width: 128, Height: 32}

	def, ok := driver.Get("ssd1305")
	if !ok {
		t.Fatal("driver not registered")
	}

	_, err := def.NewI2C(nil, cfg)
	if err == nil {
		t.Fatal("expected error for nil I2C bus")
	}
	if !strings.Contains(err.Error(), "i2c bus") {
		t.Fatalf("error should mention 'i2c bus', got: %q", err.Error())
	}
}

func TestInitSequence_SPI(t *testing.T) {
	conn := &mockSPIConn{}
	port := &mockSPIPort{conn: conn}
	dc := &mockPinOut{name: "DC"}
	rst := &mockPinOut{name: "RST"}
	bl := &mockPinOut{name: "BL"}
	cfg := driver.DriverConfig{Width: 128, Height: 32}

	def, ok := driver.Get("ssd1305")
	if !ok {
		t.Fatal("driver not registered")
	}

	_, err := def.NewSPI(port, dc, rst, bl, nil, cfg)
	if err != nil {
		t.Fatalf("NewSPI() unexpected error: %v", err)
	}

	// Expected initialization command sequence (each byte sent individually via spiCmd).
	// Default config has Rotate180=false, so uses standard orientation (0xA1, 0xC8).
	expected := []byte{
		0xAE,       // Display OFF
		0x04,       // Set lower column start address (offset 4)
		0x10,       // Set higher column start address
		0x40,       // Set display start line to 0
		0x81, 0x80, // Set contrast
		0xA1,       // Set segment remap (col 127 → SEG0)
		0xA6,       // Set normal display (not inverted)
		0xA8, 0x1F, // Set multiplex ratio to 31 (32 rows)
		0xC8,       // Set COM scan direction: remapped
		0xD3, 0x00, // Set display offset to 0
		0xD5, 0xF0, // Set clock divide ratio/oscillator frequency
		0xD8, 0x05, // Set area color mode / low power display mode
		0xD9, 0xC2, // Set pre-charge period
		0xDA, 0x12, // Set COM pins hardware configuration
		0xDB, 0x08, // Set VCOMH deselect level
		0xAF, // Display ON
	}

	// The mock captures all Tx bytes. The init sequence sends each command byte
	// individually (one byte per Tx call), so conn.txBuf should contain exactly
	// the expected bytes in order.
	if len(conn.txBuf) != len(expected) {
		t.Fatalf("init sequence length: got %d bytes, want %d bytes\ngot:  %#v\nwant: %#v",
			len(conn.txBuf), len(expected), conn.txBuf, expected)
	}

	for i, b := range expected {
		if conn.txBuf[i] != b {
			t.Fatalf("init sequence mismatch at byte %d: got 0x%02X, want 0x%02X\nfull got:  %#v\nfull want: %#v",
				i, conn.txBuf[i], b, conn.txBuf, expected)
		}
	}
}
