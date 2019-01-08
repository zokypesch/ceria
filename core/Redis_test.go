package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zokypesch/ceria/helper"
)

func TestRegisterRedis(t *testing.T) {
	config := helper.NewReadConfigService()
	config.Init()

	host := config.GetByName("redis.host")
	port := config.GetByName("redis.port")
	stat := config.GetByName("redis.status")

	if stat == "false" {
		return
	}

	t.Run("Failed Registration in Redis Core", func(t *testing.T) {
		_, err := NewServiceRedisCore(host, "5000")

		assert.Error(t, err)

	})

	t.Run("Success Registration in Redis Core", func(t *testing.T) {
		_, err := NewServiceRedisCore(host, port)

		assert.NoError(t, err)
	})

	t.Run("Reuse Registration in Redis Core", func(t *testing.T) {
		_, err := NewServiceRedisCore(host, port)

		assert.NoError(t, err)
	})

}

func TestCreateOrUpdate(t *testing.T) {
	config := helper.NewReadConfigService()
	config.Init()

	host := config.GetByName("redis.host")
	port := config.GetByName("redis.port")
	stat := config.GetByName("redis.status")

	if stat == "false" {
		return
	}

	cmd, err := NewServiceRedisCore(host, port)

	assert.NoError(t, err)

	t.Run("Failed Create Document in Redis Core", func(t *testing.T) {
		err = cmd.CreateOrUpdateDocument("", "", "", "")

		assert.Error(t, err)
	})

	t.Run("Create Document in Redis Core", func(t *testing.T) {
		err = cmd.CreateOrUpdateDocument("album", "1", "title", "welcome to the jungle")

		assert.NoError(t, err)
	})

	t.Run("Update Document in Redis Core", func(t *testing.T) {
		err = cmd.CreateOrUpdateDocument("album", "1", "title", "welcome to the city")

		assert.NoError(t, err)
	})

}

func TestDeleteDocument(t *testing.T) {
	config := helper.NewReadConfigService()
	config.Init()

	host := config.GetByName("redis.host")
	port := config.GetByName("redis.port")
	stat := config.GetByName("redis.status")

	if stat == "false" {
		return
	}
	cmd, err := NewServiceRedisCore(host, port)

	assert.NoError(t, err)

	t.Run("Failed Delete Document in Redis Core", func(t *testing.T) {
		err = cmd.DeleteDocument("", "")
		assert.Error(t, err)
	})

	t.Run("Success Delete in Redis Core", func(t *testing.T) {
		err = cmd.DeleteDocument("album", "1")

		assert.NoError(t, err)
	})

}

func TestGetDocument(t *testing.T) {
	config := helper.NewReadConfigService()
	config.Init()

	host := config.GetByName("redis.host")
	port := config.GetByName("redis.port")
	stat := config.GetByName("redis.status")

	if stat == "false" {
		return
	}

	var res string

	cmd, err := NewServiceRedisCore(host, port)
	assert.NoError(t, err)

	err = cmd.CreateOrUpdateDocument("album", "1", "title", "welcome to the city")
	assert.NoError(t, err)

	t.Run("Failed Get Document in Redis Core", func(t *testing.T) {
		res, err = cmd.GetDocument("", "", "")
		assert.Error(t, err)
		assert.Equal(t, "", res)
	})

	t.Run("Success Get Document in Redis Core", func(t *testing.T) {
		res, err = cmd.GetDocument("album", "1", "title")
		assert.NoError(t, err)

		assert.Equal(t, "welcome to the city", res)
	})

}

func TestGetAllDocument(t *testing.T) {

	config := helper.NewReadConfigService()
	config.Init()

	host := config.GetByName("redis.host")
	port := config.GetByName("redis.port")
	stat := config.GetByName("redis.status")

	if stat == "false" {
		return
	}

	var res map[string]string

	cmd, err := NewServiceRedisCore(host, port)
	defer cmd.Conn.Close()
	assert.NoError(t, err)

	t.Run("Failed GetAll Document in Redis Core", func(t *testing.T) {
		res, err = cmd.GetAllDocument("album", "")
		assert.NoError(t, err)

	})

	t.Run("Success GetAll Document in Redis Core", func(t *testing.T) {
		res, err = cmd.GetAllDocument("album", "1")
		assert.NoError(t, err)

		assert.Len(t, res, 1)
	})
}
