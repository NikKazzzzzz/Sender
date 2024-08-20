package rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
)

type Consumer struct {
	Channel *amqp.Channel
	Queue   string
}

func NewConsumer(url string, queueName string) (*Consumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &Consumer{Channel: ch, Queue: queueName}, nil
}

func (c *Consumer) Consume(handler func(string)) error {
	msgs, err := c.Channel.Consume(
		c.Queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			handler(string(d.Body))

			if err := d.Ack(false); err != nil {
				log.Printf("Failed to acknowledge message: %v", err)
			}
		}
	}()

	log.Printf("Waiting for messages. To exit press CTRL+C")
	<-forever

	return nil
}
