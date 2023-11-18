package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mattn/go-mastodon"
)

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to Mastodon
	client := mastodon.NewClient(&mastodon.Config{
		Server:      	os.Getenv("MASTODON_SERVER"),
		ClientID:    	os.Getenv("MASTODON_CLIENT_KEY"),
		ClientSecret:	os.Getenv("MASTODON_CLIENT_SECRET"),
		AccessToken:	os.Getenv("MASTODON_ACCESS_TOKEN"),
	})
	if client == nil {
		log.Fatal("Problem connecting to mastodon")
	}

	// TODO: get weather data
	status := "Test"

	postStatus(client, status)
}

func postStatus(client *mastodon.Client, status string) {
	newStatus, err := client.PostStatus(context.Background(), &mastodon.Toot{
		Status: status,
	})
	if err != nil {
		log.Println("Error posting status:", err)
		return
	}
	fmt.Println("Posted status with ID:", newStatus.ID)
}
