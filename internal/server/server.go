package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rizasghari/kalkan/internal/cfg"
	"github.com/rizasghari/kalkan/internal/handlers"
)

type Server struct {
	mux     *http.ServeMux
	handler *handlers.Handler
}

func New(handler *handlers.Handler) *Server {
	mux := http.NewServeMux()
	return &Server{
		mux:     mux,
		handler: handler,
	}
}

func (s *Server) Run() error {
	config, err := cfg.NewConfiguration()
	if err != nil {
		return err
	}

	s.RegisterRoutes()

	if err := s.RegisterProxies(config.Origins); err != nil {
		return err
	}

	addr := fmt.Sprintf(":%s", config.Server.Port)
	log.Printf("server addr: %s", addr)
	if err := http.ListenAndServe(addr, s.mux); err != nil {
		return fmt.Errorf("could not start the server: %v", err)
	}

	return nil
}
