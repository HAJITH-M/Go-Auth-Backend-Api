package authservice

import (
	autherrors "go-auth-backend-api/internal/errors"
	"go-auth-backend-api/internal/repository"
	"go-auth-backend-api/pkg/utils"
)

func ChangePasswordService(input ChangePasswordInput) error {
	userEmail := input.Email
	userOldPass := input.OldPassword
	userNewPass := input.NewPassword
	userNewPassHash, err := utils.GeneratePasswordWithHash(userNewPass)

	if err != nil {
		return err
	}

	storedPass, err := repository.GetUserPasswordRepo(userEmail)
	if err != nil {
		return err
	}

	passMatch := utils.CompareHashedPassword(storedPass, userOldPass)
	if passMatch != nil {
		return autherrors.ErrOldPasswordDoesntMatch
	}
	newPassMatch := utils.CompareHashedPassword(storedPass, userNewPass)
	if newPassMatch == nil {
		return autherrors.ErrNewPasswordShouldDifferFromOld
	}

	if err := repository.ChangePasswordRepo(userEmail, userNewPassHash); err != nil {
		return autherrors.ErrFailedToUpdatePasswordChange
	}

	return nil
}

func ForgotPasswordUpdateService(email, newPassword string) error {
	authPass, err := repository.GetUserPasswordRepo(email)
	if err != nil {
		return autherrors.ErrFailedToGetPassword
	}

	err = utils.CompareHashedPassword(authPass, newPassword)
	if err == nil {
		return autherrors.ErrNewPasswordSameAsCurrent
	}

	hashPassword, err := utils.GeneratePasswordWithHash(newPassword)
	if err != nil {
		return err
	}

	err = repository.ChangePasswordRepo(email, hashPassword)
	if err != nil {
		return autherrors.ErrFailedToUpdatePassword
	}

	err = repository.UpdateUserAccountStatusRepo(email, "active")
	if err != nil {
		return autherrors.ErrFailedToUpdateStatus
	}

	return nil
}
