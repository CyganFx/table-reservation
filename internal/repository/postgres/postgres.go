package postgres

import (
	"context"
	"fmt"
	"github.com/CyganFx/table-reservation/internal/app/config"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

func InitPool(cfg config.Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.Database.Username, cfg.Database.Password, cfg.Database.Host,
		cfg.Database.Port, cfg.Database.DBName)

	dbPool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to db")
	}

	return dbPool, nil
}
