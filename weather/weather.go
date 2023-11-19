package weather

import (
	"encoding/json"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Weather struct {
	Description string `json:"description"` // Overcast Clouds
}

type Conditions struct {
	Temp      float64 `json:"temp"`       // in KELVIN
	FeelsLike float64 `json:"feels_like"` // in KELVIN
	Humidity  int     `json:"humidity"`
	Pressure  int     `json:"pressure"`
}

type Wind struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
	Gust  float64 `json:"gust"`
}

type Sys struct {
	Country string `json:"country"`
}

type WeatherResponse struct {
	Weather    []Weather  `json:"weather"`
	Conditions Conditions `json:"main"`
	Wind       Wind       `json:"wind"`
	Sys        Sys        `json:"sys"`
	Name       string     `json:"name"`
	Dt         int        `json:"dt"`
}

const (
	AbsoluteZero float64 = -273.15
)

// GetWeatherData returns full weather data as WeatherResponse.
func GetWeatherData(city string, state string, country_code string, api_key string) WeatherResponse {
	// Handle convert spaces and symbols to URL escape versions
	city = url.PathEscape(city)
	state = url.PathEscape(state)

	// Build query
	query := city + "," + state + "," + country_code

	// Build string of url
	builder := strings.Builder{}
	builder.WriteString("https://api.openweathermap.org/data/2.5/weather?q=")
	builder.WriteString(query)
	builder.WriteString("&appid=")
	builder.WriteString(api_key)
	url := builder.String()

	// Get json response from url and return
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var response WeatherResponse
	json.Unmarshal(body, &response)

	return response
}

// GetTemperature returns the temperature in Celsius with "℃" notation as string.
func GetTemperature(weatherData WeatherResponse) string {
	temperatureFloat := math.Round(weatherData.Conditions.Temp + AbsoluteZero)
	temperature := strconv.FormatFloat(temperatureFloat, 'f', -1, 64)
	return temperature + "˚C"
}

// GetFeelsLike returns the 'feels like' temperature in Celsius with "℃" notation as string.
func GetFeelsLike(weatherData WeatherResponse) string {
	feelsLikeFloat := math.Round(weatherData.Conditions.FeelsLike + AbsoluteZero)
	feelsLike := strconv.FormatFloat(feelsLikeFloat, 'f', -1, 64)
	return feelsLike + "˚C"
}

// GetDescription returns the description of the weather in title case as string
func GetDescription(weatherData WeatherResponse) string {
	caser := cases.Title(language.English)
	return caser.String(weatherData.Weather[0].Description)
}

// GetHumidity returns humidity as percentage with "%" notation as string
func GetHumidity(weatherData WeatherResponse) string {
	humidity := weatherData.Conditions.Humidity
	return strconv.Itoa(humidity) + "%"
}

// GetPressure returns pressure as percentage with "%" notation as string
func GetPressure(weatherData WeatherResponse) string {
	pressure := weatherData.Conditions.Pressure
	return strconv.Itoa(pressure) + " hPa"
}

// GetWindSpeed returns wind speed in km/h with "km/h" notation as string
func GetWindSpeed(weatherData WeatherResponse) string {
	windSpeedFloat := weatherData.Wind.Speed * 3.6
	windSpeed := strconv.FormatFloat(windSpeedFloat, 'f', -1, 64)
	return windSpeed + " km/h"
}

// GetCardinalDirection returns string of cardinal direction given direction in degrees
func GetCardinalDirection(degrees int) string {
	// Normalize degrees to be within 0 to 360 range
	degrees = degrees % 360

	// Define cardinal directions with their ranges
	cardinalDirections := map[string][2]int{
		"N":   {0, 22},
		"NNE": {23, 67},
		"NE":  {68, 112},
		"ENE": {113, 157},
		"E":   {158, 202},
		"ESE": {203, 247},
		"SE":  {248, 292},
		"SSE": {293, 337},
		"S":   {338, 360},
	}

	// Determine the cardinal direction
	for dir, rangeDeg := range cardinalDirections {
		if degrees >= rangeDeg[0] && degrees <= rangeDeg[1] {
			return dir
		}
	}

	return "???"
}

// GetWind returns string of wind's speed and direction
func GetWind(weatherData WeatherResponse) string {
	windSpeedValue := math.Round(weatherData.Wind.Speed * 3.6)
	windSpeed := strconv.FormatFloat(windSpeedValue, 'f', -1, 64)
	windSpeed = windSpeed + " km/h"

	windDegreeValue := weatherData.Wind.Deg
	windDirection := GetCardinalDirection(windDegreeValue)

	return windSpeed + " " + windDirection
}

// GetGustSpeed returns wind speed in km/h with "km/h" notation as string
func GetGustSpeed(weatherData WeatherResponse) string {
	windSpeedValue := math.Round(weatherData.Wind.Gust * 3.6)
	windSpeed := strconv.FormatFloat(windSpeedValue, 'f', -1, 64)
	return windSpeed + " km/h"
}

// GetLocation returns comma seperated string with  city name and country code
func GetLocation(weatherData WeatherResponse) string {
	city := weatherData.Name
	countryCode := weatherData.Sys.Country
	return city + ", " + countryCode
}

func EpochToFormattedString(epochTime int64) string {
	// Convert epoch time to time.Time
	t := time.Unix(epochTime, 0)

	// Format the time as required
	formattedTime := t.Format("2006-01-02 3:04PM MST")

	return formattedTime
}

// GetUpdatedDateTime returns formatted string of when weather was last updated
func GetUpdatedDateTime(weatherData WeatherResponse) string {
	dateTimeValue := time.Unix(int64(weatherData.Dt), 0)
	return dateTimeValue.Format("2006-01-02 3:04PM NST")

}
