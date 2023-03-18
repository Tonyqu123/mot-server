package model

import (
	"sync"

	"github.com/tony/mot-server/cmd/server/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db     *gorm.DB
	dbOnce sync.Once
)

func initMysqlOrDie() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.GetDsnFromEnv()), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

var InitMySQLOrDie = func() error {
	dbOnce.Do(func() {
		db = initMysqlOrDie()
		AutoMigrate(db)
	})
	return nil
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
