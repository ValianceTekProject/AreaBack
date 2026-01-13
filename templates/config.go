package templates

import (
	"github.com/ValianceTekProject/AreaBack/action"
	"github.com/ValianceTekProject/AreaBack/reaction"
)

type (
	ActionHandler   func(config map[string]any) error
	ReactionHandler func(config map[string]any) error
)

type ActionField struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Label    string `json:"label"`
	Required bool   `json:"required"`
}

type ReactionField struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Label    string `json:"label"`
	Required bool   `json:"required"`
}

type ActionDefinition struct {
	Name        string
	Description string
	Service     string
	Config      []ActionField
	Handler     ActionHandler
}

type ReactionDefinition struct {
	Name        string
	Description string
	Service     string
	Config      []ReactionField
	Handler     ReactionHandler
}

type Service struct {
	Name      string
	Actions   map[string]*ActionDefinition
	Reactions map[string]*ReactionDefinition
}

var Services = map[string]*Service{
	"Github": {
		Name: "Github",
		Actions: map[string]*ActionDefinition{
			"github_pr_merged": {
				Name:        "github_pr_merged",
				Description: "Pr closed",
				Service:     "Github",
				Config: []ActionField{
					{
						Name:     "owner",
						Type:     "text",
						Label:    "Repository Owner",
						Required: true,
					},
					{
						Name:     "repo",
						Type:     "text",
						Label:    "Repository Name",
						Required: true,
					},
				},
				Handler: action.ExecGithubPR,
			},
			"github_new_PR": {
				Name:        "github_new_PR",
				Description: "New PR",
				Service:     "Github",
				Config: []ActionField{
					{
						Name:     "owner",
						Type:     "text",
						Label:    "Repository Owner",
						Required: true,
					},
					{
						Name:     "repo",
						Type:     "text",
						Label:    "Repository Name",
						Required: true,
					},
				},
				Handler: action.ExecGithubNewPr,
			},
		},
		Reactions: map[string]*ReactionDefinition{
			"github_create_pr": {
				Name:        "github_create_pr",
				Description: "Create Pr in TestArea Repository",
				Service:     "Github",
				Config: []ReactionField{
					{
						Name:     "channel_id",
						Type:     "text",
						Label:    "Channel ID",
						Required: true,
					},
					{
						Name:     "content",
						Type:     "text",
						Label:    "Message Content",
						Required: true,
					},
				},
				Handler: reaction.CreateGihubPr,
			},
		},
	},

	"Discord": {
		Name: "Discord",
		Actions: map[string]*ActionDefinition{
			"discord_message_received": {
				Name:        "discord_message_received",
				Description: "Message received in a specific channel",
				Service:     "Discord",
				Config: []ActionField{
					{
						Name:     "channel_id",
						Type:     "text",
						Label:    "Channel ID",
						Required: true,
					},
				},
				Handler: action.ExecDiscordNewMsg,
			},
		},
		Reactions: map[string]*ReactionDefinition{
			"discord_send_message": {
				Name:        "discord_send_message",
				Description: "Send a message to a channel",
				Service:     "Discord",
				Config: []ReactionField{
					{
						Name:     "channel_id",
						Type:     "text",
						Label:    "Channel ID",
						Required: true,
					},
					{
						Name:     "content",
						Type:     "text",
						Label:    "Message Content",
						Required: true,
					},
				},
				Handler: reaction.ReactWithDiscordMsg,
			},
		},
	},

	"Steam": {
		Name: "Steam",
		Actions: map[string]*ActionDefinition{
			"steam_valiance_playing": {
				Name:        "steam_valiance_playing",
				Description: "Valiance is playing a game",
				Service:     "Steam",
				Config: []ActionField{
					{
						Name:     "channel_id",
						Type:     "text",
						Label:    "Channel ID",
						Required: false,
					},
				},
				Handler: action.ExecSteamPlaying,
			},
		},
	},
	"Google": {
		Name: "Google",
		Actions: map[string]*ActionDefinition{
			"google_new_file_drive": {
				Name:        "google_new_file_drive",
				Description: "New file created in drive",
				Service:     "Google",
				Config: []ActionField{
					{
						Name:     "channel_id",
						Type:     "text",
						Label:    "Channel ID",
						Required: false,
					},
				},
				Handler: action.ExecGoogleDriveNewFile,
			},
		},
		Reactions: map[string]*ReactionDefinition{
			"google_Send_Email": {
				Name:        "google_Send_Email",
				Description: "Send email to itself",
				Service:     "Google",
				Config: []ReactionField{
					{
						Name:     "content",
						Type:     "text",
						Label:    "Message Content",
						Required: true,
					},
				},
				Handler: reaction.SendEmail,
			},
			"google_Create_Sheet": {
				Name:        "google_Create_Sheet",
				Description: "Create a google Sheet",
				Service:     "Google",
				Config: []ReactionField{
					{
						Name:     "content",
						Type:     "text",
						Label:    "Message Content",
						Required: true,
					},
				},
				Handler: reaction.CreateSheet,
			},
			"google_Create_Event": {
				Name:        "google_Create_Event",
				Description: "Create a google Event",
				Service:     "Google",
				Config: []ReactionField{
					{
						Name:     "content",
						Type:     "text",
						Label:    "Message Content",
						Required: true,
					},
				},
				Handler: reaction.CreateGoogleEvent,
			},
		},
	},
	"Twitch": {
		Name: "Twitch",
		Actions: map[string]*ActionDefinition{
			"twitch_streamer_live": {
				Name:        "twitch_streamer_live",
				Description: "Streamer is live",
				Service:     "Twitch",
				Config: []ActionField{
					{
						Name:     "streamer_name",
						Type:     "text",
						Label:    "Streamer Username",
						Required: true,
					},
				},
				Handler: action.ExecTwitchLive,
			},
		},
		Reactions: map[string]*ReactionDefinition{},
	},

	"OpenWeather": {
		Name: "OpenWeather",
		Actions: map[string]*ActionDefinition{
			"openweather_rain": {
				Name:        "openweather_rain",
				Description: "It rains",
				Service:     "OpenWeather",
				Config: []ActionField{
					{
						Name:     "streamer_name",
						Type:     "text",
						Label:    "Streamer Username",
						Required: true,
					},
				},
				Handler: action.ExecWeatherRain,
			},
		},
		Reactions: map[string]*ReactionDefinition{},
	},
}

func GetService(name string) (*Service, bool) {
	service, exists := Services[name]
	return service, exists
}

func GetAction(serviceName, actionName string) (*ActionDefinition, bool) {
	service, exists := Services[serviceName]
	if !exists {
		return nil, false
	}
	action, exists := service.Actions[actionName]
	return action, exists
}

func GetReaction(serviceName, reactionName string) (*ReactionDefinition, bool) {
	service, exists := Services[serviceName]
	if !exists {
		return nil, false
	}
	reaction, exists := service.Reactions[reactionName]
	return reaction, exists
}

func GetReactionHandler(serviceName, reactionName string) (ReactionHandler, bool) {
	service, exists := Services[serviceName]
	if !exists {
		return nil, false
	}
	reaction, exists := service.Reactions[reactionName]
	if !exists {
		return nil, false
	}
	return reaction.Handler, true
}

func GetActionHandler(serviceName, actionName string) (ActionHandler, bool) {
	service, exists := Services[serviceName]
	if !exists {
		return nil, false
	}
	action, exists := service.Actions[actionName]
	if !exists {
		return nil, false
	}
	return action.Handler, true
}
