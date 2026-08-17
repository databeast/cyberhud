package source

import "sync"

// Package-level override state, protected by mutex.
var testOverride struct {
	sync.RWMutex
	data *WifiData // nil = not set
}

// SetTestWifiState stores a test override. When set, GatherWifiState()
// returns this value instead of querying hardware.
func SetTestWifiState(data WifiData) {
	testOverride.Lock()
	defer testOverride.Unlock()
	testOverride.data = &data
}

// ResetTestWifiState clears the override, restoring normal hardware queries.
func ResetTestWifiState() {
	testOverride.Lock()
	defer testOverride.Unlock()
	testOverride.data = nil
}

// getTestOverride returns the override if set. Thread-safe read path.
func getTestOverride() (WifiData, bool) {
	testOverride.RLock()
	defer testOverride.RUnlock()
	if testOverride.data == nil {
		return WifiData{}, false
	}
	return *testOverride.data, true
}
