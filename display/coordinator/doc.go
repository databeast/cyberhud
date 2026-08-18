// Package coordinator is a compatibility facade over live region state.
// It exposes Set / Next / Prev operations and region snapshots for remote
// control, but delegates to display/region when the daemon has bound the live
// RegionManager. The display/catalog package is used for metadata enrichment.
package coordinator
