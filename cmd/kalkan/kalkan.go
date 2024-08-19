package kalkan

import (
	"log"

	"github.com/rizasghari/kalkan/internal/cfg"
	"github.com/rizasghari/kalkan/internal/handlers"
	"github.com/rizasghari/kalkan/internal/server"
)

func Run() error {
	log.Printf("Starting ⛊ KALKAN ⛊ Reverse Proxy Server")
	cfg, err := cfg.NewConfiguration()
	if err != nil {
		return err
	}
	if err := server.
		New(handlers.New(), cfg).
		Start(); err != nil {
		return err
	}
	return nil
}
