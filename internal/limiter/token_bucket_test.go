package limiter

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
)

func TestAllow_BasicAllow(t *testing.T) {
	// create a bucket with capacity 1, rate 1
	tb := NewTokenBucket(1, 1)

	// Background is a no-op context for tests
	allowed, _, _, err := tb.Allow(context.Background(), "client1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !allowed {
		t.Errorf("expected allowed=true, got false")
	}

}

func TestAllow_AllowExhausted(t *testing.T) {
	// create a bucket with capacity 1, rate 1
	tb := NewTokenBucket(1, 1)

	// Background is a no-op context for tests
	allowed1, _, _, err1 := tb.Allow(context.Background(), "client1")
	allowed2, _, _, err2 := tb.Allow(context.Background(), "client1")

	if err1 != nil {
		t.Fatalf("unexpected error1: %v", err1)
	}
	if !allowed1 {
		t.Errorf("expected allowed1=true, got false")
	}

	if err2 != nil {
		t.Fatalf("unexpected error2 error?: %v", err2)
	}
	if allowed2 {
		t.Errorf("expected allowed2=false, got true")
	}
}

// concurrency testing
func TestAllow_Concurrent(t *testing.T) {
	tb := NewTokenBucket(10, 100)
	var wg sync.WaitGroup
	var allowedCount atomic.Int64

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			allowed, _, _, err := tb.Allow(context.Background(), "client1")
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if allowed {
				allowedCount.Add(1)
			}
		}()
	}

	wg.Wait()

	if allowedCount.Load() > 10 {
		t.Errorf("expected at most 10 allowed, got %d", allowedCount.Load())
	}
}
