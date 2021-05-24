package domain

import (
	"github.com/gin-gonic/gin"
	"time"
)

type Table struct {
	ID              int
	Capacity        int
	CapacityForHTML []int //go can't iterate over integer, therefore creating slice
	Location        *Location
}

func NewTable() *Table {
	return &Table{Location: &Location{}}
}

type Reservation struct {
	ID                      int
	PartySize               int
	CustName                string
	CustMobile              string
	CustEmail               string
	EventDescription        string
	Date                    time.Time
	NotifyDate              time.Time
	Cafe                    *Cafe
	Table                   *Table
	Event                   *Event
	User                    *User
	HoursUntilReservation   int  // not in db
	MinutesUntilReservation int  // not in db
	IsActive                bool // not in db
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
		Table: &Table{Location: &Location{}},
		Event: &Event{},
		User:  &User{Role: &Role{}},
	}
}

type ReservationHandler interface {
	GetAvailableTables(c *gin.Context)
	BookTable(c *gin.Context)
	ReservationPage(c *gin.Context)
}
