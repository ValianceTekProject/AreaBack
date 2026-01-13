package reaction

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/v79/github"
)

func CreateGihubPr(config map[string]any) error {
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
	return execCreatePr(githubToken, reactionID, config, ctx)
}

func execCreatePr(token string, reactionID string, config map[string]any, ctx context.Context) error {
	client := github.NewClient(nil).WithAuthToken(token)

	owner := "ValianceTekProject"
	repo := "TestArea"
	title, _ := config["title"].(string)
	body, _ := config["body"].(string)
	head := "dev"
	base, _ := config["base"].(string)

	if title == "" {
		title = "Automated PR from AREA"
	}
	if body == "" {
		body = "This pull request was automatically created by AREA."
	}
	if base == "" {
		base = "main"
	}

	if owner == "" || repo == "" || head == "" {
		err := fmt.Errorf("missing required parameters: owner, repo, or head branch")
		fmt.Printf("Failed to create PR: %v\n", err)
		return err
	}

	newPR := &github.NewPullRequest{
		Title: github.Ptr(title),
		Body:  github.Ptr(body),
		Head:  github.Ptr(head),
		Base:  github.Ptr(base),
	}

	pr, resp, err := client.PullRequests.Create(ctx, owner, repo, newPR)
	if err != nil {
		fmt.Printf("Failed to create pull request: %v\n", err)
		if resp != nil {
			fmt.Printf("Response Status: %s\n", resp.Status)
		}
		return err
	}

	log.Printf("Pull request created successfully (reactionID: %s, PR #%d)", reactionID, pr.GetNumber())
	return nil
}
