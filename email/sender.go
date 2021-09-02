package email

import (
	"encoding/json"
	"github.com/valentinvstoyanov/rating-review-system/notify"
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

func Send(message Message) error {
	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	m := map[string]string{"target": "email", "message": string(bytes)}
	notification, err := json.Marshal(m)
	if err != nil {
		return err
	}

	return notify.SendPushNotification(string(notification))
}
