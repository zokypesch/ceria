package main

import (
	"bytes"
	"log"
	"time"

	"github.com/zokypesch/ceria/core"
	"github.com/zokypesch/ceria/helper"
)

func main() {
	config := helper.NewReadConfigService()
	config.Init()

	host := config.GetByName("rabbitmq.host")
	hostname := config.GetByName("rabbitmq.hostname")
	port := config.GetByName("rabbitmq.port")
	user := config.GetByName("rabbitmq.user")
	password := config.GetByName("rabbitmq.password")

	rb, errNew := core.NewServiceRabbitMQ(&core.RabbitMQConfig{
		Host:       host,
		Hostname:   hostname,
		Port:       port,
		User:       user,
		Password:   password,
		WorkerName: "my_task",
	})

	if errNew != nil {
		panic(errNew)
	}

	msgs, err := rb.RegisterWorker()

	if err != nil {
		panic(err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			dotCount := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)
			log.Printf("Done")
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}
