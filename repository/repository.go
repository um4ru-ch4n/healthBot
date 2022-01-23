package repository

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	psqlPool *pgxpool.Pool
}

func NewRepository(psqlPool *pgxpool.Pool) *Repository {
	return &Repository{
		psqlPool: psqlPool,
	}
}
