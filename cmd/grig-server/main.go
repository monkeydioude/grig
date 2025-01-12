package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"monkeydioude/grig/internal/consts"
	"monkeydioude/grig/internal/service/server/config"
	"monkeydioude/grig/pkg/server"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/oklog/run"
)

func setupServer(layout *server.Layout[config.ServerConfig]) *http.Server {
	// setup multiplexer
	mux := routing(layout)
	port := consts.DEFAULT_GRIG_SERVER_PORT
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	return &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		ReadTimeout:       3 * time.Second,
		WriteTimeout:      3 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		Handler:           mux,
	}
}

func main() {
	serverLayout := boot()
	server := setupServer(serverLayout)
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
