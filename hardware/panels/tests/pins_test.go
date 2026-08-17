package tests_test

import (
	"strings"
	"testing"

	"github.com/databeast/cyberhud/hardware/panels"
	_ "github.com/databeast/cyberhud/hardware/panels/all"
)

func TestPinNoticesWaveshare13HatGPIO13Unavailable(t *testing.T) {
	p, err := panels.Get("waveshare-1.3hat")
	if err != nil {
		t.Fatalf("Get() error=%v", err)
	}
	notices := panels.PinNotices(p)
	if len(notices) == 0 {
		t.Fatal("expected at least one pin notice")
	}
	if !strings.Contains(strings.Join(notices, "\n"), "GPIO13") {
		t.Fatalf("expected GPIO13 unavailable notice, got %v", notices)
	}
}

func TestPinNoticesTripleScreenGPIO13AndGPIO18Unavailable(t *testing.T) {
	p, err := panels.Get("waveshare-triple-screen")
	if err != nil {
		t.Fatalf("Get() error=%v", err)
	}
	notices := panels.PinNotices(p)
	joined := strings.Join(notices, "\n")
	if !strings.Contains(joined, "GPIO13") || !strings.Contains(joined, "GPIO18") {
		t.Fatalf("expected GPIO13 and GPIO18 unavailable notices, got %q", joined)
	}
}
