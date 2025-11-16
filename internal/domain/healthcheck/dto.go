package domainhealthcheck

import "time"

type CheckDependenciesOutput struct {
	Status    StatusHealthCheck
	Timestamp time.Time // Timestamp in utc
	Database  CheckDependency
}

type CheckDependency struct {
	Status       StatusDependency
	ResponseTime time.Duration
	Message      string
}
