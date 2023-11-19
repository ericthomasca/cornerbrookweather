# CornerBrook Weather App

## Overview

The CornerBrook Weather App is a Go application designed to retrieve current weather data from the OpenWeatherMap API and post it to the Mastodon account [cornerbrookweather@botsin.space](https://botsin.space/@cornerbrookweather) every hour. This application serves as a convenient way to keep the Corner Brook community updated on their local weather.

## Features

- OpenWeatherMap Integration: Utilizes the OpenWeatherMap API to fetch current weather data for Corner Brook.
- Mastodon Posting: Posts the fetched weather information to the Mastodon account [cornerbrookweather@botsin.space](https://botsin.space/@cornerbrookweather) on an hourly basis.
- Hourly Execution: The application is designed to run hourly through a scheduler or cron job.

## Configuration

- Rename .env.example to .env or make a copy.
- Edit the .env with values provided from Mastodon and OpenWeatherMap.
