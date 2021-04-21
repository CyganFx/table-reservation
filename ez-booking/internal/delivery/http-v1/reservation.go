package http_v1

import "github.com/gin-gonic/gin"

func (h *handler) initReservationRoutes(api *gin.RouterGroup) {
	menu := api.Group("/reservation")
	{
		menu.GET("/", h.ReservationPage)
	}
}

func (h *handler) ReservationPage(c *gin.Context) {
	h.render(c, "reservation.page.html", &templateData{})
}
