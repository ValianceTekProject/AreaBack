package routine

import (
	"context"
	"time"

	"github.com/ValianceTekProject/AreaBack/action"
	"github.com/ValianceTekProject/AreaBack/initializers"
)

func LaunchRoutines() {
    go func() {
        ticker := time.NewTicker(5 * time.Second)
        for range ticker.C {
            users, _ := initializers.DB.Users.FindMany().Exec(context.Background())
            
            for _, user := range users {
                action.GetGithubWebHook(user.ID)
            }
        }
    }()
}
