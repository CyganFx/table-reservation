package service

import (
	"github.com/CyganFx/table-reservation/ez-booking/pkg/domain"
	"strconv"
	"time"
)

const reservationInterval = 45

type reservation struct {
	repo ReservationRepo
}

type ReservationRepo interface {
	GetSuitableTables(cafeID, partySize, locationID int, date, minPossibleBookingTime, maxPossibleBookingTime string) ([]*domain.Table, error)
	GetAvailableLocationsByCafeID(cafeID int) ([]*domain.Location, error)
}

func NewReservation(repo ReservationRepo) *reservation {
	return &reservation{repo: repo}
}

func (r *reservation) GetLocationsByCafeID(cafeID int) ([]*domain.Location, error) {
	return r.repo.GetAvailableLocationsByCafeID(cafeID)
}

func (r *reservation) GetAvailableTables(cafeID, partySize, locationID int, date, bookTime string) ([]*domain.Table, error) {
	tempTime, _ := time.Parse("15:04", bookTime)

	tempMaxPossibleBookingTime := tempTime.Add(reservationInterval * time.Minute)
	tempMinPossibleBookingTime := tempTime.Add(-reservationInterval * time.Minute)

	maxPossibleBookingTime :=
		strconv.Itoa(tempMaxPossibleBookingTime.Hour()) + ":" + strconv.Itoa(tempMaxPossibleBookingTime.Minute())

	minPossibleBookingTime :=
		strconv.Itoa(tempMinPossibleBookingTime.Hour()) + ":" + strconv.Itoa(tempMinPossibleBookingTime.Minute())

	return r.repo.GetSuitableTables(cafeID, partySize, locationID, date, minPossibleBookingTime, maxPossibleBookingTime)
}
