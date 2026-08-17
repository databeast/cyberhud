package source

import (
	"image/color"
	"time"

	fonts "github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/widgets/marquee"
)

func CharacterPool() []rune {
	out := make([]rune, len(characterPool))
	copy(out, characterPool)
	return out
}

func DefaultMutationInterval() time.Duration { return defaultMutationInterval }

func ComputeColumnCount(panelWidth, glyphAdvance int) int {
	return computeColumnCount(panelWidth, glyphAdvance)
}

func ComputeVisibleCells(panelHeight, rowHeight int) int {
	return computeVisibleCells(panelHeight, rowHeight)
}

func ComputeActiveColumns(totalColumns int, density float64) []int {
	return computeActiveColumns(totalColumns, density)
}

func BuildColorArray(trailLength int, mono bool) []color.RGBA {
	return buildColorArray(trailLength, mono)
}

func ResetLayoutCache() {
	stripCache = nil
	stripColumnMap = nil
	lastPolicyFingerprint = ""
	lastCacheKey = ""
	firstBuildDone = false
}

func ResetFontCache() {
	cachedGlyphAdvance = 0
	cachedRowHeight = 0
	cachedFontResolved = false
	cachedSpriteFace = nil
	cachedCatalogWidth = 0
	cachedCatalogHeight = 0
}

func SetCycleStart(t time.Time) { cycleStart = t }

func ResetCycleStart() { cycleStart = time.Time{} }

func RebuildStripsForTest(p Policy, panelWidth, panelHeight, glyphAdvance, rowHeight int, mono bool, seed int64, face fonts.Face) []*marquee.Strip {
	return rebuildStrips(p, panelWidth, panelHeight, glyphAdvance, rowHeight, mono, seed, face)
}
