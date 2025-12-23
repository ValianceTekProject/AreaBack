package authentification

import (
	// "context"
	"crypto/rand"
	"encoding/base64"
	// "log"
	"os"
	"time"

	// "github.com/ValianceTekProject/AreaBack/db"
	// "github.com/ValianceTekProject/AreaBack/initializers"
	"github.com/golang-jwt/jwt/v5"
)

func LinkAllUsersToAreas() error {
	// ctx := context.Background()
	//
	// areas, err := initializers.DB.Areas.FindMany().With(
	// 	db.Areas.User.Fetch(),
	// ).Exec(ctx)
	// if err != nil {
	// 	return err
	// }
	//
	// users, err := initializers.DB.Users.FindMany().Exec(ctx)
	// if err != nil {
	// 	return err
	// }
	//
	// for _, area := range areas {
	// 	existingUsers := make(map[string]bool)
	// 	for _, u := range area.User() {
	// 		existingUsers[u.ID] = true
	// 	}
	//
	// 	var updates []db.AreasSetParam
	// 	for _, user := range users {
	// 		if !existingUsers[user.ID] {
	// 			updates = append(updates, 
	// 				db.Areas.User.Link(
	// 					db.Users.ID.Equals(user.ID),
	// 				),
	// 			)
	// 		}
	// 	}
	//
	// 	if len(updates) > 0 {
	// 		_, err = initializers.DB.Areas.FindUnique(
	// 			db.Areas.ID.Equals(area.ID),
	// 		).Update(updates...).Exec(ctx)
	//
	// 		if err != nil {
	// 			log.Printf("error linking users to area %s: %v", area.ID, err)
	// 		}
	// 	}
	// }
	//
	return nil
}

func GenerateJWT(userID string) (string, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET"))
	claim := jwt.MapClaims {
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	token_str, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return token_str, nil
}

func generateStateOauthCookie() string {
	file := make([]byte, 16)
	rand.Read(file)
	state := base64.URLEncoding.EncodeToString(file)

	return state
}

