package weather

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Main struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	Pressure  int     `json:"pressure"`
	Humidity  int     `json:"humidity"`
}

type Wind struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
	Gust  float64 `json:"gust"`
}

type WeatherResponse struct {
	Weather []Weather `json:"weather"`
	Main    Main      `json:"main"`
	Wind    Wind      `json:"wind"`
}

func GetWeather(city string, state string, country_code string, api_key string) WeatherResponse {
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

	// Get json response from url
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
