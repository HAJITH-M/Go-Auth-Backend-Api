package mailconfig

import "go-auth-backend-api/internal/config/env"

type SmtpConfig struct {
	From     string
	Password string
	Host     string
	Port     string
	BaseURL  string
}

// load config from env
func LoadSMTPConfig() SmtpConfig {
	return SmtpConfig{
		From:     env.AppEnv.SMTP_FROM,
		Password: env.AppEnv.SMTP_PASSWORD,
		Host:     env.AppEnv.SMTP_HOST,
		Port:     env.AppEnv.SMTP_PORT,
		BaseURL:  env.AppEnv.SMTP_URL,
	}
}
