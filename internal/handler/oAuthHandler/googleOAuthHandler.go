package oauthhandler

import (
	"context"
	"encoding/json"
	"errors"
	autherrors "go-auth-backend-api/internal/errors"
	"go-auth-backend-api/internal/config/authconfig"
	oauthservice "go-auth-backend-api/internal/service/oAuthService"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func GoogleLogin(c *gin.Context) {
	url := authconfig.GoogleOAuthConfig().AuthCodeURL(
		"state-token",
		oauth2.AccessTypeOffline,
		oauth2.ApprovalForce,
	)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(c *gin.Context) {
	if oauthErr := c.Query("error"); oauthErr != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             oauthErr,
			"error_description": c.Query("error_description"),
		})
		return
	}

	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code missing"})
		return
	}

	token, err := authconfig.GoogleOAuthConfig().Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token exchange failed"})
		return
	}

	client := authconfig.GoogleOAuthConfig().Client(context.Background(), token)
	userResp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": autherrors.ErrOAuthFailedToGetUserInfo.Error()})
		return
	}
	defer userResp.Body.Close()

	var user map[string]interface{}
	if err := json.NewDecoder(userResp.Body).Decode(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": autherrors.ErrOAuthDecodeFailed.Error()})
		return
	}

	googleID, _ := user["id"].(string)
	email, _ := user["email"].(string)
	name, _ := user["name"].(string)
	picture, _ := user["picture"].(string)
	if googleID == "" || email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "google user info missing required fields"})
		return
	}

	var oauthTokenExpiry *time.Time
	if !token.Expiry.IsZero() {
		expiry := token.Expiry
		oauthTokenExpiry = &expiry
	}

	loginRes, err := oauthservice.GoogleOAuthLoginService(
		oauthservice.GoogleUser{
			ID:      googleID,
			Email:   email,
			Name:    name,
			Picture: picture,
		},
		token.RefreshToken,
		oauthTokenExpiry,
		c.ClientIP(),
		c.Request.UserAgent(),
	)
	if err != nil {
		if errors.Is(err, autherrors.ErrEmailNotVerified) ||
			errors.Is(err, autherrors.ErrAccountNotActive) ||
			errors.Is(err, autherrors.ErrChangePasswordRequired) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": autherrors.ErrInvalidCredentials.Error()})
		return
	}

	c.SetCookie(
		"refresh_token",
		loginRes.RefreshToken,
		7*24*3600,
		"/api/v1/auth/refresh",
		"",
		false,
		true,
	)

	redirectURL := "http://127.0.0.1:5500/index.html?provider=google" +
		"&user_id=" + url.QueryEscape(loginRes.UserId) +
		"&email=" + url.QueryEscape(loginRes.Email) +
		"&display_name=" + url.QueryEscape(loginRes.DisplayName) +
		"&access_token=" + url.QueryEscape(loginRes.AccessToken)
	c.Redirect(http.StatusSeeOther, redirectURL)
}
