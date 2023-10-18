package telebot

import (
	"git.foxminded.ua/foxstudent106092/weather-bot/config"
	"git.foxminded.ua/foxstudent106092/weather-bot/utils/geoutils"
	"git.foxminded.ua/foxstudent106092/weather-bot/weatherapi"
	"git.foxminded.ua/foxstudent106092/weather-bot/weatherbotdb"
	"github.com/rs/zerolog/log"
	tele "gopkg.in/telebot.v3"
	"strings"
	"time"
)

var lastLatStored, lastLonStored string
var weatherForecast *weatherapi.Response

type StoredMessages struct {
	Timestamp int64 `bson:"timestamp"`
	MessageID int   `bson:"message_id"`
	ChatID    int64 `bson:"chat_id"`
	IsCommand bool  `bson:"is_command"`
}

func insertChatHistoryToDb(
	c tele.Context,
	dbClient *weatherbotdb.WeatherBotClientDb,
	cfgDb *config.DbCfg,
) error {
	doc := StoredMessages{
		Timestamp: time.Now().Unix(),
		MessageID: c.Message().ID,
		ChatID:    c.Chat().ID,
		IsCommand: true,
	}

	if err := dbClient.InsertDocToDbCollection(doc, cfgDb.MsgCollectionName); err != nil {
		return err
	}

	return nil
}

func handleOnLocation(cfg *config.Config, c tele.Context) error {
	lat, lon := c.Message().Location.Lat, c.Message().Location.Lng
	lastLatStored, lastLonStored =
		geoutils.FormatCoordinateToString(lat),
		geoutils.FormatCoordinateToString(lon)

	weatherAPI := weatherapi.NewWeatherService(&cfg.WeatherAPI)

	weatherForecastAtLocation, err := weatherAPI.GetWeatherForecast(lastLatStored, lastLonStored)
	if err != nil {
		log.Error().
			Str("service", "GetWeatherForecast").
			Err(err).
			Msg("failed to get weather forecast")

		return c.Send("Data is unavailable for this location!")
	}

	weatherForecast = weatherForecastAtLocation

	return c.Send("Choose time period to get forecast:", menu)
}

func handleWholePeriodBtn(c tele.Context) error {
	if weatherForecast == nil {
		return c.Send("Data is unavailable!\nSend location pin")
	}

	var dailyWeatherBuilder strings.Builder
	for _, dailyWeather := range weatherForecast.Daily {
		weatherTextMsg, err := dailyWeather.FormatToTextMsg()
		if err != nil {
			log.Warn().
				Str("service", "FormatToTextMsg").
				Err(err).
				Msg("Warning! Forecast could not be formatted to text message")
		}

		dailyWeatherBuilder.WriteString(weatherTextMsg)
		dailyWeatherBuilder.WriteString("\n")
	}

	return c.Send(dailyWeatherBuilder.String())
}

func handleDateBtn(c tele.Context, dtBtnIndex int) error {
	if weatherForecast == nil {
		return c.Send("Data is unavailable!\nSend location pin")
	}

	weatherTextMsg, err := weatherForecast.Daily[dtBtnIndex-1].FormatToTextMsg()
	if err != nil {
		log.Warn().
			Str("service", "FormatToTextMsg").
			Err(err).
			Msg("Warning! Forecast could not be formatted to text message")
	}

	return c.Send(weatherTextMsg)
}
