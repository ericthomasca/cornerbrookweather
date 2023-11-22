package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ericthomasca/cornerbrookweather/weather"
	// "github.com/joho/godotenv"
	"github.com/mattn/go-mastodon"
)

func main() {
	// Load .env -- OMIT for image build
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Pull environmental variables
	city := os.Getenv("CITY")
	state := os.Getenv("STATE")
	countryCode := os.Getenv("COUNRTY_CODE")
	apiKey := os.Getenv("OPENWEATHERMAP_API_KEY")
	mastodonServer := os.Getenv("MASTODON_SERVER")
	mastodonClientKey := os.Getenv("MASTODON_CLIENT_KEY")
	mastodonClientSecret := os.Getenv("MASTODON_CLIENT_SECRET")
	mastodonAccessToken := os.Getenv("MASTODON_ACCESS_TOKEN")

	for {
		weatherData, err := weather.Data(city, state, countryCode, apiKey)
		if err != nil {
			log.Fatal(err)
		}
	
		status := buildWeatherStatus(weatherData)
	
		err = postToMastodon(status, mastodonServer, mastodonClientKey, mastodonClientSecret, mastodonAccessToken)
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(time.Hour)
	}
}

func buildWeatherStatus(weatherData weather.Response) string {
	temperature := weather.Temperature(weatherData)
	feelsLike := weather.FeelsLike(weatherData)
	description := weather.Description(weatherData)
	humidity := weather.Humidity(weatherData)
	pressure := weather.Pressure(weatherData)
	wind := weather.WindSummary(weatherData)
	gusts := weather.GustSpeed(weatherData)
	updatedDateTime := weather.UpdatedDateTime(weatherData)

	statusBuilder := strings.Builder{}
	statusBuilder.WriteString(description)
	statusBuilder.WriteString(" at ")
	statusBuilder.WriteString(temperature)
	statusBuilder.WriteString(", feels like ")
	statusBuilder.WriteString(feelsLike)
	statusBuilder.WriteString(". Humidity is ")
	statusBuilder.WriteString(humidity)
	statusBuilder.WriteString(". Pressure is ")
	statusBuilder.WriteString(pressure)
	statusBuilder.WriteString(". Winds of ")
	statusBuilder.WriteString(wind)
	statusBuilder.WriteString(" with gusts of ")
	statusBuilder.WriteString(gusts)
	statusBuilder.WriteString(". Information updated at ")
	statusBuilder.WriteString(updatedDateTime)
	statusBuilder.WriteString(". #CornerBrook #Newfoundland #WeatherUpdate")

	return statusBuilder.String()
}

func postToMastodon(status, mastodonServer, mastodonClientKey, mastodonClientSecret, mastodonAccessToken string) error {
	client := mastodon.NewClient(&mastodon.Config{
		Server:       mastodonServer,
		ClientID:     mastodonClientKey,
		ClientSecret: mastodonClientSecret,
		AccessToken:  mastodonAccessToken,
	})

	_, err := client.PostStatus(context.Background(), &mastodon.Toot{
		Status: status,
	})
	if err != nil {
		return err
	}

	log.Printf("Posted weather update: %s", status)
	return nil
}
