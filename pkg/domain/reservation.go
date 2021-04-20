package domain

import "time"

// TODO(Duman) Create tables in db

type Table struct {
	ID          int
	Capacity    int
	IsAvailable bool
	Location    string
	Reservation ReservationSchema
}

type ReservationSchema struct {
	Name   string
	Mobile string
	Email  string
}

type DaySchema struct {
	Date   time.Time
	Tables []Table
}
