package tests_test

import (
	"testing"

	"github.com/databeast/cyberhud/display/modes/menu"
)

func TestItems(t *testing.T) {
	items := menu.Items()
	if len(items) != 6 {
		t.Fatalf("Items() len=%d, want 6", len(items))
	}
	if items[0] != "STEMMA QT / QWIIC" || items[1] != "GPIO Pins" || items[2] != "GPIO Control" || items[3] != "USB Bench" || items[4] != "Serial Monitor" || items[5] != "System" {
		t.Fatalf("Items()=%v, unexpected", items)
	}
}
