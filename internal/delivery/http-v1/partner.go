package http_v1

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *handler) initPartnerRoutes(api *gin.RouterGroup) {
	admin := api.Group("/partner", h.RequireRole(partnerRoleID))
	{
		admin.GET("/", h.PartnerPage)
		admin.GET("/user/:id", h.PartnerReservationPage)
		admin.POST("/reservation/confirm", h.PartnerReservationConfirm)
		admin.POST("/reservation/tables", h.GetAvailableTablesForPartner)
	}
}

func (h *handler) PartnerPage(c *gin.Context) {
	h.render(c, "partner.admin.page.html", &templateData{})
}

func (h *handler) PartnerReservationPage(c *gin.Context) {
	adminID, err := strconv.Atoi(c.Param("id"))
	if err != nil || adminID < 1 {
		h.errors.NotFound(c)
		return
	}

	cafeID, err := h.cafeService.GetCafeIDByAdminID(adminID)
	if err != nil {
		h.errors.ServerError(c, err)
		return
	}

	reservationData := &ReservationData{}
	err = h.reservationService.SetDefaultReservationData(reservationData, cafeID)
	if err != nil {
		h.errors.ServerError(c, err)
		return
	}

	h.render(c, "partner.reservation.page.html", &templateData{
		ReservationData: reservationData,
	})
}

// actually sets a table as busy (f.e. when someone just came in the cafe without reservation; should be done manually)
func (h *handler) PartnerReservationConfirm(c *gin.Context) {
	if err := c.Request.ParseForm(); err != nil {
		h.errors.ClientError(c, http.StatusBadRequest)
		return
	}

	tableID, _ := strconv.Atoi(c.Request.FormValue("table_id"))
	userID, _ := strconv.Atoi(c.Request.FormValue("userID"))

	session := sessions.Default(c)
	userChoice := session.Get("userChoice").(UserChoice)
	userChoice.TableID = tableID

	_, err := h.reservationService.BookTableManually(userChoice, userID)
	if err != nil {
		h.errors.ServerError(c, err)
		return
	}

	session.Delete("userChoice")
	session.Set("flash", "Set as busy successfully!")
	session.Save()

	http.Redirect(c.Writer, c.Request, "/", http.StatusSeeOther)
}

func (h *handler) GetAvailableTablesForPartner(c *gin.Context) {
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

	err = h.reservationService.SetDefaultReservationData(reservationData, cafeID)
	if err != nil {
		h.errors.ServerError(c, err)
		return
	}

	reservationData.CurrentDate = date

	h.render(c, "partner.reservation.page.html", &templateData{
		ReservationData: reservationData,
	})
}
