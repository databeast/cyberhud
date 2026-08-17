package displaymode

import (
	"fmt"
	"reflect"
	"testing"

	"pgregory.net/rapid"
)

// =============================================================================
//
// For any ViewData value, if allEmptyItems(vd.Items) is true then
// convertViewData(vd).Items equals []string{"(template)"}, and if
// allEmptyItems(vd.Items) is false then convertViewData(vd).Items equals
// vd.Items unchanged.

// =============================================================================

func TestPropertyViewDataConversionEmptyItemsGuard(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate a ViewData with random Items
		choice := rapid.IntRange(0, 3).Draw(t, "choice")
		var items []string
		switch choice {
		case 0:
			// nil slice
			items = nil
		case 1:
			// empty slice
			items = []string{}
		case 2:
			// slice of only empty strings
			count := rapid.IntRange(1, 5).Draw(t, "emptyCount")
			items = make([]string, count)
		case 3:
			// slice with at least one non-empty string
			count := rapid.IntRange(1, 5).Draw(t, "itemCount")
			items = make([]string, count)
			for i := range items {
				items[i] = rapid.String().Draw(t, fmt.Sprintf("item%d", i))
			}
			// Ensure at least one is non-empty
			hasNonEmpty := false
			for _, item := range items {
				if item != "" {
					hasNonEmpty = true
					break
				}
			}
			if !hasNonEmpty {
				items[0] = "nonempty"
			}
		}

		vd := ViewData{Items: items}
		result := convertViewData(vd)

		if allEmptyItems(items) {
			expected := []string{"(template)"}
			if !reflect.DeepEqual(result.Items, expected) {
				t.Fatalf("allEmptyItems=true: expected %v, got %v", expected, result.Items)
			}
		} else {
			if !reflect.DeepEqual(result.Items, items) {
				t.Fatalf("allEmptyItems=false: expected items %v preserved, got %v", items, result.Items)
			}
		}
	})
}
