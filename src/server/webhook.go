package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"time"
)

type Webhook struct {
	Certificate *Certificate `yaml:"cert,omitempty"`
	Port        int          `yaml:"port"`
	Timeout     *Timeout     `yaml:"timeout,omitempty"`
	Buffer      *Buffer      `yaml:"buffer,omitempty"`
	Secure      bool         `yaml:"secure"`
	Logger      *Logger      `yaml:"logger,omitempty"`
}

func (w Webhook) Init() error {

	cfg := fiber.Config{}
	if w.Buffer != nil {
		cfg.ReadBufferSize = w.Timeout.Read
		cfg.WriteBufferSize = w.Timeout.Write
	}
	if w.Timeout != nil {
		cfg.ReadTimeout = time.Duration(w.Timeout.Read) * time.Second
		cfg.WriteTimeout = time.Duration(w.Timeout.Write) * time.Second
		cfg.IdleTimeout = time.Duration(w.Timeout.Idle) * time.Second
	}

	app := fiber.New(cfg)
	// Logger
	if w.Logger != nil {
		app.Use(logger.New(logger.Config{
			Format:       w.Logger.Format,
			TimeFormat:   w.Logger.Time.Format,
			TimeZone:     w.Logger.Time.Zone,
			TimeInterval: time.Duration(w.Logger.Time.Interval) * time.Millisecond,
			Output:       w.Logger,
		}))
	}

	addr := fmt.Sprintf(":%d", w.Port)

	if w.Secure && w.Certificate != nil {
		return app.ListenMutualTLS(addr, w.Certificate.Cert, w.Certificate.Key, w.Certificate.ClientCert)
	}
	return app.Listen(addr)
}
