package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func newServer(mux *mux.Router, portNumber string) *http.Server {
	return &http.Server{
		Addr:         ":" + portNumber,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
		Handler:      mux,
	}
}
