package main

import (
	"git.foxminded.ua/foxstudent106092/weather-bot/config"
	"git.foxminded.ua/foxstudent106092/weather-bot/db"
	"git.foxminded.ua/foxstudent106092/weather-bot/logger"
	"git.foxminded.ua/foxstudent106092/weather-bot/telebot"
	"log"
)

func main() {
	// parse config to cfg
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	// initialize logger with config
	logger.InitLogger(cfg)

	// initialize db with db config
	dbClient, err := db.NewDatabaseClient(&cfg.Db)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	// close connection with db server
	defer dbClient.CloseConnectionToDB()

	// initialize telegram bot
	telebot.InitTelegramBot(cfg, dbClient)
}
