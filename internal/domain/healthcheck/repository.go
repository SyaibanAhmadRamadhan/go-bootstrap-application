//go:generate go tool mockgen -source=repository.go -destination=../../gen/mockgen/healthcheck_repository_mock.gen.go -package=mockgen

package domainhealthcheck

import (
	"context"
	"time"
)

type HealthCheckRepositoryDatastore interface {
	PingDatabase(ctx context.Context) (responseTime time.Duration, err error)
}
