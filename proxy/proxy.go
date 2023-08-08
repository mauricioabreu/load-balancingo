package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Balancer interface {
	NextAddress() string
}

// NewProxy implements a reverse proxy that will load balance requests.
// It receives a Balancer interface that will be used to get the next address
// working as a balancing strategy.
func NewProxy(b Balancer) *httputil.ReverseProxy {
	p := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			target, err := url.Parse(b.NextAddress())
			if err != nil {
				log.Fatal("Error parsing URL:", err)
			}
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			req.Header.Add("X-Forwarded-Host", req.Header.Get("Host"))
			req.Host = target.Host
		},
	}

	return p
}
