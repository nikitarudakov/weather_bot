package telebot

import (
	"git.foxminded.ua/foxstudent106092/weather-bot/db"
)

// SubscriptionEvent stores time to send forecast at
type SubscriptionEvent struct {
	recurringTime string `bson:"time"`
}

// SubscriptionService stores userID and subscription event
type SubscriptionService struct {
	UserID    int64             `bson:"user_id"`
	Event     SubscriptionEvent `bson:"event"`
	Processed bool              `bson:"processed"`
}

type SubscriptionManager interface {
	CheckSubscription(dbClient db.DatabaseAccessor) error
	RequestSubscription(dbClient db.DatabaseAccessor) error
	UpdateSubscription(dbClient db.DatabaseAccessor) error
}

func NewSubscriptionService(userID int64, time string, processed bool) SubscriptionManager {
	var subService SubscriptionManager = &SubscriptionService{
		UserID:    userID,
		Event:     SubscriptionEvent{time},
		Processed: processed,
	}

	return subService
}

func (subService *SubscriptionService) CheckSubscription(dbClient db.DatabaseAccessor) error {
	if err := dbClient.FindItemInDb(subService.UserID); err != nil {
		return err
	}

	return nil
}

func (subService *SubscriptionService) RequestSubscription(dbClient db.DatabaseAccessor) error {
	if err := dbClient.InsertItemToDB(subService); err != nil {
		return err
	}

	return nil
}

func (subService *SubscriptionService) UpdateSubscription(dbClient db.DatabaseAccessor) error {
	if err := dbClient.UpdateItemInDb(subService.UserID, subService.Event.recurringTime, subService.Processed); err != nil {
		return err
	}

	return nil
}
