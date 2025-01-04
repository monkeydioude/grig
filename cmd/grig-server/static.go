package main

import (
	"monkeydioude/grig/pkg/tiger/assert"
	"net/http"
)

func serveStatic(mux *http.ServeMux) {
	assert.NotNil(mux)
	staticDir := "./static"
	fs := http.FileServer(http.Dir(staticDir))
	// mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, staticDir+"/favicon.ico")
	// })
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
}
