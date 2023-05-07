package telegram

import "sevchel74botService/src/db"

type ViewBotModerator struct {
	BotId       int    `json:"bot_id,omitempty" gorm:"column:bot_id"`
	Bot         string `json:"bot" gorm:"column:bot"`
	Token       string `json:"token" gorm:"column:token"`
	Timeout     int    `json:"timeout" gorm:"column:timeout"`
	ParseMode   string `json:"parse_mode" gorm:"column:parse_mode"`
	Welcome     bool   `json:"welcome" gorm:"column:welcome"`
	WelcomeText string `json:"welcome_text" gorm:"column:welcome_text"`
	ModeratorId int    `json:"moderator_id,omitempty" gorm:"column:moderator_id"`
	Moderator   string `json:"moderator" gorm:"column:moderator"`
	Pattern     string `json:"pattern" gorm:"column:pattern"`
	Delete      bool   `json:"delete" gorm:"column:delete"`
	Warn        bool   `json:"warn" gorm:"column:warn"`
	WarnNumber  int    `json:"warn_number" gorm:"column:warn_number"`
	Words       string `json:"words" gorm:"column:words"`
}

type ViewBotModerators []*ViewBotModerator

const (
	vwBotModerator = "sevchel.vw_bot_moderator"
)

func GetViewBotModerator(query string, args ...interface{}) (s *ViewBotModerator, err error) {
	err = db.DB(vwBotModerator).GetRecord(&s, query, args...)
	return
}

func GetViewBotModerators(query string, args ...interface{}) (s ViewBotModerators, err error) {
	err = db.DB(vwBotModerator).GetRecords(&s, query, args...)
	return
}
