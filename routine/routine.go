package routine

import (
	"context"
	"log"
	"time"

	"github.com/ValianceTekProject/AreaBack/db"
	"github.com/ValianceTekProject/AreaBack/initializers"
	"github.com/ValianceTekProject/AreaBack/templates"
)

func LaunchRoutines() {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for range ticker.C {
			actions, _ := initializers.DB.Actions.FindMany().With(
				db.Actions.Service.Fetch(),
			).Exec(context.Background())

			for _, action := range actions {
				services := action.Service()

				Data := map[string]any{
					"action_id": action.ID,
				}

				handler, exists := templates.GetActionHandler(services.Name, action.Type)
				if !exists {
					log.Printf("Handler not found: %s/%s", action.Service, action.Type)
					continue
				}

				if err := handler(Data); err != nil {
					log.Printf("Error reaction: %v", err)
				}
			}
		}
	}()
}
