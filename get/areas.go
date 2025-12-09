package controller

import (
	"net/http"

	"github.com/ValianceTekProject/AreaBack/db"
	"github.com/ValianceTekProject/AreaBack/initializers"
	"github.com/gin-gonic/gin"
)

func GetUserAreas(ctx *gin.Context) {
	userCtx, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"Error": "Unauthorized"})
		return
	}
	
	user := userCtx.(*db.UsersModel)

	areas, err := initializers.DB.Areas.FindMany(
		db.Areas.User.Some(
			db.Users.ID.Equals(user.ID),
		),
	).With(
		db.Areas.Actions.Fetch(),
		db.Areas.Reactions.Fetch(),
	).Exec(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch areas"})
		return
	}

	ctx.JSON(http.StatusOK, areas)
}
