package http_v1

import (
	"errors"
	"fmt"
	"github.com/CyganFx/table-reservation/internal/domain"
	"github.com/CyganFx/table-reservation/pkg/validator/forms"
	aws_session "github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
)

var (
	defaultProfileImagePath = "/static/img/default_profile_image.png"
	fileStorageEndPoint     = "https://ez-booking-bucket.s3-eu-north-1.amazonaws.com/"
)

func (h *handler) initUserRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.GET("/contributors", h.ContributorsPage)
		users.GET("/sign-up", h.SignUpPage)
		users.POST("/sign-up", h.SignUp)
		users.GET("/login", h.LoginPage)
		users.POST("/login", h.Login)

		authenticated := users.Group("/", h.RequireAuthentication())
		{
			authenticated.GET("/profile/:id", h.ProfilePage)
			authenticated.POST("/logout", h.Logout)
			authenticated.POST("/set-image", h.UpdateImage)
			authenticated.POST("/update/:id", h.Update)
		}
	}
}

type UserService interface {
	Save(form *forms.FormValidator) (bool, error)
	SignIn(email, password string) (int, error)
	FindById(id int) (*domain.User, error)
	Update(user *domain.User) error
	UpdateImage(filePath string, userID int) error
	UploadImageToAWSBucket(awsSession *aws_session.Session, MyBucket, filename string, file multipart.File) error
	DeleteImageFromAWSBucket(awsSession *aws_session.Session, imageURL, myBucket, objectsLocationURL string, infoLog *log.Logger) error
	SetConfirmData(ctx *gin.Context, reservationData *ReservationData, tableID, eventID int, eventDescription string) (*forms.FormValidator, error)
	UpdateUserRole(userID, roleID int) error
}

func (h *handler) ProfilePage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		h.errors.NotFound(c)
		return
	}

	user, err := h.userService.FindById(id)
	if err != nil {
		if errors.Is(err, domain.ErrNoRecord) {
			h.errors.NotFound(c)
			return
		} else {
			h.errors.ServerError(c, err)
			return
		}
	}

	bookings, err := h.reservationService.GetUserBookings(id)
	if err != nil {
		h.errors.ServerError(c, err)
		return
	}

	h.render(c, "profile.page.html", &templateData{
		User:         user,
		Reservations: bookings,
	})
}

func (h *handler) ContributorsPage(c *gin.Context) {
	h.render(c, "contributors.page.html", &templateData{})
}

func (h *handler) SignUpPage(c *gin.Context) {
	h.render(c, "signup.page.html", &templateData{Form: forms.New(nil)})
}

func (h *handler) SignUp(c *gin.Context) {
	if err := c.Request.ParseForm(); err != nil {
		h.errors.ClientError(c, http.StatusBadRequest)
		return
	}

	form := forms.New(c.Request.PostForm)

	valid, err := h.userService.Save(form)
	if !valid {
		h.render(c, "signup.page.html", &templateData{Form: form})
		return
	}
	if err != nil {
		if errors.Is(err, domain.ErrDuplicateEmail) {
			form.Errors.Add("email", "Address is already in use")
			h.render(c, "signup.page.html", &templateData{Form: form})
		} else {
			h.errors.ServerError(c, err)
		}
		return
	}

	session := sessions.Default(c)
	session.Set("flash", "Your sign up was successful. Please login.")
	session.Save()

	http.Redirect(c.Writer, c.Request, "/api/users/login", http.StatusSeeOther)
}

func (h *handler) LoginPage(c *gin.Context) {
	h.render(c, "login.page.html", &templateData{Form: forms.New(nil)})
}

func (h *handler) Login(c *gin.Context) {
	if err := c.Request.ParseForm(); err != nil {
		h.errors.ClientError(c, http.StatusBadRequest)
		return
	}

	form := forms.New(c.Request.PostForm)

	id, err := h.userService.SignIn(form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Email or Password is incorrect")
			h.render(c, "login.page.html", &templateData{Form: form})
		} else {
			h.errors.ServerError(c, err)
		}
		return
	}

	user, err := h.userService.FindById(id)
	if err != nil {
		h.errors.NotFound(c)
		return
	}

	session := sessions.Default(c)
	session.Set("authenticatedUserID", id)
	session.Set("role", user.Role.ID)
	session.Set("flash", "loginned successfully!")

	if session.Get("redirectPathAfterLogin") == nil {
		session.Save()
		http.Redirect(c.Writer, c.Request, "/", http.StatusSeeOther)
		return
	}

	path := session.Get("redirectPathAfterLogin").(string)
	session.Delete("redirectPathAfterLogin")
	session.Save()

	http.Redirect(c.Writer, c.Request, path, http.StatusSeeOther)
}

func (h *handler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("authenticatedUserID")
	session.Delete("role")
	session.Set("flash", "You've been logged out successfully!")
	session.Save()

	http.Redirect(c.Writer, c.Request, "/", http.StatusSeeOther)
}

////using AWS S3 to store profile images
func (h *handler) UpdateImage(c *gin.Context) {
	awsSession := c.MustGet("awsSession").(*aws_session.Session)
	file, header, err := c.Request.FormFile("profile-image")
	if err != nil {
		h.errors.ServerError(c, fmt.Errorf("error retrieving the file %v", err))
		return
	}
	defer file.Close()

	filename := header.Filename
	userID, _ := strconv.Atoi(c.Request.FormValue("id"))

	user, err := h.userService.FindById(userID)
	if err != nil {
		if errors.Is(err, domain.ErrNoRecord) {
			h.errors.NotFound(c)
		} else {
			h.errors.ServerError(c, err)
		}
		return
	}

	if user.ImageURL != defaultProfileImagePath {
		err = h.userService.DeleteImageFromAWSBucket(awsSession, user.ImageURL, MyBucket, fileStorageEndPoint, h.infoLog)
		if err != nil {
			h.errors.ServerError(c, err)
			return
		}
	}

	err = h.userService.UploadImageToAWSBucket(awsSession, MyBucket, filename, file)
	if err != nil {
		h.errors.ServerError(c, err)
		return
	}

	filepath := fileStorageEndPoint + filename
	if err = h.userService.UpdateImage(filepath, userID); err != nil {
		h.errors.ServerError(c, err)
		return
	}

	http.Redirect(c.Writer, c.Request, fmt.Sprintf("/api/users/profile/%d", userID), http.StatusSeeOther)
}

func (h *handler) Update(c *gin.Context) {
	panic("implement me")
}
