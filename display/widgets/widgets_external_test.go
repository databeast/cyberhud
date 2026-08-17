package widgets_test

import (
	"reflect"
	"testing"

	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/textlabel"
)

// ---------------------------------------------------------------------------
// Source: widgets_test.go (external test package)
// ---------------------------------------------------------------------------

// TestWidgetInterfaceCompiles verifies that the widgets.Widget interface exists
// and is a valid Go interface type. This is primarily a compile-time check.
//

func TestWidgetInterfaceCompiles(t *testing.T) {
	// Compile-time assertion: Widget is an interface type
	var _ = (widgets.Widget)(nil)

	// Runtime reflection check
	iface := reflect.TypeOf((*widgets.Widget)(nil)).Elem()
	if iface.Kind() != reflect.Interface {
		t.Fatalf("widgets.Widget is %v, want interface", iface.Kind())
	}
}

// TestTextlabelConfigLacksTextboxFields verifies that textlabel.Config does NOT
// include textbox-specific fields: Overflow, VAlign, LineSpacing, PadX, PadY,
// Border, FontOverrides.
//

func TestTextlabelConfigLacksTextboxFields(t *testing.T) {
	typ := reflect.TypeOf(textlabel.Config{})

	forbidden := []string{
		"Overflow",
		"VAlign",
		"LineSpacing",
		"PadX",
		"PadY",
		"Border",
		"FontOverrides",
	}

	for _, name := range forbidden {
		_, found := typ.FieldByName(name)
		if found {
			t.Errorf("textlabel.Config should NOT have field %q (textbox-specific)", name)
		}
	}
}
