package http_v1

import (
	"github.com/gin-gonic/gin"
)

func (h *handler) initPartnerAdminRoutes(api *gin.RouterGroup) {
	admin := api.Group("/partner/admin", h.RequireRole(partnerRoleID))
	{
		admin.GET("/", h.AdminPage)
	}
}
