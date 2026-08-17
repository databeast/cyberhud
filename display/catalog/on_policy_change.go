package catalog

import (
	"encoding/json"

	"github.com/databeast/cyberhud/display/catalog/cmdutil"
)

// OnPolicyChange is an optional callback invoked whenever a mode's policy
// changes via the CmdHandler key=value application path. It receives the
// mode ID and the serialized policy snapshot (JSON). The daemon's PolicyStore
// registers itself as the callback target at startup.
var OnPolicyChange func(modeID string, data json.RawMessage)

func init() {
	// Wire cmdutil.OnApplied → catalog.OnPolicyChange bridge.
	// When any CmdHandler successfully applies key=value pairs, snapshot
	// the mode's policy and forward to the OnPolicyChange listener.
	cmdutil.OnApplied = func(verb string) {
		if OnPolicyChange == nil {
			return
		}
		snap := Snapshotter(verb)
		if snap == nil {
			return
		}
		m := snap.SnapshotPolicy()
		data, err := json.Marshal(m)
		if err != nil {
			return
		}
		OnPolicyChange(verb, json.RawMessage(data))
	}
}
