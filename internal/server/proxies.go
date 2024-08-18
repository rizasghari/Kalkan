package server

import (
	"net/url"

	"github.com/rizasghari/kalkan/internal/models"
	"github.com/rizasghari/kalkan/internal/proxy"
)

func (s *Server) RegisterProxies(origins []models.Origin) error {
	for _, origin := range origins {
		url, err := url.Parse(origin.Url)
		if err != nil {
			return err
		}
		proxy, err := proxy.New(url)
		if err != nil {
			return err
		}
		s.mux.HandleFunc(origin.Edge, proxy.ProxyRequestHandler(url, origin.Edge))
	}

	return nil
}
