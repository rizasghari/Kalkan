package servers

import (
	"net/url"

	"github.com/rizasghari/kalkan/internal/models"
	"github.com/rizasghari/kalkan/internal/proxy"
)

func (s *Server) RegisterProxies(origins []models.Origin) error {
	// Iterating through the origins and registering them into the router.
	for _, origin := range origins {
		url, _ := url.Parse(origin.Url)
		proxy, err := proxy.New(url)
		if err != nil {
			return err
		}
		s.mux.HandleFunc(origin.Edge, proxy.ProxyRequestHandler(url, origin.Edge))
	}

	return nil
}
