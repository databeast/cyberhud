package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/databeast/cyberhud/buildtools/fontgen/codegen"
)

// TestEmitFaceRegistration verifies that codegen.Emit produces source
// containing init() calling Register(materialIcons24Face{}) for a 24-pixel
// icon face, and that the face ID is "material-icons-24".
//

func TestEmitFaceRegistration(t *testing.T) {
	// Build a minimal glyph map (single glyph, 24 rows).
	glyphMap := map[rune][]uint32{
		0xe63e: make([]uint32, 24), // wifi icon placeholder
	}

	faceID := "material-icons-24"
	structName := idToIdentifier(faceID) + "Face"
	constName := idToIdentifier(faceID) + "ID"
	arrayName := idToIdentifier(faceID)

	cfg := codegen.EmitConfig{
		PackageName:  "font",
		FontID:       faceID,
		StructName:   structName,
		ConstName:    constName,
		ArrayName:    arrayName,
		GlyphWidth:   24,
		GlyphHeight:  24,
		GlyphAdvance: 24,
		RowHeight:    24,
		GlyphMap:     glyphMap,
		FallbackChar: '?',
	}

	var buf bytes.Buffer
	if err := codegen.Emit(&buf, cfg); err != nil {
		t.Fatalf("codegen.Emit failed: %v", err)
	}

	src := buf.String()

	// Verify init() function is present.
	if !strings.Contains(src, "func init()") {
		t.Fatal("generated source does not contain 'func init()'")
	}

	// Verify Register(materialIcons24Face{}) is called in init().
	expectedRegister := fmt.Sprintf("Register(%s{})", structName)
	if !strings.Contains(src, expectedRegister) {
		t.Fatalf("generated source does not contain %q", expectedRegister)
	}

	// Verify the face ID constant is present.
	expectedIDConst := fmt.Sprintf("%s = %q", constName, faceID)
	if !strings.Contains(src, expectedIDConst) {
		t.Fatalf("generated source does not contain ID constant %q", expectedIDConst)
	}

	// Verify square metrics line (all four fields = 24).
	expectedMetrics := "Metrics{GlyphWidth: 24, GlyphHeight: 24, GlyphAdvance: 24, RowHeight: 24}"
	if !strings.Contains(src, expectedMetrics) {
		t.Fatalf("generated source does not contain expected metrics %q", expectedMetrics)
	}

	// Verify the struct name follows the expected pattern.
	if structName != "materialIcons24Face" {
		t.Fatalf("structName = %q, want %q", structName, "materialIcons24Face")
	}
}

// TestEmitFaceRegistrationContainsID verifies the generated face source includes
// the face ID string "material-icons-24" so that fonts.Get("material-icons-24")
// will succeed after registration.
func TestEmitFaceRegistrationContainsID(t *testing.T) {
	glyphMap := map[rune][]uint32{
		'A': make([]uint32, 24),
	}

	faceID := "material-icons-24"
	cfg := codegen.EmitConfig{
		PackageName:  "font",
		FontID:       faceID,
		StructName:   "materialIcons24Face",
		ConstName:    "materialIcons24ID",
		ArrayName:    "materialIcons24",
		GlyphWidth:   24,
		GlyphHeight:  24,
		GlyphAdvance: 24,
		RowHeight:    24,
		GlyphMap:     glyphMap,
		FallbackChar: '?',
	}

	var buf bytes.Buffer
	if err := codegen.Emit(&buf, cfg); err != nil {
		t.Fatalf("codegen.Emit failed: %v", err)
	}

	src := buf.String()

	// The source must contain the literal face ID string so that
	// the ID() method returns "material-icons-24".
	if !strings.Contains(src, `"material-icons-24"`) {
		t.Fatal("generated source does not contain the face ID string \"material-icons-24\"")
	}

	// The ID() method must reference the const.
	if !strings.Contains(src, "ID() string { return materialIcons24ID }") {
		t.Fatal("generated source does not contain proper ID() method")
	}
}
