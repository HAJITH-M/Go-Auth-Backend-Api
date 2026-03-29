package authservice

import (
	"errors"
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
		return errors.New("old password doesn't match")
	}
	newPassMatch := utils.CompareHashedPassword(storedPass, userNewPass)
	if newPassMatch == nil {
		return errors.New("New password should be different from old password")
	}

	if err := repository.ChangePasswordRepo(userEmail, userNewPassHash); err != nil {
		return errors.New("Failed to update Password")
	}

	return nil
}

func ForgotPasswordUpdateService(email, newPassword string) error {
	authPass, err := repository.GetUserPasswordRepo(email)
	if err != nil {
		return errors.New("failed to get password")
	}

	err = utils.CompareHashedPassword(authPass, newPassword)
	if err == nil {
		return errors.New("new password should not be same as old password")
	}

	hashPassword, err := utils.GeneratePasswordWithHash(newPassword)
	if err != nil {
		return err
	}

	err = repository.ChangePasswordRepo(email, hashPassword)
	if err != nil {
		return errors.New("failed to update password")
	}

	err = repository.UpdateUserAccountStatusRepo(email, "active")
	if err != nil {
		return errors.New("failed to update status")
	}

	return nil
}
