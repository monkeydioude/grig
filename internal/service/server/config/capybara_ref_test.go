package config

import "testing"

func TestCapybaraName(t *testing.T) {
	cases := map[string]string{
		"/conf/capybara.yaml":        "capybara",
		"/conf/garden.capybara.yaml": "garden",
		"/conf/capybara.config.json": "capybara",
		"/conf/prod.yml":             "prod",
	}
	for path, want := range cases {
		if got := capybaraName(path); got != want {
			t.Errorf("capybaraName(%q) = %q, want %q", path, got, want)
		}
	}
}

func TestCapybaraRefsDisambiguatesCollisions(t *testing.T) {
	sc := ServerConfig{CapybaraConfigPaths: []string{
		"/srv/blue/capybara.yaml",
		"/srv/green/capybara.yaml",
		"/srv/garden.capybara.yaml",
	}}
	want := []string{"blue/capybara", "green/capybara", "garden"}
	refs := sc.CapybaraRefs()
	for i, ref := range refs {
		if ref.Name != want[i] || ref.ID != i {
			t.Errorf("ref %d = {ID:%d Name:%q}, want {ID:%d Name:%q}", i, ref.ID, ref.Name, i, want[i])
		}
	}
}
