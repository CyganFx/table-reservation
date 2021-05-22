package http_v1

import (
	"github.com/CyganFx/table-reservation/ez-booking/internal/domain"
	"github.com/CyganFx/table-reservation/ez-booking/pkg/validator/forms"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	aws_session "github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"os"
)

var (
	AccessKeyID     string
	SecretAccessKey string
	MyRegion        string
	MyBucket        string
)

type handler struct {
	userService        UserService
	reservationService ReservationService
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
	User            *domain.User
	ReservationData *ReservationData
	Reservations    []*domain.Reservation
	Form            *forms.FormValidator
	CurrentYear     int
	Flash           string
	IsAuthenticated bool
}

func NewHandler(userService UserService, reservationService ReservationService, errors Responser,
	infoLog *log.Logger, templateCache map[string]*template.Template) *handler {
	return &handler{
		userService:        userService,
		reservationService: reservationService,
		errors:             errors,
		infoLog:            infoLog,
		templateCache:      templateCache,
	}
}

func (h *handler) Init() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	sessionStore := cookie.NewStore([]byte(os.Getenv("SESSION_SECRET")))
	awsSession := ConnectAws()

	router.Use(sessions.Sessions("mySessionStore", sessionStore),
		gin.Logger(), gin.Recovery(), h.SecureHeaders(),
		func(c *gin.Context) {
			c.Set("awsSession", awsSession)
			c.Next()
		})

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

func ConnectAws() *aws_session.Session {
	AccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
	SecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	MyRegion = os.Getenv("AWS_REGION")
	MyBucket = os.Getenv("BUCKET_NAME")

	sess, err := aws_session.NewSession(
		&aws.Config{
			Region: aws.String(MyRegion),
			Credentials: credentials.NewStaticCredentials(
				AccessKeyID,
				SecretAccessKey,
				"", // a token will be created when the session it's used.
			),
		})
	if err != nil {
		panic(err)
	}

	return sess
}
