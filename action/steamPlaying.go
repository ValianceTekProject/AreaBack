package action

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/ValianceTekProject/AreaBack/db"
	"github.com/ValianceTekProject/AreaBack/initializers"
)

type SteamPlayerSummary struct {
	Response struct {
		Players []struct {
			SteamID       string `json:"steamid"`
			PersonaName   string `json:"personaname"`
			GameID        string `json:"gameid"`
			GameExtraInfo string `json:"gameextrainfo"`
		} `json:"players"`
	} `json:"response"`
}

func ExecSteamPlaying(config map[string]any) error {
	ctx := context.Background()

	actionID, ok := config["action_id"].(string)
	if !ok {
		return fmt.Errorf("Unable to retrieve actionId")
	}

	steamID := "76561198375694417"
	if !ok || steamID == "" {
		return fmt.Errorf("Unable to retrieve steam_id from config")
	}

	execSteamPlayingAction(steamID, actionID, ctx)
	return nil
}

func execSteamPlayingAction(steamID string, actionID string, ctx context.Context) {
	apiKey := os.Getenv("STEAM_API_KEY")
	if apiKey == "" {
		log.Printf("STEAM_API_KEY not set in environment")
		return
	}

	url := fmt.Sprintf("http://api.steampowered.com/ISteamUser/GetPlayerSummaries/v2/?key=%s&steamids=%s", apiKey, steamID)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to call Steam API: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Steam API returned status: %s", resp.Status)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read Steam API response: %v", err)
		return
	}

	var playerSummary SteamPlayerSummary
	err = json.Unmarshal(body, &playerSummary)
	if err != nil {
		log.Printf("Failed to parse Steam API response: %v", err)
		return
	}

	if len(playerSummary.Response.Players) == 0 {
		log.Printf("No player found with Steam ID: %s", steamID)
		return
	}

	player := playerSummary.Response.Players[0]

	if player.GameID != "" {
		log.Printf("Player %s is playing: %s (GameID: %s)", player.PersonaName, player.GameExtraInfo, player.GameID)
		_, err := initializers.DB.Actions.FindUnique(
			db.Actions.ID.Equals(actionID),
		).Update(
			db.Actions.Triggered.Set(true),
		).Exec(ctx)
		if err != nil {
			log.Printf("Error setting action trigger: %v", err)
		}
	}
}
