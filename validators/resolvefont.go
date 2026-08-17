package validators

/*


// resolveFont applies tier-based font resolution to the DashboardViewData using the
// instance's TierCatalog. Dashboard uses family "spleen" with TierNormal.
func resolveFont(state *style.ViewData) {
	if i.hints.Catalog.PixelWidth() == 0 {
		panic("dashboard: resolveFont called with zero-valued catalog")
	}

	const family = "spleen"
	const tier = tiercatalog.TierNormal

	if len(state.Tiers) > 0 {
		// Multi-tier resolution: resolve each declared tier into a concrete font ID.
		reqs := make([]tierselect.Request, len(state.Tiers))
		for idx, t := range state.Tiers {
			reqs[idx] = tierselect.Request{
				Family: family,
				Tier:   t,
			}
		}
		faces := tierselect.SelectMulti(i.hints.Catalog, reqs)
		state.FontIDs = make([]string, len(faces))
		for idx, f := range faces {
			state.FontIDs[idx] = f.ID()
		}
		// Set FontID to the first resolved ID for uniform-consumer compatibility.
		if len(state.FontIDs) > 0 {
			state.FontID = state.FontIDs[0]
		}
	} else {
		// Single-tier resolution: use the mode's configured default tier.
		face := tierselect.Select(i.hints.Catalog, tierselect.Request{
			Family: family,
			Tier:   tier,
		})
		state.FontID = face.ID()
		state.FontIDs = nil
	}
}



// resolveFont applies tier-based font resolution to the ViewData using the
// instance's TierCatalog. Thermal uses family "spleen" with TierNormal.
func (i *instance) resolveFont(state *style.ViewData) {
	if i.hints.Catalog.PixelWidth() == 0 {
		panic("thermal: resolveFont called with zero-valued catalog")
	}

	family := thermalFontFamily
	tier := thermalFontTier

	if len(state.Tiers) > 0 {
		// Multi-tier resolution: resolve each declared tier into a concrete font ID.
		reqs := make([]tierselect.Request, len(state.Tiers))
		for idx, t := range state.Tiers {
			reqs[idx] = tierselect.Request{
				Family: family,
				Tier:   t,
			}
		}
		faces := tierselect.SelectMulti(i.hints.Catalog, reqs)
		state.FontIDs = make([]string, len(faces))
		for idx, f := range faces {
			state.FontIDs[idx] = f.ID()
		}
		// Set FontID to the first resolved ID for uniform-consumer compatibility.
		if len(state.FontIDs) > 0 {
			state.FontID = state.FontIDs[0]
		}
	} else {
		// Single-tier resolution: use the thermal mode's configured default tier.
		face := tierselect.Select(i.hints.Catalog, tierselect.Request{
			Family: family,
			Tier:   tier,
		})
		state.FontID = face.ID()
		state.FontIDs = nil
	}
}

*/
