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

func ExecGithubNewPr(config map[string]any) error {
	ctx := context.Background()

	actionID, githubToken, error := GetGithubToken(ctx, config)

	if error == nil {
		return nil
	}

	if githubToken != "" {
		execGithubNewPrAction(githubToken, actionID, ctx)
	}
	return nil
}

func execGithubNewPrAction(token string, actionID string, ctx context.Context) {
	client := github.NewClient(nil).WithAuthToken(token)

	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		fmt.Printf("Token authentication failed: %v\n", err)
		return
	}

	since := time.Now().Add(-5 * time.Minute).Format(time.RFC3339)
	query := fmt.Sprintf("is:pr author:%s created:>%s", user.GetLogin(), since)

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



