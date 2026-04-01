package autherrors

import "errors"

// Login and session
var (
	ErrInvalidCredentials       = errors.New("invalid credentials")
	ErrEmailNotVerified         = errors.New("Email Not Verified")
	ErrChangePasswordRequired   = errors.New("Change password Required")
	ErrAccountNotActive         = errors.New("account is not active")
	ErrFailedToGenerateToken    = errors.New("failed to generate token")
	ErrFailedToCreateSession    = errors.New("failed to create session")
	ErrInvalidRefreshToken      = errors.New("invalid refresh token")
	ErrSessionNotFoundOrRevoked = errors.New("session not found or revoked")
	ErrFailedToRevokeSession    = errors.New("failed to revoke session")
	ErrUserNotFound             = errors.New("user not found")
	ErrUnableToGetUserToken     = errors.New("unable to get user token")
)

// Registration
var (
	ErrEmailAlreadyExists                = errors.New("Email already Exists")
	ErrFailedToProcessPassword           = errors.New("Failed to process Password")
	ErrFailedToCreateAccount             = errors.New("failed to create account")
	ErrFailedToGenerateVerificationToken = errors.New("failed to generate verification token")
	ErrFailedToSaveToken                 = errors.New("Failed to save Token")
	ErrFailedToSendVerificationEmail   = errors.New("failed to send verification email")
)

// Password
var (
	ErrOldPasswordDoesntMatch         = errors.New("old password doesn't match")
	ErrNewPasswordShouldDifferFromOld = errors.New("New password should be different from old password")
	ErrFailedToUpdatePasswordChange   = errors.New("Failed to update Password")
	ErrFailedToGetPassword            = errors.New("failed to get password")
	ErrNewPasswordSameAsCurrent       = errors.New("new password should not be same as old password")
	ErrFailedToUpdatePassword         = errors.New("failed to update password")
	ErrFailedToUpdateStatus           = errors.New("failed to update status")
)

// Forgot password / OTP
var (
	ErrOtpCooldown          = errors.New("please wait 30 seconds before requesting another OTP")
	ErrFailedToStoreOTP     = errors.New("Failed to store OTP")
	ErrTooManyOtpAttempts   = errors.New("too many attempts, please request a new OTP")
	ErrInvalidOTP           = errors.New("Invalid OTP")
	ErrFailedToSendOTPEmail = errors.New("failed to send OTP email")
)

// OAuth login
// ErrLinkRequired is returned when Google login finds an existing password account; the client should confirm linking.
var (
	ErrLinkRequired             = errors.New("link_required")
	ErrFailedToLinkOAuthAccount = errors.New("failed to link oauth account")
)

// Rate limiting
var (
	ErrRateLimiterRedis  = errors.New("failed to generate token")
	ErrRateLimitExceeded = errors.New("Rate Limit exceeded try after 1 minute")
)

// OAuth HTTP handler responses
var (
	ErrOAuthTokenExchangeFailed         = errors.New("token exchange failed")
	ErrOAuthFailedToGetUserInfo         = errors.New("failed to get user info")
	ErrOAuthDecodeFailed                = errors.New("decode failed")
	ErrOAuthGoogleUserInfoIncomplete    = errors.New("google user info missing required fields")
	ErrOAuthFailedToLoadCredentials     = errors.New("failed to load oauth credentials")
	ErrOAuthFailedToRefreshGoogleToken  = errors.New("failed to refresh token")
	ErrOAuthFailedToPersistRefreshToken = errors.New("failed to persist oauth refresh token")
)
