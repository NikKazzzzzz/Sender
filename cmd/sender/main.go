package main

import (
	"github.com/NikKazzzzzz/Sender/internal/config"
	"github.com/NikKazzzzzz/Sender/internal/notification"
	"github.com/NikKazzzzzz/Sender/internal/rabbitmq"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	"log"
)

func main() {
	cfg, err := config.LoadConfig("./config/sender.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	consumer, err := rabbitmq.NewConsumer(cfg.RabbitMQ.URL, cfg.RabbitMQ.Queue)
	if err != nil {
		log.Fatalf("failed to create consumer: %v", err)
	}
	defer consumer.Channel.Close()

	emailConfig := notification.EmailConfig{
		SMTPHost:       cfg.Email.SMTPHost,
		SMTPPort:       cfg.Email.SMTPPort,
		SenderEmail:    cfg.Email.SenderEmail,
		SenderPassword: cfg.Email.SenderPassword,
		RecipientEmail: cfg.Email.RecipientEmail,
		Subject:        cfg.Email.Subject,
	}

	notificationService := notification.NewNotificationService(emailConfig)

	err = consumer.Consume(notificationService.SendNotification)
	if err != nil {
		log.Fatalf("failed to consume messages: %v", err)
	}
}
