package compose

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func testService() Service {
	return Service{
		Name:         "my-svc",
		BuildContext: "/home/mkd/work/my-svc",
		Dockerfile:   "Dockerfile.prod",
		Port:         8400,
		HealthPath:   "/healthcheck",
		EnvFile:      "/home/mkd/envs/my-svc/.env",
	}
}

func TestServiceYAMLParses(t *testing.T) {
	var parsed map[string]map[string]map[string]any
	require.NoError(t, yaml.Unmarshal([]byte("services:\n"+testService().YAML(2)), &parsed))

	svc := parsed["services"]["my-svc"]
	assert.Equal(t, "host", svc["network_mode"])
	assert.Equal(t, map[string]any{"PORT": 8400}, svc["environment"])
	assert.Equal(t, []any{"CMD-SHELL", "wget -q --spider http://localhost:8400/healthcheck || exit 1"},
		svc["healthcheck"].(map[string]any)["test"])
	assert.Equal(t, []any{map[string]any{"path": "/home/mkd/envs/my-svc/.env", "required": false}}, svc["env_file"])
}

func TestServiceYAMLWithoutOptionalParts(t *testing.T) {
	svc := testService()
	svc.HealthPath, svc.EnvFile = "", ""

	out := svc.YAML(2)

	assert.NotContains(t, out, "healthcheck")
	assert.NotContains(t, out, "env_file")
	assert.True(t, strings.HasPrefix(out, "  my-svc:\n    build:\n"))
}

func TestValidate(t *testing.T) {
	assert.NoError(t, testService().Validate())
	for name, edit := range map[string]func(*Service){
		"name":    func(s *Service) { s.Name = "a: b" },
		"port":    func(s *Service) { s.Port = 0 },
		"context": func(s *Service) { s.BuildContext = "relative/path" },
		"health":  func(s *Service) { s.HealthPath = `/x" && rm -rf /` },
	} {
		svc := testService()
		edit(&svc)
		assert.Error(t, svc.Validate(), name)
	}
}

const withVolumes = `services:
  db:
    image: postgres # the database

  # the API
  api:
    build: .

# shared volumes
volumes:
  data:
`

func TestInsertBeforeTheNextTopLevelKey(t *testing.T) {
	out, err := Insert([]byte(withVolumes), testService())
	require.NoError(t, err)

	names, err := ServiceNames(out)
	require.NoError(t, err)
	assert.Equal(t, []string{"db", "api", "my-svc"}, names)
	assert.True(t, strings.HasPrefix(string(out), withVolumes[:strings.Index(withVolumes, "\n# shared")]),
		"what comes before is untouched")
	assert.True(t, strings.HasSuffix(string(out), "  restart: always\n\n# shared volumes\nvolumes:\n  data:\n"),
		"the volumes section and its comment follow the new service")
}

func TestInsertAtTheEndOfTheFile(t *testing.T) {
	content := "services:\n    api:\n        image: x\n\n\n"

	out, err := Insert([]byte(content), testService())
	require.NoError(t, err)

	assert.True(t, strings.HasPrefix(string(out), "services:\n    api:\n        image: x\n\n    my-svc:\n        build:\n"),
		"the existing indentation is kept")
}

func writeComposeDir(t *testing.T) Dir {
	t.Helper()
	dir := t.TempDir()
	files := map[string]string{
		"docker-compose.tools.yml": "services:\n  documentation:\n    image: x\n",
		"docker-compose.prod.yml":  withVolumes,
		"docker-compose.yml.bak":   "not a compose file",
	}
	for name, content := range files {
		require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644))
	}
	return Dir{Base: dir}
}

func TestDirTypes(t *testing.T) {
	types, err := writeComposeDir(t).Types()

	require.NoError(t, err)
	assert.Equal(t, []string{"prod", "tools"}, types)
}

func TestDirAdd(t *testing.T) {
	dir := writeComposeDir(t)

	path, err := dir.Add("tools", testService())

	require.NoError(t, err)
	assert.Equal(t, filepath.Join(dir.Base, "docker-compose.tools.yml"), path)
	content, _ := os.ReadFile(path)
	names, _ := ServiceNames(content)
	assert.Equal(t, []string{"documentation", "my-svc"}, names)
}

func TestDirAddRefusesANameUsedInAnyFile(t *testing.T) {
	dir := writeComposeDir(t)
	svc := testService()
	svc.Name = "api"

	_, err := dir.Add("tools", svc)

	assert.ErrorIs(t, err, ErrNameTaken)
	assert.ErrorContains(t, err, "docker-compose.prod.yml")
	_, err = dir.Add("staging", testService())
	assert.ErrorIs(t, err, ErrUnknownType)
}
