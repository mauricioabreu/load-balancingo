package server

import (
	"net/http"
	"time"
)

const (
	ReadTimeout  = 2 * time.Second
	WriteTimeout = 2 * time.Second
	IdleTimeout  = 2 * time.Second
)

func NewServer(address string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         address,
		Handler:      handler,
		ReadTimeout:  ReadTimeout,
		WriteTimeout: WriteTimeout,
		IdleTimeout:  IdleTimeout,
	}
}
