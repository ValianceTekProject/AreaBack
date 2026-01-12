package reaction

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/ValianceTekProject/AreaBack/db"
	"github.com/ValianceTekProject/AreaBack/initializers"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func SendEmail(config map[string]any) error {
	ctx := context.Background()

	reactionID, googleToken, err := GetGoogleToken(ctx, config)
	if err != nil {
		return fmt.Errorf("error getting Google token: %w", err)
	}

	if googleToken == "" {
		return fmt.Errorf("empty Google token")
	}

	reaction, err := initializers.DB.Reactions.FindUnique(
		db.Reactions.ID.Equals(reactionID),
	).With(
		db.Reactions.Area.Fetch().With(
			db.Areas.User.Fetch(),
		),
	).Exec(ctx)
	if err != nil {
		log.Printf("Failed to get reactions: %v", err)
	}
	area := reaction.Area()
	user := area.User()
	return execSendEmailAction(googleToken, reactionID, user, ctx)
}

func execSendEmailAction(token string, reactionID string, user *db.UsersModel, ctx context.Context) error {
	httpClient := &http.Client{
		Transport: &tokenTransport{token: token},
	}

	gmailService, err := gmail.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		fmt.Printf("Failed to create Gmail service: %v\n", err)
		return err
	}

	to := user.Email
	subject := "Action Detected"
	body := "AREA"

	emailContent := fmt.Sprintf(
		"To: %s\r\n"+
			"Subject: %s\r\n"+
			"Content-Type: text/plain; charset=utf-8\r\n\r\n"+
			"%s",
		to, subject, body,
	)

	encodedMessage := base64.URLEncoding.EncodeToString([]byte(emailContent))

	message := &gmail.Message{
		Raw: encodedMessage,
	}

	_, err = gmailService.Users.Messages.Send("me", message).Do()
	if err != nil {
		fmt.Printf("Failed to send email: %v\n", err)
		return err
	}

	log.Printf("Email sent successfully (reactionID: %s)", reactionID)
	return nil
}
