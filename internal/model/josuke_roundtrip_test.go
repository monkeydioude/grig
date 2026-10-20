package model_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"monkeydioude/grig/internal/model"
	"monkeydioude/grig/internal/service/file"
)

// Josuke configs are written in YAML with commands as sequences and keys
// grig does not model (jobs_api, queue_file). Loading and saving must
// keep all of them.
const josukeYAML = `logLevel: "INFO"
host: "localhost"
port: 8082
store: "/data/josuke"
healthcheck_route: "/josuke/healthcheck"
jobs_api:
  token: josuke:secret
queue_file: "/config/queue.json"
hook:
  - name: "github"
    type: "github"
    path: "/josuke/github"
    secret: "s3cret"
deployment:
  - repo: "monkeydioude/grig"
    proj_dir: "grig"
    base_dir: "/home/mkd/work"
    branches:
      - branch: "main"
        actions:
          - action: "push"
            commands:
              - ["cd", "%base_dir%"]
              - ["git", "pull", "--rebase"]
`

func loadJosuke(t *testing.T, path string) model.Josuke {
	t.Helper()
	jk, err := file.UnmarshalFromPath[model.Josuke](path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	return jk
}

func assertJosukeContent(t *testing.T, jk model.Josuke) {
	t.Helper()
	if jk.LogLevel != "INFO" || jk.HealthcheckRoute != "/josuke/healthcheck" || jk.Port != 8082 {
		t.Errorf("top-level fields = %q %q %d", jk.LogLevel, jk.HealthcheckRoute, jk.Port)
	}
	dep := jk.Deployment[0]
	if dep.ProjDir != "grig" || dep.BaseDir != "/home/mkd/work" {
		t.Errorf("deployment dirs = %q %q", dep.ProjDir, dep.BaseDir)
	}
	cmds := dep.Branches[0].Actions[0].Commands
	if len(cmds) != 2 || strings.Join(cmds[1].Parts, " ") != "git pull --rebase" {
		t.Errorf("commands = %+v", cmds)
	}
	if jk.Extra["queue_file"] != "/config/queue.json" || jk.Extra["jobs_api"] == nil {
		t.Errorf("extras = %v", jk.Extra)
	}
}

func TestJosukeYAMLRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "josuke.yaml")
	if err := os.WriteFile(path, []byte(josukeYAML), 0o600); err != nil {
		t.Fatal(err)
	}

	loaded := loadJosuke(t, path)
	assertJosukeContent(t, loaded)

	loaded.Path = path
	loaded.FileWriter = file.CreateAndWriteFile
	if err := loaded.Save(); err != nil {
		t.Fatalf("save: %v", err)
	}
	saved, _ := os.ReadFile(path)
	if !strings.Contains(string(saved), `- ["cd", "%base_dir%"]`) {
		t.Errorf("commands not written as YAML sequences:\n%s", saved)
	}
	assertJosukeContent(t, loadJosuke(t, path))
}
