package middleware

import (
	"net/http"
	"strings"

	"github.com/ValianceTekProject/AreaBack/db"
	"github.com/ValianceTekProject/AreaBack/initializers"
	"github.com/gin-gonic/gin"
)

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

	if err != nil || user == nil || (user.Authorized != true && user.Admin != true){
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found or not authorized"})
		return
	}

	ctx.Set("user", user)
	ctx.Next()
}
