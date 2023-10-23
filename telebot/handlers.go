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

func sendErrorMsgWithBot(b *tele.Bot, s *Subscription, err error) {
	log.Error().Err(err).Send()
	_, err = b.Send(&tele.User{ID: s.UserID}, "Something went wrong :( Try again")
	if err != nil {
		log.Error().Err(err).Send()
	}
}

func sendSubscriptionWeatherForecast(
	b *tele.Bot,
	s Subscription,
	weatherAPI weatherapi.WeatherAPI,
) {
	lat, lon := s.Event.Location.Lat, s.Event.Location.Lon

	latStr, lonStr :=
		geoutils.FormatCoordinateToString(lat),
		geoutils.FormatCoordinateToString(lon)

	weatherForecast, err := weatherAPI.GetWeatherForecast(latStr, lonStr)
	if err != nil {
		sendErrorMsgWithBot(b, &s, err)
	}

	weatherTextMsg, err := weatherForecast.Daily[0].FormatToTextMsg()
	if err != nil {
		log.Warn().Err(err).Send()
	}

	_, err = b.Send(&tele.User{ID: s.UserID}, weatherTextMsg)
	if err != nil {
		sendErrorMsgWithBot(b, &s, err)
	}
}

func handleSubscriptionWeatherForecast(
	t time.Time,
	b *tele.Bot,
	dbClient db.DatabaseAccessor,
	weatherAPI weatherapi.WeatherAPI) {

	tickerTimeUTCFormatted := t.UTC().Format(timeutils.Layout24H)

	subscriptions := FindProcessedSubscriptionsForTime(dbClient, tickerTimeUTCFormatted)

	for _, s := range subscriptions {
		go sendSubscriptionWeatherForecast(b, s, weatherAPI)
	}
}

func runTickerForSubscriptionWeatherForecast(
	b *tele.Bot,
	ticker *time.Ticker,
	dbClient db.DatabaseAccessor,
	weatherAPI weatherapi.WeatherAPI) {

	for t := range ticker.C {
		go handleSubscriptionWeatherForecast(t, b, dbClient, weatherAPI)
	}
}

func handleStartCmd(c tele.Context) error {
	return c.Send(`Welcome to <b>Weather Bot</b>!
Use this bot to see Weather Forecast in your area :)
Just send <i>location pin</i> to Weather bot and get accurate forecast for up to <b>7 days</b> forward!`)
}

func handleSubscriptionCmd(c tele.Context, dbClient db.DatabaseAccessor, cfg *config.DbCfg) error {
	if err := RequestSubscription(dbClient, cfg, c.Sender().ID); err != nil {
		return c.Send(`Something went wrong :( Try again later'`)
	}

	return c.Send(`To subscribe send time in 24h format (e.g 15:04)`)
}

func handleTimeMessageForSubscription(c tele.Context, dbClient db.DatabaseAccessor, cfg *config.Config) error {
	subscription, err := CheckSubscriptionExist(dbClient, &cfg.Db, c.Sender().ID)
	if err != nil || subscription.Event.Processed {
		return nil
	}

	_, err = timeutils.ParseTimeFormat(c.Message().Text)
	if err != nil {
		return c.Send("Time format is invalid or unsupported, try again with different format")
	}

	timeUpdate := bson.M{
		"event.time": c.Message().Text,
	}

	if err = UpdateSubscription(dbClient, c.Sender().ID, timeUpdate, &cfg.Db); err != nil {
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

	subscription, err := CheckSubscriptionExist(dbClient, &cfg.Db, c.Sender().ID)

	if err == nil && subscription != nil {
		isSubscriptionValid := !subscription.Event.Processed && subscription.Event.RecurringTime != ""

		if isSubscriptionValid {
			locUpdate := bson.M{
				"event.location.lat": lat,
				"event.location.lon": lon,
				"event.processed":    true,
			}

			if err = UpdateSubscription(dbClient, c.Sender().ID, locUpdate, &cfg.Db); err != nil {
				return c.Send("Failed to subscribe :( Try again later")
			}

			return c.Send("Subscription's active now!")
		}
	}

	weatherForecast, err := weatherAPI.GetWeatherForecast(latStr, lonStr)
	if err != nil {
		log.Error().Err(err).Send()
		return c.Send("Something went wrong :( Try again")
	}

	if err = weatherForecast.StoreUpdateWeatherForecast(dbClient, &cfg.Db, c.Sender().ID); err != nil {
		return c.Send("Something went wrong :( Try again")
	}

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
