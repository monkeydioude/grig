package urls

import (
	"context"
	"testing"
)

func TestNewNormalizes(t *testing.T) {
	for _, in := range []string{"grig", "/grig", "/grig/", " /grig/ "} {
		if got := New(in).Base(); got != "/grig" {
			t.Errorf("New(%q).Base() = %q, want %q", in, got, "/grig")
		}
	}
	for _, in := range []string{"", "/", "  "} {
		if got := New(in).Base(); got != "" {
			t.Errorf("New(%q).Base() = %q, want empty", in, got)
		}
	}
}

func TestPath(t *testing.T) {
	p := New("/grig")
	cases := map[string]string{
		"/":                  "/grig/",
		"/josuke":            "/grig/josuke",
		"/api/v1/capybara/0": "/grig/api/v1/capybara/0",
		"https://cdn/x.css":  "https://cdn/x.css",
		"//cdn/x.css":        "//cdn/x.css",
		"#anchor":            "#anchor",
	}
	for in, want := range cases {
		if got := p.Path(in); got != want {
			t.Errorf("Path(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestZeroPrefixerAndMissingContextLeavePathsAlone(t *testing.T) {
	if got := New("").Path("/josuke"); got != "/josuke" {
		t.Errorf("zero Prefixer changed path: %q", got)
	}
	if got := Path(context.Background(), "/josuke"); got != "/josuke" {
		t.Errorf("missing context changed path: %q", got)
	}
}
