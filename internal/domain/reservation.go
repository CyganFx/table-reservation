package domain

import (
	"github.com/gin-gonic/gin"
	"time"
)

type ReservationHandler interface {
	GetAvailableTables(c *gin.Context)
	BookTable(c *gin.Context)
	ReservationPage(c *gin.Context)
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
	Cafe                    Cafe
	Table                   Table
	Event                   Event
	User                    User
	HoursUntilReservation   int  // not in db
	MinutesUntilReservation int  // not in db
	IsActive                bool // not in db
}

func NewReservation() *Reservation {
	return &Reservation{}
}

//just for benchmarks
func (r *Reservation) Reset() {
	r.ID = 0
	r.PartySize = 0
	r.CustName = ""
	r.CustMobile = ""
	r.CustEmail = ""
	r.EventDescription = ""
	r.Date = time.Time{}
	r.NotifyDate = time.Time{}
	r.Cafe = Cafe{}
	r.Table = Table{}
	r.Event = Event{}
	r.User = User{}
	r.User.Role = Role{}
}

type Table struct {
	ID       int
	Capacity int
	Location Location
}

func NewTable() *Table {
	return &Table{}
}
