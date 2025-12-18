package engine

import (
	"net/http"
	"time"

	"github.com/ValianceTekProject/AreaBack/model"
	"github.com/ValianceTekProject/AreaBack/templates"
	"github.com/gin-gonic/gin"
)

func GetAbout(c *gin.Context) {
	availableServices := templates.Services
	servicesResponse := []model.ServiceResponse{}

	for _, service := range availableServices {
		services := model.ServiceResponse{}
		services.Name = service.Name;
		for _, action := range service.Actions {
			newAction := model.ActionResponse{}
			newAction.Name = action.Name
			newAction.Description = action.Description
			services.Actions = append(services.Actions, newAction)
		}
		for _, reaction := range service.Reactions {
			newReaction := model.ReactionResponse{}
			newReaction.Name = reaction.Name
			newReaction.Description = reaction.Description
			services.Reactions = append(services.Reactions, newReaction)
		}
		servicesResponse = append(servicesResponse, services)
	}

	response := model.AboutResponse{
		Client: model.ClientResponse{
			Host: c.ClientIP(),
		},
		Server: model.ServerResponse{
			CurrentTime: time.Now().Unix(),
			Services:    servicesResponse,
		},
	}

	c.JSON(http.StatusOK, response)
}
