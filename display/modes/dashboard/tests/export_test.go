package tests

import (
	"github.com/databeast/cyberhud/display/modes/dashboard"
	"github.com/databeast/cyberhud/display/modes/dashboard/source"
)

// NormalizePolicyExported exports normalizePolicy for property-based testing.
var NormalizePolicyExported = dashboard.NormalizePolicy

// FormatUptimeExported exports formatUptime for property-based testing.
var FormatUptimeExported = source.FormatUptime

// RenderCacheKeyInternalExported exports renderCacheKeyInternal for property-based testing.
var RenderCacheKeyInternalExported = dashboard.RenderCacheKey

// DashboardRegistryExported exports dashboardRegistry for property-based testing.
var DashboardRegistryExported = dashboard.Registry()

// AllowedAccentsExported exports allowedAccents for property-based testing.
var AllowedAccentsExported = source.AllowedAccents

// BuildSnapshotExported exports buildSnapshot for property-based testing.
var BuildSnapshotExported = source.BuildDashboardContent

// GetUptimeExported exports getUptime for property-based testing.
var GetUptimeExported = source.GetUptime

// GetWifiSSIDExported exports getWifiSSID for property-based testing.
var GetWifiSSIDExported = source.GetWifiSSID

// GetVersionExported exports getVersion for property-based testing.
var GetVersionExported = source.GetVersion

// GetPanelNameExported exports getPanelName for property-based testing.
var GetPanelNameExported = source.GetPanelName
