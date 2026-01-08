package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/ValianceTekProject/AreaBack/db"
	"github.com/ValianceTekProject/AreaBack/initializers"
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

func CheckUserAccess(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, bearerPrefix)

	claims, err := ValidateJWTToken(token)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	user, err := initializers.DB.Users.FindUnique(
		db.Users.ID.Equals(userID),
	).Exec(ctx)

	if err != nil || user == nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	ctx.Set("user", user)
	ctx.Next()
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

	_, err := ValidateJWTToken(token)
    if err != nil {
        ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }
	ctx.Next()
}
