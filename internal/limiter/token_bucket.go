package limiter

import (
	"context"
	"math"
	"sync"
	"time"
)

type TokenBucket struct {
	tokens     float64    // current token count
	lastRefill time.Time  // when the refill was last computed
	capacity   float64    // max tokens
	rate       float64    // tokens per second
	mu         sync.Mutex // goroutines will call Allow concurrently
}

func NewTokenBucket(capacity float64, rate float64) *TokenBucket {
	tb := new(TokenBucket)
	tb.capacity = capacity
	tb.rate = rate
	tb.tokens = capacity
	tb.lastRefill = time.Now()
	return tb
}

func (tb *TokenBucket) Allow(ctx context.Context, clientID string) (bool, int64, time.Time, error) {
	tb.mu.Lock()
	defer tb.mu.Unlock() // release lock once Allow returns

	elapsed := time.Since(tb.lastRefill).Seconds() // convert to seconds / type = float64
	// compute new tokens
	tb.tokens = math.Min(tb.tokens+elapsed*tb.rate, tb.capacity)
	tb.lastRefill = time.Now()

	if tb.tokens >= 1 {
		tb.tokens -= 1
		return true, int64(tb.tokens), time.Now(), nil
	}
	// denied, so compute how many seconds until the next token arrives
	// time = distance / speed
	timeToNextToken := (1 - tb.tokens) / tb.rate
	return false, int64(tb.tokens), time.Now().Add(time.Duration(timeToNextToken * float64(time.Second))), nil

}

func (tb *TokenBucket) Close() error {
	return nil
}
