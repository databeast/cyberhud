package main

// SplitCommands splits a CLI argument list by standalone semicolons into
// individual commands. Only standalone ";" args are treated as separators.
// No quoted-string multi-command form is supported.
func SplitCommands(args []string) [][]string {
	if len(args) == 0 {
		return nil
	}

	var result [][]string
	var current []string

	for _, arg := range args {
		if arg == ";" {
			if len(current) > 0 {
				result = append(result, current)
				current = nil
			}
		} else {
			current = append(current, arg)
		}
	}

	if len(current) > 0 {
		result = append(result, current)
	}

	return result
}
