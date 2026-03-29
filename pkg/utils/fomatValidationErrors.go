package utils

import "github.com/go-playground/validator/v10"

func FormatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {

			switch fe.Field() {

			case "Email":
				errors["email"] = "invalid email format"

			case "Password":
				errors["password"] = "password must be at least 8 characters"

			case "DisplayName":
				errors["display_name"] = "display name is required"

			default:
				errors[fe.Field()] = fe.Error()
			}
		}
	}

	return errors
}
