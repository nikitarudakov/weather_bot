package telebot

import (
	"git.foxminded.ua/foxstudent106092/weather-bot/config"
	"git.foxminded.ua/foxstudent106092/weather-bot/db"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	tele "gopkg.in/telebot.v3"
	"time"
)

// InitTelegramBot initializes Telegram Weather Bot
func InitTelegramBot(cfg *config.Config, dbClient db.DatabaseAccessor) {
	pref := tele.Settings{
		Token:     viper.GetString("telegram.token"),
		Poller:    &tele.LongPoller{Timeout: 10 * time.Second},
		ParseMode: "HTML",
	}

	// connect to telegram bot
	b, err := tele.NewBot(pref)
	if err != nil {
		log.Error().Err(err).Send()

		return
	}

	log.Info().Msg("Telegram bot was successfully initialized")

	// handles event triggered by start command
	b.Handle("/start", func(c tele.Context) error {
		return c.Send(`Welcome to <b>Weather Bot</b>!
Use this bot to see Weather Forecast in your area :)
Just send <i>location pin</i> to Weather bot and get accurate forecast for up to <b>7 days</b> forward!`)
	})

	// handles event triggered by subscribe command
	b.Handle("/subscribe", func(c tele.Context) error {
		if err = handleSubscriptionDataInsertionToDB(c, dbClient, 1); err != nil {
			log.Error().Err(err).Send()

			return err
		}

		return c.Send(`To subscribe for daily weather forecast 
send time (format: 15:04) you wish to receive it at`)
	})

	// handles event triggered by any sent text
	b.Handle(tele.OnText, func(c tele.Context) error {
		return nil
	})

	// create menu with dates starting today and ending on the day 7 days ahead
	dtBtnSlice := getDtBtnSlice()

	menu.Reply(
		menu.Row(dtBtnSlice[0]),
		menu.Row(dtBtnSlice[1], dtBtnSlice[2]),
		menu.Row(dtBtnSlice[3], dtBtnSlice[4]),
		menu.Row(dtBtnSlice[5], dtBtnSlice[6]),
		menu.Row(dtBtnSlice[7]),
	)

	// handles event triggered when Location Pin is sent
	b.Handle(tele.OnLocation, func(c tele.Context) error {
		err = handleOnLocation(cfg, c)
		if err != nil {
			log.Error().Err(err).Send()
			return err
		}

		return nil
	})

	// handles event triggered when whole period btn is pressed on menu
	b.Handle(&dtBtnSlice[0], func(c tele.Context) error {
		err = handleWholePeriodBtn(c)
		if err != nil {
			log.Error().Err(err).Send()
			return err
		}

		return nil
	})

	// handles event triggered when day btn is pressed on menu
	for i := 1; i < 8; i++ {
		dtBtnIndex := i
		b.Handle(&dtBtnSlice[dtBtnIndex], func(c tele.Context) error {
			err = handleDateBtn(c, dtBtnIndex)
			if err != nil {
				return err
			}

			return nil
		})
	}

	b.Start()
}
