package mailer

import (
	"fmt"
	"net/smtp"
)

func SendPasswordResetEmail(to, token string) error {
	from := "noreply@example.com"
	password := "password"
	smtpServer := "smtp.example.com"
	smtpPort := "587"

	body := fmt.Sprintf("To reset your password, please click the link below:\nhttp://localhost:8080/reset-password/%s", token)

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Password Reset\n\n" +
		body

	auth := smtp.PlainAuth("", from, password, smtpServer)

	err := smtp.SendMail(smtpServer+":"+smtpPort, auth, from, []string{to}, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
