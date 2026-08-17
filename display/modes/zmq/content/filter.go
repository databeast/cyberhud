package content

import "encoding/json"

// filterJSON extracts specified fields from a JSON message.
// When fields is non-empty and msg is a valid JSON object, it returns
// formatted "key: value" lines (one per found field) in the order the
// fields appear in the fields slice. The value is the raw JSON encoding
// of the field's value.
//
// Returns nil when filtering does not apply:
//   - fields is empty
//   - msg is not a valid JSON object
//   - none of the specified fields exist in the message
func filterJSON(msg string, fields []string) []string {
	if len(fields) == 0 {
		return nil
	}

	// Parse msg as a JSON object mapping keys to raw JSON values.
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(msg), &obj); err != nil {
		return nil
	}

	var result []string
	for _, key := range fields {
		raw, ok := obj[key]
		if !ok {
			continue
		}
		result = append(result, key+": "+string(raw))
	}

	if len(result) == 0 {
		return nil
	}
	return result
}

// FilterJSON exposes JSON field extraction for integration tests.
func FilterJSON(msg string, fields []string) []string { return filterJSON(msg, fields) }
