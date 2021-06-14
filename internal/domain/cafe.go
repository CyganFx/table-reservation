package domain

import (
	"github.com/gin-gonic/gin"
	"time"
)

type CafeHandler interface {
	Collaborate(c *gin.Context)
}

type AdminHandler interface {
	AdminPage(c *gin.Context)
}

type Cafe struct {
	ID          int
	Name        string
	Address     string
	ImageURL    string
	Mobile      string
	Email       string
	Description string
	City        City
	Type        Type
	Created     time.Time
}

type Type struct {
	ID   int
	Name string
}

type City struct {
	ID   int
	Name string
}

type Location struct {
	ID   int
	Name string
}

type Event struct {
	ID   int
	Name string
}
