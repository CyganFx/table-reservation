package domain

import (
	"github.com/gin-gonic/gin"
	"time"
)

type Table struct {
	ID       int
	Capacity int
	Location string
}

type Reservation struct {
	ID           int
	Table        *Table
	CustName     string
	CustMobile   string
	CustEmail    string
	Event        string
	NumOfPersons int
	Date         time.Time
}

type Location struct {
	ID   int
	Name string
}

func NewReservation() *Reservation {
	return &Reservation{
		Table: &Table{},
	}
}

type DaySchema struct {
	Date   time.Time
	Tables []*Table
}

type ReservationHandler interface {
	GetAvailableTables(c *gin.Context)
}
