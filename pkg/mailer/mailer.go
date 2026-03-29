package mailer

import (
	"fmt"
	mailconfig "go-auth-backend-api/internal/config/mailConfig"
	"net/smtp"
)

func SendEmailVerificationEmail(toEmail, displayName, token string) error {

	cfg := mailconfig.LoadSMTPConfig()

	link := fmt.Sprintf("%s/api/v1/auth/verification-email?token=%s", cfg.BaseURL, token)

	subject := "Subject: Verify your email\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	body := fmt.Sprintf(`
		<h2>Hi %s,</h2>
		<p>Click the link below to verify your email:</p>
		<a href="%s">Verify Email</a>
		<p>This link expires in 24 hours.</p>
	`, displayName, link)

	msg := []byte(subject + mime + body)

	auth := smtp.PlainAuth("", cfg.From, cfg.Password, cfg.Host)
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	return smtp.SendMail(addr, auth, cfg.From, []string{toEmail}, msg)
}

func SendOtpEmail(toEmail, otp string) error {

	cfg := mailconfig.LoadSMTPConfig()

	subject := "Subject: Your Password Reset OTP\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	body := fmt.Sprintf(`
		<h2>Hello %s,</h2>
		<p>Your password reset OTP is:</p>
		<h1 style="letter-spacing:3px;">%s</h1>
		<p>This OTP will expire in 5 minutes.</p>
		<p>If you didn't request this, please ignore this email.</p>
	`, toEmail, otp)

	msg := []byte(subject + mime + body)

	auth := smtp.PlainAuth("", cfg.From, cfg.Password, cfg.Host)
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	return smtp.SendMail(addr, auth, cfg.From, []string{toEmail}, msg)
}
