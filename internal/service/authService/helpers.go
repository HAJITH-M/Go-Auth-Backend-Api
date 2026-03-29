package authservice

import (
	tokenjwt "go-auth-backend-api/pkg/tokenJWT"
	"go-auth-backend-api/pkg/utils"
)

func generateAccessToken(userID, email string) (string, error) {
	return tokenjwt.GenerateAccessTokenJWT(userID, email)
}

func generateRefreshToken(userID, email string) (string, error) {
	return tokenjwt.GenerateRefreshTokenJWT(userID, email)
}

func hashToken(raw string) string {
	return utils.HashSecureToken(raw)
}
