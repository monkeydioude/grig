package db_schema

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// Admin creates schemas and users through a connection with the rights
// to (CREATEROLE, and CREATE on the database).
type Admin struct {
	URL string
}

// Enabled tells whether an admin URL (DB_URL) is configured.
func (a Admin) Enabled() bool {
	return a.URL != ""
}

// Result is a created schema and the URL its service connects with.
type Result struct {
	Database   string
	Schema     string
	User       string
	ServiceURL string
}

// Create runs Statements in one transaction: nothing is left half done.
func (a Admin) Create(ctx context.Context, r Request) (Result, error) {
	conn, err := pgx.Connect(ctx, a.URL)
	if err != nil {
		return Result{}, fmt.Errorf("connect: %w", err)
	}
	defer conn.Close(ctx)

	var database string
	if err := pgx.BeginFunc(ctx, conn, func(tx pgx.Tx) error {
		if err := tx.QueryRow(ctx, "SELECT current_database()").Scan(&database); err != nil {
			return err
		}
		return execAll(ctx, tx, Statements(database, r))
	}); err != nil {
		return Result{}, err
	}
	serviceURL, err := ServiceURL(a.URL, r)
	if err != nil {
		return Result{}, err
	}
	return Result{Database: database, Schema: r.Schema, User: r.User, ServiceURL: serviceURL}, nil
}

func execAll(ctx context.Context, tx pgx.Tx, statements []string) error {
	for _, statement := range statements {
		if _, err := tx.Exec(ctx, statement); err != nil {
			return err
		}
	}
	return nil
}
