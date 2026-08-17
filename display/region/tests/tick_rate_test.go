package tests_test

import (
	"testing"
	"time"

	tickermode "github.com/databeast/cyberhud/display/modes/ticker"
	"github.com/databeast/cyberhud/display/region"
)

func TestDefaultTickRateResolver_TickerMode(t *testing.T) {
	resolver := &region.DefaultTickRateResolver{}

	t.Run("AutoScrollMS=50 returns 50ms", func(t *testing.T) {
		tickermode.SetPolicy(tickermode.Policy{
			Style:        "plain",
			Font:         "auto",
			LineMode:     "truncate",
			Direction:    "vertical",
			AutoScrollMS: 50,
		})
		defer tickermode.SetPolicy(tickermode.DefaultPolicy())

		got := resolver.TickInterval("ticker")
		want := 50 * time.Millisecond
		if got != want {
			t.Errorf("TickInterval(ticker) = %v, want %v", got, want)
		}
	})

	t.Run("AutoScrollMS=200 returns 200ms", func(t *testing.T) {
		tickermode.SetPolicy(tickermode.Policy{
			Style:        "plain",
			Font:         "auto",
			LineMode:     "truncate",
			Direction:    "vertical",
			AutoScrollMS: 200,
		})
		defer tickermode.SetPolicy(tickermode.DefaultPolicy())

		got := resolver.TickInterval("ticker")
		want := 200 * time.Millisecond
		if got != want {
			t.Errorf("TickInterval(ticker) = %v, want %v", got, want)
		}
	})

	t.Run("AutoScrollMS=0 returns default 1000ms", func(t *testing.T) {
		tickermode.SetPolicy(tickermode.Policy{
			Style:        "plain",
			Font:         "auto",
			LineMode:     "truncate",
			Direction:    "vertical",
			AutoScrollMS: 0,
		})
		defer tickermode.SetPolicy(tickermode.DefaultPolicy())

		got := resolver.TickInterval("ticker")
		want := 1000 * time.Millisecond
		if got != want {
			t.Errorf("TickInterval(ticker) with AutoScrollMS=0 = %v, want %v", got, want)
		}
	})

	t.Run("negative AutoScrollMS is normalized to 0 by SetPolicy, returns default 1000ms", func(t *testing.T) {
		// SetPolicy normalizes negative values to 0.
		tickermode.SetPolicy(tickermode.Policy{
			Style:        "plain",
			Font:         "auto",
			LineMode:     "truncate",
			Direction:    "vertical",
			AutoScrollMS: -100,
		})
		defer tickermode.SetPolicy(tickermode.DefaultPolicy())

		got := resolver.TickInterval("ticker")
		want := 1000 * time.Millisecond
		if got != want {
			t.Errorf("TickInterval(ticker) with negative AutoScrollMS = %v, want %v", got, want)
		}
	})

	t.Run("AutoScrollMS exceeding 10s is clamped to 10s", func(t *testing.T) {
		tickermode.SetPolicy(tickermode.Policy{
			Style:        "plain",
			Font:         "auto",
			LineMode:     "truncate",
			Direction:    "vertical",
			AutoScrollMS: 15000,
		})
		defer tickermode.SetPolicy(tickermode.DefaultPolicy())

		got := resolver.TickInterval("ticker")
		want := 10000 * time.Millisecond
		if got != want {
			t.Errorf("TickInterval(ticker) with AutoScrollMS=15000 = %v, want %v", got, want)
		}
	})
}

func TestDefaultTickRateResolver_NonTickerModes(t *testing.T) {
	resolver := &region.DefaultTickRateResolver{}

	t.Run("dashboard mode returns 1000ms default", func(t *testing.T) {
		got := resolver.TickInterval("dashboard")
		want := 1000 * time.Millisecond
		if got != want {
			t.Errorf("TickInterval(dashboard) = %v, want %v", got, want)
		}
	})

	t.Run("unknown mode returns 1000ms default", func(t *testing.T) {
		got := resolver.TickInterval("nonexistent-mode")
		want := 1000 * time.Millisecond
		if got != want {
			t.Errorf("TickInterval(nonexistent-mode) = %v, want %v", got, want)
		}
	})

	t.Run("empty mode string returns 1000ms default", func(t *testing.T) {
		got := resolver.TickInterval("")
		want := 1000 * time.Millisecond
		if got != want {
			t.Errorf("TickInterval('') = %v, want %v", got, want)
		}
	})
}

func TestClampInterval(t *testing.T) {
	resolver := &region.DefaultTickRateResolver{}

	t.Run("value within range is unchanged", func(t *testing.T) {
		// Register a provider that returns 500ms
		region.RegisterTickRate("__test_clamp_500ms", &fixedProvider{500 * time.Millisecond})
		got := resolver.TickInterval("__test_clamp_500ms")
		want := 500 * time.Millisecond
		if got != want {
			t.Errorf("TickInterval with 500ms provider = %v, want %v", got, want)
		}
	})

	t.Run("value below min is clamped to 1ms", func(t *testing.T) {
		// Register a provider that returns 0
		region.RegisterTickRate("__test_clamp_0", &fixedProvider{0})
		got := resolver.TickInterval("__test_clamp_0")
		want := 1 * time.Millisecond
		if got != want {
			t.Errorf("TickInterval with 0 provider = %v, want %v", got, want)
		}
	})

	t.Run("value above max is clamped to 10000ms", func(t *testing.T) {
		// Register a provider that returns 20000ms
		region.RegisterTickRate("__test_clamp_20s", &fixedProvider{20000 * time.Millisecond})
		got := resolver.TickInterval("__test_clamp_20s")
		want := 10000 * time.Millisecond
		if got != want {
			t.Errorf("TickInterval with 20000ms provider = %v, want %v", got, want)
		}
	})

	t.Run("exact min boundary returns min", func(t *testing.T) {
		region.RegisterTickRate("__test_clamp_1ms", &fixedProvider{1 * time.Millisecond})
		got := resolver.TickInterval("__test_clamp_1ms")
		want := 1 * time.Millisecond
		if got != want {
			t.Errorf("TickInterval with 1ms provider = %v, want %v", got, want)
		}
	})

	t.Run("exact max boundary returns max", func(t *testing.T) {
		region.RegisterTickRate("__test_clamp_10s", &fixedProvider{10000 * time.Millisecond})
		got := resolver.TickInterval("__test_clamp_10s")
		want := 10000 * time.Millisecond
		if got != want {
			t.Errorf("TickInterval with 10000ms provider = %v, want %v", got, want)
		}
	})
}

// fixedProvider is a test helper that always returns a fixed interval.
type fixedProvider struct {
	interval time.Duration
}

func (f *fixedProvider) PreferredTickInterval() time.Duration {
	return f.interval
}
