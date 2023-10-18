package weatherbotdb

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

type WeatherBotClientDb struct {
	cfg    *config.DbCfg
	client *mongo.Client
}

type DbAccessor interface {
	InsertDocToDbCollection(collectionName string) error
	CloseConnectionToDb() error
}

func connectToDb(cfg *config.DbCfg) (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	opts := options.Client().ApplyURI(cfg.ConnectionURL).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}

	// Send a ping to confirm a successful connection
	if err := client.Database(cfg.Name).RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		return nil, err
	}

	log.Info().
		Str("service", "ConnectToMongoDB").
		Msg("Pinged your deployment. You successfully connected to MongoDB!")

	return client, nil
}

func NewWeatherBotDbClient(cfg *config.DbCfg) (*WeatherBotClientDb, error) {
	client, err := connectToDb(cfg)
	if err != nil {
		log.Error().Stack().
			Err(err).
			Str("service", "NewWeatherBotDbClient").
			Msg("failed to create new instance of WeatherBotClientDb")

		return nil, err
	}

	return &WeatherBotClientDb{
		cfg:    cfg,
		client: client,
	}, nil
}

func (wbd *WeatherBotClientDb) InsertDocToDbCollection(doc interface{}, collectionName string) error {
	collection := wbd.client.Database(wbd.cfg.Name).Collection(collectionName)

	if collection == nil {
		err := errors.New(
			fmt.Sprintf("error: weatherbotdb client failed to find collection \"%s\"", collectionName),
		)

		return err
	}

	defer wbd.client.Disconnect(context.TODO())

	_, err := collection.InsertOne(context.TODO(), doc)
	if err != nil {
		return err
	}

	return nil
}
