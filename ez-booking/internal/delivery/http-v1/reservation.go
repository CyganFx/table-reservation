package http_v1

import (
	"encoding/gob"
	"fmt"
	"github.com/CyganFx/table-reservation/ez-booking/pkg/domain"
	"github.com/CyganFx/table-reservation/ez-booking/pkg/validator/forms"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

const (
	maxDaysToBookInAdvance       = 7
	amountOfHoursRestaurantWorks = 15 // for now let every restaurant work from 11:00 to 2:00
	possibleBookingsPerHour      = 4
	restaurantOpeningTimeHours   = 11
	minutesInHour                = 60
	maxTableCapacity             = 8 //assume that max table size is 8 for all restaurants
)

func (h *handler) initReservationRoutes(api *gin.RouterGroup) {
	reservation := api.Group("/reservation")
	{
		reservation.GET("/cafe/:id", h.ReservationPage)
		reservation.POST("/tables", h.GetAvailableTables)
		reservation.POST("/confirm", h.Confirm)
		reservation.POST("/submit", h.BookTable)
	}
	gob.Register(UserChoice{}) // in order to put struct in session
}

type ReservationService interface {
	GetAvailableTables(cafeID, partySize, locationID int, date, bookTime string) ([]*domain.Table, error)
	GetLocationsByCafeID(cafeID int) ([]*domain.Location, error)
	GetEventsByCafeID(cafeID int) ([]*domain.Event, error)
	BookTable(form *forms.FormValidator, userChoice UserChoice) (int, *forms.FormValidator, error)
}

type ReservationData struct {
	CafeID            int
	CurrentDate       string
	MaxBookingDate    string
	TimeSelector      []string
	PartySizeSelector []int
	LocationSelector  []*domain.Location
	EventSelector     []*domain.Event
	Tables            []*domain.Table
	UserChoice        UserChoice
}

type UserChoice struct {
	CafeID           int
	TableID          int
	EventID          int
	EventDescription string
	PartySize        int
	Date             string
	BookTime         string
}

func (h *handler) setReservationData(data *ReservationData, cafeID int) error {
	data.CafeID = cafeID
	data.CurrentDate = time.Now().Format("2006-01-02")
	data.MaxBookingDate =
		time.Now().AddDate(0, 0, maxDaysToBookInAdvance).Format("2006-01-02")

	setTimeSelector(data)
	setPartySizeSelector(data)
	err := setLocationSelector(h, data, cafeID)
	if err != nil {
		return err
	}
	err = setEventSelector(h, data, cafeID)
	if err != nil {
		return err
	}
	return nil
}

func setTimeSelector(data *ReservationData) {
	hours := restaurantOpeningTimeHours
	minutes := 0
	addZeroToMinutes := "0"
	addZeroToHours := ""
	for i := 0; i <= amountOfHoursRestaurantWorks*possibleBookingsPerHour-possibleBookingsPerHour; i++ {
		if i%4 == 0 && i != 0 {
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
		minutes += minutesInHour / possibleBookingsPerHour
	}
}

func setPartySizeSelector(data *ReservationData) {
	for partySize := 1; partySize <= maxTableCapacity; partySize++ {
		data.PartySizeSelector = append(data.PartySizeSelector,
			partySize)
	}
}

func setLocationSelector(h *handler, data *ReservationData, cafeID int) error {
	var err error
	data.LocationSelector, err = h.reservationService.GetLocationsByCafeID(cafeID)
	if err != nil {
		return fmt.Errorf("trouble getting locations %v", err)
	}
	return nil
}

func setEventSelector(h *handler, data *ReservationData, cafeID int) error {
	var err error
	data.EventSelector, err = h.reservationService.GetEventsByCafeID(cafeID)
	if err != nil {
		return fmt.Errorf("trouble getting events %v", err)
	}
	return nil
}

func (h *handler) ReservationPage(c *gin.Context) {
	cafeID, err := strconv.Atoi(c.Param("id"))
	if err != nil || cafeID < 1 {
		h.errors.NotFound(c)
		return
	}

	reservationData := &ReservationData{}
	err = h.setReservationData(reservationData, cafeID)
	if err != nil {
		h.errors.ServerError(c, err)
		return
	}

	h.render(c, "reservation.page.html", &templateData{
		ReservationData: reservationData,
	})
}

func (h *handler) GetAvailableTables(c *gin.Context) {
	err := c.Request.ParseForm()
	if err != nil {
		h.errors.ClientError(c, http.StatusBadRequest)
		return
	}

	cafeID, _ := strconv.Atoi(c.Request.FormValue("cafe_id"))
	date := c.Request.FormValue("date")
	bookTime := c.Request.FormValue("bookTime")
	locationID, _ := strconv.Atoi(c.Request.FormValue("location_id"))
	partySize, _ := strconv.Atoi(c.Request.FormValue("party_size"))

	userChoice := &UserChoice{
		CafeID:    cafeID,
		PartySize: partySize,
		Date:      date,
		BookTime:  bookTime,
	}

	session := sessions.Default(c)
	session.Set("userChoice", userChoice)
	session.Save()

	reservationData := &ReservationData{}

	reservationData.Tables, err = h.reservationService.GetAvailableTables(cafeID, partySize, locationID, date, bookTime)
	if err != nil {
		h.errors.ServerError(c, err)
		return
	}

	for _, t := range reservationData.Tables {
		tempCapacityForHTML := make([]int, t.Capacity)
		t.CapacityForHTML = tempCapacityForHTML
	}

	err = h.setReservationData(reservationData, cafeID)
	if err != nil {
		h.errors.ServerError(c, err)
		return
	}

	h.render(c, "reservation.page.html", &templateData{
		ReservationData: reservationData,
	})
}

func (h *handler) Confirm(c *gin.Context) {
	if err := c.Request.ParseForm(); err != nil {
		h.errors.ClientError(c, http.StatusBadRequest)
		return
	}

	tableID, _ := strconv.Atoi(c.Request.FormValue("table_id"))
	eventID, _ := strconv.Atoi(c.Request.FormValue("event_id"))
	eventDescription := c.Request.FormValue("event_description")

	session := sessions.Default(c)
	userChoice := session.Get("userChoice").(UserChoice)
	userChoice.TableID = tableID
	userChoice.EventID = eventID
	userChoice.EventDescription = eventDescription
	session.Set("userChoice", userChoice)
	session.Save()

	reservationData := &ReservationData{}
	reservationData.UserChoice = userChoice

	h.render(c, "confirm.page.html", &templateData{
		ReservationData: reservationData,
		Form:            forms.New(nil),
	})
}

func (h *handler) BookTable(c *gin.Context) {
	if err := c.Request.ParseForm(); err != nil {
		h.errors.ClientError(c, http.StatusBadRequest)
		return
	}

	session := sessions.Default(c)
	userChoice := session.Get("userChoice").(UserChoice)
	reservationData := &ReservationData{}
	reservationData.UserChoice = userChoice

	form := forms.New(c.Request.PostForm)

	reservationID, formValidator, err := h.reservationService.BookTable(form, userChoice)
	if formValidator != nil {
		h.render(c, "confirm.page.html", &templateData{
			ReservationData: reservationData,
			Form:            form,
		})
		return
	} else if err != nil {
		h.errors.ServerError(c, err)
		return
	}

	session.Set("reservationID", reservationID)
	session.Set("flash", "Booked successfully! You will get notifications as your time comes")
	session.Save()

	h.MainPage(c)
}
