package model

import "time"

type Service struct {
	ID   int    `json:"id"`
	Name string `json:"name"`

	Actions   []Actions   `json:"actions"`
	Reactions []Reactions `json:"reactions"`
}

type Areas struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	IsEnabled bool      `json:"is_enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Actions   []Actions   `json:"actions"`
	Reactions []Reactions `json:"reactions"`
}

type Actions struct {
	ID              string         `json:"id"`
	AreaID          string         `json:"area_id"`
	ServiceID       int            `json:"service_id"`
	Type            string         `json:"type"`
	Config          map[string]any `json:"config"`
	LastTriggeredAt *time.Time     `json:"last_triggered_at"`
	LastState       map[string]any `json:"last_state"`

	Area    *Areas   `json:"area"`
	Service *Service `json:"service"`
}

type Reactions struct {
	ID        string         `json:"id"`
	AreaID    string         `json:"area_id"`
	ServiceID int            `json:"service_id"`
	Type      string         `json:"type"`
	Config    map[string]any `json:"config"`

	Area    *Areas   `json:"area"`
	Service *Service `json:"service"`
}

type AreaUpdateStatusPayload struct {
	IsEnabled bool `json:"is_enabled"`
}
