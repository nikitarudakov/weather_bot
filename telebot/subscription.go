package telebot

import (
	"git.foxminded.ua/foxstudent106092/weather-bot/db"
)

type Location struct {
	Lat float64 `bson:"lat"`
	Lon float64 `bson:"lon"`
}

// SubscriptionEvent stores time to send forecast at
type SubscriptionEvent struct {
	RecurringTime string   `bson:"time"`
	Location      Location `bson:"location"`
}

// SubscriptionService stores userID and subscription event
type SubscriptionService struct {
	UserID    int64             `bson:"user_id"`
	Event     SubscriptionEvent `bson:"event"`
	Processed bool              `bson:"processed"`
}

type SubscriptionManager interface {
	CheckSubscriptionExist(dbClient db.DatabaseAccessor) error
	RequestSubscription(dbClient db.DatabaseAccessor) error
	UpdateSubscription(dbClient db.DatabaseAccessor) error
}

func NewSubscriptionService(userID int64, time string, processed bool, loc Location) SubscriptionManager {
	subscriptionEvent := SubscriptionEvent{
		time,
		loc,
	}

	var subService SubscriptionManager = &SubscriptionService{
		UserID:    userID,
		Event:     subscriptionEvent,
		Processed: processed,
	}

	return subService
}

func (subService *SubscriptionService) CheckSubscriptionExist(dbClient db.DatabaseAccessor) error {
	if err := dbClient.FindItemInDb(subService.UserID); err != nil {
		return err
	}

	return nil
}

func (subService *SubscriptionService) RequestSubscription(dbClient db.DatabaseAccessor) error {
	if err := subService.CheckSubscriptionExist(dbClient); err == nil {
		return nil
	}

	if err := dbClient.InsertItemToDB(subService); err != nil {
		return err
	}

	return nil
}

func (subService *SubscriptionService) UpdateSubscription(dbClient db.DatabaseAccessor) error {
	if err := dbClient.UpdateItemInDb(subService.UserID, subService.Event.RecurringTime, subService.Processed); err != nil {
		return err
	}

	return nil
}
