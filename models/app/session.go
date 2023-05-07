package app

import (
	"fmt"
	"gorm.io/gorm"
	"sevchel74botService/src/db"
	"time"
)

type Session struct {
	Id         string     `json:"id" gorm:"column:id"`
	UserId     *int       `json:"user_id" gorm:"column:user_id"`
	Ipaddress  string     `json:"ipaddress" gorm:"column:ipaddress"`
	Created    *time.Time `json:"created,omitempty" gorm:"column:created"`
	LastTime   *time.Time `json:"last_time,omitempty" gorm:"column:last_time"`
	ExpireTime int        `json:"expire_time" gorm:"column:expire_time"`
	UserAgent  string     `json:"user_agent" gorm:"column:user_agent"`
}

type Sessions []*Session

type ViewSession struct {
	SessionId   string     `json:"session_id" gorm:"column:session_id"`
	UserId      *int       `json:"user_id" gorm:"column:user_id"`
	Ipaddress   string     `json:"ipaddress" gorm:"column:ipaddress"`
	Created     *time.Time `json:"created,omitempty" gorm:"column:created"`
	LastTime    *time.Time `json:"last_time,omitempty" gorm:"column:last_time"`
	ExpireTime  int        `json:"expire_time" gorm:"column:expire_time"`
	UserAgent   string     `json:"user_agent" gorm:"column:user_agent"`
	Login       string     `json:"login" gorm:"column:login"`
	Password    string     `json:"password" gorm:"column:password"`
	Firstname   string     `json:"firstname" gorm:"column:firstname"`
	Lastname    string     `json:"lastname" gorm:"column:lastname"`
	Middlename  string     `json:"middlename" gorm:"column:middlename"`
	Description string     `json:"description" gorm:"column:description"`
	Local       bool       `json:"local" gorm:"column:local"`
}

type ViewSessions []*ViewSession

const (
	session   = "app.session"
	vwSession = "app.vw_session"
)

func GetSession(where string) (s *Session, err error) {
	err = db.DB(session).GetRecord(where, &s)
	return
}

func GetSessions(where string) (s Sessions, err error) {
	err = db.DB(session).GetRecords(where, &s)
	return
}

func GetViewSession(where string) (s *ViewSession, err error) {
	err = db.DB(vwSession).GetRecord(where, &s)
	return
}

func GetViewSessions(where string) (s ViewSessions, err error) {
	err = db.DB(vwSession).GetRecords(where, &s)
	return
}

// Set Установка сессии
func (s Session) Set() error {
	sTemp, err := GetSession(fmt.Sprintf("user_id=%d and ipaddress='%s'", *s.UserId, s.Ipaddress))
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	err = DeleteSession(sTemp.Id)
	if err != nil {
		return err
	}
	now := time.Now()
	s.LastTime = &now
	return db.DB(session).Create(&s)
}

// Update Обновление сессии
func (s Session) Update() error {
	return db.DB(session).Update(&s, fmt.Sprintf("id='%s'", s.Id))
}

// DeleteSession Удаление сессии
func DeleteSession(id string) error {
	return db.DB(session).Delete(nil, fmt.Sprintf("id='%s'", id))
}
