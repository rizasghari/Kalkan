package server

import (
	"net/http"
	"net/url"

	"github.com/rizasghari/kalkan/internal/services/gateway"
	rl "github.com/rizasghari/kalkan/internal/services/rate_limiter"
	"github.com/rizasghari/kalkan/internal/types"
)

func (s *Server) RegisterProxies(origins []types.Origin, rateLimiter *rl.RateLimiter) error {
	for _, origin := range origins {
		url, err := url.Parse(origin.Url)
		if err != nil {
			return err
		}
		proxy, err := gateway.New(url)
		if err != nil {
			return err
		}

		proxyHandler := http.HandlerFunc(proxy.ProxyRequestHandler(url, origin.Edge))
		if rateLimiter != nil {
			rateLimitedProxyHandler := rateLimiter.RateLimiterMiddleware(proxyHandler)
			s.mux.Handle(origin.Edge, rateLimitedProxyHandler)
		} else {
			s.mux.Handle(origin.Edge, proxyHandler)
		}
	}

	return nil
}
