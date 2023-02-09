package models

import "github.com/koho/dns-board/db"

func init() {
	db.RegisterModel(&Meta{}, &Message{}, &User{})
	db.OnStartup(ensureAdmin)
	db.OnStartup(func(option db.Option) error {
		if option.Retention != 0 {
			go dropOutdatedMessages(option.Retention)
		}
		return nil
	})
}

type Meta struct {
	Key   string `gorm:"primaryKey"`
	Value string
}

func GetMeta(key string) (string, error) {
	var meta Meta
	if err := db.GetDB().Where("key = ?", key).First(&meta).Error; err != nil {
		return "", err
	}
	return meta.Value, nil
}

func SetMeta(key, value string) error {
	return db.GetDB().Save(&Meta{key, value}).Error
}
