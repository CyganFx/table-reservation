package service

import (
	http_v1 "github.com/CyganFx/table-reservation/ez-booking/internal/delivery/http-v1"
	"github.com/CyganFx/table-reservation/ez-booking/pkg/domain"
	"github.com/CyganFx/table-reservation/ez-booking/pkg/validator/forms"
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
	GetAvailableEventsByCafeID(cafeID int) ([]*domain.Event, error)
	BookTable(reservation *domain.Reservation) error
}

func NewReservation(repo ReservationRepo) *reservation {
	return &reservation{repo: repo}
}

func (r *reservation) GetLocationsByCafeID(cafeID int) ([]*domain.Location, error) {
	return r.repo.GetAvailableLocationsByCafeID(cafeID)
}

func (r *reservation) GetEventsByCafeID(cafeID int) ([]*domain.Event, error) {
	return r.repo.GetAvailableEventsByCafeID(cafeID)
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

func (r *reservation) BookTable(form *forms.FormValidator, userChoice http_v1.UserChoice) (int, *forms.FormValidator, error) {
	form.Required("name", "mobile", "email")
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("mobile", 11)
	form.MaxLength("mobile", 12)
	form.MaxLength("name", 50)
	form.MaxLength("email", 100)
	if !form.Valid() {
		return -1, form, nil
	}

	date := userChoice.Date
	bookTime := userChoice.BookTime + ":00"
	timeStampLikeDate := date + " " + bookTime

	reservation := domain.NewReservation()
	reservation.Date = timeStampLikeDate
	reservation.CustName = form.Get("name")
	reservation.CustMobile = form.Get("mobile")
	reservation.CustEmail = form.Get("email")
	reservation.Cafe.ID = userChoice.CafeID
	reservation.Table.ID = userChoice.TableID
	reservation.Event.ID = userChoice.EventID
	reservation.PartySize = userChoice.PartySize
	reservation.EventDescription = userChoice.EventDescription

	if err := r.repo.BookTable(reservation); err != nil {
		return -1, nil, err
	}

	return reservation.ID, nil, nil
}
