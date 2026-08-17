export interface Feature {
    id: string;
    title: string;
    description: string;
    mediaSrc: string;
    mediaAlt: string;
    mediaWidth: number;
    mediaHeight: number;
}

export const features: Feature[] = [
    {
        id: "multi-display",
        title: "Multi-Display, Many Modes",
        description:
            "Run dashboard, GPIO/I2C monitors, system stats, clock, and more — independently on up to 3 panels.",
        mediaSrc: "/images/features/multi-display.svg",
        mediaAlt: "Three SPI displays showing different HUD content driven by one Raspberry Pi",
        mediaWidth: 600,
        mediaHeight: 400,
    },
    {
        id: "attract-modes",
        title: "Ambient Attract Modes",
        description:
            "Animated bokeh, particles, plasma, and starfield modes for passive displays with no input needed.",
        mediaSrc: "/images/features/attract-modes.svg",
        mediaAlt: "Animated bokeh attract mode running on a small SPI display panel",
        mediaWidth: 600,
        mediaHeight: 400,
    },
    {
        id: "runtime-control",
        title: "Live Runtime Control",
        description:
            "Switch modes, push text, check status, and freeze the config to disk — all via cyberhudctl or a socket.",
        mediaSrc: "/images/features/runtime-control.svg",
        mediaAlt: "Terminal showing cyberhudctl commands controlling display modes in real time",
        mediaWidth: 600,
        mediaHeight: 400,
    },
    {
        id: "hardware-compatibility",
        title: "No-Code Hardware Support",
        description:
            "Adafruit and Waveshare panels do nothing until coded. CyberHUD drives any panel on a supported chipset out of the box.",
        mediaSrc: "/images/features/hardware-compat.svg",
        mediaAlt: "Raspberry Pi connected to multiple SPI display panels of different sizes",
        mediaWidth: 600,
        mediaHeight: 400,
    },
];
