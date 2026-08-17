package widgets

// SuppressionContext carries panel properties used by suppression rules
// to determine whether a widget should be skipped during composition.
type SuppressionContext struct {
	IsEink          bool
	AvailableWidth  int
	AvailableHeight int
}

// SuppressionRule evaluates whether a widget should be skipped based on
// its Descriptor and the current panel context. Returns true to suppress.
type SuppressionRule func(Descriptor, SuppressionContext) bool

// SuppressOnEink returns a rule that suppresses widgets lacking "eink-safe"
// in their Capabilities when the panel is an e-ink display.
func SuppressOnEink() SuppressionRule {
	return func(desc Descriptor, ctx SuppressionContext) bool {
		if !ctx.IsEink {
			return false
		}
		for _, cap := range desc.Capabilities {
			if cap == "eink-safe" {
				return false
			}
		}
		return true
	}
}

// SuppressBelow returns a rule that suppresses widgets whose minimum
// bounds exceed the given available dimensions. The widget is suppressed
// when its MinWidth > minWidth or MinHeight > minHeight.
func SuppressBelow(minWidth, minHeight int) SuppressionRule {
	return func(desc Descriptor, _ SuppressionContext) bool {
		return desc.MinWidth > minWidth || desc.MinHeight > minHeight
	}
}
