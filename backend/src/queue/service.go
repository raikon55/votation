package queue

import (
	"context"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type QueueManager interface {
	Enqueue(body string, queueName string)
	Close()
	estabillishConnection(queueName string)
}

type RabbitMQ struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      amqp.Queue
	qm         QueueManager
}

func (rmq *RabbitMQ) Enqueue(body string) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := rmq.channel.PublishWithContext(ctx,
		"",             // exchange
		rmq.queue.Name, // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
}

func (rmq *RabbitMQ) Consume() <-chan amqp.Delivery {
	msgs, err := rmq.channel.Consume(
		rmq.queue.Name, // queue
		"",             // consumer
		true,           // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	if err != nil {
		log.Println("Failed to register a consumer:", err)
	}
	return msgs
}

func (rmq *RabbitMQ) Close() {
	rmq.channel.Close()
	rmq.connection.Close()
}

func InitRabbitMQ(queueName string) (rmq RabbitMQ) {
	rmq.estabillishConnection(queueName)
	return
}

func (rmq *RabbitMQ) estabillishConnection(queueName string) {
	urlAmqp := os.Getenv("RABBITMQ_URL")

	conn, err := amqp.Dial(urlAmqp)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue:")

	rmq.queue = q
	rmq.channel = ch
	rmq.connection = conn
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
