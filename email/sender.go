package email

import (
	"encoding/json"
	"fmt"
	"github.com/valentinvstoyanov/rating-review-system/notify"
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

	notification := fmt.Sprintf("{\"target\": \"%s\", \"message\": %s}", "email", bytes)
	return notify.SendPushNotification(notification)
}
