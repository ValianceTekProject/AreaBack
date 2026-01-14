package controller

import (
	"net/http"

	"github.com/ValianceTekProject/AreaBack/db"
	"github.com/ValianceTekProject/AreaBack/initializers"
	"github.com/ValianceTekProject/AreaBack/model"
	"github.com/gin-gonic/gin"
)

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
