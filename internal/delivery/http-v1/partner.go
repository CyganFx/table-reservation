package http_v1

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *handler) initPartnerRoutes(api *gin.RouterGroup) {
	admin := api.Group("/partner", h.RequireRole(partnerRoleID))
	{
		admin.GET("/", h.PartnerPage)

		admin.GET("/busy/user/:id", h.PartnerBusyReservationPage)
		admin.POST("/reservation/available/tables", h.GetAvailableTablesForPartner)
		admin.POST("/reservation/busy/confirm", h.PartnerBusyTableConfirm)

		admin.GET("/free/user/:id", h.PartnerFreeReservationPage)
		admin.POST("/reservation/reserved/tables", h.GetReservedTablesForPartner)
		admin.POST("/reservation/free/confirm", h.PartnerFreeTableConfirm)

		admin.GET("/reportPage/:id", h.ReportPage)
		admin.POST("/report", h.Report)
	}
}

func (h *handler) PartnerPage(c *gin.Context) {
	h.render(c, "partner.admin.page.html", &templateData{})
}

func (h *handler) ReportPage(c *gin.Context) {
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

	users, err := h.userService.GetAll(cafeID)
	if err != nil {
		h.errors.ServerError(c, err)
		return
	}

	h.render(c, "partner.report.page.html", &templateData{Users: users})
}

func (h *handler) Report(c *gin.Context) {
	if err := c.Request.ParseForm(); err != nil {
		h.errors.ClientError(c, http.StatusBadRequest)
		return
	}

	adminID, _ := strconv.Atoi(c.Request.FormValue("adminID"))
	userID, _ := strconv.Atoi(c.Request.FormValue("userID"))
	cafeID, err := h.cafeService.GetCafeIDByAdminID(adminID)
	if err != nil {
		h.errors.ServerError(c, err)
		return
	}

	if err := h.userService.AddToBlacklist(userID, cafeID); err != nil {
		h.errors.ServerError(c, err)
		return
	}

	session := sessions.Default(c)
	session.Set("flash", "Reported successfully!")
	session.Save()

	http.Redirect(c.Writer, c.Request, fmt.Sprintf("/api/partner/reportPage/%d", adminID), http.StatusSeeOther)
}

func (h *handler) PartnerBusyReservationPage(c *gin.Context) {
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

	h.render(c, "partner.busy.reservation.page.html", &templateData{
		ReservationData: reservationData,
	})
}

func (h *handler) PartnerFreeReservationPage(c *gin.Context) {
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

	h.render(c, "partner.free.reservation.page.html", &templateData{
		ReservationData: reservationData,
	})
}

// Sets a table as busy (f.e. when someone just came in the cafe without reservation; should be done manually)
func (h *handler) PartnerBusyTableConfirm(c *gin.Context) {
	if err := c.Request.ParseForm(); err != nil {
		h.errors.ClientError(c, http.StatusBadRequest)
		return
	}

	tableID, _ := strconv.Atoi(c.Request.FormValue("table_id"))
	adminID, _ := strconv.Atoi(c.Request.FormValue("userID"))

	session := sessions.Default(c)
	userChoice := session.Get("userChoice").(UserChoice)
	userChoice.TableID = tableID

	_, err := h.reservationService.BookTableManually(userChoice, adminID)
	if err != nil {
		h.errors.ServerError(c, err)
		return
	}

	session.Delete("userChoice")
	session.Set("flash", "Set as busy successfully!")
	session.Save()

	http.Redirect(c.Writer, c.Request, fmt.Sprintf("/api/partner/busy/user/%d", adminID), http.StatusSeeOther)
}

func (h *handler) PartnerFreeTableConfirm(c *gin.Context) {
	if err := c.Request.ParseForm(); err != nil {
		h.errors.ClientError(c, http.StatusBadRequest)
		return
	}

	tableID, _ := strconv.Atoi(c.Request.FormValue("table_id"))
	adminID, _ := strconv.Atoi(c.Request.FormValue("userID"))

	session := sessions.Default(c)
	userChoice := session.Get("userChoice").(UserChoice)
	userChoice.TableID = tableID

	err := h.reservationService.FreeTableManually(userChoice, adminID)
	if err != nil {
		h.errors.ServerError(c, err)
		return
	}

	session.Set("flash", "Freed table successfully!")
	session.Save()

	http.Redirect(c.Writer, c.Request, fmt.Sprintf("/api/partner/free/user/%d", adminID), http.StatusSeeOther)
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

	h.render(c, "partner.busy.reservation.page.html", &templateData{
		ReservationData: reservationData,
	})
}

func (h *handler) GetReservedTablesForPartner(c *gin.Context) {
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
	reservationData.Tables, err = h.reservationService.GetBusyTables(cafeID, partySize, locationID, date, bookTime)
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

	h.render(c, "partner.free.reservation.page.html", &templateData{
		ReservationData: reservationData,
	})
}
