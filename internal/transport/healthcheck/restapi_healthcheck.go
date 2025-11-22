package transporthealthcheck

import (
	domainhealthcheck "go-bootstrap/internal/domain/healthcheck"
	"go-bootstrap/internal/gen/restapigen"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransportHealthCheckRestApi struct {
	healthcheckService domainhealthcheck.HealthCheckService
}

func NewTransportRestApi(
	healthcheckService domainhealthcheck.HealthCheckService,
) *TransportHealthCheckRestApi {
	return &TransportHealthCheckRestApi{
		healthcheckService: healthcheckService,
	}
}

func (t *TransportHealthCheckRestApi) ApiV1GetHealthCheck(c *gin.Context) {
	outputHealthcheck := t.healthcheckService.CheckDependencies(c.Request.Context())

	resp := restapigen.ApiV1GetHealthCheckResponse{
		Dependencies: restapigen.ApiV1GetHealthCheckResponseDependencies{
			Database: &restapigen.ApiV1GetHealthCheckResponseDependency{
				Message:      outputHealthcheck.Database.Message,
				ResponseTime: outputHealthcheck.Database.ResponseTime.String(),
				Status:       restapigen.ApiV1GetHealthCheckResponseDependencyStatus(outputHealthcheck.Database.Status),
			},
		},
		Status:    restapigen.ApiV1GetHealthCheckResponseStatus(outputHealthcheck.Status),
		Timestamp: outputHealthcheck.Timestamp,
	}

	c.JSON(http.StatusOK, resp)
}
