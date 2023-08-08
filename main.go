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

func startServer(address string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Origin-Server", address)
		fmt.Fprintf(w, "Hello from %s\n", address)
	})

	server := &http.Server{
		Addr:         address,
		Handler:      mux,
		ReadTimeout:  ReadTimeout,
		WriteTimeout: WriteTimeout,
		IdleTimeout:  IdleTimeout,
	}

	log.Fatal(server.ListenAndServe())
}

func main() {
	balancer := roundrobin.New(
		roundrobin.NewServer("http://127.0.0.1:8081").WithWeight(5), //nolint:gomnd // 5 seems pretty obvious
		roundrobin.NewServer("http://127.0.0.1:8082"),
	)

	rrproxy := proxy.NewRoundRobinProxy(balancer)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      nil,
		ReadTimeout:  ReadTimeout,
		WriteTimeout: WriteTimeout,
		IdleTimeout:  IdleTimeout,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rrproxy.ServeHTTP(w, r)
	})

	go startServer("127.0.0.1:8081")
	go startServer("127.0.0.1:8082")

	log.Fatal(server.ListenAndServe())
}
