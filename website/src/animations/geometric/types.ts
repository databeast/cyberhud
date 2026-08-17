/*
 * Geometric Background Animation — Type definitions and constants.
 *
 * Defines all interfaces for cluster generation, fade cycle computation,
 * fragment scheduling, performance monitoring, and animation state management.
 * Also defines color constants for the Necron tomb aesthetic (deep emerald to
 * bright green accent range) and the pseudocode fragment pool.
 *
 */

// ---------------------------------------------------------------------------
// Core color type
// ---------------------------------------------------------------------------

export interface HSLColor {
    /* Hue 0-360 */
    h: number;
    /* Saturation 0-100 */
    s: number;
    /* Lightness 0-100 */
    l: number;
}

// ---------------------------------------------------------------------------
// Square and Cluster interfaces
// ---------------------------------------------------------------------------

export interface SquareConfig {
    /* Offset from cluster center X (px) */
    offsetX: number;
    /* Offset from cluster center Y (px) */
    offsetY: number;
    /* Base size in CSS pixels [20, 120] — used as height; width = size * aspect */
    size: number;
    /* Aspect ratio (width / height). 1.0 = square, >1 = landscape, <1 = portrait */
    aspect: number;
    /* Rotation angle in degrees (0 or 45) */
    rotation: number;
    /* HSL color within green accent range */
    color: HSLColor;
    /* Phase offset in seconds [0, cycleDuration] */
    phaseOffset: number;
    /* Fade cycle duration in seconds [3, 10] */
    cycleDuration: number;
    /* Peak base opacity [0.05, 0.6] */
    peakOpacity: number;
}

export interface ClusterConfig {
    /* Cluster center X position (% of viewport width) */
    centerXPct: number;
    /* Cluster center Y position (% of viewport height) */
    centerYPct: number;
    /* Constituent squares (3-8) */
    squares: SquareConfig[];
    /* Bounding radius: distance from center to farthest square center (px) */
    boundingRadius: number;
    /* Spawn time: animation time in seconds when cluster first appears */
    spawnTime: number;
    /* Fade-in duration [1, 3] seconds */
    fadeInDuration: number;
}

// ---------------------------------------------------------------------------
// Pseudocode fragment interfaces
// ---------------------------------------------------------------------------

export interface ActiveFragment {
    /* The pseudocode text snippet */
    text: string;
    /* X position on canvas (px) */
    x: number;
    /* Y position on canvas (px) */
    y: number;
    /* Animation time when fragment was spawned (seconds) */
    startTime: number;
    /* Duration of fade-in phase [1, 2] seconds */
    fadeInDuration: number;
    /* Duration of hold phase [2, 5] seconds */
    holdDuration: number;
    /* Duration of fade-out phase [1, 2] seconds */
    fadeOutDuration: number;
    /* Font size [10, 16] px */
    fontSize: number;
    /* Text color (green or cyan range) */
    color: HSLColor;
    /* Peak opacity [0.3, 0.8] */
    peakOpacity: number;
}

export interface FragmentSchedulerState {
    /* Currently active (visible or fading) fragments */
    activeFragments: ActiveFragment[];
    /* Animation time of last fragment spawn (seconds) */
    lastSpawnTime: number;
    /* Text of the last spawned fragment (to avoid consecutive repeats) */
    lastSpawnedText: string | null;
}

// ---------------------------------------------------------------------------
// Performance monitoring
// ---------------------------------------------------------------------------

export interface PerfState {
    /* Rolling window of frame times (ms), max 10 entries */
    frameTimes: number[];
    /* Current number of squares being rendered */
    currentSquareCount: number;
    /* Original square count before any reductions */
    originalSquareCount: number;
    /* Whether performance-based reduction has occurred */
    hasReduced: boolean;
}

// ---------------------------------------------------------------------------
// Top-level animation state
// ---------------------------------------------------------------------------

export interface AnimationState {
    /* Whether the animation loop is running */
    running: boolean;
    /* Whether reduced motion is active */
    reducedMotion: boolean;
    /* Canvas width (px) */
    width: number;
    /* Canvas height (px) */
    height: number;
    /* Current animation time (seconds) */
    time: number;
    /* Active cluster configurations */
    clusters: ClusterConfig[];
    /* Fragment scheduler state */
    fragmentState: FragmentSchedulerState;
    /* Performance monitoring state */
    perfState: PerfState;
    /* Total frames rendered since start */
    frameCount: number;
    /* Current requestAnimationFrame ID, or null if stopped */
    animationFrameId: number | null;
}

// ---------------------------------------------------------------------------
// Color constants — Square accent range (Requirement 1.1, 6.1)
// ---------------------------------------------------------------------------

/* Minimum hue for square colors (deep emerald) */
export const SQUARE_HUE_MIN = 100;
/* Maximum hue for square colors (bright green) */
export const SQUARE_HUE_MAX = 140;
/* Minimum saturation for square colors (%) */
export const SQUARE_SAT_MIN = 70;
/* Minimum lightness for square colors (%) */
export const SQUARE_LIGHTNESS_MIN = 20;
/* Maximum lightness for square colors (%) */
export const SQUARE_LIGHTNESS_MAX = 55;

// ---------------------------------------------------------------------------
// Color constants — Fragment green range (Requirement 3.5, 6.2)
// ---------------------------------------------------------------------------

/* Hue range for green pseudocode fragments */
export const FRAGMENT_GREEN_HUE_RANGE = [100, 140] as const;
/* Minimum saturation for green fragments (%) */
export const FRAGMENT_GREEN_SAT_MIN = 70;
/* Lightness range for green fragments (%) */
export const FRAGMENT_GREEN_LIGHTNESS = [40, 60] as const;

// ---------------------------------------------------------------------------
// Color constants — Fragment cyan range (Requirement 3.5, 6.2)
// ---------------------------------------------------------------------------

/* Hue range for cyan pseudocode fragments */
export const FRAGMENT_CYAN_HUE_RANGE = [170, 190] as const;
/* Minimum saturation for cyan fragments (%) */
export const FRAGMENT_CYAN_SAT_MIN = 80;
/* Lightness range for cyan fragments (%) */
export const FRAGMENT_CYAN_LIGHTNESS = [45, 65] as const;

// ---------------------------------------------------------------------------
// Pseudocode fragment pool (Requirement 3.3)
// ---------------------------------------------------------------------------

/*
 * Pool of pseudocode text snippets — cyberpunk Pi/electronics/monitoring aesthetic.
 * Themed around system monitoring, radio, home automation, and hardware hacking.
 * Each entry is at most 40 characters.
 */
export const PSEUDOCODE_POOL: string[] = [
    "cpu.temp = read_sensor(0x4C);",
    "if (gpio.pin(17).HIGH) {",
    "let freq = sdr.tune(433.92);",
    "sys.uptime = 847291;",
    "while (daemon.alive) {",
    "voltage = adc.read(CH0);",
    "i2c.write(0x3C, framebuf);",
    "mqtt.publish('/home/status');",
    "return core_temp_c;",
    "for (sensor in bus) {",
    "status: NOMINAL",
    "fan.speed = pwm(0.6);",
    "spi.transfer(display_cmd);",
    "mem.free = 412MB",
    "radio.scan(87.5, 108.0);",
    "uart.baud(115200);",
    "> signal acquired",
    "light.scene('evening');",
    "cron.schedule('*/5 * * * *');",
    "thermostat.set(21.5);",
    "ping(gateway) = 2ms",
    "docker.restart('hud-core');",
    "load_avg: 0.42 0.38 0.31",
    "antenna.gain = 12dBi;",
    "zigbee.pair(device_new);",
    "disk.usage = 67%;",
    "waveform.sample(44100);",
    "$ systemctl status cyberhud",
    "relay.toggle(CH2);",
    "spectrum[i] = fft(buf);",
    "cam.snapshot('/tmp/frame');",
    "dht22.humidity = 58%;",
    "net.latency < threshold",
    "oled.draw_line(0,0,128,64);",
    "pwr.consumption = 4.2W;",
    "log('boot sequence complete');",
    "gps.fix = { lat, lon, alt };",
    "servo.angle(90);",
    "watchdog.feed();",
    "rf.tx(packet, 915MHz);",
    "led.strip.fill(0x00FF41);",
    "$ journalctl -f -u cyberhud",
    "adc.calibrate(vref=3.3);",
    "ha.state('binary_sensor.door');",
    "sdcard.mount('/data');",
    "ir.send(NEC, 0xA25D);",
    "bme280.pressure = 1013hPa;",
    "> all systems operational",
];

// ---------------------------------------------------------------------------
// Layout constants (Requirement 1.2)
// ---------------------------------------------------------------------------

/* Minimum number of clusters on desktop viewports (width >= 1024px) */
export const MIN_DESKTOP_CLUSTERS = 5;
