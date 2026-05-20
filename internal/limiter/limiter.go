package limiter

import (
	"context"
	"time"
)

type Limiter interface {
	Allow(ctx context.Context, clientID string) (bool, int64, time.Time, error)
	Close() error
}
