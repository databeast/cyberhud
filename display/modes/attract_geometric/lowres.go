package attract_geometric

import "github.com/databeast/cyberhud/display/modes/attract_geometric/source"

func tunePolicyForUltraLowRes(p source.Policy, panelWidth, panelHeight int) source.Policy {
	longEdge, shortEdge := panelWidth, panelHeight
	if panelHeight > panelWidth {
		longEdge, shortEdge = panelHeight, panelWidth
	}
	if longEdge > 16 || shortEdge > 8 {
		return p
	}
	if p.Density > 0.2 {
		p.Density = 0.2
	}
	if p.FragmentRate > 0.15 {
		p.FragmentRate = 0.15
	}
	if p.GlowIntensity > 0.6 {
		p.GlowIntensity = 0.6
	}
	if p.Speed > 0.85 {
		p.Speed = 0.85
	}
	return p
}
