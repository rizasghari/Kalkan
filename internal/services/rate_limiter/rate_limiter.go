package rl

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/rizasghari/kalkan/internal/cfg"
	"github.com/rizasghari/kalkan/internal/services/redis"
	"github.com/rizasghari/kalkan/internal/types"
)

type RateLimiter struct {
	timeframe time.Duration
	block     time.Duration
	allowed   int
	history   map[string]*types.Clinet
	mu        sync.Mutex
	redisService *redis.RedisService
}

func New(cfg *cfg.Configuration, redisService *redis.RedisService) *RateLimiter {
	return &RateLimiter{
		timeframe: time.Duration(cfg.RL.Timeframe * int(time.Second)),
		block:     time.Duration(cfg.RL.Block * int(time.Second)),
		allowed:   cfg.RL.Allowed,
		history:   make(map[string]*types.Clinet),
		mu:        sync.Mutex{},
		redisService: redisService,
	}
}

func (rl *RateLimiter) RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := strings.Split(r.RemoteAddr, ":")[0]
		if _, found := rl.history[clientIP]; !found {
			rl.history[clientIP] = &types.Clinet{Count: 0, LastAccess: time.Now()}
		}

		rl.mu.Lock()
		defer rl.mu.Unlock()

		client := rl.history[clientIP]

		if client.BlockedUntil.After(time.Now()) {
			http.Error(w, "You are temporarily blocked due to too many requests", http.StatusTooManyRequests)
			return
		}

		if time.Since(client.LastAccess) > rl.timeframe {
			rl.ResetClientHistory(client)
		}

		if client.Count >= rl.allowed {
			rl.Block(client)
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		client.Count++
		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) ResetRateLimiter() {
	rl.history = make(map[string]*types.Clinet)
}

func (rl *RateLimiter) ResetClientHistory(client *types.Clinet) {
	client.Count = 0
	client.LastAccess = time.Now()
}

func (rl *RateLimiter) Block(client *types.Clinet) {
	client.BlockedUntil = time.Now().Add(rl.block)
}