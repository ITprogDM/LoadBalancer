package ratelimiter

import (
	"LoadBalancer/internal/models"
	"LoadBalancer/internal/repository"
	"context"
	"sync"
	"time"
)

type Bucket struct {
	cap          int
	tokens       int
	refillRate   int
	lastRefilled time.Time
	mu           sync.Mutex
}

type RateLimiter struct {
	clients map[string]*Bucket
	mu      sync.RWMutex
	repo    repository.ClientRepository
}

func NewRateLimiterWithRepository(repo repository.ClientRepository) *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string]*Bucket),
		repo:    repo,
	}

	go rl.refillTokens()
	return rl
}

func (rl *RateLimiter) getClientBucket(clientID string) *Bucket {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if bucket, exists := rl.clients[clientID]; exists {
		return bucket
	}

	// Ищем лимиты клиента в БД
	client, err := rl.repo.GetClient(context.Background(), clientID)
	if err != nil {
		// Если клиента нет в БД — дефолтные значения
		client = &models.Client{
			ID:         clientID,
			Capacity:   10,
			RatePerSec: 1,
		}
	}

	bucket := &Bucket{
		cap:          client.Capacity,
		tokens:       client.Capacity,
		refillRate:   client.RatePerSec,
		lastRefilled: time.Now(),
	}

	rl.clients[clientID] = bucket
	return bucket
}

func (rl *RateLimiter) refillTokens() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.RLock()
		for _, bucket := range rl.clients {
			bucket.mu.Lock()

			elapsed := time.Since(bucket.lastRefilled).Seconds()
			tokensToAdd := int(elapsed * float64(bucket.refillRate))
			if tokensToAdd > 0 {
				bucket.tokens += tokensToAdd
				if bucket.tokens > bucket.cap {
					bucket.tokens = bucket.cap
				}
				bucket.lastRefilled = time.Now()
			}

			bucket.mu.Unlock()
		}
		rl.mu.RUnlock()
	}
}

func (rl *RateLimiter) Allow(clientID string) bool {
	bucket := rl.getClientBucket(clientID)

	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	if bucket.tokens > 0 {
		bucket.tokens--
		return true
	}

	return false
}
