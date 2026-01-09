package action

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ValianceTekProject/AreaBack/db"
	"github.com/ValianceTekProject/AreaBack/initializers"
	"github.com/google/go-github/v79/github"
)

func GetGithubWebHook(config map[string]any) error {
	ctx := context.Background()
	actionID, ok := config["action_id"].(string)
	if !ok {
		return fmt.Errorf("Unable to retrieve actionId")
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
	var githubToken string
	for _, ust := range user.ServiceTokens() {
		fmt.Println(ust.ServiceID)
		if ust.ServiceID == service.ID {
			githubToken = ust.AccessToken
			break
		}
	}
	if githubToken != "" {
		getGithubPrWebHook(githubToken, actionID, ctx)
	}
	return nil
}

func getGithubPrWebHook(token string, actionID string, ctx context.Context) {
	client := github.NewClient(nil).WithAuthToken(token)

	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		fmt.Printf("Token authentication failed: %v\n", err)
		return
	}

	since := time.Now().Add(-5 * time.Minute).Format(time.RFC3339)
	query := fmt.Sprintf("is:pr author:%s is:closed closed:>%s", user.GetLogin(), since)

	searchOpts := &github.SearchOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	results, resp, err := client.Search.Issues(ctx, query, searchOpts)
	if err != nil {
		fmt.Printf("Failed to search PRs: %v\n", err)
		if resp != nil {
			fmt.Printf("Response Status: %s\n", resp.Status)
		}
		return
	}

	if results.GetTotal() > 0 {
		_, err := initializers.DB.Actions.FindUnique(
			db.Actions.ID.Equals(actionID),
		).Update(
			db.Actions.Triggered.Set(true),
		).Exec(ctx)
		if err != nil {
			log.Printf("Error resetting action trigger: %v", err)
		}

	}
}
