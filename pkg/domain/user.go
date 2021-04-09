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

// more fields will be added later
type User struct {
	ID       int       `json:"id"`
	Name     string    `json:"username"`
	Email    string    `json:"email"`
	Mobile   string    `json:"mobile"`
	Password []byte    `json:"password"`
	Created  time.Time `json:"created"`
}

type UserHandler interface {
	SignUp(c *gin.Context)
	Login(c *gin.Context)
	ShowById(c *gin.Context)
	Update(c *gin.Context)
	Init() *gin.Engine
}
