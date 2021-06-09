package service

import "github.com/CyganFx/table-reservation/internal/domain"

type cafe struct {
	repo CafeRepo
}

func NewCafe(repo CafeRepo) *cafe {
	return &cafe{repo: repo}
}

type CafeRepo interface {
	FindLocations() ([]domain.Location, error)
	FindTypes() ([]domain.Type, error)
	FindEvents() ([]domain.Event, error)
	FindCities() ([]domain.City, error)
	FindCafes() ([]domain.Cafe, error)
	Insert(cafe *domain.Cafe) error
	SetLocationsByCafeID(cafeID int, locations []string) error
	SetEventsByCafeID(cafeID int, events []string) error
	SetTablesByCafeID(cafeID, locationID, numOfTables, capacity int) error
}

func (c *cafe) GetLocations() ([]domain.Location, error) {
	return c.repo.FindLocations()
}

func (c *cafe) GetTypes() ([]domain.Type, error) {
	return c.repo.FindTypes()
}

func (c *cafe) GetEvents() ([]domain.Event, error) {
	return c.repo.FindEvents()
}

func (c *cafe) GetCities() ([]domain.City, error) {
	return c.repo.FindCities()
}

func (c *cafe) CreateCafe(cafe *domain.Cafe) error {
	return c.repo.Insert(cafe)
}

func (c *cafe) SetLocations(cafeID int, locations []string) error {
	return c.repo.SetLocationsByCafeID(cafeID, locations)
}

func (c *cafe) SetEvents(cafeID int, events []string) error {
	return c.repo.SetEventsByCafeID(cafeID, events)
}

func (c *cafe) SetTables(cafeID, locationID, numOfTables, capacity int) error {
	return c.repo.SetTablesByCafeID(cafeID, locationID, numOfTables, capacity)
}

func (c *cafe) GetCafes() ([]domain.Cafe, error) {
	return c.repo.FindCafes()
}
