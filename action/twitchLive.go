package action

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ValianceTekProject/AreaBack/db"
	"github.com/ValianceTekProject/AreaBack/initializers"
	"github.com/nicklaw5/helix/v2"
)

func ExecTwitchLive(config map[string]any) error {
	ctx := context.Background()

	actionID, ok := config["action_id"].(string)
	if !ok {
		return fmt.Errorf("Unable to retrieve actionId")
	}

	streamerName := "Gotaga"
	execTwitchLiveAction(streamerName, actionID, ctx)
	return nil
}

func execTwitchLiveAction(streamerName string, actionID string, ctx context.Context) {
	clientID := os.Getenv("TWITCH_CLIENT_ID")
	clientSecret := os.Getenv("TWITCH_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		log.Printf("TWITCH_CLIENT_ID or TWITCH_CLIENT_SECRET not set in environment")
		return
	}

	client, err := helix.NewClient(&helix.Options{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	})
	if err != nil {
		log.Printf("Failed to create Twitch client: %v", err)
		return
	}

	resp, err := client.RequestAppAccessToken([]string{})
	if err != nil {
		log.Printf("Failed to get Twitch app access token: %v", err)
		return
	}

	client.SetAppAccessToken(resp.Data.AccessToken)

	streamsResp, err := client.GetStreams(&helix.StreamsParams{
		UserLogins: []string{streamerName},
	})
	if err != nil {
		log.Printf("Failed to get Twitch streams: %v", err)
		return
	}

	if len(streamsResp.Data.Streams) > 0 {
		stream := streamsResp.Data.Streams[0]
		log.Printf("Streamer %s is live! Title: %s, Game: %s, Viewers: %d",
			stream.UserName, stream.Title, stream.GameName, stream.ViewerCount)

		_, err := initializers.DB.Actions.FindUnique(
			db.Actions.ID.Equals(actionID),
		).Update(
			db.Actions.Triggered.Set(true),
		).Exec(ctx)
		if err != nil {
			log.Printf("Error setting action trigger: %v", err)
		}
	} else {
		log.Printf("Streamer %s is offline", streamerName)
	}
}
