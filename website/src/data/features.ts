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
            "Run dashboards, GPIO/I2C monitors, system stats, clocks, and more — independently across multiple panels.",
        mediaSrc: "/images/features/multi-display.svg",
        mediaAlt: "Three SPI displays showing different HUD content driven by one Raspberry Pi",
        mediaWidth: 600,
        mediaHeight: 400,
    },
    {
        id: "attract-modes",
        title: "Ambient Attract Modes",
        description:
            "A rotating cube is one of many attract modes—including bokeh, particles, plasma, and starfields—for passive displays.",
        mediaSrc: "/images/features/attract-modes.svg",
        mediaAlt: "A glowing cyan and magenta wireframe cube spinning on a dark display",
        mediaWidth: 600,
        mediaHeight: 400,
    },
    {
        id: "runtime-control",
        title: "Live Runtime Control",
        description:
            "Switch modes, push text, check status, and freeze the config to disk — via cyberhudctl, the web UI, or a socket.",
        mediaSrc: "/images/features/runtime-control.svg",
        mediaAlt: "Terminal showing cyberhudctl commands controlling display modes in real time",
        mediaWidth: 600,
        mediaHeight: 400,
    },
    {
        id: "hardware-compatibility",
        title: "Broad Display Support",
        description:
            "Ready-made profiles for common panels; generic configs support compatible displays across 20 controller drivers.",
        mediaSrc: "/images/features/hardware-compat.svg",
        mediaAlt: "Raspberry Pi connected to multiple SPI display panels of different sizes",
        mediaWidth: 600,
        mediaHeight: 400,
    },
];
