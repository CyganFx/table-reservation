package postgres

import (
	"context"
	"github.com/CyganFx/table-reservation/internal/domain"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

type cafe struct {
	db *pgxpool.Pool
}

func NewCafe(db *pgxpool.Pool) *cafe {
	return &cafe{db: db}
}

func (c *cafe) FindLocations() ([]domain.Location, error) {
	query := `SELECT * FROM locations`

	rows, err := c.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ll []domain.Location

	for rows.Next() {
		l := locationsPool.Get().(*domain.Location)
		err = rows.Scan(&l.ID, &l.Name)
		if err != nil {
			return nil, errors.Wrap(err, "failed to assign values to location struct from row")
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

func (c *cafe) FindTypes() ([]domain.Type, error) {
	query := `SELECT * FROM types`

	rows, err := c.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tt []domain.Type

	for rows.Next() {
		t := typesPool.Get().(*domain.Type)
		err = rows.Scan(&t.ID, &t.Name)
		if err != nil {
			return nil, errors.Wrap(err, "failed to assign values to type struct from row")
		}
		tt = append(tt, *t)

		*t = domain.Type{}
		typesPool.Put(t)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tt, nil
}

func (c *cafe) FindEvents() ([]domain.Event, error) {
	query := `SELECT * FROM events`

	rows, err := c.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ee []domain.Event

	for rows.Next() {
		e := eventsPool.Get().(*domain.Event)
		err = rows.Scan(&e.ID, &e.Name)
		if err != nil {
			return nil, errors.Wrap(err, "failed to assign values to type struct from row")
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

func (c *cafe) FindCities() ([]domain.City, error) {
	query := `SELECT * FROM cities`

	rows, err := c.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cc []domain.City

	for rows.Next() {
		c := citiesPool.Get().(*domain.City)
		err = rows.Scan(&c.ID, &c.Name)
		if err != nil {
			return nil, errors.Wrap(err, "failed to assign values to type struct from row")
		}
		cc = append(cc, *c)

		*c = domain.City{}
		citiesPool.Put(c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cc, nil
}

func (c *cafe) FindCafes() ([]domain.Cafe, error) {
	query := `SELECT c.id, c.name, c.city_id, city.name, c.type_id, type.name, c.address, c.mobile, c.email, c.created, c.image
			FROM cafes as c join types as type on c.type_id = type.id join cities as city on c.city_id = city.id`

	rows, err := c.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cc []domain.Cafe

	for rows.Next() {
		c := cafesPool.Get().(*domain.Cafe)
		err = rows.Scan(&c.ID, &c.Name, &c.City.ID, &c.City.Name, &c.Type.ID, &c.Type.Name, &c.Address, &c.Mobile, &c.Email, &c.Created, &c.ImageURL)
		if err != nil {
			return nil, errors.Wrap(err, "failed to assign values to type struct from row")
		}
		cc = append(cc, *c)

		*c = domain.Cafe{}
		cafesPool.Put(c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cc, nil
}

func (c *cafe) FindCafesFiltered(typeID, cityID int) ([]domain.Cafe, error) {
	query := `SELECT c.id, c.name, c.city_id, city.name, c.type_id, type.name, c.address, c.mobile, c.email, c.created, c.image
			FROM cafes as c join types as type on c.type_id = type.id join cities as city on c.city_id = city.id where type_id = $1 and city_id = $2`

	rows, err := c.db.Query(context.Background(), query, typeID, cityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cc []domain.Cafe

	for rows.Next() {
		c := cafesPool.Get().(*domain.Cafe)
		err = rows.Scan(&c.ID, &c.Name, &c.City.ID, &c.City.Name, &c.Type.ID, &c.Type.Name, &c.Address, &c.Mobile, &c.Email, &c.Created, &c.ImageURL)
		if err != nil {
			return nil, errors.Wrap(err, "failed to assign values to type struct from row")
		}
		cc = append(cc, *c)

		*c = domain.Cafe{}
		cafesPool.Put(c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cc, nil
}

func (c *cafe) Insert(cafe *domain.Cafe) error {
	query := `INSERT INTO cafes (name, city_id, type_id, address, mobile, email, created)
	VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err := c.db.QueryRow(context.Background(), query, cafe.Name, cafe.City.ID,
		cafe.Type.ID, cafe.Address, cafe.Mobile, cafe.Email, time.Now()).
		Scan(&cafe.ID)
	if err != nil {
		return errors.Wrap(err, "failed to insert cafe")
	}

	return nil
}

func (c *cafe) SetLocationsByCafeID(cafeID int, locations []string) error {
	query := `INSERT INTO cafes_locations (cafe_id, location_id)
	VALUES($1, $2)`
	for _, v := range locations {
		intV, err := strconv.Atoi(v)
		if err != nil {
			return errors.Wrap(err, "converting to int")
		}

		_, err = c.db.Exec(context.Background(), query, cafeID, intV)
		if err != nil {
			return errors.Wrap(err, "inserting to cafes_locations")
		}
	}

	return nil
}

func (c *cafe) SetEventsByCafeID(cafeID int, events []string) error {
	query := `INSERT INTO cafes_events (cafe_id, event_id)
	VALUES($1, $2)`
	for _, v := range events {
		intV, err := strconv.Atoi(v)
		if err != nil {
			return errors.Wrap(err, "converting to int")
		}

		_, err = c.db.Exec(context.Background(), query, cafeID, intV)
		if err != nil {
			return errors.Wrap(err, "inserting to cafes_events")
		}
	}

	return nil
}

func (c *cafe) SetTablesByCafeID(cafeID, locationID, numOfTables, capacity int) error {
	query := `INSERT INTO tables (cafe_id, capacity, location_id)
	VALUES($1, $2, $3)`
	for i := 0; i < numOfTables; i++ {
		_, err := c.db.Exec(context.Background(), query, cafeID, capacity, locationID)
		if err != nil {
			return errors.Wrap(err, "inserting to tables")
		}
	}

	return nil
}

func (c *cafe) SearchByName(name string) ([]domain.Cafe, error) {
	query := `SELECT c.id, c.name, c.city_id, city.name, c.type_id, type.name, c.address, c.mobile, c.email, c.created, c.image
			FROM cafes as c join types as type on c.type_id = type.id join cities as city on c.city_id = city.id
			where LOWER(c.name) like $1`

	rows, err := c.db.Query(context.Background(), query, "%"+name+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cc []domain.Cafe

	for rows.Next() {
		c := cafesPool.Get().(*domain.Cafe)
		err = rows.Scan(&c.ID, &c.Name, &c.City.ID, &c.City.Name, &c.Type.ID, &c.Type.Name, &c.Address, &c.Mobile, &c.Email, &c.Created, &c.ImageURL)
		if err != nil {
			return nil, errors.Wrap(err, "failed to assign values to type struct from row")
		}
		cc = append(cc, *c)

		*c = domain.Cafe{}
		cafesPool.Put(c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cc, nil
}
