package main

import "testing"

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
			args: []string{"display", "set", "main.0", "clock", ";", "status"},
			want: [][]string{{"display", "set", "main.0", "clock"}, {"status"}},
		},
		{
			name: "three commands",
			args: []string{"display", "set", "main.0", "clock", ";", "display", "next", "main.0", ";", "status"},
			want: [][]string{{"display", "set", "main.0", "clock"}, {"display", "next", "main.0"}, {"status"}},
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
			args: []string{"display", "set", "main.0", "clock", ";", ";", "status"},
			want: [][]string{{"display", "set", "main.0", "clock"}, {"status"}},
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

func slicesEqual(a, b [][]string) bool {
	if len(a) == 0 && len(b) == 0 {
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
