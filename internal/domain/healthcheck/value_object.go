package domainhealthcheck

type StatusHealthCheck string

func NewStatusHealthCheck(dependencyStatus ...StatusDependency) (output StatusHealthCheck) {
	if len(dependencyStatus) == 0 {
		return "degraded"
	}

	output = "healthy"

	errorStatus := make([]StatusDependency, 0, len(dependencyStatus))
	for _, v := range dependencyStatus {
		if v == StatusDependencyError {
			if output != "degraded" {
				output = "degraded"
			}

			errorStatus = append(errorStatus, v)
		}
	}

	if len(errorStatus) == len(dependencyStatus) {
		output = "unhealthy"
	}

	return output
}

type StatusDependency string

const (
	StatusDependencyOk    StatusDependency = "ok"
	StatusDependencyError StatusDependency = "error"
)
