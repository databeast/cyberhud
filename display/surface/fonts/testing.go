package font

// SnapshotAndClear saves the current registry state and clears it.
// The returned function restores the original state. Intended for
// property tests that need to control the registry contents.
func SnapshotAndClear() func() {
	facesMu.Lock()
	saved := faces
	faces = map[string]Face{}
	facesMu.Unlock()
	return func() {
		facesMu.Lock()
		faces = saved
		facesMu.Unlock()
	}
}
