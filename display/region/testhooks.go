package region

import (
	"strings"
	"time"

	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// TestPerRegionTicker exposes per-region ticker state for the external tests package.
type TestPerRegionTicker struct {
	Region   *Region
	Interval time.Duration
	LastFire time.Time
}

// TestValidateSpec exposes validateSpec to the external tests package.
func (rm *RegionManager) TestValidateSpec(spec RegionSpec) error {
	return rm.validateSpec(spec)
}

// TestValidateSpecCore exposes validateSpecCore to the external tests package.
func (rm *RegionManager) TestValidateSpecCore(spec RegionSpec) error {
	return rm.validateSpecCore(spec)
}

// TestAppendRegion appends a prebuilt region to the manager for external tests.
func (rm *RegionManager) TestAppendRegion(r *Region) {
	rm.regions = append(rm.regions, r)
	rm.byName[strings.ToLower(r.name)] = r
}

// TestSetMode assigns the region's current mode for external tests.
func (r *Region) TestSetMode(mode string) {
	r.mode = mode
}

// TestSetInstance assigns the region's active instance for external tests.
func (r *Region) TestSetInstance(instance ModeInstance) {
	r.instance = instance
}

// TestRenderFrame exposes renderFrame to the external tests package.
func (rl *RenderLoop) TestRenderFrame() {
	rl.renderFrame()
}

// TestRenderDueRegions exposes renderDueRegions to the external tests package.
func (rl *RenderLoop) TestRenderDueRegions(now time.Time) {
	rl.renderDueRegions(now)
}

// TestFlush exposes flush to the external tests package.
func (rl *RenderLoop) TestFlush() {
	rl.flush()
}

// TestInitRegionTickers exposes initRegionTickers to the external tests package.
func (rl *RenderLoop) TestInitRegionTickers() {
	rl.initRegionTickers()
}

// TestMinSleepDuration exposes minSleepDuration to the external tests package.
func (rl *RenderLoop) TestMinSleepDuration() time.Duration {
	return rl.minSleepDuration()
}

// TestHasTickRateResolver reports whether a per-region resolver is configured.
func (rl *RenderLoop) TestHasTickRateResolver() bool {
	return rl.tickRateResolver != nil
}

// TestSetRegionTickers installs a complete ticker set for external tests.
func (rl *RenderLoop) TestSetRegionTickers(tickers []TestPerRegionTicker) {
	if tickers == nil {
		rl.regionTickers = nil
		return
	}
	rl.regionTickers = make([]*PerRegionTicker, len(tickers))
	for i, ticker := range tickers {
		rl.regionTickers[i] = &PerRegionTicker{
			region:   ticker.Region,
			interval: ticker.Interval,
			lastFire: ticker.LastFire,
		}
	}
}

// TestSetRegionTicker updates one ticker slot for external tests.
func (rl *RenderLoop) TestSetRegionTicker(index int, interval time.Duration, lastFire time.Time) {
	rl.regionTickers[index].interval = interval
	rl.regionTickers[index].lastFire = lastFire
}

// TestRegionTickers snapshots ticker state for external tests.
func (rl *RenderLoop) TestRegionTickers() []TestPerRegionTicker {
	if rl.regionTickers == nil {
		return nil
	}
	tickers := make([]TestPerRegionTicker, len(rl.regionTickers))
	for i, ticker := range rl.regionTickers {
		tickers[i] = TestPerRegionTicker{
			Region:   ticker.region,
			Interval: ticker.interval,
			LastFire: ticker.lastFire,
		}
	}
	return tickers
}

// TestHasTickRateProvider reports whether a provider is registered for modeID.
func TestHasTickRateProvider(modeID string) bool {
	providersMu.RLock()
	defer providersMu.RUnlock()
	_, ok := providers[modeID]
	return ok
}

// TestBuildCatalogHints exposes the region hint normalization/catalog path to tests.
func TestBuildCatalogHints(hints textlayout.TextHints, ppi float64) textlayout.TextHints {
	return buildCatalogIfNeeded(hints, ppi)
}
