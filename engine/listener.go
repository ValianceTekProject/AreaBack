package engine

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ValianceTekProject/AreaBack/db"
	"github.com/ValianceTekProject/AreaBack/initializers"
)

func Listener() {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for range ticker.C {
			ctx := context.Background()

			areas, err := initializers.DB.Areas.FindMany().With(
				db.Areas.User.Fetch(),
				db.Areas.Actions.Fetch(),
			).Exec(ctx)
			if err != nil {
				log.Printf("Error fetching areas: %v", err)
				continue
			}

			for _, area := range areas {
				if area.IsEnabled {
					if area.Name == "Github_pr_to_discord" {
						users := area.User()
						actions := area.Actions()

						for _, user := range users {
							fmt.Printf("Area: %s, User: %s (%s)\n",
								area.Name, user.Email, user.ID)
						}

						for _, action := range actions {
							if action.Triggered {
								fmt.Printf("Action %s is triggered!\n", action.ID)

								_, err = initializers.DB.Actions.FindUnique(
									db.Actions.ID.Equals(action.ID),
								).Update(
									db.Actions.Triggered.Set(false),
								).Exec(ctx)
								if err != nil {
									log.Printf("Error resetting action trigger: %v", err)
								}
							}
						}
					}
				}
			}

		}
	}()
}
