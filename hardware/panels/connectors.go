package panels

import (
	"fmt"
	"strings"
)

func cyberhudConnectors() []connectorReport {
	return []connectorReport{
		{
			Name:        "STEMMA QT/Qwiic #1",
			Kind:        "I2C1",
			Pins:        "GND, 3V3, GPIO2/SDA, GPIO3/SCL",
			Status:      "shared",
			Description: "shared I2C bus",
		},
		{
			Name:        "STEMMA QT/Qwiic #2",
			Kind:        "I2C1",
			Pins:        "GND, 3V3, GPIO2/SDA, GPIO3/SCL",
			Status:      "shared",
			Description: "shared I2C bus",
		},
		{
			Name:   "3-pin connector GPIO13",
			Kind:   "GPIO",
			Pins:   "GND, 5V, GPIO13",
			Status: "free",
		},
		{
			Name:   "3-pin connector GPIO18",
			Kind:   "GPIO",
			Pins:   "GND, 5V, GPIO18",
			Status: "free",
		},
	}
}

func connectorConflicts(c connectorReport, assignments map[int][]string) []string {
	var out []string
	for _, pin := range connectorPins(c.Name) {
		if labels, ok := assignments[pin]; ok {
			for _, label := range labels {
				out = append(out, fmt.Sprintf("GPIO%d->%s", pin, label))
			}
		}
	}
	return out
}

func connectorPins(name string) []int {
	name = strings.ToUpper(strings.TrimSpace(name))
	pins, ok := connectorPinMap()[name]
	if !ok {
		return nil
	}
	return append([]int(nil), pins...)
}

func connectorPinMap() map[string][]int {
	return map[string][]int{
		"3-PIN CONNECTOR GPIO13": {13},
		"3-PIN CONNECTOR GPIO18": {18},
		"STEMMA QT/QWIIC #1":     {2, 3},
		"STEMMA QT/QWIIC #2":     {2, 3},
	}
}
