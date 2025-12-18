package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ValianceTekProject/AreaBack/db"
	"github.com/ValianceTekProject/AreaBack/initializers"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func ValidateJWTToken(tokenString string, ctx *gin.Context) bool {
	secret := []byte(os.Getenv(("JWT_SECRET")))

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return false
		}

		userID := claims["sub"].(string)

		user, err := initializers.DB.Users.FindUnique(
			db.Users.ID.Equals(userID),
		).Exec(ctx)

		if err != nil || user == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return false
		}
		ctx.Set("user", user)
		return true
	}
	ctx.AbortWithStatus(http.StatusUnauthorized)
	return false
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
	
	token := authHeader[len(bearerPrefix):]

	validated := ValidateJWTToken(token, ctx)

	if validated != true {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error": "Unauthorized access"})
		return
	}

	ctx.Next()
}
