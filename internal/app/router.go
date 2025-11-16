package app

import (
	transporthealthcheck "erp-directory-service/internal/transport/healthcheck"
)

type router struct {
	// restapigen.Unimplemented
	*transporthealthcheck.TransportHealthCheckRestApi
}
