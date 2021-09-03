package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"gopkg.in/gomail.v2"
	"log"
	"net/http"
	"os"
)

const (
	slackMessageKey     = "text"
	slackWebhookVarName = "SLACK_WEBHOOK"
	emailPassVarName    = "EMAIL_PASS"
)

type NotificationEvent struct {
	Target  string `json:"target"`
	Message string `json:"message"`
}

func Handler(snsEvent events.SNSEvent) error {
	log.Printf("SNS event: %+v\n", snsEvent)
	var resErr error

	for _, record := range snsEvent.Records {
		snsMessage := record.SNS.Message
		event := NotificationEvent{}

		if err := json.Unmarshal([]byte(snsMessage), &event); err != nil {
			log.Printf("Failed to unmarshal sns message %s, err=%s\n", snsMessage, err)
			resErr = err
		}

		if event.Target == "slack" {
			if err := sendToSlack(event.Message); err != nil {
				log.Printf("Failed to send message to slack, err=%s\n", event.Message)
				resErr = err
			}

			continue
		}

		if event.Target == "email" {
			emailMessage := Message{}

			if err := json.Unmarshal([]byte(event.Message), &emailMessage); err != nil {
				log.Printf("Failed to unmarshal event message %s, err=%s\n", event.Message, err)
				resErr = err
			}

			emailMessage.Sender = "education.purposes.test.1914@gmail.com"
			if err := SendEmail(emailMessage, os.Getenv(emailPassVarName)); err != nil {
				log.Printf("Failed to send email message, err=%s\n", event.Message)
				resErr = err
			}

			continue
		}

		err := fmt.Sprintf("Unknown target %s, message: %s.", event.Target, event.Message)
		log.Printf(err)
		resErr = errors.New(err)
	}

	return resErr
}

func sendToSlack(message string) error {
	body, err := json.Marshal(map[string]string{slackMessageKey: message})
	if err != nil {
		return err
	}

	_, err = http.Post(os.Getenv(slackWebhookVarName), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	return nil
}

const (
	from      = "From"
	to        = "To"
	subject   = "Subject"
	PlainText = "text/plain"
)

type Message struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Subject  string `json:"subject"`
	Content  string `json:"content"`
}

func SendEmail(message Message, password string) error {
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

func (b *NotificationEvent) UnmarshalJSON(data []byte) error {
	m := map[string]interface{}{}

	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	target, ok := m["target"].(string)
	if !ok {
		return errors.New("Unexpected target type, expected string")
	}

	b.Target = target

	if target == "slack" {
		message, ok := m["message"].(string)
		if !ok {
			return errors.New("Unexpected message type, expected string")
		}

		b.Message = message
		return nil
	}

	if target == "email" {
		jsonBytes, err := json.Marshal(m["message"])
		if err != nil {
			return err
		}

		b.Message = string(jsonBytes)

		return nil
	}

	return errors.New(fmt.Sprintf("Unexpected target %s, expected slack or email"))
}

func main() {
	lambda.Start(Handler)
}
