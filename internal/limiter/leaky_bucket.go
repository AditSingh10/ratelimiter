package limiter

import (
	"context"
	"time"
)

type LeakyBucket struct {
	queue chan struct{}
	rate  time.Duration
	done  chan struct{} // only used to stop goroutines
}

func NewLeakyBucket(capacity int, rate time.Duration) *LeakyBucket {
	lb := new(LeakyBucket)
	lb.queue = make(chan struct{}, capacity)
	lb.rate = rate
	lb.done = make(chan struct{})

	go func() {
		ticker := time.NewTicker(rate)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				select {
				case <-lb.queue:
				default: // do nothing if queue is already empty
				}

			case <-lb.done:
				return
			}
		}
	}()

	return lb
}

func (lb *LeakyBucket) Allow(ctx context.Context, clientID string) (bool, int64, time.Time, error) {
	select {
	case lb.queue <- struct{}{}:
		return true, int64(cap(lb.queue) - len(lb.queue)), time.Now(), nil
	default:
		return false, 0, time.Now(), nil
	}
}

func (lb *LeakyBucket) Close() error {
	close(lb.done)
	return nil
}
