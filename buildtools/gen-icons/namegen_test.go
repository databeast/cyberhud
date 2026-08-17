package main

import "testing"

func TestSnakeToPascal(t *testing.T) {
	tests := []struct{ in, want string }{
		{"signal_wifi_4_bar", "IconSignalWifi4Bar"},
		{"wifi", "IconWifi"},
		{"bluetooth", "IconBluetooth"},
		{"a__b", "IconAB"},
		{"x", "IconX"},
	}
	for _, tc := range tests {
		got := SnakeToPascal(tc.in)
		if got != tc.want {
			t.Errorf("SnakeToPascal(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestCheckCollisions_NoCollision(t *testing.T) {
	entries := []IconEntry{
		{Name: "wifi", Codepoint: 0xe63e},
		{Name: "bluetooth", Codepoint: 0xe1a7},
	}
	if err := CheckCollisions(entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheckCollisions_Collision(t *testing.T) {
	// "a_b" and "a__b" both produce "IconAB" — wait, no they don't.
	// Let's use entries that genuinely collide. Actually with our implementation
	// both "a_b" and "a_b" (same name) would collide, but they're the same name.
	// A real collision: there isn't one easily with normal names. Let's just
	// test the error path with two entries having the same name.
	entries := []IconEntry{
		{Name: "wifi", Codepoint: 0xe63e},
		{Name: "wifi", Codepoint: 0xe63f},
	}
	err := CheckCollisions(entries)
	if err == nil {
		t.Fatal("expected collision error, got nil")
	}
	if got := err.Error(); got == "" {
		t.Fatal("expected non-empty error message")
	}
}
