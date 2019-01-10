package core

import (
	"fmt"

	"github.com/zokypesch/ceria/util"

	"github.com/streadway/amqp"
)

// RabbitMQCore struct of rabbit MQ
type RabbitMQCore struct {
	Config  *RabbitMQConfig
	Conn    *amqp.Connection
	Channel *amqp.Channel
	Queue   amqp.Queue
}

// RAbbitMQCoreInter interfacing function rabbit MQ struct
type RAbbitMQCoreInter interface{}

// RabbitMQConfig struct of rabbit mq
type RabbitMQConfig struct {
	Host       string `validate:"required"`
	Hostname   string `validate:"required"`
	Port       string `validate:"required"`
	User       string `validate:"required"`
	Password   string `validate:"required"`
	WorkerName string `validate:"required"`
}

// NewServiceRabbitMQ for new service of rabbitMQ
func NewServiceRabbitMQ(config *RabbitMQConfig) (*RabbitMQCore, error) {

	utils := util.NewUtilService(config)
	err := utils.Validate()

	var erCh error

	if err != nil {
		return nil, err
	}

	fullAddress := fmt.Sprintf("%s://%s:%s@%s:%s/", config.Hostname, config.User, config.Password, config.Host, config.Port)
	conn, errConn := amqp.Dial(fullAddress)

	if errConn != nil {
		return nil, errConn
	}

	ch, errChannel := conn.Channel()
	erCh = errChannel // bind this

	// if errConn != nil {
	// 	return nil, errCh
	// }

	q, errQueue := ch.QueueDeclare(
		config.WorkerName, // name
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	erCh = errQueue

	if erCh != nil {
		return nil, errQueue
	}

	// if errQueue != nil {
	// 	return nil, errQueue
	// }

	return &RabbitMQCore{
		Config:  config,
		Conn:    conn,
		Channel: ch,
		Queue:   q,
	}, nil
}

// RegisterWorker function for declare queou
func (rb *RabbitMQCore) RegisterWorker() (<-chan amqp.Delivery, error) {

	err := rb.Channel.Qos(
		3,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	if err != nil {
		return nil, err
	}

	msgs, errConsum := rb.Channel.Consume(
		rb.Config.WorkerName, // queue
		"",                   // consumer
		false,                // auto-ack
		false,                // exclusive
		false,                // no-local
		false,                // no-wait
		nil,                  // args
	)

	return msgs, errConsum
}

// RegisterNewTask for register new task
func (rb *RabbitMQCore) RegisterNewTask(body string) error {

	if body == "" {
		return fmt.Errorf("Body cannot empty")
	}

	err := rb.Channel.Publish(
		"",                   // exchange
		rb.Config.WorkerName, // routing key
		false,                // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})

	return err
}

// https://github.com/rabbitmq/rabbitmq-tutorials/blob/master/go/new_task.go
