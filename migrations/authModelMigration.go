package migrations

import (
	authmodel "go-auth-backend-api/internal/model/AuthModel"
	"go-auth-backend-api/pkg/database"
)

func AuthMigration() error {
	return database.DB.AutoMigrate(
		&authmodel.User{},
		&authmodel.AuthenticationMethod{},
		&authmodel.Session{},
		&authmodel.UserToken{},
	)
}
