package tests

import (
	"reflect"
	"testing"

	"github.com/databeast/cyberhud/display/style"
)

// TestStyleContextFieldsUnexported confirms that all fields on StyleContext
// are unexported, enforcing the information boundary at the type level.
//

func TestStyleContextFieldsUnexported(t *testing.T) {
	typ := reflect.TypeOf(style.StyleContext{})

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.IsExported() {
			t.Errorf("field %q is exported; all StyleContext fields must be unexported", field.Name)
		}
	}
}

// TestStyleContextNoIntMethodReturns confirms that no zero-parameter method
// on StyleContext returns a bare int (not a named type), which could represent
// raw pixel dimensions. Named int-based types like Capability are allowed.
//

func TestStyleContextNoIntMethodReturns(t *testing.T) {
	typ := reflect.TypeOf(style.StyleContext{})

	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)
		// Only check methods with no arguments beyond the receiver.
		if method.Type.NumIn() != 1 {
			continue
		}
		for j := 0; j < method.Type.NumOut(); j++ {
			out := method.Type.Out(j)
			// Only flag bare int (unnamed int type). Named int-based types
			// (e.g., Capability) are intentional and allowed.
			if out.Kind() == reflect.Int && out.Name() == "int" {
				t.Errorf("method %q returns bare int — StyleContext must not expose raw pixel dimension integers", method.Name)
			}
		}
	}
}
