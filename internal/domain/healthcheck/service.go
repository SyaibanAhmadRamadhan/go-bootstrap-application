//go:generate go tool mockgen -source=service.go -destination=../../gen/mockgen/healthcheck_service_mock.gen.go -package=mockgen

package domainhealthcheck

import "context"

type HealthCheckService interface {
	CheckDependencies(ctx context.Context) (output CheckDependenciesOutput)
}
