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
	InsertItemToDB(interface{}, string) error
	FindUserInDB(int64, string) *mongo.SingleResult
	FindItemsInDB(bson.D) (*mongo.Cursor, error)
	UpdateItemInDB(int64, bson.M, string) error
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

	log.Info().Msg("MongoDB client is running ...")

	return client, nil
}

// NewDatabaseClient establishes new connection with db server
// It creates databaseClient of type DatabaseAccessor
// which can be used for managing db
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

// InsertItemToDB inserts item (doc) of any structure to db collection collectionName
func (dc *DatabaseClient) InsertItemToDB(doc interface{}, collectionName string) error {
	collection := dc.client.Database(dc.cfg.Name).Collection(collectionName)

	if collection == nil {
		err := errors.New(
			fmt.Sprintf("error: db client failed to find collection \"%s\"", dc.cfg.SubsCollectionName),
		)
		return err
	}

	insertOneResult, err := collection.InsertOne(context.TODO(), doc)
	if err != nil {
		return err
	}

	log.Trace().Msg(fmt.Sprintf("%+v\n", *insertOneResult))

	return nil
}

// FindUserInDB finds item with user_id userID and returns *mongo.SingleResult
// that represents a single document returned from an operation
func (dc *DatabaseClient) FindUserInDB(userID int64, collectionName string) *mongo.SingleResult {
	collection := dc.client.Database(dc.cfg.Name).Collection(collectionName)
	filter := bson.D{{"user_id", userID}}

	return collection.FindOne(context.TODO(), filter)
}

// FindItemsInDB finds all items that matches filter variable
// and returns *mongo.Cursor that is used to iterate over a stream of documents
// and can be parsed
func (dc *DatabaseClient) FindItemsInDB(filter bson.D) (*mongo.Cursor, error) {
	collection := dc.client.Database(dc.cfg.Name).Collection(dc.cfg.SubsCollectionName)

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	return cursor, nil
}

// UpdateItemInDB updates item with user_id equals to userID in db collection collectionName,
// and updates its fields with update of type bson.M
func (dc *DatabaseClient) UpdateItemInDB(userID int64, update bson.M, collectionName string) error {
	collection := dc.client.Database(dc.cfg.Name).Collection(collectionName)
	filter := bson.D{{"user_id", userID}}

	log.Trace().Str("filter", fmt.Sprintf("%+v\n", filter)).Send()

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

// CloseConnectionToDB interrupts connection with db server
func (dc *DatabaseClient) CloseConnectionToDB() error {
	err := dc.client.Disconnect(context.TODO())
	if err != nil {
		return err
	}

	return nil
}
