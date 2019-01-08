package core

import (
	"database/sql"

	"github.com/jinzhu/gorm"
	mocket "github.com/selvatico/go-mocket"
)

// GetTestConnection for get face connection
func GetTestConnection() *gorm.DB {

	mocket.Catcher.Register() // Safe register. Allowed multiple calls to save
	mocket.Catcher.Logging = true

	db, _ := gorm.Open(mocket.DriverName, "connection_string") // Can be any connection string

	return db

}

// GetTestConnectionSQL funtion sql test
func GetTestConnectionSQL() *sql.DB {

	mocket.Catcher.Register() // Safe register. Allowed multiple calls to save
	mocket.Catcher.Logging = true

	db, _ := sql.Open(mocket.DriverName, "connection_string")

	return db
}
