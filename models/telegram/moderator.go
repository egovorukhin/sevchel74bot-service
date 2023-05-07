package telegram

import (
	"sevchel74botService/models"
	"sevchel74botService/src/db"
)

type Moderator struct {
	Id          *int   `json:"id,omitempty" gorm:"column:id"`
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"description" gorm:"column:description"`
	Pattern     string `json:"pattern" gorm:"column:pattern"`
	Delete      bool   `json:"delete" gorm:"column:delete"`
	Warn        bool   `json:"warn" gorm:"column:warn"`
	WarnNumber  int    `json:"warn_number" gorm:"column:warn_number"`
	Words       string `json:"words" gorm:"column:words"`
	models.Model
}

type Moderators []*Moderator

const (
	moderator = "sevchel.moderator"
)

func GetModerator(query string, args ...interface{}) (s *Moderator, err error) {
	err = db.DB(moderator).GetRecord(&s, query, args...)
	return
}

func GetModerators(query string, args ...interface{}) (s Moderators, err error) {
	err = db.DB(moderator).GetRecords(&s, query, args...)
	return
}
