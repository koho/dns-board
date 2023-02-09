package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type Handler func(Option) error

var (
	db       *gorm.DB
	models   []interface{}
	handlers []Handler
)

type Option struct {
	Retention int `yaml:"retention"`
}

func Init(opt Option) {
	var err error
	db, err = gorm.Open(sqlite.Open("dns.db"), &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             0 * time.Millisecond,
			LogLevel:                  logger.Warn,
			Colorful:                  true,
			IgnoreRecordNotFoundError: true,
		}),
	})
	if err != nil {
		log.Fatal(err)
	}
	if err = AutoMigrate(models...); err != nil {
		log.Fatal(err)
	}
	for _, h := range handlers {
		if err = h(opt); err != nil {
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

func OnStartup(h Handler) {
	handlers = append(handlers, h)
}

func GetDB() *gorm.DB {
	return db
}
