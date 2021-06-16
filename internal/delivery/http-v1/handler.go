package http_v1

import (
	"github.com/CyganFx/table-reservation/internal/app/config"
	"github.com/CyganFx/table-reservation/internal/domain"
	"github.com/CyganFx/table-reservation/pkg/validator/forms"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	aws_session "github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/nosurf"
	"html/template"
	"log"
	"net/http"
	"os"
)

const StaticFilesDir = "./ui/static"

var (
	AccessKey       string
	SecretAccessKey string
	MyRegion        string
	MyBucket        string
)

const (
	adminRoleID = 1
	userRoleID = 2
	partnerRoleID = 3
)

type handler struct {
	userService        UserService
	reservationService ReservationService
	cafeService        CafeService
	notificatorService NotificatorService
	errors             Responser
	infoLog            *log.Logger
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
	CSRFToken       string
	User            *domain.User
	ReservationData *ReservationData
	CollaborateData *CollaborateData
	Cafes           []domain.Cafe
	Types           []domain.Type
	Cities          []domain.City
	Reservations    []domain.Reservation
	Form            *forms.FormValidator
	CurrentYear     int
	Flash           string
	IsAuthenticated bool
}

func NewHandler(userService UserService, reservationService ReservationService, cafeService CafeService, notificatorService NotificatorService, errors Responser,
	infoLog *log.Logger, templateCache map[string]*template.Template) *handler {
	return &handler{
		userService:        userService,
		reservationService: reservationService,
		cafeService:        cafeService,
		notificatorService: notificatorService,
		errors:             errors,
		infoLog:            infoLog,
		templateCache:      templateCache,
	}
}

//Init set ups router and routes and wraps it in csrf middleware that checks for csrf token in every request
func (h *handler) Init(cfg config.Config) http.Handler {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	csrfHandler := nosurf.New(router)
	csrfHandler.SetFailureHandler(http.HandlerFunc(csrfFailHandler))

	sessionStore := cookie.NewStore([]byte(os.Getenv("SESSION_SECRET")))
	sessionStore.Options(sessions.Options{
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   3600,
	})

	awsSession := connectAws(cfg)

	router.Use(sessions.Sessions("mySessionStore", sessionStore),
		gin.Logger(), gin.Recovery(), h.SecureHeaders(),
		func(c *gin.Context) {
			c.Set("awsSession", awsSession)
			c.Next()
		})

	router.GET("/", h.MainPage)

	api := router.Group("/api")
	{
		h.initAdminRoutes(api)
		h.initUserRoutes(api)
		h.initReservationRoutes(api)
		h.initCafeRoutes(api)
	}

	router.Static("/static/", StaticFilesDir)

	return csrfHandler
}

func connectAws(cfg config.Config) *aws_session.Session {
	AccessKey = cfg.FileStorage.AccessKey
	MyRegion = cfg.FileStorage.AwsRegion
	MyBucket = cfg.FileStorage.BucketName
	SecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")

	sess, err := aws_session.NewSession(
		&aws.Config{
			Region: aws.String(MyRegion),
			Credentials: credentials.NewStaticCredentials(
				AccessKey,
				SecretAccessKey,
				"", // a token will be created when the session it's used.
			),
		})
	if err != nil {
		panic(err)
	}

	return sess
}
