package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zokypesch/ceria/helper"
)

func TestNewService(t *testing.T) {
	config := helper.NewReadConfigService()
	config.Init()

	host := config.GetByName("rabbitmq.host")
	hostname := config.GetByName("rabbitmq.hostname")
	port := config.GetByName("rabbitmq.port")
	user := config.GetByName("rabbitmq.user")
	password := config.GetByName("rabbitmq.password")
	stat := config.GetByName("rabbitmq.status")

	if stat == "false" {
		return
	}

	t.Run("Failed Registration RabbitMQ", func(t *testing.T) {
		_, err := NewServiceRabbitMQ(&RabbitMQConfig{})

		assert.Error(t, err)
	})

	t.Run("Failed Registration wrong host RabbitMQ", func(t *testing.T) {
		_, err := NewServiceRabbitMQ(&RabbitMQConfig{"192.168.0.1", hostname, port, user, password, ""})

		assert.Error(t, err)
	})

	t.Run("register success rabbitMQ", func(t *testing.T) {
		rb, err := NewServiceRabbitMQ(&RabbitMQConfig{host, hostname, port, user, password, "my_task"})

		if rb != nil {
			rb.Channel.Close()
			defer rb.Conn.Close()
		}

		assert.NoError(t, err)
	})

	t.Run("register Worker success", func(t *testing.T) {
		rb, err := NewServiceRabbitMQ(&RabbitMQConfig{host, hostname, port, user, password, "my_task"})

		if rb != nil {
			_, errWork := rb.RegisterWorker()

			assert.NoError(t, errWork)
			rb.Channel.Close()
			defer rb.Conn.Close()
		}

		assert.NoError(t, err)
	})

	t.Run("register new task failed", func(t *testing.T) {
		rb, err := NewServiceRabbitMQ(&RabbitMQConfig{host, hostname, port, user, password, "my_task"})

		if rb != nil {
			errWork := rb.RegisterNewTask("")

			assert.Error(t, errWork)
			rb.Channel.Close()
			defer rb.Conn.Close()
		}

		assert.NoError(t, err)
	})

	t.Run("register new task success", func(t *testing.T) {
		rb, err := NewServiceRabbitMQ(&RabbitMQConfig{host, hostname, port, user, password, "my_task"})

		if rb != nil {
			errWork := rb.RegisterNewTask("Success to send our data")

			assert.NoError(t, errWork)
			rb.Channel.Close()
			defer rb.Conn.Close()
		}

		assert.NoError(t, err)
	})

}
