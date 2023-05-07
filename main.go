package main

import (
	"fmt"
	info "github.com/egovorukhin/egoappinfo"
	"github.com/egovorukhin/egoconf"
	"github.com/egovorukhin/egolog"
	"os"
	"os/signal"
	"sevchel74botService/src/db"
	"sevchel74botService/src/logger"
	"sevchel74botService/src/server"
	"sevchel74botService/src/telegram"
	"syscall"
)

type Config struct {
	Server   server.Config `yaml:"server"`
	Database db.Config     `yaml:"database"`
	Logger   logger.Config `yaml:"logger"`
}

/*type Error struct {
	Name  string
	Error error
}*/

func init() {
	info.SetName("SevChel74 Bot Service")
	info.SetVersion(0, 0, 1)
}

func main() {

	// Канал для получения ошибки, если таковая будет
	errChan := make(chan error, 2)
	go start(errChan)
	// Ждем сигнал от ОС
	go waitSignal(errChan)

	err := <-errChan
	egolog.Errorcb("%s", err)
	egolog.Infocb("Остановка приложения")
}

func start(errChan chan error) {

	// Загружаем конфигурацию приложения
	cfg := Config{}
	err := egoconf.Load("config.yml", &cfg)
	if err != nil {
		errChan <- err
		return
	}

	// Инициализация логгера
	err = logger.Init(cfg.Logger)
	if err != nil {
		errChan <- err
		return
	}

	egolog.Infocb("Стартуем приложение")

	// Инициализация бд
	err = db.Init(cfg.Database)
	if err != nil {
		errChan <- err
		return
	}
	defer db.Close()

	// Инициализация telegram клиента
	err = telegram.InitTelegramBots()
	if err != nil {
		errChan <- err
		return
	}
	defer telegram.StopTelegramBots()

	// Инициализация WebServer
	errChan <- server.Init(cfg.Server)
}

func waitSignal(errChan chan error) {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	errChan <- fmt.Errorf("%s", <-c)
}
