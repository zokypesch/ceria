package core

import (
	"strconv"
	"testing"

	mocks "github.com/zokypesch/ceria/core/mocks"
	"github.com/zokypesch/ceria/helper"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestConnection(t *testing.T) {
	assert := assert.New(t)

	config := helper.NewReadConfigService()
	config.Init()

	host := config.GetByName("db.host")
	port := config.GetByName("db.port")
	driver := config.GetByName("db.driver")
	user := config.GetByName("db.user")
	password := config.GetByName("db.password")
	dbname := config.GetByName("db.name")
	stat := config.GetByName("db.status")
	newPort, _ := strconv.Atoi(port)

	if stat == "false" {
		return
	}

	conn := NewServiceConnection(
		driver,
		host,
		newPort,
		user,
		password,
		dbname,
	)

	var g *gorm.DB

	connMock := new(mocks.ConncetionRepo)

	connMock.On("GetConn").Return(g, nil)

	connMock.GetConn()

	_, err := conn.GetConn()

	assert.NoError(err)
	// assert.Equal(exp.Value, act.Value)

	connMock.AssertExpectations(t)
	connMock.AssertNumberOfCalls(t, "GetConn", 1)
	// Test negative condition

	connNegative := NewServiceConnection(
		"postgres_wrong",
		"localhost",
		5432,
		"local",
		"local",
		"local",
	)

	_, mustErr := connNegative.GetConn()
	assert.Error(mustErr)
}
