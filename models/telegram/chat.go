package telegram

import (
	"sevchel74botService/models"
	"sevchel74botService/src/db"
)

type Chat struct {
	Id          *int   `json:"id" gorm:"column:id"`
	ChatId      int64  `json:"chat_id" gorm:"column:chat_id"`
	Title       string `json:"title" gorm:"column:title"`
	Type        string `json:"type" gorm:"column:type"`
	Description string `json:"description" gorm:"column:description"`
	models.Model
}

type Chats []*Chat

const (
	chat = "sevchel.chat"
)

func GetChat(query string, args ...interface{}) (s *Chat, err error) {
	err = db.DB(chat).GetRecord(&s, query, args...)
	return
}

func GetChats(query string, args ...interface{}) (s Chats, err error) {
	err = db.DB(chat).GetRecords(&s, query, args...)
	return
}

func (u Chat) Insert() error {
	return db.DB(chat).Create(&u)
}

func (u Chat) Update() error {
	return db.DB(chat).Update(&u, "id=?", *u.Id)
}
