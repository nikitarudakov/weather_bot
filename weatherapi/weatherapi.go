package weatherapi

import (
	"encoding/json"
	"fmt"
	"git.foxminded.ua/foxstudent106092/weather-bot/config"
	"git.foxminded.ua/foxstudent106092/weather-bot/db"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
)

// DailyTemp stores fields of max and min temperature
type DailyTemp struct {
	Min float64 `json:"Min"`
	Max float64 `json:"max"`
}

// Weather stores field of weather description
type Weather struct {
	Desc string `json:"description,omitempty"`
}

// DailyWeather stores fields of weather features at UnixDt day
type DailyWeather struct {
	UnixDt  int64     `json:"dt"`
	Sunrise int64     `json:"sunrise"`
	Sunset  int64     `json:"sunset"`
	Summary string    `json:"summary"`
	Temp    DailyTemp `json:"temp"`
	Weather []Weather `json:"weather"`
}

// ResponseWeatherAPI represents Weather API response of daily forecasts
type ResponseWeatherAPI struct {
	Daily []DailyWeather `json:"daily"`
}

type UserWeatherForecast struct {
	UserID   int64               `bson:"user_id"`
	Forecast *ResponseWeatherAPI `json:"forecast"`
}

// WeatherAPI represents the methods for interacting with a weather API
type WeatherAPI interface {
	GetWeatherForecast(lat, lon string) (*ResponseWeatherAPI, error)
	ReadWeatherForecastFromDB(dbClient db.DatabaseAccessor, dbCfg *config.DbCfg, userID int64) (*ResponseWeatherAPI, error)
}

// WeatherService represents a service that interacts with a weather API
type WeatherService struct {
	cfg *config.WeatherAPICfg
}

func (resp *ResponseWeatherAPI) StoreWeatherForecastForUser(client db.DatabaseAccessor, d *config.DbCfg, userID int64) error {
	usf := &UserWeatherForecast{
		UserID:   userID,
		Forecast: resp,
	}

	fmt.Printf("%+v\n", usf)

	var forecast UserWeatherForecast
	err := client.FindUserInDB(userID, d.ForecastCollectionName).Decode(&forecast)
	if err == nil {
		return nil
	}

	fmt.Print(err)
	fmt.Printf("%+v\n", forecast)

	if err := client.InsertItemToDB(usf, d.ForecastCollectionName); err != nil {
		return err
	}

	return nil
}

// NewWeatherAPIService initialized new WeatherService instance
func NewWeatherAPIService(cfg *config.WeatherAPICfg) WeatherAPI {
	var weatherApi WeatherAPI = &WeatherService{cfg: cfg}
	return weatherApi
}

func (wa *WeatherService) getAPIUrl(lat, lon string) string {
	return wa.cfg.Server +
		"?lat=" + lat +
		"&lon=" + lon +
		"&appid=" + wa.cfg.Token +
		"&exclude=" + wa.cfg.Exclude +
		"&units=" + wa.cfg.Units
}

// GetWeatherForecast fetches API and returns Response - which is weather forecast
func (wa *WeatherService) GetWeatherForecast(lat, lon string) (*ResponseWeatherAPI, error) {
	// Get API url for location coords lat/lon
	apiURL := wa.getAPIUrl(lat, lon)

	resp, err := http.Get(apiURL)
	if err != nil {
		log.Error().Str("service", "http.Get").Err(err).Send()
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Str("service", "io.ReadAll").Err(err).Send()
		return nil, err
	}

	weatherResp := ResponseWeatherAPI{}
	if err = json.Unmarshal(body, &weatherResp); err != nil {
		log.Error().Str("service", "json.Unmarshal").Err(err).Send()
		return nil, err
	}

	return &weatherResp, nil
}

func (wa *WeatherService) ReadWeatherForecastFromDB(
	dbClient db.DatabaseAccessor,
	dbCfg *config.DbCfg, userID int64,
) (*ResponseWeatherAPI, error) {

	var forecast *UserWeatherForecast
	if err := dbClient.FindUserInDB(userID, dbCfg.ForecastCollectionName).Decode(&forecast); err != nil {
		return nil, err
	}

	return forecast.Forecast, nil
}
