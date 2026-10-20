// Package db_schema creates a PostgreSQL schema and the user a new
// service connects with, through an admin connection (DB_URL).
package db_schema

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"regexp"

	"github.com/jackc/pgx/v5"
)

var (
	ErrInvalidSchema   = errors.New("invalid schema name: lowercase letters, digits and _, starting with a letter or _")
	ErrInvalidUser     = errors.New("invalid user name: lowercase letters, digits and _, starting with a letter or _")
	ErrInvalidPassword = errors.New("invalid password: 24 to 128 letters and digits")
)

// identifiers are kept lowercase so they never need quoting in the
// services' own SQL; 63 bytes is PostgreSQL's limit.
var (
	validIdentifier = regexp.MustCompile(`^[a-z_][a-z0-9_]{0,62}$`)
	validPassword   = regexp.MustCompile(`^[A-Za-z0-9]{24,128}$`)
)

const (
	passwordChars  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	passwordLength = 32
)

// Request is what the "+ Service" page posts.
type Request struct {
	Schema   string `json:"schema"`
	User     string `json:"user"`
	Password string `json:"password"`
}

func (r Request) Validate() error {
	if !validIdentifier.MatchString(r.Schema) {
		return ErrInvalidSchema
	}
	if !validIdentifier.MatchString(r.User) {
		return ErrInvalidUser
	}
	if !validPassword.MatchString(r.Password) {
		return ErrInvalidPassword
	}
	return nil
}

// Statements is the SQL creating the user and giving it the schema.
// Request must be valid: the password is inlined, as CREATE ROLE takes
// no parameter, which validation makes safe (letters and digits only).
//
// The user owns what it creates. The grants on existing objects and the
// default privileges cover objects the admin creates in the schema.
func Statements(database string, r Request) []string {
	db := pgx.Identifier{database}.Sanitize()
	schema := pgx.Identifier{r.Schema}.Sanitize()
	user := pgx.Identifier{r.User}.Sanitize()
	return []string{
		// fails when the user exists: its password would not be the one shown
		fmt.Sprintf("CREATE ROLE %s LOGIN PASSWORD '%s'", user, r.Password),
		fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schema),
		// CREATE on the database: services run "CREATE SCHEMA IF NOT EXISTS"
		// at start, which PostgreSQL checks even when the schema exists
		fmt.Sprintf("GRANT CONNECT, TEMPORARY, CREATE ON DATABASE %s TO %s", db, user),
		fmt.Sprintf("GRANT USAGE, CREATE ON SCHEMA %s TO %s", schema, user),
		fmt.Sprintf("GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA %s TO %s", schema, user),
		fmt.Sprintf("GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA %s TO %s", schema, user),
		fmt.Sprintf("ALTER DEFAULT PRIVILEGES IN SCHEMA %s GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO %s", schema, user),
		fmt.Sprintf("ALTER DEFAULT PRIVILEGES IN SCHEMA %s GRANT USAGE, SELECT ON SEQUENCES TO %s", schema, user),
	}
}

// ServiceURL is the DB_URL the new service connects with: the admin URL
// with the new user's credentials, and the schema as search_path.
func ServiceURL(adminURL string, r Request) (string, error) {
	u, err := url.Parse(adminURL)
	if err != nil {
		return "", fmt.Errorf("db_schema.ServiceURL: %w", err)
	}
	u.User = url.UserPassword(r.User, r.Password)
	query := u.Query()
	query.Set("search_path", r.Schema)
	u.RawQuery = query.Encode()
	return u.String(), nil
}

// Redacted is the URL without its password, safe to display.
func Redacted(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "invalid DB_URL"
	}
	return u.Redacted()
}

// GeneratePassword returns a random password Request accepts.
func GeneratePassword() (string, error) {
	max := big.NewInt(int64(len(passwordChars)))
	password := make([]byte, passwordLength)
	for i := range password {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", fmt.Errorf("db_schema.GeneratePassword: %w", err)
		}
		password[i] = passwordChars[n.Int64()]
	}
	return string(password), nil
}
