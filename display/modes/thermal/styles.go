package thermal

import (
	"fmt"
	"regexp"

	"github.com/databeast/cyberhud/display/modes/thermal/source"
	"github.com/databeast/cyberhud/display/modes/thermal/styles"
	"github.com/databeast/cyberhud/display/style"
)

// styleNamePattern is the regex that all registered thermal style names must match.
// Format: <category_prefix><WxH>[-<variant>]
var styleNamePattern = regexp.MustCompile(`^(color-|color-slow-|mono-|mono-slow-|grayscale-fast-|grayscale-slow-)\d+x\d+(-[a-z0-9-]+)?$`)

func init() {
	for _, s := range thermalRegistry.Enumerate() {
		if !styleNamePattern.MatchString(s.Name()) {
			panic(fmt.Sprintf("thermal: invalid style name %q does not match expected pattern %s", s.Name(), styleNamePattern.String()))
		}
	}
}

// thermalRegistry is the StyleRegistry for the thermal display mode.
// Per-resolution styles are registered grouped by capability class.
// Registration order determines fitness tie-breaking (first registered wins).
// Skeleton styles (nil BuildFn, adaptiveBuild fallback) are registered FIRST
// within each capability group; polished variants follow.
var thermalRegistry = style.NewRegistry[source.ThermalSnapshot, source.Policy](
	// ──────────────────────────────────────────────────────────────────────
	// MonoSlow
	// ──────────────────────────────────────────────────────────────────────
	// 32x128 (small portrait, width<64 → minimal only)
	styles.MonoSlow32x128Style,
	styles.MonoSlow32x128MinimalStyle,
	// 64x128 (portrait → portrait variants + minimal)
	styles.MonoSlow64x128Style,
	styles.MonoSlow64x128ThermometerStyle,
	styles.MonoSlow64x128SparkStyle,
	styles.MonoSlow64x128HeatmapStyle,
	styles.MonoSlow64x128LEDsStyle,
	styles.MonoSlow64x128AvgThermoStyle,
	styles.MonoSlow64x128MinimalStyle,
	// 80x160 (portrait → portrait variants + minimal)
	styles.MonoSlow80x160Style,
	styles.MonoSlow80x160ThermometerStyle,
	styles.MonoSlow80x160SparkStyle,
	styles.MonoSlow80x160HeatmapStyle,
	styles.MonoSlow80x160LEDsStyle,
	styles.MonoSlow80x160AvgThermoStyle,
	styles.MonoSlow80x160MinimalStyle,
	// 128x32 (small landscape, height<64 → minimal only)
	styles.MonoSlow128x32Style,
	styles.MonoSlow128x32MinimalStyle,
	// 128x64 (landscape → landscape variants)
	styles.MonoSlow128x64Style, // polished (explicit buildMonoOLEDCompactStyle)
	styles.MonoSlow128x64OverviewStyle,
	styles.MonoSlow128x64DetailStyle,
	styles.MonoSlow128x64GraphStyle,
	styles.MonoSlow128x64MinimalStyle,
	// 128x128 (square → square variants)
	styles.MonoSlow128x128Style,
	styles.MonoSlow128x128OverviewStyle,
	styles.MonoSlow128x128DetailStyle,
	styles.MonoSlow128x128GraphStyle,
	styles.MonoSlow128x128MinimalStyle,
	// 128x160 (portrait → portrait variants + minimal)
	styles.MonoSlow128x160Style,
	styles.MonoSlow128x160ThermometerStyle,
	styles.MonoSlow128x160SparkStyle,
	styles.MonoSlow128x160HeatmapStyle,
	styles.MonoSlow128x160LEDsStyle,
	styles.MonoSlow128x160AvgThermoStyle,
	styles.MonoSlow128x160MinimalStyle,
	// 135x240 (portrait → portrait variants + minimal)
	styles.MonoSlow135x240Style,
	styles.MonoSlow135x240ThermometerStyle,
	styles.MonoSlow135x240SparkStyle,
	styles.MonoSlow135x240HeatmapStyle,
	styles.MonoSlow135x240LEDsStyle,
	styles.MonoSlow135x240AvgThermoStyle,
	styles.MonoSlow135x240MinimalStyle,
	// 160x80 (landscape → landscape variants)
	styles.MonoSlow160x80Style,
	styles.MonoSlow160x80OverviewStyle,
	styles.MonoSlow160x80DetailStyle,
	styles.MonoSlow160x80GraphStyle,
	styles.MonoSlow160x80MinimalStyle,
	// 160x128 (landscape → landscape variants)
	styles.MonoSlow160x128Style,
	styles.MonoSlow160x128OverviewStyle,
	styles.MonoSlow160x128DetailStyle,
	styles.MonoSlow160x128GraphStyle,
	styles.MonoSlow160x128MinimalStyle,
	// 240x135 (landscape → landscape variants)
	styles.MonoSlow240x135Style,
	styles.MonoSlow240x135OverviewStyle,
	styles.MonoSlow240x135DetailStyle,
	styles.MonoSlow240x135GraphStyle,
	styles.MonoSlow240x135MinimalStyle,
	// 240x240 (square → square variants)
	styles.MonoSlow240x240Style,
	styles.MonoSlow240x240OverviewStyle,
	styles.MonoSlow240x240DetailStyle,
	styles.MonoSlow240x240GraphStyle,
	styles.MonoSlow240x240MinimalStyle,
	// 240x320 (portrait → portrait variants + minimal)
	styles.MonoSlow240x320Style,
	styles.MonoSlow240x320ThermometerStyle,
	styles.MonoSlow240x320SparkStyle,
	styles.MonoSlow240x320HeatmapStyle,
	styles.MonoSlow240x320LEDsStyle,
	styles.MonoSlow240x320AvgThermoStyle,
	styles.MonoSlow240x320MinimalStyle,
	// 320x240 (landscape → landscape variants)
	styles.MonoSlow320x240Style,
	styles.MonoSlow320x240OverviewStyle,
	styles.MonoSlow320x240DetailStyle,
	styles.MonoSlow320x240GraphStyle,
	styles.MonoSlow320x240MinimalStyle,
	// 320x480 (portrait → portrait variants + minimal)
	styles.MonoSlow320x480Style,
	styles.MonoSlow320x480ThermometerStyle,
	styles.MonoSlow320x480SparkStyle,
	styles.MonoSlow320x480HeatmapStyle,
	styles.MonoSlow320x480LEDsStyle,
	styles.MonoSlow320x480AvgThermoStyle,
	styles.MonoSlow320x480MinimalStyle,
	// 480x320 (landscape → landscape variants)
	styles.MonoSlow480x320Style,
	styles.MonoSlow480x320OverviewStyle,
	styles.MonoSlow480x320DetailStyle,
	styles.MonoSlow480x320GraphStyle,
	styles.MonoSlow480x320MinimalStyle,
	// 800x480 (landscape → landscape variants)
	styles.MonoSlow800x480Style,
	styles.MonoSlow800x480OverviewStyle,
	styles.MonoSlow800x480DetailStyle,
	styles.MonoSlow800x480GraphStyle,
	styles.MonoSlow800x480MinimalStyle,

	// ──────────────────────────────────────────────────────────────────────
	// MonoFast
	// ──────────────────────────────────────────────────────────────────────
	// 32x128 (small portrait, width<64 → minimal only)
	styles.MonoFast32x128Style,
	styles.MonoFast32x128MinimalStyle,
	// 64x128 (portrait → portrait variants + minimal)
	styles.MonoFast64x128Style,
	styles.MonoFast64x128ThermometerStyle,
	styles.MonoFast64x128SparkStyle,
	styles.MonoFast64x128HeatmapStyle,
	styles.MonoFast64x128LEDsStyle,
	styles.MonoFast64x128AvgThermoStyle,
	styles.MonoFast64x128MinimalStyle,
	// 128x32 (small landscape, height<64 → minimal only)
	styles.MonoFast128x32Style,
	styles.MonoFast128x32MinimalStyle,
	// 128x64 (landscape → landscape variants + timegraph)
	styles.MonoFast128x64Style,
	styles.MonoFast128x64OverviewStyle,
	styles.MonoFast128x64TimegraphStyle,
	styles.MonoFast128x64DetailStyle,
	styles.MonoFast128x64GraphStyle,
	styles.MonoFast128x64MinimalStyle,
	// 128x128 (square → square variants + timegraph) — polished (explicit buildDetailStyle)
	styles.Mono128x128Style,
	styles.MonoFast128x128OverviewStyle,
	styles.MonoFast128x128TimegraphStyle,
	styles.MonoFast128x128GraphStyle,
	styles.MonoFast128x128MinimalStyle,

	// ──────────────────────────────────────────────────────────────────────
	// GrayscaleSlow
	// ──────────────────────────────────────────────────────────────────────
	// 32x128 (small portrait, width<64 → minimal only)
	styles.GrayscaleSlow32x128Style,
	styles.GrayscaleSlow32x128MinimalStyle,
	// 64x128 (portrait → portrait variants + minimal)
	styles.GrayscaleSlow64x128Style,
	styles.GrayscaleSlow64x128ThermometerStyle,
	styles.GrayscaleSlow64x128SparkStyle,
	styles.GrayscaleSlow64x128HeatmapStyle,
	styles.GrayscaleSlow64x128LEDsStyle,
	styles.GrayscaleSlow64x128AvgThermoStyle,
	styles.GrayscaleSlow64x128MinimalStyle,
	// 80x160 (portrait → portrait variants + minimal)
	styles.GrayscaleSlow80x160Style,
	styles.GrayscaleSlow80x160ThermometerStyle,
	styles.GrayscaleSlow80x160SparkStyle,
	styles.GrayscaleSlow80x160HeatmapStyle,
	styles.GrayscaleSlow80x160LEDsStyle,
	styles.GrayscaleSlow80x160AvgThermoStyle,
	styles.GrayscaleSlow80x160MinimalStyle,
	// 104x212 (portrait → portrait variants + minimal)
	styles.GrayscaleSlow104x212Style,
	styles.GrayscaleSlow104x212ThermometerStyle,
	styles.GrayscaleSlow104x212SparkStyle,
	styles.GrayscaleSlow104x212HeatmapStyle,
	styles.GrayscaleSlow104x212LEDsStyle,
	styles.GrayscaleSlow104x212AvgThermoStyle,
	styles.GrayscaleSlow104x212MinimalStyle,
	// 122x250 (portrait → portrait variants + minimal)
	styles.GrayscaleSlow122x250Style,
	styles.GrayscaleSlow122x250ThermometerStyle,
	styles.GrayscaleSlow122x250SparkStyle,
	styles.GrayscaleSlow122x250HeatmapStyle,
	styles.GrayscaleSlow122x250LEDsStyle,
	styles.GrayscaleSlow122x250AvgThermoStyle,
	styles.GrayscaleSlow122x250MinimalStyle,
	// 128x32 (small landscape, height<64 → minimal only)
	styles.GrayscaleSlow128x32Style,
	styles.GrayscaleSlow128x32MinimalStyle,
	// 128x64 (landscape → landscape variants)
	styles.GrayscaleSlow128x64Style,
	styles.GrayscaleSlow128x64OverviewStyle,
	styles.GrayscaleSlow128x64DetailStyle,
	styles.GrayscaleSlow128x64GraphStyle,
	styles.GrayscaleSlow128x64MinimalStyle,
	// 128x128 (square → square variants)
	styles.GrayscaleSlow128x128Style,
	styles.GrayscaleSlow128x128OverviewStyle,
	styles.GrayscaleSlow128x128DetailStyle,
	styles.GrayscaleSlow128x128GraphStyle,
	styles.GrayscaleSlow128x128MinimalStyle,
	// 128x160 (portrait → portrait variants + minimal)
	styles.GrayscaleSlow128x160Style,
	styles.GrayscaleSlow128x160ThermometerStyle,
	styles.GrayscaleSlow128x160SparkStyle,
	styles.GrayscaleSlow128x160HeatmapStyle,
	styles.GrayscaleSlow128x160LEDsStyle,
	styles.GrayscaleSlow128x160AvgThermoStyle,
	styles.GrayscaleSlow128x160MinimalStyle,
	// 128x296 (portrait → portrait variants + minimal)
	styles.GrayscaleSlow128x296Style,
	styles.GrayscaleSlow128x296ThermometerStyle,
	styles.GrayscaleSlow128x296SparkStyle,
	styles.GrayscaleSlow128x296HeatmapStyle,
	styles.GrayscaleSlow128x296LEDsStyle,
	styles.GrayscaleSlow128x296AvgThermoStyle,
	styles.GrayscaleSlow128x296MinimalStyle,
	// 135x240 (portrait → portrait variants + minimal)
	styles.GrayscaleSlow135x240Style,
	styles.GrayscaleSlow135x240ThermometerStyle,
	styles.GrayscaleSlow135x240SparkStyle,
	styles.GrayscaleSlow135x240HeatmapStyle,
	styles.GrayscaleSlow135x240LEDsStyle,
	styles.GrayscaleSlow135x240AvgThermoStyle,
	styles.GrayscaleSlow135x240MinimalStyle,
	// 160x80 (landscape → landscape variants)
	styles.GrayscaleSlow160x80Style,
	styles.GrayscaleSlow160x80OverviewStyle,
	styles.GrayscaleSlow160x80DetailStyle,
	styles.GrayscaleSlow160x80GraphStyle,
	styles.GrayscaleSlow160x80MinimalStyle,
	// 160x128 (landscape → landscape variants)
	styles.GrayscaleSlow160x128Style,
	styles.GrayscaleSlow160x128OverviewStyle,
	styles.GrayscaleSlow160x128DetailStyle,
	styles.GrayscaleSlow160x128GraphStyle,
	styles.GrayscaleSlow160x128MinimalStyle,
	// 176x264 (portrait → portrait variants + minimal)
	styles.GrayscaleSlow176x264Style,
	styles.GrayscaleSlow176x264ThermometerStyle,
	styles.GrayscaleSlow176x264SparkStyle,
	styles.GrayscaleSlow176x264HeatmapStyle,
	styles.GrayscaleSlow176x264LEDsStyle,
	styles.GrayscaleSlow176x264AvgThermoStyle,
	styles.GrayscaleSlow176x264MinimalStyle,
	// 200x200 (square → square variants)
	styles.GrayscaleSlow200x200Style,
	styles.GrayscaleSlow200x200OverviewStyle,
	styles.GrayscaleSlow200x200DetailStyle,
	styles.GrayscaleSlow200x200GraphStyle,
	styles.GrayscaleSlow200x200MinimalStyle,
	// 212x104 (landscape → landscape variants)
	styles.GrayscaleSlow212x104Style,
	styles.GrayscaleSlow212x104OverviewStyle,
	styles.GrayscaleSlow212x104DetailStyle,
	styles.GrayscaleSlow212x104GraphStyle,
	styles.GrayscaleSlow212x104MinimalStyle,
	// 240x135 (landscape → landscape variants)
	styles.GrayscaleSlow240x135Style,
	styles.GrayscaleSlow240x135OverviewStyle,
	styles.GrayscaleSlow240x135DetailStyle,
	styles.GrayscaleSlow240x135GraphStyle,
	styles.GrayscaleSlow240x135MinimalStyle,
	// 240x240 (square → square variants)
	styles.GrayscaleSlow240x240Style,
	styles.GrayscaleSlow240x240OverviewStyle,
	styles.GrayscaleSlow240x240DetailStyle,
	styles.GrayscaleSlow240x240GraphStyle,
	styles.GrayscaleSlow240x240MinimalStyle,
	// 240x320 (portrait → portrait variants + minimal)
	styles.GrayscaleSlow240x320Style,
	styles.GrayscaleSlow240x320ThermometerStyle,
	styles.GrayscaleSlow240x320SparkStyle,
	styles.GrayscaleSlow240x320HeatmapStyle,
	styles.GrayscaleSlow240x320LEDsStyle,
	styles.GrayscaleSlow240x320AvgThermoStyle,
	styles.GrayscaleSlow240x320MinimalStyle,
	// 250x122 (landscape → landscape variants)
	styles.GrayscaleSlow250x122Style,
	styles.GrayscaleSlow250x122OverviewStyle,
	styles.GrayscaleSlow250x122DetailStyle,
	styles.GrayscaleSlow250x122GraphStyle,
	styles.GrayscaleSlow250x122MinimalStyle,
	// 264x176 (landscape → landscape variants)
	styles.GrayscaleSlow264x176Style,
	styles.GrayscaleSlow264x176OverviewStyle,
	styles.GrayscaleSlow264x176DetailStyle,
	styles.GrayscaleSlow264x176GraphStyle,
	styles.GrayscaleSlow264x176MinimalStyle,
	// 296x128 (landscape → landscape variants) — polished (explicit buildOverviewStyle)
	styles.GrayscaleSlow296x128Style,
	styles.GrayscaleSlow296x128DetailStyle,
	styles.GrayscaleSlow296x128GraphStyle,
	styles.GrayscaleSlow296x128MinimalStyle,
	// 300x400 (portrait → portrait variants + minimal)
	styles.GrayscaleSlow300x400Style,
	styles.GrayscaleSlow300x400ThermometerStyle,
	styles.GrayscaleSlow300x400SparkStyle,
	styles.GrayscaleSlow300x400HeatmapStyle,
	styles.GrayscaleSlow300x400LEDsStyle,
	styles.GrayscaleSlow300x400AvgThermoStyle,
	styles.GrayscaleSlow300x400MinimalStyle,
	// 320x240 (landscape → landscape variants)
	styles.GrayscaleSlow320x240Style,
	styles.GrayscaleSlow320x240OverviewStyle,
	styles.GrayscaleSlow320x240DetailStyle,
	styles.GrayscaleSlow320x240GraphStyle,
	styles.GrayscaleSlow320x240MinimalStyle,
	// 320x480 (portrait → portrait variants + minimal)
	styles.GrayscaleSlow320x480Style,
	styles.GrayscaleSlow320x480ThermometerStyle,
	styles.GrayscaleSlow320x480SparkStyle,
	styles.GrayscaleSlow320x480HeatmapStyle,
	styles.GrayscaleSlow320x480LEDsStyle,
	styles.GrayscaleSlow320x480AvgThermoStyle,
	styles.GrayscaleSlow320x480MinimalStyle,
	// 400x300 (landscape → landscape variants)
	styles.GrayscaleSlow400x300Style,
	styles.GrayscaleSlow400x300OverviewStyle,
	styles.GrayscaleSlow400x300DetailStyle,
	styles.GrayscaleSlow400x300GraphStyle,
	styles.GrayscaleSlow400x300MinimalStyle,
	// 480x320 (landscape → landscape variants)
	styles.GrayscaleSlow480x320Style,
	styles.GrayscaleSlow480x320OverviewStyle,
	styles.GrayscaleSlow480x320DetailStyle,
	styles.GrayscaleSlow480x320GraphStyle,
	styles.GrayscaleSlow480x320MinimalStyle,
	// 480x800 (portrait → portrait variants + minimal)
	styles.GrayscaleSlow480x800Style,
	styles.GrayscaleSlow480x800ThermometerStyle,
	styles.GrayscaleSlow480x800SparkStyle,
	styles.GrayscaleSlow480x800HeatmapStyle,
	styles.GrayscaleSlow480x800LEDsStyle,
	styles.GrayscaleSlow480x800AvgThermoStyle,
	styles.GrayscaleSlow480x800MinimalStyle,
	// 800x480 (landscape → landscape variants)
	styles.GrayscaleSlow800x480Style,
	styles.GrayscaleSlow800x480OverviewStyle,
	styles.GrayscaleSlow800x480DetailStyle,
	styles.GrayscaleSlow800x480GraphStyle,
	styles.GrayscaleSlow800x480MinimalStyle,

	// ──────────────────────────────────────────────────────────────────────
	// GrayscaleFast
	// ──────────────────────────────────────────────────────────────────────
	// 80x160 (portrait → portrait variants + minimal)
	styles.GrayscaleFast80x160Style,
	styles.GrayscaleFast80x160ThermometerStyle,
	styles.GrayscaleFast80x160SparkStyle,
	styles.GrayscaleFast80x160HeatmapStyle,
	styles.GrayscaleFast80x160LEDsStyle,
	styles.GrayscaleFast80x160AvgThermoStyle,
	styles.GrayscaleFast80x160MinimalStyle,
	// 128x128 (square → square variants + timegraph)
	styles.GrayscaleFast128x128Style,
	styles.GrayscaleFast128x128OverviewStyle,
	styles.GrayscaleFast128x128TimegraphStyle,
	styles.GrayscaleFast128x128DetailStyle,
	styles.GrayscaleFast128x128GraphStyle,
	styles.GrayscaleFast128x128MinimalStyle,
	// 128x160 (portrait → portrait variants + minimal)
	styles.GrayscaleFast128x160Style,
	styles.GrayscaleFast128x160ThermometerStyle,
	styles.GrayscaleFast128x160SparkStyle,
	styles.GrayscaleFast128x160HeatmapStyle,
	styles.GrayscaleFast128x160LEDsStyle,
	styles.GrayscaleFast128x160AvgThermoStyle,
	styles.GrayscaleFast128x160MinimalStyle,
	// 135x240 (portrait → portrait variants + minimal)
	styles.GrayscaleFast135x240Style,
	styles.GrayscaleFast135x240ThermometerStyle,
	styles.GrayscaleFast135x240SparkStyle,
	styles.GrayscaleFast135x240HeatmapStyle,
	styles.GrayscaleFast135x240LEDsStyle,
	styles.GrayscaleFast135x240AvgThermoStyle,
	styles.GrayscaleFast135x240MinimalStyle,
	// 160x80 (landscape → landscape variants + timegraph)
	styles.GrayscaleFast160x80Style,
	styles.GrayscaleFast160x80OverviewStyle,
	styles.GrayscaleFast160x80TimegraphStyle,
	styles.GrayscaleFast160x80DetailStyle,
	styles.GrayscaleFast160x80GraphStyle,
	styles.GrayscaleFast160x80MinimalStyle,
	// 160x128 (landscape → landscape variants + timegraph)
	styles.GrayscaleFast160x128Style,
	styles.GrayscaleFast160x128OverviewStyle,
	styles.GrayscaleFast160x128TimegraphStyle,
	styles.GrayscaleFast160x128DetailStyle,
	styles.GrayscaleFast160x128GraphStyle,
	styles.GrayscaleFast160x128MinimalStyle,
	// 240x135 (landscape → landscape variants + timegraph)
	styles.GrayscaleFast240x135Style,
	styles.GrayscaleFast240x135OverviewStyle,
	styles.GrayscaleFast240x135TimegraphStyle,
	styles.GrayscaleFast240x135DetailStyle,
	styles.GrayscaleFast240x135GraphStyle,
	styles.GrayscaleFast240x135MinimalStyle,
	// 240x240 (square → square variants + timegraph)
	styles.GrayscaleFast240x240Style,
	styles.GrayscaleFast240x240OverviewStyle,
	styles.GrayscaleFast240x240TimegraphStyle,
	styles.GrayscaleFast240x240DetailStyle,
	styles.GrayscaleFast240x240GraphStyle,
	styles.GrayscaleFast240x240MinimalStyle,
	// 240x320 (portrait → portrait variants + minimal)
	styles.GrayscaleFast240x320Style,
	styles.GrayscaleFast240x320ThermometerStyle,
	styles.GrayscaleFast240x320SparkStyle,
	styles.GrayscaleFast240x320HeatmapStyle,
	styles.GrayscaleFast240x320LEDsStyle,
	styles.GrayscaleFast240x320AvgThermoStyle,
	styles.GrayscaleFast240x320MinimalStyle,
	// 320x240 (landscape → landscape variants + timegraph)
	styles.GrayscaleFast320x240Style,
	styles.GrayscaleFast320x240OverviewStyle,
	styles.GrayscaleFast320x240TimegraphStyle,
	styles.GrayscaleFast320x240DetailStyle,
	styles.GrayscaleFast320x240GraphStyle,
	styles.GrayscaleFast320x240MinimalStyle,
	// 320x480 (portrait → portrait variants + minimal)
	styles.GrayscaleFast320x480Style,
	styles.GrayscaleFast320x480ThermometerStyle,
	styles.GrayscaleFast320x480SparkStyle,
	styles.GrayscaleFast320x480HeatmapStyle,
	styles.GrayscaleFast320x480LEDsStyle,
	styles.GrayscaleFast320x480AvgThermoStyle,
	styles.GrayscaleFast320x480MinimalStyle,
	// 400x300 (landscape → landscape variants + timegraph) — polished (explicit buildOverviewStyle)
	styles.GrayscaleFast400x300Style,
	styles.GrayscaleFast400x300TimegraphStyle,
	styles.GrayscaleFast400x300DetailStyle,
	styles.GrayscaleFast400x300GraphStyle,
	styles.GrayscaleFast400x300MinimalStyle,
	// 480x320 (landscape → landscape variants + timegraph)
	styles.GrayscaleFast480x320Style,
	styles.GrayscaleFast480x320OverviewStyle,
	styles.GrayscaleFast480x320TimegraphStyle,
	styles.GrayscaleFast480x320DetailStyle,
	styles.GrayscaleFast480x320GraphStyle,
	styles.GrayscaleFast480x320MinimalStyle,
	// 480x800 (portrait → portrait variants + minimal)
	styles.GrayscaleFast480x800Style,
	styles.GrayscaleFast480x800ThermometerStyle,
	styles.GrayscaleFast480x800SparkStyle,
	styles.GrayscaleFast480x800HeatmapStyle,
	styles.GrayscaleFast480x800LEDsStyle,
	styles.GrayscaleFast480x800AvgThermoStyle,
	styles.GrayscaleFast480x800MinimalStyle,
	// 800x480 (landscape → landscape variants + timegraph)
	styles.GrayscaleFast800x480Style,
	styles.GrayscaleFast800x480OverviewStyle,
	styles.GrayscaleFast800x480TimegraphStyle,
	styles.GrayscaleFast800x480DetailStyle,
	styles.GrayscaleFast800x480GraphStyle,
	styles.GrayscaleFast800x480MinimalStyle,

	// ──────────────────────────────────────────────────────────────────────
	// ColorSlow
	// ──────────────────────────────────────────────────────────────────────
	// 32x128 (small portrait, width<64 → minimal only)
	styles.ColorSlow32x128Style,
	styles.ColorSlow32x128MinimalStyle,
	// 64x128 (portrait → portrait variants + minimal)
	styles.ColorSlow64x128Style,
	styles.ColorSlow64x128ThermometerStyle,
	styles.ColorSlow64x128SparkStyle,
	styles.ColorSlow64x128HeatmapStyle,
	styles.ColorSlow64x128LEDsStyle,
	styles.ColorSlow64x128AvgThermoStyle,
	styles.ColorSlow64x128MinimalStyle,
	// 80x160 (portrait → portrait variants + minimal)
	styles.ColorSlow80x160Style,
	styles.ColorSlow80x160ThermometerStyle,
	styles.ColorSlow80x160SparkStyle,
	styles.ColorSlow80x160HeatmapStyle,
	styles.ColorSlow80x160LEDsStyle,
	styles.ColorSlow80x160AvgThermoStyle,
	styles.ColorSlow80x160MinimalStyle,
	// 104x212 (portrait → portrait variants + minimal)
	styles.ColorSlow104x212Style,
	styles.ColorSlow104x212ThermometerStyle,
	styles.ColorSlow104x212SparkStyle,
	styles.ColorSlow104x212HeatmapStyle,
	styles.ColorSlow104x212LEDsStyle,
	styles.ColorSlow104x212AvgThermoStyle,
	styles.ColorSlow104x212MinimalStyle,
	// 122x250 (portrait → portrait variants + minimal)
	styles.ColorSlow122x250Style,
	styles.ColorSlow122x250ThermometerStyle,
	styles.ColorSlow122x250SparkStyle,
	styles.ColorSlow122x250HeatmapStyle,
	styles.ColorSlow122x250LEDsStyle,
	styles.ColorSlow122x250AvgThermoStyle,
	styles.ColorSlow122x250MinimalStyle,
	// 128x32 (small landscape, height<64 → minimal only)
	styles.ColorSlow128x32Style,
	styles.ColorSlow128x32MinimalStyle,
	// 128x64 (landscape → landscape variants)
	styles.ColorSlow128x64Style,
	styles.ColorSlow128x64OverviewStyle,
	styles.ColorSlow128x64DetailStyle,
	styles.ColorSlow128x64GraphStyle,
	styles.ColorSlow128x64MinimalStyle,
	// 128x128 (square → square variants)
	styles.ColorSlow128x128Style,
	styles.ColorSlow128x128OverviewStyle,
	styles.ColorSlow128x128DetailStyle,
	styles.ColorSlow128x128GraphStyle,
	styles.ColorSlow128x128MinimalStyle,
	// 128x160 (portrait → portrait variants + minimal)
	styles.ColorSlow128x160Style,
	styles.ColorSlow128x160ThermometerStyle,
	styles.ColorSlow128x160SparkStyle,
	styles.ColorSlow128x160HeatmapStyle,
	styles.ColorSlow128x160LEDsStyle,
	styles.ColorSlow128x160AvgThermoStyle,
	styles.ColorSlow128x160MinimalStyle,
	// 128x296 (portrait → portrait variants + minimal)
	styles.ColorSlow128x296Style,
	styles.ColorSlow128x296ThermometerStyle,
	styles.ColorSlow128x296SparkStyle,
	styles.ColorSlow128x296HeatmapStyle,
	styles.ColorSlow128x296LEDsStyle,
	styles.ColorSlow128x296AvgThermoStyle,
	styles.ColorSlow128x296MinimalStyle,
	// 135x240 (portrait → portrait variants + minimal)
	styles.ColorSlow135x240Style,
	styles.ColorSlow135x240ThermometerStyle,
	styles.ColorSlow135x240SparkStyle,
	styles.ColorSlow135x240HeatmapStyle,
	styles.ColorSlow135x240LEDsStyle,
	styles.ColorSlow135x240AvgThermoStyle,
	styles.ColorSlow135x240MinimalStyle,
	// 160x80 (landscape → landscape variants)
	styles.ColorSlow160x80Style,
	styles.ColorSlow160x80OverviewStyle,
	styles.ColorSlow160x80DetailStyle,
	styles.ColorSlow160x80GraphStyle,
	styles.ColorSlow160x80MinimalStyle,
	// 160x128 (landscape → landscape variants)
	styles.ColorSlow160x128Style,
	styles.ColorSlow160x128OverviewStyle,
	styles.ColorSlow160x128DetailStyle,
	styles.ColorSlow160x128GraphStyle,
	styles.ColorSlow160x128MinimalStyle,
	// 176x264 (portrait → portrait variants + minimal)
	styles.ColorSlow176x264Style,
	styles.ColorSlow176x264ThermometerStyle,
	styles.ColorSlow176x264SparkStyle,
	styles.ColorSlow176x264HeatmapStyle,
	styles.ColorSlow176x264LEDsStyle,
	styles.ColorSlow176x264AvgThermoStyle,
	styles.ColorSlow176x264MinimalStyle,
	// 200x200 (square → square variants)
	styles.ColorSlow200x200Style,
	styles.ColorSlow200x200OverviewStyle,
	styles.ColorSlow200x200DetailStyle,
	styles.ColorSlow200x200GraphStyle,
	styles.ColorSlow200x200MinimalStyle,
	// 212x104 (landscape → landscape variants)
	styles.ColorSlow212x104Style,
	styles.ColorSlow212x104OverviewStyle,
	styles.ColorSlow212x104DetailStyle,
	styles.ColorSlow212x104GraphStyle,
	styles.ColorSlow212x104MinimalStyle,
	// 240x135 (landscape → landscape variants)
	styles.ColorSlow240x135Style,
	styles.ColorSlow240x135OverviewStyle,
	styles.ColorSlow240x135DetailStyle,
	styles.ColorSlow240x135GraphStyle,
	styles.ColorSlow240x135MinimalStyle,
	// 240x240 (square → square variants)
	styles.ColorSlow240x240Style,
	styles.ColorSlow240x240OverviewStyle,
	styles.ColorSlow240x240DetailStyle,
	styles.ColorSlow240x240GraphStyle,
	styles.ColorSlow240x240MinimalStyle,
	// 240x320 (portrait → portrait variants + minimal)
	styles.ColorSlow240x320Style,
	styles.ColorSlow240x320ThermometerStyle,
	styles.ColorSlow240x320SparkStyle,
	styles.ColorSlow240x320HeatmapStyle,
	styles.ColorSlow240x320LEDsStyle,
	styles.ColorSlow240x320AvgThermoStyle,
	styles.ColorSlow240x320MinimalStyle,
	// 250x122 (landscape → landscape variants)
	styles.ColorSlow250x122Style,
	styles.ColorSlow250x122OverviewStyle,
	styles.ColorSlow250x122DetailStyle,
	styles.ColorSlow250x122GraphStyle,
	styles.ColorSlow250x122MinimalStyle,
	// 264x176 (landscape → landscape variants)
	styles.ColorSlow264x176Style,
	styles.ColorSlow264x176OverviewStyle,
	styles.ColorSlow264x176DetailStyle,
	styles.ColorSlow264x176GraphStyle,
	styles.ColorSlow264x176MinimalStyle,
	// 296x128 (landscape → landscape variants)
	styles.ColorSlow296x128Style,
	styles.ColorSlow296x128OverviewStyle,
	styles.ColorSlow296x128DetailStyle,
	styles.ColorSlow296x128GraphStyle,
	styles.ColorSlow296x128MinimalStyle,
	// 300x400 (portrait → portrait variants + minimal)
	styles.ColorSlow300x400Style,
	styles.ColorSlow300x400ThermometerStyle,
	styles.ColorSlow300x400SparkStyle,
	styles.ColorSlow300x400HeatmapStyle,
	styles.ColorSlow300x400LEDsStyle,
	styles.ColorSlow300x400AvgThermoStyle,
	styles.ColorSlow300x400MinimalStyle,
	// 320x240 (landscape → landscape variants)
	styles.ColorSlow320x240Style,
	styles.ColorSlow320x240OverviewStyle,
	styles.ColorSlow320x240DetailStyle,
	styles.ColorSlow320x240GraphStyle,
	styles.ColorSlow320x240MinimalStyle,
	// 320x480 (portrait → portrait variants + minimal)
	styles.ColorSlow320x480Style,
	styles.ColorSlow320x480ThermometerStyle,
	styles.ColorSlow320x480SparkStyle,
	styles.ColorSlow320x480HeatmapStyle,
	styles.ColorSlow320x480LEDsStyle,
	styles.ColorSlow320x480AvgThermoStyle,
	styles.ColorSlow320x480MinimalStyle,
	// 400x300 (landscape → landscape variants)
	styles.ColorSlow400x300Style,
	styles.ColorSlow400x300OverviewStyle,
	styles.ColorSlow400x300DetailStyle,
	styles.ColorSlow400x300GraphStyle,
	styles.ColorSlow400x300MinimalStyle,
	// 480x320 (landscape → landscape variants)
	styles.ColorSlow480x320Style,
	styles.ColorSlow480x320OverviewStyle,
	styles.ColorSlow480x320DetailStyle,
	styles.ColorSlow480x320GraphStyle,
	styles.ColorSlow480x320MinimalStyle,
	// 480x800 (portrait → portrait variants + minimal)
	styles.ColorSlow480x800Style,
	styles.ColorSlow480x800ThermometerStyle,
	styles.ColorSlow480x800SparkStyle,
	styles.ColorSlow480x800HeatmapStyle,
	styles.ColorSlow480x800LEDsStyle,
	styles.ColorSlow480x800AvgThermoStyle,
	styles.ColorSlow480x800MinimalStyle,
	// 800x480 (landscape → landscape variants)
	styles.ColorSlow800x480Style,
	styles.ColorSlow800x480OverviewStyle,
	styles.ColorSlow800x480DetailStyle,
	styles.ColorSlow800x480GraphStyle,
	styles.ColorSlow800x480MinimalStyle,

	// ──────────────────────────────────────────────────────────────────────
	// ColorFast
	// ──────────────────────────────────────────────────────────────────────
	// 80x160 (portrait → portrait variants + minimal)
	styles.ColorFast80x160Style,
	styles.ColorFast80x160ThermometerStyle,
	styles.ColorFast80x160SparkStyle,
	styles.ColorFast80x160HeatmapStyle,
	styles.ColorFast80x160LEDsStyle,
	styles.ColorFast80x160AvgThermoStyle,
	styles.ColorFast80x160MinimalStyle,
	// 128x128 (square → square variants + timegraph)
	styles.ColorFast128x128Style,
	styles.ColorFast128x128OverviewStyle,
	styles.ColorFast128x128TimegraphStyle,
	styles.ColorFast128x128DetailStyle,
	styles.ColorFast128x128GraphStyle,
	styles.ColorFast128x128MinimalStyle,
	// 128x160 (portrait → portrait variants + minimal)
	styles.ColorFast128x160Style,
	styles.ColorFast128x160ThermometerStyle,
	styles.ColorFast128x160SparkStyle,
	styles.ColorFast128x160HeatmapStyle,
	styles.ColorFast128x160LEDsStyle,
	styles.ColorFast128x160AvgThermoStyle,
	styles.ColorFast128x160MinimalStyle,
	// 135x240 (portrait → portrait variants + minimal)
	styles.ColorFast135x240Style,
	styles.ColorFast135x240ThermometerStyle,
	styles.ColorFast135x240SparkStyle,
	styles.ColorFast135x240HeatmapStyle,
	styles.ColorFast135x240LEDsStyle,
	styles.ColorFast135x240AvgThermoStyle,
	styles.ColorFast135x240MinimalStyle,
	// 160x80 (landscape → landscape variants + timegraph)
	styles.ColorFast160x80Style,
	styles.ColorFast160x80OverviewStyle,
	styles.ColorFast160x80TimegraphStyle,
	styles.ColorFast160x80DetailStyle,
	styles.ColorFast160x80GraphStyle,
	styles.ColorFast160x80MinimalStyle,
	// 160x128 (landscape → landscape variants + timegraph)
	styles.ColorFast160x128Style,
	styles.ColorFast160x128OverviewStyle,
	styles.ColorFast160x128TimegraphStyle,
	styles.ColorFast160x128DetailStyle,
	styles.ColorFast160x128GraphStyle,
	styles.ColorFast160x128MinimalStyle,
	// 240x135 (landscape → landscape variants + timegraph)
	styles.ColorFast240x135Style,
	styles.ColorFast240x135OverviewStyle,
	styles.ColorFast240x135TimegraphStyle,
	styles.ColorFast240x135DetailStyle,
	styles.ColorFast240x135GraphStyle,
	styles.ColorFast240x135MinimalStyle,
	// 240x240 (square → square variants + timegraph)
	styles.ColorFast240x240Style,
	styles.ColorFast240x240OverviewStyle,
	styles.ColorFast240x240TimegraphStyle,
	styles.ColorFast240x240DetailStyle,
	styles.ColorFast240x240GraphStyle,
	styles.ColorFast240x240MinimalStyle,
	// 240x320 (portrait → portrait variants + minimal) — polished variants in styles.go
	styles.ColorFast240x320Style,
	styles.Color240x320ThermometerStyle,
	styles.Color240x320SparkStyle,
	styles.Color240x320HeatmapStyle,
	styles.Color240x320LEDsStyle,
	styles.Color240x320AvgThermoStyle,
	styles.ColorFast240x320MinimalStyle,
	// 320x240 (landscape → landscape variants + timegraph) — polished variants in styles.go
	styles.ColorFast320x240Style,
	styles.Color320x240OverviewStyle,
	styles.Color320x240TimegraphStyle,
	styles.ColorFast320x240DetailStyle,
	styles.ColorFast320x240GraphStyle,
	styles.ColorFast320x240MinimalStyle,
	// 320x480 (portrait → portrait variants + minimal)
	styles.ColorFast320x480Style,
	styles.ColorFast320x480ThermometerStyle,
	styles.ColorFast320x480SparkStyle,
	styles.ColorFast320x480HeatmapStyle,
	styles.ColorFast320x480LEDsStyle,
	styles.ColorFast320x480AvgThermoStyle,
	styles.ColorFast320x480MinimalStyle,
	// 480x320 (landscape → landscape variants + timegraph)
	styles.ColorFast480x320Style,
	styles.ColorFast480x320OverviewStyle,
	styles.ColorFast480x320TimegraphStyle,
	styles.ColorFast480x320DetailStyle,
	styles.ColorFast480x320GraphStyle,
	styles.ColorFast480x320MinimalStyle,
	// 480x800 (portrait → portrait variants + minimal)
	styles.ColorFast480x800Style,
	styles.ColorFast480x800ThermometerStyle,
	styles.ColorFast480x800SparkStyle,
	styles.ColorFast480x800HeatmapStyle,
	styles.ColorFast480x800LEDsStyle,
	styles.ColorFast480x800AvgThermoStyle,
	styles.ColorFast480x800MinimalStyle,
	// 800x480 (landscape → landscape variants + timegraph)
	styles.ColorFast800x480Style,
	styles.ColorFast800x480OverviewStyle,
	styles.ColorFast800x480TimegraphStyle,
	styles.ColorFast800x480DetailStyle,
	styles.ColorFast800x480GraphStyle,
	styles.ColorFast800x480MinimalStyle,
)

// allowedStyleNames returns the list of registered style names for use in validators.
func allowedStyleNames() []string {
	styles := thermalRegistry.Enumerate()
	names := make([]string, len(styles))
	for i, s := range styles {
		names[i] = s.Name()
	}
	return names
}

// allEmptyItems returns true if items is nil, empty, or contains only empty strings.
func allEmptyItems(items []string) bool {
	if len(items) == 0 {
		return true
	}
	for _, item := range items {
		if item != "" {
			return false
		}
	}
	return true
}
