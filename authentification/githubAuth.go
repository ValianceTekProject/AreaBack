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
	"golang.org/x/oauth2/github"
)

var githubOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8080/auth/github/callback",
	ClientID:     os.Getenv("GH_BASIC_CLIENT_ID"),
	ClientSecret: os.Getenv("GH_BASIC_SECRET_ID"),
	Scopes: []string{
		"user:email",
	},
	Endpoint: github.Endpoint,
}

func GithubLogin(c *gin.Context) {
	state := generateStateOauthCookie()
	url := githubOauthConfig.AuthCodeURL(state)
	c.SetCookie("oauth_state", state, 300, "/", "", false, true)
	c.Redirect(302, url)
}

func GithubCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	randomState, err := c.Cookie("oauth_state")
	if err != nil || state != randomState {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid State"})
	}
	c.SetCookie("oauth_state", "", -1, "/", "", false, true)
	token, err := githubOauthConfig.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "OAuth exchange failed"})
		return
	}
	data, err := getUserData(token)

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to get user information"})
	}

	var githubUser model.GithubUserInfo
	if err := json.Unmarshal(data, &githubUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to parse user data"})
		return
	}

	user, err := saveOrUpdateGithubUser(githubUser, token)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to save user: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}

func getUserData(token *oauth2.Token) ([]byte, error){
	ctx := context.Background()

	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
	userDataURL := "https://api.github.com/user"
	req, _ := http.NewRequest("GET", userDataURL, nil)
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print("Request Failed:", err)
		return nil, fmt.Errorf("failed to send request: %s", err.Error());
	}
    defer resp.Body.Close()
    contents, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed read response: %s", err.Error())
    }
	return contents, nil
}

func saveOrUpdateGithubUser(info model.GithubUserInfo, token *oauth2.Token) (*db.UsersModel, error) {
	ctx := context.Background()
	serviceName := "Github"

	service, _ := initializers.DB.Services.FindUnique(
		db.Services.Name.Equals(serviceName),
	).Exec(ctx)

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
