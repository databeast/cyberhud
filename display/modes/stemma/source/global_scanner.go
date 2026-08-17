package source

import (
	"sync"
)

// globalScannerState holds the package-level scanner singleton.
// It is set by the runtime wiring (cmd/cyberhudd) at startup via SetGlobalScanner
// and read by the stemma ModeInstance factory via GlobalScanner.
var globalScannerState struct {
	sync.RWMutex
	scanner *Scanner
}

// SetGlobalScanner sets the package-level scanner singleton.
// Call this from the daemon wiring (cmd/cyberhudd) after constructing the Scanner,
// before any mode instances are constructed.
func SetGlobalScanner(s *Scanner) {
	globalScannerState.Lock()
	defer globalScannerState.Unlock()
	globalScannerState.scanner = s
}

// GlobalScanner returns the package-level scanner singleton, or nil if
// no scanner has been wired (e.g., headless mode or no I2C available).
func GlobalScanner() *Scanner {
	globalScannerState.RLock()
	defer globalScannerState.RUnlock()
	return globalScannerState.scanner
}
