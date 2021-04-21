package http_v1

import (
	"errors"
	"github.com/CyganFx/table-reservation/ez-booking/pkg/domain"
	"github.com/CyganFx/table-reservation/ez-booking/pkg/validator/forms"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserService interface {
	Save(form *forms.Form) (bool, error)
	SignIn(email, password string) (int, error)
	FindById(id int) (*domain.User, error)
	Update(user *domain.User) error
}

func (h *handler) initUserRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.GET("/sign-up", h.SignUpPage)
		users.POST("/sign-up", h.SignUp)

		users.GET("/login", h.LoginPage)
		users.POST("/login", h.Login)

		users.POST("/logout", h.Logout)

		users.GET("/show/:id", h.ShowById)
		users.POST("/update/:id", h.Update)
	}
}

func (h *handler) MainPage(c *gin.Context) {
	h.render(c, "main.page.html", &templateData{})
}

func (h *handler) SignUpPage(c *gin.Context) {
	h.render(c, "signup.page.html", &templateData{
		Form: forms.New(nil),
	})
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

func (h *handler) ShowById(c *gin.Context) {
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

	h.render(c, "test.page.html", &templateData{User: user})
}

func (h *handler) Update(c *gin.Context) {
	panic("implement me")
}
