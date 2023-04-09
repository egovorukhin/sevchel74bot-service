package telegram

import (
	"github.com/egovorukhin/egorest"
	"time"
)

type Config struct {
	Url     string `yaml:"url"`
	Timeout int    `yaml:"timeout"`
}

type Response struct {
	OK     bool        `json:"ok"`
	Result interface{} `json:"result"`
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

// GetMe Инфа обо мне
func GetMe() (Me, error) {
	req := egorest.NewRequest("/getMe")
	resp := Response{
		Result: Me{},
	}
	err := client.ExecuteGet(req, &resp)
	return resp.Result.(Me), err
}

func GetUpdates() (Response, error) {
	req := egorest.NewRequest("/getUpdates")
	resp := Response{
		Result: Me{},
	}
	err := client.ExecuteGet(req, &resp)
	return resp.Result.(Me), err
}
