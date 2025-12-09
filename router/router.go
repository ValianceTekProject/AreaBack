package router

import (
	"time"

	"github.com/ValianceTekProject/AreaBack/authentification"
	controller "github.com/ValianceTekProject/AreaBack/get"
	"github.com/ValianceTekProject/AreaBack/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func setupCORS(router *gin.Engine) {
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:8081"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	
	router.Use(cors.New(config))
}

func setupAuthRouter(router *gin.Engine) *gin.Engine {
	router.POST("/auth/login", authentification.LoginHandler)
	router.POST("/auth/register", authentification.RegisterHandler)

	router.GET("/auth/google/login", authentification.GoogleLogin)
	router.GET("/auth/google/callback", authentification.GoogleCallback)

	router.GET("/auth/github/login", authentification.GithubLogin)
	router.GET("/auth/github/callback", authentification.GithubCallback)

	router.GET("/auth/discord/login", authentification.DiscordLogin)
	router.GET("/auth/discord/callback", authentification.DiscordCallback)

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

		protectedRoute.GET("/areas", controller.GetUserAreas)
	}
	return router
}

func SetupRouting() {
	router := gin.Default()

	setupCORS(router)
	router = setupAuthRouter(router)
	router = setupProtectedRouter(router)

	router.Run(":8080")
}
