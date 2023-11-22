package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
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

type Response struct {
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

// Data returns full weather data as Response with error.
func Data(city string, state string, country_code string, api_key string) (Response, error) {
	// Handle convert spaces and symbols to URL escape versions
	city = url.PathEscape(city)
	state = url.PathEscape(state)

	// Build query
	query := city + "," + state + "," + country_code

	// Build URL string
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s", query, api_key)

	// Get json response from url
	resp, err := http.Get(url)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return Response{}, err
	}

	return response, nil
}

// Temperature returns the temperature in Celsius with "℃" notation as string.
func Temperature(weatherData Response) string {
	temperatureFloat := math.Round(weatherData.Conditions.Temp + AbsoluteZero)
	temperature := strconv.FormatFloat(temperatureFloat, 'f', -1, 64)
	return temperature + "˚C"
}

// FeelsLike returns the 'feels like' temperature in Celsius with "℃" notation as string.
func FeelsLike(weatherData Response) string {
	feelsLikeFloat := math.Round(weatherData.Conditions.FeelsLike + AbsoluteZero)
	feelsLike := strconv.FormatFloat(feelsLikeFloat, 'f', -1, 64)
	return feelsLike + "˚C"
}

// Description returns the description of the weather in title case as string
func Description(weatherData Response) string {
	caser := cases.Title(language.English)
	return caser.String(weatherData.Weather[0].Description)
}

// Humidity returns humidity as percentage with "%" notation as string
func Humidity(weatherData Response) string {
	humidity := weatherData.Conditions.Humidity
	return strconv.Itoa(humidity) + "%"
}

// Pressure returns pressure as percentage with "%" notation as string
func Pressure(weatherData Response) string {
	pressure := weatherData.Conditions.Pressure
	return strconv.Itoa(pressure) + " hPa"
}

// WindSpeed returns wind speed in km/h with "km/h" notation as string
func WindSpeed(weatherData Response) string {
	windSpeedFloat := weatherData.Wind.Speed * 3.6
	windSpeed := strconv.FormatFloat(windSpeedFloat, 'f', -1, 64)
	return windSpeed + " km/h"
}

// DegreeToCardinalDirection returns string of cardinal direction given direction in degrees
func DegreeToCardinalDirection(degrees int) string {
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

	return ""
}

// WindSummary returns string of wind's speed and direction
func WindSummary(weatherData Response) string {
	windSpeedValue := math.Round(weatherData.Wind.Speed * 3.6)
	windSpeed := strconv.FormatFloat(windSpeedValue, 'f', -1, 64)
	windSpeed = windSpeed + " km/h"

	windDegreeValue := weatherData.Wind.Deg
	windDirection := DegreeToCardinalDirection(windDegreeValue)

	return windSpeed + " " + windDirection
}

// GustSpeed returns wind speed in km/h with "km/h" notation as string
func GustSpeed(weatherData Response) string {
	windSpeedValue := math.Round(weatherData.Wind.Gust * 3.6)
	windSpeed := strconv.FormatFloat(windSpeedValue, 'f', -1, 64)
	return windSpeed + " km/h"
}

// Location returns comma seperated string with  city name and country code
func Location(weatherData Response) string {
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

// UpdatedDateTime returns formatted string of when weather was last updated
func UpdatedDateTime(weatherData Response) string {
	dateTimeValue := time.Unix(int64(weatherData.Dt), 0)
	return dateTimeValue.Format("2006-01-02 3:04PM NST")
}
