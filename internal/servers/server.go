package servers

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

func New(
	handler *handlers.Handler,
) *Server {
	mux := http.NewServeMux()
	return &Server{
		mux:     mux,
		handler: handler,
	}
}

func (s *Server) Run() error {
	log.Printf("Starting http server")

	// load configurations from config file
	config, err := cfg.NewConfiguration()
	if err != nil {
		return err
	}

	// Registering the normal routes
	s.RegisterRoutes()

	// Registering the proxies
	if err := s.RegisterProxies(config.Origins); err != nil {
		return err
	}

	// Running proxy server
	addr := fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port)
	log.Printf("server addr: %s", addr)
	if err := http.ListenAndServe(addr, s.mux); err != nil {
		return fmt.Errorf("could not start the server: %v", err)
	}

	return nil
}
