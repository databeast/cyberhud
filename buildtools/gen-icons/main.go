package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	var (
		codepoints   string
		ttfPath      string
		constout     string
		faceout      string
		pkg          string
		targetHeight int
		faceID       string
	)

	flag.StringVar(&codepoints, "codepoints", "", "path to codepoints file")
	flag.StringVar(&ttfPath, "ttf", "", "path to Material Symbols TTF file")
	flag.StringVar(&constout, "constout", "", "output path for icon constants Go file")
	flag.StringVar(&faceout, "faceout", "", "output path for icon face Go file")
	flag.StringVar(&pkg, "pkg", "", "Go package name (e.g. \"font\")")
	flag.IntVar(&targetHeight, "targetheight", 0, "target pixel height (e.g. 24)")
	flag.StringVar(&faceID, "faceid", "", "font.Face ID (e.g. \"material-icons-24\")")

	flag.Parse()

	if err := run(codepoints, ttfPath, constout, faceout, pkg, targetHeight, faceID); err != nil {
		fmt.Fprintf(os.Stderr, "gen-icons: %v\n", err)
		os.Exit(1)
	}
}

func run(codepointsPath, ttfPath, constout, faceout, pkg string, targetHeight int, faceID string) error {
	// Step 1: Validate all required flags.
	if codepointsPath == "" {
		return fmt.Errorf("missing required flag: -codepoints")
	}
	if ttfPath == "" {
		return fmt.Errorf("missing required flag: -ttf")
	}
	if constout == "" {
		return fmt.Errorf("missing required flag: -constout")
	}
	if faceout == "" {
		return fmt.Errorf("missing required flag: -faceout")
	}
	if pkg == "" {
		return fmt.Errorf("missing required flag: -pkg")
	}
	if targetHeight <= 0 {
		return fmt.Errorf("missing or invalid flag: -targetheight")
	}
	if faceID == "" {
		return fmt.Errorf("missing required flag: -faceid")
	}

	// Step 2: Open codepoints file.
	cpFile, err := os.Open(codepointsPath)
	if err != nil {
		return fmt.Errorf("%s: %v", codepointsPath, err)
	}
	defer cpFile.Close()

	// Step 3: Parse codepoints.
	entries, err := ParseCodepoints(cpFile)
	if err != nil {
		return fmt.Errorf("parsing codepoints: %w", err)
	}

	// Step 4: Check for naming collisions.
	if err := CheckCollisions(entries); err != nil {
		return err
	}

	// Step 5: Emit constants file.
	constFile, err := os.Create(constout)
	if err != nil {
		return fmt.Errorf("cannot create %s: %w", constout, err)
	}
	defer constFile.Close()

	if err := EmitConstants(constFile, pkg, entries); err != nil {
		return fmt.Errorf("emitting constants: %w", err)
	}

	// Step 6: Open TTF file and emit face via EmitFace.
	ttfFile, err := os.Open(ttfPath)
	if err != nil {
		return fmt.Errorf("%s: %v", ttfPath, err)
	}
	defer ttfFile.Close()

	faceFile, err := os.Create(faceout)
	if err != nil {
		return fmt.Errorf("cannot create %s: %w", faceout, err)
	}
	defer faceFile.Close()

	if err := EmitFace(faceFile, ttfFile, entries, pkg, faceID, targetHeight); err != nil {
		return err
	}

	return nil
}

// idToIdentifier converts a face ID like "material-icons-24" to "materialIcons24" (unexported camelCase).
// Removes dashes and capitalizes the character following each dash.
func idToIdentifier(id string) string {
	parts := strings.Split(id, "-")
	if len(parts) == 0 {
		return id
	}
	var b strings.Builder
	b.WriteString(parts[0])
	for _, p := range parts[1:] {
		if len(p) == 0 {
			continue
		}
		b.WriteString(strings.ToUpper(p[:1]) + p[1:])
	}
	return b.String()
}
