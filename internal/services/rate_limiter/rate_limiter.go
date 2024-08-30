package rl

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/rizasghari/kalkan/internal/cfg"
	"github.com/rizasghari/kalkan/internal/services/cache"
	"github.com/rizasghari/kalkan/internal/services/geolocation"
	"github.com/rizasghari/kalkan/internal/types"
	"github.com/rizasghari/kalkan/internal/utils"
)

type RateLimiter struct {
	timeframe time.Duration
	block     time.Duration
	allowed   int
	mu        sync.Mutex
	cache     cache.Cacher
	cfg       *cfg.Configuration
}

func New(cfg *cfg.Configuration, cache cache.Cacher) *RateLimiter {
	return &RateLimiter{
		timeframe: time.Duration(cfg.RL.Timeframe * int(time.Second)),
		block:     time.Duration(cfg.RL.Block * int(time.Second)),
		allowed:   cfg.RL.Allowed,
		mu:        sync.Mutex{},
		cache:     cache,
		cfg:       cfg,
	}
}

func (rl *RateLimiter) RateLimiterCacherMiddleware(next http.Handler, geoLocation *geolocation.GeoLocation) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		location, err := geoLocation.GetLocation(r.RemoteAddr)
		if err != nil {
			log.Printf("error getting location: %v", err)
		}
		log.Printf("location for %s IP Address: %+v", r.RemoteAddr, location)

		clientIP, err := utils.GetIP(r.RemoteAddr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		key := clientIP.String()
		log.Printf("key: %s", key)

		var client = &types.Clinet{
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
			err = json.Unmarshal([]byte(clinetStr), client)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			log.Printf("client: %+v", client)
		}

		if client.BlockedUntil.After(time.Now()) {
			http.Error(w, "You are temporarily blocked due to too many requests", http.StatusTooManyRequests)
			return
		}

		if time.Since(client.LastAccess) > rl.timeframe {
			client, err = rl.ResetClientHistory(r.Context(), key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

		if client.Count >= rl.allowed {
			client, err = rl.Block(r.Context(), client, key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		err = rl.IncreaseClientUsage(r.Context(), key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) ResetRateLimiter(ctx context.Context) {
	rl.cache.ResetAll(ctx)
}

func (rl *RateLimiter) ResetClientHistory(ctx context.Context, key string) (*types.Clinet, error) {
	var client types.Clinet
	clinetData, err := rl.cache.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	clinetStr, ok := clinetData.(string)
	if !ok {
		return nil, errors.New("error reading cache: clinetData type assertion to string")
	}
	err = json.Unmarshal([]byte(clinetStr), &client)
	if err != nil {
		return nil, err
	}

	client.Count = 0
	client.LastAccess = time.Now()
	value, err := json.Marshal(client)
	if err != nil {
		return nil, err
	}
	err = rl.cache.Set(
		ctx,
		key,
		value,
		time.Duration(rl.cfg.Redis.Expiration*int(time.Second)),
	)
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (rl *RateLimiter) Block(ctx context.Context, client *types.Clinet, key string) (*types.Clinet, error) {
	client.BlockedUntil = time.Now().Add(rl.block)
	value, err := json.Marshal(client)
	if err != nil {
		return nil, err
	}
	err = rl.cache.Set(
		ctx,
		key,
		value,
		time.Duration(rl.cfg.Redis.Expiration*int(time.Second)),
	)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (rl *RateLimiter) IncreaseClientUsage(ctx context.Context, key string) error {
	var client types.Clinet
	clinetData, err := rl.cache.Get(ctx, key)
	if err != nil {
		return err
	}
	clinetStr, ok := clinetData.(string)
	if !ok {
		return errors.New("error reading cache: clinetData type assertion to string")
	}
	err = json.Unmarshal([]byte(clinetStr), &client)
	if err != nil {
		return err
	}

	client.Count += 1
	value, err := json.Marshal(client)
	if err != nil {
		return err
	}
	return rl.cache.Set(
		ctx,
		key,
		value,
		time.Duration(rl.cfg.Redis.Expiration*int(time.Second)),
	)
}
