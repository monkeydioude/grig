package model_test

import (
	"os"
	"path/filepath"
	"testing"

	"monkeydioude/grig/internal/model"
	"monkeydioude/grig/internal/service/file"
	"monkeydioude/grig/internal/service/payload"
)

// Shapes taken from real capybara configs: a service routed only by
// "host", one with a "redirect" grig does not model, and one with no
// "method" at all.
const realWorldYAML = `proxy:
    port: 443
    email: someone@example.com
    tls_hosts:
        - a.example.com
    acme_dir: /var/lib/acme
services:
    - id: sb_front
      method: string
      pattern: /
      port: 4057
      host: app.example.com
    - id: sb_front_home
      method: string
      pattern: /
      port: 4056
      host: example.com
    - id: favicon
      pattern: ^/favicon.*$
      redirect: /
      port: 7057
log_level: debug
`

const realWorldJSON = `{
  "proxy": {"port": 443, "tls_hosts": ["a.example.com"], "acme_dir": "/var/lib/acme"},
  "services": [
    {"id": "sb_front", "method": "string", "pattern": "/", "port": 4057, "host": "app.example.com"},
    {"id": "favicon", "pattern": "^/favicon.*$", "redirect": "/", "port": 7057}
  ],
  "log_level": "debug"
}`

func saveAndReload(t *testing.T, name, body string) model.Capybara {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	loaded, err := file.UnmarshalFromPath[model.Capybara](path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	loaded.Path = path
	loaded.FileWriter = file.CreateAndWriteFile
	// go through the same validation the save endpoint applies
	if err := payload.VerifyAndSanitizeCapybara(&loaded); err != nil {
		t.Fatalf("verify: %v", err)
	}
	if err := loaded.Save(); err != nil {
		t.Fatalf("save: %v", err)
	}
	reloaded, err := file.UnmarshalFromPath[model.Capybara](path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	return reloaded
}

func TestSaveKeepsHostRedirectAndUnknownKeys(t *testing.T) {
	for _, tc := range []struct {
		name, file, body string
		services         int
	}{
		{"yaml", "capybara.yaml", realWorldYAML, 3},
		{"json", "capybara.json", realWorldJSON, 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := saveAndReload(t, tc.file, tc.body)

			if len(got.Services) != tc.services {
				t.Fatalf("services = %d, want %d (a service was deleted)", len(got.Services), tc.services)
			}
			// host distinguishes services that share a pattern
			if got.Services[0].Host != "app.example.com" {
				t.Errorf("service host lost: %q", got.Services[0].Host)
			}
			// redirect is not modelled, but must survive
			last := got.Services[len(got.Services)-1]
			if last.Extra["redirect"] != "/" {
				t.Errorf("service redirect lost: %v", last.Extra)
			}
			// a service with no method must not be dropped
			if last.ID != "favicon" {
				t.Errorf("service without method was dropped, last is %q", last.ID)
			}
			// unknown keys at proxy and root level
			if got.Proxy.Extra["acme_dir"] != "/var/lib/acme" {
				t.Errorf("proxy unknown key lost: %v", got.Proxy.Extra)
			}
			if got.Extra["log_level"] != "debug" {
				t.Errorf("root unknown key lost: %v", got.Extra)
			}
		})
	}
}

// Clearing a field in the UI must not be undone by a stale extra.
func TestModelledFieldsWinOverExtras(t *testing.T) {
	cp := model.Capybara{
		Proxy:    model.Proxy{Port: 80, Email: "", Extra: map[string]any{"email": "stale@example.com"}},
		Services: []model.ServiceDefinition{{ID: "a", Pattern: "/a", Port: 1}},
	}
	raw, err := cp.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	var back model.Capybara
	if err := back.UnmarshalJSON(raw); err != nil {
		t.Fatal(err)
	}
	if back.Proxy.Email != "" {
		t.Errorf("cleared email was resurrected from extras: %q", back.Proxy.Email)
	}
}

// The edit form posts only the displayed fields, so this is what actually
// reaches the save endpoint. Everything else must come back from disk.
func TestMergeExtrasFromDiskSurvivesFormShapedPayload(t *testing.T) {
	onDisk := model.Capybara{
		Proxy: model.Proxy{Port: 443, Extra: map[string]any{"acme_dir": "/var/lib/acme"}},
		Services: []model.ServiceDefinition{
			{ID: "dashboard", Pattern: "/dashboard", Port: 8202},
			{ID: "grafana", Pattern: "/grafana", Port: 3000, Extra: map[string]any{"weight": 5}},
			{ID: "dashboard", Pattern: "/", Port: 8202, Extra: map[string]any{"redirect": "/dashboard"}},
		},
		Extra: map[string]any{"log_level": "debug"},
	}
	// as posted by the form: same order, no extras, one field edited
	posted := model.Capybara{
		Proxy: model.Proxy{Port: 443},
		Services: []model.ServiceDefinition{
			{ID: "dashboard", Pattern: "/dashboard", Port: 8202},
			{ID: "grafana", Pattern: "/grafana-new", Port: 3000},
			{ID: "dashboard", Pattern: "/", Port: 8202},
		},
	}
	posted.MergeExtrasFrom(onDisk)

	if posted.Extra["log_level"] != "debug" {
		t.Errorf("root extras lost: %v", posted.Extra)
	}
	if posted.Proxy.Extra["acme_dir"] != "/var/lib/acme" {
		t.Errorf("proxy extras lost: %v", posted.Proxy.Extra)
	}
	// the two services sharing an id must not swap their keys
	if posted.Services[0].Extra["redirect"] != nil {
		t.Errorf("extras attached to the wrong dashboard: %v", posted.Services[0].Extra)
	}
	if posted.Services[2].Extra["redirect"] != "/dashboard" {
		t.Errorf("redirect lost on the second dashboard: %v", posted.Services[2].Extra)
	}
	// an edited pattern falls back to matching by id
	if posted.Services[1].Extra["weight"] != 5 {
		t.Errorf("extras lost after editing a displayed field: %v", posted.Services[1].Extra)
	}
}

// The modal posts every extra back as text. Editing one must not change
// the types of the others, and keys the modal shows read-only must stay.
func TestEditedExtrasKeepTheirTypes(t *testing.T) {
	onDisk := model.Capybara{
		Services: []model.ServiceDefinition{{
			ID: "api", Pattern: "/api", Port: 80,
			Extra: map[string]any{
				"redirect": "/old",
				"weight":   5,
				"enabled":  true,
				"headers":  map[string]any{"X-Real-IP": "$remote_addr"},
				"stale":    "drop me",
			},
		}},
	}
	posted := model.Capybara{
		Services: []model.ServiceDefinition{{
			ID: "api", Pattern: "/api", Port: 80,
			// as the form posts them: all strings, nested key absent
			Extra: map[string]any{
				"redirect": "/new",
				"weight":   "7",
				"enabled":  "false",
				"stale":    "",
			},
		}},
	}
	posted.MergeExtrasFrom(onDisk)
	got := posted.Services[0].Extra

	if got["redirect"] != "/new" {
		t.Errorf("redirect = %v, want /new", got["redirect"])
	}
	if got["weight"] != 7 {
		t.Errorf("weight = %#v, want int 7 (type must survive the form)", got["weight"])
	}
	if got["enabled"] != false {
		t.Errorf("enabled = %#v, want bool false", got["enabled"])
	}
	if _, ok := got["headers"].(map[string]any); !ok {
		t.Errorf("nested key shown read-only was dropped: %#v", got["headers"])
	}
	if _, present := got["stale"]; present {
		t.Errorf("key cleared in the modal was not removed: %#v", got["stale"])
	}
}
