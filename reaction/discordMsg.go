package reaction

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ValianceTekProject/AreaBack/model"
)

func ReactWithDiscordMsg(conf map[string]any) error {
	var config model.DiscordConfig

	config.WebhookURL = "https://discord.com/api/webhooks/1447899360145965127/NNlPLGO08_jbI1Pz0SED7wFqtsYvK4H93xG_v6Z25VZEuDd090IkvfpicUOFak9U5bPz"

	if config.WebhookURL == "" {
		return fmt.Errorf("Error in webhook url")
	}

	payload := model.DiscordMessage{
		Content: "PR Merge",
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Error unmarshaling message: %s\n", err.Error())
	}

	resp, err := http.Post(config.WebhookURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("Error: %s\n", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Printf("Error: %s\n", resp.Status)
	}
	return nil
}
