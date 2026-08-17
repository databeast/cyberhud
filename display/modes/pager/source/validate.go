package source

import "errors"

// errSourceNotAbsolute is the error returned when a source path does not begin
// with "/".
var errSourceNotAbsolute = errors.New("source path must be absolute (begin with /)")

// validateSource checks that the given path is a valid absolute source path.
// It accepts only non-empty strings beginning with "/". Empty strings and
// non-absolute paths are rejected with a descriptive error.
func ValidateSource(path string) error {
	if len(path) == 0 || path[0] != '/' {
		return errSourceNotAbsolute
	}
	return nil
}
