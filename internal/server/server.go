package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rizasghari/kalkan/internal/cfg"
	"github.com/rizasghari/kalkan/internal/handlers"
	rl "github.com/rizasghari/kalkan/internal/services/rate_limiter"
	"github.com/rizasghari/kalkan/internal/services/redis"
)

type Server struct {
	mux          *http.ServeMux
	handler      *handlers.Handler
	cfg          *cfg.Configuration
	redisService *redis.RedisService
}

func New(
	handler *handlers.Handler,
	cfg *cfg.Configuration,
	redisService *redis.RedisService,
) *Server {
	mux := http.NewServeMux()
	return &Server{
		mux:          mux,
		handler:      handler,
		cfg:          cfg,
		redisService: redisService,
	}
}

func (s *Server) Start() error {
	var rateLimiter *rl.RateLimiter
	if s.cfg.RL.Enabled {
		log.Println("Rate limiter enabled")
		rateLimiter = rl.New(s.cfg, s.redisService)
	}

	s.RegisterRoutes(rateLimiter)

	if err := s.RegisterProxies(s.cfg.Origins, rateLimiter); err != nil {
		return err
	}

	addr := fmt.Sprintf(":%s", s.cfg.Server.Port)
	if err := http.ListenAndServe(addr, s.mux); err != nil {
		return fmt.Errorf("could not start the server: %v", err)
	}
	log.Printf("KALKAN HTTP Server started on port %s", addr)

	return nil
}
