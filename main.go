package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/ericthomasca/cornerbrookweather/weather"
	"github.com/joho/godotenv"
	"github.com/mattn/go-mastodon"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
	country_code := os.Getenv("COUNRTY_CODE")
	api_key := os.Getenv("OPENWEATHERMAP_API_KEY")

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

	// Get current weather
	current_weather := weather.GetWeather(city, state, country_code, api_key)

	// Get current temp
	// Convert to celcius
	current_temperature_value := math.Round(current_weather.Main.Temp - 273.15)
	// Convert float to string
	current_temperature_value_string := strconv.FormatFloat(current_temperature_value, 'f', -1, 64)
	// Add degree celcius notation
	current_temperature := current_temperature_value_string + "°C"
	fmt.Println(current_temperature)

	// Get description
	caser := cases.Title(language.English)
	current_description := caser.String(current_weather.Weather[0].Description)
	fmt.Println(current_description)

	// status := "Test"

	// postStatus(client, status)
}

// func postStatus(client *mastodon.Client, status string) {
// 	newStatus, err := client.PostStatus(context.Background(), &mastodon.Toot{
// 		Status: status,
// 	})
// 	if err != nil {
// 		log.Println("Error posting status:", err)
// 		return
// 	}
// 	fmt.Println("Posted status with ID:", newStatus.ID)
// }
