package model

import (
	"fmt"
	"log/slog"
	customErr "monkeydioude/grig/internal/errors"
	"monkeydioude/grig/internal/service/file"
	"monkeydioude/grig/pkg/errors"
	"monkeydioude/grig/pkg/trans_types"
	"os"
	"strconv"
	"strings"
)

type Proxy struct {
	Port trans_types.StringInt `json:"port" yaml:"port"`
	// TLSHosts are the domains the proxy terminates TLS for. Optional:
	// a proxy serving plain HTTP has none.
	TLSHosts []string `json:"tls_hosts,omitempty" yaml:"tls_hosts,omitempty"`
	// Email is the contact address used for certificate issuance.
	Email string `json:"email,omitempty" yaml:"email,omitempty"`
	// Extra holds proxy keys grig does not model, kept so a save here
	// never drops them.
	Extra map[string]any `json:"-" yaml:",inline"`
}

func (p *Proxy) UnmarshalJSON(data []byte) error {
	type alias Proxy
	extras, err := decodeWithExtras(data, (*alias)(p))
	if err != nil {
		return err
	}
	p.Extra = extras
	return nil
}

func (p Proxy) MarshalJSON() ([]byte, error) {
	type alias Proxy
	return encodeWithExtras(alias(p), p.Extra)
}

// Verify only requires a port. TLS hosts and the contact email are
// optional, since a proxy may serve plain HTTP.
func (p Proxy) Verify() error {
	if p.Port <= 0 {
		return fmt.Errorf("Proxy.Verify(): Port: %d: %w", p.Port, customErr.ErrModelVerifyInvalidValue)
	}
	return nil
}

// sanitizeTLSHosts drops the blank entries an edit form submits for
// rows the user left empty.
func (p *Proxy) sanitizeTLSHosts() {
	hosts := make([]string, 0, len(p.TLSHosts))
	for _, host := range p.TLSHosts {
		if trimmed := strings.TrimSpace(host); trimmed != "" {
			hosts = append(hosts, trimmed)
		}
	}
	p.TLSHosts = hosts
}

type ServiceDefinition struct {
	ID       string                `json:"id" yaml:"id"`
	Method   string                `json:"method" yaml:"method"`
	Pattern  string                `json:"pattern" yaml:"pattern"`
	Port     trans_types.StringInt `json:"port" yaml:"port"`
	Protocol string                `json:"protocol,omitempty" yaml:"protocol,omitempty"`
	// Host routes by virtual host. Several services may share a pattern
	// and be told apart only by this, so it must survive a save.
	Host string `json:"host,omitempty" yaml:"host,omitempty"`
	// Extra holds service keys grig does not model, such as "redirect".
	Extra map[string]any `json:"-" yaml:",inline"`
}

func (sd *ServiceDefinition) UnmarshalJSON(data []byte) error {
	type alias ServiceDefinition
	extras, err := decodeWithExtras(data, (*alias)(sd))
	if err != nil {
		return err
	}
	sd.Extra = extras
	return nil
}

func (sd ServiceDefinition) MarshalJSON() ([]byte, error) {
	type alias ServiceDefinition
	return encodeWithExtras(alias(sd), sd.Extra)
}

// isBlank reports whether this is an untouched row from the edit form
// rather than a real service.
func (sd ServiceDefinition) isBlank() bool {
	return sd.ID == "" && sd.Pattern == "" && sd.Host == "" && sd.Port == 0
}

func wrapErr(method, id string) error {
	return errors.Wrapf(customErr.ErrModelVerifyInvalidValue, "ServiceDefinition.%s(): ID: %q", method, id)
}

func (sd ServiceDefinition) Verify() error {
	if sd.ID == "" {
		return wrapErr("Verify", sd.ID)
	}
	if sd.Method == "" {
		return wrapErr("Method", sd.ID)
	}
	if sd.Pattern == "" {
		return wrapErr("Pattern", sd.ID)
	}
	if sd.Port <= 0 {
		return wrapErr("Port", sd.ID)
	}
	return nil
}

func (sd ServiceDefinition) PortString() string {
	if sd.Port == 0 {
		return ""
	}
	return strconv.Itoa(int(sd.Port))
}

type Capybara struct {
	Proxy      Proxy                                   `json:"proxy" yaml:"proxy"`
	Services   []ServiceDefinition                     `json:"services" yaml:"services"`
	Path       string                                  `json:"-" yaml:"-"`
	FileWriter func(string, []byte, os.FileMode) error `json:"-" yaml:"-"`
	// Extra holds top-level keys grig does not model.
	Extra map[string]any `json:"-" yaml:",inline"`
}

func (c *Capybara) UnmarshalJSON(data []byte) error {
	type alias Capybara
	extras, err := decodeWithExtras(data, (*alias)(c))
	if err != nil {
		return err
	}
	c.Extra = extras
	return nil
}

func (c Capybara) MarshalJSON() ([]byte, error) {
	type alias Capybara
	return encodeWithExtras(alias(c), c.Extra)
}

func (c Capybara) Save() error {
	// the on-disk format follows the config file extension (yaml or json)
	data, err := file.MarshalForPath(c.Path, c)
	if err != nil {
		return errors.Wrapf(err, "Capybara.Save(): %w", customErr.ErrMarshaling)
	}
	if c.FileWriter == nil {
		return errors.Wrapf(customErr.ErrNilPointer, "Capybara.FileWriter()")
	}
	if err := c.FileWriter(c.Path, data, os.ModePerm); err != nil {
		return errors.Wrapf(err, "Capybara.Save(): %w", customErr.ErrWritingFile)
	}
	return nil
}

func (c *Capybara) Sanitize() {
	c.Proxy.sanitizeTLSHosts()
	// c.Services = slices.DeleteFunc(c.Services, func (sd ServiceDefinition) bool {
	// 	return sd.Verify() != nil
	// })
	// Only blank rows submitted by the edit form are dropped. A service
	// that is merely incomplete by grig's rules is kept: capybara decides
	// what is valid, and silently deleting a service the user never
	// touched is worse than writing it back unchanged.
	services := make([]ServiceDefinition, 0, len(c.Services))
	for _, sd := range c.Services {
		if sd.isBlank() {
			slog.Debug("Capybara.Sanitize: dropping blank service row")
			continue
		}
		services = append(services, sd)
	}
	c.Services = services
}

func (c Capybara) CloneBase() Capybara {
	return Capybara{
		Path:       c.Path,
		FileWriter: c.FileWriter,
	}
}

// MergeExtrasFrom copies the unknown keys of the config currently on disk
// onto this one. The edit form submits only the fields grig displays, so
// a save built from that payload carries no extras; without this merge it
// would rewrite the file with everything else stripped out.
func (c *Capybara) MergeExtrasFrom(existing Capybara) {
	c.Extra = mergeExtras(existing.Extra, c.Extra)
	c.Proxy.Extra = mergeExtras(existing.Proxy.Extra, c.Proxy.Extra)
	c.mergeServiceExtras(existing.Services)
}

// mergeExtras overlays the extras posted by the edit form onto the ones
// already in the file. Keys the form does not render, such as nested
// values, are kept; a key posted empty is removed; and a posted value is
// converted back to the type it had on disk, since a form posts every
// value as a string. For brand-new keys the form also posts a companion
// __type__keyname hint so grig knows whether to store a number, bool or
// string.
func mergeExtras(onDisk, posted map[string]any) map[string]any {
	if len(posted) == 0 {
		return onDisk
	}

	// collect type hints before processing values
	typeHints := make(map[string]string)
	for key, value := range posted {
		if after, ok := strings.CutPrefix(key, "__type__"); ok {
			if hint, isStr := value.(string); isStr {
				typeHints[after] = hint
			}
		}
	}

	merged := make(map[string]any, len(onDisk)+len(posted))
	for key, value := range onDisk {
		merged[key] = value
	}
	for key, value := range posted {
		// skip type hint companion fields, they are metadata only
		if strings.HasPrefix(key, "__type__") {
			continue
		}
		if text, isText := value.(string); isText && strings.TrimSpace(text) == "" {
			delete(merged, key)
			continue
		}
		merged[key] = coerceToType(value, onDisk[key], typeHints[key])
	}
	return merged
}

// coerceToType restores the type a value had on disk when the form
// posted it back as text. YAML decodes whole numbers as int and JSON as
// float64, so both are handled. For brand-new keys that have no on-disk
// original, the typeHint ("bool", "number" or "string") from the form's
// type selector is used instead.
func coerceToType(posted, original any, typeHint string) any {
	text, ok := posted.(string)
	if !ok {
		return posted
	}

	// existing key: coerce to the type already on disk
	if original != nil {
		switch original.(type) {
		case bool:
			if parsed, err := strconv.ParseBool(text); err == nil {
				return parsed
			}
		case int:
			if parsed, err := strconv.Atoi(text); err == nil {
				return parsed
			}
		case float64:
			if parsed, err := strconv.ParseFloat(text, 64); err == nil {
				return parsed
			}
		}
		return posted
	}

	// new key: use the type hint from the form
	switch typeHint {
	case "bool":
		if parsed, err := strconv.ParseBool(text); err == nil {
			return parsed
		}
	case "number":
		if parsed, err := strconv.Atoi(text); err == nil {
			return parsed
		}
		if parsed, err := strconv.ParseFloat(text, 64); err == nil {
			return parsed
		}
	}
	return posted
}


// mergeServiceExtras pairs each incoming service with the one it came
// from. Real configs contain several services sharing an id, so an exact
// match is tried for every service first, and only then the looser match
// by id alone, which lets a service keep its extra keys after one of its
// displayed fields is edited. Each existing service is consumed once.
func (c *Capybara) mergeServiceExtras(existing []ServiceDefinition) {
	used := make([]bool, len(existing))
	claim := func(match func(ServiceDefinition) bool) (map[string]any, bool) {
		for i, sd := range existing {
			if used[i] || !match(sd) {
				continue
			}
			used[i] = true
			return sd.Extra, true
		}
		return nil, false
	}
	pending := make([]int, 0, len(c.Services))
	for i, incoming := range c.Services {
		extra, ok := claim(func(sd ServiceDefinition) bool {
			return sd.ID == incoming.ID && sd.Pattern == incoming.Pattern && sd.Port == incoming.Port
		})
		if ok {
			c.Services[i].Extra = mergeExtras(extra, incoming.Extra)
			continue
		}
		pending = append(pending, i)
	}
	for _, i := range pending {
		if extra, ok := claim(func(sd ServiceDefinition) bool { return sd.ID == c.Services[i].ID }); ok {
			c.Services[i].Extra = mergeExtras(extra, c.Services[i].Extra)
		}
	}
}
