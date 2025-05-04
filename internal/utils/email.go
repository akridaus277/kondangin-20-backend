package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"path/filepath"
)

func SendVerificationEmail(to, subject, verificationLink string) error {
	from := os.Getenv("SMTP_SENDER_EMAIL")
	senderName := os.Getenv("SMTP_SENDER_NAME")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	// Parse template
	tmplPath := filepath.Join("email_templates", "verification_email.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	var bodyBuffer bytes.Buffer
	if err := tmpl.Execute(&bodyBuffer, map[string]string{
		"VerificationLink": verificationLink,
	}); err != nil {
		return fmt.Errorf("failed to render email template: %w", err)
	}

	// Format full message
	message := []byte(fmt.Sprintf("From: %s <%s>\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		senderName, from, to, subject, bodyBuffer.String()))

	auth := smtp.PlainAuth("", username, password, smtpHost)
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	return smtp.SendMail(addr, auth, from, []string{to}, message)
}

func SendResetPasswordEmail(to, subject, token string) error {
	from := os.Getenv("SMTP_SENDER_EMAIL")
	senderName := os.Getenv("SMTP_SENDER_NAME")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	resetLink := os.Getenv("APP_URL") + "/reset-password?token=" + token

	// Parse template
	tmplPath := filepath.Join("email_templates", "reset_password_email.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	var bodyBuffer bytes.Buffer
	if err := tmpl.Execute(&bodyBuffer, map[string]string{
		"ResetPasswordLink": resetLink,
	}); err != nil {
		return fmt.Errorf("failed to render email template: %w", err)
	}

	// Format full message
	message := []byte(fmt.Sprintf("From: %s <%s>\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		senderName, from, to, subject, bodyBuffer.String()))

	auth := smtp.PlainAuth("", username, password, smtpHost)
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	return smtp.SendMail(addr, auth, from, []string{to}, message)

}
