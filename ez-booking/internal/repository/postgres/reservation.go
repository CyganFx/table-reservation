package postgres

import (
	"context"
	"fmt"
	"github.com/CyganFx/table-reservation/ez-booking/pkg/domain"
	"github.com/jackc/pgx/v4/pgxpool"
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

	var ll []*domain.Location

	for rows.Next() {
		l := &domain.Location{}
		err = rows.Scan(&l.ID, &l.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to assign values to location struct from row %v", err)
		}
		ll = append(ll, l)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ll, nil
}

func (r *reservation) GetAvailableEventsByCafeID(cafeID int) ([]*domain.Event, error) {
	query := `SELECT DISTINCT e.id, e.name FROM events e
				JOIN cafes_events ce on e.id = ce.event_id WHERE ce.cafe_id = $1`
	rows, err := r.db.Query(context.Background(), query, cafeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ee []*domain.Event

	for rows.Next() {
		e := &domain.Event{}
		err = rows.Scan(&e.ID, &e.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to assign values to event struct from row %v", err)
		}
		ee = append(ee, e)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ee, nil
}

func (r *reservation) GetSuitableTables(cafeID, partySize, locationID int, date, minPossibleBookingTime, maxPossibleBookingTime string) ([]*domain.Table, error) {
	query := `	select id, capacity, location_id
				from tables
				where cafe_id = $1
				  and capacity >= $2
				  and location_id = $3
				  and id not in (
					select table_id
					from reservations
					where to_char(date, 'YYYY-MM-DD') = $4
					  and to_char(date, 'HH24:MI') between $5 and $6
					);`
	rows, err := r.db.Query(context.Background(),
		query, cafeID, partySize, locationID, date, minPossibleBookingTime, maxPossibleBookingTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tt []*domain.Table

	for rows.Next() {
		t := &domain.Table{}
		err = rows.Scan(&t.ID, &t.Capacity, &t.LocationID)
		if err != nil {
			return nil, fmt.Errorf("failed to assign values to table struct from row %v", err)
		}
		tt = append(tt, t)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tt, nil
}

func (r *reservation) BookTable(reservation *domain.Reservation) error {
	query := `INSERT INTO reservations(cafe_id, table_id, event_id, event_description,
				cust_name, cust_mobile, cust_email, num_of_persons, date)
				VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id;`
	err := r.db.QueryRow(context.Background(), query, reservation.Cafe.ID, reservation.Table.ID, reservation.Event.ID,
		reservation.EventDescription, reservation.CustName, reservation.CustMobile,
		reservation.CustEmail, reservation.PartySize, reservation.Date).
		Scan(&reservation.ID)
	if err != nil {
		return fmt.Errorf("failed to book table: %v", err)
	}

	return nil
}
