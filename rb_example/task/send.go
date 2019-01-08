package main

import (
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

	rb.RegisterNewTask("Hello world")
}
