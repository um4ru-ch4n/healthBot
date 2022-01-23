package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgreConfig struct {
	Host     string
	Port     string
	Username string
	DBName   string
	SSLMode  string
	Password string
}

func NewPostgreConnPool(cfg PostgreConfig) (*pgxpool.Pool, error) {
	url := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.SSLMode, cfg.Password)

	dbpool, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		return nil, fmt.Errorf("pg connect: %v", err)
	}

	if err = dbpool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("pg ping: %v", err)
	}
	return dbpool, nil
}
