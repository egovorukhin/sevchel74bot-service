package telegram

import (
	"fmt"
	"github.com/egovorukhin/egorest"
	"time"
)

type Config struct {
	Url     string `yaml:"url"`
	Timeout int    `yaml:"timeout"`
}

type GetResponse struct {
	OK     bool        `json:"ok"`
	Result interface{} `json:"result"`
}

type PostResponse struct {
	OK          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

var client *egorest.Client

// Init Инициализация клиента
func Init(c Config) {
	client = egorest.NewClient(egorest.Config{
		BaseUrl: egorest.BaseUrl{
			Url: c.Url,
		},
		Timeout: time.Second * time.Duration(c.Timeout),
	})
}

func Watch() error {
	u, err := GetUpdates()
	if err != nil {
		return err
	}
	for _, update := range u {
		if update.Message.Text == AboutCommand {
			err = About(update.Message.Chat.Id)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
	return nil
}

func ExecuteGet(path string, v interface{}) error {
	req := egorest.NewRequest(path)
	resp := GetResponse{
		Result: v,
	}
	return client.ExecuteGet(req, &resp)
}

func ExecutePost(path string, data interface{}) (*PostResponse, error) {
	req := egorest.NewRequest(path).JSON(data)
	resp := &PostResponse{}
	err := client.ExecutePost(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
