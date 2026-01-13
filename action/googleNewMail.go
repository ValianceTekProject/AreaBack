package action

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ValianceTekProject/AreaBack/db"
	"github.com/ValianceTekProject/AreaBack/initializers"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func ExecGmailNewEmail(config map[string]any) error {
	ctx := context.Background()

	actionID, googleToken, err := GetGoogleToken(ctx, config)

	if err != nil {
		return err
	}

	if googleToken != "" {
		execGmailNewEmailAction(googleToken, actionID, ctx)
	}
	return nil
}

func execGmailNewEmailAction(token string, actionID string, ctx context.Context) {
	httpClient := &http.Client{
		Transport: &tokenTransport{token: token},
	}

	gmailService, err := gmail.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		fmt.Printf("Failed to create Gmail service: %v\n", err)
		return
	}

	since := time.Now().Add(-1 * time.Minute).Unix()
	query := fmt.Sprintf("after:%d", since)

	messageList, err := gmailService.Users.Messages.List("me").
		Q(query).
		MaxResults(100).
		Do()

	if err != nil {
		fmt.Printf("Failed to list messages: %v\n", err)
		return
	}

	if len(messageList.Messages) > 0 {
		_, err := initializers.DB.Actions.FindUnique(
			db.Actions.ID.Equals(actionID),
		).Update(
			db.Actions.Triggered.Set(true),
		).Exec(ctx)
		if err != nil {
			log.Printf("Error updating action trigger: %v", err)
		}
	}
}
