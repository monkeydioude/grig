package main

import (
	"context"
	"log"
	"syscall"

	"github.com/oklog/run"
)

func main() {
	serverLayout := boot()
	jsonApiServer := setupJsonApiServer(&serverLayout)
	// Create the run group, and add each actor.
	var runGroup run.Group
	// JSON API goroutine
	runGroup.Add(func() error {
		// Start the server on port
		log.Println("API starting on", jsonApiServer.Addr)
		return jsonApiServer.ListenAndServe()
	}, func(_ error) {
		log.Println("closing API server")
		if err := jsonApiServer.Close(); err != nil {
			log.Println("failed to stop web server", "err", err)
		}
	})

	runGroup.Add(run.SignalHandler(context.Background(), syscall.SIGINT, syscall.SIGTERM))
	if err := runGroup.Run(); err != nil {
		log.Fatal(err)
	}
}
