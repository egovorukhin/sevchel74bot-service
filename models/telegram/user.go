package telegram

import (
	model "sevchel74botService/models"
	"sevchel74botService/src/db"
)

type User struct {
	Id        *int   `json:"id" gorm:"column:id"`
	UserId    int64  `json:"user_id" gorm:"column:user_id"`
	Username  string `json:"username" gorm:"column:username"`
	Firstname string `json:"firstname" gorm:"column:firstname"`
	Lastname  string `json:"lastname" gorm:"column:lastname"`
	WarnCount int    `json:"warn_count" gorm:"column:warn_count"`
	model.Model
}

type Users []*User

const (
	user = "sevchel.user"
)

func GetUser(query string, args ...interface{}) (s *User, err error) {
	err = db.DB(user).GetRecord(&s, query, args...)
	return
}

func GetUsers(query string, args ...interface{}) (s Users, err error) {
	err = db.DB(user).GetRecords(&s, query, args...)
	return
}

func (u User) Insert() error {
	return db.DB(user).Create(&u)
}

func (u User) Update() error {
	return db.DB(user).Update(&u, "id=?", *u.Id)
}
