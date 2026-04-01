package oauthhandler

import (
	"context"
	autherrors "go-auth-backend-api/internal/errors"
	"go-auth-backend-api/internal/config/authconfig"
	"go-auth-backend-api/internal/config/env"
	"go-auth-backend-api/internal/repository"
	tokenjwt "go-auth-backend-api/pkg/tokenJWT"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func RefreshAccessToken(c *gin.Context) {
	rawRefreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token missing"})
		return
	}

	claims, err := tokenjwt.ValidateToken(rawRefreshToken, env.AppEnv.JWT_REFRESH_SECRET)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	method, err := repository.GetAuthenticationMethodByUserAndProviderRepo(claims.UserID, "google")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": autherrors.ErrOAuthFailedToLoadCredentials.Error()})
		return
	}
	if method == nil || method.OauthRefreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "google refresh token missing"})
		return
	}

	token := &oauth2.Token{RefreshToken: method.OauthRefreshToken}
	newToken, err := authconfig.GoogleOAuthConfig().TokenSource(context.Background(), token).Token()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": autherrors.ErrOAuthFailedToRefreshGoogleToken.Error()})
		return
	}

	var expiryPtr *time.Time
	if !newToken.Expiry.IsZero() {
		exp := newToken.Expiry
		expiryPtr = &exp
	}

	if err := repository.UpdateAuthenticationMethodOauthTokensByMethodIDRepo(method.ID, newToken.RefreshToken, expiryPtr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": autherrors.ErrOAuthFailedToPersistRefreshToken.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": newToken.AccessToken,
		"expiry":       newToken.Expiry,
	})
}
