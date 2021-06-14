package http_v1

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *handler) initAdminRoutes(api *gin.RouterGroup) {
	admin := api.Group("/admin", h.RequireAdmin())
	{
		admin.GET("/", h.AdminPage)
		admin.GET("/collabs", h.CollabRequestsPage)
		admin.GET("/approve/:id", h.ApproveCollabRequest)
		admin.GET("/disapprove/:id", h.DisapproveCollabRequest)
	}
}

func (h *handler) AdminPage(c *gin.Context) {
	h.render(c, "admin.page.html", nil)
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
	cafeID, err := strconv.Atoi(c.Param("id"))
	if err != nil || cafeID < 1 {
		h.errors.NotFound(c)
		return
	}

	if err := h.cafeService.Approve(cafeID); err != nil {
		h.errors.ServerError(c, err)
		return
	}

	session := sessions.Default(c)
	session.Set("flash", fmt.Sprintf("Approved cafe #%d successfully!", cafeID))
	session.Save()

	http.Redirect(c.Writer, c.Request, "/api/admin/collabs", http.StatusSeeOther)
}

func (h *handler) DisapproveCollabRequest(c *gin.Context) {
	cafeID, err := strconv.Atoi(c.Param("id"))
	if err != nil || cafeID < 1 {
		h.errors.NotFound(c)
		return
	}

	if err := h.cafeService.Disapprove(cafeID); err != nil {
		h.errors.ServerError(c, err)
		return
	}

	session := sessions.Default(c)
	session.Set("flash", fmt.Sprintf("Disapproved cafe #%d", cafeID))
	session.Save()

	http.Redirect(c.Writer, c.Request, "/api/admin/collabs", http.StatusSeeOther)
}
