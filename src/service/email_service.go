package service

import (
	"app/src/config"
	"app/src/utils"
	"fmt"

	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type EmailService interface {
	SendEmail(to, subject, body string) error
	SendResetPasswordEmail(to, token string) error
	SendVerificationEmail(to, token string) error
}

type emailService struct {
	Log    *logrus.Logger
	Dialer *gomail.Dialer
}

func NewEmailService() EmailService {
	return &emailService{
		Log: utils.Log,
		Dialer: gomail.NewDialer(
			config.SMTPHost,
			config.SMTPPort,
			config.SMTPUsername,
			config.SMTPPassword,
		),
	}
}

func (s *emailService) SendEmail(to, subject, body string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", config.EmailFrom)
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/plain", body)

	if err := s.Dialer.DialAndSend(mailer); err != nil {
		s.Log.Errorf("Failed to send email: %v", err)
		return err
	}

	return nil
}

func (s *emailService) SendResetPasswordEmail(to, token string) error {
	subject := "Reset password"

	// TODO: replace this url with the link to the reset password page of your front-end app
	resetPasswordURL := fmt.Sprintf("http://link-to-app/reset-password?token=%s", token)
	body := fmt.Sprintf(`Dear user,

To reset your password, click on this link: %s

If you did not request any password resets, then ignore this email.`, resetPasswordURL)
	return s.SendEmail(to, subject, body)
}

func (s *emailService) SendVerificationEmail(to, token string) error {
	subject := "Email Verification"

	// TODO: replace this url with the link to the email verification page of your front-end app
	verificationEmailURL := fmt.Sprintf("http://link-to-app/verify-email?token=%s", token)
	body := fmt.Sprintf(`Dear user,

To verify your email, click on this link: %s

If you did not create an account, then ignore this email.`, verificationEmailURL)
	return s.SendEmail(to, subject, body)
}
