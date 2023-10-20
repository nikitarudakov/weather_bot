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
)

// DatabaseClient represents a service that interacts with a bot db
type DatabaseClient struct {
	cfg    *config.DbCfg
	client *mongo.Client
}

// DatabaseAccessor represents the methods for interacting with a bot db
type DatabaseAccessor interface {
	InsertItemToDB(doc interface{}, collectionName string) error
	FindUserInDB(userID int64, collectionName string) *mongo.SingleResult
	UpdateItemInDB(userID int64, update bson.M) error
	CloseConnectionToDB() error
	FindItemsInDB(filter bson.D) (*mongo.Cursor, error)
}

func NewDatabaseClient(cfg *config.DbCfg) (DatabaseAccessor, error) {
	client, err := connectToDB(cfg)

	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}

	var databaseClient DatabaseAccessor = &DatabaseClient{
		cfg:    cfg,
		client: client,
	}

	return databaseClient, nil
}

func (dc *DatabaseClient) FindItemsInDB(filter bson.D) (*mongo.Cursor, error) {
	collection := dc.client.Database(dc.cfg.Name).Collection(dc.cfg.SubsCollectionName)

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	return cursor, nil
}

func (dc *DatabaseClient) UpdateItemInDB(userID int64, update bson.M) error {
	collection := dc.client.Database(dc.cfg.Name).Collection(dc.cfg.SubsCollectionName)
	filter := bson.D{{"user_id", userID}}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (dc *DatabaseClient) FindUserInDB(userID int64, collectionName string) *mongo.SingleResult {
	collection := dc.client.Database(dc.cfg.Name).Collection(collectionName)
	filter := bson.D{{"user_id", userID}}

	return collection.FindOne(context.TODO(), filter)
}

func (dc *DatabaseClient) InsertItemToDB(doc interface{}, collectionName string) error {
	collection := dc.client.Database(dc.cfg.Name).Collection(collectionName)

	if collection == nil {
		err := errors.New(
			fmt.Sprintf("error: db client failed to find collection \"%s\"", dc.cfg.SubsCollectionName),
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

func (dc *DatabaseClient) CloseConnectionToDB() error {
	err := dc.client.Disconnect(context.TODO())
	if err != nil {
		return err
	}

	return nil
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
