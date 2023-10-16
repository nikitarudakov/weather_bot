package weatherapi

import (
	"fmt"
	"git.foxminded.ua/foxstudent106092/weather-bot/utils/timeutils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"time"
)

func (d *DailyWeather) FormatToTextMsg() string {
	dateStr := timeutils.ConvertUnixTimestampToDate(d.UnixDt, time.DateOnly)
	sunriseDateStr := timeutils.ConvertUnixTimestampToDate(d.Sunrise, time.Kitchen)
	sunsetDateStr := timeutils.ConvertUnixTimestampToDate(d.Sunset, time.Kitchen)

	weatherDescStr := ""

	if len(d.Weather) > 0 {
		var weatherDesc []string
		for _, weather := range d.Weather {
			weatherDesc = append(weatherDesc, weather.Desc)
		}

		caser := cases.Title(language.English)

		weatherDescStr = caser.String(d.Weather[0].Desc)
	}

	// round weather temperature
	return fmt.Sprintf("Weather on *%v*:\n%s: %s\nH:**%.0f**° L:**%.0f**°\nSunrise at %s\nSunset at %s\n",
		dateStr, weatherDescStr, d.Summary, d.Temp.Max, d.Temp.Min,
		sunriseDateStr, sunsetDateStr)
}
