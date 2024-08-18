package kalkan

import (
	"log"

	"github.com/rizasghari/kalkan/internal/handlers"
	"github.com/rizasghari/kalkan/internal/server"
)

type Kalkan struct {
}

func New() *Kalkan {
	return &Kalkan{}
}

func (k *Kalkan) Run() error {
	log.Printf("Starting ⛊ KALKAN ⛊ reverse proxy server")
	if err := server.
		New(handlers.New()).
		Run(); err != nil {
		return err
	}
	return nil
}
