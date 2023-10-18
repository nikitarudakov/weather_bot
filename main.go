package main

import (
	"git.foxminded.ua/foxstudent106092/weather-bot/config"
	"git.foxminded.ua/foxstudent106092/weather-bot/logger"
	"git.foxminded.ua/foxstudent106092/weather-bot/telebot"
	"git.foxminded.ua/foxstudent106092/weather-bot/weatherbotdb"
	"log"
)

func main() {
	// parse config to cfg
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Fatal error: %s", err)
	}

	// initialize logger with config
	logger.InitLogger(cfg)

	// initialize db with db config
	dbClient, err := weatherbotdb.NewWeatherBotDbClient(&cfg.Db)
	if err != nil {
		panic(err)
	}

	// initialize telegram bot
	telebot.InitTelegramBot(cfg, dbClient)
}
