package panels

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Canonical GPIO pin-name constants used by panel definitions.
const (
	GPIO2  = "GPIO2"
	GPIO3  = "GPIO3"
	GPIO4  = "GPIO4"
	GPIO5  = "GPIO5"
	GPIO6  = "GPIO6"
	GPIO7  = "GPIO7"
	GPIO8  = "GPIO8"
	GPIO9  = "GPIO9"
	GPIO10 = "GPIO10"
	GPIO11 = "GPIO11"
	GPIO12 = "GPIO12"
	GPIO13 = "GPIO13"
	GPIO16 = "GPIO16"
	GPIO17 = "GPIO17"
	GPIO18 = "GPIO18"
	GPIO19 = "GPIO19"
	GPIO20 = "GPIO20"
	GPIO21 = "GPIO21"
	GPIO22 = "GPIO22"
	GPIO23 = "GPIO23"
	GPIO24 = "GPIO24"
	GPIO25 = "GPIO25"
	GPIO26 = "GPIO26"
	GPIO27 = "GPIO27"
)

type pinAssignment struct {
	Label string
	GPIO  string
}

type connectorReport struct {
	Name        string
	Kind        string
	Pins        string
	Status      string
	Conflicts   []string
	Description string
}

// BuildPinReport renders a connector/pin conflict report for one panel.
func BuildPinReport(def Definition) string {
	assignments := panelPinAssignments(def)
	assignmentByNum := map[int][]string{}
	for _, a := range assignments {
		if n, ok := gpioNumber(a.GPIO); ok {
			assignmentByNum[n] = append(assignmentByNum[n], a.Label)
		}
	}

	connectors := cyberhudConnectors()
	for i := range connectors {
		connectors[i].Conflicts = connectorConflicts(connectors[i], assignmentByNum)
		if len(connectors[i].Conflicts) > 0 {
			connectors[i].Status = "conflict"
		}
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "OK pin report\n")
	fmt.Fprintf(&sb, "  panel=%s controller=%s\n", def.Name, def.Controller)
	if def.Monochrome {
		sb.WriteString("  display_mode=monochrome\n")
	}
	if len(assignments) == 0 {
		sb.WriteString("  display_pins=(none)\n")
	} else {
		sb.WriteString("  display_pins:\n")
		sort.Slice(assignments, func(i, j int) bool {
			iNum, iOK := gpioNumber(assignments[i].GPIO)
			jNum, jOK := gpioNumber(assignments[j].GPIO)
			if iOK && jOK && iNum != jNum {
				return iNum < jNum
			}
			return assignments[i].Label < assignments[j].Label
		})
		for _, a := range assignments {
			fmt.Fprintf(&sb, "    %-12s %s\n", a.Label+":", a.GPIO)
		}
	}
	if len(connectors) > 0 {
		sb.WriteString("  connectors:\n")
		for _, c := range connectors {
			fmt.Fprintf(&sb, "    %-24s kind=%s pins=%s status=%s", c.Name+":", c.Kind, c.Pins, c.Status)
			if c.Description != "" {
				fmt.Fprintf(&sb, " (%s)", c.Description)
			}
			if len(c.Conflicts) > 0 {
				fmt.Fprintf(&sb, " conflict_with=%s", strings.Join(c.Conflicts, ","))
			}
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

// PinNotices returns informational messages about cyberhud GPIO output
// connectors that are unavailable because the active display panel is using
// those pins.
func PinNotices(def Definition) []string {
	assignments := panelPinAssignments(def)
	assignmentByNum := map[int][]string{}
	for _, a := range assignments {
		if n, ok := gpioNumber(a.GPIO); ok {
			assignmentByNum[n] = append(assignmentByNum[n], a.Label)
		}
	}
	var out []string
	for _, c := range cyberhudConnectors() {
		if c.Kind != "GPIO" {
			continue
		}
		conflicts := connectorConflicts(c, assignmentByNum)
		if len(conflicts) == 0 {
			continue
		}
		out = append(out, fmt.Sprintf("cyberhud output %s unavailable (GPIO claimed by %s)", c.Name, strings.Join(conflicts, ",")))
	}
	return out
}

func panelPinAssignments(def Definition) []pinAssignment {
	assignments := []pinAssignment{}
	add := func(label, gpioName string) {
		gpioName = strings.TrimSpace(gpioName)
		if gpioName == "" {
			return
		}
		assignments = append(assignments, pinAssignment{Label: label, GPIO: gpioName})
	}

	if def.DCPin != "" {
		add("display_dc", def.DCPin)
	}
	if def.RSTPin != "" {
		add("display_rst", def.RSTPin)
	}
	if def.BLPin != "" {
		add("display_bl", def.BLPin)
	}
	if def.BusyPin != "" {
		add("display_busy", def.BusyPin)
	}
	for _, vp := range def.Virtual {
		tag := fmt.Sprintf("display%d", vp.Index)
		for _, sp := range spiPinsForDevice(vp.SPI) {
			add(tag+"_"+sp.Label, sp.GPIO)
		}
		if vp.DCPin != "" {
			add(tag+"_dc", vp.DCPin)
		}
		if vp.RSTPin != "" {
			add(tag+"_rst", vp.RSTPin)
		}
		if vp.BLPin != "" {
			add(tag+"_bl", vp.BLPin)
		}
		if vp.BusyPin != "" {
			add(tag+"_busy", vp.BusyPin)
		}
	}

	add("input_key1", def.Inputs.Key1)
	add("input_key2", def.Inputs.Key2)
	add("input_key3", def.Inputs.Key3)
	add("input_up", def.Inputs.JoyUp)
	add("input_down", def.Inputs.JoyDown)
	add("input_left", def.Inputs.JoyLeft)
	add("input_right", def.Inputs.JoyRight)
	add("input_press", def.Inputs.JoyPressed)
	return assignments
}

func spiPinsForDevice(name string) []pinAssignment {
	name = strings.ToUpper(strings.TrimSpace(name))
	switch name {
	case "SPI0.0":
		return []pinAssignment{{Label: "mosi", GPIO: GPIO10}, {Label: "sclk", GPIO: GPIO11}, {Label: "cs", GPIO: GPIO8}}
	case "SPI0.1":
		return []pinAssignment{{Label: "mosi", GPIO: GPIO10}, {Label: "sclk", GPIO: GPIO11}, {Label: "cs", GPIO: GPIO7}}
	case "SPI1.0":
		return []pinAssignment{{Label: "mosi", GPIO: GPIO20}, {Label: "sclk", GPIO: GPIO21}, {Label: "cs", GPIO: GPIO18}}
	case "SPI1.1":
		return []pinAssignment{{Label: "mosi", GPIO: GPIO20}, {Label: "sclk", GPIO: GPIO21}, {Label: "cs", GPIO: GPIO17}}
	default:
		return nil
	}
}

func gpioNumber(name string) (int, bool) {
	name = strings.ToUpper(strings.TrimSpace(name))
	if !strings.HasPrefix(name, "GPIO") {
		return 0, false
	}
	if n, err := strconv.Atoi(strings.TrimPrefix(name, "GPIO")); err == nil {
		return n, true
	}
	return 0, false
}
