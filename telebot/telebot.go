package telebot

import (
	"fmt"
	"git.foxminded.ua/foxstudent106092/weather-bot/config"
	"git.foxminded.ua/foxstudent106092/weather-bot/db"
	"git.foxminded.ua/foxstudent106092/weather-bot/utils/timeutils"
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
		subscriptionService := NewSubscriptionService(c.Sender().ID, "", false)

		if err = subscriptionService.RequestSubscription(dbClient); err != nil {
			return err
		}

		return c.Send(`Send time in 24h (15:04) format to finish subscription`)
	})

	// handles event triggered by any sent text after subscribe cmd was triggered
	b.Handle(tele.OnText, func(c tele.Context) error {
		_, err = timeutils.ParseTimeFormat(c.Message().Text)
		if err != nil {
			return c.Send("Time format is invalid or unsupported, try again with different format")
		}

		subscriptionService := NewSubscriptionService(c.Sender().ID, c.Message().Text, true)

		if err = subscriptionService.CheckSubscription(dbClient); err != nil {
			return c.Send("Send subscribe command first!")
		}

		if err = subscriptionService.UpdateSubscription(dbClient); err != nil {
			return c.Send("Subscription was unsuccessful")
		}

		return c.Send(fmt.Sprintf("You were subscribed for time %s", c.Message().Text))
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
