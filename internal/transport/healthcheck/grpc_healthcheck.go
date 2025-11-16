package transporthealthcheck

import (
	"context"
	"erp-directory-service/gen/grpcgen/healthcheck"
	domainhealthcheck "erp-directory-service/internal/domain/healthcheck"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type TransportHealthCheckGrpc struct {
	healthcheckService domainhealthcheck.HealthCheckService
	healthcheck.UnimplementedHealthCheckServiceServer
}

func NewTransportGrpc(
	healthcheckService domainhealthcheck.HealthCheckService,
) *TransportHealthCheckGrpc {
	return &TransportHealthCheckGrpc{
		healthcheckService: healthcheckService,
	}
}

func (t *TransportHealthCheckGrpc) ApiV1HealthCheck(ctx context.Context, _ *healthcheck.ApiV1HealthCheckRequest) (*healthcheck.ApiV1HealthCheckResponse, error) {
	outputHealthcheck := t.healthcheckService.CheckDependencies(ctx)

	statusHealthCheck := healthcheck.ServiceStatus_HEALTH_CHECK_SERVICE_STATUS_HEALTHY
	switch outputHealthcheck.Status {
	case domainhealthcheck.StatusHealthCheckDegraded:
		statusHealthCheck = healthcheck.ServiceStatus_HEALTH_CHECK_SERVICE_STATUS_DEGRADED
	case domainhealthcheck.StatusHealthCheckUnhealthy:
		statusHealthCheck = healthcheck.ServiceStatus_HEALTH_CHECK_SERVICE_STATUS_UNHEALTHY
	}

	databaseStatus := healthcheck.DependencyStatus_HEALTH_CHECK_DEPENDENCY_STATUS_OK
	if outputHealthcheck.Database.Status == domainhealthcheck.StatusDependencyError {
		databaseStatus = healthcheck.DependencyStatus_HEALTH_CHECK_DEPENDENCY_STATUS_ERROR
	}

	resp := healthcheck.ApiV1HealthCheckResponse{
		Dependencies: &healthcheck.ApiV1HealthCheckDependencies{
			Database: &healthcheck.ApiV1HealthCheckDependency{
				Message:      outputHealthcheck.Database.Message,
				ResponseTime: outputHealthcheck.Database.ResponseTime.String(),
				Status:       databaseStatus,
			},
		},
		Status:    statusHealthCheck,
		Timestamp: timestamppb.New(outputHealthcheck.Timestamp),
	}

	return &resp, nil
}
