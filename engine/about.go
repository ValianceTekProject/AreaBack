package engine

import (
	"context"
	"net/http"
	"time"

	"github.com/ValianceTekProject/AreaBack/db"
	"github.com/ValianceTekProject/AreaBack/initializers"
	"github.com/ValianceTekProject/AreaBack/model"
	"github.com/gin-gonic/gin"
)

func GetAbout(c *gin.Context) {
	ctx := context.Background()

	services, err := initializers.DB.Services.FindMany().With(
		db.Services.Actions.Fetch(),
		db.Services.Reactions.Fetch(),
	).Exec(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch services"})
		return
	}

	var serviceResponses []model.ServiceResponse
	for _, service := range services {
		var actions []model.ActionResponse
		var reactions []model.ReactionResponse

		if service.Actions() != nil {
			for _, action := range service.Actions() {
				actions = append(actions, model.ActionResponse{
					Name:        action.ID,
					Description: "Action description",
				})
			}
		}

		if service.Reactions() != nil {
			for _, reaction := range service.Reactions() {
				reactions = append(reactions, model.ReactionResponse{
					Name:        reaction.ID,
					Description: "Reaction description",
				})
			}
		}

		serviceResponses = append(serviceResponses, model.ServiceResponse{
			Name:      service.Name,
			Actions:   actions,
			Reactions: reactions,
		})
	}

	response := model.AboutResponse{
		Client: model.ClientResponse{
			Host: c.ClientIP(),
		},
		Server: model.ServerResponse{
			CurrentTime: time.Now().Unix(),
			Services:    serviceResponses,
		},
	}

	c.JSON(http.StatusOK, response)
}
