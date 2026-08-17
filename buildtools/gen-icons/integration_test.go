package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"sort"
	"strings"
	"testing"

	"github.com/databeast/cyberhud/buildtools/fontgen/codegen"
)

// TestIntegration_EndToEnd exercises the full icon generation pipeline from
// codepoints parsing through constants and face file emission, verifying that
// both outputs are valid Go source with the expected structural content.
//
// This test uses Approach A (codegen.Emit with synthetic glyph data) because
// the TTF rasterization step requires a real font file which may not be available
// in all environments. The TTF pipeline is validated separately by property tests
// and the Makefile-driven generation pipeline.
//

func TestIntegration_EndToEnd(t *testing.T) {
	// --- Step 1: Parse a synthetic codepoints input ---
	input := "wifi e63e\nbluetooth e1a7\nsignal_wifi_4_bar e1d8\n"
	entries, err := ParseCodepoints(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ParseCodepoints returned error: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}

	// --- Step 2: Check for naming collisions ---
	if err := CheckCollisions(entries); err != nil {
		t.Fatalf("CheckCollisions returned error: %v", err)
	}

	// --- Step 3: Emit constants file ---
	var constBuf bytes.Buffer
	if err := EmitConstants(&constBuf, "font", entries); err != nil {
		t.Fatalf("EmitConstants returned error: %v", err)
	}

	// --- Step 4: Parse constants file with go/parser ---
	fset := token.NewFileSet()
	constFile, parseErr := parser.ParseFile(fset, "constants.go", constBuf.Bytes(), 0)
	if parseErr != nil {
		t.Fatalf("constants file not parseable by go/parser:\n%s\nerror: %v", constBuf.String(), parseErr)
	}

	// Extract constant names from AST.
	var constNames []string
	for _, decl := range constFile.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.CONST {
			continue
		}
		for _, spec := range genDecl.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for _, ident := range vs.Names {
				constNames = append(constNames, ident.Name)
			}
		}
	}

	// Verify exactly 3 const declarations.
	if len(constNames) != 3 {
		t.Fatalf("expected 3 const declarations, got %d: %v", len(constNames), constNames)
	}

	// Verify constants are in lexicographic order.
	if !sort.StringsAreSorted(constNames) {
		t.Fatalf("constants not in lexicographic order: %v", constNames)
	}

	// Verify expected constant names: IconBluetooth, IconSignalWifi4Bar, IconWifi
	expectedConsts := []string{"IconBluetooth", "IconSignalWifi4Bar", "IconWifi"}
	for i, expected := range expectedConsts {
		if constNames[i] != expected {
			t.Fatalf("constant[%d]: expected %q, got %q", i, expected, constNames[i])
		}
	}

	// --- Step 5: Build a synthetic glyph map (simulating TTF rasterization output) ---
	targetHeight := 24
	glyphMap := make(map[rune][]uint32, 3)
	for _, e := range entries {
		rows := make([]uint32, targetHeight)
		// Fill with non-zero data to simulate real glyph bitmaps.
		for row := 0; row < targetHeight; row++ {
			rows[row] = uint32(e.Codepoint) + uint32(row)
		}
		glyphMap[e.Codepoint] = rows
	}

	// --- Step 6: Emit face file using codegen.Emit ---
	faceID := "material-icons-24"
	structName := idToIdentifier(faceID) + "Face"
	constName := idToIdentifier(faceID) + "ID"
	arrayName := idToIdentifier(faceID)

	emitCfg := codegen.EmitConfig{
		PackageName:  "font",
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

	var faceBuf bytes.Buffer
	if err := codegen.Emit(&faceBuf, emitCfg); err != nil {
		t.Fatalf("codegen.Emit returned error: %v", err)
	}

	// --- Step 7: Parse face file with go/parser ---
	faceAst, parseErr := parser.ParseFile(fset, "face.go", faceBuf.Bytes(), 0)
	if parseErr != nil {
		t.Fatalf("face file not parseable by go/parser:\n%s\nerror: %v", faceBuf.String(), parseErr)
	}

	// Verify face file contains init() function.
	hasInit := false
	hasRegister := false
	for _, decl := range faceAst.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if funcDecl.Name.Name == "init" {
			hasInit = true
		}
	}

	if !hasInit {
		t.Fatal("face file missing init() function")
	}

	// Verify face file contains "Register(" call.
	faceSrc := faceBuf.String()
	if !strings.Contains(faceSrc, "Register(") {
		t.Fatal("face file missing Register( call")
	}
	hasRegister = strings.Contains(faceSrc, "Register(")
	if !hasRegister {
		t.Fatal("face file missing Register( call")
	}

	// Verify face file contains ID method with correct face ID.
	if !strings.Contains(faceSrc, fmt.Sprintf("%q", faceID)) {
		t.Fatalf("face file does not contain face ID %q", faceID)
	}

	// Verify face file contains Metrics with square dimensions.
	expectedMetrics := fmt.Sprintf(
		"Metrics{GlyphWidth: %d, GlyphHeight: %d, GlyphAdvance: %d, RowHeight: %d}",
		targetHeight, targetHeight, targetHeight, targetHeight,
	)
	if !strings.Contains(faceSrc, expectedMetrics) {
		t.Fatalf("face file does not contain expected metrics: %s", expectedMetrics)
	}

	// Verify face file contains GlyphRow method.
	if !strings.Contains(faceSrc, "GlyphRow(") {
		t.Fatal("face file missing GlyphRow method")
	}

	// Verify generated code comment header.
	if !strings.Contains(faceSrc, "// Code generated") {
		t.Fatal("face file missing generated code comment header")
	}
	constSrc := constBuf.String()
	if !strings.Contains(constSrc, "// Code generated") {
		t.Fatal("constants file missing generated code comment header")
	}
}
