package db

import (
	"git.foxminded.ua/foxstudent106092/weather-bot/config"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

type NestedItemToInsert struct {
	Field1 int `bson:"field_one"`
	Field2 int `bson:"field_two"`
}

type ItemToInsert struct {
	UserID  int64              `bson:"user_id"`
	ItemObj NestedItemToInsert `bson:"item_obj"`
}

func TestDbCRUDFunctions(t *testing.T) {
	cfg, err := config.GetConfig()
	if err != nil {
		t.Errorf("error initalizing config %s", err)
	}

	dbClient, err := NewDatabaseClient(&cfg.Db)
	if err != nil {
		t.Errorf("error connecting client to db server %s", err)
	}

	var itemToInsert = ItemToInsert{
		UserID:  123456,
		ItemObj: NestedItemToInsert{5, 10},
	}

	t.Run("CREATE", func(t *testing.T) {
		if err = dbClient.InsertItemToDB(itemToInsert, cfg.Db.ForecastCollectionName); err != nil {
			t.Errorf("error inserting item to db %s", err)
		}
	})

	t.Run("READ", func(t *testing.T) {
		var itemToRead ItemToInsert
		if err = dbClient.FindUserInDB(itemToInsert.UserID, cfg.Db.ForecastCollectionName).Decode(&itemToRead); err != nil {
			t.Errorf("error reading item from db %s", err)
		}

		t.Logf("%+v\n", itemToRead)
	})

	t.Run("UPDATE", func(t *testing.T) {
		update := bson.M{"$set": bson.M{
			"item_obj.field_one": 2,
		}}

		if err = dbClient.UpdateItemInDB(itemToInsert.UserID, update, cfg.Db.ForecastCollectionName); err != nil {
			t.Errorf("error updating item to db %s", err)
		}
	})
}
