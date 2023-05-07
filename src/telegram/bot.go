package telegram

import (
	"fmt"
	info "github.com/egovorukhin/egoappinfo"
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
	Pattern        string                 `json:"pattern"`
	Words          map[string]interface{} `json:"words"`
	Warn           bool                   `json:"warn"`
	WarnNumber     int                    `json:"warn_number"`
	UntilDate      int64                  `json:"until_date"`
	RevokeMessages bool                   `json:"revoke_messages"`
}

// Start Запуск бота
func (b *Bot) Start() (err error) {

	egolog.Infofn(b.Name, "start updates")
	defer egolog.Infofn(b.Name, "stop updates")

	b.bot, err = tbApi.NewBotAPI(b.Token)
	if err != nil {
		return
	}

	updateConfig := tbApi.NewUpdate(1)
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
	case fmt.Sprintf("/start@%s", b.Name):
	case fmt.Sprintf("/about@%s", b.Name):
		if message.Chat != nil {
			app := info.GetApplication()
			return b.SendMessage(message.Chat.ID, fmt.Sprintf("%s v.%s", app.Name, app.Version.String()))
		}
	}

	// Проверка на нового пользователя в группе
	err = b.NewChatMembers(message.NewChatMembers, message.Chat)
	if err != nil {
		return err
	}

	// Проверка на покидание пользователя в группе
	err = b.LeftChatMember(message.LeftChatMember)
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

func (b *Bot) moderator(message *tbApi.Message) (err error) {

	var isWarning, ok bool
	if b.Moderator.Pattern != "" {
		if ok, err = regexp.MatchString(b.Moderator.Pattern, message.Text); err != nil {
			return err
		}
		if ok {
			err = b.DeleteMessage(message.Chat.ID, message.MessageID)
			if err != nil {
				return err
			}
			isWarning = true
		}
	}

	if len(b.Moderator.Words) > 0 {
		for _, word := range strings.Split(message.Text, " ") {
			if _, ok = b.Moderator.Words[strings.ToLower(word)]; ok {
				err = b.DeleteMessage(message.Chat.ID, message.MessageID)
				if err != nil {
					return err
				}
				isWarning = true
				break
			}
		}
	}

	if isWarning && b.Moderator.Warn && message.From != nil {
		u, err := telegram.SetUser(message.From.ID, message.From.UserName, message.From.FirstName, message.From.LastName, true)
		if err != nil {
			return err
		}
		if u.WarnCount >= b.Moderator.WarnNumber {
			return b.BanChatMember(u.UserId, message.Chat, message.Date)
		}
	}

	return nil
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

func (b *Bot) BanChatMember(userId int64, chat *tbApi.Chat, messageDate int) error {
	ban := tbApi.BanChatMemberConfig{
		ChatMemberConfig: tbApi.ChatMemberConfig{
			ChatID:             chat.ID,
			SuperGroupUsername: chat.UserName,
			ChannelUsername:    "",
			UserID:             userId,
		},
		UntilDate:      int64(messageDate) + b.Moderator.UntilDate,
		RevokeMessages: b.Moderator.RevokeMessages,
	}
	r, err := b.bot.Request(ban)
	if err != nil {
		return err
	}
	if string(r.Result) == "false" {
		return fmt.Errorf("Не удалось забанить участника[%s]", userId)
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
func (b *Bot) NewChatMembers(users []tbApi.User, chat *tbApi.Chat) error {
	for _, user := range users {
		if !user.IsBot {
			if b.Welcome {
				err := b.SendMessage(user.ID, b.WelcomeText)
				if err != nil {
					return err
				}
			}
			_, err := telegram.SetUser(user.ID, user.UserName, user.FirstName, user.LastName, false)
			if err != nil {
				return err
			}
		} else if user.UserName == b.Name && chat != nil {
			err := b.SetChat(chat)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// LeftChatMember Удаление участника из чата
func (b *Bot) LeftChatMember(user *tbApi.User) error {
	if user == nil {
		return nil
	}
	return telegram.RemoveUser(user.ID)
}

// SetChat Добавить чат
func (b *Bot) SetChat(chat *tbApi.Chat) error {
	c, err := telegram.GetChat("chat_id=?", chat.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c = &telegram.Chat{
				ChatId:      chat.ID,
				Title:       chat.Title,
				Type:        chat.Type,
				Description: chat.Description,
			}
			err = c.Insert()
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	return c.Update()
}
