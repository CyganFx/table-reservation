package http_v1

import (
	"errors"
	"fmt"
	"github.com/CyganFx/table-reservation/ez-booking/pkg/domain"
	"github.com/CyganFx/table-reservation/ez-booking/pkg/validator/forms"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	aws_session "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

var (
	DefaultProfileImagePath = "/static/img/default_profile_image.png"
)

type UserService interface {
	Save(form *forms.FormValidator) (bool, error)
	SignIn(email, password string) (int, error)
	FindById(id int) (*domain.User, error)
	Update(user *domain.User) error
	UpdateImage(filePath string, userID int) error
}

func (h *handler) initUserRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.GET("/profile/:id", h.ProfilePage)
		users.GET("/sign-up", h.SignUpPage)
		users.POST("/sign-up", h.SignUp)

		users.GET("/login", h.LoginPage)
		users.POST("/login", h.Login)

		users.POST("/logout", h.Logout)
		users.POST("/set-image", h.UpdateImage)
		users.POST("/update/:id", h.Update)
	}
}

func (h *handler) MainPage(c *gin.Context) {
	session := sessions.Default(c)
	session.Get("redirectPathAfterLogin")
	h.render(c, "main.page.html", &templateData{})
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

	h.render(c, "profile.page.html", &templateData{User: user})
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
	session.Set("flash", "You've been logged out successfully!")
	session.Save()
	http.Redirect(c.Writer, c.Request, "/", http.StatusSeeOther)
}

//using aws s3 to store profile images
//TODO seperate logic in service
func (h *handler) UpdateImage(c *gin.Context) {
	objectsLocationURL := "https://" + MyBucket + "." + "s3-" + MyRegion + ".amazonaws.com/"
	awsSession := c.MustGet("awsSession").(*aws_session.Session)
	uploader := s3manager.NewUploader(awsSession)
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
			return
		} else {
			h.errors.ServerError(c, err)
			return
		}
	}

	if user.ImageURL != DefaultProfileImagePath {
		objectName := strings.ReplaceAll(user.ImageURL, objectsLocationURL, "")
		//delete from s3 bucket
		svc := s3.New(awsSession)
		input := &s3.DeleteObjectInput{
			Bucket: aws.String(MyBucket),
			Key:    aws.String(objectName),
		}

		result, err := svc.DeleteObject(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				default:
					h.errors.ServerError(c,
						fmt.Errorf("failed to delete object from aws bucket, here is aws error code %v",
							aerr.Error()))
					return
				}
			} else {
				h.errors.ServerError(c, fmt.Errorf("failed to delete object from aws bucket %v", err))
				return
			}
		}
		h.infoLog.Printf("Object deleted from aws s3 bucket: %v", result)
	}

	//upload to the s3 bucket
	up, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(MyBucket),
		ACL:    aws.String("public-read"),
		Key:    aws.String(filename),
		Body:   file,
	})
	if err != nil {
		h.errors.ServerError(c, fmt.Errorf("failed to upload file, uploader: %v / error: %v ", up, err))
		return
	}

	filepath := objectsLocationURL + filename

	if err = h.userService.UpdateImage(filepath, userID); err != nil {
		h.errors.ServerError(c, err)
		return
	}

	http.Redirect(c.Writer, c.Request, fmt.Sprintf("/api/users/profile/%d", userID), http.StatusSeeOther)
}

func (h *handler) Update(c *gin.Context) {
	panic("implement me")
}
