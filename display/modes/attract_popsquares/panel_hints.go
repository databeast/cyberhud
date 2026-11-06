package attract_popsquares

// fitnessNotesPostApply is the PostApply hook for the attract_popsquares
// CmdHandler. This mode auto-resolves styles from panel hints; no policy key
// affects style selection, so no fitness notes are generated.
func fitnessNotesPostApply(appliedKeys []string) []string { return nil }
