package model

type DiscordConfig struct {
	WebhookURL string `json:"webhook_url"`
}

type DiscordMessage struct {
	Content string `json:"content"`
}
