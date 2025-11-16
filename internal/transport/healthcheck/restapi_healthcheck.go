package transporthealthcheck

import (
	domainhealthcheck "erp-directory-service/internal/domain/healthcheck"
	"erp-directory-service/internal/gen/restapigen"
	"net/http"

	httpx "github.com/SyaibanAhmadRamadhan/go-foundation-kit/http"
	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/http/server/chix"
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

func (t *TransportHealthCheckRestApi) ApiV1GetHealthCheck(w http.ResponseWriter, r *http.Request) {
	outputHealthcheck := t.healthcheckService.CheckDependencies(r.Context())

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

	chix.Write(w, http.StatusOK, httpx.ContentTypeJSON, resp)
}
