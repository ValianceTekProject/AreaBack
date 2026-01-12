package reaction

import (
	"context"
	"fmt"
	"log"

	"github.com/ValianceTekProject/AreaBack/db"
	"github.com/ValianceTekProject/AreaBack/initializers"
)

func GetGithubToken(ctx context.Context, config map[string]any) (string, string, error) {
	reactionID, ok := config["reaction_id"].(string)
	if !ok {
		return "", "", fmt.Errorf("Unable to retrieve reactionId")
	}
	action, err := initializers.DB.Reactions.FindUnique(
		db.Reactions.ID.Equals(reactionID),
	).With(
		db.Reactions.Area.Fetch().With(
			db.Areas.User.Fetch().With(
				db.Users.ServiceTokens.Fetch(),
			),
		),
		db.Reactions.Service.Fetch(),
	).Exec(ctx)
	if err != nil {
		log.Printf("Failed to get reactions: %v", err)
	}
	area := action.Area()
	user := area.User()
	service := action.Service()
	var githubToken string
	for _, ust := range user.ServiceTokens() {
		if ust.ServiceID == service.ID {
			githubToken = ust.AccessToken
			break
		}
	}
	return reactionID, githubToken, nil
}
