package service

import (
	"github.com/CyganFx/table-reservation/ez-booking/pkg/domain"
	"time"
)

type reservation struct {
	repo ReservationRepo
}

type ReservationRepo interface {
	GetSuitableTables(cafeID, capacity, locationID int, date time.Time, time string) ([]*domain.Table, error)
	GetAvailableLocationsByCafeID(cafeID int) ([]*domain.Location, error)
}

func NewReservation(repo ReservationRepo) *reservation {
	return &reservation{repo: repo}
}

func (r *reservation) GetLocationsByCafeID(cafeID int) ([]*domain.Location, error) {
	return r.repo.GetAvailableLocationsByCafeID(cafeID)
}

func (r *reservation) GetAvailableTables(cafeID, capacity, locationID int, date, bookTime string) ([]*domain.Table, error) {
	panic("implement me")
	//currentDate := time.Now().Format("2006-01-02")
	//return r.repo.GetSuitableTables(cafeID, capacity, locationID, date, bookTime), nil
}
