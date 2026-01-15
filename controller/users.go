package controller

import (
	"net/http"

	"github.com/ValianceTekProject/AreaBack/db"
	"github.com/ValianceTekProject/AreaBack/initializers"
	"github.com/ValianceTekProject/AreaBack/model"
	"github.com/gin-gonic/gin"
)

func GetSpecificUser(ctx *gin.Context) {

	userId := ctx.Param("userId")

	if userId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	users, err := initializers.DB.Users.FindMany(
		db.Users.ID.Equals(userId),
	).With().Exec(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch users"})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func GetUsers(ctx *gin.Context) {
	users, err := initializers.DB.Users.FindMany(
	).Exec(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch users"})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func UpdateUserStatus(ctx *gin.Context) {
	userId := ctx.Param("userId")

	if userId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}
	var payload model.UserUpdateStatusPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	updatedUser, err := initializers.DB.Users.FindUnique(
		db.Users.ID.Equals(userId),
	).Update(
		db.Users.Authorized.Set(payload.Authorized),
	).Exec(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user authorization", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedUser)
}

func GetUserServices(ctx *gin.Context) {
	userId := ctx.Param("userId")

	if userId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	user, err := initializers.DB.Users.FindUnique(
		db.Users.ID.Equals(userId),
	).Exec(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	serviceTokens, err := initializers.DB.UserServiceTokens.FindMany(
		db.UserServiceTokens.UserID.Equals(user.ID),
	).With(
		db.UserServiceTokens.Service.Fetch(),
	).Exec(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user services"})
		return
	}

	ctx.JSON(http.StatusOK, serviceTokens)
}

func GetSelfUserId(ctx *gin.Context) {
	userInterface, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user := userInterface.(*db.UsersModel)
	ctx.JSON(http.StatusOK, gin.H{"userId": user.ID})
}
