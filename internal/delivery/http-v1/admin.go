package http_v1

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *handler) initAdminRoutes(api *gin.RouterGroup) {
	admin := api.Group("/admin", h.RequireRole(adminRoleID))
	{
		admin.GET("/", h.AdminPage)
		admin.GET("/collabs", h.CollabRequestsPage)
		admin.POST("/approve", h.ApproveCollabRequest)
		admin.POST("/disapprove", h.DisapproveCollabRequest)
	}
}

func (h *handler) AdminPage(c *gin.Context) {
	// this is slow, because we need only counter of collabs, but im lazy
	cafes, err := h.cafeService.GetCollabRequests()
	if err != nil {
		h.errors.ServerError(c, err)
		return
	}

	h.render(c, "admin.page.html", &templateData{Cafes: cafes})
}

func (h *handler) CollabRequestsPage(c *gin.Context) {
	cafes, err := h.cafeService.GetCollabRequests()
	if err != nil {
		h.errors.ServerError(c, err)
		return
	}

	h.render(c, "admin-collabRequests.page.html", &templateData{Cafes: cafes})
}

func (h *handler) ApproveCollabRequest(c *gin.Context) {
	cafeID, _ := strconv.Atoi(c.Request.FormValue("cafeID"))
	partnerID, _ := strconv.Atoi(c.Request.FormValue("adminID"))
	email := c.Request.FormValue("email")

	if err := h.cafeService.Approve(cafeID); err != nil {
		h.errors.ServerError(c, err)
		return
	}

	if err := h.userService.UpdateUserRole(partnerID, partnerRoleID); err != nil {
		h.errors.ServerError(c, err)
		return
	}

	if err := h.notificatorService.AdminResponseToPartnership(email, true); err != nil {
		h.errors.ServerError(c, err)
		return
	}

	session := sessions.Default(c)
	session.Set("flash", fmt.Sprintf("Approved cafe #%d successfully!", cafeID))
	session.Save()

	http.Redirect(c.Writer, c.Request, "/api/admin/collabs", http.StatusSeeOther)
}

func (h *handler) DisapproveCollabRequest(c *gin.Context) {
	cafeID, _ := strconv.Atoi(c.Request.FormValue("cafeID"))
	email := c.Request.FormValue("email")

	if err := h.cafeService.Disapprove(cafeID); err != nil {
		h.errors.ServerError(c, err)
		return
	}

	if err := h.notificatorService.AdminResponseToPartnership(email, false); err != nil {
		h.errors.ServerError(c, err)
		return
	}

	session := sessions.Default(c)
	session.Set("flash", fmt.Sprintf("Disapproved cafe #%d", cafeID))
	session.Save()

	http.Redirect(c.Writer, c.Request, "/api/admin/collabs", http.StatusSeeOther)
}
