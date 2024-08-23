package kalkan

import (
	"log"

	"github.com/rizasghari/kalkan/internal/cfg"
	"github.com/rizasghari/kalkan/internal/handlers"
	"github.com/rizasghari/kalkan/internal/server"
	"github.com/rizasghari/kalkan/internal/services/redis"
)

func Run() error {
	log.Printf("Starting ⛊ KALKAN ⛊ Reverse Proxy Server")

	cfg, err := cfg.NewConfiguration()
	if err != nil {
		return err
	}

	redisService := redis.Initialize(cfg)
	handler := handlers.New()

	if err := server.New(handler, cfg, redisService).
		Start(); err != nil {
		return err
	}
	return nil
}
