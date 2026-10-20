package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"monkeydioude/grig/internal/consts"
	"monkeydioude/grig/internal/service/compose"
	"monkeydioude/grig/internal/service/db_schema"
	"monkeydioude/grig/internal/service/env_file"
	"monkeydioude/grig/internal/service/forge"
	"monkeydioude/grig/internal/service/restart"
	"monkeydioude/grig/internal/service/server/config"
	grigOs "monkeydioude/grig/pkg/os"
	"monkeydioude/grig/pkg/server"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/oklog/run"
)

func setupForges() forge.Registry {
	forges := forge.FromEnv(os.Getenv, &http.Client{Timeout: 10 * time.Second})
	slog.Info("forges available", "forges", forges.Names())
	return forges
}

func setupSchemas() db_schema.Admin {
	schemas := db_schema.Admin{URL: os.Getenv("DB_URL")}
	slog.Info("database schemas", "enabled", schemas.Enabled())
	return schemas
}

func setupRestarter(conf config.ServerConfig) restart.Restarter {
	return restart.Restarter{
		Units:   conf.RestartOnSave,
		Restart: grigOs.ServiceRestart,
		Client:  &http.Client{Timeout: 5 * time.Second},
	}
}

func setupServer(layout *server.Layout[config.ServerConfig], deps dependencies) *http.Server {
	// setup multiplexer
	mux := routing(layout, deps)
	port := consts.DEFAULT_GRIG_SERVER_PORT
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	return &http.Server{
		Addr:        fmt.Sprintf(":%s", port),
		ReadTimeout: 3 * time.Second,
		// creating a repository waits on two forge API calls
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		Handler:           mux,
	}
}

func main() {
	serverLayout := boot()
	server := setupServer(serverLayout, dependencies{
		Forges:    setupForges(),
		Schemas:   setupSchemas(),
		Restarter: setupRestarter(serverLayout.ServerConfig),
		Envs:      env_file.Dir{Base: serverLayout.ServerConfig.EnvBasePath},
		Compose:   compose.Dir{Base: serverLayout.ServerConfig.DockerComposeBasePath},
	})
	// Create the run group, and add each actor.
	var runGroup run.Group
	// JSON API goroutine
	runGroup.Add(func() error {
		// Start the server on port
		slog.Info("API started", "addr", server.Addr)
		return server.ListenAndServe()
	}, func(_ error) {
		slog.Info("closing API server")
		if err := server.Close(); err != nil {
			slog.Info("failed to stop web server", "err", err)
		}
	})

	runGroup.Add(run.SignalHandler(context.Background(), syscall.SIGINT, syscall.SIGTERM))
	if err := runGroup.Run(); err != nil {
		log.Fatal(err)
	}
}
