package config

import (
	"fmt"
	"path/filepath"
	"strings"
)

// CapybaraRef identifies one capybara config: its index in the server
// config list, the name shown in the UI and the file it was loaded from.
type CapybaraRef struct {
	ID   int
	Name string
	Path string
}

// URL is the page editing this config.
func (r CapybaraRef) URL() string {
	return fmt.Sprintf("/capybara/%d", r.ID)
}

// APIURL is the endpoint this config is saved to.
func (r CapybaraRef) APIURL() string {
	return fmt.Sprintf("/api/v1/capybara/%d", r.ID)
}

// capybaraName derives a display name from a config file name:
// "garden.capybara.yaml" -> "garden", "capybara.config.json" -> "capybara".
func capybaraName(path string) string {
	name := filepath.Base(path)
	name = strings.TrimSuffix(name, filepath.Ext(name))
	name = strings.TrimSuffix(name, ".config")
	name = strings.TrimSuffix(name, ".capybara")
	if name == "" || name == "." || name == string(filepath.Separator) {
		return "capybara"
	}
	return name
}

// CapybaraRefs lists every configured capybara, named after its file.
// Names that would collide keep their parent directory as a prefix.
func (sc ServerConfig) CapybaraRefs() []CapybaraRef {
	count := make(map[string]int, len(sc.CapybaraConfigPaths))
	for _, path := range sc.CapybaraConfigPaths {
		count[capybaraName(path)]++
	}
	refs := make([]CapybaraRef, 0, len(sc.CapybaraConfigPaths))
	for id, path := range sc.CapybaraConfigPaths {
		name := capybaraName(path)
		if count[name] > 1 {
			name = filepath.Join(filepath.Base(filepath.Dir(path)), name)
		}
		refs = append(refs, CapybaraRef{ID: id, Name: name, Path: path})
	}
	return refs
}

// CapybaraRefByID resolves the config an URL's {id} refers to.
func (sc ServerConfig) CapybaraRefByID(id int) (CapybaraRef, bool) {
	if id < 0 || id >= len(sc.CapybaraConfigPaths) {
		return CapybaraRef{}, false
	}
	return sc.CapybaraRefs()[id], true
}
