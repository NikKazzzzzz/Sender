package notification

import (
	"crypto/tls"
	"fmt"
	"github.com/NikKazzzzzz/Sender/monitoring"
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
	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s",
		service.emailConfig.SenderEmail,
		service.emailConfig.RecipientEmail,
		service.emailConfig.Subject,
		event,
	)

	addr := fmt.Sprintf("%s:%d", service.emailConfig.SMTPHost, service.emailConfig.SMTPPort)
	client, err := smtp.Dial(addr)
	if err != nil {
		log.Printf("Failed to connect to SMPT server: %v", err)
		return
	}
	defer client.Close()

	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         service.emailConfig.SMTPHost,
	}
	if err = client.StartTLS(tlsConfig); err != nil {
		log.Printf("Failed to start TLS: %v", err)
		return
	}

	auth := smtp.PlainAuth("", service.emailConfig.SenderEmail, service.emailConfig.SenderPassword, service.emailConfig.SMTPHost)
	if err = client.Auth(auth); err != nil {
		log.Printf("Failed to authenticate: %v", err)
		return
	}
	log.Printf("Authentication successful")

	if err = client.Mail(service.emailConfig.SenderEmail); err != nil {
		log.Printf("Failed to set sender: %v", err)
		return
	}

	if err = client.Rcpt(service.emailConfig.RecipientEmail); err != nil {
		log.Printf("Failed to set recipient: %v", err)
		return
	}

	writer, err := client.Data()
	if err != nil {
		log.Printf("Failed to get SMTP writer: %v", err)
		return
	}
	defer writer.Close()

	_, err = writer.Write([]byte(msg))
	if err != nil {
		log.Printf("Failed to write message: %v", err)
		return
	}

	log.Printf("Email sent successfully to %s", service.emailConfig.RecipientEmail)

	monitoring.SentMessagesCounter.Inc()
}
