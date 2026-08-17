package source

import (
	"sync"
	"time"
)

// knownDeviceName returns a human-readable name for well-known STEMMA QT /
// QWIIC I2C device addresses.  Returns "unknown" for unrecognised addresses.
func knownDeviceName(addr uint16) string {
	if name, ok := knownDevices[addr]; ok {
		return name
	}
	return "unknown"
}

// knownDevices maps I2C 7-bit addresses to product names for common devices
// in the Adafruit STEMMA QT and SparkFun QWIIC ecosystems.
var knownDevices = map[uint16]string{
	// ── Adafruit ─────────────────────────────────────────────────────────────
	0x18: "LIS3DH accelerometer",
	0x19: "LIS3DH accelerometer (alt)",
	0x1C: "MMA8451 accelerometer",
	0x1D: "MMA8451 accelerometer (alt)",
	0x20: "MCP23017 I/O expander",
	0x21: "MCP23017 I/O expander (alt)",
	0x23: "BH1750 light sensor",
	0x28: "BNO055 IMU",
	0x29: "TSL2561/VL53L0X/APDS9960",
	0x38: "VEML6070 UV sensor",
	0x39: "APDS9960/TSL2561 light sensor",
	0x3C: "SSD1306 OLED display (128×64)",
	0x3D: "SSD1306 OLED display (alt)",
	0x40: "INA219 power monitor",
	0x41: "INA219 power monitor (alt)",
	0x44: "SHT30/SHT31 humidity sensor",
	0x45: "SHT30/SHT31 humidity sensor (alt)",
	0x48: "ADS1015/ADS1115 ADC",
	0x49: "ADS1015/ADS1115 ADC (alt)",
	0x4A: "ADS1015/ADS1115 ADC (alt2)",
	0x4B: "ADS1015/ADS1115 ADC (alt3)",
	0x60: "Si7021 humidity/temp / VCNL4010",
	0x68: "MPU-6050 IMU / DS3231 RTC",
	0x69: "MPU-6050 IMU (alt)",
	0x70: "HT16K33 LED backpack",
	0x71: "HT16K33 LED backpack (alt1)",
	0x72: "HT16K33 LED backpack (alt2) / Qwiic Twist",
	0x73: "HT16K33 LED backpack (alt3)",
	0x74: "HT16K33 LED backpack (alt4)",
	0x75: "HT16K33 LED backpack (alt5)",
	0x76: "BME280/BMP280 env sensor",
	0x77: "BME280/BMP280 env sensor (alt)",

	// ── SparkFun QWIIC ───────────────────────────────────────────────────────
	0x10: "VCNL4040 proximity/light",
	0x52: "APDS9960 gesture/light (SparkFun)",
	0x5B: "CCS811 air quality",
	0x5D: "ZX distance sensor",
	0x6F: "Qwiic Keypad",

	// ── Misc ─────────────────────────────────────────────────────────────────
	0x50: "24CXX EEPROM",
	0x57: "MAX17048 fuel gauge",
}

// scannerConfig holds the package-level scanner configuration provided by
// daemon init (cmd/cyberhudd) via SetScannerConfig. These values are copied
// into each instance at construction time so that Activate() can create a
// scanner with the correct parameters.
var scannerConfig struct {
	sync.RWMutex
	buses    []string
	interval time.Duration
}

// SetScannerConfig stores I2C bus names and scan interval at the package level.
// Call this from daemon init (cmd/cyberhudd) before any mode instances are
// constructed, so that Activate() can create a scanner with correct parameters.
func SetScannerConfig(buses []string, interval time.Duration) {
	scannerConfig.Lock()
	defer scannerConfig.Unlock()
	scannerConfig.buses = buses
	scannerConfig.interval = interval
}

// getScannerConfig returns the current package-level scanner configuration.
func GetScannerConfig() ([]string, time.Duration) {
	scannerConfig.RLock()
	defer scannerConfig.RUnlock()
	return scannerConfig.buses, scannerConfig.interval
}
