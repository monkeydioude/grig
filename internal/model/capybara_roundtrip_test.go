package model_test

import (
	"os"
	"path/filepath"
	"testing"

	"monkeydioude/grig/internal/model"
	"monkeydioude/grig/internal/service/file"
)

// Real capybara configs carry proxy fields grig did not model, which a
// save then silently dropped. Loading and saving must preserve them.
const capybaraYAML = `proxy:
    port: 443
    email: someone@example.com
    tls_hosts:
        - a.example.com
        - b.example.com
services:
    - id: api
      method: string
      pattern: /api
      port: 8080
    - id: rpc
      method: string
      pattern: /rpc
      port: 9393
      protocol: rpc
`

func TestCapybaraYAMLRoundTripKeepsProxyFields(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "capybara.yaml")
	if err := os.WriteFile(path, []byte(capybaraYAML), 0o600); err != nil {
		t.Fatal(err)
	}

	loaded, err := file.UnmarshalFromPath[model.Capybara](path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if got := len(loaded.Proxy.TLSHosts); got != 2 {
		t.Errorf("TLSHosts = %v, want 2 entries", loaded.Proxy.TLSHosts)
	}
	if loaded.Proxy.Email != "someone@example.com" {
		t.Errorf("Email = %q, want someone@example.com", loaded.Proxy.Email)
	}

	loaded.Path = path
	loaded.FileWriter = file.CreateAndWriteFile
	if err := loaded.Save(); err != nil {
		t.Fatalf("save: %v", err)
	}

	reloaded, err := file.UnmarshalFromPath[model.Capybara](path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if reloaded.Proxy.Email != loaded.Proxy.Email {
		t.Errorf("Email lost on save: %q", reloaded.Proxy.Email)
	}
	if len(reloaded.Proxy.TLSHosts) != len(loaded.Proxy.TLSHosts) {
		t.Errorf("TLSHosts lost on save: %v", reloaded.Proxy.TLSHosts)
	}
	if len(reloaded.Services) != len(loaded.Services) {
		t.Errorf("Services lost on save: %d want %d", len(reloaded.Services), len(loaded.Services))
	}
	if reloaded.Services[1].Protocol != "rpc" {
		t.Errorf("Protocol lost on save: %q", reloaded.Services[1].Protocol)
	}
}

// A proxy with no TLS is valid, and blank rows submitted by the edit form
// must not survive into the saved config.
func TestProxyWithoutTLSVerifiesAndBlankHostsAreDropped(t *testing.T) {
	cp := model.Capybara{
		Proxy:    model.Proxy{Port: 7000, TLSHosts: []string{"keep.example.com", "", "   "}},
		Services: []model.ServiceDefinition{{ID: "a", Method: "string", Pattern: "/a", Port: 1}},
	}
	if err := cp.Proxy.Verify(); err != nil {
		t.Errorf("proxy without TLS hosts should verify, got %v", err)
	}
	cp.Sanitize()
	if len(cp.Proxy.TLSHosts) != 1 || cp.Proxy.TLSHosts[0] != "keep.example.com" {
		t.Errorf("blank TLS hosts not dropped: %v", cp.Proxy.TLSHosts)
	}
}
