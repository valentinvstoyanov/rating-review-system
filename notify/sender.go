package notify

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"log"
)

const (
	Topic      = "arn:aws:sns:us-east-1:299047036407:rating-alerts"
	RegionName = "us-east-1"
)

func SendPushNotification(message string) error {
	// Initialize a session in REGION_NAME that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(RegionName)},
	)

	if err != nil {
		log.Printf("Unable to send notification %s to topic %s\n", message, Topic)
		return err
	}

	svc := sns.New(sess)
	topic := Topic

	result, err := svc.Publish(&sns.PublishInput{
		Message:  &message,
		TopicArn: &topic,
	})

	if err != nil {
		log.Printf("Failed to send push notification to topic %s: %s\n", Topic, err.Error())
		return err
	}

	log.Printf("Sent push notification to topic %s with messageId %s\n", Topic, *result.MessageId)
	return nil
}
