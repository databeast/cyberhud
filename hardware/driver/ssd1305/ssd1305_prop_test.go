package ssd1305

import (
	"image"
	"testing"

	conn "periph.io/x/conn/v3"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"pgregory.net/rapid"

	"github.com/databeast/cyberhud/hardware/driver"
)

// mockSPIConn is a minimal spi.Conn implementation that succeeds on all operations.
type mockSPIConn struct{}

func (m *mockSPIConn) String() string                 { return "mockSPIConn" }
func (m *mockSPIConn) Halt() error                    { return nil }
func (m *mockSPIConn) Tx(_, _ []byte) error           { return nil }
func (m *mockSPIConn) Duplex() conn.Duplex            { return conn.Full }
func (m *mockSPIConn) TxPackets(_ []spi.Packet) error { return nil }

// mockSPIPort is a minimal spi.Port that returns a mockSPIConn on Connect.
type mockSPIPort struct{}

func (m *mockSPIPort) String() string { return "mockSPIPort" }
func (m *mockSPIPort) Connect(_ physic.Frequency, _ spi.Mode, _ int) (spi.Conn, error) {
	return &mockSPIConn{}, nil
}

// mockPin is a minimal gpio.PinOut implementation that succeeds on all operations.
type mockPin struct{ name string }

func (p *mockPin) String() string                            { return p.name }
func (p *mockPin) Halt() error                               { return nil }
func (p *mockPin) Name() string                              { return p.name }
func (p *mockPin) Number() int                               { return 0 }
func (p *mockPin) Function() string                          { return "Out" }
func (p *mockPin) Out(_ gpio.Level) error                    { return nil }
func (p *mockPin) PWM(_ gpio.Duty, _ physic.Frequency) error { return nil }

// mockI2CBus is a minimal i2c.Bus implementation that succeeds on all operations.
type mockI2CBus struct{}

func (m *mockI2CBus) String() string                    { return "mockI2CBus" }
func (m *mockI2CBus) Tx(_ uint16, _, _ []byte) error    { return nil }
func (m *mockI2CBus) SetSpeed(_ physic.Frequency) error { return nil }

func TestProperty1_BoundsCorrectnessForValidDriverConfig(t *testing.T) {

	rapid.Check(t, func(t *rapid.T) {
		width := rapid.IntRange(1, 512).Draw(t, "width")
		// Height must be a multiple of 8 for page-column format.
		heightMult := rapid.IntRange(1, 16).Draw(t, "heightMult")
		height := heightMult * 8

		cfg := driver.DriverConfig{
			Width:  width,
			Height: height,
		}

		expectedBounds := image.Rect(0, 0, width, height)

		// Test via NewSPI path.
		port := &mockSPIPort{}
		dc := &mockPin{name: "DC"}
		rst := &mockPin{name: "RST"}
		bl := &mockPin{name: "BL"}

		spiDev, err := NewSPI(port, dc, rst, bl, cfg)
		if err != nil {
			t.Fatalf("NewSPI() returned unexpected error: %v", err)
		}
		if got := spiDev.Bounds(); got != expectedBounds {
			t.Fatalf("SPI Bounds() = %v, want %v", got, expectedBounds)
		}

		// Test via NewI2C path.
		bus := &mockI2CBus{}

		i2cDev, err := NewI2C(bus, cfg)
		if err != nil {
			t.Fatalf("NewI2C() returned unexpected error: %v", err)
		}
		if got := i2cDev.Bounds(); got != expectedBounds {
			t.Fatalf("I2C Bounds() = %v, want %v", got, expectedBounds)
		}
	})
}
