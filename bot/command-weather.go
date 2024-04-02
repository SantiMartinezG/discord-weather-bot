package bot

import (
	"encoding/json"
	"fmt"
	// "io"
	"net/http"
	// "regexp"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const URL string = "https://api.openweathermap.org/data/2.5/weather?"

type WeatherData struct {
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
	Name string `json:"name"`
}

func getCurrentWeather(message string) *discordgo.MessageSend {
	// Extract city name from message
	city := strings.TrimSpace(message)

	// If city name not provided, return an error
	if city == "" {
		return &discordgo.MessageSend{
			Content: "Sorry, you need to provide a city name",
		}
	}

	// Build full URL to query OpenWeather
	weatherURL := fmt.Sprintf("%sq=%s&units=metric&appid=%s", URL, url.QueryEscape(city), OpenWeatherToken)

	// Create new HTTP client & set timeout
	client := http.Client{Timeout: 5 * time.Second}

	// Query OpenWeather API
	response, err := client.Get(weatherURL)
	if err != nil {
		return &discordgo.MessageSend{
			Content: "Sorry, there was an error trying to get the weather",
		}
	}
	defer response.Body.Close()

	// Check if response is successful
	if response.StatusCode != http.StatusOK {
		return &discordgo.MessageSend{
			Content: fmt.Sprintf("Sorry, couldn't find weather for the city %s", city),
		}
	}

	// Convert JSON
	var data WeatherData
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return &discordgo.MessageSend{
			Content: "Sorry, there was an error processing the server's response",
		}
	}

	// Pull out desired weather info & Convert to string if necessary
	conditions := data.Weather[0].Description
	temperature := strconv.FormatFloat(data.Main.Temp, 'f', 2, 64)
	humidity := strconv.Itoa(data.Main.Humidity)
	wind := strconv.FormatFloat(data.Wind.Speed, 'f', 2, 64)

	// Build Discord embed response
	embed := &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Type:        discordgo.EmbedTypeRich,
				Title:       "Current Weather",
				Description: "Weather for " + city,
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Conditions",
						Value:  conditions,
						Inline: true,
					},
					{
						Name:   "Temperature",
						Value:  temperature + "°C",
						Inline: true,
					},
					{
						Name:   "Humidity",
						Value:  humidity + "%",
						Inline: true,
					},
					{
						Name:   "Wind",
						Value:  wind + " km/h",
						Inline: true,
					},
				},
			},
		},
	}

	return embed
}
