package authservice

import (
	"go-auth-backend-api/internal/config/env"
	autherrors "go-auth-backend-api/internal/errors"
	authmodel "go-auth-backend-api/internal/model/AuthModel"
	"go-auth-backend-api/internal/repository"
	tokenjwt "go-auth-backend-api/pkg/tokenJWT"
	"go-auth-backend-api/pkg/utils"
	"log"
	"time"

	"github.com/google/uuid"
)

func LoginService(input LoginInput) (*LoginResult, error) {
	user, err := repository.GetAuthenticationMethodUserRepo(input.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, autherrors.ErrInvalidCredentials
	}

	var passwordMethod *authmodel.AuthenticationMethod
	for i := range user.AuthMethods {
		if user.AuthMethods[i].ProviderType == "password" {
			passwordMethod = &user.AuthMethods[i]
			break
		}
	}

	if passwordMethod == nil {
		return nil, autherrors.ErrInvalidCredentials
	}

	err = utils.CompareHashedPassword(passwordMethod.PasswordHash, input.Password)
	if err != nil {
		return nil, autherrors.ErrInvalidCredentials
	}

	if user.EmailVerified != true {
		return nil, autherrors.ErrEmailNotVerified
	}

	if user.AccountStatus == "pending_otp" {
		return nil, autherrors.ErrChangePasswordRequired
	}

	if user.AccountStatus != "active" {
		return nil, autherrors.ErrAccountNotActive
	}

	accessToken, err := generateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, autherrors.ErrFailedToGenerateToken
	}
	refreshToken, err := generateRefreshToken(user.ID, user.Email)
	if err != nil {
		return nil, autherrors.ErrFailedToGenerateToken
	}

	session := &authmodel.Session{
		SessionID:        uuid.NewString(),
		UserID:           user.ID,
		RefreshTokenHash: hashToken(refreshToken),
		IPAddress:        input.IPAddress,
		DeviceInfo:       input.Device,
		ExpiresAt:        time.Now().Add(7 * 24 * time.Hour),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := repository.CreateSessionUserRepo(session); err != nil {
		return nil, autherrors.ErrFailedToCreateSession
	}

	go repository.LastLoginUserUpdateRepo(user.ID)

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserId:       user.ID,
		Email:        user.Email,
		DisplayName:  user.DisplayName,
	}, nil
}

func RefreshTokenService(rawRefreshToken string) (*LoginResult, error) {
	claims, err := tokenjwt.ValidateToken(rawRefreshToken, env.AppEnv.JWT_REFRESH_SECRET)
	if err != nil {
		log.Println("Refresh token invalid:", err)
		return nil, autherrors.ErrInvalidRefreshToken
	}

	hash := hashToken(rawRefreshToken)
	session, err := repository.GetSessionByTokenHashRepo(hash)
	if err != nil {
		return nil, err
	}
	if session == nil {
		log.Println("Session not found or revoked")
		return nil, autherrors.ErrSessionNotFoundOrRevoked
	}

	if err := repository.RevokeSessionRepo(session.SessionID); err != nil {
		return nil, autherrors.ErrFailedToRevokeSession
	}

	user, err := repository.GetUserByIDRepo(claims.UserID)
	if err != nil || user == nil {
		return nil, autherrors.ErrUserNotFound
	}

	accessToken, err := generateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, autherrors.ErrFailedToGenerateToken
	}
	refreshToken, err := generateRefreshToken(user.ID, user.Email)
	if err != nil {
		return nil, autherrors.ErrFailedToGenerateToken
	}

	session = &authmodel.Session{
		SessionID:        uuid.NewString(),
		UserID:           user.ID,
		RefreshTokenHash: hashToken(refreshToken),
		IPAddress:        session.IPAddress,
		DeviceInfo:       session.DeviceInfo,
		ExpiresAt:        time.Now().Add(7 * 24 * time.Hour),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	if err := repository.CreateSessionUserRepo(session); err != nil {
		return nil, autherrors.ErrFailedToCreateSession
	}

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserId:       user.ID,
		Email:        user.Email,
		DisplayName:  user.DisplayName,
	}, nil
}

func LogoutService(rawRefreshToken string) error {
	hash := hashToken(rawRefreshToken)
	session, err := repository.GetSessionByTokenHashRepo(hash)
	if err != nil || session == nil {
		return nil
	}
	return repository.RevokeSessionRepo(session.SessionID)
}

func UpdateVerificationEmailService(rawToken string) error {
	tokenHash := utils.HashSecureToken(rawToken)

	userToken, err := repository.GetUserTokenRepo(tokenHash, "email_verification")
	if err != nil {
		return err
	}
	if userToken == nil {
		return autherrors.ErrUnableToGetUserToken
	}

	if err := repository.UpdateVerificationUserTokenStatusRepo(userToken.UserID, userToken.TokenID); err != nil {
		return err
	}

	return nil
}
