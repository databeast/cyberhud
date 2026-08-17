package region

// ModeChangeCommand is a ModeSwitch-layer type that represents a request to
// change a Region's display mode at runtime. It is passed to [ModeSwitch.Execute]
// which routes the request through the [RegionManager].
type ModeChangeCommand struct {
	Target string // region name (case-insensitive) or numeric index (as string)
	Mode   string // display mode ID to switch to
}

// ModeSwitch is the ModeSwitch architectural component: the command handler for
// runtime mode changes. It provides a clean command interface layer that routes
// [ModeChangeCommand] requests through the [RegionManager], which in turn invokes
// [Region.SetMode] with full lifecycle management.
type ModeSwitch struct {
	rm *RegionManager
}

// NewModeSwitch creates a [ModeSwitch] handler backed by the given
// [RegionManager]. The rm parameter is the RegionManager that owns the Regions
// whose modes may be changed. Returns the initialized ModeSwitch.
func NewModeSwitch(rm *RegionManager) *ModeSwitch {
	return &ModeSwitch{rm: rm}
}

// Execute processes a [ModeChangeCommand] by resolving the target Region (by
// name or index) and delegating to [RegionManager.SetMode]. The cmd parameter
// carries the target region identifier and the desired mode ID. Returns nil on
// success, or an error if the target region is not found or the mode is not
// registered. On error the current mode remains unchanged.
func (ms *ModeSwitch) Execute(cmd ModeChangeCommand) error {
	return ms.rm.SetMode(cmd.Target, cmd.Mode)
}
