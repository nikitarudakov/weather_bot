package main

import (
	"git.foxminded.ua/foxstudent106092/weather-bot/config"
	"git.foxminded.ua/foxstudent106092/weather-bot/logger"
	"git.foxminded.ua/foxstudent106092/weather-bot/telebot"
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

	// initialize telegram bot
	telebot.InitTelegramBot(cfg)
}
