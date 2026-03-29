package authservice

import (
	"errors"
	authmodel "go-auth-backend-api/internal/model/AuthModel"
	"go-auth-backend-api/internal/repository"
	"go-auth-backend-api/pkg/mailer"
	"go-auth-backend-api/pkg/utils"
	"log"
	"time"

	"github.com/google/uuid"
)

func RegisterService(input RegisterInput) (*RegisterResult, error) {
	existing, err := repository.GetUserByEmailRepo(input.Email)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		return nil, errors.New("Email already Exists")
	}

	hashPassword, err := utils.GeneratePasswordWithHash(input.Password)
	if err != nil {
		return nil, errors.New("Failed to process Password")
	}

	userID := uuid.NewString()
	now := time.Now()

	user := &authmodel.User{
		ID:            userID,
		Email:         input.Email,
		DisplayName:   input.DisplayName,
		EmailVerified: false,
		AccountStatus: "active",
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	method := &authmodel.AuthenticationMethod{
		ID:           uuid.NewString(),
		UserID:       userID,
		ProviderType: "password",
		PasswordHash: string(hashPassword),
		IsPrimary:    true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err = repository.CreateUserWithAuthRepo(user, method)
	if err != nil {
		return nil, errors.New("failed to create account")
	}

	rawToken, err := utils.GenerateSecureToken(userID)
	if err != nil {
		return nil, errors.New("failed to generate verification token")
	}

	rawTokenHash := utils.HashSecureToken(rawToken)
	tokenID := uuid.NewString()

	userToken := &authmodel.UserToken{
		TokenID:   tokenID,
		UserID:    userID,
		TokenHash: rawTokenHash,
		TokenType: "email_verification",
		ExpiresAt: time.Now().Add(10 * time.Minute),
		IsUsed:    false,
		CreatedAt: now,
	}

	if err := repository.UserTokenCreationRepo(userToken); err != nil {
		return nil, errors.New("Failed to save Token")
	}

	go func() {
		if err := mailer.SendEmailVerificationEmail(user.Email, user.DisplayName, rawToken); err != nil {
			log.Println("Failed to send verification email:", err)
		}
	}()

	return &RegisterResult{
		Email:       user.Email,
		DisplayName: user.DisplayName,
		UserId:      user.ID,
	}, nil
}
