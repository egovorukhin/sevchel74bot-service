package server

import (
	"fmt"
	info "github.com/egovorukhin/egoappinfo"
	"github.com/ewa-go/ewa"
	f "github.com/ewa-go/ewa-fiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"sevchel74botService/controllers"
	"time"
)

type Http struct {
	Root        string       `yaml:"root"`
	Certificate *Certificate `yaml:"cert,omitempty"`
	Port        int          `yaml:"port"`
	Timeout     *Timeout     `yaml:"timeout,omitempty"`
	Buffer      *Buffer      `yaml:"buffer,omitempty"`
	Secure      bool         `yaml:"secure"`
	Logger      *Logger      `yaml:"logger,omitempty"`
}

func (w Http) Init() error {

	wsCfg := ewa.Config{
		Port: w.Port,
		/*Authorization: security.Authorization{
			ApiKey: &security.ApiKey{
				KeyName: "Token",
				Param:   security.ParamHeader,
				Handler: auth.CheckToken,
			},
		},*/
		ContextHandler: w.contextHandler,
	}
	if w.Secure && w.Certificate != nil {
		wsCfg.Secure = &ewa.Secure{
			Key:        w.Certificate.Key,
			Cert:       w.Certificate.Cert,
			ClientCert: w.Certificate.ClientCert,
		}
	}
	fiberConfig := fiber.Config{}
	if w.Buffer != nil {
		fiberConfig.ReadBufferSize = w.Timeout.Read
		fiberConfig.WriteBufferSize = w.Timeout.Write
	}
	if w.Timeout != nil {
		fiberConfig.ReadTimeout = time.Duration(w.Timeout.Read) * time.Second
		fiberConfig.WriteTimeout = time.Duration(w.Timeout.Write) * time.Second
		fiberConfig.IdleTimeout = time.Duration(w.Timeout.Idle) * time.Second
	}

	app := fiber.New(fiberConfig)
	// CORS
	app.Use(cors.New())
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
	server := ewa.New(&f.Server{App: app}, wsCfg)

	//api
	server.Register(new(controllers.Api)).NotShow()

	// Описываем swagger
	server.Swagger.SetModels(ewa.Models{})
	server.Swagger.SetBasePath("/api").SetInfo(fmt.Sprintf("%s:%d", info.GetApplication().Hostname, w.Port), &ewa.Info{
		Description: "Бот 'Северный человек. Челябинск'",
		Version:     info.GetVersion().String(),
		Title:       info.GetApplicationName(),
		Contact: &ewa.Contact{
			Email: "yegor.govorukhin@mail.ru",
		},
		License: &ewa.License{
			Name: "Sevchel74Bot",
		},
	}, nil)

	//Запускаем сервер
	return server.Start()
}

func (Http) contextHandler(handler ewa.Handler) interface{} {
	return func(ctx *fiber.Ctx) error {
		return handler(ewa.NewContext(&f.Context{Ctx: ctx}))
	}
}
