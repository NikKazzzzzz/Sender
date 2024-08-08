package notification

import (
	"fmt"
	"log"
	"net/smtp"
)

type EmailConfig struct {
	SMTPHost       string
	SMTPPort       int
	SenderEmail    string
	SenderPassword string
	RecipientEmail string
	Subject        string
}

type NotificationService struct {
	emailConfig EmailConfig
}

func NewNotificationService(emailConfig EmailConfig) *NotificationService {
	return &NotificationService{emailConfig: emailConfig}
}

func (service *NotificationService) SendNotification(event string) {
	msg := fmt.Sprintf("Subject: %s\n\n%s", service.emailConfig.Subject, event)
	auth := smtp.PlainAuth("", service.emailConfig.SenderEmail, service.emailConfig.SenderPassword, service.emailConfig.SMTPHost)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", service.emailConfig.SMTPHost, service.emailConfig.SMTPPort),
		auth,
		service.emailConfig.SenderEmail,
		[]string{service.emailConfig.RecipientEmail},
		[]byte(msg),
	)

	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return
	}
	log.Printf("Sending notification for event: %s", event)
}
