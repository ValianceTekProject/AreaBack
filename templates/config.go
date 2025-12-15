package templates

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
}

type ReactionDefinition struct {
	Name        string
	Description string
	Service     string
	Config      []ReactionField
}

type Service struct {
	Name      string
	Actions   []ActionDefinition
	Reactions []ReactionDefinition
}

var Services = []Service{
	{
		Name: "Github",
		Actions: []ActionDefinition{
			{
				Name:        "github_new_issue",
				Description: "New issue",
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
			},
			{
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
			},
		},
		Reactions: []ReactionDefinition{},
	},
	{
		Name: "Discord",
		Actions: []ActionDefinition{
			{
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
			},
		},
		Reactions: []ReactionDefinition{
			{
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
			},
		},
	},
}
