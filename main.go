package main

import (
	"discord-weather-bot/bot"
	"log"
	"os" /* this will load the env variables */
)

func main() {
	// Load environment variables
	botToken, ok := os.LookupEnv("BOT_TOKEN")
	if !ok {
		log.Fatal("Must set Discord token as env variable: BOT_TOKEN")
	}
	openWeatherToken, ok := os.LookupEnv("OPENWEATHER_TOKEN")
	if !ok {
		log.Fatal("Must set Open Weather token as env variable: OPENWEATHER_TOKEN")
	}

	// Start the bot
	bot.BotToken = botToken
	bot.OpenWeatherToken = openWeatherToken
	bot.Run()
}
