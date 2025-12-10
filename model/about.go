package model

type ActionResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ReactionResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ServiceResponse struct {
	Name      string             `json:"name"`
	Actions   []ActionResponse   `json:"actions"`
	Reactions []ReactionResponse `json:"reactions"`
}

type ServerResponse struct {
	CurrentTime int64             `json:"current_time"`
	Services    []ServiceResponse `json:"services"`
}

type ClientResponse struct {
	Host string `json:"host"`
}

type AboutResponse struct {
	Client ClientResponse `json:"client"`
	Server ServerResponse `json:"server"`
}
