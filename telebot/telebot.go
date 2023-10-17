package telebot

import (
	"git.foxminded.ua/foxstudent106092/weather-bot/config"
	"git.foxminded.ua/foxstudent106092/weather-bot/utils/geoutils"
	"git.foxminded.ua/foxstudent106092/weather-bot/weatherapi"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	tele "gopkg.in/telebot.v3"
	"strings"
	"time"
)

var menu = &tele.ReplyMarkup{ResizeKeyboard: true}
var lastLatStored, lastLonStored string
var weatherForecast *weatherapi.Response

func getDtBtnSlice() []tele.Btn {
	var dtBtnSlice []tele.Btn

	dtBtnSlice = append(dtBtnSlice, menu.Text("All 7 days"))

	dtToday := time.Now()

	for dayPlus := 0; dayPlus < 7; dayPlus++ {
		dtStr := dtToday.AddDate(0, 0, dayPlus).Format("02/01/2006")
		dtBtnSlice = append(dtBtnSlice, menu.Text(dtStr))
	}

	return dtBtnSlice
}

// HandleDtBtn initialize handler for menu btn
// that send back weather forecast on specified date
func HandleDtBtn(b *tele.Bot, dtBtnSlice []tele.Btn, dtBtnIndex int) {
	b.Handle(&dtBtnSlice[dtBtnIndex], func(c tele.Context) error {
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
	})
}

// InitTelegramBot initializes Telegram Weather Bot
func InitTelegramBot(cfg *config.Config) {
	pref := tele.Settings{
		Token:     viper.GetString("telegram.token"),
		Poller:    &tele.LongPoller{Timeout: 10 * time.Second},
		ParseMode: "HTML",
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Error().
			Str("service", "NewBot").
			Err(err).
			Msg("error initializing bot")

		return
	}

	log.Info().Msg("telegram bot was successfully initialized")

	b.Handle("/start", func(c tele.Context) error {
		return c.Send("Welcome to <b>Weather Bot</b>!\nUse this bot to see Weather Forecast in your area :)\nJust send <i>location pin</i> to Weather bot and get accurate forecast for up to <b>7 days</b> forward!")
	})

	// Create menu with dates starting today and ending on the day 7 days ahead
	dtBtnSlice := getDtBtnSlice()

	menu.Reply(
		menu.Row(dtBtnSlice[0]),
		menu.Row(dtBtnSlice[1], dtBtnSlice[2]),
		menu.Row(dtBtnSlice[3], dtBtnSlice[4]),
		menu.Row(dtBtnSlice[5], dtBtnSlice[6]),
		menu.Row(dtBtnSlice[7]),
	)
	// --------------------------------------------------------------------------

	b.Handle(tele.OnLocation, func(c tele.Context) error {
		lat, lon := c.Message().Location.Lat, c.Message().Location.Lng
		lastLatStored, lastLonStored =
			geoutils.FormatCoordinateToString(lat),
			geoutils.FormatCoordinateToString(lon)

		apiURL := weatherapi.GetAPIUrl(cfg, lastLatStored, lastLonStored)

		weatherForecast, err = weatherapi.GetWeatherForecast(apiURL)
		if err != nil {
			log.Error().
				Str("service", "GetWeatherForecast").
				Err(err).
				Msg("failed to get weather forecast")

			return c.Send("Data is unavailable for this location!")
		}

		return c.Send("Choose time period to get forecast:", menu)
	})

	b.Handle(&dtBtnSlice[0], func(c tele.Context) error {
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
	})

	for i := 1; i < 8; i++ {
		HandleDtBtn(b, dtBtnSlice, i)
	}

	b.Start()
}
