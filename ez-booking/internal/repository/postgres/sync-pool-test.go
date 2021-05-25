package postgres

import (
	"github.com/CyganFx/table-reservation/ez-booking/internal/domain"
	"time"
)

func WithPoolReseting() {
	var rr []domain.Reservation
	testRR := []domain.Reservation{
		{
			ID:                      1,
			PartySize:               1,
			CustName:                "asd",
			CustMobile:              "dsa",
			CustEmail:               "zxc",
			EventDescription:        "zxc",
			Date:                    time.Time{},
			NotifyDate:              time.Time{},
			Cafe:                    domain.Cafe{},
			Table:                   domain.Table{},
			Event:                   domain.Event{},
			User:                    domain.User{},
			HoursUntilReservation:   0,
			MinutesUntilReservation: 0,
			IsActive:                false,
		},
		{
			ID:                      2,
			PartySize:               2,
			CustName:                "asd",
			CustMobile:              "dsa",
			CustEmail:               "zxc",
			EventDescription:        "zxc",
			Date:                    time.Time{},
			NotifyDate:              time.Time{},
			Cafe:                    domain.Cafe{},
			Table:                   domain.Table{},
			Event:                   domain.Event{},
			User:                    domain.User{},
			HoursUntilReservation:   0,
			MinutesUntilReservation: 0,
			IsActive:                false,
		},
		{
			ID:                      3,
			PartySize:               3,
			CustName:                "asd",
			CustMobile:              "dsa",
			CustEmail:               "zxc",
			EventDescription:        "zxc",
			Date:                    time.Time{},
			NotifyDate:              time.Time{},
			Cafe:                    domain.Cafe{},
			Table:                   domain.Table{},
			Event:                   domain.Event{},
			User:                    domain.User{},
			HoursUntilReservation:   0,
			MinutesUntilReservation: 0,
			IsActive:                false,
		},
	}

	for _, v := range testRR {
		myR := reservationsPool.Get().(*domain.Reservation)
		*myR = v
		rr = append(rr, *myR)

		myR.Reset()
		reservationsPool.Put(myR)
	}
}

func WithPoolAssigning() {
	var rr []domain.Reservation
	testRR := []domain.Reservation{
		{
			ID:                      1,
			PartySize:               1,
			CustName:                "asd",
			CustMobile:              "dsa",
			CustEmail:               "zxc",
			EventDescription:        "zxc",
			Date:                    time.Time{},
			NotifyDate:              time.Time{},
			Cafe:                    domain.Cafe{},
			Table:                   domain.Table{},
			Event:                   domain.Event{},
			User:                    domain.User{},
			HoursUntilReservation:   0,
			MinutesUntilReservation: 0,
			IsActive:                false,
		},
		{
			ID:                      2,
			PartySize:               2,
			CustName:                "asd",
			CustMobile:              "dsa",
			CustEmail:               "zxc",
			EventDescription:        "zxc",
			Date:                    time.Time{},
			NotifyDate:              time.Time{},
			Cafe:                    domain.Cafe{},
			Table:                   domain.Table{},
			Event:                   domain.Event{},
			User:                    domain.User{},
			HoursUntilReservation:   0,
			MinutesUntilReservation: 0,
			IsActive:                false,
		},
		{
			ID:                      3,
			PartySize:               3,
			CustName:                "asd",
			CustMobile:              "dsa",
			CustEmail:               "zxc",
			EventDescription:        "zxc",
			Date:                    time.Time{},
			NotifyDate:              time.Time{},
			Cafe:                    domain.Cafe{},
			Table:                   domain.Table{},
			Event:                   domain.Event{},
			User:                    domain.User{},
			HoursUntilReservation:   0,
			MinutesUntilReservation: 0,
			IsActive:                false,
		},
	}

	for _, v := range testRR {
		myR := reservationsPool.Get().(*domain.Reservation)
		*myR = v

		rr = append(rr, *myR)

		*myR = domain.Reservation{}
		reservationsPool.Put(myR)
	}
}

func WithoutPool() {
	var rr []domain.Reservation
	testRR := []domain.Reservation{
		{
			ID:                      1,
			PartySize:               1,
			CustName:                "asd",
			CustMobile:              "dsa",
			CustEmail:               "zxc",
			EventDescription:        "zxc",
			Date:                    time.Time{},
			NotifyDate:              time.Time{},
			Cafe:                    domain.Cafe{},
			Table:                   domain.Table{},
			Event:                   domain.Event{},
			User:                    domain.User{},
			HoursUntilReservation:   0,
			MinutesUntilReservation: 0,
			IsActive:                false,
		},
		{
			ID:                      2,
			PartySize:               2,
			CustName:                "asd",
			CustMobile:              "dsa",
			CustEmail:               "zxc",
			EventDescription:        "zxc",
			Date:                    time.Time{},
			NotifyDate:              time.Time{},
			Cafe:                    domain.Cafe{},
			Table:                   domain.Table{},
			Event:                   domain.Event{},
			User:                    domain.User{},
			HoursUntilReservation:   0,
			MinutesUntilReservation: 0,
			IsActive:                false,
		},
		{
			ID:                      3,
			PartySize:               3,
			CustName:                "asd",
			CustMobile:              "dsa",
			CustEmail:               "zxc",
			EventDescription:        "zxc",
			Date:                    time.Time{},
			NotifyDate:              time.Time{},
			Cafe:                    domain.Cafe{},
			Table:                   domain.Table{},
			Event:                   domain.Event{},
			User:                    domain.User{},
			HoursUntilReservation:   0,
			MinutesUntilReservation: 0,
			IsActive:                false,
		},
	}

	for _, v := range testRR {
		myR := reservationsPool.Get().(*domain.Reservation)
		*myR = v
		rr = append(rr, *myR)
	}
}
