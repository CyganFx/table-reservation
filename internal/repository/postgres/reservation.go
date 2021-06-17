package postgres

import (
	"context"
	"fmt"
	"github.com/CyganFx/table-reservation/internal/domain"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"time"
)

type reservation struct {
	db *pgxpool.Pool
}

func NewReservation(db *pgxpool.Pool) *reservation {
	return &reservation{db: db}
}

func (r *reservation) GetAvailableLocationsByCafeID(cafeID int) ([]domain.Location, error) {
	query := `SELECT DISTINCT t.location_id, l.name FROM locations l
				JOIN tables t on t.location_id = l.id WHERE t.cafe_id = $1`
	rows, err := r.db.Query(context.Background(), query, cafeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ll []domain.Location

	for rows.Next() {
		l := locationsPool.Get().(*domain.Location)
		err = rows.Scan(&l.ID, &l.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to assign values to location struct from row %v", err)
		}
		ll = append(ll, *l)

		*l = domain.Location{}
		locationsPool.Put(l)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ll, nil
}

func (r *reservation) GetAvailableEventsByCafeID(cafeID int) ([]domain.Event, error) {
	query := `SELECT DISTINCT e.id, e.name FROM events e
				JOIN cafes_events ce on e.id = ce.event_id WHERE ce.cafe_id = $1`
	rows, err := r.db.Query(context.Background(), query, cafeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ee []domain.Event

	for rows.Next() {
		e := eventsPool.Get().(*domain.Event)
		err = rows.Scan(&e.ID, &e.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to assign values to event struct from row %v", err)
		}
		ee = append(ee, *e)

		*e = domain.Event{}
		eventsPool.Put(e)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ee, nil
}

func (r *reservation) GetSuitableTables(cafeID, partySize, locationID int, date, minPossibleBookingTime, maxPossibleBookingTime string) ([]domain.Table, error) {
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

	var tt []domain.Table

	for rows.Next() {
		t := tablesPool.Get().(*domain.Table)
		err = rows.Scan(&t.ID, &t.Capacity, &t.Location.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to assign values to table struct from row %v", err)
		}
		tt = append(tt, *t)

		*t = domain.Table{}
		tablesPool.Put(t)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tt, nil
}

func (r *reservation) BookTable(reservation *domain.Reservation) error {
	query := `INSERT INTO reservations(cafe_id, user_id, table_id, event_id, event_description,
				cust_name, cust_mobile, cust_email, num_of_persons, date, notify_date)
				VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id;`

	err := r.db.QueryRow(context.Background(), query, reservation.Cafe.ID, newNullInt(int32(reservation.User.ID)), reservation.Table.ID, reservation.Event.ID,
		newNullString(reservation.EventDescription), reservation.CustName, reservation.CustMobile,
		reservation.CustEmail, reservation.PartySize, reservation.Date, reservation.NotifyDate).
		Scan(&reservation.ID)
	if err != nil {
		return fmt.Errorf("failed to book table: %v", err)
	}

	return nil
}

func (r *reservation) GetUserReservations(userID int) ([]domain.Reservation, error) {
	query := `SELECT DISTINCT r.id, r.cafe_id, c.name, r.table_id,
			t.location_id, l.name, r.event_id, e.name, r.num_of_persons, r.date
			from reservations r
			left join cafes c on r.cafe_id = c.id
			left join tables t on r.table_id = t.id
			left join locations l on t.location_id = l.id
			left join events e on r.event_id = e.id
			WHERE r.user_id = $1;`

	rows, err := r.db.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rr []domain.Reservation

	for rows.Next() {
		r := reservationsPool.Get().(*domain.Reservation)
		err = rows.Scan(&r.ID, &r.Cafe.ID, &r.Cafe.Name, &r.Table.ID, &r.Table.Location.ID,
			&r.Table.Location.Name, &r.Event.ID, &r.Event.Name, &r.PartySize, &r.Date)
		if err != nil {
			return nil, fmt.Errorf("failed to assign values to Reservation struct from row: %v", err)
		}
		rr = append(rr, *r)

		*r = domain.Reservation{}
		reservationsPool.Put(r)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return rr, nil
}

func (r *reservation) GetReservationsByNotifyDate(now time.Time) ([]domain.Reservation, error) {
	query := `SELECT r.id, r.cafe_id, c.name, r.table_id,
			t.location_id, l.name, r.event_id, e.name, r.num_of_persons, r.date, r.cust_name, r.cust_email, r.notify_date
			from reservations r
			join cafes c on r.cafe_id = c.id
			join tables t on r.table_id = t.id
			join locations l on t.location_id = l.id
			join events e on r.event_id = e.id
			WHERE r.notify_date = $1;`

	rows, err := r.db.Query(context.Background(), query, now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rr []domain.Reservation

	for rows.Next() {
		r := reservationsPool.Get().(*domain.Reservation)
		err = rows.Scan(&r.ID, &r.Cafe.ID, &r.Cafe.Name, &r.Table.ID, &r.Table.Location.ID,
			&r.Table.Location.Name, &r.Event.ID, &r.Event.Name, &r.PartySize, &r.Date, &r.CustName, &r.CustEmail, &r.NotifyDate)
		if err != nil {
			return nil, fmt.Errorf("failed to assign values to Reservation struct from row: %v", err)
		}
		rr = append(rr, *r)

		*r = domain.Reservation{}
		reservationsPool.Put(r)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return rr, nil
}

func (r *reservation) GetBusyTables(cafeID, partySize, locationID int, date, minPossibleBookingTime, maxPossibleBookingTime string) ([]domain.Table, error) {
	query := `	select id, capacity, location_id
				from tables
				where cafe_id = $1
				  and capacity >= $2
				  and location_id = $3
				  and id in (
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

	var tt []domain.Table

	for rows.Next() {
		t := tablesPool.Get().(*domain.Table)
		err = rows.Scan(&t.ID, &t.Capacity, &t.Location.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to assign values to table struct from row %v", err)
		}
		tt = append(tt, *t)

		*t = domain.Table{}
		tablesPool.Put(t)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tt, nil
}

func (r *reservation) FreeTable(reservation *domain.Reservation, minBookDate, maxBookDate time.Time) error {
	query := `DELETE FROM reservations where cafe_id = $1 and table_id = $2 and date between $3 and $4;`

	_, err := r.db.Exec(context.Background(), query, reservation.Cafe.ID, reservation.Table.ID, minBookDate, maxBookDate)
	if err != nil {
		return errors.Wrap(err, "removing table from reservation")
	}

	return nil
}
