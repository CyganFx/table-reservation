package domain

import "time"

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

func NewReservation() *Reservation {
	return &Reservation{
		Table: &Table{},
	}
}

type DaySchema struct {
	Date   time.Time
	Tables []*Table
}
