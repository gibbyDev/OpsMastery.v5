package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func SendVerificationEmail(to string, token string) error {
	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	verificationLink := fmt.Sprintf("%s/api/v1/verify/%s", os.Getenv("APP_URL"), token)

	subject := "Verify Your Email"
	body := fmt.Sprintf(`
        <html>
            <body>
                <h2>Welcome to OpsMastery!</h2>
                <p>Please verify your email address by clicking the link below:</p>
                <a href="%s">Verify Email</a>
                <p>If you didn't create this account, please ignore this email.</p>
            </body>
        </html>
    `, verificationLink)

	message := fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", to, subject, body)

	auth := smtp.PlainAuth("", from, password, smtpHost)
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	err := smtp.SendMail(addr, auth, from, []string{to}, []byte(message))
	if err != nil {
		log.Printf("SMTP error: %v", err) // Add this line for more details
	}
	return err
}

func SendPasswordResetEmail(to string, token string) error {
	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	//resetLink := fmt.Sprintf("%s/auth/reset-password?reset_token=%s", os.Getenv("FRONTEND_URL"), token)

	resetLink := fmt.Sprintf("http://localhost:3000/auth/reset-password?reset_token=%s", token)

	log.Printf("Generated reset link: %s", resetLink)

	subject := "Reset Your Password"
	body := fmt.Sprintf(`
        <html>
            <body>
                <h2>Password Reset Request</h2>
                <p>Click the link below to reset your password:</p>
                <a href="%s">Reset Password</a>
                <p>If you didn't request this, please ignore this email.</p>
                <p>This link will expire in 1 hour.</p>
            </body>
        </html>
    `, resetLink)

	message := fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", to, subject, body)

	auth := smtp.PlainAuth("", from, password, smtpHost)
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	return smtp.SendMail(addr, auth, from, []string{to}, []byte(message))
}
