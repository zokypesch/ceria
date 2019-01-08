package core

import (
	"database/sql"
	"testing"

	"github.com/jinzhu/gorm"
	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/assert"
)

func TestFakeConnection(t *testing.T) {
	assert := assert.New(t)

	t.Run("check fake connection Gorm", func(t *testing.T) {

		faker, _ := gorm.Open(mocket.DriverName, "connection_string")

		assert.EqualValues(GetTestConnection().Dialect().GetName(), faker.Dialect().GetName())

	})

	// Test Sql Fake Connection
	t.Run("check fake connection Gorm", func(t *testing.T) {

		fakerSQL, _ := sql.Open(mocket.DriverName, "connection_string")

		assert.EqualValues(GetTestConnectionSQL().Driver(), fakerSQL.Driver())
	})

}
