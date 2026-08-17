package main

import (
	"fmt"
	"io"
	"os"

	"github.com/databeast/cyberhud/buildtools/fontgen/codegen"
	"github.com/databeast/cyberhud/buildtools/fontgen/ttf"
)

// EmitFace handles the TTF rasterization and face file generation.
// It reads a TTF font from ttfReader, rasterizes glyphs for the given entries,
// builds the glyph map, and emits the face source file using codegen.Emit.
func EmitFace(w io.Writer, ttfReader io.Reader, entries []IconEntry, pkg string, faceID string, targetHeight int) error {
	// Convert entries to single-rune CodepointRanges.
	ranges := make([]ttf.CodepointRange, len(entries))
	for i, e := range entries {
		ranges[i] = ttf.CodepointRange{Low: e.Codepoint, High: e.Codepoint}
	}

	// Rasterize TTF glyphs.
	font, err := ttf.Parse(ttfReader, ttf.ParseConfig{
		Ranges:       ranges,
		TargetHeight: targetHeight,
	})
	if err != nil {
		return fmt.Errorf("parsing TTF: %w", err)
	}

	// Build glyph map, warn on missing glyphs.
	glyphMap := make(map[rune][]uint32, len(entries))
	for _, e := range entries {
		gd, ok := font.Glyphs[e.Codepoint]
		if !ok {
			fmt.Fprintf(os.Stderr, "gen-icons: warning: no glyph for %q (U+%04X), skipping\n", e.Name, e.Codepoint)
			continue
		}
		glyphMap[e.Codepoint] = gd.Rows
	}

	// Derive struct/const/array names from faceID.
	structName := idToIdentifier(faceID) + "Face"
	constName := idToIdentifier(faceID) + "ID"
	arrayName := idToIdentifier(faceID)

	// Compose EmitConfig with square metrics.
	emitCfg := codegen.EmitConfig{
		PackageName:  pkg,
		FontID:       faceID,
		StructName:   structName,
		ConstName:    constName,
		ArrayName:    arrayName,
		GlyphWidth:   targetHeight,
		GlyphHeight:  targetHeight,
		GlyphAdvance: targetHeight,
		RowHeight:    targetHeight,
		GlyphMap:     glyphMap,
		FallbackChar: '?',
	}

	// Emit face Go source file.
	if err := codegen.Emit(w, emitCfg); err != nil {
		return fmt.Errorf("emitting face: %w", err)
	}

	return nil
}
