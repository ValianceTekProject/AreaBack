package router

import (
	"time"

	"github.com/ValianceTekProject/AreaBack/authentification"
	"github.com/ValianceTekProject/AreaBack/controller"
	"github.com/ValianceTekProject/AreaBack/download"
	"github.com/ValianceTekProject/AreaBack/engine"
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
	oauthRoute := router.Group("/")

	oauthRoute.Use(middleware.VerifyOauthUser)
	{
		oauthRoute.GET("/auth/google/login", authentification.GoogleLogin)
		oauthRoute.GET("/auth/github/login", authentification.GithubLogin)
		oauthRoute.GET("/auth/discord/login", authentification.DiscordLogin)
	}

	router.POST("/auth/login", authentification.LoginHandler)
	router.POST("/auth/register", authentification.RegisterHandler)

	router.GET("/auth/google/callback", authentification.GoogleCallback)

	router.GET("/auth/github/callback", authentification.GithubCallback)

	router.GET("/auth/discord/callback", authentification.DiscordCallback)

	router.GET("/about.json", engine.GetAbout)
	router.GET("/client.apk", download.DownloadApk)

	return router
}

func setupProtectedRouter(router *gin.Engine) *gin.Engine {
	protectedRoute := router.Group("/")

	protectedRoute.Use(middleware.CheckUserAccess)
	{
		protectedRoute.GET("/zebi", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"message": "zebi I'm the route zebi",
			})
		})
		protectedRoute.POST("/areas/create", controller.CreateArea)

		protectedRoute.GET("/areas", controller.GetUserAreas)
		protectedRoute.GET("/me/userId", controller.GetSelfUserId)
		protectedRoute.GET("/users/:userId", controller.GetSpecificUser)
		protectedRoute.GET("/services/:userId", controller.GetUserServices)
		protectedRoute.PATCH("/areas/:areaId/status", controller.UpdateAreaStatus)
		protectedRoute.DELETE("/areas/:areaId/delete", controller.DeleteArea)
		protectedRoute.POST("/areas/:areaId/action/add", controller.LinkAction)
		protectedRoute.POST("/areas/:areaId/reaction/add", controller.LinkReaction)
	}
	protectedRoute.Use(middleware.CheckAdminAccess)
	{
		protectedRoute.PATCH("/users/:userId/status", controller.UpdateUserStatus)
		protectedRoute.GET("/users", controller.GetUsers)
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
