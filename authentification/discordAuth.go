package authentification

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ValianceTekProject/AreaBack/db"
	"github.com/ValianceTekProject/AreaBack/initializers"
	"github.com/ValianceTekProject/AreaBack/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func getDiscordOAuthConfig() *oauth2.Config {
	clientID := os.Getenv("DISCORD_CLIENT_ID")
	clientSecret := os.Getenv("DISCORD_CLIENT_SECRET")
	redirectURL := os.Getenv("DISCORD_REDIRECT_URL")

	return &oauth2.Config{
		RedirectURL:  redirectURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes: []string{
			"identify",
			"email",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://discord.com/api/oauth2/authorize",
			TokenURL: "https://discord.com/api/oauth2/token",
		},
	}
}

var discordOauthConfig *oauth2.Config = getDiscordOAuthConfig()

func getUserDataFromDiscord(token *oauth2.Token) ([]byte, error) {
	userInfoURL := "https://discord.com/api/users/@me"

	client := discordOauthConfig.Client(context.Background(), token)

	response, err := client.Get(userInfoURL)
	if err != nil {
		return nil, fmt.Errorf("Failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed read response: %s", err.Error())
	}

	return contents, nil
}

func saveOrUpdateDiscordUser(info model.DiscordUserInfo, token *oauth2.Token) (*db.UsersModel, error) {
	ctx := context.Background()
	serviceName := "Discord"

	service, err := initializers.DB.Services.FindUnique(
		db.Services.Name.Equals(serviceName),
	).Exec(ctx)

	if err != nil {
		return nil, fmt.Errorf("Service '%s' not found in db", serviceName)
	}

	existingUser, err := initializers.DB.Users.FindUnique(
		db.Users.Email.Equals(info.Email),
	).Exec(ctx)

	if err == nil && existingUser != nil {
		existingToken, _ := initializers.DB.UserServiceTokens.FindFirst(
			db.UserServiceTokens.UserID.Equals(existingUser.ID),
			db.UserServiceTokens.ServiceID.Equals(service.ID),
		).Exec(ctx)

		if existingToken != nil {
			_, err = initializers.DB.UserServiceTokens.FindUnique(
				db.UserServiceTokens.ID.Equals(existingToken.ID),
			).Update(
				db.UserServiceTokens.AccessToken.Set(token.AccessToken),
				db.UserServiceTokens.RefreshToken.Set(token.RefreshToken),
				db.UserServiceTokens.ExpiresAt.Set(token.Expiry),
			).Exec(ctx)
		} else {
			_, err = initializers.DB.UserServiceTokens.CreateOne(
				db.UserServiceTokens.AccessToken.Set(token.AccessToken),
				db.UserServiceTokens.RefreshToken.Set(token.RefreshToken),
				db.UserServiceTokens.ExpiresAt.Set(token.Expiry),
				db.UserServiceTokens.User.Link(db.Users.ID.Equals(existingUser.ID)),
				db.UserServiceTokens.Service.Link(db.Services.ID.Equals(service.ID)),
			).Exec(ctx)
		}

		if err != nil {
			return nil, err
		}

		return existingUser, nil
	}

	newUser, err := initializers.DB.Users.CreateOne(
		db.Users.Email.Set(info.Email),
		db.Users.PasswordHash.Set(""),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	_, err = initializers.DB.UserServiceTokens.CreateOne(
		db.UserServiceTokens.AccessToken.Set(token.AccessToken),
		db.UserServiceTokens.RefreshToken.Set(token.RefreshToken),
		db.UserServiceTokens.ExpiresAt.Set(token.Expiry),
		db.UserServiceTokens.User.Link(db.Users.ID.Equals(newUser.ID)),
		db.UserServiceTokens.Service.Link(db.Services.ID.Equals(service.ID)),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func DiscordLogin(ctx *gin.Context) {
	oauthState := generateStateOauthCookie()
	ctx.SetCookie("oauthState", oauthState, 3600, "/", "", false, true)
	
	url := discordOauthConfig.AuthCodeURL(oauthState, oauth2.AccessTypeOffline)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func DiscordCallback(ctx *gin.Context) {
	state := ctx.Query("state")
	code := ctx.Query("code")

	cookieState, err := ctx.Cookie("oauthState")
	if err != nil || cookieState != state {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid oauth state"})
		return
	}

	token, err := discordOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": "Code exchange failed"})
		return
	}

	data, err := getUserDataFromDiscord(token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	var discordUser model.DiscordUserInfo
	if err := json.Unmarshal(data, &discordUser); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to parse user data"})
		return
	}

	user, err := saveOrUpdateDiscordUser(discordUser, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to save user: " + err.Error()})
		return
	}

	tokenJWT, err := GenerateJWT(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": "Token generation failed: " + err.Error()})
		return
	}

	ctx.SetCookie("Authorization", tokenJWT, 3600*24*7, "/", "", false, true)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   tokenJWT,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}
