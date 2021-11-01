package controller

import (
	"github.com/DieGopherLT/refaccionaria-backend/internal/repository"
)

var Repo *Repository

type Repository struct {
	db repository.DatabaseRepo
}

func NewRepo(db repository.DatabaseRepo) *Repository {
	return &Repository{
		db: db,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}
