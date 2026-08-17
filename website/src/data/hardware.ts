export interface SupportedPanel {
    name: string;
    controller: string;
    family: string;
    resolution: {
        width: number;
        height: number;
    };
    description: string;
}

export const supportedPanels: SupportedPanel[] = [
    // Waveshare LCD HATs
    {
        name: "Waveshare 1.3\" LCD HAT",
        controller: "ST7789",
        family: "Color LCD",
        resolution: {width: 240, height: 240},
        description: "240x240 ST7789 with buttons and joystick",
    },
    {
        name: "Waveshare 1.44\" LCD HAT",
        controller: "ST7735S",
        family: "Color LCD",
        resolution: {width: 128, height: 128},
        description: "128x128 ST7735S with three buttons and joystick",
    },
    {
        name: "Waveshare 2.2\" SPI Display",
        controller: "ST7789",
        family: "Color LCD",
        resolution: {width: 320, height: 240},
        description: "320x240 ST7789 without onboard inputs",
    },
    {
        name: "Waveshare Triple Screen HAT",
        controller: "ST7789 + ST7735S",
        family: "Color LCD",
        resolution: {width: 240, height: 240},
        description: "Main 1.3\" ST7789 (240x240) + dual 0.96\" ST7735S (160x80)",
    },
    // Waveshare OLED HATs
    {
        name: "Waveshare 1.3\" OLED HAT",
        controller: "SH1106",
        family: "Monochrome OLED",
        resolution: {width: 128, height: 64},
        description: "128x64 monochrome OLED with three buttons and D-pad",
    },
    {
        name: "Waveshare 2.23\" OLED HAT (I2C)",
        controller: "SSD1305",
        family: "Monochrome OLED",
        resolution: {width: 128, height: 32},
        description: "128x32 monochrome OLED via I2C",
    },
    {
        name: "Waveshare 2.23\" OLED HAT (SPI)",
        controller: "SSD1305",
        family: "Monochrome OLED",
        resolution: {width: 128, height: 32},
        description: "128x32 monochrome OLED via SPI",
    },
    // Generic ST7789 panels
    {
        name: "Generic ST7789 240x135",
        controller: "ST7789",
        family: "Color LCD",
        resolution: {width: 240, height: 135},
        description: "Generic SPI panel without input controls",
    },
    {
        name: "Generic ST7789 240x240",
        controller: "ST7789",
        family: "Color LCD",
        resolution: {width: 240, height: 240},
        description: "Generic SPI panel without input controls",
    },
    {
        name: "Generic ST7789 320x240",
        controller: "ST7789",
        family: "Color LCD",
        resolution: {width: 320, height: 240},
        description: "Generic SPI panel without input controls",
    },
    // Adafruit E-Ink
    {
        name: "Adafruit 2.13\" E-Ink Bonnet",
        controller: "SSD1680",
        family: "E-Ink",
        resolution: {width: 250, height: 122},
        description: "Monochrome E-Ink in landscape mode with two buttons",
    },
    {
        name: "Adafruit 2.13\" E-Ink Bonnet (2-btn)",
        controller: "SSD1680",
        family: "E-Ink",
        resolution: {width: 250, height: 122},
        description: "E-Ink with two-button menu mapping",
    },
    // Adafruit CharliePlex LED Matrices
    {
        name: "Adafruit 15x7 CharliePlex FeatherWing",
        controller: "IS31FL3731",
        family: "LED Matrix",
        resolution: {width: 15, height: 7},
        description: "LED matrix via I2C",
    },
    {
        name: "Adafruit 8x16 CharliePlex Bonnet",
        controller: "IS31FL3731",
        family: "LED Matrix",
        resolution: {width: 16, height: 8},
        description: "Green LED matrix bonnet via I2C",
    },
];
