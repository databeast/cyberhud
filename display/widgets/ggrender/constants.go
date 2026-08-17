package ggrender

// MinBoundsWidth is the minimum pixel width for gg-based rendering.
// Render functions return nil for bounds narrower than this threshold.
const MinBoundsWidth = 16

// MinBoundsHeight is the minimum pixel height for gg-based rendering.
// Render functions return nil for bounds shorter than this threshold.
const MinBoundsHeight = 16

// maxLabelLen is the maximum number of runes allowed in a widget label
// before truncation occurs.
const maxLabelLen = 128
