package controller

import (
	"net/http"

	"github.com/ValianceTekProject/AreaBack/db"
	"github.com/ValianceTekProject/AreaBack/initializers"
	"github.com/ValianceTekProject/AreaBack/model"
	"github.com/ValianceTekProject/AreaBack/templates"
	"github.com/gin-gonic/gin"
)

func GetUserAreas(ctx *gin.Context) {
	userInterface, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	user, ok := userInterface.(*db.UsersModel)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type"})
		return
	}

	areas, err := initializers.DB.Areas.FindMany(
		db.Areas.UserID.Equals(user.ID),
	).With(
		db.Areas.Action.Fetch(),
		db.Areas.Reaction.Fetch(),
	).Exec(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch areas"})
		return
	}

	ctx.JSON(http.StatusOK, areas)
}

func CreateArea(ctx *gin.Context) {
	userInterface, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	user, ok := userInterface.(*db.UsersModel)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type"})
		return
	}

	var input model.AreaInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	area, err := initializers.DB.Areas.CreateOne(
		db.Areas.Name.Set(input.Name),
		db.Areas.User.Link(
			db.Users.ID.Equals(user.ID),
		),
	).Exec(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create area"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Area created successfully",
		"area":    area,
	})
}

func LinkAction(ctx *gin.Context) {
	areaID := ctx.Param("areaId")
	if areaID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Area ID is required"})
		return
	}

	userInterface, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	user, ok := userInterface.(*db.UsersModel)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type"})
		return
	}

	var addAction model.AddAction

	if err := ctx.ShouldBindJSON(&addAction); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	action, exists := templates.GetAction(addAction.ServiceName, addAction.Name)
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Action not found"})
		return
	}

	area, err := initializers.DB.Areas.FindFirst(
		db.Areas.ID.Equals(areaID),
		db.Areas.UserID.Equals(user.ID),
	).Exec(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Area not found or unauthorized"})
		return
	}

	initializers.DB.Actions.CreateOne(
		db.Actions.Type.Set(action.Name),
		db.Actions.Area.Link(
			db.Areas.ID.Equals(area.ID),
		),
		db.Actions.Service.Link(
			db.Services.Name.Equals(addAction.ServiceName),
		),
		db.Actions.Triggered.Set(false),
	).Exec(ctx)
}

func LinkReaction(ctx *gin.Context) {
	areaID := ctx.Param("areaId")
	if areaID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Area ID is required"})
		return
	}

	userInterface, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	user, ok := userInterface.(*db.UsersModel)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type"})
		return
	}

	var addAction model.AddAction

	if err := ctx.ShouldBindJSON(&addAction); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	action, exists := templates.GetReaction(addAction.ServiceName, addAction.Name)
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Action not found"})
		return
	}

	area, err := initializers.DB.Areas.FindFirst(
		db.Areas.ID.Equals(areaID),
		db.Areas.UserID.Equals(user.ID),
	).Exec(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Area not found or unauthorized"})
		return
	}

	initializers.DB.Reactions.CreateOne(
		db.Reactions.Type.Set(action.Name),
		db.Reactions.Area.Link(
			db.Areas.ID.Equals(area.ID),
		),
		db.Reactions.Service.Link(
			db.Services.Name.Equals(addAction.ServiceName),
		),
	).Exec(ctx)
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
