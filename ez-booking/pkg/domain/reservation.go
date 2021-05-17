package domain

import (
	"github.com/gin-gonic/gin"
)

type Table struct {
	ID              int
	Capacity        int
	CapacityForHTML []int //go can't iterate over integer, therefore creating slice
	LocationID      int
}

type Reservation struct {
	ID               int
	PartySize        int
	CustName         string
	CustMobile       string
	CustEmail        string
	Date             string //using string for convenience
	EventDescription string
	Cafe             *Cafe
	Table            *Table
	Event            *Event
	User             *User
}

type Location struct {
	ID   int
	Name string
}

type Event struct {
	ID   int
	Name string
}

type Cafe struct {
	ID   int
	Name string
}

func NewReservation() *Reservation {
	return &Reservation{
		Cafe:  &Cafe{},
		Table: &Table{},
		Event: &Event{},
		User:  &User{},
	}
}

type ReservationHandler interface {
	GetAvailableTables(c *gin.Context)
	BookTable(c *gin.Context)
}
