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
	"golang.org/x/oauth2/google"
)

func getGoogleOAuthConfig() *oauth2.Config {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")

	return &oauth2.Config{
		RedirectURL:  redirectURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/gmail.send",
			"https://www.googleapis.com/auth/spreadsheets",
			"https://www.googleapis.com/auth/calendar",
		},
		Endpoint: google.Endpoint,
	}
}

var googleOauthConfig *oauth2.Config = getGoogleOAuthConfig()

func getUserDataFromGoogle(token *oauth2.Token) ([]byte, error) {
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

func saveOrUpdateGoogleUser(info model.GoogleUserInfo, token *oauth2.Token, existingUserID *string) (*db.UsersModel, error) {
	ctx := context.Background()
	serviceName := "Google"

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
			db.UserServiceTokens.ProviderID.Equals(info.ID),
		).With(
			db.UserServiceTokens.User.Fetch(),
		).Exec(ctx)

		if existingToken != nil && existingToken.UserID != *existingUserID {
			return nil, fmt.Errorf("this Google account is already linked to another user")
		}

		userGoogleToken, _ := initializers.DB.UserServiceTokens.FindFirst(
			db.UserServiceTokens.UserID.Equals(*existingUserID),
			db.UserServiceTokens.ServiceID.Equals(service.ID),
		).Exec(ctx)

		if userGoogleToken != nil {
			_, err = initializers.DB.UserServiceTokens.FindUnique(
				db.UserServiceTokens.ID.Equals(userGoogleToken.ID),
			).Update(
				db.UserServiceTokens.AccessToken.Set(token.AccessToken),
				db.UserServiceTokens.RefreshToken.Set(token.RefreshToken),
				db.UserServiceTokens.ExpiresAt.Set(token.Expiry),
				db.UserServiceTokens.ProviderID.Set(info.ID),
			).Exec(ctx)
		} else {
			_, err = initializers.DB.UserServiceTokens.CreateOne(
				db.UserServiceTokens.AccessToken.Set(token.AccessToken),
				db.UserServiceTokens.RefreshToken.Set(token.RefreshToken),
				db.UserServiceTokens.ExpiresAt.Set(token.Expiry),
				db.UserServiceTokens.User.Link(db.Users.ID.Equals(*existingUserID)),
				db.UserServiceTokens.Service.Link(db.Services.ID.Equals(service.ID)),
				db.UserServiceTokens.ProviderID.Set(info.ID),
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
				db.UserServiceTokens.ProviderID.Set(info.ID),
			).Exec(ctx)
		} else {
			_, err = initializers.DB.UserServiceTokens.CreateOne(
				db.UserServiceTokens.AccessToken.Set(token.AccessToken),
				db.UserServiceTokens.RefreshToken.Set(token.RefreshToken),
				db.UserServiceTokens.ExpiresAt.Set(token.Expiry),
				db.UserServiceTokens.User.Link(db.Users.ID.Equals(existingUser.ID)),
				db.UserServiceTokens.Service.Link(db.Services.ID.Equals(service.ID)),
				db.UserServiceTokens.ProviderID.Set(info.ID),
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
		db.UserServiceTokens.ProviderID.Set(info.ID),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func GoogleLogin(ctx *gin.Context) {
	oauthState := generateStateOauthCookie()
	userID, exist := ctx.Get("userID")

	if exist {
		ctx.SetCookie("oauth_user_id", userID.(string), 300, "/", "", false, true)
	}

	redirectURI := ctx.Query("redirect_uri")
	if redirectURI != "" {
		ctx.SetCookie("oauth_redirect_uri", redirectURI, 300, "/", "", false, true)
	}

	ctx.SetCookie("oauthState", oauthState, 3600, "/", "", false, true)
	url := googleOauthConfig.AuthCodeURL(oauthState, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(ctx *gin.Context) {
	state := ctx.Query("state")
	code := ctx.Query("code")

	currentUserID, err := ctx.Cookie("oauth_user_id")
	var userIDPtr *string
	if err == nil && currentUserID != "" {
		userIDPtr = &currentUserID
	}
	ctx.SetCookie("oauth_user_id", "", -1, "/", "", false, true)

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

	data, err := getUserDataFromGoogle(token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	var googleUser model.GoogleUserInfo
	if err := json.Unmarshal(data, &googleUser); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to parse user data"})
		return
	}

	user, err := saveOrUpdateGoogleUser(googleUser, token, userIDPtr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to save user: " + err.Error()})
		return
	}

	customRedirectURI, err := ctx.Cookie("oauth_redirect_uri")
	if err == nil && customRedirectURI != "" {
		ctx.SetCookie("oauth_redirect_uri", "", -1, "/", "", false, true)
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:8081"
	}

	if userIDPtr != nil && *userIDPtr != "" {
		if customRedirectURI != "" {
			ctx.Redirect(http.StatusTemporaryRedirect, customRedirectURI)
		} else {
			redirectURL := fmt.Sprintf("%s/dashboard", frontendURL)
			ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)
		}
	} else {
		tokenJWT, err := GenerateJWT(user.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"Error": "Token generation failed: " + err.Error()})
			return
		}

		if customRedirectURI != "" {
			redirectURL := fmt.Sprintf("%s?token=%s", customRedirectURI, tokenJWT)
			ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)
		} else {
			redirectURL := fmt.Sprintf("%s/login?token=%s", frontendURL, tokenJWT)
			ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)
		}
	}
}
