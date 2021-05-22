package service

import (
	"fmt"
	http_v1 "github.com/CyganFx/table-reservation/ez-booking/internal/delivery/http-v1"
	"github.com/CyganFx/table-reservation/ez-booking/internal/domain"
	"github.com/CyganFx/table-reservation/ez-booking/pkg/validator/forms"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

const (
	reservationInterval          = 90
	maxDaysToBookInAdvance       = 7
	amountOfHoursRestaurantWorks = 15 // for now let every restaurant work from 11:00 to 2:00
	possibleBookingsPerHour      = 4
	restaurantOpeningTimeHours   = 11
	maxTableCapacity             = 8 //assume that max table size is 8 for all restaurants
	bookTimeSelectInterval       = 15
)

type reservation struct {
	repo ReservationRepo
}

type ReservationRepo interface {
	GetSuitableTables(cafeID, partySize, locationID int, date, minPossibleBookingTime, maxPossibleBookingTime string) ([]*domain.Table, error)
	GetAvailableLocationsByCafeID(cafeID int) ([]*domain.Location, error)
	GetAvailableEventsByCafeID(cafeID int) ([]*domain.Event, error)
	BookTable(reservation *domain.Reservation) error
	GetUserReservations(userID int) ([]*domain.Reservation, error)
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

	// postgres between operation is inclusive for greater value, therefore we need to extract to fit in range
	tempMaxPossibleBookingTime := tempTime.Add(reservationInterval - 1*time.Minute)
	tempMinPossibleBookingTime := tempTime.Add(-reservationInterval * time.Minute)

	maxPossibleBookingTime :=
		strconv.Itoa(tempMaxPossibleBookingTime.Hour()) + ":" + strconv.Itoa(tempMaxPossibleBookingTime.Minute())

	minPossibleBookingTime :=
		strconv.Itoa(tempMinPossibleBookingTime.Hour()) + ":" + strconv.Itoa(tempMinPossibleBookingTime.Minute())

	return r.repo.GetSuitableTables(cafeID, partySize, locationID, date, minPossibleBookingTime, maxPossibleBookingTime)
}

func (r *reservation) BookTable(form *forms.FormValidator, userChoice http_v1.UserChoice, userID interface{}) (int, *forms.FormValidator, error) {
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
	if userID != nil {
		reservation.User.ID = userID.(int)
	} else {
		reservation.User.ID = -1
	}
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

func (r *reservation) GetUserBookings(userID int) ([]*domain.Reservation, error) {
	rr, err := r.repo.GetUserReservations(userID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	for _, r := range rr {
		if r.TimeStampDate.Year() == now.Year() &&
			int(r.TimeStampDate.Month()) == int(now.Month()) &&
			r.TimeStampDate.Day() == now.Day() {

			if r.TimeStampDate.Hour() > now.Hour() {
				r.IsActive = true
				r.HoursUntilReservation = r.TimeStampDate.Hour() - now.Hour()
				r.MinutesUntilReservation = r.TimeStampDate.Minute() - now.Minute()
			}
		}
	}

	return rr, nil
}

func (r *reservation) SetDefaultReservationData(data *http_v1.ReservationData, cafeID int) error {
	data.CafeID = cafeID
	data.CurrentDate = time.Now().Format("2006-01-02")
	data.MaxBookingDate =
		time.Now().AddDate(0, 0, maxDaysToBookInAdvance).Format("2006-01-02")

	if err := r.setTimeSelector(data); err != nil {
		return err
	}
	r.setPartySizeSelector(data)
	err := r.setLocationSelector(data, cafeID)
	if err != nil {
		return err
	}
	err = r.setEventSelector(data, cafeID)
	if err != nil {
		return err
	}
	return nil
}

func (r *reservation) setTimeSelector(data *http_v1.ReservationData) error {
	hours := restaurantOpeningTimeHours
	minutes := 0
	addZeroToMinutes := "0"
	addZeroToHours := ""

	now := time.Now()

	if data.UserChoice.Date == now.Format("2006-01-02") &&
		now.Hour()-restaurantOpeningTimeHours >= 0 {
		var err error

		hours, err = strconv.Atoi(data.UserChoice.BookTime[:2])
		if err != nil {
			return errors.Wrap(err, "Failed to get bookTime hours")
		}
		if hours == 0 {
			addZeroToHours = "0"
		}

		minutes, err = strconv.Atoi(data.UserChoice.BookTime[3:5])
		if err != nil {
			return errors.Wrap(err, "Failed to get bookTime minutes")
		}
		if minutes != 0 {
			addZeroToMinutes = ""
		}
	}

	var possibleBookingIntervalsLength int
	considerMinutes := 0

	if minutes != 0 {
		considerMinutes = minutes / bookTimeSelectInterval
	}

	if hours == restaurantOpeningTimeHours {
		possibleBookingIntervalsLength = amountOfHoursRestaurantWorks*possibleBookingsPerHour - considerMinutes - (reservationInterval / bookTimeSelectInterval)
	} else {
		possibleBookingIntervalsLength = (amountOfHoursRestaurantWorks-(hours-restaurantOpeningTimeHours))*possibleBookingsPerHour - considerMinutes - (reservationInterval / bookTimeSelectInterval)
	}

	for i := 0; i <= possibleBookingIntervalsLength; i++ {
		if minutes == 60 {
			minutes = 0
			hours++
			addZeroToMinutes = "0"
		}
		if hours == 24 {
			hours = 0
			addZeroToHours = "0"
		}

		data.TimeSelector = append(data.TimeSelector,
			fmt.Sprintf("%s%d:%d%s", addZeroToHours, hours, minutes, addZeroToMinutes))

		addZeroToMinutes = ""
		minutes += bookTimeSelectInterval
	}

	return nil
}

func (r *reservation) setPartySizeSelector(data *http_v1.ReservationData) {
	for partySize := 1; partySize <= maxTableCapacity; partySize++ {
		data.PartySizeSelector = append(data.PartySizeSelector, partySize)
	}
}

func (r *reservation) setLocationSelector(data *http_v1.ReservationData, cafeID int) error {
	var err error
	data.LocationSelector, err = r.GetLocationsByCafeID(cafeID)
	if err != nil {
		return fmt.Errorf("trouble getting locations %v", err)
	}
	return nil
}

func (r *reservation) setEventSelector(data *http_v1.ReservationData, cafeID int) error {
	var err error
	data.EventSelector, err = r.GetEventsByCafeID(cafeID)
	if err != nil {
		return fmt.Errorf("trouble getting events %v", err)
	}
	return nil
}
