export interface VersionEntry {
    version: string;
    date: string;
    summary: string;
}

/*
 * Version history shown on the public marketing site.
 * Add newest entries to the top of the list.
 */
export const versionHistory: VersionEntry[] = [
    {
        version: "0.1.0",
        date: "August 2026",
        summary: "First public release of CyberHUD.",
    },
];
