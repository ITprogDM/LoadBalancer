package ratelimiter

type Limiter interface {
	Allow(clientID string) bool
}
