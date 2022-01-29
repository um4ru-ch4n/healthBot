package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/um4aru-ch4n/healthBot/domain"
)

type Repository struct {
	psqlPool *pgxpool.Pool
}

func NewRepository(psqlPool *pgxpool.Pool) *Repository {
	return &Repository{
		psqlPool: psqlPool,
	}
}

func (r *Repository) GetChatInfoPoll(ctx context.Context) (chatInfo map[int64]*domain.ChatInfo, pollChat map[string]int64, err error) {
	tx, err := r.psqlPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, nil, fmt.Errorf("BeginTx: %w", err)
	}

	const query = `
	`

	rows, err := tx.Query(ctx, query)
	if err != nil {
		return nil, nil, fmt.Errorf("tx.Query: %w", err)
	}

	defer rows.Close()

	return nil, nil, nil
}
