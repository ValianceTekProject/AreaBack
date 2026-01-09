package authentification

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

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
		"repo",
	},
	Endpoint: github.Endpoint,
}

func GithubLogin(c *gin.Context) {
	state := generateStateOauthCookie()
	userID, exist := c.Get("userID")

	if exist {
		c.SetCookie("oauth_user_id", userID.(string), 300, "/", "", false, true)
	}
	url := githubOauthConfig.AuthCodeURL(state)
	c.SetCookie("oauth_state", state, 300, "/", "", false, true)
	c.Redirect(302, url)
}

func GithubCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	currentUserID, err := c.Cookie("oauth_user_id")
	var userIDPtr *string
	if err == nil && currentUserID != "" {
		userIDPtr = &currentUserID
	}
	c.SetCookie("oauth_user_id", "", -1, "/", "", false, true)

	randomState, err := c.Cookie("oauth_state")
	if err != nil || state != randomState {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid State"})
		return
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
		return
	}

	emailData, err := getUserEmail(token)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to get user mail"})
		return
	}

	var githubUser model.GithubUserInfo
	if json.Unmarshal(data, &githubUser) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user data"})
		return
	}

	var emails []struct {
		Email      string `json:"email"`
		Primary    bool   `json:"primary"`
		Verified   bool   `json:"verified"`
		Visibility string `json:"visibility"`
	}

	if err := json.Unmarshal(emailData, &emails); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user email"})
		return
	}

	for _, e := range emails {
		if e.Primary && e.Verified {
			githubUser.Email = e.Email
			break
		}
	}

	if githubUser.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No verified primary email found"})
		return
	}

	user, err := saveOrUpdateGithubUser(githubUser, token, userIDPtr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to save user: " + err.Error()})
		return
	}

	if userIDPtr != nil && *userIDPtr != "" {
		c.JSON(http.StatusOK, gin.H{
			"message": "Link Successful",
			"user": gin.H{
				"id":    user.ID,
				"email": user.Email,
			},
		})
	} else {
		tokenJWT, err := GenerateJWT(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Token generation failed: " + err.Error()})
			return
		}

		frontendURL := os.Getenv("FRONTEND_URL")
		if frontendURL == "" {
			frontendURL = "http://localhost:8081"
		}

		redirectURL := fmt.Sprintf("%s/login?token=%s", frontendURL, tokenJWT)
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
	}
}

func getUserData(token *oauth2.Token) ([]byte, error) {
	ctx := context.Background()

	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
	userDataURL := "https://api.github.com/user"
	req, _ := http.NewRequest("GET", userDataURL, nil)
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print("Request Failed:", err)
		return nil, fmt.Errorf("failed to send request: %s", err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	contents, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, nil
}

func getUserEmail(token *oauth2.Token) ([]byte, error) {
	ctx := context.Background()

	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
	userDataURL := "https://api.github.com/user/emails"
	req, _ := http.NewRequest("GET", userDataURL, nil)
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print("Request Failed:", err)
		return nil, fmt.Errorf("failed to send request: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	contents, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, nil
}

func saveOrUpdateGithubUser(info model.GithubUserInfo, token *oauth2.Token, existingUserID *string) (*db.UsersModel, error) {
	ctx := context.Background()
	serviceName := "Github"

	service, err := initializers.DB.Services.FindUnique(
		db.Services.Name.Equals(serviceName),
	).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("service not found: %s", err.Error())
	}

	if existingUserID != nil && *existingUserID != "" {
		existingUser, err := initializers.DB.Users.FindUnique(
			db.Users.ID.Equals(*existingUserID),
		).Exec(ctx)
		if err != nil {
			return nil, fmt.Errorf("user not found: %s", err.Error())
		}

		existingToken, _ := initializers.DB.UserServiceTokens.FindFirst(
			db.UserServiceTokens.ServiceID.Equals(service.ID),
			db.UserServiceTokens.ProviderID.Equals(strconv.Itoa(info.ID)),
		).With(
			db.UserServiceTokens.User.Fetch(),
		).Exec(ctx)

		if existingToken != nil && existingToken.UserID != *existingUserID {
			return nil, fmt.Errorf("this GitHub account is already linked to another user")
		}

		userGithubToken, _ := initializers.DB.UserServiceTokens.FindFirst(
			db.UserServiceTokens.UserID.Equals(*existingUserID),
			db.UserServiceTokens.ServiceID.Equals(service.ID),
		).Exec(ctx)

		if userGithubToken != nil {
			_, err = initializers.DB.UserServiceTokens.FindUnique(
				db.UserServiceTokens.ID.Equals(userGithubToken.ID),
			).Update(
				db.UserServiceTokens.AccessToken.Set(token.AccessToken),
				db.UserServiceTokens.RefreshToken.Set(token.RefreshToken),
				db.UserServiceTokens.ExpiresAt.Set(token.Expiry),
				db.UserServiceTokens.ProviderID.Set(strconv.Itoa(info.ID)),
			).Exec(ctx)
		} else {
			_, err = initializers.DB.UserServiceTokens.CreateOne(
				db.UserServiceTokens.AccessToken.Set(token.AccessToken),
				db.UserServiceTokens.RefreshToken.Set(token.RefreshToken),
				db.UserServiceTokens.ExpiresAt.Set(token.Expiry),
				db.UserServiceTokens.User.Link(db.Users.ID.Equals(*existingUserID)),
				db.UserServiceTokens.Service.Link(db.Services.ID.Equals(service.ID)),
				db.UserServiceTokens.ProviderID.Set(strconv.Itoa(info.ID)),
			).Exec(ctx)
		}

		if err != nil {
			return nil, err
		}

		return existingUser, nil
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
				db.UserServiceTokens.ProviderID.Set(strconv.Itoa(info.ID)),
			).Exec(ctx)
		} else {
			_, err = initializers.DB.UserServiceTokens.CreateOne(
				db.UserServiceTokens.AccessToken.Set(token.AccessToken),
				db.UserServiceTokens.RefreshToken.Set(token.RefreshToken),
				db.UserServiceTokens.ExpiresAt.Set(token.Expiry),
				db.UserServiceTokens.User.Link(db.Users.ID.Equals(existingUser.ID)),
				db.UserServiceTokens.Service.Link(db.Services.ID.Equals(service.ID)),
				db.UserServiceTokens.ProviderID.Set(strconv.Itoa(info.ID)),
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
		db.UserServiceTokens.ProviderID.Set(strconv.Itoa(info.ID)),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}
