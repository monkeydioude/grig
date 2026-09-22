package model

import (
	"encoding/json"
	"reflect"
	"strings"
)

// Capybara config files carry keys grig does not model. Editing a config
// here must never drop them, so each type keeps the unknown keys it was
// loaded with and writes them back untouched.
//
// YAML handles this with a `yaml:",inline"` map. JSON has no equivalent,
// so the helpers below do the same job for the JSON codec.

// jsonKeys lists the JSON object keys a struct declares, which is how an
// unknown key is told apart from one the model owns.
func jsonKeys(v any) map[string]struct{} {
	t := reflect.TypeOf(v)
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	keys := make(map[string]struct{}, t.NumField())
	for i := range t.NumField() {
		field := t.Field(i)
		name, _, _ := strings.Cut(field.Tag.Get("json"), ",")
		if name == "-" {
			continue
		}
		if name == "" {
			name = field.Name
		}
		keys[name] = struct{}{}
	}
	return keys
}

// decodeWithExtras fills v from data and returns the keys v does not
// declare. v must be a pointer to a type without its own UnmarshalJSON,
// so callers pass a local alias of their struct.
func decodeWithExtras(data []byte, v any) (map[string]any, error) {
	if err := json.Unmarshal(data, v); err != nil {
		return nil, err
	}
	all := map[string]any{}
	if err := json.Unmarshal(data, &all); err != nil {
		return nil, err
	}
	for key := range jsonKeys(v) {
		delete(all, key)
	}
	if len(all) == 0 {
		return nil, nil
	}
	return all, nil
}

// encodeWithExtras encodes v and splices extras back in. Keys the model
// declares always win, so clearing a field in the UI is not undone by a
// stale extra. The check is against the declared keys rather than the
// encoded output, since an omitempty field disappears from the output
// exactly when it has been cleared.
func encodeWithExtras(v any, extras map[string]any) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	if len(extras) == 0 {
		return data, nil
	}
	merged := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &merged); err != nil {
		return nil, err
	}
	declared := jsonKeys(v)
	for key, value := range extras {
		if _, owned := declared[key]; owned {
			continue
		}
		raw, err := json.Marshal(value)
		if err != nil {
			return nil, err
		}
		merged[key] = raw
	}
	return json.Marshal(merged)
}
