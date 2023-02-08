package models

import (
	dnstap "github.com/dnstap/golang-dnstap"
	"github.com/koho/dnstap-web/db"
	"gorm.io/gorm"
	"time"
)

type Message struct {
	ID       int                 `json:"id"`
	Time     *time.Time          `json:"time"`
	Type     dnstap.Message_Type `json:"type"`
	IP       string              `json:"ip"`
	Port     uint32              `json:"port"`
	Protocol string              `json:"protocol"`
	Domain   string              `json:"domain"`
	QType    string              `json:"qtype"`
	Answer   string              `json:"answer"`
	RCode    int                 `json:"rcode"`
	Duration int64               `json:"duration"`
}

func init() {
	db.RegisterModel(&Message{}, &User{})
	db.OnStartup(ensureAdmin)
}

func AddTimeClause(m *gorm.DB, duration string) *gorm.DB {
	if duration == "x" {
		return m
	}
	d, err := time.ParseDuration(duration + "h")
	if err != nil {
		return m
	}
	end := time.Now()
	start := end.Add(-d)
	return m.Where("time between ? and ?", start, end)
}

func GetMessagesByType(t dnstap.Message_Type, duration string) ([]Message, error) {
	m := db.GetDB().Where("type = ?", t)
	var msg []Message
	if err := AddTimeClause(m, duration).Find(&msg).Error; err != nil {
		return nil, err
	}
	return msg, nil
}

func GetMessagesByKeyword(kw string) ([]Message, error) {
	var msg []Message
	if err := db.GetDB().Where("domain like ? or answer like ?", "%"+kw+"%", "%"+kw+"%").Find(&msg).Error; err != nil {
		return nil, err
	}
	return msg, nil
}
