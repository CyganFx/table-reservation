package domain

import (
	"errors"
	"github.com/gin-gonic/gin"
	"time"
)

var (
	ErrNoRecord           = errors.New("domain: no matching record found")
	ErrInvalidCredentials = errors.New("domain: invalid credentials")
	ErrDuplicateEmail     = errors.New("domain: duplicate email")
)

type UserHandler interface {
	SignUp(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
	ProfilePage(c *gin.Context)
	UpdateImage(c *gin.Context)
	Update(c *gin.Context)
}

type User struct {
	ID       int       `json:"id"`
	Name     string    `json:"username"`
	Email    string    `json:"email"`
	Mobile   string    `json:"mobile"`
	ImageURL string    `json:"imageURL"`
	Password []byte    `json:"password"`
	Created  time.Time `json:"created"`
	Role     Role      `json:"role"`
}

func NewUser() *User {
	return &User{}
}

type Role struct {
	ID   int
	Name string
}
