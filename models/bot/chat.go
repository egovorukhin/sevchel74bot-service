package bot

import "sevchel74botService/models"

type Chat struct {
	Id          *int   `json:"id"`
	Name        string `json:"name"`
	ChatId      int    `json:"chat_id"`
	Description string `json:"description"`
	models.Model
}

type Chats []*Chat

func GetChat() (c *Chat, err error) {
	return nil, err
}
