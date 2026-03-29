package repository

import (
	"errors"
	authmodel "go-auth-backend-api/internal/model/AuthModel"
	"go-auth-backend-api/pkg/database"
	"time"

	"gorm.io/gorm"
)

func GetUserByEmailRepo(email string) (*authmodel.User, error) {

	var user authmodel.User

	err := database.DB.Where("email= ?", email).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, err
}

func GetUserByIDRepo(uid string) (*authmodel.User, error) {
	var user authmodel.User
	err := database.DB.Where("id = ?", uid).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, nil
}

func CreateUserWithAuthRepo(user *authmodel.User, method *authmodel.AuthenticationMethod) error {

	// using transaction if one step fails also it will revert to previous steps
	err := database.DB.Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(user).Error; err != nil {
			return err
		}

		if err := tx.Create(method).Error; err != nil {
			return err
		}

		return nil // commit
	})

	return err
}

func GetAuthenticationMethodUserRepo(email string) (*authmodel.User, error) {

	var user authmodel.User

	err := database.DB.Preload("AuthMethods").Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, err
}

func LastLoginUserUpdateRepo(userID string) error {
	now := time.Now()
	return database.DB.
		Model(&authmodel.User{}).
		Where("id = ?", userID).
		Update("last_login_at", now).Error
}

func CreateSessionUserRepo(session *authmodel.Session) error {
	return database.DB.Create(session).Error
}

func GetSessionByTokenHashRepo(hash string) (*authmodel.Session, error) {
	var session authmodel.Session

	err := database.DB.
		Where("refresh_token_hash = ? AND is_revoked = false and expires_at > NOW()", hash).
		First(&session).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &session, err
}

func RevokeSessionRepo(sessionID string) error {
	return database.DB.
		Model(&authmodel.Session{}).
		Where("session_id = ?", sessionID).
		Update("is_revoked", true).Error
}

func UserTokenCreationRepo(userToken *authmodel.UserToken) error {
	return database.DB.Create(userToken).Error
}

func GetUserTokenRepo(hashedToken, tokenType string) (*authmodel.UserToken, error) {
	var userToken authmodel.UserToken
	err := database.DB.
		Where("token_hash = ? AND token_type = ? AND is_used = false AND expires_at > NOW()", hashedToken, tokenType).
		First(&userToken).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &userToken, nil
}

func UpdateVerificationUserTokenStatusRepo(userID, tokenId string) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&authmodel.UserToken{}).Where("token_id = ?", tokenId).Update("is_used", true).Error; err != nil {
			return errors.New("failed to mark token as used")
		}

		if err := tx.Model(&authmodel.User{}).Where("id = ?", userID).Update("email_verified", true).Error; err != nil {
			return errors.New("failed to verify user email")
		}
		return nil
	})
}

func GetUserPasswordRepo(email string) (string, error) {

	var authMethod authmodel.AuthenticationMethod

	err := database.DB.
		Joins("JOIN users ON users.id = authentication_methods.user_id").
		Where("users.email = ?", email).
		Select("authentication_methods.password_hash").
		First(&authMethod).Error

	if err != nil {
		return "", err
	}

	return authMethod.PasswordHash, nil
}

func ChangePasswordRepo(email, newPassword string) error {
	// 1. Find the user by email
	var user authmodel.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// 2. Update password for this user's password auth method
	return database.DB.
		Model(&authmodel.AuthenticationMethod{}).
		Where("user_id = ? AND provider_type = ?", user.ID, "password").
		Update("password_hash", newPassword).Error
}

func UpdateUserAccountStatusRepo(email, status string) error {
	// 1. Find the user by email
	var user authmodel.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}
	return database.DB.
		Model(&authmodel.User{}).
		Where("id = ?", user.ID).
		Update("account_status", status).Error
}

// GetAuthenticationMethodByProviderRepo fetches an authentication method using the provider identity.
func GetAuthenticationMethodByProviderRepo(providerType, providerUserID string) (*authmodel.AuthenticationMethod, error) {
	var method authmodel.AuthenticationMethod
	err := database.DB.
		Where("provider_type = ? AND provider_user_id = ?", providerType, providerUserID).
		First(&method).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &method, err
}

func CreateAuthenticationMethodRepo(method *authmodel.AuthenticationMethod) error {
	return database.DB.Create(method).Error
}

// UpdateAuthenticationMethodOauthTokensByMethodIDRepo updates stored Google OAuth refresh token + expiry.
// If oauthRefreshToken is empty, it will not overwrite the existing refresh token.
func UpdateAuthenticationMethodOauthTokensByMethodIDRepo(methodID string, oauthRefreshToken string, oauthTokenExpiry *time.Time) error {
	updates := map[string]interface{}{}

	if oauthRefreshToken != "" {
		updates["oauth_refresh_token"] = oauthRefreshToken
	}
	if oauthTokenExpiry != nil {
		updates["oauth_token_expiry"] = oauthTokenExpiry
	}

	if len(updates) == 0 {
		return nil
	}

	return database.DB.
		Model(&authmodel.AuthenticationMethod{}).
		Where("id = ?", methodID).
		Updates(updates).Error
}

// UpdateUserForOAuthRepo marks email as verified and updates display picture/name.
func UpdateUserForOAuthRepo(userID, displayName, profilePictureURL string) error {
	updates := map[string]interface{}{
		"email_verified": true,
	}
	if displayName != "" {
		updates["display_name"] = displayName
	}
	if profilePictureURL != "" {
		updates["profile_picture_url"] = profilePictureURL
	}

	return database.DB.
		Model(&authmodel.User{}).
		Where("id = ?", userID).
		Updates(updates).Error
}

// GetAuthenticationMethodByUserAndProviderRepo finds the first auth method for a user/provider.
func GetAuthenticationMethodByUserAndProviderRepo(userID, providerType string) (*authmodel.AuthenticationMethod, error) {
	var method authmodel.AuthenticationMethod
	err := database.DB.
		Where("user_id = ? AND provider_type = ?", userID, providerType).
		First(&method).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &method, err
}
