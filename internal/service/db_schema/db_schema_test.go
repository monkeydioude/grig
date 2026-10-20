package db_schema

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testPassword = "abcdefghijklmnopqrstuvwxyz012345"

func TestValidate(t *testing.T) {
	assert.NoError(t, Request{Schema: "kabane", User: "kabane_app", Password: testPassword}.Validate())

	assert.ErrorIs(t, Request{Schema: "Kabane", User: "u", Password: testPassword}.Validate(), ErrInvalidSchema)
	assert.ErrorIs(t, Request{Schema: "s", User: "u;drop", Password: testPassword}.Validate(), ErrInvalidUser)
	assert.ErrorIs(t, Request{Schema: "s", User: "u", Password: "short"}.Validate(), ErrInvalidPassword)
	assert.ErrorIs(t, Request{Schema: "s", User: "u", Password: testPassword + "'"}.Validate(), ErrInvalidPassword)
}

func TestStatementsGrantCreateOnTheDatabase(t *testing.T) {
	sql := strings.Join(Statements("prod", Request{Schema: "s", User: "u", Password: testPassword}), ";\n")

	assert.Contains(t, sql, `CREATE ROLE "u" LOGIN PASSWORD '`+testPassword+`'`)
	assert.Contains(t, sql, `GRANT CONNECT, TEMPORARY, CREATE ON DATABASE "prod" TO "u"`)
	assert.Contains(t, sql, `GRANT USAGE, CREATE ON SCHEMA "s" TO "u"`)
	assert.NotContains(t, sql, "IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles", "an existing user must fail")
}

func TestServiceURL(t *testing.T) {
	got, err := ServiceURL("postgres://postgres:secret@127.0.0.1:5432/prod?sslmode=disable",
		Request{Schema: "s", User: "u", Password: "pw"})

	require.NoError(t, err)
	assert.Equal(t, "postgres://u:pw@127.0.0.1:5432/prod?search_path=s&sslmode=disable", got)
}

func TestRedactedHidesThePassword(t *testing.T) {
	assert.Equal(t, "postgres://postgres:xxxxx@db:5432/prod", Redacted("postgres://postgres:secret@db:5432/prod"))
}

func TestGeneratePasswordIsValid(t *testing.T) {
	password, err := GeneratePassword()

	require.NoError(t, err)
	assert.NoError(t, Request{Schema: "s", User: "u", Password: password}.Validate())
}

// TestCreate runs against a real PostgreSQL when TEST_DB_URL is set
// (an admin URL), e.g. postgres://postgres:postgres@127.0.0.1:5432/postgres.
func TestCreate(t *testing.T) {
	adminURL := os.Getenv("TEST_DB_URL")
	if adminURL == "" {
		t.Skip("TEST_DB_URL not set")
	}
	ctx := context.Background()
	req := Request{Schema: "grig_test_schema", User: "grig_test_user", Password: testPassword}

	res, err := Admin{URL: adminURL}.Create(ctx, req)
	require.NoError(t, err)

	// the service can create its schema again, then a table in it, and use it
	conn, err := pgx.Connect(ctx, res.ServiceURL)
	require.NoError(t, err)
	defer conn.Close(ctx)
	for _, sql := range []string{
		"CREATE SCHEMA IF NOT EXISTS grig_test_schema",
		"CREATE TABLE items (id serial PRIMARY KEY, name text)",
		"INSERT INTO items (name) VALUES ('a')",
	} {
		_, err := conn.Exec(ctx, sql)
		require.NoError(t, err, sql)
	}

	_, err = Admin{URL: adminURL}.Create(ctx, req)
	assert.ErrorContains(t, err, "already exists", "an existing user must not be reused")
}
