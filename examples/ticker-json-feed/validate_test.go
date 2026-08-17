package ticker_json_feed_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/databeast/cyberhud/display/modes/ticker/source"
)

// Property 10: Example file validity
// For all JSON files in examples/ticker-json-feed/, read file content,
// call ParseJSONFeed, assert non-nil slice and nil error.

func TestExampleFilesParseSuccessfully(t *testing.T) {
	dir := "."
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read dir: %v", err)
	}

	jsonFiles := 0
	for _, e := range entries {
		if filepath.Ext(e.Name()) != ".json" {
			continue
		}
		jsonFiles++
		t.Run(e.Name(), func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join(dir, e.Name()))
			if err != nil {
				t.Fatalf("read file: %v", err)
			}
			directives, err := source.ParseJSONFeed(string(data))
			if err != nil {
				t.Fatalf("ParseJSONFeed failed: %v", err)
			}
			if len(directives) == 0 {
				t.Fatalf("expected at least one directive, got 0")
			}
			t.Logf("parsed %d directives", len(directives))
		})
	}

	if jsonFiles == 0 {
		t.Fatal("no JSON files found in examples/ticker-json-feed/")
	}
	if jsonFiles != 7 {
		t.Fatalf("expected 7 JSON files, found %d", jsonFiles)
	}
}
