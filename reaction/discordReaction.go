package reaction

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

type DiscordMessageResponse struct {
	ID        string `json:"id"`
	ChannelID string `json:"channel_id"`
	Reactions []struct {
		Emoji struct {
			Name string `json:"name"`
		} `json:"emoji"`
		Me    bool `json:"me"`
		Count int  `json:"count"`
	} `json:"reactions"`
}

func ReactWithDiscordReaction(conf map[string]any) error {
	botToken := os.Getenv("DISCORD_BOT_TOKEN")
	if botToken == "" {
		return fmt.Errorf("DISCORD_BOT_TOKEN environment variable is not set")
	}

	channelID := "1442506377674096690"

	emoji := "🦐"

	latestMessageID, err := getLatestMessage(botToken, channelID)
	if err != nil {
		return fmt.Errorf("error getting latest message: %w", err)
	}

	if latestMessageID == "" {
		return nil
	}

	err = addReaction(botToken, channelID, latestMessageID, emoji)
	if err != nil {
		return fmt.Errorf("error adding reaction: %w", err)
	}

	return nil
}

func getLatestMessage(botToken, channelID string) (string, error) {
	url := fmt.Sprintf("https://discord.com/api/v10/channels/%s/messages?limit=1", channelID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bot %s", botToken))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Discord API error (status %d): %s", resp.StatusCode, string(body))
	}

	var messages []DiscordMessageResponse
	if err := json.NewDecoder(resp.Body).Decode(&messages); err != nil {
		return "", err
	}

	if len(messages) == 0 {
		return "", nil
	}

	return messages[0].ID, nil
}

func addReaction(botToken, channelID, messageID, emoji string) error {
	encodedEmoji := url.PathEscape(emoji)
	apiURL := fmt.Sprintf("https://discord.com/api/v10/channels/%s/messages/%s/reactions/%s/@me",
		channelID, messageID, encodedEmoji)

	req, err := http.NewRequest("PUT", apiURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bot %s", botToken))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Discord API error (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}
