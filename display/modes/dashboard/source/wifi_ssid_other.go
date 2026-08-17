//go:build !linux

package source

func GetWifiSSID() string {
	return "(no wifi)"
}
