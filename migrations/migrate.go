package migrations

func Migrate() error {
	return AuthMigration()
}