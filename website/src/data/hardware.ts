export interface SupportedControllerFamily {
    family: string;
    controllers: string[];
    note: string;
}

export const supportedControllerFamilies: SupportedControllerFamily[] = [
    {
        family: "Color LCD",
        controllers: ["GC9A01A", "ILI9341", "ILI9486", "ST7735S", "ST7789"],
        note: "SPI TFT panels, from compact 128x128 displays to 480x320 screens.",
    },
    {
        family: "OLED",
        controllers: ["SH1106", "SSD1305", "SSD1306", "SSD1325", "SSD1351"],
        note: "Monochrome, grayscale, and full-color OLEDs over SPI or I2C.",
    },
    {
        family: "E-Ink",
        controllers: ["EPD4in26", "SSD1680", "SSD1681", "SSD1681 tri-color", "UC8179"],
        note: "Monochrome and tri-color e-paper, from compact panels to large dashboards.",
    },
    {
        family: "LED Matrix",
        controllers: ["HT16K33", "IS31FL3731", "IS31FL3741", "MAX7219"],
        note: "I2C and SPI LED matrices, including monochrome and RGB displays.",
    },
    {
        family: "Memory LCD",
        controllers: ["Sharp Memory LCD"],
        note: "Reflective, low-power Sharp memory displays with fast updates.",
    },
];
