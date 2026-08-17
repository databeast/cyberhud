package source

import "time"

// ThermalSnapshot captures the state of all thermal zones at a point in time.
type ThermalSnapshot struct {
	Zones     []ZoneReading
	Timestamp time.Time
}
