package database

import (
	"gorm.io/gorm"
)

var Db *gorm.DB

func Init(_db *gorm.DB) {
	Db = _db
}
