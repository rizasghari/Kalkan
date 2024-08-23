package rl

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/rizasghari/kalkan/internal/cfg"
	"github.com/rizasghari/kalkan/internal/services/cache"
	"github.com/rizasghari/kalkan/internal/types"
	"github.com/rizasghari/kalkan/internal/utils"
)

type RateLimiter struct {
	timeframe time.Duration
	block     time.Duration
	allowed   int
	history   map[string]*types.Clinet
	mu        sync.Mutex
	cache     cache.Cacher
	cfg       *cfg.Configuration
}

func New(cfg *cfg.Configuration, cache cache.Cacher) *RateLimiter {
	return &RateLimiter{
		timeframe: time.Duration(cfg.RL.Timeframe * int(time.Second)),
		block:     time.Duration(cfg.RL.Block * int(time.Second)),
		allowed:   cfg.RL.Allowed,
		history:   make(map[string]*types.Clinet),
		mu:        sync.Mutex{},
		cache:     cache,
		cfg:       cfg,
	}
}

func (rl *RateLimiter) RateLimiterCacherMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP, err := utils.GetIP(r.RemoteAddr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		key := clientIP.String()
		log.Printf("key: %s", key)

		var client types.Clinet = types.Clinet{
			Count:      0,
			LastAccess: time.Now(),
		}

		if clinetData, err := rl.cache.Get(r.Context(), key); err != nil {
			log.Printf("error reading cache: %v", err)

			json, err := json.Marshal(client)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			err = rl.cache.Set(
				r.Context(),
				key,
				json,
				time.Duration(rl.cfg.Redis.Expiration*int(time.Second)),
			)
			if err != nil {
				log.Printf("error writing cache: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			clinetStr, ok := clinetData.(string)
			if !ok {
				http.Error(w, "error reading cache: clinetData type assertion to string", http.StatusInternalServerError)
			}
			err = json.Unmarshal([]byte(clinetStr), &client)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			log.Printf("client: %+v", client)
		}

		// if _, found := rl.history[clientIP.String()]; !found {
		// 	rl.history[clientIP.String()] = &types.Clinet{Count: 0, LastAccess: time.Now()}
		// }

		// client := rl.history[clientIP.String()]

		// if client.BlockedUntil.After(time.Now()) {
		// 	http.Error(w, "You are temporarily blocked due to too many requests", http.StatusTooManyRequests)
		// 	return
		// }

		// if time.Since(client.LastAccess) > rl.timeframe {
		// 	rl.ResetClientHistory(client)
		// }

		// if client.Count >= rl.allowed {
		// 	rl.Block(client)
		// 	http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		// 	return
		// }

		// client.Count++
		// next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) RateLimiterMapperMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP, err := utils.GetIP(r.RemoteAddr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if _, found := rl.history[clientIP.String()]; !found {
			rl.history[clientIP.String()] = &types.Clinet{Count: 0, LastAccess: time.Now()}
		}

		rl.mu.Lock()
		defer rl.mu.Unlock()

		client := rl.history[clientIP.String()]

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
