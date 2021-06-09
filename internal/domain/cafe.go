package domain

import (
	"github.com/gin-gonic/gin"
	"time"
)

type CafeHandler interface {
	Collaborate(c *gin.Context)
}

type Cafe struct {
	ID      int
	Name    string
	Address string
	Mobile  string
	Email   string
	City    City
	Type    Type
	Created time.Time
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
