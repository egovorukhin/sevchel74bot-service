package telegram

import (
	"fmt"
	"github.com/egovorukhin/egolog"
	tbApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"regexp"
	"sevchel74botService/models/telegram"
	"strings"
)

type Bot struct {
	Name        string     `json:"name" yaml:"name"`
	Token       string     `json:"token" yaml:"token"`
	Timeout     int        `json:"timeout" yaml:"timeout"`
	ParseMode   string     `json:"parse_mode" yaml:"parseMode"`
	Welcome     bool       `json:"welcome" yaml:"welcome"`
	WelcomeText string     `json:"welcome_text" yaml:"welcomeText"`
	Moderator   *Moderator `json:"moderator,omitempty"`
	bot         *tbApi.BotAPI
}

type Moderator struct {
	Pattern    string                 `json:"pattern"`
	Words      map[string]interface{} `json:"words"`
	Warn       bool                   `json:"warn"`
	WarnNumber int                    `json:"warn_number"`
}

// Start Запуск бота
func (b *Bot) Start() (err error) {

	egolog.Infofn(b.Name, "start updates")
	defer egolog.Infofn(b.Name, "stop updates")

	bot.bot, err = tbApi.NewBotAPI(b.Token)
	if err != nil {
		return
	}

	updateConfig := tbApi.NewUpdate(0)
	updateConfig.Timeout = b.Timeout

	for update := range b.bot.GetUpdatesChan(updateConfig) {
		if update.Message == nil {
			continue
		}
		if err = b.receive(update.Message); err != nil {
			egolog.Errorfn(b.Name, err)
		}
	}

	return nil
}

// Stop Остановка бота
func (b *Bot) Stop() {
	b.bot.StopReceivingUpdates()
}

func (b *Bot) receive(message *tbApi.Message) (err error) {

	// Проверка на команды
	switch message.Text {
	case "/start":
	}

	// Проверка на нового пользователя в группе
	err = b.NewChatMembers(message.NewChatMembers)
	if err != nil {
		return err
	}

	// Если модератор
	if b.Moderator != nil {
		err = b.moderator(message)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Bot) moderator(message *tbApi.Message) error {

	if b.Moderator.Pattern != "" {
		ok, err := regexp.MatchString(b.Moderator.Pattern, message.Text)
		if err != nil {
			return err
		}
		if ok {
			err = b.DeleteMessage(message.Chat.ID, message.MessageID)
			if err != nil {
				return err
			}
		}
	}

	if len(b.Moderator.Words) > 0 {
		for _, word := range strings.Split(message.Text, " ") {
			if _, ok := b.Moderator.Words[strings.ToLower(word)]; ok {
				err := b.DeleteMessage(message.Chat.ID, message.MessageID)
				if err != nil {
					return err
				}
				if b.Moderator.Warn && message.From != nil {
					err = b.SetUser(*message.From, true)
					if err != nil {
						return err
					}
				}
				break
			}
		}
	}

	return nil
}

func (b *Bot) SetUser(user tbApi.User, incrementWarnCount bool) error {
	u, err := telegram.GetUser("user_id=?", user.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			u = &telegram.User{
				Username:  user.UserName,
				Firstname: user.FirstName,
				Lastname:  user.LastName,
				UserId:    user.ID,
			}
			if incrementWarnCount {
				u.WarnCount = 1
			}
			err = u.Insert()
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	if incrementWarnCount {
		u.WarnCount++
	}
	return u.Update()
}

// SendMessage Отправить сообщение
func (b *Bot) SendMessage(chatId int64, text string) error {
	m, err := b.bot.Send(tbApi.NewMessage(chatId, text))
	if err != nil {
		return err
	}
	if m.Chat != nil {
		fmt.Printf("chat_id: %d, msg: %s", m.Chat.ID, m.Text)
	}
	return nil
}

// SendMessageToChannel Отправить сообщение пользователю
func (b *Bot) SendMessageToChannel(username, text string) error {
	m, err := b.bot.Send(tbApi.NewMessageToChannel(username, text))
	if err != nil {
		return err
	}
	if m.Chat != nil {
		fmt.Printf("chat_id: %d, msg: %s", m.Chat.ID, m.Text)
	}
	return nil
}

// DeleteMessage Удалить сообщение
func (b *Bot) DeleteMessage(chatId int64, msgId int) error {
	r, err := b.bot.Request(tbApi.NewDeleteMessage(chatId, msgId))
	if err != nil {
		return err
	}
	if string(r.Result) == "false" {
		return fmt.Errorf("Не удалось удалить сообщение")
	}
	return nil
}

// NewChatMembers Добавленный пользователь в чат
func (b *Bot) NewChatMembers(users []tbApi.User) error {
	for _, user := range users {
		if !user.IsBot {
			if b.Welcome {
				err := b.SendMessage(user.ID, b.WelcomeText)
				if err != nil {
					return err
				}
			}
			return b.SetUser(user, false)
		}
	}
	return nil
}
