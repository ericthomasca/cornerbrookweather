package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/ericthomasca/cornerbrookweather/weather"
	"github.com/joho/godotenv"
	"github.com/mattn/go-mastodon"
)

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Pull environmental variables
	city := os.Getenv("CITY")
	state := os.Getenv("STATE")
	countryCode := os.Getenv("COUNRTY_CODE")
	apiKey := os.Getenv("OPENWEATHERMAP_API_KEY")

	// Get weather data
	weatherData := weather.GetWeatherData(city, state, countryCode, apiKey)
	temperature := weather.GetTemperature(weatherData)
	feelsLike := weather.GetFeelsLike(weatherData)
	description := weather.GetDescription(weatherData)
	humidity := weather.GetHumidity(weatherData)
	pressure := weather.GetPressure(weatherData)
	wind := weather.GetWind(weatherData)
	gusts := weather.GetGustSpeed(weatherData)
	location := weather.GetLocation(weatherData)
	updatedDateTime := weather.GetUpdatedDateTime(weatherData)

	// Build toot
	builder := strings.Builder{}
	builder.WriteString(location)
	builder.WriteString(" Weather Update: ")
	builder.WriteString(description)
	builder.WriteString(" at ")
	builder.WriteString(temperature)
	builder.WriteString(", feels like ")
	builder.WriteString(feelsLike)
	builder.WriteString(". Humidity is ")
	builder.WriteString(humidity)
	builder.WriteString(". Pressure is ")
	builder.WriteString(pressure)
	builder.WriteString(". Winds of ")
	builder.WriteString(wind)
	builder.WriteString(" with gusts of ")
	builder.WriteString(gusts)
	builder.WriteString(". Information updated at ")
	builder.WriteString(updatedDateTime)
	builder.WriteString(". #CornerBrook #Newfoundland #WeatherUpdate")
	status := builder.String()

	// Connect to Mastodon
	client := mastodon.NewClient(&mastodon.Config{
		Server:       os.Getenv("MASTODON_SERVER"),
		ClientID:     os.Getenv("MASTODON_CLIENT_KEY"),
		ClientSecret: os.Getenv("MASTODON_CLIENT_SECRET"),
		AccessToken:  os.Getenv("MASTODON_ACCESS_TOKEN"),
	})
	if client == nil {
		log.Fatal("Problem connecting to mastodon")
	}
	newStatus, err := client.PostStatus(context.Background(), &mastodon.Toot{
		Status: status,
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("Posted status with ID:", newStatus.ID)
}
