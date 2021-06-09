package notificator

import (
	"crypto/tls"
	"fmt"
	"github.com/CyganFx/table-reservation/internal/app/config"
	"github.com/CyganFx/table-reservation/internal/domain"
	"github.com/pkg/errors"
	gomail "gopkg.in/mail.v2"
	"time"
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

func (n *notificator) CollaborationNotify(cafe domain.Cafe) error {
	m := gomail.NewMessage()
	// Settings for SMTP server
	d := gomail.NewDialer(n.Host, n.Port, n.From, n.Pass)
	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	// Set E-Mail sender
	m.SetHeader("From", n.From)

	// Set E-Mail receivers
	m.SetHeader("To", "duman_ishanov@mail.ru")
	// Set E-Mail subject
	m.SetHeader("Subject", "Collaboration")
	// Set E-Mail body.
	m.SetBody("text/plain", fmt.Sprintf(
		`You have new request for collaboration,
					Name: %s
					Email: %s
					Date: %v
					Time: %v`,
		cafe.Name, cafe.Email, time.Now().Format(dateLayout), time.Now().Format(timeLayoutWithoutSeconds)))

	fmt.Println("Message: ", m)
	if err := d.DialAndSend(m); err != nil {
		return errors.Wrap(err, "sending email")
	}

	return nil
}
