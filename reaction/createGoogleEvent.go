package reaction

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func CreateGoogleEvent(config map[string]any) error {
	ctx := context.Background()

	reactionID, googleToken, err := GetGoogleToken(ctx, config)
	if err != nil {
		return fmt.Errorf("error getting Google token: %w", err)
	}

	if googleToken == "" {
		return fmt.Errorf("empty Google token")
	}

	if err != nil {
		log.Printf("Failed to get reactions: %v", err)
	}
	return execCreateEventAction(googleToken, reactionID, config, ctx)
}

func execCreateEventAction(token string, reactionID string, config map[string]any, ctx context.Context) error {
	httpClient := &http.Client{
		Transport: &tokenTransport{token: token},
	}

	calendarService, err := calendar.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		fmt.Printf("Failed to create Calendar service: %v\n", err)
		return err
	}

	summary, _ := config["summary"].(string)
	description, _ := config["description"].(string)
	location, _ := config["location"].(string)

	if summary == "" {
		summary = "AREA - Automated Event"
	}
	if description == "" {
		description = "This event was automatically created by AREA."
	}

	startTime := time.Date(2026, 4, 1, 10, 0, 0, 0, time.UTC)
	endTime := startTime.Add(1 * time.Hour)

	event := &calendar.Event{
		Summary:     summary,
		Location:    location,
		Description: description,
		Start: &calendar.EventDateTime{
			DateTime: startTime.Format(time.RFC3339),
			TimeZone: "UTC",
		},
		End: &calendar.EventDateTime{
			DateTime: endTime.Format(time.RFC3339),
			TimeZone: "UTC",
		},
	}

	createdEvent, err := calendarService.Events.Insert("primary", event).Do()
	if err != nil {
		fmt.Printf("Failed to create calendar event: %v\n", err)
		return err
	}

	log.Printf("Calendar event created successfully (reactionID: %s, EventID: %s)", reactionID, createdEvent.Id)
	return nil
}
