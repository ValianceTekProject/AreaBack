package action

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ValianceTekProject/AreaBack/db"
	"github.com/ValianceTekProject/AreaBack/initializers"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func ExecGoogleDriveNewFile(config map[string]any) error {
	ctx := context.Background()

	actionID, googleToken, err := GetGoogleToken(ctx, config)

	if err != nil {
		return err
	}

	if googleToken != "" {
		execGoogleDriveNewFileAction(googleToken, actionID, ctx)
	}
	return nil
}

func execGoogleDriveNewFileAction(token string, actionID string, ctx context.Context) {
	httpClient := &http.Client{
		Transport: &tokenTransport{token: token},
	}

	driveService, err := drive.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		fmt.Printf("Failed to create Drive service: %v\n", err)
		return
	}

	since := time.Now().Add(-1 * time.Minute).Format(time.RFC3339)
	query := fmt.Sprintf("createdTime > '%s' and trashed = false", since)

	fileList, err := driveService.Files.List().
		Q(query).
		PageSize(100).
		Fields("files(id, name, createdTime, mimeType)").
		Do()

	if err != nil {
		fmt.Printf("Failed to list files: %v\n", err)
		return
	}

	if len(fileList.Files) > 0 {
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

type tokenTransport struct {
	token string
}

func (t *tokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+t.token)
	return http.DefaultTransport.RoundTrip(req)
}
