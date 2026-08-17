export interface SiteConfig {
    projectName: string;
    tagline: string;
    metaDescription: string;
    ogImage: string;
    links: {
        install: string;
        github: string;
        docs: string;
    };
}

export const siteConfig: SiteConfig = {
    projectName: "CyberHUD",
    tagline: "The coolest\n" +
        "Heads-Up Display system\n" +
        "for Small Screens\n" +
        "and Cyberdecks",
    metaDescription:
        "CyberHUD turns SPI/I2C displays from Adafruit, Waveshare, Elecrow and more into working screens with no code — dashboards, monitors, and more on Linux GPIO.",
    ogImage: "/images/og-preview.png",
    links: {
        install: "https://databeast.github.io/cyberhud/installation/",
        github: "https://github.com/databeast/cyberhud",
        docs: "https://databeast.github.io/cyberhud",
    },
};
