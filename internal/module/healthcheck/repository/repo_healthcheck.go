package healthcheckrepository

import "erp-directory-service/internal/infrastructure"

type repository struct {
	db infrastructure.DB
}

func NewRepository(db infrastructure.DB) *repository {
	return &repository{
		db: db,
	}
}
