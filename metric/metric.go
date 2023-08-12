package metric

import (
	"log"
	"net/http"

	"github.com/mauricioabreu/load-balancingo/server"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	algorithmUsedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "algorithm_used_total",
			Help: "Number of times a server was chosen including the given algorithm",
		},
		[]string{"server", "algorithm"},
	)
)

func CountUsedAlgorithm(srv, algorithm string) {
	algorithmUsedTotal.WithLabelValues(srv, algorithm).Inc()
}

func StartServer() {
	prometheus.MustRegister(algorithmUsedTotal)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	srv := server.NewServer(":9394", mux)
	log.Fatal(srv.ListenAndServe())
}
