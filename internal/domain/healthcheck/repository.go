package domainhealthcheck

import (
	"context"
	"time"
)

type HealthCheckRepositoryDatastore interface {
	PingDatabase(ctx context.Context) (responseTime time.Duration, err error)
}
