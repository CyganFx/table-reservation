package http_v1

import (
	"encoding/gob"
	"github.com/CyganFx/table-reservation/internal/domain"
	"github.com/CyganFx/table-reservation/pkg/validator/forms"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *handler) initReservationRoutes(api *gin.RouterGroup) {
	reservation := api.Group("/reservation")
	{
		reservation.GET("/cafe/:id", h.ReservationPage)
		reservation.POST("/tables", h.GetAvailableTables)
		reservation.POST("/confirm", h.Confirm)
		reservation.POST("/submit", h.BookTable)
	}
	gob.RegisterName("github.com/CyganFx/table-reservation/ez-booking/internal/delivery/http-v1.UserChoice",
		UserChoice{}) // in order to put struct in session
}

type ReservationService interface {
	GetAvailableTables(cafeID, partySize, locationID int, date, bookTime string) ([]domain.Table, error)
	GetLocationsByCafeID(cafeID int) ([]domain.Location, error)
	GetEventsByCafeID(cafeID int) ([]domain.Event, error)
	BookTable(form *forms.FormValidator, userChoice UserChoice, userID interface{}) (int, *forms.FormValidator, error)
	GetUserBookings(userID int) ([]domain.Reservation, error)
	SetDefaultReservationData(data *ReservationData, cafeID int) error
}

type ReservationData struct {
	CafeID            int
	CurrentDate       string
	MaxBookingDate    string
	CustName          string
	CustEmail         string
	CustMobile        string
	TimeSelector      []string
	PartySizeSelector []int
	LocationSelector  []domain.Location
	EventSelector     []domain.Event
	Tables            []domain.Table
	UserChoice        UserChoice
}

type UserChoice struct {
	CafeID           int
	TableID          int
	EventID          int
	LocationID       int
	PartySize        int
	EventDescription string
	Date             string
	BookTime         string
}

func (h *handler) ReservationPage(c *gin.Context) {
	cafeID, err := strconv.Atoi(c.Param("id"))
	if err != nil || cafeID < 1 {
		h.errors.NotFound(c)
		return
	}

	reservationData := &ReservationData{}
	err = h.reservationService.SetDefaultReservationData(reservationData, cafeID)
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
		CafeID:     cafeID,
		PartySize:  partySize,
		Date:       date,
		BookTime:   bookTime,
		LocationID: locationID,
	}

	session := sessions.Default(c)
	session.Set("userChoice", userChoice)
	session.Save()

	reservationData := &ReservationData{}
	reservationData.UserChoice = *userChoice
	reservationData.Tables, err = h.reservationService.GetAvailableTables(cafeID, partySize, locationID, date, bookTime)
	if err != nil {
		h.errors.ServerError(c, err)
		return
	}

	for idx, t := range reservationData.Tables {
		tempCapacityForHTML := make([]int, t.Capacity)
		reservationData.Tables[idx].CapacityForHTML = tempCapacityForHTML
	}

	err = h.reservationService.SetDefaultReservationData(reservationData, cafeID)
	if err != nil {
		h.errors.ServerError(c, err)
		return
	}
	reservationData.CurrentDate = date

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

	reservationData := &ReservationData{}
	form, err := h.userService.SetConfirmData(c, reservationData, tableID, eventID, eventDescription)
	if err != nil {
		h.errors.NotFound(c)
		return
	}

	h.render(c, "confirm.page.html", &templateData{
		ReservationData: reservationData,
		Form:            form,
	})
}

func (h *handler) BookTable(c *gin.Context) {
	if err := c.Request.ParseForm(); err != nil {
		h.errors.ClientError(c, http.StatusBadRequest)
		return
	}

	session := sessions.Default(c)
	userChoice := session.Get("userChoice").(UserChoice)
	userID := session.Get("authenticatedUserID")

	reservationData := &ReservationData{}
	reservationData.UserChoice = userChoice

	form := forms.New(c.Request.PostForm)

	reservationID, formValidator, err := h.reservationService.BookTable(form, userChoice, userID)
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

	session.Delete("userChoice")
	session.Set("reservationID", reservationID)
	session.Set("flash", "Booked successfully! You will get notifications as your time comes")
	session.Save()

	h.MainPage(c)
}
