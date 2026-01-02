package engine

import (
	"context"
	"log"
	"time"

	"github.com/ValianceTekProject/AreaBack/db"
	"github.com/ValianceTekProject/AreaBack/initializers"
	"github.com/ValianceTekProject/AreaBack/templates"
)

func Listener() {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for range ticker.C {
			ctx := context.Background()

			areas, err := initializers.DB.Areas.FindMany().With(
				db.Areas.User.Fetch(),
				db.Areas.Reaction.Fetch().With(
					db.Reactions.Service.Fetch(),
				),
				db.Areas.Action.Fetch().With(
					db.Actions.Service.Fetch(),
				),
			).Exec(ctx)
			if err != nil {
				log.Printf("Error fetching areas: %v", err)
				continue
			}

			for _, area := range areas {
				if area.IsEnabled {
					handleEnableArea(&area, ctx)
				}
			}

		}
	}()
}

func handleEnableArea(area *db.AreasModel, ctx context.Context) {
	action, exist := area.Action()
	if !exist {
		return
	}
	
	reaction, exist := area.Reaction()
	if !exist {
		return
	}
	service := reaction.Service()
	if action.Triggered {
		handler, exists := templates.GetReactionHandler(service.Name, reaction.Type)

		if !exists {
			log.Printf("Handler not found: %s/%s", action.Service, action.Type)
		}

		if err := handler(nil); err != nil {
			log.Printf("Error reaction: %v", err)
		}

		_, err := initializers.DB.Actions.FindUnique(
			db.Actions.ID.Equals(action.ID),
		).Update(
			db.Actions.Triggered.Set(false),
		).Exec(ctx)
		if err != nil {
			log.Printf("Error resetting action trigger: %v", err)
		}
	}
}
