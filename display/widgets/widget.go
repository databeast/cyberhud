package widgets

// Widget defines the common rendering contract for documentation purposes.
// It is NOT used for polymorphic dispatch at runtime because Config types
// differ between packages. Use the Renderable interface for runtime dispatch.
//
// Conforming packages:
//   - borderframe
//   - gradient
//   - led
//   - progressbar
//   - scaledtextbox
//   - scrollbar
//   - sparkline
//   - textbox
//   - textlabel
//
// Non-conforming packages (different patterns):
//   - icons (asset registry)
type Widget interface {
	// Render is documented here to show the expected contract.
	// Actual implementations use package-specific Config types.
}
