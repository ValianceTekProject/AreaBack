package authentification

import (
	"os"
	"time"
	"crypto/rand"
	"encoding/base64"

	"github.com/golang-jwt/jwt/v5"
)


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

