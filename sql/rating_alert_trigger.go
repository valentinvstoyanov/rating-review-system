package sql

import (
	"errors"
	"fmt"
	rrs "github.com/valentinvstoyanov/rating-review-system"
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

	if reviewsCount == 0 {
		return ratingAlert, nil
	}

	firstReview, lastReview := reviews[0], reviews[reviewsCount-1]
	percentageChange, period := calculateRatingPercentageChange(&firstReview, &lastReview), lastReview.CreatedAt.Sub(firstReview.CreatedAt)

	if percentageChange >= ratingAlert.PercentageChange {
		//TODO: Send alarm to Slack
		fmt.Printf("Rating alarm for entity '%s' with id %d: percentage change is %f over %s period\n", entity.Name, entity.Id, percentageChange, period)
	}

	return ratingAlert, nil
}

func calculateRatingPercentageChange(firstReview *rrs.Review, lastReview *rrs.Review) float32 {
	ratingDiff := lastReview.Rating - firstReview.Rating
	return (ratingDiff / firstReview.Rating) * 100
}
