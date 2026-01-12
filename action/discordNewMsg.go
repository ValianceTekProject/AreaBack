package action

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ValianceTekProject/AreaBack/db"
	"github.com/ValianceTekProject/AreaBack/initializers"
	"io"
	"net/http"
	"os"
	"time"
)

type DiscordMessage struct {
	ID        string    `json:"id"`
	ChannelID string    `json:"channel_id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	Author    struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	} `json:"author"`
}

func ExecDiscordNewMsg(config map[string]any) error {
	ctx := context.Background()

	actionID, ok := config["action_id"].(string)
	if !ok {
		return fmt.Errorf("unable to retrieve action_id")
	}

	botToken := os.Getenv("DISCORD_BOT_TOKEN")
	if botToken == "" {
		return fmt.Errorf("DISCORD_BOT_TOKEN environment variable is not set")
	}

	channelID := "1442506377674096690"

	execDiscordNewMsgAction(botToken, channelID, actionID, ctx)

	return nil
}

func execDiscordNewMsgAction(botToken string, channelID string, actionID string, ctx context.Context) {
	url := fmt.Sprintf("https://discord.com/api/v10/channels/%s/messages?limit=10", channelID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bot %s", botToken))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error calling Discord API:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Discord API error (status %d): %s\n", resp.StatusCode, string(body))
		return
	}

	var messages []DiscordMessage
	if err := json.NewDecoder(resp.Body).Decode(&messages); err != nil {
		fmt.Println("Error decoding Discord response:", err)
		return
	}

	now := time.Now()
	hasNewMessage := false

	for _, msg := range messages {
		timeSinceMessage := now.Sub(msg.Timestamp)
		if timeSinceMessage <= 1*time.Minute {
			hasNewMessage = true
			fmt.Printf("New message detected in channel %s: %s (by %s)\n",
				channelID, msg.Content, msg.Author.Username)
			break
		}
	}

	if hasNewMessage {
		_, err := initializers.DB.Actions.FindUnique(
			db.Actions.ID.Equals(actionID),
		).Update(
			db.Actions.Triggered.Set(true),
		).Exec(ctx)

		if err != nil {
			fmt.Println("Error updating action trigger:", err)
		}
	}
}
