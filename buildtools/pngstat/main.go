// Command pngstat reports content statistics for snapshot PNGs.
// Temporary diagnostic used to measure whether rendered frames contain text.
package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"sort"
)

// rowModal returns the most frequent color in a pixel row, breaking ties
// toward the darkest color.
func rowModal(img image.Image, y int) (uint32, uint32, uint32) {
	type key struct{ r, g, b uint32 }
	counts := map[key]int{}
	b := img.Bounds()
	for x := b.Min.X; x < b.Max.X; x++ {
		r, g, bl, _ := img.At(x, y).RGBA()
		counts[key{r >> 8, g >> 8, bl >> 8}]++
	}
	best := key{}
	bestN := -1
	for k, n := range counts {
		if n > bestN || (n == bestN && k.r+k.g+k.b < best.r+best.g+best.b) {
			best, bestN = k, n
		}
	}
	return best.r, best.g, best.b
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func main() {
	dir := os.Args[1]
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".png" {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)

	for _, name := range names {
		f, err := os.Open(filepath.Join(dir, name))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		img, err := png.Decode(f)
		f.Close()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		bounds := img.Bounds()
		fg := 0
		bands := 0
		inBand := false
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			mr, mg, mb := rowModal(img, y)
			occupied := false
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, _ := img.At(x, y).RGBA()
				if abs(int(r>>8)-int(mr)) > 8 || abs(int(g>>8)-int(mg)) > 8 || abs(int(b>>8)-int(mb)) > 8 {
					fg++
					occupied = true
				}
			}
			if occupied && !inBand {
				bands++
			}
			inBand = occupied
		}
		area := bounds.Dx() * bounds.Dy()
		fmt.Printf("%-40s %4dx%-4d fg=%-7d (%5.2f%%) bands=%d\n",
			name, bounds.Dx(), bounds.Dy(), fg, 100*float64(fg)/float64(area), bands)
	}
}
