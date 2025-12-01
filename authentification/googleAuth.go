package authentification

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ValianceTekProject/AreaBack/db"
	"github.com/ValianceTekProject/AreaBack/initializers"
	"github.com/ValianceTekProject/AreaBack/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8080/auth/google/callback",
	ClientID:     "{PATTERN}.apps.googleusercontent.com",
	ClientSecret: "{SECRET}",
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/photoslibrary.readonly",
	},
	Endpoint: google.Endpoint,
}

func generateStateOauthCookie() string {
	file := make([]byte, 16)
	rand.Read(file)
	state := base64.URLEncoding.EncodeToString(file)

	return state
}

func getUserDataFromGoogle(code string) ([]byte, error) {
	userInfoURL := "https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken
	response, err := http.Get(userInfoURL)
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

func saveOrUpdateGoogleUser(info model.GoogleUserInfo, token *oauth2.Token) (*db.UsersModel, error) {

	ctx := context.Background()
    
    existingUser, err := initializers.DB.Users.FindUnique(
        db.Users.Email.Equals(info.Email),
    ).Exec(ctx)

	if err != nil {
    	return nil, err
	}

	if existingUser != nil {
		updatedServiceToken, err := initializers.DB.UserServiceTokens.FindUnique(
			db.UserServiceTokens.UserIDServiceIDEq(existingUser.ID, serviceID),
		).Update(
			db.UserServiceTokens.AccessToken.Set(token.AccessToken),
			db.UserServiceTokens.RefreshToken.Set(token.RefreshToken),
			db.UserServiceTokens.ExpiresAt.Set(time.Now().Add(time.Hour)),
		).Exec(ctx)
	}

    newUser, err := initializers.DB.Users.CreateOne(
        db.Users.Email.Set(googleUser.Email),
        db.Users.GoogleID.Set(googleUser.ID),
        db.Users.Name.Set(googleUser.Name),
        db.Users.Picture.Set(googleUser.Picture),
        db.Users.AccessToken.Set(token.AccessToken),
        db.Users.RefreshToken.Set(token.RefreshToken),
    ).Exec(ctx)

	return newUser, err
}

func GoogleLogin(ctx *gin.Context) {
	oauthState := generateStateOauthCookie()
	ctx.SetCookie("oauthState", oauthState, 3600, "/", "", true, true)
	url := googleOauthConfig.AuthCodeURL(oauthState, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(ctx *gin.Context) {
	state := ctx.Query("state")
	code := ctx.Query("code")

	cookieState, err := ctx.Cookie("oauthState")
	if err != nil || cookieState != state {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid oauth state"})
		return
	}

	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": "Code exchange failed"})
		return
	}

	data, err := getUserDataFromGoogle(code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	fmt.Fprintf(ctx.Writer, "UserInfo: %s\n", data)

	var googleUser model.GoogleUserInfo
	if err := json.Unmarshal(data, &googleUser); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to parse user data"})
		return
	}

	user, err := saveOrUpdateGoogleUser(googleUser, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to save user: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}
