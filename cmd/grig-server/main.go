package main

import (
	"context"
	"fmt"
	"log"
	"monkeydioude/grig/internal/consts"
	"monkeydioude/grig/internal/service/server"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/oklog/run"
)

func setupServer(layout *server.Layout) *http.Server {
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
		log.Println("API starting on", server.Addr)
		return server.ListenAndServe()
	}, func(_ error) {
		log.Println("closing API server")
		if err := server.Close(); err != nil {
			log.Println("failed to stop web server", "err", err)
		}
	})

	runGroup.Add(run.SignalHandler(context.Background(), syscall.SIGINT, syscall.SIGTERM))
	if err := runGroup.Run(); err != nil {
		log.Fatal(err)
	}
}
