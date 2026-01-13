package action

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/ValianceTekProject/AreaBack/db"
	"github.com/ValianceTekProject/AreaBack/initializers"
)

func GetGoogleToken(ctx context.Context, config map[string]any) (string, string, error) {
	actionID, ok := config["action_id"].(string)
	if !ok {
		return "", "", fmt.Errorf("Unable to retrieve actionId")
	}
	action, err := initializers.DB.Actions.FindUnique(
		db.Actions.ID.Equals(actionID),
	).With(
		db.Actions.Area.Fetch().With(
			db.Areas.User.Fetch().With(
				db.Users.ServiceTokens.Fetch(),
			),
		),
		db.Actions.Service.Fetch(),
	).Exec(ctx)
	if err != nil {
		log.Printf("Failed to get Actions: %v", err)
	}
	area := action.Area()
	user := area.User()
	service := action.Service()
	var googleToken string
	for _, ust := range user.ServiceTokens() {
		if ust.ServiceID == service.ID {
			googleToken = ust.AccessToken
			break
		}
	}
	return actionID, googleToken, nil
}

type tokenTransport struct {
	token string
}

func (t *tokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+t.token)
	return http.DefaultTransport.RoundTrip(req)
}
