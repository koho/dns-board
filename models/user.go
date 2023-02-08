package models

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/koho/dnstap-web/db"
)

type User struct {
	ID       int
	User     string `gorm:"uniqueIndex"`
	Password string
}

func GetUser(name string) (*User, error) {
	var user User
	if err := db.GetDB().Where("user = ?", name).Limit(1).Find(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetEncodedPassword(password string) string {
	encodedPass := ""
	if password != "" {
		b := md5.Sum([]byte(password))
		encodedPass = hex.EncodeToString(b[:])
	}
	return encodedPass
}

func CreateUser(name string, password string) error {
	user := &User{
		User:     name,
		Password: GetEncodedPassword(password),
	}
	return db.GetDB().Create(user).Error
}

func UpdateUserPassword(name string, password string) error {
	return db.GetDB().Model(&User{}).Where("user = ?", name).Update("password", GetEncodedPassword(password)).Error
}

func ensureAdmin() error {
	if _, err := GetUser("admin"); err != nil {
		return CreateUser("admin", "")
	}
	return nil
}
