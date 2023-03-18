package model

import (
	"github.com/tony/mot-server/cmd/server/config"
	"gorm.io/gorm"
)

var db = getDB()
var getDB = func() *gorm.DB {
	instance := config.InitMysql()
	AutoMigrate(instance)
	return instance
}

func AutoMigrate(db *gorm.DB) {
	var err error

	if err = db.AutoMigrate(&File{}); err != nil {
		panic(err)
	}
	if err = db.AutoMigrate(&FileStatus{}); err != nil {
		panic(err)
	}
	if err = db.AutoMigrate(&User{}); err != nil {
		panic(err)
	}

}
