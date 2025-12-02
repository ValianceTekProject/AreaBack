package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ValidateJWTToken(tokenString string) (string, error) {
	return "", nil
}

func CheckUserAccess(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")

	if authHeader == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error": "Authorization header is missing"})
		return
	}

	tokenString := strings.Split(authHeader, "Bearer ")[1]

	_, err := ValidateJWTToken(tokenString)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error": "Unauthorized access"})
		return
	}

	ctx.Next()
}
