// Command collect copies PNG snapshot images from display mode test output
// directories into the MkDocs documentation source tree.
//
// Usage:
//
//	go run ./tools/docsnap/collect [flags]
//	  -src string    Source root containing mode directories (default "display/modes")
//	  -dst string    Destination root for collected images (default "ghpages/docs/display-modes/img")
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/databeast/cyberhud/buildtools/docsnap/internal/collector"
)

func main() {
	src := flag.String("src", "display/modes", "Source root containing mode directories")
	dst := flag.String("dst", "ghpages/docs/display-modes/img", "Destination root for collected images")
	flag.Parse()

	summary, err := collector.CollectAll(*src, *dst)
	if err != nil {
		fmt.Fprintf(os.Stderr, "collect: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Modes processed: %d\n", summary.ModesProcessed)
	fmt.Printf("Modes skipped:   %d\n", summary.ModesSkipped)
	fmt.Printf("Files copied:    %d\n", summary.FilesCopied)
}
