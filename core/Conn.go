package core

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Connection struct
type Connection struct {
	driver   string
	host     string
	port     int
	user     string
	password string
	dbname   string
}

// ConncetionRepo for interfacing
type ConncetionRepo interface {
	GetConn() (*gorm.DB, error)
}

// GetConn for get connection
func (conn *Connection) GetConn() (*gorm.DB, error) {

	myConn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		conn.host, conn.port, conn.user, conn.dbname, conn.password)

	db, err := gorm.Open(conn.driver, myConn)

	return db, err
}

// NewServiceConnection return a struct connection
func NewServiceConnection(driver, host string, port int, user, password, dbname string) *Connection {
	return &Connection{
		driver:   driver,
		host:     host,
		port:     port,
		user:     user,
		password: password,
		dbname:   dbname,
	}
}
