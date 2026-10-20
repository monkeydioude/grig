package env_file

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPathStaysInTheBase(t *testing.T) {
	dir := Dir{Base: "/envs"}

	path, err := dir.Path("my-service")
	require.NoError(t, err)
	assert.Equal(t, "/envs/my-service/.env", path)

	for _, project := range []string{"", "..", "../etc", "a/b", ".hidden", "a..b"} {
		_, err := dir.Path(project)
		assert.ErrorIs(t, err, ErrInvalidProject, project)
	}
}

func TestAddCreatesTheFileOwnerOnly(t *testing.T) {
	path := filepath.Join(t.TempDir(), "svc", ".env")

	require.NoError(t, Add(path, "DB_URL", "postgres://u:p@h/db?search_path=s&sslmode=disable"))

	content, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, "DB_URL=postgres://u:p@h/db?search_path=s&sslmode=disable\n", string(content))
	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
}

func TestAddKeepsExistingVariables(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	require.NoError(t, os.WriteFile(path, []byte("# app\nPORT=8080"), 0o600))

	require.NoError(t, Add(path, "DB_URL", "x"))

	content, _ := os.ReadFile(path)
	assert.Equal(t, "# app\nPORT=8080\nDB_URL=x\n", string(content))
}

func TestAddRefusesAKeyAlreadySet(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	require.NoError(t, os.WriteFile(path, []byte("export DB_URL = old\n"), 0o600))

	has, err := HasKey(path, "DB_URL")
	require.NoError(t, err)
	assert.True(t, has)
	assert.ErrorIs(t, Add(path, "DB_URL", "new"), ErrKeyExists)

	content, _ := os.ReadFile(path)
	assert.Equal(t, "export DB_URL = old\n", string(content))
}
