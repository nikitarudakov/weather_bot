package db

import (
	"fmt"
	"git.foxminded.ua/foxstudent106092/weather-bot/config"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDbConnection(t *testing.T) {
	cfg, err := config.GetConfig()
	if err != nil {
		t.Errorf("error parsing config %s", err)
	}

	_, err = NewDBClient(&cfg.Db)
	if err != nil {
		t.Error(err)
	}
}

func TestDbClean(t *testing.T) {
	cfg, err := config.GetConfig()
	if err != nil {
		t.Errorf("error parsing config %s", err)
	}

	dbClient, err := NewDBClient(&cfg.Db)
	if err != nil {
		t.Error(err)
	}

	if err := dbClient.CleanSubscriptionDataFromDB(); err != nil {
		t.Error(err)
	}
}

func TestDbInsertion(t *testing.T) {
	cfg, err := config.GetConfig()
	if err != nil {
		t.Fatalf("error parsing config %s", err)
	}

	dbClient, err := NewDBClient(&cfg.Db)
	if err != nil {
		t.Fatal(err)
	}

	doc := &SubscriptionData{
		Timestamp: time.Now().Unix(),
		SenderID:  123,
		MessageID: 123,
		ChatID:    123,
		Status:    1,
	}

	if err := dbClient.InsertSubscriptionDataToDB(doc); err != nil {
		t.Error(err)
	}

	subscriptionData, err := dbClient.FindSubscriptionDataInDB(doc.SenderID)
	if err != nil {
		t.Error(err)
	}

	t.Log(fmt.Sprintf("%+v\n", *subscriptionData))
	assert.Equal(t, doc, subscriptionData)

	if err := dbClient.CleanSubscriptionDataFromDB(); err != nil {
		t.Error(err)
	}
}
