package healthcheckrepository

import "erp-directory-service/internal/provider"

type repository struct {
	db provider.DB
}

func NewRepository(db provider.DB) *repository {
	return &repository{
		db: db,
	}
}
