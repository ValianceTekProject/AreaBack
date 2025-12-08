package authentification

import (
	"net/http"

	"github.com/ValianceTekProject/AreaBack/db"
	"github.com/ValianceTekProject/AreaBack/initializers"
	"github.com/ValianceTekProject/AreaBack/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(ctx *gin.Context) {
	var user model.User

	if ctx.ShouldBindJSON(&user) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid Request Body"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to hash password"})
		return
	}

	user.Password = string(hashedPassword)

	createdUser, err := initializers.DB.Users.CreateOne(
		db.Users.Email.Set(user.Email),
		db.Users.PasswordHash.Set(string(hashedPassword)),
	).Exec(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	tokenJWT, err := GenerateJWT(createdUser.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": "Token generation failed"})
		return
	}

	ctx.SetCookie("Authorization", tokenJWT, 3600 * 24 * 7, "/", "", false, true)

	ctx.Status(http.StatusCreated)
}
