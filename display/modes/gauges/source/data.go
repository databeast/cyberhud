package source

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Gauge is the canonical per-widget data model used by gauges mode.
type Gauge struct {
	Label    string  `json:"label"`
	Value    float64 `json:"value"`
	Min      float64 `json:"min"`
	Max      float64 `json:"max"`
	Percent  float64 `json:"percent"`
	Shape    string  `json:"shape,omitempty"`
	Accent   string  `json:"accent,omitempty"`
	Status   string  `json:"status,omitempty"`
	RawValue string  `json:"raw_value,omitempty"`
}

// GaugeSet is the canonical snapshot rendered by the mode.
type GaugeSet struct {
	Gauges []Gauge `json:"gauges"`
	Source string  `json:"source,omitempty"`
}

func CloneGaugeSet(s GaugeSet) GaugeSet {
	out := s
	if len(s.Gauges) > 0 {
		out.Gauges = append([]Gauge(nil), s.Gauges...)
	}
	return out
}

func (s GaugeSet) Fingerprint() string {
	if len(s.Gauges) == 0 {
		return "gauges:empty"
	}
	var b strings.Builder
	for _, g := range s.Gauges {
		fmt.Fprintf(&b, "|%s|%g|%g|%g|%g|%s|%s|%s",
			g.Label, g.Value, g.Min, g.Max, g.Percent, g.Shape, g.Accent, g.Status)
	}
	return b.String()
}

func SerializeSnapshot(s GaugeSet) (string, error) {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return "", fmt.Errorf("serialize gauges snapshot: %v", err)
	}
	return string(data), nil
}

// ParsePayload accepts a top-level object, array, or number and returns a normalized snapshot.
func ParsePayload(payload string, pol Policy) (GaugeSet, error) {
	payload = strings.TrimSpace(payload)
	if payload == "" {
		return GaugeSet{}, nil
	}

	switch payload[0] {
	case '[':
		return parseArrayPayload(payload, pol)
	case '{':
		return parseObjectPayload(payload, pol)
	default:
		if v, err := strconv.ParseFloat(payload, 64); err == nil {
			return GaugeSet{
				Gauges: []Gauge{normalizedGauge(Gauge{Label: "gauge-1", Value: v}, pol, 0, false, false)},
			}, nil
		}
		return GaugeSet{}, fmt.Errorf("expected JSON object, array, or number")
	}
}

func parseArrayPayload(payload string, pol Policy) (GaugeSet, error) {
	var raw []json.RawMessage
	if err := json.Unmarshal([]byte(payload), &raw); err != nil {
		return GaugeSet{}, fmt.Errorf("invalid JSON array: %v", err)
	}

	out := GaugeSet{}
	for i, elem := range raw {
		g, ok, err := parseGaugeElement(elem, pol, i)
		if err != nil {
			return GaugeSet{}, err
		}
		if ok {
			out.Gauges = append(out.Gauges, g)
		}
	}
	if len(out.Gauges) == 0 {
		return GaugeSet{}, fmt.Errorf("no valid gauges found")
	}
	return out, nil
}

func parseObjectPayload(payload string, pol Policy) (GaugeSet, error) {
	var env map[string]json.RawMessage
	if err := json.Unmarshal([]byte(payload), &env); err != nil {
		return GaugeSet{}, fmt.Errorf("invalid JSON object: %v", err)
	}
	if gaugesRaw, ok := env["gauges"]; ok {
		return parseArrayPayload(string(gaugesRaw), pol)
	}

	g, ok, err := parseGaugeObject(env, pol, 0)
	if err != nil {
		return GaugeSet{}, err
	}
	if !ok {
		return GaugeSet{}, fmt.Errorf("missing gauge data")
	}
	return GaugeSet{Gauges: []Gauge{g}}, nil
}

func parseGaugeElement(raw json.RawMessage, pol Policy, index int) (Gauge, bool, error) {
	raw = json.RawMessage(strings.TrimSpace(string(raw)))
	if len(raw) == 0 {
		return Gauge{}, false, nil
	}
	switch raw[0] {
	case '{':
		var env map[string]json.RawMessage
		if err := json.Unmarshal(raw, &env); err != nil {
			return Gauge{}, false, fmt.Errorf("gauge %d: %v", index, err)
		}
		return parseGaugeObject(env, pol, index)
	case '"':
		var text string
		if err := json.Unmarshal(raw, &text); err != nil {
			return Gauge{}, false, fmt.Errorf("gauge %d: %v", index, err)
		}
		if strings.TrimSpace(text) == "" {
			return Gauge{}, false, nil
		}
		return normalizedGauge(Gauge{Label: text, Value: 0}, pol, index, false, false), true, nil
	default:
		var v float64
		if err := json.Unmarshal(raw, &v); err != nil {
			return Gauge{}, false, fmt.Errorf("gauge %d: %v", index, err)
		}
		return normalizedGauge(Gauge{Label: fmt.Sprintf("gauge-%d", index+1), Value: v}, pol, index, false, false), true, nil
	}
}

func parseGaugeObject(env map[string]json.RawMessage, pol Policy, index int) (Gauge, bool, error) {
	var g Gauge
	var hasValue bool
	var hasMin bool
	var hasMax bool

	if raw, ok := env["label"]; ok {
		_ = json.Unmarshal(raw, &g.Label)
	}
	if g.Label == "" {
		if raw, ok := env["name"]; ok {
			_ = json.Unmarshal(raw, &g.Label)
		}
	}
	if raw, ok := env["value"]; ok {
		hasValue = true
		_ = json.Unmarshal(raw, &g.Value)
	}
	if !hasValue {
		if raw, ok := env["current"]; ok {
			hasValue = true
			_ = json.Unmarshal(raw, &g.Value)
		}
	}
	if raw, ok := env["min"]; ok {
		hasMin = true
		_ = json.Unmarshal(raw, &g.Min)
	}
	if raw, ok := env["max"]; ok {
		hasMax = true
		_ = json.Unmarshal(raw, &g.Max)
	}
	if raw, ok := env["shape"]; ok {
		_ = json.Unmarshal(raw, &g.Shape)
	}
	if g.Shape == "" {
		if raw, ok := env["style"]; ok {
			_ = json.Unmarshal(raw, &g.Shape)
		}
	}
	if raw, ok := env["accent"]; ok {
		_ = json.Unmarshal(raw, &g.Accent)
	}
	if g.Accent == "" {
		if raw, ok := env["color"]; ok {
			_ = json.Unmarshal(raw, &g.Accent)
		}
	}
	if raw, ok := env["status"]; ok {
		_ = json.Unmarshal(raw, &g.Status)
	}

	if !hasValue {
		return Gauge{}, false, nil
	}
	if strings.TrimSpace(g.Label) == "" {
		g.Label = fmt.Sprintf("gauge-%d", index+1)
	}
	return normalizedGauge(g, pol, index, hasMin, hasMax), true, nil
}

func normalizedGauge(g Gauge, pol Policy, index int, hasMin, hasMax bool) Gauge {
	min := pol.DefaultMin
	max := pol.DefaultMax
	if hasMin {
		if g.Min != 0 || min == 0 {
			min = g.Min
		}
	}
	if hasMax {
		if g.Max != 0 || max == 0 {
			max = g.Max
		}
	}
	if max <= min {
		max = min + 1
	}

	value := g.Value
	if math.IsNaN(value) || math.IsInf(value, 0) {
		value = min
	}
	if value < min {
		value = min
	}
	if value > max {
		value = max
	}
	percent := (value - min) / (max - min)
	if percent < 0 {
		percent = 0
	}
	if percent > 1 {
		percent = 1
	}

	label := strings.TrimSpace(g.Label)
	if label == "" {
		label = fmt.Sprintf("gauge-%d", index+1)
	}

	shape := strings.ToLower(strings.TrimSpace(g.Shape))
	if !ValidShape(shape) || shape == "auto" {
		shape = pol.Shape
	}
	if shape == "auto" || shape == "" {
		shape = "linear"
	}

	accent := strings.ToLower(strings.TrimSpace(g.Accent))
	if !ValidAccent(accent) || accent == "" {
		accent = pol.Accent
	}

	return Gauge{
		Label:    label,
		Value:    value,
		Min:      min,
		Max:      max,
		Percent:  percent,
		Shape:    shape,
		Accent:   accent,
		Status:   strings.TrimSpace(g.Status),
		RawValue: strconv.FormatFloat(g.Value, 'f', -1, 64),
	}
}
