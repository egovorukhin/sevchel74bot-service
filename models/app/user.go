package app

import (
	model "sevchel74botService/models"
	"sevchel74botService/src/db"
)

type User struct {
	Id          *int   `json:"id" gorm:"column:id"`
	Login       string `json:"login" gorm:"column:login"`
	Password    string `json:"password" gorm:"column:password"`
	Firstname   string `json:"firstname" gorm:"column:firstname"`
	Lastname    string `json:"lastname" gorm:"column:lastname"`
	Middlename  string `json:"middlename" gorm:"column:middlename"`
	Description string `json:"description" gorm:"column:description"`
	Local       bool   `json:"local" gorm:"column:local"`
	model.Model
}

type Users []*User

const (
	user = "app.user"
)

func GetUser(where string) (s *User, err error) {
	err = db.DB(user).GetRecord(where, &s)
	return
}

func GetUsers(where string) (s Users, err error) {
	err = db.DB(user).GetRecords(where, &s)
	return
}
