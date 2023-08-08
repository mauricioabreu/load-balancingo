package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/mauricioabreu/load-balancingo/proxy"
	"github.com/mauricioabreu/load-balancingo/roundrobin"
)

const (
	ReadTimeout  = 2 * time.Second
	WriteTimeout = 2 * time.Second
	IdleTimeout  = 2 * time.Second
)

func newServer(address string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         address,
		Handler:      handler,
		ReadTimeout:  ReadTimeout,
		WriteTimeout: WriteTimeout,
		IdleTimeout:  IdleTimeout,
	}
}

func startServer(address string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Origin-Server", address)
		fmt.Fprintf(w, "Hello from %s\n", address)
	})

	server := newServer(address, mux)

	log.Fatal(server.ListenAndServe())
}

func main() {
	balancer := roundrobin.New(
		roundrobin.NewServer("http://127.0.0.1:8081").WithWeight(5), //nolint:gomnd // 5 seems pretty obvious
		roundrobin.NewServer("http://127.0.0.1:8082"),
	)

	rrproxy := proxy.NewRoundRobinProxy(balancer)

	server := newServer(":8080", nil)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rrproxy.ServeHTTP(w, r)
	})

	go startServer("127.0.0.1:8081")
	go startServer("127.0.0.1:8082")

	log.Fatal(server.ListenAndServe())
}
