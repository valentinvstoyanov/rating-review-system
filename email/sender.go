package email

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
	"log"
)

const (
	from      = "From"
	to        = "To"
	subject   = "Subject"
	PlainText = "text/plain"
)

type Message struct {
	Sender   string
	Receiver string
	Subject  string
	Content  string
}

func Send(message Message, password string) error {
	m := gomail.NewMessage()

	m.SetHeader(from, message.Sender)
	m.SetHeader(to, message.Receiver)
	m.SetHeader(subject, message.Subject)
	m.SetBody(PlainText, message.Content)

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, message.Sender, password)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email message %+v, err=%s\n", message, err)
		return err
	}

	return nil
}
