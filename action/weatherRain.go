package action

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ValianceTekProject/AreaBack/db"
	"github.com/ValianceTekProject/AreaBack/initializers"
)

type WeatherResponse struct {
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
	Name string `json:"name"`
}

func ExecWeatherRain(config map[string]any) error {
	ctx := context.Background()

	actionID, ok := config["action_id"].(string)
	if !ok {
		return fmt.Errorf("Unable to retrieve actionId")
	}

	city := "Bordeaux"
	if cityConfig, ok := config["city"].(string); ok && cityConfig != "" {
		city = cityConfig
	}

	execWeatherRainAction(city, actionID, ctx)
	return nil
}

func execWeatherRainAction(city string, actionID string, ctx context.Context) {
	apiKey := os.Getenv("OPENWEATHER_API_KEY")

	if apiKey == "" {
		log.Printf("OPENWEATHER_API_KEY not set in environment")
		return
	}

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s", city, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to call OpenWeatherMap API: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("OpenWeatherMap API returned status %d", resp.StatusCode)
		return
	}

	var weatherResp WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResp); err != nil {
		log.Printf("Failed to decode weather response: %v", err)
		return
	}

	isRaining := false
	for _, weather := range weatherResp.Weather {
		if weather.ID >= 200 && weather.ID < 700 {
			isRaining = true
			break
		}
	}

	if isRaining {
		_, err := initializers.DB.Actions.FindUnique(
			db.Actions.ID.Equals(actionID),
		).Update(
			db.Actions.Triggered.Set(true),
		).Exec(ctx)
		if err != nil {
			log.Printf("Error setting action trigger: %v", err)
		}
	} else {
		log.Printf("It's not raining in %s", weatherResp.Name)
	}
}
