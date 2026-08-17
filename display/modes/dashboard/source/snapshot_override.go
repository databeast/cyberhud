package source

import "sync"

// Package-level override state, protected by mutex.
var testOverride struct {
	sync.RWMutex
	data *DashboardContent // nil = not set
}

// SetTestDashboardContent stores a test override. When set, BuildDashboardContent()
// returns this value instead of querying system state.
func SetTestDashboardContent(data DashboardContent) {
	testOverride.Lock()
	defer testOverride.Unlock()
	testOverride.data = &data
}

// ResetTestDashboardContent clears the override, restoring normal system queries.
func ResetTestDashboardContent() {
	testOverride.Lock()
	defer testOverride.Unlock()
	testOverride.data = nil
}

// getTestOverride returns the override if set. Thread-safe read path.
func getTestOverride() (DashboardContent, bool) {
	testOverride.RLock()
	defer testOverride.RUnlock()
	if testOverride.data == nil {
		return DashboardContent{}, false
	}
	return *testOverride.data, true
}
