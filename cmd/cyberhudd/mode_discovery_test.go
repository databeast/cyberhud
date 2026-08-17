package main

import (
	"testing"

	"github.com/databeast/cyberhud/display/catalog"
	displaymodes "github.com/databeast/cyberhud/display/modes"
)

func TestModeDiscovery_Testicons(t *testing.T) {
	def, ok := catalog.Describe("testicons")
	if !ok {
		t.Fatal("testicons not found in catalog")
	}
	if def.ID != "testicons" {
		t.Fatalf("catalog ID = %q, want %q", def.ID, "testicons")
	}

	inst, ok := displaymodes.GetInstance("testicons")
	if !ok {
		t.Fatal("testicons mode factory not found")
	}
	if inst.ID() != "testicons" {
		t.Fatalf("instance ID = %q, want %q", inst.ID(), "testicons")
	}
}

func TestModeDiscovery_Testwidgets(t *testing.T) {
	def, ok := catalog.Describe("testwidgets")
	if !ok {
		t.Fatal("testwidgets not found in catalog")
	}
	if def.ID != "testwidgets" {
		t.Fatalf("catalog ID = %q, want %q", def.ID, "testwidgets")
	}

	inst, ok := displaymodes.GetInstance("testwidgets")
	if !ok {
		t.Fatal("testwidgets mode factory not found")
	}
	if inst.ID() != "testwidgets" {
		t.Fatalf("instance ID = %q, want %q", inst.ID(), "testwidgets")
	}
}
