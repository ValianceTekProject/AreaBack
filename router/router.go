package router

import (
	"github.com/ValianceTekProject/AreaBack/authentification"
	"github.com/ValianceTekProject/AreaBack/middleware"
	"github.com/gin-gonic/gin"
)

func setupOAuth2Router(router *gin.Engine) *gin.Engine {
	router.GET("/auth/google/login", authentification.GoogleLogin)
	router.GET("/auth/google/callback", authentification.GoogleCallback)


	return router
}

func setupProtectedRouter(router *gin.Engine) *gin.Engine {
	protectedRoute := router.Group("/", )

	protectedRoute.Use(middleware.CheckUserAccess)
	{
		protectedRoute.GET("/zebi", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"message": "zebi I'm the route zebi",
			})
		})
	}
	return router
}

func SetupRouting() {
	router := gin.Default()

	router = setupOAuth2Router(router)
	router = setupProtectedRouter(router)

	router.Run(":8080")
}
