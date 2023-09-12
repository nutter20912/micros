package database

import "gorm.io/gorm"

type Database interface {
	migrate()
}

var db *gorm.DB

func init() {

}

func Get() *gorm.DB {
	return db
}
