package reaction

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ValianceTekProject/AreaBack/db"
	"github.com/ValianceTekProject/AreaBack/model"
)

func ReactWithDiscordMsg(user db.UsersModel, reaction db.ReactionsModel) {
	// configJSON, ok := reaction.Config()
	// if !ok {
	// 	fmt.Println("Error: No config found")
	// 	return
	// }
	//
	var config model.DiscordConfig
	// if err := json.Unmarshal(configJSON, &config); err != nil {
	// 	fmt.Printf("Error unmarshaling config %s\n", err.Error())
	// 	return
	// }

	config.WebhookURL = "https://discord.com/api/webhooks/1447899360145965127/NNlPLGO08_jbI1Pz0SED7wFqtsYvK4H93xG_v6Z25VZEuDd090IkvfpicUOFak9U5bPz"

	if config.WebhookURL == "" {
		fmt.Println("Error in webhook url")
		return
	}

	payload := model.DiscordMessage{
		Content: "PR Merge",
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error unmarshaling message: %s\n", err.Error())
		return
	}

	resp, err := http.Post(config.WebhookURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Printf("Error: %s\n", resp.Status)
	}
}
