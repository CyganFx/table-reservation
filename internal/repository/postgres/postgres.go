package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/CyganFx/table-reservation/internal/app/config"
	"github.com/CyganFx/table-reservation/internal/domain"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"sync"
)

var (
	tablesPool = sync.Pool{
		New: func() interface{} {
			return domain.NewTable()
		},
	}
	eventsPool = sync.Pool{
		New: func() interface{} {
			return &domain.Event{}
		},
	}
	locationsPool = sync.Pool{
		New: func() interface{} {
			return &domain.Location{}
		},
	}
	cafesPool = sync.Pool{
		New: func() interface{} {
			return &domain.Cafe{}
		},
	}
	typesPool = sync.Pool{
		New: func() interface{} {
			return &domain.Type{}
		},
	}
	citiesPool = sync.Pool{
		New: func() interface{} {
			return &domain.City{}
		},
	}
	reservationsPool = sync.Pool{
		New: func() interface{} {
			return domain.NewReservation()
		},
	}
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

func newNullInt(n int32) sql.NullInt32 {
	if n == -1 {
		return sql.NullInt32{}
	}
	return sql.NullInt32{
		Int32: n,
		Valid: true,
	}
}

func newNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
