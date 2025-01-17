package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Store Bot API Tokens:
var (
	OpenWeatherToken string
	BotToken         string
)

func Run() {
	// Create new Discord Session
	discord, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatal(err)
	}

	// Add event handler for general messages
	discord.AddHandler(newMessage)

	// Open session
	discord.Open()
	defer discord.Close()

	// Run until code is terminated
	fmt.Println("Bot running...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {
	// Ignore bot message
	if message.Author.ID == discord.State.User.ID {
		return
	}

	// Respond to messages
	switch {
	case strings.Contains(message.Content, "weather"):
		discord.ChannelMessageSend(message.ChannelID, "I can help with that! Use '!city <city name>'")
	case strings.Contains(message.Content, "tempy"):
		discord.ChannelMessageSend(message.ChannelID, "I am the weather god")
	case strings.HasPrefix(message.Content, "!city"):
		city := strings.TrimSpace(strings.TrimPrefix(message.Content, "!city"))
		currentWeather := getCurrentWeather(city)
		discord.ChannelMessageSendComplex(message.ChannelID, currentWeather)
	}
}
