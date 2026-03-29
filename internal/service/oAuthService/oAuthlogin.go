package oauthservice

import (
	"errors"
	authmodel "go-auth-backend-api/internal/model/AuthModel"
	"go-auth-backend-api/internal/repository"
	authservice "go-auth-backend-api/internal/service/authService"
	tokenjwt "go-auth-backend-api/pkg/tokenJWT"
	"go-auth-backend-api/pkg/utils"
	"time"

	"github.com/google/uuid"
)

// ErrLinkRequired is returned when a Google login attempt finds an existing password account.
// The frontend should prompt the user to confirm linking via password.
var ErrLinkRequired = errors.New("link_required")

func GoogleOAuthLoginService(input GoogleUser, oauthRefreshToken string, oauthTokenExpiry *time.Time, ipAddress string, device string) (*authservice.LoginResult, error) {
	const providerType = "google"

	// 1) Try to locate the user by provider user id (strong match).
	method, err := repository.GetAuthenticationMethodByProviderRepo(providerType, input.ID)
	if err != nil {
		return nil, err
	}

	var user *authmodel.User

	if method != nil {
		// Existing account linked to this Google identity.
		user, err = repository.GetUserByIDRepo(method.UserID)
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, errors.New("invalid credentials")
		}

		// Update OAuth tokens (refresh token can be empty if Google doesn't return it this time).
		if oauthRefreshToken != "" || oauthTokenExpiry != nil {
			if err := repository.UpdateAuthenticationMethodOauthTokensByMethodIDRepo(method.ID, oauthRefreshToken, oauthTokenExpiry); err != nil {
				return nil, err
			}
		}

		if err := repository.UpdateUserForOAuthRepo(user.ID, input.Name, input.Picture); err != nil {
			return nil, err
		}
		user.EmailVerified = true
	} else {
		// No provider link yet — find by email or create new user.
		userByEmail, err := repository.GetUserByEmailRepo(input.Email)
		if err != nil {
			return nil, err
		}

		if userByEmail == nil {
			// Brand new user — create account with Google provider.
			now := time.Now()
			newUserID := uuid.NewString()
			user = &authmodel.User{
				ID:                newUserID,
				Email:             input.Email,
				DisplayName:       input.Name,
				ProfilePictureURL: input.Picture,
				EmailVerified:     true,
				AccountStatus:     "active",
				CreatedAt:         now,
				UpdatedAt:         now,
			}

			newMethod := &authmodel.AuthenticationMethod{
				ID:                uuid.NewString(),
				UserID:            newUserID,
				ProviderType:      providerType,
				ProviderUserID:    input.ID,
				OauthRefreshToken: oauthRefreshToken,
				OauthTokenExpiry:  oauthTokenExpiry,
				IsPrimary:         true,
				CreatedAt:         now,
				UpdatedAt:         now,
			}

			if err := repository.CreateUserWithAuthRepo(user, newMethod); err != nil {
				return nil, errors.New("failed to create account")
			}
		} else {
			user = userByEmail

			// Check if this user has a password provider.
			existingPasswordMethod, err := repository.GetAuthenticationMethodByUserAndProviderRepo(user.ID, "password")
			if err != nil {
				return nil, err
			}

			if existingPasswordMethod != nil {
				// Password account exists — signal frontend to show link confirmation.
				// We do NOT link silently. User must confirm with their password.
				return nil, ErrLinkRequired
			}

			// No password provider — safe to link Google directly.
			newMethod := &authmodel.AuthenticationMethod{
				ID:                uuid.NewString(),
				UserID:            user.ID,
				ProviderType:      providerType,
				ProviderUserID:    input.ID,
				OauthRefreshToken: oauthRefreshToken,
				OauthTokenExpiry:  oauthTokenExpiry,
				IsPrimary:         false,
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			}

			if err := repository.CreateAuthenticationMethodRepo(newMethod); err != nil {
				return nil, errors.New("failed to link oauth account")
			}

			if err := repository.UpdateUserForOAuthRepo(user.ID, input.Name, input.Picture); err != nil {
				return nil, err
			}
			user.EmailVerified = true
		}
	}

	// 2) Enforce your existing login rules (same error messages as password login).
	if user.EmailVerified != true {
		return nil, errors.New("Email Not Verified")
	}
	if user.AccountStatus == "pending_otp" {
		return nil, errors.New("Change password Required")
	}
	if user.AccountStatus != "active" {
		return nil, errors.New("account is not active")
	}

	// 3) Create app JWT tokens + DB session (same behavior as LoginService).
	accessToken, err := tokenjwt.GenerateAccessTokenJWT(user.ID, user.Email)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}
	refreshToken, err := tokenjwt.GenerateRefreshTokenJWT(user.ID, user.Email)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	session := &authmodel.Session{
		SessionID:        uuid.NewString(),
		UserID:           user.ID,
		RefreshTokenHash: utils.HashSecureToken(refreshToken),
		IPAddress:        ipAddress,
		DeviceInfo:       device,
		ExpiresAt:        time.Now().Add(7 * 24 * time.Hour),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := repository.CreateSessionUserRepo(session); err != nil {
		return nil, errors.New("failed to create session")
	}
	go repository.LastLoginUserUpdateRepo(user.ID)

	return &authservice.LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserId:       user.ID,
		Email:        user.Email,
		DisplayName:  user.DisplayName,
	}, nil
}
