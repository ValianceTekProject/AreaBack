package reaction

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func CreateSheet(config map[string]any) error {
	ctx := context.Background()

	reactionID, googleToken, err := GetGoogleToken(ctx, config)
	if err != nil {
		return fmt.Errorf("error getting Google token: %w", err)
	}

	if googleToken == "" {
		return fmt.Errorf("empty Google token")
	}

	return execCreateSheetAction(googleToken, reactionID, ctx)
}

func execCreateSheetAction(token string, reactionID string, ctx context.Context) error {
	httpClient := &http.Client{
		Transport: &tokenTransport{token: token},
	}

	sheetsService, err := sheets.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		fmt.Printf("Failed to create Sheets service: %v\n", err)
		return err
	}

	spreadsheet := &sheets.Spreadsheet{
		Properties: &sheets.SpreadsheetProperties{
			Title: "Action Detected - AREA",
		},
		Sheets: []*sheets.Sheet{
			{
				Properties: &sheets.SheetProperties{
					Title: "Sheet1",
				},
				Data: []*sheets.GridData{
					{
						RowData: []*sheets.RowData{
							{
								Values: []*sheets.CellData{
									{
										UserEnteredValue: &sheets.ExtendedValue{
											StringValue: func(s string) *string { return &s }("Action Detected"),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	createdSpreadsheet, err := sheetsService.Spreadsheets.Create(spreadsheet).Do()
	if err != nil {
		fmt.Printf("Failed to create spreadsheet: %v\n", err)
		return err
	}

	log.Printf("Spreadsheet created successfully (ID: %s, URL: %s) for reactionID: %s",
		createdSpreadsheet.SpreadsheetId,
		createdSpreadsheet.SpreadsheetUrl,
		reactionID)

	return nil
}
