package authentification

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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

func GithubLogin(context *gin.Context) {
	state := "wefioweiofgeoprgerogoperojpg"
	url := githubOauthConfig.AuthCodeURL(state)
	context.SetCookie("oauth_state", state, 300, "/", "", false, true)
	context.Redirect(302, url)
}

func GithubCallback(context *gin.Context) {
	code := context.Query("code")
	state := context.Query("state")

	randomState, err := context.Cookie("oauth_state")
	if err != nil || state != randomState {
		context.JSON(http.StatusNotFound, gin.H{"error": "Invalid State"})
	}
	context.SetCookie("oauth_state", "", -1, "/", "", false, true)
	token, err := githubOauthConfig.Exchange(context, code)
	if err != nil {
		context.JSON(500, gin.H{"error": "OAuth exchange failed"})
		return
	}
	getUserData(token)
	
}

func getUserData(token *oauth2.Token) {
	ctx := context.Background()

	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
	userDataURL := "https://api.github.com/user"
	req, _ := http.NewRequest("GET", userDataURL, nil)
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print("Request Failed:", err)
		return;
	}
    defer resp.Body.Close()

	var user map[string]any
    json.NewDecoder(resp.Body).Decode(&user)
    fmt.Println("User info:", user)
}
