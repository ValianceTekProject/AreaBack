package authentification

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8000/auth/google/callback",
	ClientID:     "{PATTERN}.apps.googleusercontent.com",
	ClientSecret: "{SECRET}",
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/photoslibrary.readonly",
	},
	Endpoint: google.Endpoint,
}

func GoogleLogin(context *gin.Context) {
}

func GoogleCallback(context *gin.Context) {
}
