package middleware

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func verifyToken(secret []byte) jwt.Keyfunc {
	return func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	}
}

func ValidateJWTToken(tokenString string) (jwt.MapClaims, error) {
	secret := []byte(os.Getenv(("JWT_SECRET")))

	token, err := jwt.Parse(tokenString, verifyToken(secret))
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func VerifyOauthUser(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.Next()
		return
	}

	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		ctx.Next()
		return
	}

	token := strings.TrimPrefix(authHeader, bearerPrefix)

	claims, err := ValidateJWTToken(token)
    if err != nil {
		ctx.Next()
        return
    }
	userID, ok := claims["sub"].(string)
	if !ok {
		ctx.Next()
		return
	}
	ctx.Set("userID", userID)
	ctx.Next()
}
