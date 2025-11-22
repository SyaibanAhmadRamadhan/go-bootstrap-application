package healthcheckservice

import (
	"context"
	domainhealthcheck "go-bootstrap/internal/domain/healthcheck"
	"time"
)

type service struct {
	healthcheckRepo domainhealthcheck.HealthCheckRepositoryDatastore
}

func NewService(
	healthcheckRepo domainhealthcheck.HealthCheckRepositoryDatastore,
) *service {
	return &service{
		healthcheckRepo: healthcheckRepo,
	}
}

func (s *service) CheckDependencies(ctx context.Context) (output domainhealthcheck.CheckDependenciesOutput) {
	database := domainhealthcheck.CheckDependency{}

	{
		database = domainhealthcheck.CheckDependency{
			Status:  domainhealthcheck.StatusDependencyOk,
			Message: "Ping Database Successfully",
		}
		responseTimePingDatabase, err := s.healthcheckRepo.PingDatabase(ctx)
		if err != nil {
			database.Message = err.Error()
			database.Status = domainhealthcheck.StatusDependencyError
		}
		database.ResponseTime = responseTimePingDatabase
	}

	return domainhealthcheck.CheckDependenciesOutput{
		Status: domainhealthcheck.NewStatusHealthCheck(
			database.Status,
		),
		Timestamp: time.Now().UTC(),
		Database:  database,
	}
}
