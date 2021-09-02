package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/valentinvstoyanov/rating-review-system/notify"
	"net/http"
)

const (
	slackWebhookVarName = "SLACK_WEBHOOK"
	messageKey          = "text"
)

func Send(message string) error {
	m := map[string]string{"target": "slack", "message": message}
	bytes, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return notify.SendPushNotification(string(bytes))
}

func sendToUrl(url, message string) error {
	if len(url) == 0 {
		msg := fmt.Sprintf("Cannot send slack notification: missing %s environment variable.", slackWebhookVarName)
		return errors.New(msg)
	}

	body, err := json.Marshal(map[string]string{messageKey: message})
	if err != nil {
		return err
	}

	_, err = http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	return nil
}
