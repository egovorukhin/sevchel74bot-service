package telegram

import (
	"sevchel74botService/models/telegram"
	"strings"
)

var bots = make(map[string]*Bot)

func InitTelegramBots() error {

	b, err := telegram.GetBots("enabled=?", true)
	if err != nil {
		return err
	}
	for _, bb := range b {

		tgBot := &Bot{
			Name:        bb.Name,
			Token:       bb.Token,
			Timeout:     bb.Timeout,
			ParseMode:   bb.ParseMode,
			Welcome:     bb.Welcome,
			WelcomeText: bb.WelcomeText,
		}
		m, _ := telegram.GetViewBotModerators("bot_id=?", *bb.Id)
		if m != nil {
			tgBot.Moderator = &Moderator{
				Words: make(map[string]interface{}),
			}
			for _, mm := range m {
				tgBot.Moderator.Pattern = mm.Pattern
				tgBot.Moderator.Warn = mm.Warn
				tgBot.Moderator.WarnNumber = mm.WarnNumber
				for _, word := range strings.Split(mm.Words, " ") {
					tgBot.Moderator.Words[word] = nil
				}
			}
		}
		bots[bb.Name] = tgBot
		go tgBot.Start()
	}

	return nil
}

func StopTelegramBots() {
	for _, value := range bots {
		value.Stop()
	}
}
