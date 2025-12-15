package controller

import (
	// "net/http"
	//
	// "github.com/ValianceTekProject/AreaBack/db"
	// "github.com/ValianceTekProject/AreaBack/initializers"
	// "github.com/ValianceTekProject/AreaBack/model"
	"github.com/gin-gonic/gin"
)

func GetUserAreas(ctx *gin.Context) {
	// userCtx, exists := ctx.Get("user")
	// if !exists {
	// 	ctx.JSON(http.StatusUnauthorized, gin.H{"Error": "Unauthorized"})
	// 	return
	// }
	//
	// user := userCtx.(*db.UsersModel)
	//
	// areas, err := initializers.DB.Areas.FindMany(
	// 	db.Areas.User.Some(
	// 		db.Users.ID.Equals(user.ID),
	// 	),
	// ).With(
	// 	db.Areas.Actions.Fetch(),
	// 	db.Areas.Reactions.Fetch(),
	// ).Exec(ctx)
	//
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch areas"})
	// 	return
	// }
	//
	// ctx.JSON(http.StatusOK, areas)
}

func UpdateAreaStatus(ctx *gin.Context) {
	// areaID := ctx.Param("areaId")
	// if areaID == "" {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Area ID is required"})
	// 	return
	// }
	//
	// var payload model.AreaUpdateStatusPayload
	// if err := ctx.ShouldBindJSON(&payload); err != nil {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
	// 	return
	// }
	//
	// updatedArea, err := initializers.DB.Areas.FindUnique(
	// 	db.Areas.ID.Equals(areaID),
	// ).Update(
	// 	db.Areas.IsEnabled.Set(payload.IsEnabled),
	// ).Exec(ctx)
	//
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update area status", "details": err.Error()})
	// 	return
	// }
	//
	// ctx.JSON(http.StatusOK, updatedArea)
}
