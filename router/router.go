package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ValianceTekProject/AreaBack/authentification"
)

func setupOAuth2Router(router *gin.Engine) *gin.Engine {
	router.GET("/auth/google/login", authentification.GoogleLogin)
	router.GET("/auth/google/callback", authentification.GoogleCallback)

	return router
}

func SetupRouting() {
	router := gin.Default()

	router = setupOAuth2Router(router)
	router.Run(":8080")
}
