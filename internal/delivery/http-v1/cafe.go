package http_v1

import (
	"fmt"
	"github.com/CyganFx/table-reservation/internal/domain"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *handler) initCafeRoutes(api *gin.RouterGroup) {
	cafe := api.Group("/cafe")
	{
		authenticated := cafe.Group("/", h.RequireAuthentication())
		{
			authenticated.GET("/collaborate", h.CollaboratePage)
			authenticated.POST("/collaborate", h.Collaborate)
		}
	}
}

type CafeService interface {
	GetLocations() ([]domain.Location, error)
	GetTypes() ([]domain.Type, error)
	GetEvents() ([]domain.Event, error)
	GetCities() ([]domain.City, error)
	GetCafes() ([]domain.Cafe, error)
	CreateCafe(cafe *domain.Cafe) error
	SetLocations(cafeID int, locations []string) error
	SetEvents(cafeID int, events []string) error
	SetTables(cafeID, locationID, numOfTables, capacity int) error
}

type CollaborateData struct {
	Locations []domain.Location
	Events    []domain.Event
	Types     []domain.Type
	Cities    []domain.City
}

func (h *handler) MainPage(c *gin.Context) {
	cafes, err := h.cafeService.GetCafes()
	if err != nil {
		h.errors.ServerError(c, err)
		return
	}
	h.render(c, "main.page.html", &templateData{Cafes: cafes})
}

func (h *handler) CollaboratePage(c *gin.Context) {
	fmt.Print("I am here!")
	locations, err := h.cafeService.GetLocations()
	if err != nil {
		h.errors.ServerError(c, err)
		return
	}
	events, err := h.cafeService.GetEvents()
	if err != nil {
		h.errors.ServerError(c, err)
		return
	}
	types, err := h.cafeService.GetTypes()
	if err != nil {
		h.errors.ServerError(c, err)
		return
	}
	cities, err := h.cafeService.GetCities()
	if err != nil {
		h.errors.ServerError(c, err)
		return
	}

	data := &CollaborateData{
		Locations: locations,
		Events:    events,
		Types:     types,
		Cities:    cities,
	}

	h.render(c, "collaborate.page.html", &templateData{
		CollaborateData: data,
	})
}

func (h *handler) Collaborate(c *gin.Context) {
	if err := c.Request.ParseForm(); err != nil {
		h.errors.ClientError(c, http.StatusBadRequest)
		return
	}

	name := c.Request.FormValue("name")
	address := c.Request.FormValue("address")
	mobile := c.Request.FormValue("mobile")
	email := c.Request.FormValue("email")
	cityID, _ := strconv.Atoi(c.Request.FormValue("city"))
	typeID, _ := strconv.Atoi(c.Request.FormValue("type"))

	cafe := &domain.Cafe{
		Name:    name,
		Address: address,
		Mobile:  mobile,
		Email:   email,
		City:    domain.City{},
		Type:    domain.Type{},
	}
	cafe.City.ID = cityID
	cafe.Type.ID = typeID

	if err := h.cafeService.CreateCafe(cafe); err != nil {
		h.errors.ServerError(c, err)
		return
	}

	locations := c.Request.Form["locations"]
	if err := h.cafeService.SetLocations(cafe.ID, locations); err != nil {
		h.errors.ServerError(c, err)
		return
	}
	events := c.Request.Form["events"]
	if err := h.cafeService.SetEvents(cafe.ID, events); err != nil {
		h.errors.ServerError(c, err)
		return
	}

	tableTypesCounter, _ := strconv.Atoi(c.Request.FormValue("tableTypesCounter")) //indicates how many different tables with different size and location set

	for i := 0; i < tableTypesCounter; i++ {
		locationID, _ := strconv.Atoi(c.Request.FormValue(fmt.Sprintf("location%d", i)))
		numOfTables, _ := strconv.Atoi(c.Request.FormValue(fmt.Sprintf("number%d", i)))
		capacity, _ := strconv.Atoi(c.Request.FormValue(fmt.Sprintf("capacity%d", i)))
		if err := h.cafeService.SetTables(cafe.ID, locationID, numOfTables, capacity); err != nil {
			h.errors.ServerError(c, err)
			return
		}
	}

	if err := h.notificatorService.CollaborationNotify(*cafe); err != nil {
		h.errors.ServerError(c, err)
		return
	}

	session := sessions.Default(c)
	session.Set("flash", "Thanks, your request is being considered")
	session.Save()

	h.MainPage(c)
}
