package service

import (
	"fmt"
	http_v1 "github.com/CyganFx/table-reservation/internal/delivery/http-v1"
	"github.com/CyganFx/table-reservation/internal/domain"
	"github.com/CyganFx/table-reservation/pkg/validator/forms"
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

	dateLayout          = "2006-01-02"
	timeLayout          = "15:04:05"
	minutesBeforeNotify = 60
)

type reservation struct {
	repo ReservationRepo
}

func NewReservation(repo ReservationRepo) *reservation {
	return &reservation{repo: repo}
}

type ReservationRepo interface {
	GetSuitableTables(cafeID, partySize, locationID int, date, minPossibleBookingTime, maxPossibleBookingTime string) ([]domain.Table, error)
	GetBusyTables(cafeID, partySize, locationID int, date, minPossibleBookingTime, maxPossibleBookingTime string) ([]domain.Table, error)
	GetAvailableLocationsByCafeID(cafeID int) ([]domain.Location, error)
	GetAvailableEventsByCafeID(cafeID int) ([]domain.Event, error)
	BookTable(reservation *domain.Reservation) error
	GetUserReservations(userID int) ([]domain.Reservation, error)
	GetReservationsByNotifyDate(now time.Time) ([]domain.Reservation, error)
	FreeTable(reservation *domain.Reservation, minBookDate, maxBookDate time.Time) error
}

func (r *reservation) GetLocationsByCafeID(cafeID int) ([]domain.Location, error) {
	return r.repo.GetAvailableLocationsByCafeID(cafeID)
}

func (r *reservation) GetEventsByCafeID(cafeID int) ([]domain.Event, error) {
	return r.repo.GetAvailableEventsByCafeID(cafeID)
}

func (r *reservation) GetAvailableTables(cafeID, partySize, locationID int, date, bookTime string) ([]domain.Table, error) {
	tempTime, _ := time.Parse("15:04", bookTime)
	// postgres between operation is inclusive for greater value, therefore we need to extract to fit in range
	tempMaxPossibleBookingTime := tempTime.Add((reservationInterval - 1) * time.Minute)
	tempMinPossibleBookingTime := tempTime.Add(-(reservationInterval - 1) * time.Minute)

	maxPossibleBookingTime :=
		strconv.Itoa(tempMaxPossibleBookingTime.Hour()) + ":" + strconv.Itoa(tempMaxPossibleBookingTime.Minute())

	minPossibleBookingTime :=
		strconv.Itoa(tempMinPossibleBookingTime.Hour()) + ":" + strconv.Itoa(tempMinPossibleBookingTime.Minute())

	return r.repo.GetSuitableTables(cafeID, partySize, locationID, date, minPossibleBookingTime, maxPossibleBookingTime)
}

func (r *reservation) GetBusyTables(cafeID, partySize, locationID int, date, bookTime string) ([]domain.Table, error) {
	tempTime, _ := time.Parse("15:04", bookTime)
	// postgres between operation is inclusive for greater value, therefore we need to extract to fit in range
	tempMaxPossibleBookingTime := tempTime.Add((reservationInterval - 1) * time.Minute)
	tempMinPossibleBookingTime := tempTime.Add(-(reservationInterval - 1) * time.Minute)

	maxPossibleBookingTime :=
		strconv.Itoa(tempMaxPossibleBookingTime.Hour()) + ":" + strconv.Itoa(tempMaxPossibleBookingTime.Minute())

	minPossibleBookingTime :=
		strconv.Itoa(tempMinPossibleBookingTime.Hour()) + ":" + strconv.Itoa(tempMinPossibleBookingTime.Minute())

	return r.repo.GetBusyTables(cafeID, partySize, locationID, date, minPossibleBookingTime, maxPossibleBookingTime)
}

func setBookAndNotifyDate(strBookDate, strBookTime string, bookDate, notifyDate *time.Time) error {
	partialDate, err := time.Parse(dateLayout, strBookDate)
	if err != nil {
		return errors.Wrap(err, "failed to parse string date")
	}

	bookTime, err := time.Parse(timeLayout, strBookTime)
	if err != nil {
		return errors.Wrap(err, "failed to parse string bookTime")
	}

	*bookDate = partialDate.Add(time.Hour*time.Duration(bookTime.Hour()) +
		time.Minute*time.Duration(bookTime.Minute()) +
		time.Second*time.Duration(bookTime.Second()))

	*notifyDate = bookDate.Add(time.Minute * time.Duration(-minutesBeforeNotify))

	return nil
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

	var bookDate, notifyDate time.Time
	strBookTime := userChoice.BookTime + ":00"

	err := setBookAndNotifyDate(userChoice.Date, strBookTime, &bookDate, &notifyDate)
	if err != nil {
		return -1, form, err
	}

	reservation := domain.NewReservation()
	if userID != nil {
		reservation.User.ID = userID.(int)
	} else {
		reservation.User.ID = -1
	}
	reservation.Date = bookDate
	reservation.NotifyDate = notifyDate
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

func (r *reservation) BookTableManually(userChoice http_v1.UserChoice, userID interface{}) (int, error) {
	var bookDate, notifyDate time.Time
	strBookTime := userChoice.BookTime + ":00"

	err := setBookAndNotifyDate(userChoice.Date, strBookTime, &bookDate, &notifyDate)
	if err != nil {
		return -1, err
	}

	reservation := domain.NewReservation()
	if userID != nil {
		reservation.User.ID = userID.(int)
	} else {
		reservation.User.ID = -1
	}

	reservation.Date = bookDate
	reservation.Cafe.ID = userChoice.CafeID
	reservation.Table.ID = userChoice.TableID
	reservation.Event.ID = 1 // default value

	if err := r.repo.BookTable(reservation); err != nil {
		return -1, err
	}

	return reservation.ID, nil
}

func (r *reservation) FreeTableManually(userChoice http_v1.UserChoice, userID interface{}) error {
	var bookDate, notifyDate time.Time
	strBookTime := userChoice.BookTime + ":00"

	err := setBookAndNotifyDate(userChoice.Date, strBookTime, &bookDate, &notifyDate)
	if err != nil {
		return err
	}

	reservation := domain.NewReservation()
	if userID != nil {
		reservation.User.ID = userID.(int)
	} else {
		reservation.User.ID = -1
	}

	minBookDate := bookDate.Add(-(reservationInterval - 1) * time.Minute)
	maxBookDate := bookDate.Add((reservationInterval - 1) * time.Minute)

	reservation.Cafe.ID = userChoice.CafeID
	reservation.Table.ID = userChoice.TableID
	reservation.Event.ID = 1 // default value

	if err := r.repo.FreeTable(reservation, minBookDate, maxBookDate); err != nil {
		return err
	}

	return nil
}

func (r *reservation) GetUserBookings(userID int) ([]domain.Reservation, error) {
	rr, err := r.repo.GetUserReservations(userID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	for idx, r := range rr {
		if r.Date.Year() == now.Year() &&
			int(r.Date.Month()) == int(now.Month()) &&
			r.Date.Day() == now.Day() {
			if r.Date.Hour() > now.Hour() {
				// using rr[idx] because we have slice of values, not pointers
				rr[idx].IsActive = true
				rr[idx].HoursUntilReservation = r.Date.Hour() - now.Hour()
				rr[idx].MinutesUntilReservation = r.Date.Minute() - now.Minute()
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

	//now := time.Now()
	//
	//if now.Hour()-restaurantOpeningTimeHours >= 0 &&
	//	(data.UserChoice == http_v1.UserChoice{} ||
	//		data.UserChoice.Date == now.Format("2006-01-02")) {
	//	hours = now.Hour()
	//	if hours == 0 {
	//		addZeroToHours = "0"
	//	}
	//
	//	minutes = time.Now().Minute()
	//	if minutes != 0 {
	//		addZeroToMinutes = ""
	//		for minutes % 15 != 0 {
	//			if minutes == 60 {
	//				hours++
	//				if hours == 24 {
	//					hours = 0
	//					addZeroToHours = "0"
	//				}
	//				minutes = 0
	//				break
	//			}
	//			minutes++
	//		}
	//	}
	//}
	//
	//var possibleBookingIntervalsLength int
	//considerMinutes := 0
	//
	//if minutes != 0 {
	//	considerMinutes = minutes / bookTimeSelectInterval
	//}
	//
	//if hours == restaurantOpeningTimeHours {
	//	possibleBookingIntervalsLength = amountOfHoursRestaurantWorks*possibleBookingsPerHour - considerMinutes - (reservationInterval / bookTimeSelectInterval)
	//} else {
	//	possibleBookingIntervalsLength = (amountOfHoursRestaurantWorks-(hours-restaurantOpeningTimeHours))*possibleBookingsPerHour - considerMinutes - (reservationInterval / bookTimeSelectInterval)
	//}

	possibleBookingIntervalsLength := amountOfHoursRestaurantWorks*possibleBookingsPerHour - (reservationInterval / bookTimeSelectInterval)

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

func (r *reservation) CheckNotifyDate(now time.Time, notificator http_v1.NotificatorService) error {
	nowNoSeconds := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, time.UTC)
	reservations, err := r.repo.GetReservationsByNotifyDate(nowNoSeconds)
	if err != nil {
		return errors.Wrap(err, "getting emails")
	}
	if len(reservations) == 0 {
		return nil
	}

	if err = notificator.UsersBooking(reservations); err != nil {
		return err
	}

	return nil
}
