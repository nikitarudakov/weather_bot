package telebot

import (
	"context"
	"git.foxminded.ua/foxstudent106092/weather-bot/config"
	"git.foxminded.ua/foxstudent106092/weather-bot/db"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	tele "gopkg.in/telebot.v3"
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
}

// SubscriptionService stores userID and subscription event
type SubscriptionService struct {
	UserID    int64             `bson:"user_id"`
	UserObj   tele.User         `bson:"user"`
	Event     SubscriptionEvent `bson:"event"`
	Processed bool              `bson:"processed"`
}

// FindProcessedSubscriptions searches for all subscription in db that
// has been processed and being active
func FindProcessedSubscriptions(dbClient db.DatabaseAccessor) []SubscriptionService {
	filter := bson.D{{"processed", true}}

	var results []SubscriptionService
	cursor, err := dbClient.FindItemsInDB(filter)
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Warn().Err(err).Send()
		return nil
	}

	return results
}

// CheckSubscriptionExist checks weather subscription item
// of user with userID is currently present in db
func CheckSubscriptionExist(dbClient db.DatabaseAccessor, dbCfg *config.DbCfg, userID int64) (*SubscriptionService, error) {
	var subService SubscriptionService
	if err := dbClient.FindUserInDB(userID, dbCfg.SubsCollectionName).Decode(&subService); err != nil {
		return nil, err
	}

	return &subService, nil
}

// RequestSubscription inserts initial subscription item for userID unless such item is already present in db,
// in that case it updates that item with new SubscriptionEvent and processed status
func RequestSubscription(dbClient db.DatabaseAccessor, dbCfg *config.DbCfg, userID int64, userOBJ tele.User) error {
	subscriptionService, err := CheckSubscriptionExist(dbClient, dbCfg, userID)
	if err == nil && !subscriptionService.Processed {
		return nil
	}

	subService := &SubscriptionService{
		UserID:    userID,
		UserObj:   userOBJ,
		Event:     SubscriptionEvent{},
		Processed: false,
	}

	if err == nil && subscriptionService.Processed {
		update := bson.M{
			"$set": bson.M{
				"event":     SubscriptionEvent{},
				"processed": false,
			},
		}

		if err = dbClient.UpdateItemInDB(userID, update, dbCfg.SubsCollectionName); err != nil {
			return err
		}

		return nil
	}

	if err = dbClient.InsertItemToDB(subService, dbCfg.SubsCollectionName); err != nil {
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
