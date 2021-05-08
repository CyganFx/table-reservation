package http_v1

import "github.com/gin-gonic/gin"

func (h *handler) initMenuRoutes(api *gin.RouterGroup) {
	menu := api.Group("/menu")
	{
		menu.GET("/", h.MenuPage)
	}
}

// Static data (not important)
func (h *handler) MenuPage(c *gin.Context) {
	h.render(c, "menu.page.html", &templateData{})
}
