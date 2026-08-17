// Package input handles the three push-buttons and five-way joystick on the
// Waveshare 1.3inch LCD HAT connected to a Raspberry Pi via GPIO.
//
// Default GPIO mapping (Waveshare 1.3inch LCD HAT):
//
//	KEY1    GPIO5  (pin 29)
//	KEY2    GPIO6  (pin 31)
//	KEY3    GPIO13 (pin 33)
//	Joy-Up  GPIO21 (pin 40)
//	Joy-Dn  GPIO19 (pin 35)
//	Joy-Lt  GPIO16 (pin 36)
//	Joy-Rt  GPIO20 (pin 38)
//	Joy-Pr  GPIO26 (pin 37)
//
// All inputs are active-low (pulled high when idle).
package input

import (
	"sync"
	"time"

	"periph.io/x/conn/v3/gpio"
)

// Event represents a single button or joystick action.
type Event struct {
	Key  Key
	Type EventType
}

// EventType distinguishes a press from a release.
type EventType uint8

const (
	// Press is emitted when a key transitions from idle to active.
	Press EventType = iota
	// Release is emitted when a key transitions from active to idle.
	Release
)

// Key identifies a single input source.
type Key uint8

const (
	Key1       Key = iota // Left hardware button
	Key2                  // Middle hardware button
	Key3                  // Right hardware button
	JoyUp                 // Joystick pushed up
	JoyDown               // Joystick pushed down
	JoyLeft               // Joystick pushed left
	JoyRight              // Joystick pushed right
	JoyPressed            // Joystick pressed inward
)

func (k Key) String() string {
	switch k {
	case Key1:
		return "KEY1"
	case Key2:
		return "KEY2"
	case Key3:
		return "KEY3"
	case JoyUp:
		return "JOY_UP"
	case JoyDown:
		return "JOY_DOWN"
	case JoyLeft:
		return "JOY_LEFT"
	case JoyRight:
		return "JOY_RIGHT"
	case JoyPressed:
		return "JOY_PRESS"
	default:
		return "UNKNOWN"
	}
}

// Handler polls all configured GPIO pins and emits Events over a channel.
type Handler struct {
	pins     []pinEntry
	events   chan Event
	stopCh   chan struct{}
	wg       sync.WaitGroup
	debounce time.Duration
}

type pinEntry struct {
	pin  gpio.PinIn
	key  Key
	last gpio.Level
}

// Config holds the GPIO pin assignment for each input.
type Config struct {
	Key1       gpio.PinIn
	Key2       gpio.PinIn
	Key3       gpio.PinIn
	JoyUp      gpio.PinIn
	JoyDown    gpio.PinIn
	JoyLeft    gpio.PinIn
	JoyRight   gpio.PinIn
	JoyPressed gpio.PinIn
}

// New creates a Handler from c.  Nil entries in c are silently ignored.
// debounce sets the debounce interval between polls.
func New(c Config, debounce time.Duration) (*Handler, error) {
	if debounce <= 0 {
		debounce = 20 * time.Millisecond
	}

	type kv struct {
		pin gpio.PinIn
		key Key
	}
	candidates := []kv{
		{c.Key1, Key1},
		{c.Key2, Key2},
		{c.Key3, Key3},
		{c.JoyUp, JoyUp},
		{c.JoyDown, JoyDown},
		{c.JoyLeft, JoyLeft},
		{c.JoyRight, JoyRight},
		{c.JoyPressed, JoyPressed},
	}

	var pins []pinEntry
	for _, kv := range candidates {
		if kv.pin == nil {
			continue
		}
		// Configure as input with pull-up; all signals are active-low.
		if err := kv.pin.In(gpio.PullUp, gpio.NoEdge); err != nil {
			return nil, err
		}
		pins = append(pins, pinEntry{
			pin:  kv.pin,
			key:  kv.key,
			last: gpio.High, // idle state
		})
	}

	return &Handler{
		pins:     pins,
		events:   make(chan Event, 32),
		stopCh:   make(chan struct{}),
		debounce: debounce,
	}, nil
}

// Events returns the channel on which input events are delivered.
func (h *Handler) Events() <-chan Event {
	return h.events
}

// Start begins background polling of all GPIO pins.
func (h *Handler) Start() {
	h.wg.Add(1)
	go h.poll()
}

// Stop shuts down the polling goroutine and closes the events channel.
func (h *Handler) Stop() {
	close(h.stopCh)
	h.wg.Wait()
	close(h.events)
}

func (h *Handler) poll() {
	defer h.wg.Done()
	ticker := time.NewTicker(h.debounce)
	defer ticker.Stop()

	for {
		select {
		case <-h.stopCh:
			return
		case <-ticker.C:
			h.sample()
		}
	}
}

func (h *Handler) sample() {
	for i := range h.pins {
		current := h.pins[i].pin.Read()
		if current == h.pins[i].last {
			continue
		}
		h.pins[i].last = current
		evt := Event{Key: h.pins[i].key}
		if current == gpio.Low {
			evt.Type = Press
		} else {
			evt.Type = Release
		}
		select {
		case h.events <- evt:
		default:
			// drop if consumer is slow
		}
	}
}
