package main

// State is the combined read model the web UI renders: the list of regions with
// their current and available modes, plus the virtual-display geometry.
type State struct {
	Regions  []Region
	Geometry Geometry
}

// RegionByID returns the region with the given addressing ID (e.g. "main.0"),
// or nil if not found.
func (s *State) RegionByID(id string) *Region {
	for i := range s.Regions {
		if s.Regions[i].ID == id {
			return &s.Regions[i]
		}
	}
	return nil
}

// Fetcher is the minimal daemon interface the HTTP handlers depend on. It is
// satisfied by *Client and by test fakes.
type Fetcher interface {
	// FetchState assembles the full read model from the daemon.
	FetchState() (*State, error)
	// Send issues a raw control command (e.g. "display set main.0 clock") and
	// returns the response lines; used by the fire-and-confirm control handlers.
	Send(cmd string) ([]string, error)
}

// FetchState implements Fetcher for the live daemon Client. It issues the three
// read commands and assembles them into a single State. "display status" is the
// authoritative region list; "display modes" enriches mode metadata; "display
// geometry" provides placement. Geometry/modes failures degrade gracefully (the
// page still renders a region list) but a status failure is fatal.
func (c *Client) FetchState() (*State, error) {
	statusLines, err := c.Send("display status")
	if err != nil {
		return nil, err
	}
	if IsErr(statusLines) {
		return nil, &daemonError{msg: ErrText(statusLines)}
	}
	regions := parseStatus(statusLines)

	if modeLines, err := c.Send("display modes"); err == nil && !IsErr(modeLines) {
		regions = mergeModes(regions, parseModes(modeLines))
	}

	var geo Geometry
	if geoLines, err := c.Send("display geometry"); err == nil && !IsErr(geoLines) {
		geo = parseGeometry(geoLines)
	}

	// Fetch per-mode policy values so each carousel card can render its own
	// config editor prefilled with current values. The active mode's values come
	// from the region's live policy ("display policy <region>"); every other mode
	// with options is queried by its global policy snapshot ("display policy
	// <mode-id>"). Modes with no options are skipped. Failures degrade gracefully
	// (the editor falls back to option defaults). Snapshots are cached per mode ID
	// so a mode appearing in several regions is only fetched once.
	modePolicyCache := map[string]map[string]string{}
	for i := range regions {
		for j := range regions[i].Modes {
			m := &regions[i].Modes[j]
			if len(m.Options) == 0 {
				continue
			}
			if regions[i].IsActive(m.ID) {
				if polLines, err := c.Send("display policy " + regions[i].ID); err == nil && !IsErr(polLines) {
					m.Policy = parsePolicy(polLines)
				}
				continue
			}
			if cached, ok := modePolicyCache[m.ID]; ok {
				m.Policy = cached
				continue
			}
			if polLines, err := c.Send("display policy " + m.ID); err == nil && !IsErr(polLines) {
				m.Policy = parsePolicy(polLines)
				modePolicyCache[m.ID] = m.Policy
			}
		}
	}

	return &State{Regions: regions, Geometry: geo}, nil
}

// daemonError wraps an ERR response from the daemon.
type daemonError struct{ msg string }

func (e *daemonError) Error() string { return "daemon: " + e.msg }
