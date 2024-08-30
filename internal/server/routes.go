package server

import (
	"net/http"

	"github.com/rizasghari/kalkan/internal/services/geolocation"
	rl "github.com/rizasghari/kalkan/internal/services/rate_limiter"
)

func (s *Server) RegisterRoutes(rl *rl.RateLimiter, geoLocation *geolocation.GeoLocation) {
	healthHandler := http.HandlerFunc(s.handler.Health)
	if rl != nil {
		rateLimitedHealthHandler := rl.RateLimiterCacherMiddleware(healthHandler, geoLocation)
		s.mux.Handle("/health", rateLimitedHealthHandler)
	} else {
		s.mux.Handle("/health", healthHandler)
	}
}