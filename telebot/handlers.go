package telebot

import (
	"git.foxminded.ua/foxstudent106092/weather-bot/config"
	"git.foxminded.ua/foxstudent106092/weather-bot/db"
	"git.foxminded.ua/foxstudent106092/weather-bot/utils/geoutils"
	"git.foxminded.ua/foxstudent106092/weather-bot/utils/timeutils"
	"git.foxminded.ua/foxstudent106092/weather-bot/weatherapi"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	tele "gopkg.in/telebot.v3"
	"strings"
	"time"
)

func sendErrorMsgWithBot(b *tele.Bot, subService *SubscriptionService, err error) {
	log.Error().Err(err).Send()
	_, err = b.Send(&subService.UserObj, "Something went wrong :( Try again")
	if err != nil {
		log.Error().Err(err).Send()
	}
}

func HandleWeatherForecastSubscription(
	b *tele.Bot,
	dbClient db.DatabaseAccessor,
	subService SubscriptionService,
	weatherAPI weatherapi.WeatherAPI,
	dbCfg *config.DbCfg,
) {
	lat, lon := subService.Event.Location.Lat, subService.Event.Location.Lon

	latStr, lonStr :=
		geoutils.FormatCoordinateToString(lat),
		geoutils.FormatCoordinateToString(lon)

	weatherForecast, err := weatherAPI.GetWeatherForecast(latStr, lonStr)
	if err != nil {
		sendErrorMsgWithBot(b, &subService, err)
	}

	if err = weatherForecast.StoreWeatherForecastForUser(dbClient, dbCfg, subService.UserID); err != nil {
		sendErrorMsgWithBot(b, &subService, err)
	}

	msg := "Weather forecast is ready for you, just click date button in menu"
	_, err = b.Send(&subService.UserObj, msg, menu)
	if err != nil {
		sendErrorMsgWithBot(b, &subService, err)
	}
}

func RecurrentWeatherForecast(
	b *tele.Bot,
	ticker *time.Ticker,
	dbClient db.DatabaseAccessor,
	dbCfg *config.DbCfg,
	weatherAPI weatherapi.WeatherAPI) {

	for t := range ticker.C {
		tickerTimeUTCFormatted := t.UTC().Format(timeutils.Layout24H)

		subscriptions := FindProcessedSubscriptions(dbClient)

		for _, subscriptionService := range subscriptions {
			isRecurrentTime := subscriptionService.Event.RecurringTime == tickerTimeUTCFormatted

			if subscriptionService.Processed && isRecurrentTime {
				go HandleWeatherForecastSubscription(b, dbClient, subscriptionService, weatherAPI, dbCfg)
			}
		}
	}
}

func handleStartCmd(c tele.Context) error {
	return c.Send(`Welcome to <b>Weather Bot</b>!
Use this bot to see Weather Forecast in your area :)
Just send <i>location pin</i> to Weather bot and get accurate forecast for up to <b>7 days</b> forward!`)
}

func handleSubscriptionCmd(c tele.Context, dbClient db.DatabaseAccessor, cfg *config.DbCfg) error {
	if err := RequestSubscription(dbClient, cfg, c.Sender().ID, *c.Sender()); err != nil {
		return c.Send(`Something went wrong :( Try again later'`)
	}

	return c.Send(`To subscribe send time in 24h format (e.g 15:04)`)
}

func handleTimeMessageForSubscription(c tele.Context, dbClient db.DatabaseAccessor, cfg *config.Config) error {
	subscriptionService, err := CheckSubscriptionExist(dbClient, &cfg.Db, c.Sender().ID)
	if err != nil || subscriptionService.Processed {
		return nil
	}

	_, err = timeutils.ParseTimeFormat(c.Message().Text)
	if err != nil {
		return c.Send("Time format is invalid or unsupported, try again with different format")
	}

	timeUpdate := bson.M{
		"event.time": c.Message().Text,
	}

	if err = UpdateSubscription(dbClient, c.Sender().ID, timeUpdate); err != nil {
		return c.Send("Ahh... It didn't work... Check input time and try again")
	}

	return c.Send(`Good! You are one step closer :)
Now all I need to know  is area to look up forecast at.
Send Location Pin to proceed`)
}

func handleLocationPinMessage(
	c tele.Context,
	cfg *config.Config,
	dbClient db.DatabaseAccessor,
	weatherAPI weatherapi.WeatherAPI,
) error {
	lat, lon := c.Message().Location.Lat, c.Message().Location.Lng

	latStr, lonStr :=
		geoutils.FormatCoordinateToString(lat),
		geoutils.FormatCoordinateToString(lon)

	subscriptionService, err := CheckSubscriptionExist(dbClient, &cfg.Db, c.Sender().ID)
	if err == nil && !subscriptionService.Processed && subscriptionService.Event.RecurringTime != "" {
		locUpdate := bson.M{
			"event.location.lat": lat,
			"event.location.lon": lon,
			"processed":          true,
		}

		if err = UpdateSubscription(dbClient, c.Sender().ID, locUpdate); err != nil {
			return c.Send("Failed to subscribe :( Try again later")
		}

		return c.Send("Subscription's active now!")
	}

	weatherForecast, err := weatherAPI.GetWeatherForecast(latStr, lonStr)
	if err != nil {
		log.Error().Err(err).Send()
		return c.Send("Something went wrong :( Try again")
	}

	log.Info().Msg("GOT WEATHER FORECAST")

	if err = weatherForecast.StoreWeatherForecastForUser(dbClient, &cfg.Db,
		subscriptionService.UserID); err != nil {
		return c.Send("Something went wrong :( Try again")
	}

	log.Info().Msg("STORED WEATHER FORECAST")

	return c.Send("Weather forecast is ready for you, just click date button in menu", menu)
}

func handleDateWholePeriodBtn(c tele.Context, dbClient db.DatabaseAccessor, cfg *config.Config, weatherAPI weatherapi.WeatherAPI) error {
	weatherForecast, err := weatherAPI.ReadWeatherForecastFromDB(dbClient, &cfg.Db, c.Sender().ID)
	if err != nil {
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

func handleDateBtn(c tele.Context, dbClient db.DatabaseAccessor, cfg *config.Config,
	weatherAPI weatherapi.WeatherAPI, btnMenuIndex int) error {

	weatherForecast, err := weatherAPI.ReadWeatherForecastFromDB(dbClient, &cfg.Db, c.Sender().ID)
	if err != nil {
		return c.Send("Data is unavailable!\nSend location pin")
	}

	weatherTextMsg, err := weatherForecast.Daily[btnMenuIndex-1].FormatToTextMsg()
	if err != nil {
		log.Warn().Err(err).Send()
	}

	return c.Send(weatherTextMsg)
}
