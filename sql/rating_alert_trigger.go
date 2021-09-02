package sql

import (
	"errors"
	"fmt"
	rrs "github.com/valentinvstoyanov/rating-review-system"
	"github.com/valentinvstoyanov/rating-review-system/slack"
	"log"
	"math"
	"time"
)

type PersistentRatingAlertTriggerService struct {
	entityService      rrs.EntityService
	reviewService      rrs.ReviewService
	ratingAlertService rrs.RatingAlertService
}

func NewRatingAlertTriggerService(entityService rrs.EntityService, reviewService rrs.ReviewService, ratingAlertService rrs.RatingAlertService) *PersistentRatingAlertTriggerService {
	return &PersistentRatingAlertTriggerService{entityService, reviewService, ratingAlertService}
}

func (s *PersistentRatingAlertTriggerService) Trigger(ratingAlert *rrs.RatingAlert) (*rrs.RatingAlert, error) {
	entity, err := s.entityService.GetById(ratingAlert.EntityId)
	if err != nil {
		return nil, errors.New("failed to find entity with such id: " + err.Error())
	}

	now := time.Now()
	startTime, endTime := now.Add(-time.Minute*time.Duration(ratingAlert.PeriodMinutes)), now

	if err := s.ratingAlertService.UpdateLastTriggeredAtById(ratingAlert.Id, now); err == nil {
		ratingAlert.LastTriggeredAt = now
	}

	reviews := s.reviewService.GetByEntityIdInPeriod(entity.Id, startTime, endTime)
	reviewsCount := len(reviews)

	log.Printf("Checking for rating alarm with reviews: %#v\n", reviews)

	if reviewsCount == 0 {
		return ratingAlert, nil
	}

	firstReview, lastReview := reviews[0], reviews[reviewsCount-1]
	percentageChange, period := calculateRatingPercentageChange(&firstReview, &lastReview), lastReview.CreatedAt.Sub(firstReview.CreatedAt)

	if percentageChange >= ratingAlert.PercentageChange {
		message := fmt.Sprintf("Rating alarm for entity '%s' with id %d: percentage change is %f over %s period\n", entity.Name, entity.Id, percentageChange, period)
		if err := slack.Send(message); err != nil {
			log.Printf("Failed to send Slack notification, message=%s, err=%s\n", message, err)
		} else {
			log.Printf("Sent Slack notification: %s", message)
		}
	}

	return ratingAlert, nil
}

func calculateRatingPercentageChange(firstReview *rrs.Review, lastReview *rrs.Review) float32 {
	ratingDiff := lastReview.Rating - firstReview.Rating
	return float32(math.Abs(float64(ratingDiff/firstReview.Rating)) * 100)
}
