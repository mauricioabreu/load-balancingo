package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mauricioabreu/load-balancingo/metric"
	"github.com/mauricioabreu/load-balancingo/p2c"
	"github.com/mauricioabreu/load-balancingo/proxy"
	"github.com/mauricioabreu/load-balancingo/roundrobin"
	"github.com/mauricioabreu/load-balancingo/server"
)

func startServer(address string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Origin-Server", address)
		fmt.Fprintf(w, "Hello from %s\n", address)
	})

	srv := server.NewServer(address, mux)

	log.Fatal(srv.ListenAndServe())
}

func main() {
	rr := roundrobin.New(
		roundrobin.NewServer("http://127.0.0.1:8081").WithWeight(5), //nolint:gomnd // 5 seems pretty obvious
		roundrobin.NewServer("http://127.0.0.1:8082"),
	)

	pwr := p2c.New(
		p2c.NewServer("http://127.0.0.1:8081", &p2c.RandomLoadFetcher{}),
		p2c.NewServer("http://127.0.0.1:8082", &p2c.RandomLoadFetcher{}),
	)

	rrproxy := proxy.NewProxy(rr)
	pwrproxy := proxy.NewProxy(pwr)

	srv := server.NewServer(":8080", nil)

	http.HandleFunc("/rr", func(w http.ResponseWriter, r *http.Request) {
		rrproxy.ServeHTTP(w, r)
	})
	http.HandleFunc("/p2c", func(w http.ResponseWriter, r *http.Request) {
		pwrproxy.ServeHTTP(w, r)
	})

	go startServer("127.0.0.1:8081")
	go startServer("127.0.0.1:8082")
	go metric.StartServer()

	log.Fatal(srv.ListenAndServe())
}
