package notificator

import (
	"crypto/tls"
	"fmt"
	"github.com/CyganFx/table-reservation/internal/app/config"
	"github.com/CyganFx/table-reservation/internal/domain"
	"github.com/pkg/errors"
	gomail "gopkg.in/mail.v2"
)

const (
	timeLayoutWithoutSeconds = "15:04"
	dateLayout               = "2006-01-02"
)

type notificator struct {
	Host string
	Port int
	From string
	Pass string
}

func New(cfg config.Config) *notificator {
	return &notificator{
		Host: cfg.SMTP.Host,
		Port: cfg.SMTP.Port,
		From: cfg.SMTP.From,
		Pass: cfg.SMTP.Pass,
	}
}

// TODO add cafe address and phone, Test
func (n *notificator) UsersBooking(reservations []domain.Reservation) error {
	m := gomail.NewMessage()
	// Settings for SMTP server
	d := gomail.NewDialer(n.Host, n.Port, n.From, n.Pass)
	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	// Set E-Mail sender
	m.SetHeader("From", n.From)

	for _, data := range reservations {

		// Set E-Mail receivers
		m.SetHeader("To", data.CustEmail)
		// Set E-Mail subject
		m.SetHeader("Subject", "Reservation")
		// Set E-Mail body.
		m.SetBody("text/plain", fmt.Sprintf(
			`%s, we remind you about reservation:

					Date: %v
					Time: %v
					Place: %s

					Your visit is %v minutes away

					Thank you and look forward to seeing you again!
					
					With gratitude,
					Ez-Booking application`,
			data.CustName, data.Date.Format(dateLayout), data.Date.Format(timeLayoutWithoutSeconds),
			data.Cafe.Name, data.Date.Sub(data.NotifyDate).Minutes()))

		fmt.Println("Message: ", m)

		if err := d.DialAndSend(m); err != nil {
			return errors.Wrap(err, "sending email")
		}
	}

	return nil
}
