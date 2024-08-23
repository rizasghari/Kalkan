package server

import (
	"net/http"

	rl "github.com/rizasghari/kalkan/internal/services/rate_limiter"
)

func (s *Server) RegisterRoutes(rl *rl.RateLimiter) {
	healthHandler := http.HandlerFunc(s.handler.Health)
	if rl != nil {
		rateLimitedHealthHandler := rl.RateLimiterCacherMiddleware(healthHandler)
		s.mux.Handle("/health", rateLimitedHealthHandler)
	} else {
		s.mux.Handle("/health", healthHandler)
	}
}