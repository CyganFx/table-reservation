package http_v1

import (
	"fmt"
	"github.com/CyganFx/table-reservation/ez-booking/pkg/domain"
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
		reservation.POST("/tables", h.GetAvailableTables)
		reservation.GET("/cafe/:id", h.ReservationPage)
	}
}

type ReservationService interface {
	GetAvailableTables(cafeID, capacity, locationID int, date, bookTime string) ([]*domain.Table, error)
	GetLocationsByCafeID(cafeID int) ([]*domain.Location, error)
}

type ReservationData struct {
	CurrentDate       string
	MaxBookingDate    string
	TimeSelector      []string
	LocationSelector  []*domain.Location
	PartySizeSelector []int
}

func (h *handler) setReservationData(data *ReservationData, cafeID int) error {
	data.CurrentDate = time.Now().Format("2006-01-02")
	data.MaxBookingDate =
		time.Now().AddDate(0, 0, maxDaysToBookInAdvance).Format("2006-01-02")

	setTimeSelector(data)
	setPartySizeSelector(data)
	err := setLocationSelector(h, data, cafeID)
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
	for capacity := 1; capacity < maxTableCapacity; capacity++ {
		data.PartySizeSelector = append(data.PartySizeSelector,
			capacity)
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
	}

	fmt.Println(reservationData.CurrentDate)
	fmt.Println(reservationData.MaxBookingDate)
	fmt.Println(reservationData.PartySizeSelector)
	fmt.Println(reservationData.LocationSelector)
	fmt.Println(reservationData.TimeSelector)

	h.render(c, "reservation.page.html", &templateData{
		ReservationData: reservationData,
	})
}

func (h *handler) GetAvailableTables(c *gin.Context) {
	if err := c.Request.ParseForm(); err != nil {
		h.errors.ClientError(c, http.StatusBadRequest)
		return
	}

	cafeID, _ := strconv.Atoi(c.Request.FormValue("cafe_id"))
	date := c.Request.FormValue("date")
	bookTime := c.Request.FormValue("bookTime")
	locationID, _ := strconv.Atoi(c.Request.FormValue("location_id"))
	capacity, _ := strconv.Atoi(c.Request.FormValue("party_size"))

	_, _ = h.reservationService.GetAvailableTables(cafeID, capacity, locationID, date, bookTime)

	panic("implement me")
}
