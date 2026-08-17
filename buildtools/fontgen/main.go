package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/databeast/cyberhud/buildtools/fontgen/bdf"
	"github.com/databeast/cyberhud/buildtools/fontgen/clone"
	"github.com/databeast/cyberhud/buildtools/fontgen/codegen"
	"github.com/databeast/cyberhud/buildtools/fontgen/ttf"
)

// Known upstream repository URLs.
var repoURLs = map[string]string{
	"spleen":   "https://github.com/fcambus/spleen.git",
	"cozette":  "https://github.com/slavfox/Cozette.git",
	"terminus": "https://github.com/mikebeaton/terminus-font-4.49.1.git",
}

func main() {
	var (
		bdfPath      string
		ttfPath      string
		outPath      string
		fontID       string
		pkg          string
		width        int
		height       int
		advance      int
		rowheight    int
		cloneName    string
		downloadURL  string
		targetHeight int
		ranges       string
	)

	flag.StringVar(&bdfPath, "bdf", "", "path to BDF file (relative to cloned repo root, or filename for -download)")
	flag.StringVar(&ttfPath, "ttf", "", "path to TTF file (mutually exclusive with -clone/-download)")
	flag.StringVar(&outPath, "out", "", "output Go source file path")
	flag.StringVar(&fontID, "id", "", "font ID string (e.g. spleen-5x8)")
	flag.StringVar(&pkg, "pkg", "", "Go package name for generated file")
	flag.IntVar(&width, "width", 0, "glyph width in pixels")
	flag.IntVar(&height, "height", 0, "glyph height in pixels")
	flag.IntVar(&advance, "advance", 0, "glyph advance in pixels")
	flag.IntVar(&rowheight, "rowheight", 0, "row height in pixels")
	flag.StringVar(&cloneName, "clone", "", "upstream repo to clone (spleen, terminus)")
	flag.StringVar(&downloadURL, "download", "", "URL to download BDF file directly (alternative to -clone)")
	flag.IntVar(&targetHeight, "targetheight", 0, "target glyph height in pixels (ppem) for TTF rasterization")
	flag.StringVar(&ranges, "ranges", "", "comma-separated low-high codepoint ranges (e.g. \"33-126,65382-65437\")")

	flag.Parse()

	if err := run(bdfPath, ttfPath, outPath, fontID, pkg, width, height, advance, rowheight, cloneName, downloadURL, targetHeight, ranges); err != nil {
		fmt.Fprintf(os.Stderr, "fontgen: %v\n", err)
		os.Exit(1)
	}
}

func run(bdfPath, ttfPath, outPath, fontID, pkg string, width, height, advance, rowheight int, cloneName, downloadURL string, targetHeight int, ranges string) error {
	// Validate required flags.
	if outPath == "" {
		return fmt.Errorf("missing required flag: -out")
	}
	if fontID == "" {
		return fmt.Errorf("missing required flag: -id")
	}
	if pkg == "" {
		return fmt.Errorf("missing required flag: -pkg")
	}
	if width <= 0 {
		return fmt.Errorf("missing or invalid flag: -width")
	}
	if advance <= 0 {
		return fmt.Errorf("missing or invalid flag: -advance")
	}
	if rowheight <= 0 {
		return fmt.Errorf("missing or invalid flag: -rowheight")
	}

	// Determine input mode: TTF or BDF.
	useTTF := ttfPath != ""
	useBDF := cloneName != "" || downloadURL != ""

	if useTTF && useBDF {
		return fmt.Errorf("-ttf is mutually exclusive with -clone/-download")
	}
	if !useTTF && !useBDF {
		return fmt.Errorf("must specify either -ttf or -clone/-download")
	}

	// Parse -ranges flag if provided.
	var parsedRanges []codepointRange
	if ranges != "" {
		var err error
		parsedRanges, err = parseRangesFlag(ranges)
		if err != nil {
			return fmt.Errorf("invalid -ranges: %w", err)
		}
	}

	// Derive struct/const/array names from font ID.
	structName := idToIdentifier(fontID) + "Face"
	constName := idToExportedIdentifier(fontID) + "ID"
	arrayName := idToIdentifier(fontID) + "Rows"

	if useTTF {
		return runTTF(ttfPath, outPath, fontID, pkg, width, height, advance, rowheight, targetHeight, parsedRanges, structName, constName, arrayName)
	}
	return runBDF(bdfPath, outPath, fontID, pkg, width, height, advance, rowheight, cloneName, downloadURL, parsedRanges, structName, constName, arrayName)
}

// codepointRange is the internal representation of a parsed range pair.
type codepointRange struct {
	low  rune
	high rune
}

// parseRangesFlag parses a comma-separated string of "low-high" decimal pairs.
func parseRangesFlag(s string) ([]codepointRange, error) {
	parts := strings.Split(s, ",")
	var result []codepointRange
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		dash := strings.Index(part, "-")
		if dash < 0 {
			return nil, fmt.Errorf("expected low-high pair, got %q", part)
		}
		lowStr := part[:dash]
		highStr := part[dash+1:]
		low, err := strconv.ParseInt(lowStr, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid low codepoint %q: %w", lowStr, err)
		}
		high, err := strconv.ParseInt(highStr, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid high codepoint %q: %w", highStr, err)
		}
		if low > high {
			return nil, fmt.Errorf("invalid range: low %d > high %d", low, high)
		}
		result = append(result, codepointRange{low: rune(low), high: rune(high)})
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no valid ranges found")
	}
	return result, nil
}

// runTTF handles the TTF input path.
func runTTF(ttfPath, outPath, fontID, pkg string, width, height, advance, rowheight, targetHeight int, ranges []codepointRange, structName, constName, arrayName string) error {
	if targetHeight <= 0 {
		return fmt.Errorf("-targetheight is required when using -ttf")
	}
	if len(ranges) == 0 {
		return fmt.Errorf("-ranges is required when using -ttf")
	}

	// Default height to targetheight if not explicitly provided.
	if height <= 0 {
		height = targetHeight
	}

	// Convert ranges to ttf.CodepointRange.
	ttfRanges := make([]ttf.CodepointRange, len(ranges))
	for i, r := range ranges {
		ttfRanges[i] = ttf.CodepointRange{Low: r.low, High: r.high}
	}

	// Open and parse the TTF file.
	f, err := os.Open(ttfPath)
	if err != nil {
		return fmt.Errorf("cannot open TTF file %s: %w", ttfPath, err)
	}
	defer f.Close()

	font, err := ttf.Parse(f, ttf.ParseConfig{
		Ranges:       ttfRanges,
		TargetHeight: targetHeight,
	})
	if err != nil {
		return fmt.Errorf("parsing TTF: %w", err)
	}

	// Convert font.Glyphs map[rune]*ttf.GlyphData to map[rune][]uint32.
	glyphMap := make(map[rune][]uint32, len(font.Glyphs))
	for cp, gd := range font.Glyphs {
		glyphMap[cp] = gd.Rows
	}

	// Emit the Go source file.
	outFile, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("cannot write %s: %w", outPath, err)
	}
	defer outFile.Close()

	emitCfg := codegen.EmitConfig{
		PackageName:  pkg,
		FontID:       fontID,
		StructName:   structName,
		ConstName:    constName,
		ArrayName:    arrayName,
		GlyphWidth:   width,
		GlyphHeight:  height,
		GlyphAdvance: advance,
		RowHeight:    rowheight,
		GlyphMap:     glyphMap,
	}

	if err := codegen.Emit(outFile, emitCfg); err != nil {
		return fmt.Errorf("emit: %w", err)
	}

	return nil
}

// runBDF handles the BDF input path.
func runBDF(bdfPath, outPath, fontID, pkg string, width, height, advance, rowheight int, cloneName, downloadURL string, ranges []codepointRange, structName, constName, arrayName string) error {
	if bdfPath == "" {
		return fmt.Errorf("missing required flag: -bdf")
	}
	if height <= 0 {
		return fmt.Errorf("missing or invalid flag: -height")
	}
	if cloneName != "" && downloadURL != "" {
		return fmt.Errorf("-clone and -download are mutually exclusive")
	}

	// Resolve the module root so cache paths work regardless of cwd.
	modRoot, err := findModuleRoot()
	if err != nil {
		return fmt.Errorf("cannot locate module root: %w", err)
	}

	// Step 1: Ensure the BDF source is available.
	var fullBDFPath string

	if downloadURL != "" {
		// Download mode: fetch BDF file directly from URL.
		cacheDir := filepath.Join(modRoot, "tools/fontgen/.cache/downloads")
		if err := clone.DownloadFile(downloadURL, cacheDir, bdfPath); err != nil {
			return fmt.Errorf("download: %w", err)
		}
		fullBDFPath = filepath.Join(cacheDir, bdfPath)
	} else {
		// Clone mode: clone/pull upstream repo.
		repoURL, ok := repoURLs[cloneName]
		if !ok {
			return fmt.Errorf("unknown clone target %q (valid: spleen, terminus)", cloneName)
		}
		cacheDir := filepath.Join(modRoot, "tools/fontgen/.cache", cloneName)
		cfg := clone.RepoConfig{
			URL:      repoURL,
			CacheDir: cacheDir,
		}
		if err := clone.EnsureRepo(cfg); err != nil {
			return fmt.Errorf("clone %s: %w", cloneName, err)
		}
		fullBDFPath = filepath.Join(cacheDir, bdfPath)
	}
	f, err := os.Open(fullBDFPath)
	if err != nil {
		return fmt.Errorf("cannot open %s: %w", fullBDFPath, err)
	}
	defer f.Close()

	// Parse with ranges — default to ASCII 32-126 when no ranges provided.
	var bdfCfg bdf.ParseConfig
	if len(ranges) > 0 {
		bdfRanges := make([]bdf.CodepointRange, len(ranges))
		for i, r := range ranges {
			bdfRanges[i] = bdf.CodepointRange{Low: r.low, High: r.high}
		}
		bdfCfg.Ranges = bdfRanges
	} else {
		bdfCfg.Ranges = []bdf.CodepointRange{{Low: 32, High: 126}}
	}

	font, err := bdf.ParseWithConfig(f, bdfCfg)
	if err != nil {
		return fmt.Errorf("%s: %w", fullBDFPath, err)
	}

	// Step 3: Emit the Go source file using the map-based path.
	outFile, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("cannot write %s: %w", outPath, err)
	}
	defer outFile.Close()

	glyphMap := make(map[rune][]uint32, len(font.Glyphs))
	for cp, g := range font.Glyphs {
		glyphMap[cp] = g.Rows
	}

	emitCfg := codegen.EmitConfig{
		PackageName:  pkg,
		FontID:       fontID,
		StructName:   structName,
		ConstName:    constName,
		ArrayName:    arrayName,
		GlyphWidth:   width,
		GlyphHeight:  height,
		GlyphAdvance: advance,
		RowHeight:    rowheight,
		GlyphMap:     glyphMap,
	}

	if err := codegen.Emit(outFile, emitCfg); err != nil {
		return fmt.Errorf("emit: %w", err)
	}

	return nil
}

// idToIdentifier converts a font ID like "spleen-5x8" to "spleen5x8" (unexported).
func idToIdentifier(id string) string {
	var result []byte
	for i := 0; i < len(id); i++ {
		if id[i] == '-' {
			continue
		}
		result = append(result, id[i])
	}
	return string(result)
}

// idToExportedIdentifier converts a font ID like "spleen-5x8" to "Spleen5x8" (exported).
func idToExportedIdentifier(id string) string {
	var result []byte
	capitalize := true
	for i := 0; i < len(id); i++ {
		if id[i] == '-' {
			capitalize = true
			continue
		}
		if capitalize && id[i] >= 'a' && id[i] <= 'z' {
			result = append(result, id[i]-32)
			capitalize = false
		} else {
			result = append(result, id[i])
			capitalize = false
		}
	}
	return string(result)
}

// findModuleRoot walks up from the current working directory looking for go.mod.
func findModuleRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found in any parent directory")
		}
		dir = parent
	}
}
