package http_v1

import (
	"github.com/CyganFx/table-reservation/ez-booking/pkg/domain"
	"github.com/CyganFx/table-reservation/ez-booking/pkg/validator/forms"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"html/template"
	"os"
)

type handler struct {
	userService        UserService
	reservationService ReservationService
	errors             Responser
	templateCache      map[string]*template.Template
}

// handler uses Responser interface to handle errors
// Responser is used in all handlers, therefore we put it here
type Responser interface {
	ServerError(c *gin.Context, err error)
	ClientError(c *gin.Context, status int)
	NotFound(c *gin.Context)
}

// Passing templateData in html pages at render
type templateData struct {
	User            *domain.User
	ReservationData *ReservationData
	Form            *forms.Form
	CurrentYear     int
	Flash           string
	IsAuthenticated bool
}

func NewHandler(userService UserService, reservationService ReservationService, errors Responser,
	templateCache map[string]*template.Template) *handler {
	return &handler{
		userService:        userService,
		reservationService: reservationService,
		errors:             errors,
		templateCache:      templateCache,
	}
}

func (h *handler) Init() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	sessionStore := cookie.NewStore([]byte(os.Getenv("SESSION_SECRET")))

	router.Use(sessions.Sessions("mySessionStore", sessionStore),
		gin.Logger(), gin.Recovery(), SecureHeaders())

	router.GET("/", h.MainPage)

	api := router.Group("/api")
	{
		h.initUserRoutes(api)
		h.initReservationRoutes(api)
		h.initMenuRoutes(api)
	}

	router.Static("/static/", "./ui/static")

	return router
}
