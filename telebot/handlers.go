package telebot

import (
	"fmt"
	"git.foxminded.ua/foxstudent106092/weather-bot/config"
	"git.foxminded.ua/foxstudent106092/weather-bot/db"
	"git.foxminded.ua/foxstudent106092/weather-bot/utils/geoutils"
	"git.foxminded.ua/foxstudent106092/weather-bot/utils/timeutils"
	"git.foxminded.ua/foxstudent106092/weather-bot/weatherapi"
	"github.com/rs/zerolog/log"
	tele "gopkg.in/telebot.v3"
	"strings"
)

type TelegramMessageService struct {
}

func handleStartCmd(c tele.Context) error {
	return c.Send(`Welcome to <b>Weather Bot</b>!
Use this bot to see Weather Forecast in your area :)
Just send <i>location pin</i> to Weather bot and get accurate forecast for up to <b>7 days</b> forward!`)
}

func handleSubscriptionCmd(c tele.Context, dbClient db.DatabaseAccessor) error {
	subscriptionService := NewSubscriptionService(c.Sender().ID, "", false, Location{})

	if err := subscriptionService.RequestSubscription(dbClient); err != nil {
		return c.Send(`Something went wrong :( Try again later'`)
	}

	return c.Send(`To subscribe send time in 24h format (e.g 15:04)`)
}

func handleTimeMessageForSubscription(c tele.Context, dbClient db.DatabaseAccessor) error {
	_, err := timeutils.ParseTimeFormat(c.Message().Text)
	if err != nil {
		return c.Send("Time format is invalid or unsupported, try again with different format")
	}

	subscriptionService := NewSubscriptionService(c.Sender().ID, c.Message().Text, true, Location{})

	if err = subscriptionService.CheckSubscriptionExist(dbClient); err != nil {
		return c.Send("Hit that /subscribe command first :)")
	}

	if err = subscriptionService.UpdateSubscription(dbClient); err != nil {
		return c.Send("Ahh... It didn't work... Check input time and try again")
	}

	return c.Send(fmt.Sprintf("Congrats! Subscription is now active (time: %s)", c.Message().Text))
}

func handleLocationPinMessage(c tele.Context, cfg *config.Config) error {
	lat, lon := c.Message().Location.Lat, c.Message().Location.Lng
	lastLatStored, lastLonStored =
		geoutils.FormatCoordinateToString(lat),
		geoutils.FormatCoordinateToString(lon)

	weatherAPI := weatherapi.NewWeatherService(&cfg.WeatherAPI)

	weatherForecastAtLocation, err := weatherAPI.GetWeatherForecast(lastLatStored, lastLonStored)
	if err != nil {
		log.Error().Err(err).Send()

		return c.Send("Data is unavailable for this location!")
	}

	weatherForecast = weatherForecastAtLocation

	return c.Send("Weather forecast is ready for you, just click date button in menu", menu)
}

func handleDateWholePeriodBtn(c tele.Context) error {
	if weatherForecast == nil {
		return c.Send("Data is unavailable!\nSend location pin")
	}

	var dailyWeatherBuilder strings.Builder
	for _, dailyWeather := range weatherForecast.Daily {
		weatherTextMsg, err := dailyWeather.FormatToTextMsg()
		if err != nil {
			log.Warn().Err(err).Send()
		}

		dailyWeatherBuilder.WriteString(weatherTextMsg)
		dailyWeatherBuilder.WriteString("\n")
	}

	return c.Send(dailyWeatherBuilder.String())
}

func handleDateBtn(c tele.Context, btnMenuIndex int) error {
	if weatherForecast == nil {
		return c.Send("Data is unavailable!\nSend location pin")
	}

	weatherTextMsg, err := weatherForecast.Daily[btnMenuIndex-1].FormatToTextMsg()
	if err != nil {
		log.Warn().Err(err).Send()
	}

	return c.Send(weatherTextMsg)
}
