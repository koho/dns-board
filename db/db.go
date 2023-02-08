package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var (
	db       *gorm.DB
	models   []interface{}
	handlers []func() error
)

func Init() {
	var err error
	db, err = gorm.Open(sqlite.Open("dns.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	if err = AutoMigrate(models...); err != nil {
		log.Fatal(err)
	}
	for _, h := range handlers {
		if err = h(); err != nil {
			log.Fatal(err)
		}
	}
}

func AutoMigrate(dst ...interface{}) error {
	return db.AutoMigrate(dst...)
}

func RegisterModel(obj ...interface{}) {
	models = append(models, obj...)
}

func OnStartup(h func() error) {
	handlers = append(handlers, h)
}

func GetDB() *gorm.DB {
	return db
}
