package app

import (
	"github.com/egovorukhin/egolog"
	"sevchel74botService/src/db"
	"time"
)

type Audit struct {
	Table      string    `json:"table" gorm:"column:table"`
	Path       string    `json:"path" gorm:"column:path"`
	Author     string    `json:"author" gorm:"column:author"`
	Datetime   time.Time `json:"datetime" gorm:"column:datetime"`
	Action     string    `json:"action" gorm:"column:action"`
	StatusCode int       `json:"status_code" gorm:"column:status_code"`
	Result     string    `json:"result" gorm:"column:result"`
	Request    string    `json:"request" gorm:"column:request"`
}

type Audits []*Audit

type Action string

const (
	SettingsAudit = "app.audit"
	ModelAudit    = "ModelAudit"

	Read   Action = "read"
	Create Action = "create"
	Update Action = "update"
	Delete Action = "delete"
	Add    Action = "add"
	Remove Action = "remove"
)

func NewAudit(action Action, tableName, author, path string) *Audit {
	if author == "" {
		author = "unknown"
	}
	return &Audit{
		Action: action.String(),
		Author: author,
		Path:   path,
		Table:  tableName,
	}
}

func GetAudit(where string) (s *Audit, err error) {
	err = db.DB(SettingsAudit).GetRecord(where, &s)
	return
}

func GetAudits(where string) (s Audits, err error) {
	err = db.DB(SettingsAudit).GetRecords(where, &s)
	return
}

func (a *Audit) Insert() {
	a.Datetime = time.Now()
	err := db.DB(SettingsAudit).Create(a)
	if err != nil {
		egolog.Errorfn("audit", "table: %s, action: %s, author: %s, path: %s, err: %s", a.Table, a.Action, a.Author, a.Path, err)
	}
}

func (c Action) String() string {
	return string(c)
}
