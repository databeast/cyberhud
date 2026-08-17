package main

import (
	"testing"

	"github.com/databeast/cyberhud/display/regionid"
)

func TestSplitCommands(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want [][]string
	}{
		{
			name: "nil input",
			args: nil,
			want: nil,
		},
		{
			name: "empty input",
			args: []string{},
			want: nil,
		},
		{
			name: "single command no semicolons",
			args: []string{"display", "status"},
			want: [][]string{{"display", "status"}},
		},
		{
			name: "two commands separated by semicolon",
			args: []string{"region", "main.0", ";", "mode", "clock"},
			want: [][]string{{"region", "main.0"}, {"mode", "clock"}},
		},
		{
			name: "three commands",
			args: []string{"region", "main.0", ";", "mode", "attract_matrix", ";", "config", "density=0.8"},
			want: [][]string{{"region", "main.0"}, {"mode", "attract_matrix"}, {"config", "density=0.8"}},
		},
		{
			name: "leading semicolon ignored",
			args: []string{";", "status"},
			want: [][]string{{"status"}},
		},
		{
			name: "trailing semicolon ignored",
			args: []string{"status", ";"},
			want: [][]string{{"status"}},
		},
		{
			name: "consecutive semicolons produce no empty commands",
			args: []string{"region", "main.0", ";", ";", "mode", "clock"},
			want: [][]string{{"region", "main.0"}, {"mode", "clock"}},
		},
		{
			name: "semicolons embedded in args are not separators",
			args: []string{"display", "config", "main.0", "label=hello;world"},
			want: [][]string{{"display", "config", "main.0", "label=hello;world"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SplitCommands(tt.args)
			if !slicesEqual(got, tt.want) {
				t.Errorf("SplitCommands(%v) = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}

func TestRegionContext_ResolveCommand(t *testing.T) {
	mainRegion := regionid.ID{Surface: "main", Index: 0}
	leftAux := regionid.ID{Surface: "left-aux", Index: 0}

	tests := []struct {
		name    string
		active  *regionid.ID
		args    []string
		want    string
		wantErr bool
	}{
		// region command sets context
		{
			name:   "region sets context",
			active: nil,
			args:   []string{"region", "main.0"},
			want:   "",
		},
		{
			name:   "region with bare integer",
			active: nil,
			args:   []string{"region", "0"},
			want:   "",
		},
		{
			name:    "region without id errors",
			active:  nil,
			args:    []string{"region"},
			wantErr: true,
		},
		// mode command
		{
			name:   "mode expands to display set",
			active: &mainRegion,
			args:   []string{"mode", "attract_matrix"},
			want:   "display set main.0 attract_matrix",
		},
		{
			name:   "mode with extra key=value args",
			active: &mainRegion,
			args:   []string{"mode", "attract_matrix", "density=0.8", "speed=5"},
			want:   "display set main.0 attract_matrix density=0.8 speed=5",
		},
		{
			name:    "mode without context errors",
			active:  nil,
			args:    []string{"mode", "clock"},
			wantErr: true,
		},
		{
			name:    "mode without mode_name errors",
			active:  &mainRegion,
			args:    []string{"mode"},
			wantErr: true,
		},
		// config command
		{
			name:   "config expands to display config",
			active: &mainRegion,
			args:   []string{"config", "density=0.8", "trail_length=20"},
			want:   "display config main.0 density=0.8 trail_length=20",
		},
		{
			name:   "config with left-aux region",
			active: &leftAux,
			args:   []string{"config", "style=bold"},
			want:   "display config left-aux.0 style=bold",
		},
		{
			name:    "config without context errors",
			active:  nil,
			args:    []string{"config", "density=0.8"},
			wantErr: true,
		},
		// next command
		{
			name:   "next expands to display next",
			active: &mainRegion,
			args:   []string{"next"},
			want:   "display next main.0",
		},
		{
			name:    "next without context errors",
			active:  nil,
			args:    []string{"next"},
			wantErr: true,
		},
		// prev command
		{
			name:   "prev expands to display prev",
			active: &mainRegion,
			args:   []string{"prev"},
			want:   "display prev main.0",
		},
		{
			name:    "prev without context errors",
			active:  nil,
			args:    []string{"prev"},
			wantErr: true,
		},
		// status command
		{
			name:   "status expands to display policy with region",
			active: &mainRegion,
			args:   []string{"status"},
			want:   "display policy main.0",
		},
		{
			name:   "status with left-aux",
			active: &leftAux,
			args:   []string{"status"},
			want:   "display policy left-aux.0",
		},
		{
			name:    "status without context errors",
			active:  nil,
			args:    []string{"status"},
			wantErr: true,
		},
		// unknown command
		{
			name:    "unknown scoped command errors",
			active:  &mainRegion,
			args:    []string{"foobar"},
			wantErr: true,
		},
		// empty command
		{
			name:    "empty args errors",
			active:  &mainRegion,
			args:    []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &RegionContext{Active: tt.active}
			got, err := ctx.ResolveCommand(tt.args)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil (result=%q)", got)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("ResolveCommand(%v) = %q, want %q", tt.args, got, tt.want)
			}
		})
	}
}

func TestRegionContext_RegionSetsActive(t *testing.T) {
	ctx := &RegionContext{}

	// Initially no active region.
	if ctx.Active != nil {
		t.Fatal("expected nil Active initially")
	}

	// Set region via fully-qualified ID.
	_, err := ctx.ResolveCommand([]string{"region", "left-aux.0"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ctx.Active == nil {
		t.Fatal("expected Active to be set")
	}
	if ctx.Active.Surface != "left-aux" || ctx.Active.Index != 0 {
		t.Errorf("got Active=%+v, want Surface=left-aux Index=0", ctx.Active)
	}

	// Subsequent scoped commands use the context.
	cmd, err := ctx.ResolveCommand([]string{"next"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd != "display next left-aux.0" {
		t.Errorf("got %q, want %q", cmd, "display next left-aux.0")
	}
}

func TestRegionContext_BareIntegerRegion(t *testing.T) {
	ctx := &RegionContext{}

	_, err := ctx.ResolveCommand([]string{"region", "2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Bare integer should be preserved for daemon-side resolution.
	cmd, err := ctx.ResolveCommand([]string{"mode", "clock"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd != "display set 2 clock" {
		t.Errorf("got %q, want %q", cmd, "display set 2 clock")
	}
}

// slicesEqual compares two [][]string for equality.
func slicesEqual(a, b [][]string) bool {
	if len(a) == 0 && len(b) == 0 {
		// Treat nil and empty the same.
		return true
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}
