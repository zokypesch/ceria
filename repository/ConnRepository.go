package repository

import "github.com/jinzhu/gorm"

// ConnRepository interface
type ConnRepository interface {
	GetConn() (*gorm.DB, error)
}
