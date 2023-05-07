package telegram

import (
	"sevchel74botService/models"
	"sevchel74botService/src/db"
)

type Bot struct {
	Id          *int   `json:"id,omitempty" gorm:"column:id"`
	Name        string `json:"name" gorm:"column:name"`
	Token       string `json:"token" gorm:"column:token"`
	Description string `json:"description" gorm:"column:description"`
	Timeout     int    `json:"timeout" gorm:"column:timeout"`
	ParseMode   string `json:"parse_mode" gorm:"column:parse_mode"`
	Welcome     bool   `json:"welcome" gorm:"column:welcome"`
	WelcomeText string `json:"welcome_text" gorm:"column:welcome_text"`
	models.Model
}

type Bots []*Bot

const (
	bot = "sevchel.bot"
)

func GetBot(query string, args ...interface{}) (s *Bot, err error) {
	err = db.DB(bot).GetRecord(&s, query, args...)
	return
}

func GetBots(query string, args ...interface{}) (s Bots, err error) {
	err = db.DB(bot).GetRecords(&s, query, args...)
	return
}
