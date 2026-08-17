// Command gallery reads collected snapshot images and generates markdown
// gallery sections in each display mode's documentation page.
//
// Usage:
//
//	go run ./tools/docsnap/gallery [flags]
//	  -img string    Root directory containing collected mode images (default "ghpages/docs/display-modes/img")
//	  -pages string  Root directory containing mode markdown pages (default "ghpages/docs/display-modes")
package main

import (
	"flag"
	"fmt"
	"os"

	gallery "github.com/databeast/cyberhud/buildtools/docsnap/internal/gallery"
)

func main() {
	img := flag.String("img", "ghpages/docs/display-modes/img", "Root directory containing collected mode images")
	pages := flag.String("pages", "ghpages/docs/display-modes", "Root directory containing mode markdown pages")
	flag.Parse()

	err := gallery.Run(*img, *pages)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gallery: %v\n", err)
		os.Exit(1)
	}
}
