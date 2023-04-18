package telegram

import (
	"fmt"
	info "github.com/egovorukhin/egoappinfo"
)

const (
	AboutCommand = "/about"
)

func About(chatId int) error {
	app := info.GetApplication()
	m := &Message{
		ChatId: -900974177,
		Text:   fmt.Sprintf("%s v%s", app.Name, app.Version.String()),
	}
	return m.Send()
}
