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

func GetGithubWebHook(userID string) {
	ctx := context.Background()
	userServiceToken, err := initializers.DB.UserServiceTokens.FindMany(
		db.UserServiceTokens.UserID.Equals(userID),
	).Exec(ctx)
	if err != nil {
		fmt.Printf("Error fetching tokens for user %s: %v\n", userID, err)
		return
	}
	githubServiceNumber, err := initializers.DB.Services.FindUnique(
		db.Services.Name.Equals("Github"),
	).Exec(ctx)
	if err != nil {
		fmt.Printf("Error fetching Github service: %v\n", err)
		return
	}
	var githubToken string
	for _, ust := range userServiceToken {
		fmt.Println(ust.ServiceID)
		if ust.ServiceID == githubServiceNumber.ID {
			githubToken = ust.AccessToken
			break
		}
	}
	if githubToken != "" {
		getGithubPrWebHook(githubToken, ctx)
	}
}

func getGithubPrWebHook(token string, ctx context.Context) {
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
		area, err := initializers.DB.Areas.FindUnique(
			db.Areas.Name.Equals("Github_pr_to_discord"),
		).With(
			db.Areas.Actions.Fetch(),
		).Exec(ctx)
		if err != nil {
			log.Printf("Error fetching area: %v", err)
			return
		}

		actions := area.Actions()
		for _, action := range actions {
			_, err = initializers.DB.Actions.FindUnique(
				db.Actions.ID.Equals(action.ID),
			).Update(
				db.Actions.Triggered.Set(true),
			).Exec(ctx)

			if err != nil {
				log.Printf("Error triggering action %s: %v", action.ID, err)
			} else {
				fmt.Printf("Action %s triggered!\n", action.ID)
			}
		}
	}
}
