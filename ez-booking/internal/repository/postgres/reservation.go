package postgres

import (
	"context"
	"github.com/CyganFx/table-reservation/ez-booking/pkg/domain"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type reservation struct {
	db *pgxpool.Pool
}

func NewReservation(db *pgxpool.Pool) *reservation {
	return &reservation{db: db}
}

func (r *reservation) GetAvailableLocationsByCafeID(cafeID int) ([]*domain.Location, error) {
	query := `SELECT DISTINCT t.location_id, l.name FROM locations l
				JOIN tables t on t.location_id = l.id WHERE t.cafe_id = $1`
	rows, err := r.db.Query(context.Background(), query, cafeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []*domain.Location

	for rows.Next() {
		location := &domain.Location{}
		err = rows.Scan(
			&location.ID, &location.Name)
		if err != nil {
			return nil, err
		}
		locations = append(locations, location)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return locations, nil
}

func (r *reservation) GetSuitableTables(cafeID, capacity, locationID int, date time.Time, time string) ([]*domain.Table, error) {

	panic("implement me")
}
