package reaction

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/v79/github"
)

func CreateGihubTeam(config map[string]any) error {
	ctx := context.Background()

	reactionID, githubToken, err := GetGithubToken(ctx, config)
	if err != nil {
		return fmt.Errorf("error getting Github token: %w", err)
	}

	if githubToken == "" {
		return fmt.Errorf("empty Github token")
	}

	if err != nil {
		log.Printf("Failed to get reactions: %v", err)
	}
	return execCreateGithubTeam(githubToken, reactionID, config, ctx)
}

func execCreateGithubTeam(token string, reactionID string, config map[string]any, ctx context.Context) error {
	client := github.NewClient(nil).WithAuthToken(token)

	org := "ValianceTekProject"
	newRepo := &github.NewTeam{
		Name: "New Area Team",
	}

	_, resp, err := client.Teams.CreateTeam(ctx, org, *newRepo)
	if err != nil {
		fmt.Printf("Failed to create team: %v\n", err)
		if resp != nil {
			fmt.Printf("Response Status: %s\n", resp.Status)
		}
		return err
	}

	log.Printf("Team created")
	return nil
}

