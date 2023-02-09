package models

import (
	dnstap "github.com/dnstap/golang-dnstap"
	"github.com/koho/dns-board/db"
	"gorm.io/gorm"
	"log"
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
	Size     int                 `json:"size"`
	Duration int64               `json:"duration"`
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

func dropOutdatedMessages(day int) {
	period := 24 * time.Hour * time.Duration(day)
	ticker := time.NewTicker(period)
	defer ticker.Stop()
	for range ticker.C {
		var last Message
		if err := db.GetDB().Order("id desc").Limit(1).Find(&last).Error; err == nil && last.Time != nil {
			checkpoint := last.Time.Add(-period)
			if err = db.GetDB().Where("time < ?", checkpoint).Delete(&Message{}).Error; err != nil {
				log.Printf("drop message error: %v\n", err)
			}
		}
	}
}
