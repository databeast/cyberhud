package tests

import "github.com/databeast/cyberhud/display/modes/serial"

var SerialRegistryExported = serial.Registry()

func SeedTestState(styleName string) {
	p := serial.Policy{
		Port:       "/dev/ttyUSB0",
		Baud:       115200,
		MaxLines:   24,
		AutoSelect: true,
		ScanMS:     500,
		Style:      styleName,
		Font:       "auto",
	}
	serial.SetSnapshotForTest(serial.Snapshot{
		Connected: true,
		Port:      "/dev/ttyUSB0",
		Baud:      115200,
		Lines: []string{
			"[14:30:05] sensor: temp=23.4C humidity=61%",
			"[14:30:06] motor: rpm=1200 torque=0.8Nm",
			"\033[32m[OK]\033[0m system boot complete",
			"[14:30:07] adc: ch0=2.41V ch1=1.12V",
			"[14:30:08] i2c: scan found 4 devices",
			"[14:30:09] heartbeat OK uptime=14d7h",
		},
		Sequence: 42,
	}, p)
}

func ResetTestState() { serial.ResetForTest() }
