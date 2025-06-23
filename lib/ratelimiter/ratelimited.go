package ratelimited

import (
	"time"

	"golang.org/x/time/rate"
)

type RateLimited interface {
	IsAllowed([]byte) bool
}

type RateLimiter struct {
	limiter *rate.Limiter
}

func NewRateLimiter(rateLimit int) *RateLimiter {
	return &RateLimiter{
		limiter: rate.NewLimiter(rate.Limit(rateLimit), rateLimit),
	}
}
func (rl *RateLimiter) IsAllowed(data []byte) bool {
	return rl.limiter.AllowN(time.Now(), len(data))
}
