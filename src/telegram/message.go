package telegram

import (
	"fmt"
)

type Message struct {
	ChatId                int         `json:"chat_id"`
	Text                  string      `json:"text"`
	ParseMode             string      `json:"parse_mode,omitempty"`
	DisableWebPagePreview bool        `json:"disable_web_page_preview"`
	DisableNotification   bool        `json:"disable_notification"`
	ReplyToMessageId      *int        `json:"reply_to_message_id,omitempty"`
	ReplyMarkup           interface{} `json:"reply_markup,omitempty"`
}

func (m *Message) Send() error {

	r, err := ExecutePost("/sendMessage", m)
	if err != nil {
		return err
	}

	if r.OK {
		return nil
	}

	return fmt.Errorf("statusCode: %d, content: %s", r.ErrorCode, r.Description)
}
