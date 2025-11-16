package domainhealthcheck

import "context"

type HealthCheckService interface {
	CheckDependencies(ctx context.Context) (output CheckDependenciesOutput)
}
