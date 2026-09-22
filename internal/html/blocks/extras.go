package blocks

import (
	"encoding/json"
	"sort"
	"strconv"
)

// ExtraField is one key a capybara config carries that grig does not
// model, shown in the service modal so nothing is hidden from the user.
type ExtraField struct {
	Key   string
	Value string
	// Type is "string", "number" or "bool", derived from the Go value
	// so the template can render the right input widget.
	Type string
}

// scalarText renders a scalar extra for a text input, or reports false
// for values that cannot round-trip through one. It also returns the
// type name so the UI can pick the right widget.
func scalarText(value any) (text string, typeName string, ok bool) {
	switch v := value.(type) {
	case string:
		return v, "string", true
	case bool:
		return strconv.FormatBool(v), "bool", true
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), "number", true
	case int:
		return strconv.Itoa(v), "number", true
	default:
		return "", "", false
	}
}

// ScalarExtras are the unmodelled keys editable as plain text, sorted so
// the field order is stable between renders.
func ScalarExtras(extra map[string]any) []ExtraField {
	fields := make([]ExtraField, 0, len(extra))
	for key, value := range extra {
		if text, typeName, ok := scalarText(value); ok {
			fields = append(fields, ExtraField{Key: key, Value: text, Type: typeName})
		}
	}
	sort.Slice(fields, func(i, j int) bool { return fields[i].Key < fields[j].Key })
	return fields
}


// ComplexExtras are the unmodelled keys whose values are objects or
// lists. They are shown read-only and written back untouched, rather
// than flattened into a text input.
func ComplexExtras(extra map[string]any) []ExtraField {
	fields := make([]ExtraField, 0)
	for key, value := range extra {
		if _, _, ok := scalarText(value); ok {
			continue
		}
		encoded, err := json.Marshal(value)
		if err != nil {
			continue
		}
		fields = append(fields, ExtraField{Key: key, Value: string(encoded)})
	}
	sort.Slice(fields, func(i, j int) bool { return fields[i].Key < fields[j].Key })
	return fields
}
