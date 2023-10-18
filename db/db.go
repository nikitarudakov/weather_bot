package db

import (
	"context"
	"errors"
	"fmt"
	"git.foxminded.ua/foxstudent106092/weather-bot/config"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	tele "gopkg.in/telebot.v3"
	"time"
)

// SubscriptionData represents data about subscriptions such as
// a time subscription request was sent, chat and message ids,
// user id and status of subscription
type SubscriptionData struct {
	Timestamp int64 `bson:"timestamp"`
	MessageID int   `bson:"message_id"`
	ChatID    int64 `bson:"chat_id"`
	SenderID  int64 `bson:"sender_id"`
	Status    int8  `bson:"status"`
}

func NewSubscriptionData(c tele.Context, status int8) *SubscriptionData {
	return &SubscriptionData{
		Timestamp: time.Now().Unix(),
		SenderID:  c.Sender().ID,
		MessageID: c.Message().ID,
		ChatID:    c.Chat().ID,
		Status:    status,
	}
}

// DatabaseClient represents a service that interacts with a bot db
type DatabaseClient struct {
	cfg    *config.DbCfg
	client *mongo.Client
}

// DatabaseAccessor represents the methods for interacting with a bot db
type DatabaseAccessor interface {
	InsertSubscriptionDataToDB(doc interface{}) error
	FindSubscriptionDataInDB(senderID int64) (*SubscriptionData, error)
	CleanSubscriptionDataFromDB() error
	CloseConnectionToDB() error
}

func connectToDB(cfg *config.DbCfg) (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	opts := options.Client().ApplyURI(cfg.ConnectionURL).SetServerAPIOptions(serverAPI)

	// create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}

	// send a ping to confirm a successful connection
	if err = client.Database(cfg.Name).RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		return nil, err
	}

	log.Info().Msg("MongoDB client was successfully connected to the server!")

	return client, nil
}

// NewDBClient initializes new DatabaseClient instance
func NewDBClient(cfg *config.DbCfg) (DatabaseAccessor, error) {
	client, err := connectToDB(cfg)

	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}

	var dbAccessor DatabaseAccessor = &DatabaseClient{
		cfg:    cfg,
		client: client,
	}

	return dbAccessor, nil
}

func (wbd *DatabaseClient) FindSubscriptionDataInDB(senderID int64) (*SubscriptionData, error) {
	collection := wbd.client.Database(wbd.cfg.Name).Collection(wbd.cfg.SubsCollectionName)
	filter := bson.D{{"sender_id", senderID}}

	var subscriptionData *SubscriptionData

	if err := collection.FindOne(context.TODO(), filter).Decode(&subscriptionData); err != nil {
		return nil, err
	}

	return subscriptionData, nil
}

func (wbd *DatabaseClient) CleanSubscriptionDataFromDB() error {
	collection := wbd.client.Database(wbd.cfg.Name).Collection(wbd.cfg.SubsCollectionName)

	_, err := collection.DeleteMany(context.TODO(), bson.D{})
	if err != nil {
		return err
	}

	return nil
}

// CloseConnectionToDB closes clients connection to the server
func (wbd *DatabaseClient) CloseConnectionToDB() error {
	err := wbd.client.Disconnect(context.TODO())
	if err != nil {
		return err
	}

	return nil
}

// InsertSubscriptionDataToDB inserts subscription data to db
// !!!!!!!ADD EXISTS METHOD TO INTERFACE AND CALL InsertSubscriptionDataToDB
// ONLY IF EXISTS METHOD RETURN FALSE
func (wbd *DatabaseClient) InsertSubscriptionDataToDB(doc interface{}) error {
	collection := wbd.client.Database(wbd.cfg.Name).Collection(wbd.cfg.SubsCollectionName)

	if collection == nil {
		err := errors.New(
			fmt.Sprintf("error: db client failed to find collection \"%s\"", wbd.cfg.SubsCollectionName),
		)
		return err
	}

	_, err := collection.InsertOne(context.TODO(), doc)
	if err != nil {
		return err
	}

	log.Info().Msg("Subscription was successfully stored in database!")

	return nil
}
