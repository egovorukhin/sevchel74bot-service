package models

import "time"

type Model struct {
	Created  *time.Time `json:"created"`
	Modified *time.Time `json:"modified"`
	Enabled  bool       `json:"enabled"`
	Author   string     `json:"author"`
}
