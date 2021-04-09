package http_v1

import (
	"github.com/CyganFx/table-reservation/pkg/domain"
	"github.com/CyganFx/table-reservation/pkg/validator/forms"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"html/template"
	"os"
)

type handler struct {
	userService   UserService
	errors        Responser
	templateCache map[string]*template.Template
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
	Users           []*domain.User
	Form            *forms.Form
	CurrentYear     int
	Flash           string
	IsAuthenticated bool
}

func NewHandler(userService UserService, errors Responser,
	templateCache map[string]*template.Template) domain.UserHandler {
	return &handler{
		userService:   userService,
		errors:        errors,
		templateCache: templateCache,
	}
}

func (h *handler) Init() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	sessionStore := cookie.NewStore([]byte(os.Getenv("SESSION_SECRET")))

	router.Use(sessions.Sessions("mySessionStore", sessionStore),
		gin.Logger(), gin.Recovery(), SecureHeaders())

	api := router.Group("/api")
	{
		h.initUserRoutes(api)

		api.GET("/main", h.MainPage)
	}

	router.Static("/static/", "./ui/static")

	return router
}
