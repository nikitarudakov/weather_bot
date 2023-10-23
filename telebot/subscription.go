package telebot

import (
	"context"
	"git.foxminded.ua/foxstudent106092/weather-bot/config"
	"git.foxminded.ua/foxstudent106092/weather-bot/db"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
)

// Location represents coordinates of location pin
type Location struct {
	Lat float32 `bson:"lat"`
	Lon float32 `bson:"lon"`
}

// SubscriptionEvent stores time to send forecast at
type SubscriptionEvent struct {
	RecurringTime string   `bson:"time"`
	Location      Location `bson:"location"`
	Processed     bool     `bson:"processed"`
}

// Subscription stores data about user and its Event object
type Subscription struct {
	UserID int64             `bson:"user_id"`
	Event  SubscriptionEvent `bson:"event"`
}

// FindProcessedSubscriptionsForTime searches for all subscription in db that
// has been processed and being active
func FindProcessedSubscriptionsForTime(dbClient db.DatabaseAccessor, time string) []Subscription {
	filter := bson.M{
		"event.time":      time,
		"event.processed": true,
	}

	var results []Subscription
	cursor, err := dbClient.FindItemsInDB(filter)
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Warn().Err(err).Send()
		return nil
	}

	return results
}

// CheckSubscriptionExist checks weather subscription item
// of user with userID is currently present in db
func CheckSubscriptionExist(dbClient db.DatabaseAccessor, dbCfg *config.DbCfg, userID int64) (*Subscription, error) {
	var subscription Subscription
	if err := dbClient.FindUserInDB(userID, dbCfg.SubsCollectionName).Decode(&subscription); err != nil {
		return nil, err
	}

	return &subscription, nil
}

// RequestSubscription inserts initial subscription item for userID unless such item is already present in db,
// in that case it updates that item with new SubscriptionEvent and processed status
func RequestSubscription(dbClient db.DatabaseAccessor, dbCfg *config.DbCfg, userID int64) error {
	subscription, err := CheckSubscriptionExist(dbClient, dbCfg, userID)
	if err == nil && !subscription.Event.Processed {
		return nil
	}

	s := &Subscription{
		UserID: userID,
		Event:  SubscriptionEvent{Processed: false},
	}

	if err == nil && subscription.Event.Processed {
		update := bson.M{
			"$set": bson.M{
				"event.processed": false,
			},
		}

		if err = dbClient.UpdateItemInDB(userID, update, dbCfg.SubsCollectionName); err != nil {
			return err
		}

		return nil
	}

	if err = dbClient.InsertItemToDB(s, dbCfg.SubsCollectionName); err != nil {
		return err
	}

	return nil
}

// UpdateSubscription updates that item with new updateBsonObj
func UpdateSubscription(dbClient db.DatabaseAccessor, userID int64, updateBsonObj bson.M, dbCfg *config.DbCfg) error {
	update := bson.M{
		"$set": updateBsonObj,
	}

	if err := dbClient.UpdateItemInDB(userID, update, dbCfg.SubsCollectionName); err != nil {
		return err
	}

	return nil
}
