// Package coordinator tracks per-region display mode state and provides
// Set / Next / Prev operations for remote-controlling which mode is active on
// each configured display region. It depends on the display/catalog package for
// enriching region metadata with mode definitions via catalog.DescribeMany.
package coordinator
