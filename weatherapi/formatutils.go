package weatherapi

import (
	"bytes"
	"git.foxminded.ua/foxstudent106092/weather-bot/utils/timeutils"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"math"
	"strings"
	"text/template"
	"time"
)

const tpl = `Weather on <b>{{.Date}}</b>:
<i>{{.WeatherDesc}}</i>: {{.Summary}}
H:<b>{{.MaxTemp}}</b>° L:<b>{{.MinTemp}}</b>°
Sunrise at <b>{{.SunriseDate}}</b> Sunset at <b>{{.SunsetDate}}</b>
`

// CombineWeatherDesc combines all weather descriptions
// throughout a day as the weather changes
func CombineWeatherDesc(weatherSlice []Weather) string {
	weatherDescStr := ""
	defaultDelimiter := ", "

	if len(weatherSlice) > 0 {
		var weatherDescBuilder strings.Builder

		for i, weather := range weatherSlice {
			weatherDescBuilder.WriteString(weather.Desc)

			if i != len(weatherSlice)-1 {
				weatherDescBuilder.WriteString(defaultDelimiter)
				continue
			}

			weatherDescBuilder.WriteString(" ")
		}

		casesTitleEnglish := cases.Title(language.English)
		weatherDescStr = casesTitleEnglish.String(weatherDescBuilder.String())
	}

	return weatherDescStr
}

// FormatToTextMsg method takes values of DailyWeather and return formatted string
func (d *DailyWeather) FormatToTextMsg() (string, error) {
	dateStr := timeutils.ConvertUnixTimestampToDate(d.UnixDt, "02/01/2006")
	sunriseDateStr := timeutils.ConvertUnixTimestampToDate(d.Sunrise, time.Kitchen)
	sunsetDateStr := timeutils.ConvertUnixTimestampToDate(d.Sunset, time.Kitchen)

	weatherDescStr := CombineWeatherDesc(d.Weather)

	numMax, numMin := int(math.Round(d.Temp.Max)), int(math.Round(d.Temp.Min))

	weatherMsgTemplate := template.Must(template.New("weatherTemplate").Parse(tpl))
	weatherData := struct {
		Date        string
		WeatherDesc string
		Summary     string
		MaxTemp     int
		MinTemp     int
		SunriseDate string
		SunsetDate  string
	}{
		Date:        dateStr,
		WeatherDesc: weatherDescStr,
		Summary:     d.Summary,
		MaxTemp:     numMax,
		MinTemp:     numMin,
		SunriseDate: sunriseDateStr,
		SunsetDate:  sunsetDateStr,
	}

	var buf bytes.Buffer
	if err := weatherMsgTemplate.Execute(&buf, weatherData); err != nil {
		log.Error().
			Str("service", "weatherMsgTemplate.Execute").
			Msg("Could not parse template for weather message")

		return "Something went wrong while formatting! Try again...", err
	}

	log.Trace().Str("service", "FormatToTextMsg").Msg(buf.String())
	return buf.String(), nil
}
