package main

import (
	"github.com/NikKazzzzzz/Sender/internal/config"
	"github.com/NikKazzzzzz/Sender/internal/notification"
	"github.com/NikKazzzzzz/Sender/internal/rabbitmq"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
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

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Println("Starting Prometheus metrics server on: 2112")
		if err := http.ListenAndServe(":2112", nil); err != nil {
			log.Fatalf("failed to start Prometheus metrics server: %v", err)
		}
	}()
}
