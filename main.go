package main

import (
	"fmt"
	info "github.com/egovorukhin/egoappinfo"
	"github.com/egovorukhin/egoconf"
	"github.com/egovorukhin/egolog"
	"os"
	"os/signal"
	logger "sevchel74botService/src/logger"
	"sevchel74botService/src/server"
	"sevchel74botService/src/telegram"
	"syscall"
)

type Config struct {
	Server   server.Server   `yaml:"server"`
	Logger   logger.Config   `yaml:"logger"`
	Telegram telegram.Config `yaml:"telegram"`
}

type Error struct {
	Name  string
	Error error
}

func init() {
	info.SetName("Sevchel74Bot Service")
	info.SetVersion(0, 0, 1)
}

func main() {

	// Канал для получения ошибки, если таковая будет
	errChan := make(chan Error, 2)
	go start(errChan)
	// Ждем сигнал от ОС
	go waitSignal(errChan)

	err := <-errChan
	egolog.Errorcb("%s: %s", err.Name, err.Error)
	egolog.Infocb("Остановка приложения")
}

func start(errChan chan Error) {

	// Загружаем конфигурацию приложения
	cfg := Config{}
	err := egoconf.Load("config.yml", &cfg)
	if err != nil {
		errChan <- Error{
			Name:  "config",
			Error: err,
		}
		return
	}

	// Инициализация логгера
	err = logger.Init(cfg.Logger)
	if err != nil {
		errChan <- Error{
			Name:  "logger",
			Error: err,
		}
		return
	}

	egolog.Infocb("Стартуем приложение")

	// Инициализация telegram клиента
	telegram.Init(cfg.Telegram)

	// Инициализация WebSocketServer
	go func() {
		errChan <- Error{
			Name:  "webhook",
			Error: cfg.Server.Webhook.Init(),
		}
	}()

	// Инициализация WebServer
	errChan <- Error{
		Name:  "http",
		Error: cfg.Server.Http.Init(),
	}
}

func waitSignal(errChan chan Error) {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	errChan <- Error{
		Name:  "system",
		Error: fmt.Errorf("%s", <-c),
	}
}
