package source

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"

	nlwifi "github.com/mdlayher/wifi"
)

func TestSignalInfoStationAndStickyFallback(t *testing.T) {
	var g Gatherer
	calls := 0
	read := func() ([]*nlwifi.StationInfo, error) {
		calls++
		if calls == 1 {
			return []*nlwifi.StationInfo{{Signal: -47, TransmitBitrate: 54000000}}, nil
		}
		return nil, errors.New("station failure")
	}

	signal, speed, unavailable := g.signalInfo("wlan0", -7250, read)
	if signal != -47 || speed != 54 || unavailable {
		t.Fatalf("station signal = %d, speed = %d, unavailable = %t", signal, speed, unavailable)
	}
	signal, speed, unavailable = g.signalInfo("wlan0", -7250, read)
	if signal != -72 || speed != 0 || unavailable {
		t.Fatalf("fallback signal = %d, speed = %d, unavailable = %t", signal, speed, unavailable)
	}
	signal, speed, unavailable = g.signalInfo("wlan0", -6100, read)
	if signal != -61 || speed != 0 || unavailable || calls != 2 {
		t.Fatalf("sticky fallback signal = %d, speed = %d, unavailable = %t, calls = %d", signal, speed, unavailable, calls)
	}

	g.Reset()
	g.signalInfo("wlan0", -6100, read)
	if calls != 3 {
		t.Fatalf("calls after reset = %d, want 3", calls)
	}
}

func TestSignalInfoEmptyAndMissingBSSSignal(t *testing.T) {
	var g Gatherer
	calls := 0
	read := func() ([]*nlwifi.StationInfo, error) {
		calls++
		return nil, nil
	}
	signal, speed, unavailable := g.signalInfo("wlan0", 0, read)
	if signal != -100 || speed != 0 || !unavailable {
		t.Fatalf("missing BSS signal = %d, speed = %d, unavailable = %t", signal, speed, unavailable)
	}
	signal, _, unavailable = g.signalInfo("wlan0", -8000, read)
	if signal != -80 || unavailable || calls != 1 {
		t.Fatalf("recovered BSS signal = %d, unavailable = %t, calls = %d", signal, unavailable, calls)
	}
}

func TestSignalInfoConcurrentFailureOnlyQueriesOnce(t *testing.T) {
	var g Gatherer
	var calls atomic.Int32
	read := func() ([]*nlwifi.StationInfo, error) {
		calls.Add(1)
		return nil, errors.New("station failure")
	}
	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			g.signalInfo("wlan0", -7000, read)
		}()
	}
	wg.Wait()
	if got := calls.Load(); got != 1 {
		t.Fatalf("station calls = %d, want 1", got)
	}
}
