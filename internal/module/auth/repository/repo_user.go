package authrepository

import "go-bootstrap/internal/infrastructure"

type userRepository struct {
	db infrastructure.DB
}

func NewUserRepository(db infrastructure.DB) *repository {
	return &repository{
		db: db,
	}
}
